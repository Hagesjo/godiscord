package discordgo

import (
	"fmt"
	"net/http"

	"github.com/hagesjo/webgockets"
)

func NewBot(token string) (*Bot, error) {
	wsClient, err := webgockets.NewClient("wss://gateway.discord.gg")
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}

	restClient, err := newRestClient(http.DefaultClient, token)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate rest client: %w", err)
	}

	return &Bot{
		token:      token,
		wsClient:   wsClient,
		restClient: restClient,
	}, nil
}

type Bot struct {
	token      string
	wsClient   *webgockets.Client
	restClient *restClient
}

func (b *Bot) Run() error {
	b.handshake()
	// if err := b.wsClient.Connect(); err != nil {
	// 	return fmt.Errorf("failed to connect to discord: %w", err)
	// }

	// for {
	// 	s, err := b.wsClient.ReadString()
	// 	fmt.Println(err)
	// 	fmt.Println(s)
	// }
	return nil
}

func (b *Bot) handshake() error {
	fmt.Println(b.restClient.GetGatewayURL())

	return nil
}
