package godiscord

import (
	"bytes"
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

type transport struct {
	authToken           string
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", t.authToken))
	return t.underlyingTransport.RoundTrip(req)
}

type RestClient interface {
}

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

	restClient.httpClient.Transport = &transport{
		underlyingTransport: http.DefaultTransport,
		authToken:           authToken,
	}

	return restClient, nil
}

type restClient struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

type apiError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("discord api error: %d: %s", e.Code, e.Message)
}

func (c *restClient) get(path string, resp any) error {
	return c.do(http.MethodGet, path, nil, resp)
}

func (c *restClient) delete(path string, resp any) error {
	return c.do(http.MethodDelete, path, nil, resp)
}

func (c *restClient) patch(path string, req, resp any) error {
	return c.do(http.MethodPatch, path, req, resp)
}

func (c *restClient) put(path string, req, resp any) error {
	return c.do(http.MethodPut, path, req, resp)
}

func (c *restClient) post(path string, req, resp any) error {
	return c.do(http.MethodPost, path, req, resp)
}

// do is a helper function for doing a get request.
// The jsonResp sent in should be a pointer to the json struct.
func (c *restClient) do(method string, path string, reqStruct any, respStruct any) error {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return fmt.Errorf("failed to build url: %w", err)
	}

	slog.Info("Making request.", "method", method, "path", path)

	var body *bytes.Reader
	if reqStruct != nil {
		bs, err := json.Marshal(reqStruct)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		body = bytes.NewReader(bs)
	}

	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, u, body)
		if req != nil {
			req.Header.Add("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(method, u, nil)
	}
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("User-Agent", "DiscordBot (https://github.com/hagesjo/godiscord, dev)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		var apiErr apiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("failed to decode error response, status code %d: %w", resp.StatusCode, err)
		}

		return &apiErr
	}

	if respStruct != nil {
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(respStruct); err != nil {
			return fmt.Errorf("failed to decode json response: %w", err)
		}
	}

	return nil
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
	var resp GetGatewayURLResp
	err := c.get("/gateway/bot", &resp)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?v=%d&encoding=json", resp.URL, apiVersion), nil
}

// MessageCreateRequest is the request used for creating a message.
// At least one of content, embeds, sticker_ids, components, or files[n] is required.
// An example:
//
//	{
//	  "content": "Hello, World!",
//	  "tts": false,
//	  "embeds": [{
//	    "title": "Hello, Embed!",
//	    "description": "This is an embedded message."
//	  }]
//	}
type MessageCreateRequest struct {
	Content          string              `json:"content,omitempty"`           // Message contents (up to 2000 characters)
	Nonce            string              `json:"nonce,omitempty"`             // Can be used to verify a message was sent (up to 25 characters). Value will appear in the Message Create event.
	TTS              bool                `json:"tts,omitempty"`               // true if this is a TTS message
	Embeds           []Embed             `json:"embeds,omitempty"`            // Up to 10 rich embeds (up to 6000 characters)
	AllowedMentions  *AllowedMentions    `json:"allowed_mentions,omitempty"`  // Allowed mentions for the message
	MessageReference *MessageReference   `json:"message_reference,omitempty"` // Include to make your message a reply
	Components       []MessageActionType `json:"components,omitempty"`        // Components to include with the message
	StickerIDs       []string            `json:"sticker_ids,omitempty"`       // IDs of up to 3 stickers in the server to send in the message
	Files            map[string][]byte   `json:"files,omitempty"`             // Contents of the file being sent. See Uploading Files
	PayloadJSON      string              `json:"payload_json,omitempty"`      // JSON-encoded body of non-file params, only for multipart/form-data requests. See Uploading Files
	Attachments      []MessageAttachment `json:"attachments,omitempty"`       // Attachment objects with filename and description. See Uploading Files
	Flags            int                 `json:"flags,omitempty"`             // Message flags combined as a bitfield (only SUPPRESS_EMBEDS and SUPPRESS_NOTIFICATIONS can be set)
}

// Do is for using any discord endpoint not implemented here.
func (c *restClient) Do(path string, req any) error {
	err := c.post(path, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *restClient) GetGuild(guildID string) error {
	path := fmt.Sprintf("/guilds/%s", guildID)
	var resp Guild
	if err := c.get(path, &resp); err != nil {
		return err
	}

	return nil
}

// TODO: ModifyGuild
// TODO: GetGuildPreview
// TODO: DeleteGuild

func (c *restClient) MessageSend(channelID string, req MessageCreateRequest) error {
	path := fmt.Sprintf("/channels/%s/messages", channelID)
	err := c.post(path, req, nil)
	if err != nil {
		return err
	}

	return nil
}

type CreateChannelRequest struct {
	Name      *string     `json:"name,omitempty"`
	Type      ChannelType `json:"type"`
	Topic     *string     `json:"topic,omitempty"`
	Bitrate   *int        `json:"bitrate,omitempty"`
	UserLimit *int        `json:"user_limit,omitempty"`
	// RateLimitPerUser also applies to thread creation. Users can send one message and create one thread during each rate_limit_per_user interval.
	RateLimitPerUser     *int                  `json:"rate_limit_per_user,omitempty"`
	Position             *int                  `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	// ParentID is the parent id of the channel.
	// For guild channels: id of the parent category for a channel (each parent category can contain up to 50 channels).
	// For threads: id of the text channel this thread was created.
	ParentID                      *string             `json:"parent_id,omitempty"`
	NSFW                          bool                `json:"nsfw"`
	RTCRegion                     *string             `json:"rtc_region,omitempty"`
	VideoQualityMode              *VideoQualityMode   `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji          *DefaultReaction    `json:"default_reaction_emoji,omitempty"`
	AvailableTags                 []ForumTag          `json:"available_tags,omitempty"`
	DefaultSortOrder              *ChannelSortOrder   `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *ChannelForumLayout `json:"default_forum_layout,omitempty"`
	DefaultThreadRateLimitPerUser *int                `json:"default_thread_rate_limit_per_user,omitempty"`
}

func (c *restClient) CreateChannel(guildID string, req CreateChannelRequest) error {
	path := fmt.Sprintf("/guilds/%s/channels", guildID)
	err := c.post(path, req, nil)
	if err != nil {
		return err
	}

	return nil
}

type ModifyChannelOrder struct {
	ID              string  `json:"id"`               // snowflake
	ChannelID       string  `json:"channel_id"`       // snowflake
	Position        *int    `json:"position"`         // ?integer
	LockPermissions *bool   `json:"lock_permissions"` // ?boolean
	ParentID        *string `json:"parent_id"`        // ?snowflake
}

type ModifyChannelOrderRequest []ModifyChannelOrder

func (c *restClient) ModifyChannelOrder(channelID string, req CreateChannelRequest) error {
	path := fmt.Sprintf("/guilds/%s/channels", channelID)
	err := c.post(path, req, nil)
	if err != nil {
		return err
	}

	return nil
}
