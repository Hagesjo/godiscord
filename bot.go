package discordgo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
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

		token:        token,
		gatewayURL:   gatewayURL, // NOTE: discord doesn't want this cached too long.
		prefix:       prefix,
		textCommands: make(map[string]func(args []string) error),

		guilds:            make(map[string]Guild),
		unavailableGuilds: make(map[string]Guild),
		fetchersByGuild:   make(map[string]*Fetcher),
	}, nil
}

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

	textCommands map[string]func(args []string) error

	unavailableGuilds map[string]Guild
	guilds            map[string]Guild
	fetchersByGuild   map[string]*Fetcher
}

func (b *Bot) RegisterTextCommand(command string, handler func(args []string) error) {
	b.textCommands[command] = handler
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
		// Do nothing.
	case "AUTO_MODERATION_RULE_CREATE", "AUTO_MODERATION_RULE_UPDATE", "AUTO_MODERATION_RULE_DELETE":
		// Do nothing.
	case "CHANNEL_CREATE", "CHANNEL_UPDATE", "CHANNEL_DELETE":
		channel, err := UnmarshalJSON[Channel](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal channel: %w", err)
		}

		if channel.GuildID == nil {
			break
		}

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("CHANNEL_CREATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		if eventType == "CHANNEL_DELETE" {
			delete(f.channelsByID, channel.ID)
		} else {
			f.channelsByID[channel.ID] = channel
		}
	case "ENTITLEMENT_CREATE", "ENTITLEMENT_UPDATE", "ENTITLEMENT_DELETE":
		// Do nothing.
	case "GUILD_CREATE":
		guildEvent, err := UnmarshalJSON[GuildCreate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal create guild json: %w", err)
		}

		if guildEvent.Unavailable != nil && *guildEvent.Unavailable {
			b.unavailableGuilds[guildEvent.ID] = guildEvent.Guild
		} else {
			b.guilds[guildEvent.ID] = guildEvent.Guild
		}

		b.fetchersByGuild[guildEvent.ID] = newFetcher(guildEvent)
	case "GUILD_UPDATE", "GUILD_DELETE":
		guild, err := UnmarshalJSON[Guild](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild json: %w", err)
		}

		if eventType == "GUILD_UPDATE" {
			b.guilds[guild.ID] = guild
		} else {
			delete(b.fetchersByGuild, guild.ID)
			delete(b.guilds, guild.ID)
		}
	case "THREAD_CREATE", "THREAD_UPDATE", "THREAD_DELETE":
		channel, err := UnmarshalJSON[Channel](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal channel: %w", err)
		}

		if channel.GuildID == nil {
			break
		}

		f, ok := b.fetchersByGuild[*channel.GuildID]
		if !ok {
			slog.Warn("THREAD_CREATE received with a channel outside of a known guild", "guild", *channel.GuildID)
			break
		}

		if eventType == "THREAD_DELETE" {
			delete(f.threadsByID, channel.ID)
		} else {
			f.threadsByID[channel.ID] = channel
		}
	case "THREAD_LIST_SYNC":
		listSyncEvent, err := UnmarshalJSON[ThreadListSync](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal thread list sync event: %w", err)
		}

		f, ok := b.fetchersByGuild[listSyncEvent.GuildID]
		if !ok {
			slog.Warn("THREAD_CREATE received with a channel outside of a known guild", "guild", listSyncEvent.GuildID)
			break
		}

		for _, thread := range listSyncEvent.Threads {
			f.threadsByID[thread.ID] = thread
		}

	case "THREAD_MEMBER_UPDATE":
		// Do nothing.
	case "GUILD_AUDIT_LOG_ENTRY_CREATE":
		// Do nothing.
	case "GUILD_EMOJIS_UPDATE":
		emojisUpdate, err := UnmarshalJSON[GuildEmojisUpdate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal thread list sync event: %w", err)
		}

		guild, ok := b.guilds[emojisUpdate.GuildID]
		if !ok {
			slog.Warn("GUILD_EMOJIS_UPDATE sent guild_id outside of a known guild", "guild", emojisUpdate.GuildID)
		}

		guild.Emojis = emojisUpdate.Emojis
		b.guilds[guild.ID] = guild
	case "GUILD_STICKERS_UPDATE":
		stickersUpdate, err := UnmarshalJSON[GuildStickersUpdate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal thread list sync event: %w", err)
		}

		guild, ok := b.guilds[stickersUpdate.GuildID]
		if !ok {
			slog.Warn("GUILD_stickerS_UPDATE sent guild_id outside of a known guild", "guild", stickersUpdate.GuildID)
		}

		guild.Stickers = stickersUpdate.Stickers
		b.guilds[guild.ID] = guild
	case "GUILD_MEMBER_ADD":
		memberAdd, err := UnmarshalJSON[GuildMemberAdd](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild member update event: %w", err)
		}

		f, ok := b.fetchersByGuild[memberAdd.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBER_ADD received with a channel outside of a known guild", "guild", memberAdd.GuildID)
			break
		}

		f.membersByID[memberAdd.User.ID] = memberAdd.GuildMember
	case "GUILD_MEMBER_UPDATE":
		memberUpdate, err := UnmarshalJSON[GuildMemberUpdate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild member update event: %w", err)
		}

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
	case "GUILD_MEMBER_REMOVE":
		memberRemove, err := UnmarshalJSON[GuildMemberRemove](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild member update event: %w", err)
		}

		f, ok := b.fetchersByGuild[memberRemove.GuildID]
		if !ok {
			slog.Warn("GUILD_MEMBER_REMOVE received with a channel outside of a known guild", "guild", memberRemove.GuildID)
			break
		}

		delete(f.membersByID, memberRemove.User.ID)
	case "GUILD_MEMBERS_CHUNK":
		chunk, err := UnmarshalJSON[GuildMembersChunk](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild members chunk event: %w", err)
		}

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
	case "GUILD_ROLE_CREATE":
		create, err := UnmarshalJSON[GuildRoleCreate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild role create: %w", err)
		}

		guild, ok := b.guilds[create.GuildID]
		if !ok {
			slog.Warn("GUILD_ROLE_CREATE sent guild_id outside of a known guild", "guild", create.GuildID)
		}

		guild.Roles = append(guild.Roles, create.Role)
		b.guilds[guild.ID] = guild
	case "GUILD_ROLE_UPDATE":
		update, err := UnmarshalJSON[GuildRoleCreate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild role update: %w", err)
		}

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
	case "GUILD_ROLE_DELETE":
		delete, err := UnmarshalJSON[GuildRoleCreate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild role update: %w", err)
		}

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
	case "GUILD_SCHEDULED_EVENT_CREATE", "GUILD_SCHEDULED_EVENT_UPDATE", "GUILD_SCHEDULED_EVENT_DELETE", "GUILD_SCHEDULED_EVENT_USER_REMOVE_EVENT", "GUILD_SCHEDULED_EVENT_USER_ADD_EVENT":
		// Do nothing.
	case "INTEGRATION_CREATE", "INTEGRATION_UPDATE", "INTEGRATION_DELETE":
		// Do nothing.
	case "INVITE_CREATE", "INVITE_DELETE":
		// Do nothing.
	case "MESSAGE_CREATE", "MESSAGE_UPDATE", "MESSAGE_DELETE", "MESSAGE_DELETE_BULK":
		// Do nothing.
	case "MESSAGE_REACTION_ADD", "MESSAGE_REACTION_REMOVE", "MESSAGE_REACTION_REMOVE_ALL", "MESSAGE_REACTION_REMOVE_EMOJI":
		// Do nothing.
	case "PRESENCE_UPDATE":
		// Do nothing.
	case "STAGE_INSTANCE_CREATE", "STAGE_INSTANCE_UPDATE", "STAGE_INSTANCE_DELETE":
		// Do nothing.
	case "TYPING_START":
		// Do nothing.
	case "USER_UPDATE":
		userUpdate, err := UnmarshalJSON[UserUpdate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal guild member update event: %w", err)
		}

		for _, f := range b.fetchersByGuild {
			member, ok := f.membersByID[userUpdate.User.ID]
			if !ok {
				continue
			}

			member.User = &userUpdate.User
			f.membersByID[userUpdate.User.ID] = member
		}

	case "VOICE_STATE_UPDATE", "VOICE_SERVER_UPDATE":
		// Do nothing.
	case "WEBHOOKS_UPDATE":
		// Do nothing.
	default:
		slog.Warn("Unparsed dispatch event", "type", eventType)
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
				Device:  "discordgo",
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
