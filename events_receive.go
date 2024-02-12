package discordgo

import (
	"fmt"
	"time"
)

// Hello represents the message sent on connection to the websocket, defining the heartbeat interval.
type Hello struct {
	HeartbeatInterval int `json:"heartbeat_interval"` // Interval (in milliseconds) an app should heartbeat with.
}

// Ready represents the ready event dispatched when a client has completed the initial handshake with the gateway.
type Ready struct {
	APIVersion       int         `json:"v"`                  // API version
	User             User        `json:"user"`               // Information about the user including email
	Guilds           []Guild     `json:"guilds"`             // Guilds the user is in
	SessionID        string      `json:"session_id"`         // Used for resuming connections
	ResumeGatewayURL string      `json:"resume_gateway_url"` // Gateway URL for resuming connections
	Shard            []int       `json:"shard,omitempty"`    // Shard information associated with this session, if sent when identifying
	Application      Application `json:"application"`        // Contains id and flags
}

// ApplicationCommandPermissionsUpdate represents the permissions for an application's command(s) in a guild.
type ApplicationCommandPermissionsUpdate struct {
	ID            string                          `json:"id"`             // ID of the command or the application ID
	ApplicationID string                          `json:"application_id"` // ID of the application the command belongs to
	GuildID       string                          `json:"guild_id"`       // ID of the guild
	Permissions   []ApplicationCommandPermissions `json:"permissions"`    // Permissions for the command in the guild, max of 100
}

func (e *ApplicationCommandPermissionsUpdate) Name() string {
	return "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
}

// AutoModerationRuleCreate is received when a rule is created.
type AutoModerationRuleCreate struct {
	AutoModerationRule
}

func (e *AutoModerationRuleCreate) Name() string {
	return "AUTO_MODERATION_RULE_CREATE"
}

// AutoModerationRuleUpdate is received when a rule is updated.
type AutoModerationRuleUpdate struct {
	AutoModerationRule
}

func (e *AutoModerationRuleUpdate) Name() string {
	return "AUTO_MODERATION_RULE_UPDATE"
}

// AutoModerationRuleDelete is received when a rule is deleted.
type AutoModerationRuleDelete struct {
	AutoModerationRule
}

func (e *AutoModerationRuleDelete) Name() string {
	return "AUTO_MODERATION_RULE_DELETE"
}

// ChannelCreate is received when a channel is created.
type ChannelCreate struct {
	Channel
}

func (e *ChannelCreate) Name() string {
	return "CHANNEL_CREATE"
}

// ChannelUpdate is received when a channel is updated.
type ChannelUpdate struct {
	Channel
}

func (e *ChannelUpdate) Name() string {
	return "CHANNEL_UPDATE"
}

// ChannelDelete is received when a channel is deleted.
type ChannelDelete struct {
	Channel
}

func (e *ChannelDelete) Name() string {
	return "CHANNEL_DELETE"
}

// ChannelPinsUpdate is sent when a message is pinned or unpinned in a text channel.
// This is not sent when a pinned message is deleted.
type ChannelPinsUpdate struct {
	GuildID          string     `json:"guild_id,omitempty"`           // ID of the guild
	ChannelID        string     `json:"channel_id"`                   // ID of the channel
	LastPinTimestamp *time.Time `json:"last_pin_timestamp,omitempty"` // Time at which the most recent pinned message was pinned. If nil means that the message was unpinned.
}

func (e *ChannelPinsUpdate) Name() string {
	return "CHANNEL_PINS_UPDATE"
}

// ThreadCreate is received when a thread is created.
type ThreadCreate struct {
	Channel Channel `json:"channel"`
}

func (e *ThreadCreate) Name() string {
	return "THREAD_CREATE"
}

// ThreadUpdate is received when a thread is updated.
type ThreadUpdate struct {
	Channel Channel `json:"channel"`
}

func (e *ThreadUpdate) Name() string {
	return "THREAD_UPDATE"
}

// ThreadDelete is received when a thread is deleted.
type ThreadDelete struct {
	Channel Channel `json:"channel"`
}

func (e *ThreadDelete) Name() string {
	return "THREAD_DELETE"
}

// ThreadListSync is received when gaining access to a channel, contains all active threads in that channel.
type ThreadListSync struct {
	GuildID    string         `json:"guild_id"`              // ID of the guild
	ChannelIDs []string       `json:"channel_ids,omitempty"` // Parent channel IDs whose threads are being synced
	Threads    []Channel      `json:"threads"`               // All active threads in the given channels that the current user can access
	Members    []ThreadMember `json:"members"`               // All thread member objects from the synced threads for the current user
}

func (e *ThreadListSync) Name() string {
	return "THREAD_LIST_SYNC"
}

// ThreadMemberUpdate is received when thread member for the current user is updated.
type ThreadMemberUpdate struct {
	ThreadMember

	GuildID string `json:"guild_id"` // ID of the guild
}

func (e *ThreadMemberUpdate) Name() string {
	return "THREAD_MEMBER_UPDATE"
}

// ThreadMembersUpdate is received when some user(s) were added to or removed from a thread.
type ThreadMembersUpdate struct {
	ID               string         `json:"id"`                           // ID of the thread
	GuildID          string         `json:"guild_id"`                     // ID of the guild
	MemberCount      int            `json:"member_count"`                 // Approximate number of members in the thread, capped at 50
	AddedMembers     []ThreadMember `json:"added_members,omitempty"`      // Users who were added to the thread
	RemovedMemberIDs []string       `json:"removed_member_ids,omitempty"` // ID of the users who were removed from the thread
}

func (e *ThreadMembersUpdate) Name() string {
	return "THREAD_MEMBERS_UPDATE"
}

// EntitlementCreate is received when an entitlement is created.
type EntitlementCreate struct {
	Entitlement
}

func (e *EntitlementCreate) Name() string {
	return "ENTITLEMENT_CREATE"
}

// EntitlementUpdate is received when an entitlement is updated.
type EntitlementUpdate struct {
	Entitlement
}

func (e *EntitlementUpdate) Name() string {
	return "ENTITLEMENT_UPDATE"
}

// EntitlementDelete is received when an entitlement is deleted.
type EntitlementDelete struct {
	Entitlement
}

func (e *EntitlementDelete) Name() string {
	return "ENTITLEMENT_DELETE"
}

// GuildCreate represents the event sent when a user initially connects or when a guild becomes available again to the client, or when the current user joins a new guild.
// If your bot does not have the GUILD_PRESENCES Gateway Intent,
// or if the guild has over 75k members,
// members and presences returned in this event will only contain your bot and users in voice channels.
type GuildCreate struct {
	Guild
	JoinedAt             time.Time             `json:"joined_at"`              // When this guild was joined at
	Large                bool                  `json:"large"`                  // true if this is considered a large guild
	Unavailable          *bool                 `json:"unavailable,omitempty"`  // true if this guild is unavailable due to an outage
	MemberCount          int                   `json:"member_count"`           // Total number of members in this guild
	VoiceStates          []VoiceState          `json:"voice_states"`           // States of members currently in voice channels; lacks the guild_id key
	Members              []GuildMember         `json:"members"`                // Users in the guild
	Channels             []Channel             `json:"channels"`               // Channels in the guild
	Threads              []Channel             `json:"threads"`                // All active threads in the guild that current user has permission to view
	Presences            []Presence            `json:"presences"`              // Presences of the members in the guild, will only include non-offline members if the size is greater than large threshold
	StageInstances       []StageInstance       `json:"stage_instances"`        // Stage instances in the guild
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"` // Scheduled events in the guild
}

func (e *GuildCreate) Name() string {
	return "GUILD_CREATE"
}

// GuildUpdate is received when a guild is updated
type GuildUpdate struct {
	Guild Guild `json:"guild"`
}

func (e *GuildUpdate) Name() string {
	return "GUILD_UPDATE"
}

// GuildDelete is received when a guild is deleted.
type GuildDelete struct {
	Guild Guild `json:"guild"`
}

func (e *GuildDelete) Name() string {
	return "GUILD_DELETE"
}

// GuildBanAdd is received when a user was banned from a guild.
type GuildBanAdd struct {
	GuildID string `json:"guild_id"` // ID of the guild.
	User    User   `json:"user"`     // User who was banned.
}

func (e *GuildBanAdd) Name() string {
	return "GUILD_BAN_ADD"
}

// GuildBanRemove is received when a user was unbanned from a guild.
type GuildBanRemove struct {
	GuildID string `json:"guild_id"` // ID of the guild.
	User    User   `json:"user"`     // User who was unbanned.
}

func (e *GuildBanRemove) Name() string {
	return "GUILD_BAN_REMOVE"
}

// GuildEmojisUpdate is received when guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID string  `json:"guild_id"` // ID of the guild
	Emojis  []Emoji `json:"emojis"`   // Array of emojis
}

func (e *GuildEmojisUpdate) Name() string {
	return "GUILD_EMOJIS_UPDATE"
}

// GuildStickersUpdate is received when guild stickers were updated.
type GuildStickersUpdate struct {
	GuildID  string    `json:"guild_id"` // ID of the guild.
	Stickers []Sticker `json:"stickers"` // Array of stickers.
}

func (e *GuildStickersUpdate) Name() string {
	return "GUILD_STICKERS_UPDATE"
}

// GuildMemberAdd is received when a user joins a guild.
type GuildMemberAdd struct {
	GuildMember

	GuildID string `json:"guild_id"` // ID of the guild.
}

func (e *GuildMemberAdd) Name() string {
	return "GUILD_MEMBER_ADD"
}

// GuildMemberUpdate is received when a guild member was updated.
type GuildMemberUpdate struct {
	GuildID                    string     `json:"guild_id"`                     // ID of the guild
	Roles                      []string   `json:"roles"`                        // User role ids
	User                       User       `json:"user"`                         // User
	Nick                       *string    `json:"nick,omitempty"`               // Nickname of the user in the guild
	Avatar                     *string    `json:"avatar,omitempty"`             // Member's guild avatar hash
	JoinedAt                   *time.Time `json:"joined_at,omitempty"`          // When the user joined the guild
	PremiumSince               *time.Time `json:"premium_since,omitempty"`      // When the user starting boosting the guild
	Deaf                       *bool      `json:"deaf"`                         // Whether the user is deafened in voice channels
	Mute                       *bool      `json:"mute"`                         // Whether the user is muted in voice channels
	Pending                    *bool      `json:"pending"`                      // Whether the user has not yet passed the guild's Membership Screening requirements
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"` // When the user's timeout will expire and the user will be able to communicate in the guild again, null or a time in the past if the user is not timed out
}

func (e *GuildMemberUpdate) Name() string {
	return "GUILD_MEMBER_UPDATE"
}

// GuildMemberRemove is received when a user was removed from the guild.
type GuildMemberRemove struct {
	GuildID string `json:"guild_id"` // ID of the guild
	User    User   `json:"user"`     // User who was removed
}

func (e *GuildMemberRemove) Name() string {
	return "GUILD_MEMBER_REMOVE"
}

// GuildMembersChunk is received in response to Guild Request Members. You can use the chunk_index and chunk_count to calculate how many chunks are left for your request.
type GuildMembersChunk struct {
	GuildID    string        `json:"guild_id"`            // ID of the guild
	Members    []GuildMember `json:"members"`             // Set of guild members
	ChunkIndex int           `json:"chunk_index"`         // Chunk index in the expected chunks for this response (0 <= chunk_index < chunk_count)
	ChunkCount int           `json:"chunk_count"`         // Total number of expected chunks for this response
	NotFound   []string      `json:"not_found,omitempty"` // When passing an invalid ID to REQUEST_GUILD_MEMBERS, it will be returned here
	Presences  []Presence    `json:"presences,omitempty"` // When passing true to REQUEST_GUILD_MEMBERS, presences of the returned members will be here
	Nonce      string        `json:"nonce,omitempty"`     // Nonce used in the Guild Members Request
}

func (e *GuildMembersChunk) Name() string {
	return "GUILD_MEMBERS_CHUNK"
}

// GuildRoleCreate is received when a new role was created.
type GuildRoleCreate struct {
	GuildID string `json:"guild_id"` // ID of the guild
	Role    Role   `json:"role"`     // Role that was created
}

func (e *GuildRoleCreate) Name() string {
	return "GUILD_ROLE_CREATE"
}

// GuildRoleUpdate is received when a role was updated.
type GuildRoleUpdate struct {
	GuildID string `json:"guild_id"` // ID of the guild
	Role    Role   `json:"role"`     // Role that was updated
}

func (e *GuildRoleUpdate) Name() string {
	return "GUILD_ROLE_UPDATE"
}

// GuildRoleDelete is received when a role was deleted.
type GuildRoleDelete struct {
	GuildID string `json:"guild_id"` // ID of the guild
	RoleID  string `json:"role_id"`  // ID of the role
}

func (e *GuildRoleDelete) Name() string {
	return "GUILD_ROLE_DELETE"
}

// GuildScheduledEventCreate is received when a new scheduled event is created.
type GuildScheduledEventCreate struct {
	GuildScheduledEvent
}

func (e *GuildScheduledEventCreate) Name() string {
	return "GUILD_SCHEDULED_EVENT_CREATE"
}

// GuildScheduledEventUpdate is received when a new scheduled event is updated.
type GuildScheduledEventUpdate struct {
	GuildScheduledEvent
}

func (e *GuildScheduledEventUpdate) Name() string {
	return "GUILD_SCHEDULED_EVENT_UPDATE"
}

// GuildScheduledEventDelete is received when a new scheduled event is deleted.
type GuildScheduledEventDelete struct {
	GuildScheduledEvent
}

func (e *GuildScheduledEventDelete) Name() string {
	return "GUILD_SCHEDULED_EVENT_DELETE"
}

// GuildScheduledEventUserAddEvent is received when a user has subscribed to a guild scheduled event.
type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

func (e *GuildScheduledEventUserAddEvent) Name() string {
	return "GUILD_SCHEDULED_EVENT_USER_ADD_EVENT"
}

// GuildScheduledEventUserRemoveEvent is received when a user has unsubscribed from a guild scheduled event.
type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

func (e *GuildScheduledEventUserRemoveEvent) Name() string {
	return "GUILD_SCHEDULED_EVENT_USER_REMOVE_EVENT"
}

// IntegrationCreate is received when a new integration is created.
type IntegrationCreate struct {
	Integration

	GuildID string `json:"guild_id"`
}

func (e *IntegrationCreate) Name() string {
	return "INTEGRATION_CREATE"
}

// IntegrationCreate is received when a new integration is updated.
type IntegrationUpdate struct {
	Integration

	GuildID string `json:"guild_id"`
}

func (e *IntegrationUpdate) Name() string {
	return "INTEGRATION_UPDATE"
}

// IntegrationDelete is sent when an integration is deleted.
type IntegrationDelete struct {
	ID            string  `json:"id"`                       // Integration ID.
	GuildID       string  `json:"guild_id"`                 // ID of the guild.
	ApplicationID *string `json:"application_id,omitempty"` // ID of the bot/OAuth2 application for this Discord integration.
}

func (e *IntegrationDelete) Name() string {
	return "INTEGRATION_DELETE"
}

// InviteCreate is sent when an invite to a channel is created.
type InviteCreate struct {
	ChannelID         string           `json:"channel_id"`                   // Channel the invite is for.
	Code              string           `json:"code"`                         // Unique invite code.
	CreatedAt         time.Time        `json:"created_at"`                   // Time at which the invite was created.
	GuildID           string           `json:"guild_id,omitempty"`           // Guild of the invite.
	Inviter           *User            `json:"inviter,omitempty"`            // User that created the invite.
	MaxAge            int              `json:"max_age"`                      // How long the invite is valid for (in seconds).
	MaxUses           int              `json:"max_uses"`                     // Maximum number of times the invite can be used.
	TargetType        InviteTargetType `json:"target_type,omitempty"`        // Type of target for this voice channel invite.
	TargetUser        *User            `json:"target_user,omitempty"`        // User whose stream to display for this voice channel stream invite.
	TargetApplication *Application     `json:"target_application,omitempty"` // Embedded application to open for this voice channel embedded application invite.
	Temporary         bool             `json:"temporary"`                    // Whether or not the invite is temporary (invited users will be kicked on disconnect unless they're assigned a role).
	Uses              int              `json:"uses"`                         // How many times the invite has been used (always will be 0).
}

func (e *InviteCreate) Name() string {
	return "INVITE_CREATE"
}

// InviteDelete is sent when an invite is deleted.
type InviteDelete struct {
	ChannelID string `json:"channel_id"`         // Channel of the invite.
	GuildID   string `json:"guild_id,omitempty"` // Guild of the invite.
	Code      string `json:"code"`               // Unique invite code.
}

func (e *InviteDelete) Name() string {
	return "INVITE_DELETE"
}

// MessageCreate is sent when a message is created.
type MessageCreate struct {
	Message

	ID        string       `json:"id"`                 // The ID of the message.
	ChannelID string       `json:"channel_id"`         // The ID of the channel the message was sent in.
	GuildID   string       `json:"guild_id,omitempty"` // ID of the guild the message was sent in - unless it is an ephemeral message.
	Content   string       `json:"content"`            // The content of the message.
	Author    User         `json:"author"`             // The user who sent the message.
	Member    *GuildMember `json:"member,omitempty"`   // Member properties for this message's author. Missing for ephemeral messages and messages from webhooks.
	Mentions  []User       `json:"mentions,omitempty"` // Users specifically mentioned in the message.
}

func (e *MessageCreate) Name() string {
	return "MESSAGE_CREATE"
}

// MessageUpdate is sent when a message is created.
type MessageUpdate struct {
	Message

	ID        string       `json:"id"`                 // The ID of the message.
	ChannelID string       `json:"channel_id"`         // The ID of the channel the message was sent in.
	GuildID   string       `json:"guild_id,omitempty"` // ID of the guild the message was sent in - unless it is an ephemeral message.
	Content   string       `json:"content"`            // The content of the message.
	Author    User         `json:"author"`             // The user who sent the message.
	Member    *GuildMember `json:"member,omitempty"`   // Member properties for this message's author. Missing for ephemeral messages and messages from webhooks.
	Mentions  []User       `json:"mentions,omitempty"` // Users specifically mentioned in the message.
}

func (e *MessageUpdate) Name() string {
	return "MESSAGE_UPDATE"
}

// MessageDelete is sent with a message delete event.
type MessageDelete struct {
	ID        string `json:"id"`                 // ID of the message.
	ChannelID string `json:"channel_id"`         // ID of the channel.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (e *MessageDelete) Name() string {
	return "MESSAGE_DELETE"
}

// MessageDeleteBulk is sent with a bulk message delete event.
type MessageDeleteBulk struct {
	IDs       []string `json:"ids"`                // IDs of the deleted messages.
	ChannelID string   `json:"channel_id"`         // ID of the channel.
	GuildID   string   `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (e *MessageDeleteBulk) Name() string {
	return "MESSAGE_DELETE_BULK"
}

// MessageReactionAdd is sent when a user adds a reaction to a message.
type MessageReactionAdd struct {
	UserID          string       `json:"user_id"`                     // ID of the user.
	ChannelID       string       `json:"channel_id"`                  // ID of the channel.
	MessageID       string       `json:"message_id"`                  // ID of the message.
	GuildID         string       `json:"guild_id,omitempty"`          // ID of the guild (optional).
	Member          *GuildMember `json:"member,omitempty"`            // Member who reacted if this happened in a guild.
	Emoji           *Emoji       `json:"emoji"`                       // Emoji used to react.
	MessageAuthorID string       `json:"message_author_id,omitempty"` // ID of the user who authored the message which was reacted to.
}

type messageReactionAddHandler struct {
	f func(*Fetcher, MessageReactionAdd) error
}

func (e messageReactionAddHandler) Name() string {
	return "MESSAGE_REACTION_ADD"
}

func (e messageReactionAddHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageReactionAdd))
}

// MessageReactionRemove is sent when a user removes a reaction from a message.
type MessageReactionRemove struct {
	UserID    string `json:"user_id"`            // ID of the user.
	ChannelID string `json:"channel_id"`         // ID of the channel.
	MessageID string `json:"message_id"`         // ID of the message.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
	Emoji     *Emoji `json:"emoji"`              // Emoji used to react.
}

func (e *MessageReactionRemove) Name() string {
	return "MESSAGE_REACTION_REMOVE"
}

// MessageReactionRemoveAll is sent when a user explicitly removes all reactions from a message.
type MessageReactionRemoveAll struct {
	ChannelID string `json:"channel_id"`         // ID of the channel.
	MessageID string `json:"message_id"`         // ID of the message.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (e *MessageReactionRemoveAll) Name() string {
	return "MESSAGE_REACTION_REMOVE_ALL"
}

// MessageReactionRemoveEmoji is sent when a user removes a reaction with a specific emoji from a message.
type MessageReactionRemoveEmoji struct {
	ChannelID string `json:"channel_id"`         // ID of the channel.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
	MessageID string `json:"message_id"`         // ID of the message.
	Emoji     Emoji  `json:"emoji"`              // Emoji that was removed.
}

func (e *MessageReactionRemoveEmoji) Name() string {
	return "MESSAGE_REACTION_REMOVE_EMOJI"
}

// PresenceUpdate is sent when a user's presence in a guild is updated.
type PresenceUpdate struct {
	User         User         `json:"user"`          // User whose presence is being updated
	GuildID      string       `json:"guild_id"`      // ID of the guild
	Status       string       `json:"status"`        // Either "idle", "dnd", "online", or "offline"
	Activities   []Activity   `json:"activities"`    // User's current activities
	ClientStatus ClientStatus `json:"client_status"` // User's platform-dependent status
}

func (e *PresenceUpdate) Name() string {
	return "PRESENCE_UPDATE"
}

// StageInstanceCreate is received when a new stage instance is created.
type StageInstanceCreate struct {
	StageInstance
}

func (e *StageInstanceCreate) Name() string {
	return "STAGE_INSTANCE_CREATE"
}

// StageInstanceCreate is received when a new stage instance is updated.
type StageInstanceUpdate struct {
	StageInstance
}

func (e *StageInstanceUpdate) Name() string {
	return "STAGE_INSTANC_EUPDATE"
}

// StageInstanceDelete is received when an stage instance is deleted.
type StageInstanceDelete struct {
	StageInstance
}

func (e *StageInstanceDelete) Name() string {
	return "STAGE_INSTANCE_DELETE"
}

// TypingStart is received when a user starts typing in a channel.
type TypingStart struct {
	ChannelID string       `json:"channel_id"`         // ID of the channel
	GuildID   *string      `json:"guild_id,omitempty"` // ID of the guild.
	UserID    string       `json:"user_id"`            // ID of the user
	Timestamp int          `json:"timestamp"`          // Unix time (in seconds) of when the user started typing
	Member    *GuildMember `json:"member,omitempty"`   // Member who started typing if this happened in a guild
}

func (e *TypingStart) Name() string {
	return "TYPING_START"
}

// UserUpdate is received when a user is updated.
type UserUpdate struct {
	User
}

func (e *UserUpdate) Name() string {
	return "USER_UPDATE"
}

// VoiceStateUpdate is received when someone joins/leaves/moves voice channels.
type VoiceStateUpdate struct {
	VoiceState
}

func (e *VoiceStateUpdate) Name() string {
	return "VOICE_STATE_UPDATE"
}

// VoiceServerUpdate is received when a guild's voice server is updated. This is sent when initially connecting to voice, and when the current voice instance fails over to a new server.
// A null endpoint means that the voice server allocated has gone away and is trying to be reallocated. You should attempt to disconnect from the currently connected voice server, and not attempt to reconnect until a new voice server is allocated.
type VoiceServerUpdate struct {
	Token    string  `json:"token"`              // Voice connection token
	GuildID  string  `json:"guild_id"`           // Guild this voice server update is for
	Endpoint *string `json:"endpoint,omitempty"` // Voice server host
}

func (e *VoiceServerUpdate) Name() string {
	return "VOICE_SERVER_UPDATE"
}

// WebhooksUpdate is received when a guild channel's webhook is created, updated, or deleted.
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`   // ID of the guild
	ChannelID string `json:"channel_id"` // ID of the channel
}

func (e *WebhooksUpdate) Name() string {
	return "WEBHOOKS_UPDATE"
}

type eventHandler interface {
	run(*Fetcher, any) error
	name() string
}

func eventHandlerFromInterface(iface any) (eventHandler, error) {
	switch v := iface.(type) {
	case func(*Fetcher, MessageReactionAdd) error:
		return messageReactionAddHandler{f: v}, nil

	default:
		return nil, fmt.Errorf("unknown event")
	}
}
