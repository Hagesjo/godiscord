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

	"github.com/hagesjo/discordgo/events"
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

		guilds:            make(map[string]events.Guild),
		unavailableGuilds: make(map[string]events.Guild),
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

	unavailableGuilds map[string]events.Guild
	guilds            map[string]events.Guild
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
		event, err := webgockets.ReadJSON[events.Event](ctx, b.wsClient)
		if errors.As(err, &closeErr) {
			return events.ResumeableCloseEvents[int(closeErr.Code)], nil
		} else if errors.Is(err, context.Canceled) {
			return true, nil
		} else if err != nil {
			return false, fmt.Errorf("failed to read event: %w", err)
		}

		if event.SequenceNumber != nil {
			b.lastSequence = event.SequenceNumber
		}

		switch event.OpCode {
		case events.OpCodeDispatch:
			slog.Info("Got dispatch.", "type", *event.Type, "event", event)
			if err := b.handleDispatch(*event); err != nil {
				return false, fmt.Errorf("failed to handle dispatch event: %w", err)
			}
		case events.OpCodeResume:
			// Do nothing for now
		// case events.OpCodeReconnect:
		// TODO: Reconnect logic
		case events.OpCodeInvalidSession:
			slog.Info("Got invalid session.", "event", event)
			if canResume, err := UnmarshalJSON[bool](*event.Data); err != nil {
				return false, fmt.Errorf("discord sent invalid 'invalid session': %w", err)
			} else {
				return canResume, nil
			}
		case events.OpCodeHello:
			hello, err := UnmarshalJSON[events.Hello](*event.Data)
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

		case events.OpCodeHeartbeatAck:
			slog.Info("Got heartbeat ack event.")
		default:
			slog.Info("Got unknown event.", "event", fmt.Sprintf("%#v", event))
		}
	}
}

func (b *Bot) handleDispatch(event events.Event) error {
	if event.Type == nil {
		return fmt.Errorf("discord sent dispatch without type set")
	}

	eventType := *event.Type

	switch eventType {
	case "GUILD_CREATE":
		guildEvent, err := UnmarshalJSON[events.GuildCreate](*event.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal create guild json: %w", err)
		}

		if guildEvent.Unavailable != nil && *guildEvent.Unavailable {
			b.unavailableGuilds[guildEvent.ID] = guildEvent.Guild
		} else {
			b.guilds[guildEvent.ID] = guildEvent.Guild
		}

		b.fetchersByGuild[guildEvent.ID] = newFetcher(guildEvent)
	case "READY":
		readyEvent, err := UnmarshalJSON[events.Ready](*event.Data)
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

	default:
		slog.Warn("Unparsed dispatch event", "type", eventType)
	}

	return nil
}

func (b *Bot) identify() error {
	identifyPayload := events.Identify{
		OpCode: events.OpCodeIdentity,
		Payload: events.IdentifyPayload{
			Token: b.token,
			Intents: IntentGuilds | IntentGuildMembers | IntentGuildModeration | IntentGuildEmojisAndStickers |
				IntentGuildIntegrations | IntentGuildWebhooks | IntentGuildInvites | IntentGuildVoiceStates |
				IntentGuildPresences | IntentGuildMessages | IntentGuildMessageReactions | IntentGuildMessageTyping |
				IntentDirectMessages | IntentDirectMessageReactions | IntentDirectMessageTyping | IntentMessageContent |
				IntentGuildScheduledEvents | IntentAutoModerationConfiguration | IntentAutoModerationExecution,
			Properties: events.Properties{
				OS:      "raspbian",
				Browser: "webgockets",
				Device:  "discordgo",
			},
			Presence: &events.Presence{
				Activities: []*events.Activity{
					{
						Name: "developing simulator or something",
						Type: events.ActivityGame,
						URL:  "https://github.com/Hagesjo/webgockets",
					},
				},
				Status: events.UserStatusOnline,
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
	resume := events.Resume{
		OpCode: events.OpCodeResume,
		Payload: events.ResumePayload{
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
			if _, err := b.wsClient.WriteJSON(events.Heartbeat{
				OpCode: events.OpCodeHeartbeat,
			}); err != nil {
				return fmt.Errorf("failed to send heartbeat: %w", err)
			}

			slog.Info("Sent heartbeat")

			ticker.Reset(intervalWithJitter(b.heartbeatInterval))
		}
	}
}
