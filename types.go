package godiscord

import (
	"encoding/json"
	"fmt"
	"time"
)

// TODO: These don't belong to the events package. Can probably create a types package, but that's for later when the api is a bit more covered.

type Guild struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Icon            *string `json:"icon,omitempty"`
	IconHash        *string `json:"icon_hash,omitempty"`
	Splash          *string `json:"splash,omitempty"`
	DiscoverySplash *string `json:"discovery_splash,omitempty"`
	Owner           *bool   `json:"owner,omitempty"`
	OwnerID         string  `json:"owner_id"`
	// Permissions is the old format of bitwise int. This field should only be used when sending a permission.
	// When v6 is removed, this will most likely be used for both sending and receiving
	Permissions int `json:"permissions"`
	// Permissions is the new format, where it's send as a string. This field should only be used when receiving a permission.
	// Will most likely be deprecated when v6 is removed.
	PermissionsNew string `json:"permissions_new"`
	AFKChannelID   string `json:"afk_channel_id"`
	// AFKTimeout is in seconds.
	AFKTimeout                  int                             `json:"afk_timeout"`
	WidgetEnabled               *bool                           `json:"widget_enabled,omitempty"`
	WidgetChannelID             *string                         `json:"widget_channel_id,omitempty"`
	VerificationLevel           VerificationLevel               `json:"verification_level"`
	DefaultMessageNotifications DefaultMessageNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel      `json:"explicit_content_filter"`
	Roles                       []Role                          `json:"roles"`
	Emojis                      []Emoji                         `json:"emojis"`
	Features                    []GuildFeature                  `json:"features"`
	MfaLevel                    MFALevel                        `json:"mfa_level"`
	ApplicationID               *string                         `json:"application_id,omitempty"`
	SystemChannelID             *string                         `json:"system_channel_id,omitempty"`
	// SystemChannelFlags are flags combined as a bitfield.
	SystemChannelFlags        int            `json:"system_channel_flags"`
	RulesIChannelID           string         `json:"rules_channel_id"`
	MaxPresences              *int           `json:"max_presences,omitempty"`
	MaxMembers                *int           `json:"max_members,omitempty"`
	VanityURLCode             *string        `json:"vanity_url_code,omitempty"`
	Description               *string        `json:"description,omitempty"`
	Banner                    *string        `json:"banner,omitempty"`
	PremiumTier               int            `json:"premium_tier"`
	PremiumSubscriptionCount  *int           `json:"premium_subscription_count,omitempty"`
	PreferredLocale           string         `json:"preferred_locale"`
	PublicUpdatesChannelID    *string        `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers      *int           `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount    *int           `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount  *int           `json:"approximate_presence_count,omitempty"`
	WelcomeScreen             WelcomeScreen  `json:"welcome_screen"`
	NSFWLevel                 GuildNSFWLevel `json:"nsfw_level"`
	Stickers                  []Sticker      `json:"stickers"`
	PremiumProgressBarEnabled bool           `json:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID     *string        `json:"safety_alerts_channel_id,omitempty,omitempty"`
	Unavailable               bool           `json:"unavailable"`
}

type VerificationLevel uint32

const (
	// VerificationLevelNone is unrestricted.
	VerificationLevelNone VerificationLevel = iota
	// VerificationLevelLow enforces that users must have verified email on account.
	VerificationLevelLow
	// VerificationLevelMedium enforces that users must be registered on Discord for longer than 5 minutes.
	VerificationLevelMedium
	// VerificationLevelHigh enforces that users must be registered on Discord for longer than 10 minutes.
	VerificationLevelHigh
	// VerificationLevelVeryHigh enforces that users must have a verified phone number.
	VerificationLevelVeryHigh
)

type DefaultMessageNotificationLevel uint32

const (
	// DefaultMessageNotificationLevelAllMessages represents the level where
	// members will receive notifications for all messages by default.
	DefaultMessageNotificationLevelAllMessages DefaultMessageNotificationLevel = iota

	// DefaultMessageNotificationLevelOnlyMentions represents the level where
	// members will receive notifications only for messages that @mention them by default.
	DefaultMessageNotificationLevelOnlyMentions
)

type ExplicitContentFilterLevel uint32

const (
	// ExplicitContentFilterLevelDisabled represents the level where media content will not be scanned.
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota

	// ExplicitContentFilterLevelMembersWithoutRoles represents the level where media content sent by members without roles will be scanned.
	ExplicitContentFilterLevelMembersWithoutRoles

	// ExplicitContentFilterLevelAllMembers represents the level where media content sent by all members will be scanned.
	ExplicitContentFilterLevelAllMembers
)

type MFALevel uint32

const (
	// MFALevelNone represents the level where the guild has no MFA/2FA requirement for moderation actions.
	MFALevelNone MFALevel = iota

	// MFALevelElevated represents the level where the guild has a 2FA requirement for moderation actions.
	MFALevelElevated
)

type GuildNSFWLevel uint32

const (
	// GuildNSFWLevelDefault represents the default NSFW level where the guild has no explicit content restrictions.
	GuildNSFWLevelDefault GuildNSFWLevel = iota

	// GuildNSFWLevelExplicit represents the NSFW level where explicit content is allowed.
	GuildNSFWLevelExplicit

	// GuildNSFWLevelSafe represents the NSFW level where explicit content is allowed but limited.
	GuildNSFWLevelSafe

	// GuildNSFWLevelAgeRestricted represents the NSFW level where explicit content is allowed only for members who are age-restricted.
	GuildNSFWLevelAgeRestricted
)

type GuildFeature string

const (
	// GuildFeatureAnimatedBanner indicates that the guild has access to set an animated guild banner image.
	GuildFeatureAnimatedBanner GuildFeature = "ANIMATED_BANNER"

	// GuildFeatureAnimatedIcon indicates that the guild has access to set an animated guild icon.
	GuildFeatureAnimatedIcon GuildFeature = "ANIMATED_ICON"

	// GuildFeatureApplicationCommandPermissionsV2 indicates that the guild is using the old permissions configuration behavior.
	GuildFeatureApplicationCommandPermissionsV2 GuildFeature = "APPLICATION_COMMAND_PERMISSIONS_V2"

	// GuildFeatureAutoModeration indicates that the guild has set up auto moderation rules.
	GuildFeatureAutoModeration GuildFeature = "AUTO_MODERATION"

	// GuildFeatureBanner indicates that the guild has access to set a guild banner image.
	GuildFeatureBanner GuildFeature = "BANNER"

	// GuildFeatureCommunity indicates that the guild can enable the welcome screen, Membership Screening, stage channels and discovery, and receives community updates.
	GuildFeatureCommunity GuildFeature = "COMMUNITY"

	// GuildFeatureCreatorMonetizableProvisional indicates that the guild has enabled monetization.
	GuildFeatureCreatorMonetizableProvisional GuildFeature = "CREATOR_MONETIZABLE_PROVISIONAL"

	// GuildFeatureCreatorStorePage indicates that the guild has enabled the role subscription promo page.
	GuildFeatureCreatorStorePage GuildFeature = "CREATOR_STORE_PAGE"

	// GuildFeatureDeveloperSupportServer indicates that the guild has been set as a support server on the App Directory.
	GuildFeatureDeveloperSupportServer GuildFeature = "DEVELOPER_SUPPORT_SERVER"

	// GuildFeatureDiscoverable indicates that the guild is able to be discovered in the directory.
	GuildFeatureDiscoverable GuildFeature = "DISCOVERABLE"

	// GuildFeatureFeaturable indicates that the guild is able to be featured in the directory.
	GuildFeatureFeaturable GuildFeature = "FEATURABLE"

	// GuildFeatureInvitesDisabled indicates that the guild has paused invites, preventing new users from joining.
	GuildFeatureInvitesDisabled GuildFeature = "INVITES_DISABLED"

	// GuildFeatureInviteSplash indicates that the guild has access to set an invite splash background.
	GuildFeatureInviteSplash GuildFeature = "INVITE_SPLASH"

	// GuildFeatureMemberVerificationGateEnabled indicates that the guild has enabled Membership Screening.
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"

	// GuildFeatureMoreStickers indicates that the guild has increased custom sticker slots.
	GuildFeatureMoreStickers GuildFeature = "MORE_STICKERS"

	// GuildFeatureNews indicates that the guild has access to create announcement channels.
	GuildFeatureNews GuildFeature = "NEWS"

	// GuildFeaturePartnered indicates that the guild is partnered.
	GuildFeaturePartnered GuildFeature = "PARTNERED"

	// GuildFeaturePreviewEnabled indicates that the guild can be previewed before joining via Membership Screening or the directory.
	GuildFeaturePreviewEnabled GuildFeature = "PREVIEW_ENABLED"

	// GuildFeatureRaidAlertsDisabled indicates that the guild has disabled alerts for join raids in the configured safety alerts channel.
	GuildFeatureRaidAlertsDisabled GuildFeature = "RAID_ALERTS_DISABLED"

	// GuildFeatureRoleIcons indicates that the guild is able to set role icons.
	GuildFeatureRoleIcons GuildFeature = "ROLE_ICONS"

	// GuildFeatureRoleSubscriptionsAvailableForPurchase indicates that the guild has role subscriptions that can be purchased.
	GuildFeatureRoleSubscriptionsAvailableForPurchase GuildFeature = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"

	// GuildFeatureRoleSubscriptionsEnabled indicates that the guild has enabled role subscriptions.
	GuildFeatureRoleSubscriptionsEnabled GuildFeature = "ROLE_SUBSCRIPTIONS_ENABLED"

	// GuildFeatureTicketedEventsEnabled indicates that the guild has enabled ticketed events.
	GuildFeatureTicketedEventsEnabled GuildFeature = "TICKETED_EVENTS_ENABLED"

	// GuildFeatureVanityURL indicates that the guild has access to set a vanity URL.
	GuildFeatureVanityURL GuildFeature = "VANITY_URL"

	// GuildFeatureVerified indicates that the guild is verified.
	GuildFeatureVerified GuildFeature = "VERIFIED"

	// GuildFeatureVIPRegions indicates that the guild has access to set 384kbps bitrate in voice (previously VIP voice servers).
	GuildFeatureVIPRegions GuildFeature = "VIP_REGIONS"

	// GuildFeatureWelcomeScreenEnabled indicates that the guild has enabled the welcome screen.
	GuildFeatureWelcomeScreenEnabled GuildFeature = "WELCOME_SCREEN_ENABLED"
)

type SystemChannelFlag uint32

const (
	// SystemChannelFlagSuppressJoinNotifications suppresses member join notifications.
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlag = 1 << iota

	// SystemChannelFlagSuppressPremiumSubscriptions suppresses server boost notifications.
	SystemChannelFlagSuppressPremiumSubscriptions

	// SystemChannelFlagSuppressGuildReminderNotifications suppresses server setup tips.
	SystemChannelFlagSuppressGuildReminderNotifications

	// SystemChannelFlagSuppressJoinNotificationReplies hides member join sticker reply buttons.
	SystemChannelFlagSuppressJoinNotificationReplies

	// SystemChannelFlagSuppressRoleSubscriptionPurchaseNotifications suppresses role subscription purchase and renewal notifications.
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotifications

	// SystemChannelFlagSuppressRoleSubscriptionPurchaseNotificationReplies hides role subscription sticker reply buttons.
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotificationReplies
)

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Color is integer representation of hex color.
	// Roles without colors (color == 0) do not count towards the final computed color in the user list.
	Color int `json:"color"`
	// Hoist is if this user is pinned in the user listing.
	Hoist        bool     `json:"hoist"`
	Icon         *string  `json:"icon,omitempty"`
	UnicodeEmoji *string  `json:"icon_hash,omitempty"`
	Position     int      `json:"position"`
	Permissions  string   `json:"permissions"`
	Managed      bool     `json:"managed"`
	Mentionable  bool     `json:"mentionable"`
	Tags         *RoleTag `json:"tags,omitempty"`
	// RoleFlags are flags combined as a bitfield.
	RoleFlags int `json:"flags"`
}

type RoleTag struct {
	BotID                 *string `json:"bot_id,omitempty"`
	IntegrationID         *string `json:"integration_id,omitempty"`
	PremiumSubscriber     *bool   `json:"premium_subscriber,omitempty"`
	SubscriptionListingID *string `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  *bool   `json:"available_for_purchase,omitempty"`
	GuildConnections      *bool   `json:"guild_connections,omitempty"`
}

type RoleFlag uint32

const (
	// RoleFlagInPrompt represents the flag indicating that the role can be selected by members in an onboarding prompt.
	RoleFlagInPrompt RoleFlag = 1 << 0
)

type Emoji struct {
	ID *string `json:"id,omitempty"`
	// Name can be null only in reaction emoji objects.
	Name          *string  `json:"name,omitempty"`
	Roles         []string `json:"roles"`
	User          *User    `json:"user,omitempty"`
	RequireColons *bool    `json:"require_colons,omitempty"`
	Managed       *bool    `json:"managed,omitempty"`
	Animated      *bool    `json:"animated,omitempty"`
	Available     *bool    `json:"available,omitempty"`
}

// VoiceState represents a user's voice connection status.
type VoiceState struct {
	GuildID                 *string      `json:"guild_id,omitempty"`                   // The guild id this voice state is for.
	ChannelID               *string      `json:"channel_id,omitempty"`                 // The channel id this user is connected to.
	UserID                  string       `json:"user_id"`                              // The user id this voice state is for.
	Member                  *GuildMember `json:"member,omitempty"`                     // The guild member this voice state is for.
	SessionID               string       `json:"session_id"`                           // The session id for this voice state.
	Deaf                    bool         `json:"deaf"`                                 // Whether this user is deafened by the server.
	Mute                    bool         `json:"mute"`                                 // Whether this user is muted by the server.
	SelfDeaf                bool         `json:"self_deaf"`                            // Whether this user is locally deafened.
	SelfMute                bool         `json:"self_mute"`                            // Whether this user is locally muted.
	SelfStream              *bool        `json:"self_stream,omitempty"`                // Whether this user is streaming using "Go Live".
	SelfVideo               bool         `json:"self_video"`                           // Whether this user's camera is enabled.
	Suppress                bool         `json:"suppress"`                             // Whether this user's permission to speak is denied.
	RequestToSpeakTimestamp *time.Time   `json:"request_to_speak_timestamp,omitempty"` // The time at which the user requested to speak.
}

type GuildMember struct {
	// User can be nil if the user associated with the guild member is deleted or otherwise inaccessible.
	// This situation might occur, for example, if a user deletes their account or is banned from the server.
	User         *User      `json:"user,omitempty"`
	Nick         *string    `json:"nick,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	Roles        []string   `json:"roles"`
	JoinedAt     time.Time  `json:"joined_at"`
	PremiumSince *time.Time `json:"premium_since,omitempty"`
	// Flags are guild member flags represented as a bit set, defaults to 0.
	Flags                      int        `json:"flags"`
	Deaf                       bool       `json:"deaf"`
	Mute                       bool       `json:"mute"`
	Pending                    *bool      `json:"pending,omitempty"`
	Permissions                *int       `json:"permissions,omitempty"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
}

type GuildMemberFlag int

const (
	// GuildMemberFlagDidRejoin indicates that the member has left and rejoined the guild.
	GuildMemberFlagDidRejoin GuildMemberFlag = 1 << 0

	// GuildMemberFlagCompletedOnboarding indicates that the member has completed onboarding.
	GuildMemberFlagCompletedOnboarding GuildMemberFlag = 1 << 1

	// GuildMemberFlagBypassesVerification indicates that the member is exempt from guild verification requirements.
	GuildMemberFlagBypassesVerification GuildMemberFlag = 1 << 2

	// GuildMemberFlagStartedOnboarding indicates that the member has started onboarding.
	GuildMemberFlagStartedOnboarding GuildMemberFlag = 1 << 3
)

type Channel struct {
	ID                   string                `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              *string               `json:"guild_id,omitempty"`
	Position             *int                  `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 *string               `json:"name,omitempty"`
	Topic                *string               `json:"topic,omitempty"`
	NSFW                 bool                  `json:"nsfw"`
	LastMessageID        *string               `json:"last_message_id,omitempty"`
	Bitrate              *int                  `json:"bitrate,omitempty"`
	UserLimit            *int                  `json:"user_limit,omitempty"`
	// RateLimitPerUser also applies to thread creation. Users can send one message and create one thread during each rate_limit_per_user interval.
	RateLimitPerUser *int    `json:"rate_limit_per_user,omitempty"`
	Recipients       []User  `json:"recipients,omitempty"`
	Icon             *string `json:"icon,omitempty"`
	OwnerID          *string `json:"owner_id,omitempty"`
	ApplicationID    *string `json:"application_id,omitempty"`
	Managed          bool    `json:"managed"`
	// ParentID is the parent id of the channel.
	// For guild channels: id of the parent category for a channel (each parent category can contain up to 50 channels).
	// For threads: id of the text channel this thread was created.
	ParentID                      *string             `json:"parent_id,omitempty"`
	LastPinTimestamp              *time.Time          `json:"last_pin_timestamp,omitempty"`
	RTCRegion                     *string             `json:"rtc_region,omitempty"`
	VideoQualityMode              *VideoQualityMode   `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                `json:"message_count,omitempty"`
	MemberCount                   *int                `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata     `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember       `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int                `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *int                `json:"permissions,omitempty"`
	Flags                         *int                `json:"flags,omitempty"`
	TotalMessagesSent             *int                `json:"total_messages_sent,omitempty"`
	AvailableTags                 []ForumTag          `json:"available_tags,omitempty"`
	AppliedTags                   []string            `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReaction    `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *ChannelSortOrder   `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *ChannelForumLayout `json:"default_forum_layout,omitempty"`
}

type ChannelType int

const (
	ChannelTypeGuildText          ChannelType = 0  // a text channel within a server
	ChannelTypeDM                 ChannelType = 1  // a direct message between users
	ChannelTypeGuildVoice         ChannelType = 2  // a voice channel within a server
	ChannelTypeGroupDM            ChannelType = 3  // a direct message between multiple users
	ChannelTypeGuildCategory      ChannelType = 4  // an organizational category that contains up to 50 channels
	ChannelTypeGuildAnnouncement  ChannelType = 5  // a channel that users can follow and crosspost into their own server (formerly news channels)
	ChannelTypeAnnouncementThread ChannelType = 10 // a temporary sub-channel within a GUILD_ANNOUNCEMENT channel
	ChannelTypePublicThread       ChannelType = 11 // a temporary sub-channel within a GUILD_TEXT or GUILD_FORUM channel
	ChannelTypePrivateThread      ChannelType = 12 // a temporary sub-channel within a GUILD_TEXT channel that is only viewable by those invited and those with the MANAGE_THREADS permission
	ChannelTypeGuildStageVoice    ChannelType = 13 // a voice channel for hosting events with an audience
	ChannelTypeGuildDirectory     ChannelType = 14 // the channel in a hub containing the listed servers
	ChannelTypeGuildForum         ChannelType = 15 // a channel that can only contain threads
	ChannelTypeGuildMedia         ChannelType = 16 // a channel that can only contain threads, similar to GUILD_FORUM channels
)

// StageInstance represents a live stage instance.
type StageInstance struct {
	ID                    string                    `json:"id"`                       // The id of this Stage instance.
	GuildID               string                    `json:"guild_id"`                 // The guild id of the associated Stage channel.
	ChannelID             string                    `json:"channel_id"`               // The id of the associated Stage channel.
	Topic                 string                    `json:"topic"`                    // The topic of the Stage instance (1-120 characters).
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level"`            // The privacy level of the Stage instance.
	DiscoverableDisabled  bool                      `json:"discoverable_disabled"`    // Whether or not Stage Discovery is disabled (deprecated).
	GuildScheduledEventID string                    `json:"guild_scheduled_event_id"` // The id of the scheduled event for this Stage instance.
}

// StageInstancePrivacyLevel represents the privacy level of a Stage instance.
type StageInstancePrivacyLevel int

const (
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2 // The Stage instance is visible to only guild members.
)

// VideoQualityMode represents the quality modes for voice channels.
type VideoQualityMode int

const (
	// VideoQualityAuto indicates Discord chooses the quality for optimal performance.
	VideoQualityAuto VideoQualityMode = 1
	// VideoQualityFull indicates 720p quality.
	VideoQualityFull VideoQualityMode = 2
)

// ChannelFlag represents flags associated with channels.
type ChannelFlag int

const (
	// ChannelFlagPinned indicates the thread is pinned to the top of its parent GUILD_FORUM or GUILD_MEDIA channel.
	ChannelFlagPinned ChannelFlag = 1 << 1
	// ChannelFlagRequireTag indicates whether a tag is required to be specified when creating a thread in a GUILD_FORUM or a GUILD_MEDIA channel. Tags are specified in the applied_tags field.
	ChannelFlagRequireTag ChannelFlag = 1 << 4
	// ChannelFlagHideMediaDownloadOptions hides the embedded media download options. Available only for media channels.
	ChannelFlagHideMediaDownloadOptions ChannelFlag = 1 << 15
)

// ChannelSortOrder represents the types of sort orders for forum posts.
type ChannelSortOrder int

const (
	// ChannelSortOrderLatestActivity sorts forum posts by activity.
	ChannelSortOrderLatestActivity ChannelSortOrder = 0
	// ChannelSortOrderCreationDate sorts forum posts by creation time (from most recent to oldest).
	ChannelSortOrderCreationDate ChannelSortOrder = 1
)

// ChannelForumLayout represents the types of layout views for forum channels.
type ChannelForumLayout int

const (
	// ChannelForumLayoutNotSet represents no default set for forum channel layout.
	ChannelForumLayoutNotSet ChannelForumLayout = 0
	// ChannelForumLayoutListView represents displaying posts as a list.
	ChannelForumLayoutListView ChannelForumLayout = 1
	// ChannelForumLayoutGalleryView represents displaying posts as a collection of tiles.
	ChannelForumLayoutGalleryView ChannelForumLayout = 2
)

type Application struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Icon                *string  `json:"icon,omitempty"`
	Description         string   `json:"description"`
	RPCOrigins          []string `json:"rpc_origins,omitempty"`
	BotPublic           bool     `json:"bot_public"`
	BotRequireCodeGrant bool     `json:"bot_require_code_grant"`
	// Bot is a partial user.
	Bot               *User   `json:"bot,omitempty"`
	TermsOfServiceURL *string `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL  *string `json:"privacy_policy_url,omitempty"`
	// Owner is a partial user.
	Owner     *User   `json:"owner,omitempty"`
	VerifyKey string  `json:"verify_key"`
	Team      *Team   `json:"team,omitempty"`
	GuildID   *string `json:"guild_id,omitempty"`
	// Guild is a partial guild.
	Guild                          *Guild         `json:"guild,omitempty"`
	PrimarySkuID                   *string        `json:"primary_sku_id,omitempty"`
	Slug                           *string        `json:"slug,omitempty"`
	CoverImage                     *string        `json:"cover_image,omitempty"`
	Flags                          *int           `json:"flags,omitempty"`
	ApproximateGuildCount          *int           `json:"approximate_guild_count,omitempty"`
	RedirectURIs                   []string       `json:"redirect_uris,omitempty"`
	InteractionsEndpointURL        *string        `json:"interactions_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL *string        `json:"role_connections_verification_url,omitempty"`
	Tags                           []string       `json:"tags,omitempty"`
	InstallParams                  *InstallParams `json:"install_params,omitempty"`
	CustomInstallURL               *string        `json:"custom_install_url,omitempty"`
}

const (
	ApplicationFlagAutoModerationRuleCreateBadge = 1 << 6
	ApplicationFlagGatewayPresence               = 1 << 12
	ApplicationFlagGatewayPresenceLimited        = 1 << 13
	ApplicationFlagGatewayGuildMembers           = 1 << 14
	ApplicationFlagGatewayGuildMembersLimited    = 1 << 15
	ApplicationFlagVerificationPendingGuildLimit = 1 << 16
	ApplicationFlagEmbedded                      = 1 << 17
	ApplicationFlagGatewayMessageContent         = 1 << 18
	ApplicationFlagGatewayMessageContentLimited  = 1 << 19
	ApplicationFlagApplicationCommandBadge       = 1 << 23
)

type Team struct {
	// Members are partial users.
	Members []User `json:"members"`
}

type InstallParams struct {
	Scope []string `json:"scope"`
	// Permissions is the old format of bitwise int. This field should only be used when sending a permission.
	// When v6 is removed, this will most likely be used for both sending and receiving
	Permissions int `json:"permissions"`
	// Permissions is the new format, where it's send as a string. This field should only be used when receiving a permission.
	// Will most likely be deprecated when v6 is removed.
	PermissionsNew string `json:"permissions_new"`
}

type PermissionOverwrite struct {
	ID   string                  `json:"id"`
	Type PermissionOverwriteType `json:"type"`
	// Allow is a permission bit set.
	Allow string `json:"allow"`
	// Deny is a permission bit set.
	Deny string `json:"deny"`
}

type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole   = 0
	PermissionOverwriteTypeMember = 1
)

type ThreadMetadata struct {
	Archived bool `json:"archived"`
	// AutoArchiveDuration is in minutes.
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time  `json:"archive_timestamp"`
	Locked              bool       `json:"locked"`
	Invitable           *bool      `json:"invitable,omitempty"`
	CreateTimestamp     *time.Time `json:"create_timestamp,omitempty"`
}

type ThreadMember struct {
	ID *string `json:"id,omitempty"`
	// UserID is omitted on the member sent within each thread in the GUILD_CREATE event.
	UserID        *string   `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         int       `json:"flags"`
	// Member is omitted on the member sent within each thread in the GUILD_CREATE event.
	// The member field is only present when with_member is set to true when calling List Thread Members or Get Thread Member.
	Member *GuildMember `json:"member,omitempty"`
}

// ForumTag is a forum tag.
// When updating a GUILD_FORUM or a GUILD_MEDIA channel, tag objects in available_tags only require the name field.
type ForumTag struct {
	ID string `json:"id"`
	// Name is 0-20 characters.
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`

	DefaultReaction
}

// DefaultReaction describes the default reaction.
// At most one of emoji_id and emoji_name may be set to a non-null value.
type DefaultReaction struct {
	// EmojiID is the id of a guild's custom emjoi.
	EmojiID *string `json:"emoji_id,omitempty"`
	// EmojiName is the unicode character for the emoji.
	EmojiName *string `json:"emoji_name,omitempty"`
}

type PremiumTier int

const (
	// PremiumTierNone represents the premium tier where the guild has not unlocked any Server Boost perks.
	PremiumTierNone PremiumTier = 0

	// PremiumTier1 represents the premium tier where the guild has unlocked Server Boost level 1 perks.
	PremiumTier1 PremiumTier = 1

	// PremiumTier2 represents the premium tier where the guild has unlocked Server Boost level 2 perks.
	PremiumTier2 PremiumTier = 2

	// PremiumTier3 represents the premium tier where the guild has unlocked Server Boost level 3 perks.
	PremiumTier3 PremiumTier = 3
)

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
	GlobalName    string `json:"global_name"`
}

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           string       `json:"url"`
	CreatedAt     int          `json:"created_at"`
	Timestamps    Timestamps   `json:"timestamps"`
	ApplicationID string       `json:"application_id"`
	Details       string       `json:"details"`
	State         string       `json:"state"`
	Emoji         Emoji        `json:"emoji"`
	Party         Party        `json:"party"`
	Assets        Assets       `json:"assets"`
	Secrets       Secrets      `json:"secrets"`
	Instance      bool         `json:"instance"`
	Flags         int          `json:"flags"`
	Buttons       []string     `json:"buttons"`
}

type ActivityType uint32

const (
	ActivityGame ActivityType = iota
	ActivityStreaming
	ActivityListening
	ActivityWatching
	ActivityCustom
	ActivityCompeting
)

type Timestamps struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Party struct {
	ID   string `json:"id"`
	Size []int  `json:"size"`
}

type Assets struct {
	LargeImage string `json:"large_image"`
	LargeText  string `json:"large_text"`
	SmallImage string `json:"small_image"`
	SmallText  string `json:"small_text"`
}

type Secrets struct {
	Join     string `json:"join"`
	Spectate string `json:"spectate"`
	Match    string `json:"match"`
}

type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}

type WelcomeScreen struct {
	Description     string                 `json:"description"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels"`
}

type WelcomeScreenChannel struct {
	ChannelID   string  `json:"channel_id"`
	Description string  `json:"description"`
	EmojiID     *string `json:"emoji_id,omitempty"`
	EmojiName   *string `json:"emoji_name,omitempty"`
}

type Sticker struct {
	ID string `json:"id"`
	// PackID is comma separated list of keywords is the format used in this field by standard stickers, but this is just a convention.
	// Incidentally the client will always use a name generated from an emoji as the value of this field when creating or modifying a guild sticker.
	PackID      *string           `json:"pack_id,omitempty"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Tags        string            `json:"tags"`
	Asset       *string           `json:"asset,omitempty"`
	Type        StickerType       `json:"type"`
	FormatType  StickerFormatType `json:"format_type"`
	Available   *bool             `json:"available,omitempty"`
	GuildID     *string           `json:"guild_id,omitempty"`
	User        *User             `json:"user,omitempty"`
	SortValue   *int              `json:"sort_value,omitempty"`
}

type StickerType int

const (
	// StickerTypeStandard represents an official sticker in a pack.
	StickerTypeStandard StickerType = 1

	// StickerTypeGuild represents a sticker uploaded to a guild for the guild's members.
	StickerTypeGuild StickerType = 2
)

type StickerFormatType int

const (
	// StickerFormatTypePNG represents the PNG format for a sticker.
	StickerFormatTypePNG StickerFormatType = 1

	// StickerFormatTypeAPNG represents the APNG format for a sticker.
	StickerFormatTypeAPNG StickerFormatType = 2

	// StickerFormatTypeLOTTIE represents the LOTTIE format for a sticker.
	StickerFormatTypeLOTTIE StickerFormatType = 3

	// StickerFormatTypeGIF represents the GIF format for a sticker.
	StickerFormatTypeGIF StickerFormatType = 4
)

type Properties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type Presence struct {
	Since      *int64      `json:"since"`
	Activities []*Activity `json:"activities,omitempty"`
	Status     string      `json:"status"`
	AFK        bool        `json:"afk"`
}

type ApplicationCommandPermissions struct {
	// ID is id of the role, user, or channel.
	// It can also be a permission constant (@everyone: guild_id, all_channels: guild_id - 1).
	ID         string                           `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

type ApplicationCommandPermissionType int

const (
	PermissionTypeRole    ApplicationCommandPermissionType = 1
	PermissionTypeUser    ApplicationCommandPermissionType = 2
	PermissionTypeChannel ApplicationCommandPermissionType = 3
)

type AutoModerationRule struct {
	ID              string                            `json:"id"`
	GuildID         string                            `json:"guild_id"`
	Name            string                            `json:"name"`
	CreatorID       string                            `json:"creator_id"`
	EventType       AutoModerationRuleEventType       `json:"event_type"`
	TriggerType     AutoModerationRuleTriggerType     `json:"trigger_type"`
	TriggerMetadata AutoModerationRuleTriggerMetadata `json:"trigger_metadata"`
	Actions         []AutoModerationAction            `json:"actions"`
	Enabled         bool                              `json:"enabled"`
	ExemptRoles     []string                          `json:"exempt_roles"`
	ExemptChannels  []string                          `json:"exempt_channels"`
}

type AutoModerationRuleEventType int

const (
	AutoModerationRuleEventTypeMessageSend AutoModerationRuleEventType = 1
)

type AutoModerationRuleTriggerType int

const (
	AutoModerationRuleTriggerTypeKeyword       AutoModerationRuleTriggerType = 1
	AutoModerationRuleTriggerTypeSpam          AutoModerationRuleTriggerType = 3
	AutoModerationRuleTriggerTypeKeywordPreset AutoModerationRuleTriggerType = 4
	AutoModerationRuleTriggerTypeMentionSpam   AutoModerationRuleTriggerType = 5
)

// AutoModerationRuleTriggerMetadata contains additional data used to determine whether a rule should be triggered.
type AutoModerationRuleTriggerMetadata struct {
	// KeywordFilter is associated with KEYWORD trigger type. Substrings to search for in content.
	// A keyword can be a phrase which contains multiple words. Wildcard symbols can be used to customize how each keyword will be matched. Each keyword must be 60 characters or less.
	KeywordFilter []string `json:"keyword_filter,omitempty"`
	// RegexPatterns is associated with KEYWORD trigger type. Regular expression patterns to match against content.
	// Only Rust flavored regex is currently supported, which can be tested in online editors such as Rustexp. Each regex pattern must be 260 characters or less.
	RegexPatterns []string `json:"regex_patterns,omitempty"`
	// Presets is associated with KEYWORD_PRESET trigger type. Internally pre-defined wordsets to search for in content.
	Presets []AutoModerationRuleKeywordPresetType `json:"presets,omitempty"`
	// AllowList is associated with KEYWORD and KEYWORD_PRESET trigger types. Substrings that should not trigger the rule.
	// Each allow_list keyword can be a phrase which contains multiple words.
	// Wildcard symbols can be used to customize how each keyword will be matched.
	// Rules with KEYWORD trigger_type accept a maximum of 100 keywords.
	// Rules with KEYWORD_PRESET trigger_type accept a maximum of 1000 keywords.
	AllowList []string `json:"allow_list,omitempty"`
	// MentionTotalLimit is associated with MENTION_SPAM trigger type. Total number of unique role and user mentions allowed per message.
	MentionTotalLimit int `json:"mention_total_limit"`
	// MentionRaidProtection is associated with MENTION_SPAM trigger type. Whether to automatically detect mention raids.
	MentionRaidProtection bool `json:"mention_raid_protection_enabled"`
}

type AutoModerationRuleKeywordPresetType int

const (
	// AutoModerationRuleKeywordPresetProfanity represents words that may be considered forms of swearing or cursing.
	AutoModerationRuleKeywordPresetProfanity AutoModerationRuleKeywordPresetType = 1

	// AutoModerationRuleKeywordPresetSexualContent represents words that refer to sexually explicit behavior or activity.
	AutoModerationRuleKeywordPresetSexualContent AutoModerationRuleKeywordPresetType = 2

	// AutoModerationRuleKeywordPresetSlurs represents personal insults or words that may be considered hate speech.
	AutoModerationRuleKeywordPresetSlurs AutoModerationRuleKeywordPresetType = 3
)

// AutoModerationAction represents an action to be taken in auto moderation rules.
type AutoModerationAction struct {
	Type     AutoModerationRuleActionType `json:"type"`               // The type of action.
	Metadata any                          `json:"metadata,omitempty"` // Additional metadata needed during execution for this specific action type.
}

// AutoModerationRuleActionType represents the action types for auto moderation rules.
type AutoModerationRuleActionType int

const (
	// AutoModerationRuleActionBlockMessage blocks a member's message and prevents it from being posted.
	// A custom explanation can be specified and shown to members whenever their message is blocked.
	AutoModerationRuleActionBlockMessage AutoModerationRuleActionType = 1

	// AutoModerationRuleActionSendAlertMessage logs user content to a specified channel.
	AutoModerationRuleActionSendAlertMessage AutoModerationRuleActionType = 2

	// AutoModerationRuleActionTimeout times out a user for a specified duration.
	AutoModerationRuleActionTimeout AutoModerationRuleActionType = 3
)

// AutoModerationActionMetadata represents additional data used when an action is executed.
type AutoModerationActionMetadata struct {
	ChannelID       string  `json:"channel_id"`               // Channel to which user content should be logged for SEND_ALERT_MESSAGE action type.
	DurationSeconds int     `json:"duration_seconds"`         // Timeout duration in seconds for TIMEOUT action type.
	CustomMessage   *string `json:"custom_message,omitempty"` // Additional explanation that will be shown to members whenever their message is blocked for BLOCK_MESSAGE action type.
}

// AutoModerationActionExecutionEvent is that occurs when an action is executed in auto moderation.
type AutoModerationActionExecutionEvent struct {
	GuildID              string                        `json:"guild_id"`                          // ID of the guild in which action was executed.
	Action               AutoModerationAction          `json:"action"`                            // Action which was executed.
	RuleID               string                        `json:"rule_id"`                           // ID of the rule which action belongs to.
	RuleTriggerType      AutoModerationRuleTriggerType `json:"rule_trigger_type"`                 // Trigger type of rule which was triggered.
	UserID               string                        `json:"user_id"`                           // ID of the user which generated the content which triggered the rule.
	ChannelID            *string                       `json:"channel_id,omitempty"`              // ID of the channel in which user content was posted.
	MessageID            *string                       `json:"message_id,omitempty"`              // ID of any user message which content belongs to.
	AlertSystemMessageID *string                       `json:"alert_system_message_id,omitempty"` // ID of any system auto moderation messages posted as a result of this action.
	Content              string                        `json:"content"`                           // User-generated text content.
	MatchedKeyword       *string                       `json:"matched_keyword,omitempty"`         // Word or phrase configured in the rule that triggered the rule.
	MatchedContent       *string                       `json:"matched_content,omitempty"`         // Substring in content that triggered the rule.
}

// Entitlement represents an entitlement object.
type Entitlement struct {
	ID            string          `json:"id"`                  // ID of the entitlement.
	SKU_ID        string          `json:"sku_id"`              // ID of the SKU.
	ApplicationID string          `json:"application_id"`      // ID of the parent application.
	UserID        *string         `json:"user_id,omitempty"`   // ID of the user that is granted access to the entitlement's sku.
	Type          EntitlementType `json:"type"`                // Type of entitlement.
	Deleted       bool            `json:"deleted"`             // Entitlement was deleted.
	StartsAt      *time.Time      `json:"starts_at,omitempty"` // Start date at which the entitlement is valid. Not present when using test entitlements.
	EndsAt        *time.Time      `json:"ends_at,omitempty"`   // Date at which the entitlement is no longer valid. Not present when using test entitlements.
	GuildID       *string         `json:"guild_id,omitempty"`  // ID of the guild that is granted access to the entitlement's sku.
}

type EntitlementType int

const (
	EntitlementTypeApplicationSubscription EntitlementType = 8 // EntitlementType was purchased as an app subscription.
)

type AuditLogEntry struct {
	TargetID   *string          `json:"target_id,omitempty"` // ID of the affected entity (webhook, user, role, etc.).
	Changes    []AuditLogChange `json:"changes,omitempty"`   // Changes made to the target_id.
	UserID     *string          `json:"user_id,omitempty"`   // User or app that made the changes.
	ID         string           `json:"id"`                  // ID of the entry.
	ActionType AuditLogEvent    `json:"action_type"`         // Type of action that occurred.
	Options    *AuditEntryInfo  `json:"options,omitempty"`   // Additional info for certain event types.
	Reason     *string          `json:"reason,omitempty"`    // Reason for the change (1-512 characters).
}

type AuditLogChange struct {
	NewValue any    `json:"new_value,omitempty"` // New value of the key. TODO: avoid any if possible.
	OldValue any    `json:"old_value,omitempty"` // Old value of the key. TODO: avoid any if possible.
	Key      string `json:"key"`                 // Key representing the changed property.
}

type AuditLogEvent int

const (
	AuditLogEventGuildUpdate                             AuditLogEvent = 1
	AuditLogEventChannelCreate                           AuditLogEvent = 10
	AuditLogEventChannelUpdate                           AuditLogEvent = 11
	AuditLogEventChannelDelete                           AuditLogEvent = 12
	AuditLogEventChannelOverwriteCreate                  AuditLogEvent = 13
	AuditLogEventChannelOverwriteUpdate                  AuditLogEvent = 14
	AuditLogEventChannelOverwriteDelete                  AuditLogEvent = 15
	AuditLogEventMemberKick                              AuditLogEvent = 20
	AuditLogEventMemberPrune                             AuditLogEvent = 21
	AuditLogEventMemberBanAdd                            AuditLogEvent = 22
	AuditLogEventMemberBanRemove                         AuditLogEvent = 23
	AuditLogEventMemberUpdate                            AuditLogEvent = 24
	AuditLogEventMemberRoleUpdate                        AuditLogEvent = 25
	AuditLogEventMemberMove                              AuditLogEvent = 26
	AuditLogEventMemberDisconnect                        AuditLogEvent = 27
	AuditLogEventBotAdd                                  AuditLogEvent = 28
	AuditLogEventRoleCreate                              AuditLogEvent = 30
	AuditLogEventRoleUpdate                              AuditLogEvent = 31
	AuditLogEventRoleDelete                              AuditLogEvent = 32
	AuditLogEventInviteCreate                            AuditLogEvent = 40
	AuditLogEventInviteUpdate                            AuditLogEvent = 41
	AuditLogEventInviteDelete                            AuditLogEvent = 42
	AuditLogEventWebhookCreate                           AuditLogEvent = 50
	AuditLogEventWebhookUpdate                           AuditLogEvent = 51
	AuditLogEventWebhookDelete                           AuditLogEvent = 52
	AuditLogEventEmojiCreate                             AuditLogEvent = 60
	AuditLogEventEmojiUpdate                             AuditLogEvent = 61
	AuditLogEventEmojiDelete                             AuditLogEvent = 62
	AuditLogEventMessageDelete                           AuditLogEvent = 72
	AuditLogEventMessageBulkDelete                       AuditLogEvent = 73
	AuditLogEventMessagePin                              AuditLogEvent = 74
	AuditLogEventMessageUnpin                            AuditLogEvent = 75
	AuditLogEventIntegrationCreate                       AuditLogEvent = 80
	AuditLogEventIntegrationUpdate                       AuditLogEvent = 81
	AuditLogEventIntegrationDelete                       AuditLogEvent = 82
	AuditLogEventStageInstanceCreate                     AuditLogEvent = 83
	AuditLogEventStageInstanceUpdate                     AuditLogEvent = 84
	AuditLogEventStageInstanceDelete                     AuditLogEvent = 85
	AuditLogEventStickerCreate                           AuditLogEvent = 90
	AuditLogEventStickerUpdate                           AuditLogEvent = 91
	AuditLogEventStickerDelete                           AuditLogEvent = 92
	AuditLogEventGuildScheduledEventCreate               AuditLogEvent = 100
	AuditLogEventGuildScheduledEventUpdate               AuditLogEvent = 101
	AuditLogEventGuildScheduledEventDelete               AuditLogEvent = 102
	AuditLogEventThreadCreate                            AuditLogEvent = 110
	AuditLogEventThreadUpdate                            AuditLogEvent = 111
	AuditLogEventThreadDelete                            AuditLogEvent = 112
	AuditLogEventApplicationCommandPermissionUpdate      AuditLogEvent = 121
	AuditLogEventAutoModerationRuleCreate                AuditLogEvent = 140
	AuditLogEventAutoModerationRuleUpdate                AuditLogEvent = 141
	AuditLogEventAutoModerationRuleDelete                AuditLogEvent = 142
	AuditLogEventAutoModerationBlockMessage              AuditLogEvent = 143
	AuditLogEventAutoModerationFlagToChannel             AuditLogEvent = 144
	AuditLogEventAutoModerationUserCommunicationDisabled AuditLogEvent = 145
	AuditLogEventCreatorMonetizationRequestCreated       AuditLogEvent = 150
	AuditLogEventCreatorMonetizationTermsAccepted        AuditLogEvent = 151
)

type AuditEntryInfo struct {
	ApplicationID                 string `json:"application_id"`                    // ID of the app whose permissions were targeted (APPLICATION_COMMAND_PERMISSION_UPDATE).
	AutoModerationRuleName        string `json:"auto_moderation_rule_name"`         // Name of the Auto Moderation rule that was triggered (AUTO_MODERATION_BLOCK_MESSAGE, AUTO_MODERATION_FLAG_TO_CHANNEL, AUTO_MODERATION_USER_COMMUNICATION_DISABLED).
	AutoModerationRuleTriggerType string `json:"auto_moderation_rule_trigger_type"` // Trigger type of the Auto Moderation rule that was triggered (AUTO_MODERATION_BLOCK_MESSAGE, AUTO_MODERATION_FLAG_TO_CHANNEL, AUTO_MODERATION_USER_COMMUNICATION_DISABLED).
	ChannelID                     string `json:"channel_id"`                        // Channel in which the entities were targeted (MEMBER_MOVE, MESSAGE_PIN, MESSAGE_UNPIN, MESSAGE_DELETE, STAGE_INSTANCE_CREATE, STAGE_INSTANCE_UPDATE, STAGE_INSTANCE_DELETE, AUTO_MODERATION_BLOCK_MESSAGE, AUTO_MODERATION_FLAG_TO_CHANNEL, AUTO_MODERATION_USER_COMMUNICATION_DISABLED).
	Count                         string `json:"count"`                             // Number of entities that were targeted (MESSAGE_DELETE, MESSAGE_BULK_DELETE, MEMBER_DISCONNECT, MEMBER_MOVE).
	DeleteMemberDays              string `json:"delete_member_days"`                // Number of days after which inactive members were kicked (MEMBER_PRUNE).
	ID                            string `json:"id"`                                // ID of the overwritten entity (CHANNEL_OVERWRITE_CREATE, CHANNEL_OVERWRITE_UPDATE, CHANNEL_OVERWRITE_DELETE).
	MembersRemoved                string `json:"members_removed"`                   // Number of members removed by the prune (MEMBER_PRUNE).
	MessageID                     string `json:"message_id"`                        // ID of the message that was targeted (MESSAGE_PIN, MESSAGE_UNPIN).
	RoleName                      string `json:"role_name"`                         // Name of the role if type is "0" (not present if type is "1") (CHANNEL_OVERWRITE_CREATE, CHANNEL_OVERWRITE_UPDATE, CHANNEL_OVERWRITE_DELETE).
	Type                          string `json:"type"`                              // Type of overwritten entity - role ("0") or member ("1") (CHANNEL_OVERWRITE_CREATE, CHANNEL_OVERWRITE_UPDATE, CHANNEL_OVERWRITE_DELETE).
	IntegrationType               string `json:"integration_type"`                  // The type of integration which performed the action (MEMBER_KICK, MEMBER_ROLE_UPDATE).
}

type GuildScheduledEvent struct {
	ID                 string                             `json:"id"`                           // The id of the scheduled event.
	GuildID            string                             `json:"guild_id"`                     // The guild id which the scheduled event belongs to.
	ChannelID          *string                            `json:"channel_id,omitempty"`         // The channel id in which the scheduled event will be hosted, or null if scheduled entity type is EXTERNAL.
	CreatorID          *string                            `json:"creator_id,omitempty"`         // The id of the user that created the scheduled event.
	Name               string                             `json:"name"`                         // The name of the scheduled event (1-100 characters).
	Description        *string                            `json:"description,omitempty"`        // The description of the scheduled event (1-1000 characters).
	ScheduledStartTime string                             `json:"scheduled_start_time"`         // The time the scheduled event will start.
	ScheduledEndTime   *string                            `json:"scheduled_end_time,omitempty"` // The time the scheduled event will end, required if entity_type is EXTERNAL.
	PrivacyLevel       GuildScheduledEventPrivacyLevel    `json:"privacy_level"`                // The privacy level of the scheduled event.
	Status             GuildScheduledEventStatus          `json:"status"`                       // The status of the scheduled event.
	EntityType         GuildScheduledEventEntityType      `json:"entity_type"`                  // The type of the scheduled event.
	EntityID           *string                            `json:"entity_id,omitempty"`          // The id of an entity associated with a guild scheduled event.
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`    // Additional metadata for the guild scheduled event.
	Creator            *User                              `json:"creator,omitempty"`            // The user that created the scheduled event.
	UserCount          *int                               `json:"user_count,omitempty"`         // The number of users subscribed to the scheduled event.
	Image              *string                            `json:"image,omitempty"`              // The cover image hash of the scheduled event.
}

// GuildScheduledEventPrivacyLevel represents the privacy level of a guild scheduled event.
type GuildScheduledEventPrivacyLevel int

const (
	GuildScheduledEVentPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = 2 // The scheduled event is only accessible to guild members.
)

// GuildScheduledEventEntityType represents the type of a guild scheduled event.
type GuildScheduledEventEntityType int

const (
	GuildScheduledEventEntityTypeStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledEventEntityTypeVoice         GuildScheduledEventEntityType = 2
	GuildScheduledEventEntityTypeExternal      GuildScheduledEventEntityType = 3
)

// GuildScheduledEventStatus represents the status of a guild scheduled event.
type GuildScheduledEventStatus int

const (
	ScheduledStatus GuildScheduledEventStatus = 1
	ActiveStatus    GuildScheduledEventStatus = 2
	CompletedStatus GuildScheduledEventStatus = 3
	CanceledStatus  GuildScheduledEventStatus = 4
)

// GuildScheduledEventEntityMetadata represents additional metadata for a guild scheduled event.
type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location"`
}

// Integration represents an integration within a guild.
type Integration struct {
	ID                string                     `json:"id"`
	Name              string                     `json:"name"`
	Type              string                     `json:"type"`
	Enabled           bool                       `json:"enabled"`
	Syncing           *bool                      `json:"syncing,omitempty"`
	RoleID            *string                    `json:"role_id,omitempty"`
	EnableEmoticons   *bool                      `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *IntegrationExpireBehavior `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int                       `json:"expire_grace_period,omitempty"`
	User              *User                      `json:"user,omitempty"`
	Account           *IntegrationAccount        `json:"account"`
	SyncedAt          *time.Time                 `json:"synced_at,omitempty"`
	SubscriberCount   *int                       `json:"subscriber_count,omitempty"`
	Revoked           *bool                      `json:"revoked,omitempty"`
	Application       *IntegrationApplication    `json:"application,omitempty"`
	Scopes            []string                   `json:"scopes,omitempty"`
}

// IntegrationExpireBehavior represents the behavior of expiring subscribers for an integration.
type IntegrationExpireBehavior int

const (
	// IntegrationExpireRemoveRole removes the role of expiring subscribers.
	IntegrationExpireRemoveRole IntegrationExpireBehavior = iota

	// IntegrationExpireKick kicks expiring subscribers.
	IntegrationExpireKick
)

// IntegrationAccount represents the account information for an integration.
type IntegrationAccount struct {
	ID   string `json:"id"`   // ID of the account.
	Name string `json:"name"` // Name of the account.
}

// IntegrationApplication represents the application associated with an integration.
type IntegrationApplication struct {
	ID          string  `json:"id"`             // ID of the app.
	Name        string  `json:"name"`           // Name of the app.
	Icon        *string `json:"icon,omitempty"` // Icon hash of the app.
	Description string  `json:"description"`    // Description of the app.
	Bot         *User   `json:"bot,omitempty"`  // Bot associated with this application.
}

// InviteTargetType represents the type of target for an invite.
type InviteTargetType int

const (
	// InviteTargetStream represents a target type of a voice channel stream.
	InviteTargetTypeStream InviteTargetType = iota + 1

	// InviteTargetEmbeddedApplication represents a target type of an embedded application.
	InviteTargetTypeEmbeddedApplication
)

// Message represents a message sent in a channel within Discord.
type Message struct {
	ID                   string                       `json:"id"`                               // The ID of the message.
	ChannelID            string                       `json:"channel_id"`                       // The ID of the channel the message was sent in.
	GuildID              *string                      `json:"guild_id,omitempty"`               // ID of the guild the message was sent in - unless it is an ephemeral message.
	Author               User                         `json:"author"`                           // The author of this message. If the message was created by a webhook, the user represents the webhook's id, username and avatar.
	Member               *GuildMember                 `json:"member,omitempty"`                 // Member properties for this message's author. Missing for ephemeral messages and messages from webhooks.
	Content              string                       `json:"content"`                          // Contents of the message.
	Timestamp            time.Time                    `json:"timestamp"`                        // When this message was sent.
	EditedTimestamp      *time.Time                   `json:"edited_timestamp,omitempty"`       // When this message was edited (or null if never).
	TTS                  bool                         `json:"tts"`                              // Whether this was a TTS message.
	MentionEveryone      bool                         `json:"mention_everyone"`                 // Whether this message mentions everyone.
	Mentions             []User                       `json:"mentions,omitempty"`               // Users specifically mentioned in the message.
	MentionRoles         []string                     `json:"mention_roles,omitempty"`          // Roles specifically mentioned in this message.
	MentionChannels      []ChannelMention             `json:"mention_channels,omitempty"`       // Channels specifically mentioned in this message.
	Attachments          []MessageAttachment          `json:"attachments,omitempty"`            // Any attached files.
	Embeds               []Embed                      `json:"embeds,omitempty"`                 // Any embedded content.
	Reactions            []Reaction                   `json:"reactions,omitempty"`              // Reactions to the message.
	Nonce                *string                      `json:"nonce,omitempty"`                  // Used for validating a message was sent.
	Pinned               bool                         `json:"pinned"`                           // Whether this message is pinned.
	WebhookID            *string                      `json:"webhook_id,omitempty"`             // If the message is generated by a webhook, this is the webhook's id.
	Type                 MessageType                  `json:"type"`                             // Type of message.
	Activity             *MessageActivity             `json:"activity,omitempty"`               // Sent with Rich Presence-related chat embeds.
	Application          *Application                 `json:"application,omitempty"`            // Sent with Rich Presence-related chat embeds.
	ApplicationID        *string                      `json:"application_id,omitempty"`         // If the message is an Interaction or application-owned webhook, this is the id of the application.
	MessageReference     *MessageReference            `json:"message_reference,omitempty"`      // Data showing the source of a crosspost, channel follow add, pin, or reply message.
	Flags                *int                         `json:"flags,omitempty"`                  // Message flags combined as a bitfield.
	ReferencedMessage    *Message                     `json:"referenced_message,omitempty"`     // The message associated with the message_reference.
	Interaction          *MessageInteraction          `json:"interaction,omitempty"`            // Sent if the message is a response to an Interaction.
	Thread               *Channel                     `json:"thread,omitempty"`                 // The thread that was started from this message, includes thread member object.
	ActionTypes          []MessageActionType          `json:"components,omitempty"`             // Sent if the message contains components like buttons, action rows, or other interactive components. To get these fully marshaled, use the GetComponents() method on the ActionType.
	StickerItems         []MessageStickerItem         `json:"sticker_items,omitempty"`          // Sent if the message contains stickers.
	Position             *int                         `json:"position,omitempty"`               // The approximate position of the message in a thread.
	RoleSubscriptionData *MessageRoleSubscriptionData `json:"role_subscription_data,omitempty"` // Data of the role subscription purchase or renewal that prompted this ROLE_SUBSCRIPTION_PURCHASE message.
	Resolved             *MessageResolvedData         `json:"resolved,omitempty"`               // Data for users, members, channels, and roles in the message's auto-populated select menus.
}

// ChannelMention represents a mentioned channel in a message.
type ChannelMention struct {
	ID      string      `json:"id"`       // The ID of the channel.
	GuildID string      `json:"guild_id"` // The ID of the guild containing the channel.
	Type    ChannelType `json:"type"`     // The type of channel.
	Name    string      `json:"name"`     // The name of the channel.
}

// MessageAttachment represents an attachment in a message.
type MessageAttachment struct {
	ID           string   `json:"id"`                      // The ID of the attachment.
	Filename     string   `json:"filename"`                // The name of the file attached.
	Description  *string  `json:"description,omitempty"`   // Description for the file (max 1024 characters).
	ContentType  *string  `json:"content_type,omitempty"`  // The attachment's media type.
	Size         int      `json:"size"`                    // Size of the file in bytes.
	URL          string   `json:"url"`                     // Source URL of the file.
	ProxyURL     string   `json:"proxy_url"`               // A proxied URL of the file.
	Height       *int     `json:"height,omitempty"`        // Height of the file (if image).
	Width        *int     `json:"width,omitempty"`         // Width of the file (if image).
	Ephemeral    *bool    `json:"ephemeral,omitempty"`     // Whether this attachment is ephemeral.
	DurationSecs *float64 `json:"duration_secs,omitempty"` // The duration of the audio file (currently for voice messages).
	Waveform     *string  `json:"waveform,omitempty"`      // Base64 encoded bytearray representing a sampled waveform (currently for voice messages).
	Flags        *int     `json:"flags,omitempty"`         // Attachment flags combined as a bitfield.
}

// Embed represents an embedded content in a message.
type Embed struct {
	Title       *string         `json:"title,omitempty"`       // Title of the embed.
	Type        *string         `json:"type,omitempty"`        // Type of the embed (always "rich" for webhook embeds).
	Description *string         `json:"description,omitempty"` // Description of the embed.
	URL         *string         `json:"url,omitempty"`         // URL of the embed.
	Timestamp   *string         `json:"timestamp,omitempty"`   // Timestamp of embed content in ISO8601 format.
	Color       *int            `json:"color,omitempty"`       // Color code of the embed.
	Footer      *EmbedFooter    `json:"footer,omitempty"`      // Footer information.
	Image       *EmbedImage     `json:"image,omitempty"`       // Image information.
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`   // Thumbnail information.
	Video       *EmbedVideo     `json:"video,omitempty"`       // Video information.
	Provider    *EmbedProvider  `json:"provider,omitempty"`    // Provider information.
	Author      *EmbedAuthor    `json:"author,omitempty"`      // Author information.
	Fields      []*EmbedField   `json:"fields,omitempty"`      // Fields information, max of 25.
}

// EmbedThumbnail represents the thumbnail of an embed.
type EmbedThumbnail struct {
	URL      string  `json:"url"`                 // Source URL of thumbnail (only supports http(s) and attachments).
	ProxyURL *string `json:"proxy_url,omitempty"` // A proxied URL of the thumbnail.
	Height   *int    `json:"height,omitempty"`    // Height of thumbnail.
	Width    *int    `json:"width,omitempty"`     // Width of thumbnail.
}

// EmbedVideo represents the video content of an embed.
type EmbedVideo struct {
	URL      *string `json:"url,omitempty"`       // Source URL of video.
	ProxyURL *string `json:"proxy_url,omitempty"` // A proxied URL of the video.
	Height   *int    `json:"height,omitempty"`    // Height of video.
	Width    *int    `json:"width,omitempty"`     // Width of video.
}

// EmbedImage represents the image content of an embed.
type EmbedImage struct {
	URL      string  `json:"url"`                 // Source URL of image (only supports http(s) and attachments).
	ProxyURL *string `json:"proxy_url,omitempty"` // A proxied URL of the image.
	Height   *int    `json:"height,omitempty"`    // Height of image.
	Width    *int    `json:"width,omitempty"`     // Width of image.
}

// EmbedProvider represents the provider information of an embed.
type EmbedProvider struct {
	Name *string `json:"name,omitempty"` // Name of provider.
	URL  *string `json:"url,omitempty"`  // URL of provider.
}

// EmbedAuthor represents the author information of an embed.
type EmbedAuthor struct {
	Name         string  `json:"name"`                     // Name of author.
	URL          *string `json:"url,omitempty"`            // URL of author (only supports http(s)).
	IconURL      *string `json:"icon_url,omitempty"`       // URL of author icon (only supports http(s) and attachments).
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"` // A proxied URL of author icon.
}

// EmbedFooter represents the footer information of an embed.
type EmbedFooter struct {
	Text         string  `json:"text"`                     // Footer text.
	IconURL      *string `json:"icon_url,omitempty"`       // URL of footer icon (only supports http(s) and attachments).
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"` // A proxied URL of footer icon.
}

// EmbedField represents a field in an embed.
type EmbedField struct {
	Name   string `json:"name"`             // Name of the field.
	Value  string `json:"value"`            // Value of the field.
	Inline *bool  `json:"inline,omitempty"` // Whether or not this field should display inline.
}

// Reaction represents a reaction to a message.
type Reaction struct {
	Count        int                   `json:"count"`                   // Total number of times this emoji has been used to react (including super reacts).
	CountDetails *ReactionCountDetails `json:"count_details,omitempty"` // Reaction count details object.
	Me           bool                  `json:"me"`                      // Whether the current user reacted using this emoji.
	MeBurst      bool                  `json:"me_burst"`                // Whether the current user super-reacted using this emoji.
	Emoji        *Emoji                `json:"emoji"`                   // Emoji information.
	BurstColors  []string              `json:"burst_colors,omitempty"`  // HEX colors used for super reaction.
}

// ReactionCountDetails represents the breakdown of normal and super reaction counts for an emoji.
type ReactionCountDetails struct {
	Burst  int `json:"burst"`  // Count of super reactions.
	Normal int `json:"normal"` // Count of normal reactions.
}

// MessageMentionedUser is a user, with a partial guild member in it.
// I have no idea why discord decides to suddenly nestle guild member <-> user the OTHER WAY, but I'm sure there's a logical explanation for that.
type MessageMentionedUser struct {
	User

	Member GuildMember `json:"member"`
}

type MessageType int

const (
	DefaultMessageType                                 MessageType = iota // 0 (Deletable: true)
	RecipientAddMessageType                                               // 1 (Deletable: false)
	RecipientRemoveMessageType                                            // 2 (Deletable: false)
	CallMessageType                                                       // 3 (Deletable: false)
	ChannelNameChangeMessageType                                          // 4 (Deletable: false)
	ChannelIconChangeMessageType                                          // 5 (Deletable: false)
	ChannelPinnedMessageMessageType                                       // 6 (Deletable: true)
	UserJoinMessageType                                                   // 7 (Deletable: true)
	GuildBoostMessageType                                                 // 8 (Deletable: true)
	GuildBoostTier1MessageType                                            // 9 (Deletable: true)
	GuildBoostTier2MessageType                                            // 10 (Deletable: true)
	GuildBoostTier3MessageType                                            // 11 (Deletable: true)
	ChannelFollowAddMessageType                                           // 12 (Deletable: true)
	GuildDiscoveryDisqualifiedMessageType                                 // 14 (Deletable: false)
	GuildDiscoveryRequalifiedMessageType                                  // 15 (Deletable: false)
	GuildDiscoveryGracePeriodInitialWarningMessageType                    // 16 (Deletable: false)
	GuildDiscoveryGracePeriodFinalWarningMessageType                      // 17 (Deletable: false)
	ThreadCreatedMessageType                                              // 18 (Deletable: true)
	ReplyMessageType                                                      // 19 (Deletable: true)
	ChatInputCommandMessageType                                           // 20 (Deletable: true)
	ThreadStarterMessageMessageType                                       // 21 (Deletable: false)
	GuildInviteReminderMessageType                                        // 22 (Deletable: true)
	ContextMenuCommandMessageType                                         // 23 (Deletable: true)
	AutoModerationActionMessageType                                       // 24 (Deletable: true)
	RoleSubscriptionPurchaseMessageType                                   // 25 (Deletable: true)
	InteractionPremiumUpsellMessageType                                   // 26 (Deletable: true)
	StageStartMessageType                                                 // 27 (Deletable: true)
	StageEndMessageType                                                   // 28 (Deletable: true)
	StageSpeakerMessageType                                               // 29 (Deletable: true)
	StageTopicMessageType                                                 // 31 (Deletable: true)
	GuildApplicationPremiumSubscriptionMessageType                        // 32 (Deletable: false)
)

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`     // type of message activity
	PartyID *string             `json:"party_id"` // party_id from a Rich Presence event.
}

// MessageReference represents a reference to a message.
type MessageReference struct {
	MessageID       *string `json:"message_id"`         // ID of the originating message.
	ChannelID       *string `json:"channel_id"`         // ID of the originating message's channel. Channel_id is optional when creating a reply, but will always be present when receiving an event/response that includes this data model.
	GuildID         *string `json:"guild_id"`           // ID of the originating message's guild.
	FailIfNotExists *bool   `json:"fail_if_not_exists"` // When sending, whether to error if the referenced message doesn't exist instead of sending as a normal (non-reply) message, default true.
}

type MessageActivityType int

const (
	MessageActivityTypeJoin        MessageActivityType = 1
	MessageActivityTypeSpectate    MessageActivityType = 2
	MessageActivityTypeListen      MessageActivityType = 3
	MessageActivityTypeJoinRequest MessageActivityType = 5
)

// MessageInteraction represents an interaction that occurs within a message.
type MessageInteraction struct {
	ID     string                 `json:"id"`               // ID of the interaction.
	Type   MessageInteractionType `json:"type"`             // Type of interaction.
	Name   string                 `json:"name"`             // Name of the application command, including subcommands and subcommand groups.
	User   User                   `json:"user"`             // User who invoked the interaction.
	Member *GuildMember           `json:"member,omitempty"` // Member who invoked the interaction in the guild.
}

// MessageInteractionType represents the type of message interaction.
type MessageInteractionType int

// Constants representing different message interaction types.
const (
	MessageInteractionPing                    MessageInteractionType = 1
	MessageInteractionApplicationCommand      MessageInteractionType = 2
	MessageInteractionMessageComponent        MessageInteractionType = 3
	MessageInteractionApplicationAutocomplete MessageInteractionType = 4
	MessageInteractionModalSubmit             MessageInteractionType = 5
)

type baseComponent struct {
	Type ComponentType `json:"type"`
}

// MessageActionType is what discord calls components.
// The reason why it's called action type here is because the object is nestled in itself, calling itself an action type in that context.
type MessageActionType struct {
	Type ComponentType
	// Components is a list of Button/SelectMenu/TextInput objects.
	Components []json.RawMessage `json:"components"`
}

// GetComponents unmarshals the components based on the type into a array.
// To access these, use a type switch.
func (mat *MessageActionType) GetComponents() ([]any, error) {
	var ret []any
	for _, component := range mat.Components {
		base, err := UnmarshalJSON[baseComponent](component)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal component type: %w", err)
		}

		switch base.Type {
		case ComponentTypeButton:
			button, err := UnmarshalJSON[Button](component)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal component type: %w", err)
			}

			ret = append(ret, button)
		}
	}

	return ret, nil
}

// ComponentType represents the type of a UI component.
type ComponentType int

// Component types.
const (
	ComponentTypeActionRow ComponentType = iota + 1
	ComponentTypeButton
	ComponentTypeStringSelect
	ComponentTypeTextInput
	ComponentTypeUserSelect
	ComponentTypeRoleSelect
	ComponentTypeMentionableSelect
	ComponentTypeChannelSelect
)

// Button represents a button component within a message.
type Button struct {
	Type     ComponentType `json:"type"`                // Type of component (always 2 for a button).
	Style    ButtonStyle   `json:"style"`               // Style of the button.
	Label    *string       `json:"label,omitempty"`     // Text that appears on the button; max 80 characters.
	Emoji    *Emoji        `json:"emoji,omitempty"`     // Emoji for the button.
	CustomID *string       `json:"custom_id,omitempty"` // Developer-defined identifier for the button; max 100 characters.
	URL      *string       `json:"url,omitempty"`       // URL for link-style buttons.
	Disabled *bool         `json:"disabled,omitempty"`  // Whether the button is disabled (defaults to false).
}

// ButtonStyle represents the style of a button component.
type ButtonStyle int

const (
	// ButtonStylePrimary represents a primary button (color: blurple).
	ButtonStylePrimary ButtonStyle = 1

	// ButtonStyleSecondary represents a secondary button (color: grey).
	ButtonStyleSecondary ButtonStyle = 2

	// ButtonStyleSuccess represents a success button (color: green).
	ButtonStyleSuccess ButtonStyle = 3

	// ButtonStyleDanger represents a danger button (color: red).
	ButtonStyleDanger ButtonStyle = 4

	// ButtonStyleLink represents a link button (color: grey, navigates to a URL).
	ButtonStyleLink ButtonStyle = 5
)

// SelectMenu represents a select menu component.
type SelectMenu struct {
	Type          ComponentType         `json:"type"`                     // Type of select menu component (text: 3, user: 5, role: 6, mentionable: 7, channels: 8).
	CustomID      string                `json:"custom_id"`                // ID for the select menu. Max 100 characters.
	Options       []*SelectOption       `json:"options,omitempty"`        // Specified choices in a select menu. Only required and available for string selects (type 3). Max 25.
	ChannelTypes  []int                 `json:"channel_types,omitempty"`  // List of channel types to include in the channel select component (type 8).
	Placeholder   *string               `json:"placeholder,omitempty"`    // Placeholder text if nothing is selected. Max 150 characters.
	DefaultValues []*SelectDefaultValue `json:"default_values,omitempty"` // List of default values for auto-populated select menu components. The number of default values must be in the range defined by MinValues and MaxValues.
	MinValues     *int                  `json:"min_values,omitempty"`     // Minimum number of items that must be chosen. Defaults to 1. Min 0, Max 25.
	MaxValues     *int                  `json:"max_values,omitempty"`     // Maximum number of items that can be chosen. Defaults to 1. Max 25.
	Disabled      *bool                 `json:"disabled,omitempty"`       // Whether the select menu is disabled.
}

// SelectOption represents an option in a select menu.
type SelectOption struct {
	Label       string  `json:"label"`                 // User-facing name of the option. Max 100 characters.
	Value       string  `json:"value"`                 // Dev-defined value of the option. Max 100 characters.
	Description *string `json:"description,omitempty"` // Additional description of the option. Max 100 characters.
	Emoji       *Emoji  `json:"emoji,omitempty"`       // Emoji object containing ID, name, and animated status.
	Default     *bool   `json:"default,omitempty"`     // Indicates whether this option is selected by default.
}

// SelectDefaultValue represents a default value for auto-populated select menu components.
type SelectDefaultValue struct {
	ID   string `json:"id"`   // ID of a user, role, or channel.
	Type string `json:"type"` // Type of value that ID represents. Either "user", "role", or "channel".
}

// TextInput represents a text input component in an interactive message.
type TextInput struct {
	Type        ComponentType  `json:"type"`                  // Type of the text input component (4 for text input).
	CustomID    string         `json:"custom_id"`             // Developer-defined identifier for the input; max 100 characters.
	Style       TextInputStyle `json:"style"`                 // The Text Input Style.
	Label       string         `json:"label"`                 // Label for this component; max 45 characters.
	MinLength   *int           `json:"min_length,omitempty"`  // Minimum input length for a text input; min 0, max 4000.
	MaxLength   *int           `json:"max_length,omitempty"`  // Maximum input length for a text input; min 1, max 4000.
	Required    *bool          `json:"required,omitempty"`    // Whether this component is required to be filled (defaults to true).
	Value       *string        `json:"value,omitempty"`       // Pre-filled value for this component; max 4000 characters.
	Placeholder *string        `json:"placeholder,omitempty"` // Custom placeholder text if the input is empty; max 100 characters.
}

// TextInputStyle represents the style of a text input component.
type TextInputStyle int

// Constants representing different styles of text input components.
const (
	TextInputStyleShort     TextInputStyle = 1 // Single-line input.
	TextInputStyleParagraph TextInputStyle = 2 // Multi-line input.
)

// MessageStickerItem represents the smallest amount of data required to render a sticker.
type MessageStickerItem struct {
	ID         string `json:"id"`          // ID of the sticker.
	Name       string `json:"name"`        // Name of the sticker.
	FormatType int    `json:"format_type"` // Type of sticker format.
}

// MessageRoleSubscriptionData represents the data associated with a role subscription.
type MessageRoleSubscriptionData struct {
	RoleSubscriptionListingID string `json:"role_subscription_listing_id"` // ID of the SKU and listing that the user is subscribed to.
	TierName                  string `json:"tier_name"`                    // Name of the tier that the user is subscribed to.
	TotalMonthsSubscribed     int    `json:"total_months_subscribed"`      // Cumulative number of months that the user has been subscribed for.
	IsRenewal                 bool   `json:"is_renewal"`                   // Whether this notification is for a renewal rather than a new purchase.
}

// MessageResolvedData represents the resolved data for users, members, channels, roles, messages, and attachments.
type MessageResolvedData struct {
	Users       map[string]User              `json:"users"`       // The IDs and User objects.
	Members     map[string]GuildMember       `json:"members"`     // The IDs and partial Member objects.
	Roles       map[string]Role              `json:"roles"`       // The IDs and Role objects.
	Channels    map[string]Channel           `json:"channels"`    // The IDs and partial Channel objects.
	Messages    map[string]Message           `json:"messages"`    // The IDs and partial Message objects.
	Attachments map[string]MessageAttachment `json:"attachments"` // The IDs and attachment objects.
}

// AllowedMentions represents allowed mentions for the message.
type AllowedMentions struct {
	Parse       []string `json:"parse,omitempty"`        // An array of allowed mention types to parse from the content.
	Roles       []string `json:"roles,omitempty"`        // Array of role_ids to mention (Max size of 100).
	Users       []string `json:"users,omitempty"`        // Array of user_ids to mention (Max size of 100).
	RepliedUser bool     `json:"replied_user,omitempty"` // For replies, whether to mention the author of the message being replied to (default false).
}

// AllowedMentionType represents the type of allowed mentions.
type AllowedMentionType string

const (
	// AllowedMentionRole controls role mentions.
	AllowedMentionRole AllowedMentionType = "roles"
	// AllowedMentionUser controls user mentions.
	AllowedMentionUser AllowedMentionType = "users"
	// AllowedMentionEveryone controls @everyone and @here mentions.
	AllowedMentionEveryone AllowedMentionType = "everyone"
)
