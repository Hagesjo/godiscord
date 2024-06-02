package godiscord

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hagesjo/webgockets"
)

func NewBot(token, prefix string) (*Bot, error) {
	restClient, err := newRestClient(http.DefaultClient, token)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate rest client: %w", err)
	}

	gatewayURL, err := restClient.GetGatewayURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway url: %w", err)
	}

	wsClient, err := webgockets.NewClient(gatewayURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}

	return &Bot{
		wsClient:   wsClient,
		restClient: restClient,

		token:          token,
		gatewayURL:     gatewayURL, // NOTE: discord doesn't want this cached too long.
		prefix:         prefix,
		textCommands:   make(map[string]TextCommandFunc),
		eventListeners: make(map[string]eventHandler),

		guilds:            make(map[string]Guild),
		unavailableGuilds: make(map[string]Guild),
		fetchersByGuild:   make(map[string]*Fetcher),
	}, nil
}

type TextCommandFunc func(*Fetcher, []string, Channel) error
type EventListenFunc func(*Fetcher, DispatchEvent) error

type Bot struct {
	wsClient   *webgockets.Client
	restClient *restClient
	prefix     string

	token             string
	gatewayURL        string
	heartbeatInterval int
	lastSequence      *uint64
	resumeGatewayURL  string
	sessionID         string

	textCommands   map[string]TextCommandFunc
	eventListeners map[string]eventHandler

	unavailableGuilds map[string]Guild
	guilds            map[string]Guild
	fetchersByGuild   map[string]*Fetcher
}

// RegisterTextCommand registers a text command.
// handler will be called when a text is received in either DM:s or in a channel.
// command must be a single word, only include alphanumeric and -_, and it should start with a letter.
func (b *Bot) RegisterTextCommand(command string, handler TextCommandFunc) error {
	command = strings.ToLower(strings.TrimPrefix(command, b.prefix))
	re := regexp.MustCompile(`[a-z][a-z0-9-_]`)
	if !re.MatchString(command) {
		return fmt.Errorf("invalid command name")
	}

	b.textCommands[command] = handler

	return nil
}

func (b *Bot) RegisterEventListener(handler any) error {
	eventHandler, err := eventHandlerFromInterface(handler)
	if err != nil {
		return fmt.Errorf("failed to register event: %w", err)
	}
	b.eventListeners[eventHandler.name()] = eventHandler

	return nil
}

func (b *Bot) ListGuilds() (guilds []Guild) {
	for _, guild := range b.guilds {
		guilds = append(guilds, guild)
	}

	return guilds
}

func (b *Bot) GetGuild(name string) (Guild, error) {
	guild, ok := b.guilds[name]
	if !ok {
		return guild, fmt.Errorf("guild not found")
	}

	return guild, nil
}

func (b *Bot) GetVoiceStates(guildID string) ([]VoiceState, error) {
	f, ok := b.fetchersByGuild[guildID]
	if !ok {
		return nil, fmt.Errorf("no such guild")
	}

	return f.GetVoiceStates(), nil
}

func (b *Bot) GetMembers(guildID string) ([]GuildMember, error) {
	f, ok := b.fetchersByGuild[guildID]
	if !ok {
		return nil, fmt.Errorf("no such guild")
	}

	return f.GetMembers(), nil
}

func (b *Bot) Run() error {
	canReconnect := false
	connectRetries := 0
	for {
		if err := b.wsClient.Connect(); err != nil {
			slog.Info("Trying to connect.", "attempt", connectRetries+1)

			if connectRetries > 10 {
				return fmt.Errorf("maximum connect attempts tried: %w", err)
			}

			// TODO: Make it less arbitrary with the sleep time and max attempts.
			time.Sleep(5 * time.Second)
			connectRetries += 1
			continue
		}

		connectRetries = 0

		var err error
		canReconnect, err = b.handler(canReconnect)
		if err != nil {
			return fmt.Errorf("handler failed: %w", err)
		}

		if canReconnect {
			b.wsClient, err = webgockets.NewClient(b.resumeGatewayURL)
			if err != nil {
				return fmt.Errorf("failed to create websocket client: %w", err)
			}
		}
	}
}

// handler is the main listening loop for the bot, blocking until a disconnect happens.
// The returned bool describes if a reconnect is possible or not.
func (b *Bot) handler(isResuming bool) (bool, error) {
	heartbeatCtx, cancelHeartbeatCtx := context.WithCancel(context.Background())
	defer cancelHeartbeatCtx()

	ctx, cancel := context.WithCancel(heartbeatCtx)
	defer cancel()

	heartbeatStarted := false
	var heartbeatErr error

	for {
		var closeErr *webgockets.ErrClose
		event, err := webgockets.ReadJSON[Event](ctx, b.wsClient)
		if errors.As(err, &closeErr) {
			return ResumeableCloseEvents[int(closeErr.Code)], nil
		} else if errors.Is(err, context.Canceled) {
			return true, nil
		} else if err != nil {
			return false, fmt.Errorf("failed to read event: %w", err)
		}

		if event.SequenceNumber != nil {
			b.lastSequence = event.SequenceNumber
		}

		switch event.OpCode {
		case OpCodeDispatch:
			slog.Info("Got dispatch.", "type", *event.Type, "event", event)
			if err := b.handleDispatch(*event); err != nil {
				return false, fmt.Errorf("failed to handle dispatch event: %w", err)
			}
		case OpCodeResume:
			// Do nothing for now
		case OpCodeReconnect:
			return true, nil
		case OpCodeInvalidSession:
			slog.Info("Got invalid session.", "event", event)
			if canResume, err := UnmarshalJSON[bool](*event.Data); err != nil {
				return false, fmt.Errorf("discord sent invalid 'invalid session': %w", err)
			} else {
				return canResume, nil
			}
		case OpCodeHello:
			hello, err := UnmarshalJSON[Hello](*event.Data)
			if err != nil {
				return false, fmt.Errorf("discord sent invalid hello '%s': %w", *event.Data, err)
			}

			slog.Info("Got hello event.", "event", fmt.Sprintf("%#v", event), "hello", hello)

			b.heartbeatInterval = hello.HeartbeatInterval

			if !isResuming {
				if err := b.identify(); err != nil {
					return false, fmt.Errorf("failed to identify: %w", err)
				}
			} else {
				if err := b.resume(); err != nil {
					return false, fmt.Errorf("failed to resume: %w", err)
				}
			}

			if heartbeatStarted {
				continue
			}

			go func(ctx context.Context) {
				heartbeatErr = b.heartbeater(ctx)
				if heartbeatErr != nil {
					cancelHeartbeatCtx()
				}
			}(heartbeatCtx)
			heartbeatStarted = true

		case OpCodeHeartbeatAck:
			slog.Info("Got heartbeat ack event.")
		default:
			slog.Info("Got unknown event.", "event", fmt.Sprintf("%#v", event))
		}
	}
}

// handleDispatch handles all dispatch events sent which the cache needs to act on.
// For the majority of the event, it will do nothing.
// It will log a warning if an unknown dispatch event is received.
func (b *Bot) handleDispatch(event Event) error {
	if event.Type == nil {
		return fmt.Errorf("discord sent dispatch without type set")
	}

	eventType := *event.Type

	var ev any

	switch eventType {
	case "READY":
		readyEvent, err := UnmarshalJSON[Ready](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal create guild json: %w", err)
		}

		for _, guild := range b.guilds {
			if guild.Unavailable {
				b.unavailableGuilds[guild.ID] = guild
			} else {
				b.guilds[guild.ID] = guild
			}
		}

		b.resumeGatewayURL = fmt.Sprintf("%s?v=%d&encoding=json", readyEvent.ResumeGatewayURL, apiVersion)
		b.sessionID = readyEvent.SessionID
	case "RESUMED":
		// TODO: Use this to ensure a resume worked.
	case "APPLICATION_COMMAND_PERMISSIONS_UPDATE":
		ev = MustUnmarshalJSON[ApplicationCommandPermissionsUpdate](*event.Data)
	case "AUTO_MODERATION_RULE_CREATE":
		ev = MustUnmarshalJSON[AutoModerationRuleCreate](*event.Data)
	case "AUTO_MODERATION_RULE_UPDATE":
		ev = MustUnmarshalJSON[AutoModerationRuleUpdate](*event.Data)
	case "AUTO_MODERATION_RULE_DELETE":
		ev = MustUnmarshalJSON[AutoModerationRuleDelete](*event.Data)
	case "CHANNEL_CREATE":
		channelCreate := MustUnmarshalJSON[ChannelUpdate](*event.Data)
		channel := channelCreate.Channel

		if channel.GuildID == nil {
			break
		}

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("CHANNEL_CREATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		f.channelsByID[channel.ID] = channel

		ev = channelCreate
	case "CHANNEL_UPDATE":
		channelUpdate := MustUnmarshalJSON[ChannelUpdate](*event.Data)
		channel := channelUpdate.Channel

		if channel.GuildID == nil {
			break
		}

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("CHANNEL_UPDATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		f.channelsByID[channel.ID] = channel

		ev = channelUpdate
	case "CHANNEL_DELETE":
		channelDelete := MustUnmarshalJSON[ChannelDelete](*event.Data)
		channel := channelDelete.Channel

		if channel.GuildID == nil {
			break
		}

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("CHANNEL_DELETE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		delete(f.channelsByID, channel.ID)

		ev = channelDelete
	case "ENTITLEMENT_CREATE":
		ev = MustUnmarshalJSON[EntitlementCreate](*event.Data)
	case "ENTITLEMENT_UPDATE":
		ev = MustUnmarshalJSON[EntitlementUpdate](*event.Data)
	case "ENTITLEMENT_DELETE":
		ev = MustUnmarshalJSON[EntitlementDelete](*event.Data)
	case "GUILD_CREATE":
		guildEvent := MustUnmarshalJSON[GuildCreate](*event.Data)

		if guildEvent.Unavailable != nil && *guildEvent.Unavailable {
			b.unavailableGuilds[guildEvent.ID] = guildEvent.Guild
		} else {
			b.guilds[guildEvent.ID] = guildEvent.Guild
		}

		b.fetchersByGuild[guildEvent.ID] = newFetcher(guildEvent, b.restClient)

		ev = guildEvent
	case "GUILD_UPDATE":
		guildUpdate := MustUnmarshalJSON[GuildUpdate](*event.Data)

		b.guilds[guildUpdate.Guild.ID] = guildUpdate.Guild

		ev = guildUpdate
	case "GUILD_DELETE":
		guildDelete := MustUnmarshalJSON[GuildDelete](*event.Data)

		delete(b.fetchersByGuild, guildDelete.Guild.ID)
		delete(b.guilds, guildDelete.Guild.ID)

		ev = guildDelete
	case "THREAD_CREATE":
		thread := MustUnmarshalJSON[ThreadCreate](*event.Data)

		if thread.Channel.GuildID == nil {
			break
		}

		channel := thread.Channel

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("THREAD_CREATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		f.threadsByID[channel.ID] = channel

		ev = thread
	case "THREAD_UPDATE":
		thread := MustUnmarshalJSON[ThreadUpdate](*event.Data)

		if thread.Channel.GuildID == nil {
			break
		}

		channel := thread.Channel

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("THREAD_UPDATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		f.threadsByID[channel.ID] = channel

		ev = thread
	case "THREAD_DELETE":
		thread := MustUnmarshalJSON[ThreadDelete](*event.Data)

		if thread.Channel.GuildID == nil {
			break
		}

		channel := thread.Channel

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("THREAD_DELETE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		delete(f.threadsByID, channel.ID)

		ev = thread
	case "THREAD_LIST_SYNC":
		listSyncEvent := MustUnmarshalJSON[ThreadListSync](*event.Data)

		f, ok := b.fetchersByGuild[listSyncEvent.GuildID]
		if !ok {
			slog.Warn("THREAD_LIST_SYNC received with a channel outside of a known guild", "guild", listSyncEvent.GuildID)
			break
		}

		for _, thread := range listSyncEvent.Threads {
			f.threadsByID[thread.ID] = thread
		}

		ev = listSyncEvent
	case "THREAD_MEMBERS_UPDATE":
		ev = MustUnmarshalJSON[ThreadMembersUpdate](*event.Data)
	case "THREAD_MEMBER_UPDATE":
		ev = MustUnmarshalJSON[ThreadMemberUpdate](*event.Data)
	case "GUILD_AUDIT_LOG_ENTRY_CREATE":
		// Do nothing
	case "GUILD_EMOJIS_UPDATE":
		emojisUpdate := MustUnmarshalJSON[GuildEmojisUpdate](*event.Data)

		guild, ok := b.guilds[emojisUpdate.GuildID]
		if !ok {
			slog.Warn("GUILD_EMOJIS_UPDATE sent guild_id outside of a known guild", "guild", emojisUpdate.GuildID)
		}

		guild.Emojis = emojisUpdate.Emojis
		b.guilds[guild.ID] = guild

		ev = emojisUpdate
	case "GUILD_STICKERS_UPDATE":
		stickersUpdate := MustUnmarshalJSON[GuildStickersUpdate](*event.Data)

		guild, ok := b.guilds[stickersUpdate.GuildID]
		if !ok {
			slog.Warn("GUILD_STICKERS_UPDATE sent guild_id outside of a known guild", "guild", stickersUpdate.GuildID)
		}

		guild.Stickers = stickersUpdate.Stickers
		b.guilds[guild.ID] = guild

		ev = stickersUpdate
	case "GUILD_MEMBER_ADD":
		memberAdd := MustUnmarshalJSON[GuildMemberAdd](*event.Data)

		f, ok := b.fetchersByGuild[memberAdd.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBER_ADD received with a channel outside of a known guild", "guild", memberAdd.GuildID)
			break
		}

		f.membersByID[memberAdd.User.ID] = memberAdd.GuildMember

		ev = memberAdd
	case "GUILD_MEMBER_UPDATE":
		memberUpdate := MustUnmarshalJSON[GuildMemberUpdate](*event.Data)

		f, ok := b.fetchersByGuild[memberUpdate.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBER_UPDATE received with a channel outside of a known guild", "guild", memberUpdate.GuildID)
			break
		}

		member, ok := f.membersByID[memberUpdate.User.ID]
		if !ok {
			slog.Warn("GUILD_MEMBER_UPDATE received with unknown user", "user_id", memberUpdate.User.ID)
			break
		}

		f.membersByID[memberUpdate.User.ID] = GuildMember{
			User:                       &memberUpdate.User,
			Nick:                       memberUpdate.Nick,
			Avatar:                     memberUpdate.Avatar,
			Roles:                      memberUpdate.Roles,
			JoinedAt:                   *memberUpdate.JoinedAt,
			PremiumSince:               memberUpdate.PremiumSince,
			Flags:                      member.Flags,
			Deaf:                       memberUpdate.Deaf != nil && *memberUpdate.Deaf,
			Mute:                       memberUpdate.Mute != nil && *memberUpdate.Mute,
			Pending:                    memberUpdate.Pending,
			Permissions:                member.Permissions,
			CommunicationDisabledUntil: memberUpdate.CommunicationDisabledUntil,
		}

		ev = memberUpdate
	case "GUILD_MEMBER_REMOVE":
		memberRemove := MustUnmarshalJSON[GuildMemberRemove](*event.Data)

		f, ok := b.fetchersByGuild[memberRemove.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBER_REMOVE received with a channel outside of a known guild", "guild", memberRemove.GuildID)
			break
		}

		delete(f.membersByID, memberRemove.User.ID)

		ev = memberRemove
	case "GUILD_MEMBERS_CHUNK":
		chunk := MustUnmarshalJSON[GuildMembersChunk](*event.Data)

		f, ok := b.fetchersByGuild[chunk.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBERS_CHUNK received with a channel outside of a known guild", "guild", chunk.GuildID)
			break
		}

		for _, member := range chunk.Members {
			if member.User == nil {
				continue
			}

			f.membersByID[member.User.ID] = member
		}

		ev = chunk
	case "GUILD_ROLE_CREATE":
		create := MustUnmarshalJSON[GuildRoleCreate](*event.Data)

		guild, ok := b.guilds[create.GuildID]
		if !ok {
			slog.Warn("GUILD_ROLE_CREATE sent guild_id outside of a known guild", "guild", create.GuildID)
		}

		guild.Roles = append(guild.Roles, create.Role)
		b.guilds[guild.ID] = guild

		ev = create
	case "GUILD_ROLE_UPDATE":
		update := MustUnmarshalJSON[GuildRoleCreate](*event.Data)

		guild, ok := b.guilds[update.GuildID]
		if !ok {
			slog.Warn("GUILD_ROLE_UPDATE sent guild_id outside of a known guild", "guild", update.GuildID)
		}

		for i, role := range guild.Roles {
			if role.ID != update.Role.ID {
				continue
			}

			guild.Roles[i] = update.Role
			break
		}

		guild.Roles = append(guild.Roles, update.Role)
		b.guilds[guild.ID] = guild

		ev = update
	case "GUILD_ROLE_DELETE":
		delete := MustUnmarshalJSON[GuildRoleCreate](*event.Data)

		guild, ok := b.guilds[delete.GuildID]
		if !ok {
			slog.Warn("GUILD_ROLE_DELETE sent guild_id outside of a known guild", "guild", delete.GuildID)
		}

		for i, role := range guild.Roles {
			if role.ID != delete.Role.ID {
				continue
			}

			guild.Roles = append(guild.Roles[:i], guild.Roles[i+1:]...)
			break
		}

		b.guilds[guild.ID] = guild

		ev = delete
	case "GUILD_SCHEDULED_EVENT_CREATE":
		ev = MustUnmarshalJSON[GuildScheduledEventCreate](*event.Data)
	case "GUILD_SCHEDULED_EVENT_UPDATE":
		ev = MustUnmarshalJSON[GuildScheduledEventUpdate](*event.Data)
	case "GUILD_SCHEDULED_EVENT_DELETE":
		ev = MustUnmarshalJSON[GuildScheduledEventDelete](*event.Data)
	case "GUILD_SCHEDULED_EVENT_USER_REMOVE_EVENT":
		ev = MustUnmarshalJSON[guildScheduledEventUserRemoveEventHandler](*event.Data)
	case "GUILD_SCHEDULED_EVENT_USER_ADD_EVENT":
		ev = MustUnmarshalJSON[guildScheduledEventUserAddEventHandler](*event.Data)
	case "INTEGRATION_CREATE":
		ev = MustUnmarshalJSON[IntegrationCreate](*event.Data)
	case "INTEGRATION_UPDATE":
		ev = MustUnmarshalJSON[IntegrationUpdate](*event.Data)
	case "INTEGRATION_DELETE":
		ev = MustUnmarshalJSON[IntegrationDelete](*event.Data)
	case "INVITE_CREATE":
		ev = MustUnmarshalJSON[InviteCreate](*event.Data)
	case "INVITE_DELETE":
		ev = MustUnmarshalJSON[InviteDelete](*event.Data)
	case "MESSAGE_CREATE":
		messageCreate := MustUnmarshalJSON[MessageCreate](*event.Data)

		if strings.HasPrefix(messageCreate.Content, b.prefix) {
			s := strings.Fields(messageCreate.Content[1:])
			command, ok := b.textCommands[s[0]]
			if ok {
				fetcher := b.fetchersByGuild[messageCreate.GuildID]
				channel := fetcher.channelsByID[messageCreate.ChannelID]

				if err := command(fetcher, s[1:], channel); err != nil {

				}
			}
		}

		ev = messageCreate
	case "MESSAGE_UPDATE":
		ev = MustUnmarshalJSON[MessageUpdate](*event.Data)
	case "MESSAGE_DELETE":
		ev = MustUnmarshalJSON[MessageDelete](*event.Data)
	case "MESSAGE_DELETE_BULK":
		ev = MustUnmarshalJSON[MessageDeleteBulk](*event.Data)
	case "MESSAGE_REACTION_ADD":
		ev = MustUnmarshalJSON[MessageReactionAdd](*event.Data)
	case "MESSAGE_REACTION_REMOVE":
		ev = MustUnmarshalJSON[MessageReactionRemove](*event.Data)
	case "MESSAGE_REACTION_REMOVE_ALL":
		ev = MustUnmarshalJSON[MessageReactionRemoveAll](*event.Data)
	case "MESSAGE_REACTION_REMOVE_EMOJI":
		ev = MustUnmarshalJSON[MessageReactionRemoveEmoji](*event.Data)
	case "PRESENCE_UPDATE":
		// Do nothing.
	case "STAGE_INSTANCE_CREATE":
		ev = MustUnmarshalJSON[StageInstanceCreate](*event.Data)
	case "STAGE_INSTANCE_UPDATE":
		ev = MustUnmarshalJSON[StageInstanceUpdate](*event.Data)
	case "STAGE_INSTANCE_DELETE":
		ev = MustUnmarshalJSON[StageInstanceDelete](*event.Data)
		// Do nothing.
	case "TYPING_START":
		ev = MustUnmarshalJSON[TypingStart](*event.Data)
		// Do nothing.
	case "USER_UPDATE":
		userUpdate := MustUnmarshalJSON[UserUpdate](*event.Data)

		for _, f := range b.fetchersByGuild {
			member, ok := f.membersByID[userUpdate.User.ID]
			if !ok {
				continue
			}

			member.User = &userUpdate.User
			f.membersByID[userUpdate.User.ID] = member
		}

		ev = userUpdate
	case "VOICE_STATE_UPDATE":
		voiceEvent := MustUnmarshalJSON[VoiceStateUpdate](*event.Data)
		if voiceEvent.GuildID == nil {
			break
		}

		fetcher, ok := b.fetchersByGuild[*voiceEvent.GuildID]
		if !ok {
			slog.Warn("VOICE_STATE_UPDATE received with a channel outside of a known guild", "guild", *voiceEvent.GuildID)
		}

		if voiceEvent.ChannelID == nil {
			delete(fetcher.voiceStatesByID, voiceEvent.UserID)
		} else {
			fetcher.voiceStatesByID[voiceEvent.UserID] = voiceEvent.VoiceState
		}

	case "VOICE_SERVER_UPDATE":
		ev = MustUnmarshalJSON[VoiceServerUpdate](*event.Data)
	case "WEBHOOKS_UPDATE":
		ev = MustUnmarshalJSON[WebhooksUpdate](*event.Data)
	default:
		slog.Warn("Unparsed dispatch event", "type", eventType)
	}

	switch e := ev.(type) {
	case guildEvent:
		if eventHandler, ok := b.eventListeners[eventType]; ok {
			fetcher, ok := b.fetchersByGuild[e.guild()]
			if !ok {
				return fmt.Errorf("no fetcher found for guild")
			}

			return eventHandler.run(fetcher, ev)
		}
	}

	return nil
}

func (b *Bot) identify() error {
	identifyPayload := Identify{
		OpCode: OpCodeIdentity,
		Payload: IdentifyPayload{
			Token: b.token,
			Intents: IntentGuilds | IntentGuildMembers | IntentGuildModeration | IntentGuildEmojisAndStickers |
				IntentGuildIntegrations | IntentGuildWebhooks | IntentGuildInvites | IntentGuildVoiceStates |
				IntentGuildPresences | IntentGuildMessages | IntentGuildMessageReactions | IntentGuildMessageTyping |
				IntentDirectMessages | IntentDirectMessageReactions | IntentDirectMessageTyping | IntentMessageContent |
				IntentGuildScheduledEvents | IntentAutoModerationConfiguration | IntentAutoModerationExecution,
			// Intents: IntentGuilds | IntentGuildMessages | IntentDirectMessages,
			Properties: Properties{
				OS:      "raspbian",
				Browser: "webgockets",
				Device:  "godiscord",
			},
			Presence: &Presence{
				Activities: []*Activity{
					{
						Name: "developing simulator or something",
						Type: ActivityGame,
						URL:  "https://github.com/Hagesjo/webgockets",
					},
				},
				Status: UserStatusOnline,
				AFK:    false,
			},
		},
	}

	if _, err := b.wsClient.WriteJSON(identifyPayload); err != nil {
		return fmt.Errorf("failed to identify: %w", err)
	}

	return nil
}

func (b *Bot) resume() error {
	resume := Resume{
		OpCode: OpCodeResume,
		Payload: ResumePayload{
			Token:     b.token,
			SessionID: b.sessionID,
		},
	}

	if b.lastSequence != nil {
		resume.Payload.LastSequence = *b.lastSequence
	}

	if _, err := b.wsClient.WriteJSON(resume); err != nil {
		return fmt.Errorf("failed to resume: %w", err)
	}

	return nil
}

func (b *Bot) heartbeater(ctx context.Context) error {
	if b.heartbeatInterval == 0 {
		return fmt.Errorf("heartbeatInterval must be > 0")
	}

	intervalWithJitter := func(interval int) time.Duration {
		return time.Duration(math.Round(rand.Float64()*float64(interval))) * time.Millisecond
	}

	ticker := time.NewTicker(intervalWithJitter(b.heartbeatInterval))

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if _, err := b.wsClient.WriteJSON(Heartbeat{
				OpCode: OpCodeHeartbeat,
			}); err != nil {
				return fmt.Errorf("failed to send heartbeat: %w", err)
			}

			slog.Info("Sent heartbeat")

			ticker.Reset(intervalWithJitter(b.heartbeatInterval))
		}
	}
}
