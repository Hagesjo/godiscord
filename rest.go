package discordgo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	restBaseURL = "https://discord.com"
	apiVersion  = 10
)

func newRestClient(client *http.Client, authToken string) (*restClient, error) {
	baseURL, err := url.JoinPath(restBaseURL, fmt.Sprintf("api/v%d", apiVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to build base url: %w", err)
	}

	restClient := &restClient{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
		authToken:  authToken,
	}

	if client != nil {
		restClient.httpClient = client
	}

	return restClient, nil
}

type restClient struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

func (c *restClient) get(path string) (*http.Response, error) {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", c.authToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return resp, nil
}

type GetGatewayURLResp struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

func (c *restClient) GetGatewayURL() (string, error) {
	resp, err := c.get("/gateway/bot")
	if err != nil {
		return "", err
	}

	var jsonResp GetGatewayURLResp
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return "", fmt.Errorf("failed to decode json response: %w", err)
	}

	slog.Info("got gateway response", "response", jsonResp)

	return jsonResp.URL, nil
}
