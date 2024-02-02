package discordgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/hagesjo/discordgo/events"
	"github.com/hagesjo/webgockets"
)

func NewBot(token string) (*Bot, error) {
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
		token:      token,
		gatewayURL: gatewayURL, // NOTE: discord doesn't want this cached too long.

		wsClient:   wsClient,
		restClient: restClient,
	}, nil
}

type Bot struct {
	wsClient   *webgockets.Client
	restClient *restClient

	token             string
	gatewayURL        string
	heartbeatInterval int
	lastSequence      *int
}

func (b *Bot) Run() error {
	if err := b.wsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to discord: %w", err)
	}
	return b.handler()
}

func (b *Bot) handler() error {
	heartbeatCtx, cancelHeartbeatCtx := context.WithCancel(context.Background())
	defer cancelHeartbeatCtx()
	heartbeatStarted := false

	for {
		event, err := webgockets.ReadJSON[events.Event](b.wsClient)
		if err == io.EOF { // TODO: reconnect logic.
			break
		} else if err != nil {
			return fmt.Errorf("failed to read event: %w", err)
		}

		switch event.OpCode {
		case events.OpCodeDispatch:

		case events.OpCodeHello:
			var hello events.Hello

			if err := json.Unmarshal(*event.Data, &hello); err != nil {
				fmt.Errorf("discord sent invalid hello '%s': %w", *event.Data, err)
			}

			slog.Info("got hello", "event", fmt.Sprintf("%#v", event), "hello", hello)

			b.heartbeatInterval = hello.HeartbeatInterval

			if heartbeatStarted {
				continue
			}

			go b.heartbeater(heartbeatCtx)
			heartbeatStarted = true

			if err := b.identify(); err != nil {
				return fmt.Errorf("failed to identify: %w", err)
			}
		case events.OpCodeHeartbeatAck:
			slog.Info("Got heartbeat ack event")
		default:
			slog.Info("got unknown event", "event", fmt.Sprintf("%#v", event))
		}
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
