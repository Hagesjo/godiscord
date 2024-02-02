package discordgo

import (
	"fmt"
	"net/http"

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
		wsClient:   wsClient,
		restClient: restClient,
		gatewayURL: gatewayURL, // NOTE: discord doesn't want this cached too long.
	}, nil
}

type Bot struct {
	token      string
	wsClient   *webgockets.Client
	restClient *restClient
	gatewayURL string
}

func (b *Bot) Run() error {
	b.handshake()
	if err := b.wsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to discord: %w", err)
	}

	// for {
	// 	s, err := b.wsClient.ReadString()
	// 	fmt.Println(err)
	// 	fmt.Println(s)
	// }

	return nil
}

func (b *Bot) handshake() error {

	return nil
}
