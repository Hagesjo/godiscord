package godiscord

import (
	"fmt"
	"reflect"
	"time"
)

type guildEvent interface {
	guild() string
}

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

func (m ApplicationCommandPermissionsUpdate) guild() string {
	return m.GuildID
}

type applicationCommandPermissionsUpdateHandler struct {
	f func(*Fetcher, ApplicationCommandPermissionsUpdate) error
}

func (e applicationCommandPermissionsUpdateHandler) name() string {
	return "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
}

func (e applicationCommandPermissionsUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ApplicationCommandPermissionsUpdate))
}

// AutoModerationRuleCreate is received when a rule is created.
type AutoModerationRuleCreate struct {
	AutoModerationRule
}

func (m AutoModerationRuleCreate) guild() string {
	return m.GuildID
}

type autoModerationRuleCreateHandler struct {
	f func(*Fetcher, AutoModerationRuleCreate) error
}

func (e autoModerationRuleCreateHandler) name() string {
	return "AUTO_MODERATION_RULE_CREATE"
}

func (e autoModerationRuleCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(AutoModerationRuleCreate))
}

// AutoModerationRuleUpdate is received when a rule is updated.
type AutoModerationRuleUpdate struct {
	AutoModerationRule
}

func (m AutoModerationRuleUpdate) guild() string {
	return m.GuildID
}

type autoModerationRuleUpdateHandler struct {
	f func(*Fetcher, AutoModerationRuleUpdate) error
}

func (e autoModerationRuleUpdateHandler) name() string {
	return "AUTO_MODERATION_RULE_UPDATE"
}

func (e autoModerationRuleUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(AutoModerationRuleUpdate))
}

// AutoModerationRuleDelete is received when a rule is deleted.
type AutoModerationRuleDelete struct {
	AutoModerationRule
}

func (m AutoModerationRuleDelete) guild() string {
	return m.GuildID
}

type autoModerationRuleDeleteHandler struct {
	f func(*Fetcher, AutoModerationRuleDelete) error
}

func (e autoModerationRuleDeleteHandler) name() string {
	return "AUTO_MODERATION_RULE_DELETE"
}

func (e autoModerationRuleDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(AutoModerationRuleDelete))
}

// ChannelCreate is received when a channel is created.
type ChannelCreate struct {
	Channel
}

func (m ChannelCreate) guild() string {
	if m.GuildID == nil {
		return ""
	}

	return *m.GuildID
}

type channelCreateHandler struct {
	f func(*Fetcher, ChannelCreate) error
}

func (e channelCreateHandler) name() string {
	return "CHANNEL_CREATE"
}

func (e channelCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ChannelCreate))
}

// ChannelUpdate is received when a channel is updated.
type ChannelUpdate struct {
	Channel
}

func (m ChannelUpdate) guild() string {
	if m.GuildID == nil {
		return ""
	}

	return *m.GuildID
}

type channelUpdateHandler struct {
	f func(*Fetcher, ChannelUpdate) error
}

func (e channelUpdateHandler) name() string {
	return "CHANNEL_UPDATE"
}

func (e channelUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ChannelUpdate))
}

// ChannelDelete is received when a channel is deleted.
type ChannelDelete struct {
	Channel
}

func (m ChannelDelete) guild() string {
	if m.GuildID == nil {
		return ""
	}

	return *m.GuildID
}

type channelDeleteHandler struct {
	f func(*Fetcher, ChannelDelete) error
}

func (e channelDeleteHandler) name() string {
	return "CHANNEL_DELETE"
}

func (e channelDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ChannelDelete))
}

// ChannelPinsUpdate is sent when a message is pinned or unpinned in a text channel.
// This is not sent when a pinned message is deleted.
type ChannelPinsUpdate struct {
	GuildID          string     `json:"guild_id,omitempty"`           // ID of the guild
	ChannelID        string     `json:"channel_id"`                   // ID of the channel
	LastPinTimestamp *time.Time `json:"last_pin_timestamp,omitempty"` // Time at which the most recent pinned message was pinned. If nil means that the message was unpinned.
}

func (m ChannelPinsUpdate) guild() string {
	return m.GuildID
}

type channelPinsUpdateHandler struct {
	f func(*Fetcher, ChannelPinsUpdate) error
}

func (e channelPinsUpdateHandler) name() string {
	return "CHANNEL_PINS_UPDATE"
}

func (e channelPinsUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ChannelPinsUpdate))
}

// ThreadCreate is received when a thread is created.
type ThreadCreate struct {
	Channel Channel `json:"channel"`
}

func (m ThreadCreate) guild() string {
	if m.Channel.GuildID == nil {
		return ""
	}

	return *m.Channel.GuildID
}

type threadCreateHandler struct {
	f func(*Fetcher, ThreadCreate) error
}

func (e threadCreateHandler) name() string {
	return "THREAD_CREATE"
}

func (e threadCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadCreate))
}

// ThreadUpdate is received when a thread is updated.
type ThreadUpdate struct {
	Channel Channel `json:"channel"`
}

func (m ThreadUpdate) guild() string {
	if m.Channel.GuildID == nil {
		return ""
	}

	return *m.Channel.GuildID
}

type threadUpdateHandler struct {
	f func(*Fetcher, ThreadUpdate) error
}

func (e threadUpdateHandler) name() string {
	return "THREAD_UPDATE"
}

func (e threadUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadUpdate))
}

// ThreadDelete is received when a thread is deleted.
type ThreadDelete struct {
	Channel Channel `json:"channel"`
}

func (m ThreadDelete) guild() string {
	if m.Channel.GuildID == nil {
		return ""
	}

	return *m.Channel.GuildID
}

type threadDeleteHandler struct {
	f func(*Fetcher, ThreadDelete) error
}

func (e threadDeleteHandler) name() string {
	return "THREAD_DELETE"
}

func (e threadDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadDelete))
}

// ThreadListSync is received when gaining access to a channel, contains all active threads in that channel.
type ThreadListSync struct {
	GuildID    string         `json:"guild_id"`              // ID of the guild
	ChannelIDs []string       `json:"channel_ids,omitempty"` // Parent channel IDs whose threads are being synced
	Threads    []Channel      `json:"threads"`               // All active threads in the given channels that the current user can access
	Members    []ThreadMember `json:"members"`               // All thread member objects from the synced threads for the current user
}

func (m ThreadListSync) guild() string {
	return m.GuildID
}

type threadListSyncHandler struct {
	f func(*Fetcher, ThreadListSync) error
}

func (e threadListSyncHandler) name() string {
	return "THREAD_LIST_SYNC"
}

func (e threadListSyncHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadListSync))
}

// ThreadMemberUpdate is received when thread member for the current user is updated.
type ThreadMemberUpdate struct {
	ThreadMember

	GuildID string `json:"guild_id"` // ID of the guild
}

func (m ThreadMemberUpdate) guild() string {
	return m.GuildID
}

type threadMemberUpdateHandler struct {
	f func(*Fetcher, ThreadMemberUpdate) error
}

func (e threadMemberUpdateHandler) name() string {
	return "THREAD_MEMBER_UPDATE"
}

func (e threadMemberUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadMemberUpdate))
}

// ThreadMembersUpdate is received when some user(s) were added to or removed from a thread.
type ThreadMembersUpdate struct {
	ID               string         `json:"id"`                           // ID of the thread
	GuildID          string         `json:"guild_id"`                     // ID of the guild
	MemberCount      int            `json:"member_count"`                 // Approximate number of members in the thread, capped at 50
	AddedMembers     []ThreadMember `json:"added_members,omitempty"`      // Users who were added to the thread
	RemovedMemberIDs []string       `json:"removed_member_ids,omitempty"` // ID of the users who were removed from the thread
}

func (m ThreadMembersUpdate) guild() string {
	return m.GuildID
}

type threadMembersUpdateHandler struct {
	f func(*Fetcher, ThreadMembersUpdate) error
}

func (e threadMembersUpdateHandler) name() string {
	return "THREAD_MEMBERS_UPDATE"
}

func (e threadMembersUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(ThreadMembersUpdate))
}

// EntitlementCreate is received when an entitlement is created.
type EntitlementCreate struct {
	Entitlement
}

func (m EntitlementCreate) guild() string {
	if m.Entitlement.GuildID == nil {
		return ""
	}

	return *m.Entitlement.GuildID
}

type entitlementCreateHandler struct {
	f func(*Fetcher, EntitlementCreate) error
}

func (e entitlementCreateHandler) name() string {
	return "ENTITLEMENT_CREATE"
}

func (e entitlementCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(EntitlementCreate))
}

// EntitlementUpdate is received when an entitlement is updated.
type EntitlementUpdate struct {
	Entitlement
}

func (m EntitlementUpdate) guild() string {
	if m.Entitlement.GuildID == nil {
		return ""
	}

	return *m.Entitlement.GuildID
}

type entitlementUpdateHandler struct {
	f func(*Fetcher, EntitlementUpdate) error
}

func (e entitlementUpdateHandler) name() string {
	return "ENTITLEMENT_UPDATE"
}

func (e entitlementUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(EntitlementUpdate))
}

// EntitlementDelete is received when an entitlement is deleted.
type EntitlementDelete struct {
	Entitlement
}

func (m EntitlementDelete) guild() string {
	if m.Entitlement.GuildID == nil {
		return ""
	}

	return *m.Entitlement.GuildID
}

type entitlementDeleteHandler struct {
	f func(*Fetcher, EntitlementDelete) error
}

func (e entitlementDeleteHandler) name() string {
	return "ENTITLEMENT_DELETE"
}

func (e entitlementDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(EntitlementDelete))
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

func (m GuildCreate) guild() string {
	return m.Guild.ID
}

type guildCreateHandler struct {
	f func(*Fetcher, GuildCreate) error
}

func (e guildCreateHandler) name() string {
	return "GUILD_CREATE"
}

func (e guildCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildCreate))
}

// GuildUpdate is received when a guild is updated
type GuildUpdate struct {
	Guild Guild `json:"guild"`
}

func (m GuildUpdate) guild() string {
	return m.Guild.ID
}

type guildUpdateHandler struct {
	f func(*Fetcher, GuildUpdate) error
}

func (e guildUpdateHandler) name() string {
	return "GUILD_UPDATE"
}

func (e guildUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildUpdate))
}

// GuildDelete is received when a guild is deleted.
type GuildDelete struct {
	Guild Guild `json:"guild"`
}

func (m GuildDelete) guild() string {
	return m.Guild.ID
}

type guildDeleteHandler struct {
	f func(*Fetcher, GuildDelete) error
}

func (e guildDeleteHandler) name() string {
	return "GUILD_DELETE"
}

func (e guildDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildDelete))
}

// GuildBanAdd is received when a user was banned from a guild.
type GuildBanAdd struct {
	GuildID string `json:"guild_id"` // ID of the guild.
	User    User   `json:"user"`     // User who was banned.
}

func (m GuildBanAdd) guild() string {
	return m.GuildID
}

type guildBanAddHandler struct {
	f func(*Fetcher, GuildBanAdd) error
}

func (e guildBanAddHandler) name() string {
	return "GUILD_BAN_ADD"
}

func (e guildBanAddHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildBanAdd))
}

// GuildBanRemove is received when a user was unbanned from a guild.
type GuildBanRemove struct {
	GuildID string `json:"guild_id"` // ID of the guild.
	User    User   `json:"user"`     // User who was unbanned.
}

func (m GuildBanRemove) guild() string {
	return m.GuildID
}

type guildBanRemoveHandler struct {
	f func(*Fetcher, GuildBanRemove) error
}

func (e guildBanRemoveHandler) name() string {
	return "GUILD_BAN_REMOVE"
}

func (e guildBanRemoveHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildBanRemove))
}

// GuildEmojisUpdate is received when guild emojis were updated
type GuildEmojisUpdate struct {
	GuildID string  `json:"guild_id"` // ID of the guild
	Emojis  []Emoji `json:"emojis"`   // Array of emojis
}

func (m GuildEmojisUpdate) guild() string {
	return m.GuildID
}

type guildEmojisUpdateHandler struct {
	f func(*Fetcher, GuildEmojisUpdate) error
}

func (e guildEmojisUpdateHandler) name() string {
	return "GUILD_EMOJIS_UPDATE"
}

func (e guildEmojisUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildEmojisUpdate))
}

// GuildStickersUpdate is received when guild stickers were updated.
type GuildStickersUpdate struct {
	GuildID  string    `json:"guild_id"` // ID of the guild.
	Stickers []Sticker `json:"stickers"` // Array of stickers.
}

func (m GuildStickersUpdate) guild() string {
	return m.GuildID
}

type guildStickersUpdateHandler struct {
	f func(*Fetcher, GuildStickersUpdate) error
}

func (e guildStickersUpdateHandler) name() string {
	return "GUILD_STICKERS_UPDATE"
}

func (e guildStickersUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildStickersUpdate))
}

// GuildMemberAdd is received when a user joins a guild.
type GuildMemberAdd struct {
	GuildMember

	GuildID string `json:"guild_id"` // ID of the guild.
}

func (m GuildMemberAdd) guild() string {
	return m.GuildID
}

type guildMemberAddHandler struct {
	f func(*Fetcher, GuildMemberAdd) error
}

func (e guildMemberAddHandler) name() string {
	return "GUILD_MEMBER_ADD"
}

func (e guildMemberAddHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildMemberAdd))
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

func (m GuildMemberUpdate) guild() string {
	return m.GuildID
}

type guildMemberUpdateHandler struct {
	f func(*Fetcher, GuildMemberUpdate) error
}

func (e guildMemberUpdateHandler) name() string {
	return "GUILD_MEMBER_UPDATE"
}

func (e guildMemberUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildMemberUpdate))
}

// GuildMemberRemove is received when a user was removed from the guild.
type GuildMemberRemove struct {
	GuildID string `json:"guild_id"` // ID of the guild
	User    User   `json:"user"`     // User who was removed
}

func (m GuildMemberRemove) guild() string {
	return m.GuildID
}

type guildMemberRemoveHandler struct {
	f func(*Fetcher, GuildMemberRemove) error
}

func (e guildMemberRemoveHandler) name() string {
	return "GUILD_MEMBER_REMOVE"
}

func (e guildMemberRemoveHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildMemberRemove))
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

func (m GuildMembersChunk) guild() string {
	return m.GuildID
}

type guildMembersChunkHandler struct {
	f func(*Fetcher, GuildMembersChunk) error
}

func (e guildMembersChunkHandler) name() string {
	return "GUILD_MEMBERS_CHUNK"
}

func (e guildMembersChunkHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildMembersChunk))
}

// GuildRoleCreate is received when a new role was created.
type GuildRoleCreate struct {
	GuildID string `json:"guild_id"` // ID of the guild
	Role    Role   `json:"role"`     // Role that was created
}

func (m GuildRoleCreate) guild() string {
	return m.GuildID
}

type guildRoleCreateHandler struct {
	f func(*Fetcher, GuildRoleCreate) error
}

func (e guildRoleCreateHandler) name() string {
	return "GUILD_ROLE_CREATE"
}

func (e guildRoleCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildRoleCreate))
}

// GuildRoleUpdate is received when a role was updated.
type GuildRoleUpdate struct {
	GuildID string `json:"guild_id"` // ID of the guild
	Role    Role   `json:"role"`     // Role that was updated
}

func (m GuildRoleUpdate) guild() string {
	return m.GuildID
}

type guildRoleUpdateHandler struct {
	f func(*Fetcher, GuildRoleUpdate) error
}

func (e guildRoleUpdateHandler) name() string {
	return "GUILD_ROLE_UPDATE"
}

func (e guildRoleUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildRoleUpdate))
}

// GuildRoleDelete is received when a role was deleted.
type GuildRoleDelete struct {
	GuildID string `json:"guild_id"` // ID of the guild
	RoleID  string `json:"role_id"`  // ID of the role
}

func (m GuildRoleDelete) guild() string {
	return m.GuildID
}

type guildRoleDeleteHandler struct {
	f func(*Fetcher, GuildRoleDelete) error
}

func (e guildRoleDeleteHandler) name() string {
	return "GUILD_ROLE_DELETE"
}

func (e guildRoleDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildRoleDelete))
}

// GuildScheduledEventCreate is received when a new scheduled event is created.
type GuildScheduledEventCreate struct {
	GuildScheduledEvent
}

func (m GuildScheduledEventCreate) guild() string {
	return m.GuildID
}

type guildScheduledEventCreateHandler struct {
	f func(*Fetcher, GuildScheduledEventCreate) error
}

func (e guildScheduledEventCreateHandler) name() string {
	return "GUILD_SCHEDULED_EVENT_CREATE"
}

func (e guildScheduledEventCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildScheduledEventCreate))
}

// GuildScheduledEventUpdate is received when a new scheduled event is updated.
type GuildScheduledEventUpdate struct {
	GuildScheduledEvent
}

func (m GuildScheduledEventUpdate) guild() string {
	return m.GuildID
}

type guildScheduledEventUpdateHandler struct {
	f func(*Fetcher, GuildScheduledEventUpdate) error
}

func (e guildScheduledEventUpdateHandler) name() string {
	return "GUILD_SCHEDULED_EVENT_UPDATE"
}

func (e guildScheduledEventUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildScheduledEventUpdate))
}

// GuildScheduledEventDelete is received when a new scheduled event is deleted.
type GuildScheduledEventDelete struct {
	GuildScheduledEvent
}

func (m GuildScheduledEventDelete) guild() string {
	return m.GuildID
}

type guildScheduledEventDeleteHandler struct {
	f func(*Fetcher, GuildScheduledEventDelete) error
}

func (e guildScheduledEventDeleteHandler) name() string {
	return "GUILD_SCHEDULED_EVENT_DELETE"
}

func (e guildScheduledEventDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildScheduledEventDelete))
}

// GuildScheduledEventUserAddEvent is received when a user has subscribed to a guild scheduled event.
type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

func (m GuildScheduledEventUserAddEvent) guild() string {
	return m.GuildID
}

type guildScheduledEventUserAddEventHandler struct {
	f func(*Fetcher, GuildScheduledEventUserAddEvent) error
}

func (e guildScheduledEventUserAddEventHandler) name() string {
	return "GUILD_SCHEDULED_EVENT_USER_ADD_EVENT"
}

func (e guildScheduledEventUserAddEventHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildScheduledEventUserAddEvent))
}

// GuildScheduledEventUserRemoveEvent is received when a user has unsubscribed from a guild scheduled event.
type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

func (m GuildScheduledEventUserRemoveEvent) guild() string {
	return m.GuildID
}

type guildScheduledEventUserRemoveEventHandler struct {
	f func(*Fetcher, GuildScheduledEventUserRemoveEvent) error
}

func (e guildScheduledEventUserRemoveEventHandler) name() string {
	return "GUILD_SCHEDULED_EVENT_USER_REMOVE_EVENT"
}

func (e guildScheduledEventUserRemoveEventHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(GuildScheduledEventUserRemoveEvent))
}

// IntegrationCreate is received when a new integration is created.
type IntegrationCreate struct {
	Integration

	GuildID string `json:"guild_id"`
}

func (m IntegrationCreate) guild() string {
	return m.GuildID
}

type integrationCreateHandler struct {
	f func(*Fetcher, IntegrationCreate) error
}

func (e integrationCreateHandler) name() string {
	return "INTEGRATION_CREATE"
}

func (e integrationCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(IntegrationCreate))
}

// IntegrationUpdate is received when a new integration is updated.
type IntegrationUpdate struct {
	Integration

	GuildID string `json:"guild_id"`
}

func (m IntegrationUpdate) guild() string {
	return m.GuildID
}

type integrationUpdateHandler struct {
	f func(*Fetcher, IntegrationUpdate) error
}

func (e integrationUpdateHandler) name() string {
	return "INTEGRATION_UPDATE"
}

func (e integrationUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(IntegrationUpdate))
}

// IntegrationDelete is sent when an integration is deleted.
type IntegrationDelete struct {
	ID            string  `json:"id"`                       // Integration ID.
	GuildID       string  `json:"guild_id"`                 // ID of the guild.
	ApplicationID *string `json:"application_id,omitempty"` // ID of the bot/OAuth2 application for this Discord integration.
}

func (m IntegrationDelete) guild() string {
	return m.GuildID
}

type integrationDeleteHandler struct {
	f func(*Fetcher, IntegrationDelete) error
}

func (e integrationDeleteHandler) name() string {
	return "INTEGRATION_DELETE"
}

func (e integrationDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(IntegrationDelete))
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

func (m InviteCreate) guild() string {
	return m.GuildID
}

type inviteCreateHandler struct {
	f func(*Fetcher, InviteCreate) error
}

func (e inviteCreateHandler) name() string {
	return "INVITE_CREATE"
}

func (e inviteCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(InviteCreate))
}

// InviteDelete is sent when an invite is deleted.
type InviteDelete struct {
	ChannelID string `json:"channel_id"`         // Channel of the invite.
	GuildID   string `json:"guild_id,omitempty"` // Guild of the invite.
	Code      string `json:"code"`               // Unique invite code.
}

func (m InviteDelete) guild() string {
	return m.GuildID
}

type inviteDeleteHandler struct {
	f func(*Fetcher, InviteDelete) error
}

func (e inviteDeleteHandler) name() string {
	return "INVITE_DELETE"
}

func (e inviteDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(InviteDelete))
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

func (m MessageCreate) guild() string {
	return m.GuildID
}

type messageCreateHandler struct {
	f func(*Fetcher, MessageCreate) error
}

func (e messageCreateHandler) name() string {
	return "MESSAGE_CREATE"
}

func (e messageCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageCreate))
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

func (m MessageUpdate) guild() string {
	return m.GuildID
}

type messageUpdateHandler struct {
	f func(*Fetcher, MessageUpdate) error
}

func (e messageUpdateHandler) name() string {
	return "MESSAGE_UPDATE"
}

func (e messageUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageUpdate))
}

// MessageDelete is sent with a message delete event.
type MessageDelete struct {
	ID        string `json:"id"`                 // ID of the message.
	ChannelID string `json:"channel_id"`         // ID of the channel.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (m MessageDelete) guild() string {
	return m.GuildID
}

type messageDeleteHandler struct {
	f func(*Fetcher, MessageDelete) error
}

func (e messageDeleteHandler) name() string {
	return "MESSAGE_DELETE"
}

func (e messageDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageDelete))
}

// MessageDeleteBulk is sent with a bulk message delete event.
type MessageDeleteBulk struct {
	IDs       []string `json:"ids"`                // IDs of the deleted messages.
	ChannelID string   `json:"channel_id"`         // ID of the channel.
	GuildID   string   `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (m MessageDeleteBulk) guild() string {
	return m.GuildID
}

type messageDeleteBulkHandler struct {
	f func(*Fetcher, MessageDeleteBulk) error
}

func (e messageDeleteBulkHandler) name() string {
	return "MESSAGE_DELETE_BULK"
}

func (e messageDeleteBulkHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageDeleteBulk))
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

func (m MessageReactionAdd) guild() string {
	return m.GuildID
}

type messageReactionAddHandler struct {
	f func(*Fetcher, MessageReactionAdd) error
}

func (e messageReactionAddHandler) name() string {
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

func (m MessageReactionRemove) guild() string {
	return m.GuildID
}

type messageReactionRemoveHandler struct {
	f func(*Fetcher, MessageReactionRemove) error
}

func (e messageReactionRemoveHandler) name() string {
	return "MESSAGE_REACTION_REMOVE"
}

func (e messageReactionRemoveHandler) run(fetcher *Fetcher, ev any) error {
	fmt.Println("real type", reflect.TypeOf(ev))
	return e.f(fetcher, ev.(MessageReactionRemove))
}

// MessageReactionRemoveAll is sent when a user explicitly removes all reactions from a message.
type MessageReactionRemoveAll struct {
	ChannelID string `json:"channel_id"`         // ID of the channel.
	MessageID string `json:"message_id"`         // ID of the message.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
}

func (m MessageReactionRemoveAll) guild() string {
	return m.GuildID
}

type messageReactionRemoveAllHandler struct {
	f func(*Fetcher, MessageReactionRemoveAll) error
}

func (e messageReactionRemoveAllHandler) name() string {
	return "MESSAGE_REACTION_REMOVE_ALL"
}

func (e messageReactionRemoveAllHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageReactionRemoveAll))
}

// MessageReactionRemoveEmoji is sent when a user removes a reaction with a specific emoji from a message.
type MessageReactionRemoveEmoji struct {
	ChannelID string `json:"channel_id"`         // ID of the channel.
	GuildID   string `json:"guild_id,omitempty"` // ID of the guild (optional).
	MessageID string `json:"message_id"`         // ID of the message.
	Emoji     Emoji  `json:"emoji"`              // Emoji that was removed.
}

func (m MessageReactionRemoveEmoji) guild() string {
	return m.GuildID
}

type messageReactionRemoveEmojiHandler struct {
	f func(*Fetcher, MessageReactionRemoveEmoji) error
}

func (e messageReactionRemoveEmojiHandler) name() string {
	return "MESSAGE_REACTION_REMOVE_EMOJI"
}

func (e messageReactionRemoveEmojiHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(MessageReactionRemoveEmoji))
}

// PresenceUpdate is sent when a user's presence in a guild is updated.
type PresenceUpdate struct {
	User         User         `json:"user"`          // User whose presence is being updated
	GuildID      string       `json:"guild_id"`      // ID of the guild
	Status       string       `json:"status"`        // Either "idle", "dnd", "online", or "offline"
	Activities   []Activity   `json:"activities"`    // User's current activities
	ClientStatus ClientStatus `json:"client_status"` // User's platform-dependent status
}

func (m PresenceUpdate) guild() string {
	return m.GuildID
}

type presenceUpdateHandler struct {
	f func(*Fetcher, PresenceUpdate) error
}

func (e presenceUpdateHandler) name() string {
	return "PRESENCE_UPDATE"
}

func (e presenceUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(PresenceUpdate))
}

// StageInstanceCreate is received when a new stage instance is created.
type StageInstanceCreate struct {
	StageInstance
}

func (m StageInstanceCreate) guild() string {
	return m.GuildID
}

type stageInstanceCreateHandler struct {
	f func(*Fetcher, StageInstanceCreate) error
}

func (e stageInstanceCreateHandler) name() string {
	return "STAGE_INSTANCE_CREATE"
}

func (e stageInstanceCreateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(StageInstanceCreate))
}

// StageInstanceCreate is received when a new stage instance is updated.
type StageInstanceUpdate struct {
	StageInstance
}

func (m StageInstanceUpdate) guild() string {
	return m.GuildID
}

type stageInstanceUpdateHandler struct {
	f func(*Fetcher, StageInstanceUpdate) error
}

func (e stageInstanceUpdateHandler) name() string {
	return "STAGE_INSTANCE_UPDATE"
}

func (e stageInstanceUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(StageInstanceUpdate))
}

// StageInstanceDelete is received when an stage instance is deleted.
type StageInstanceDelete struct {
	StageInstance
}

func (m StageInstanceDelete) guild() string {
	return m.GuildID
}

type stageInstanceDeleteHandler struct {
	f func(*Fetcher, StageInstanceDelete) error
}

func (e stageInstanceDeleteHandler) name() string {
	return "STAGE_INSTANCE_DELETE"
}

func (e stageInstanceDeleteHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(StageInstanceDelete))
}

// TypingStart is received when a user starts typing in a channel.
type TypingStart struct {
	ChannelID string       `json:"channel_id"`         // ID of the channel
	GuildID   *string      `json:"guild_id,omitempty"` // ID of the guild.
	UserID    string       `json:"user_id"`            // ID of the user
	Timestamp int          `json:"timestamp"`          // Unix time (in seconds) of when the user started typing
	Member    *GuildMember `json:"member,omitempty"`   // Member who started typing if this happened in a guild
}

func (m TypingStart) guild() string {
	if m.GuildID == nil {
		return ""
	}

	return *m.GuildID
}

type typingStartHandler struct {
	f func(*Fetcher, TypingStart) error
}

func (e typingStartHandler) name() string {
	return "TYPING_START"
}

func (e typingStartHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(TypingStart))
}

// UserUpdate is received when a user is updated.
type UserUpdate struct {
	User
}

func (m UserUpdate) guild() string {
	return ""
}

type userUpdateHandler struct {
	f func(*Fetcher, UserUpdate) error
}

func (e userUpdateHandler) name() string {
	return "USER_UPDATE"
}

func (e userUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(UserUpdate))
}

// VoiceStateUpdate is received when someone joins/leaves/moves voice channels.
type VoiceStateUpdate struct {
	VoiceState
}

func (m VoiceStateUpdate) guild() string {
	if m.VoiceState.GuildID == nil {
		return ""
	}

	return *m.VoiceState.GuildID
}

type voiceStateUpdateHandler struct {
	f func(*Fetcher, VoiceStateUpdate) error
}

func (e voiceStateUpdateHandler) name() string {
	return "VOICE_STATE_UPDATE"
}

func (e voiceStateUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(VoiceStateUpdate))
}

// VoiceServerUpdate is received when a guild's voice server is updated. This is sent when initially connecting to voice, and when the current voice instance fails over to a new server.
// A null endpoint means that the voice server allocated has gone away and is trying to be reallocated. You should attempt to disconnect from the currently connected voice server, and not attempt to reconnect until a new voice server is allocated.
type VoiceServerUpdate struct {
	Token    string  `json:"token"`              // Voice connection token
	GuildID  string  `json:"guild_id"`           // Guild this voice server update is for
	Endpoint *string `json:"endpoint,omitempty"` // Voice server host
}

func (m VoiceServerUpdate) guild() string {
	return m.GuildID
}

type voiceServerUpdateHandler struct {
	f func(*Fetcher, VoiceServerUpdate) error
}

func (e voiceServerUpdateHandler) name() string {
	return "VOICE_SERVER_UPDATE"
}

func (e voiceServerUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(VoiceServerUpdate))
}

// WebhooksUpdate is received when a guild channel's webhook is created, updated, or deleted.
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`   // ID of the guild
	ChannelID string `json:"channel_id"` // ID of the channel
}

func (m WebhooksUpdate) guild() string {
	return m.GuildID
}

type webhooksUpdateHandler struct {
	f func(*Fetcher, WebhooksUpdate) error
}

func (e webhooksUpdateHandler) name() string {
	return "WEBHOOKS_UPDATE"
}

func (e webhooksUpdateHandler) run(fetcher *Fetcher, ev any) error {
	return e.f(fetcher, ev.(WebhooksUpdate))
}

type eventHandler interface {
	run(*Fetcher, any) error
	name() string
}

func eventHandlerFromInterface(iface any) (eventHandler, error) {
	switch v := iface.(type) {
	case func(*Fetcher, ApplicationCommandPermissionsUpdate) error:
		return applicationCommandPermissionsUpdateHandler{f: v}, nil

	case func(*Fetcher, AutoModerationRuleCreate) error:
		return autoModerationRuleCreateHandler{f: v}, nil

	case func(*Fetcher, AutoModerationRuleUpdate) error:
		return autoModerationRuleUpdateHandler{f: v}, nil

	case func(*Fetcher, AutoModerationRuleDelete) error:
		return autoModerationRuleDeleteHandler{f: v}, nil

	case func(*Fetcher, ChannelCreate) error:
		return channelCreateHandler{f: v}, nil

	case func(*Fetcher, ChannelUpdate) error:
		return channelUpdateHandler{f: v}, nil

	case func(*Fetcher, ChannelDelete) error:
		return channelDeleteHandler{f: v}, nil

	case func(*Fetcher, ChannelPinsUpdate) error:
		return channelPinsUpdateHandler{f: v}, nil

	case func(*Fetcher, ThreadCreate) error:
		return threadCreateHandler{f: v}, nil

	case func(*Fetcher, ThreadUpdate) error:
		return threadUpdateHandler{f: v}, nil

	case func(*Fetcher, ThreadDelete) error:
		return threadDeleteHandler{f: v}, nil

	case func(*Fetcher, ThreadListSync) error:
		return threadListSyncHandler{f: v}, nil

	case func(*Fetcher, ThreadMemberUpdate) error:
		return threadMemberUpdateHandler{f: v}, nil

	case func(*Fetcher, ThreadMembersUpdate) error:
		return threadMembersUpdateHandler{f: v}, nil

	case func(*Fetcher, EntitlementCreate) error:
		return entitlementCreateHandler{f: v}, nil

	case func(*Fetcher, EntitlementUpdate) error:
		return entitlementUpdateHandler{f: v}, nil

	case func(*Fetcher, EntitlementDelete) error:
		return entitlementDeleteHandler{f: v}, nil

	case func(*Fetcher, GuildCreate) error:
		return guildCreateHandler{f: v}, nil

	case func(*Fetcher, GuildUpdate) error:
		return guildUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildDelete) error:
		return guildDeleteHandler{f: v}, nil

	case func(*Fetcher, GuildBanAdd) error:
		return guildBanAddHandler{f: v}, nil

	case func(*Fetcher, GuildBanRemove) error:
		return guildBanRemoveHandler{f: v}, nil

	case func(*Fetcher, GuildEmojisUpdate) error:
		return guildEmojisUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildStickersUpdate) error:
		return guildStickersUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildMemberAdd) error:
		return guildMemberAddHandler{f: v}, nil

	case func(*Fetcher, GuildMemberUpdate) error:
		return guildMemberUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildMemberRemove) error:
		return guildMemberRemoveHandler{f: v}, nil

	case func(*Fetcher, GuildMembersChunk) error:
		return guildMembersChunkHandler{f: v}, nil

	case func(*Fetcher, GuildRoleCreate) error:
		return guildRoleCreateHandler{f: v}, nil

	case func(*Fetcher, GuildRoleUpdate) error:
		return guildRoleUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildRoleDelete) error:
		return guildRoleDeleteHandler{f: v}, nil

	case func(*Fetcher, GuildScheduledEventCreate) error:
		return guildScheduledEventCreateHandler{f: v}, nil

	case func(*Fetcher, GuildScheduledEventUpdate) error:
		return guildScheduledEventUpdateHandler{f: v}, nil

	case func(*Fetcher, GuildScheduledEventDelete) error:
		return guildScheduledEventDeleteHandler{f: v}, nil

	case func(*Fetcher, GuildScheduledEventUserAddEvent) error:
		return guildScheduledEventUserAddEventHandler{f: v}, nil

	case func(*Fetcher, GuildScheduledEventUserRemoveEvent) error:
		return guildScheduledEventUserRemoveEventHandler{f: v}, nil

	case func(*Fetcher, IntegrationCreate) error:
		return integrationCreateHandler{f: v}, nil

	case func(*Fetcher, IntegrationUpdate) error:
		return integrationUpdateHandler{f: v}, nil

	case func(*Fetcher, IntegrationDelete) error:
		return integrationDeleteHandler{f: v}, nil

	case func(*Fetcher, InviteCreate) error:
		return inviteCreateHandler{f: v}, nil

	case func(*Fetcher, InviteDelete) error:
		return inviteDeleteHandler{f: v}, nil

	case func(*Fetcher, MessageCreate) error:
		return messageCreateHandler{f: v}, nil

	case func(*Fetcher, MessageUpdate) error:
		return messageUpdateHandler{f: v}, nil

	case func(*Fetcher, MessageDelete) error:
		return messageDeleteHandler{f: v}, nil

	case func(*Fetcher, MessageDeleteBulk) error:
		return messageDeleteBulkHandler{f: v}, nil

	case func(*Fetcher, MessageReactionAdd) error:
		return messageReactionAddHandler{f: v}, nil

	case func(*Fetcher, MessageReactionRemove) error:
		return messageReactionRemoveHandler{f: v}, nil

	case func(*Fetcher, MessageReactionRemoveAll) error:
		return messageReactionRemoveAllHandler{f: v}, nil

	case func(*Fetcher, MessageReactionRemoveEmoji) error:
		return messageReactionRemoveEmojiHandler{f: v}, nil

	case func(*Fetcher, PresenceUpdate) error:
		return presenceUpdateHandler{f: v}, nil

	case func(*Fetcher, StageInstanceCreate) error:
		return stageInstanceCreateHandler{f: v}, nil

	case func(*Fetcher, StageInstanceUpdate) error:
		return stageInstanceUpdateHandler{f: v}, nil

	case func(*Fetcher, StageInstanceDelete) error:
		return stageInstanceDeleteHandler{f: v}, nil

	case func(*Fetcher, TypingStart) error:
		return typingStartHandler{f: v}, nil

	case func(*Fetcher, UserUpdate) error:
		return userUpdateHandler{f: v}, nil

	case func(*Fetcher, VoiceStateUpdate) error:
		return voiceStateUpdateHandler{f: v}, nil

	case func(*Fetcher, VoiceServerUpdate) error:
		return voiceServerUpdateHandler{f: v}, nil

	case func(*Fetcher, WebhooksUpdate) error:
		return webhooksUpdateHandler{f: v}, nil
	default:
		return nil, fmt.Errorf("unknown event")
	}
}
