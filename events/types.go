package events

import "time"

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
	PremiumProgressBarEnabled bool           `json:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID     *string        `json:"safety_alerts_channel_id,omitempty,omitempty"`
	Unavailable               bool           `json:"unavailable"`
}

type VerificationLevel uint32

const (
	// VerificationLevelNone is unrestricted
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
	BotID                 string `json:"bot_id,omitempty"`
	IntegrationID         string `json:"integration_id,omitempty"`
	PremiumSubscriber     bool   `json:"premium_subscriber,omitempty"`
	SubscriptionListingID string `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  bool   `json:"available_for_purchase,omitempty"`
	GuildConnections      bool   `json:"guild_connections,omitempty"`
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

type VoiceState struct {
	GuildID          *string      `json:"guild_id,omitempty"`
	ChannelID        *string      `json:"channel_id,omitempty"`
	UserID           string       `json:"user_id"`
	Member           *GuildMember `json:"member,omitempty"`
	SessionID        string       `json:"session_id"`
	Deaf             bool         `json:"deaf"`
	Mute             bool         `json:"mute"`
	SelfDeaf         bool         `json:"self_deaf"`
	SelfMute         bool         `json:"self_mute"`
	SelfStream       *bool        `json:"self_stream,omitempty"`
	SelfVideo        bool         `json:"self_video"`
	Suppress         bool         `json:"suppress"`
	RequestToSpeakTS *time.Time   `json:"request_to_speak_timestamp,omitempty"`
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
	Type                 int                   `json:"type"`
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
	RateLimitPerUser              *int             `json:"rate_limit_per_user,omitempty"`
	Recipients                    []User           `json:"recipients,omitempty"`
	Icon                          *string          `json:"icon,omitempty"`
	OwnerID                       *string          `json:"owner_id,omitempty"`
	ApplicationID                 *string          `json:"application_id,omitempty"`
	Managed                       bool             `json:"managed"`
	ParentID                      *string          `json:"parent_id,omitempty"`
	LastPinTimestamp              *time.Time       `json:"last_pin_timestamp,omitempty"`
	RTCRegion                     *string          `json:"rtc_region,omitempty"`
	VideoQualityMode              *int             `json:"video_quality_mode,omitempty"`
	MessageCount                  *int             `json:"message_count,omitempty"`
	MemberCount                   *int             `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata  `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember    `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int             `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *int             `json:"permissions,omitempty"`
	Flags                         *int             `json:"flags,omitempty"`
	TotalMessagesSent             *int             `json:"total_messages_sent,omitempty"`
	AvailableTags                 []ForumTag       `json:"available_tags,omitempty"`
	AppliedTags                   []string         `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReaction `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int             `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *int             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *int             `json:"default_forum_layout,omitempty"`
}

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
	// EmojiID is the id of a guild's custom emjoi
	EmojiID *string `json:"emoji_id,omitempty"`
	// EmojiName is the unicode character for the emoji.
	EmojiName *string `json:"emoji_name,omitempty"`
}

type PresenceUpdate struct {
	User         User         `json:"user"`
	GuildID      string       `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
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
	Buttons       []Button     `json:"buttons"`
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
	Desktop string `json:"desktop"`
	Mobile  string `json:"mobile"`
	Web     string `json:"web"`
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
	Tags        string            `json:"tags,omitempty"`
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

type Button struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Style int    `json:"style"`
}

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
