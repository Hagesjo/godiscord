package events

import "time"

type Hello struct {
	HeartbeatInterval int `json:"heartbeat_interval,omitempty"`
}

type GuildCreate struct {
	Guild

	JoinedAt    time.Time    `json:"joined_at"`
	Large       bool         `json:"large"`
	Unavailable *bool        `json:"unavailable,omitempty"`
	MemberCount int          `json:"member_count"`
	VoiceStates []VoiceState `json:"voice_states"`

	// Members return all guild members.
	// If your bot does not have the GUILD_PRESENCES Gateway Intent,
	// or if the guild has over 75k members,
	// members and presences returned in this event will only contain your bot and users in voice channels.
	Members   []GuildMember    `json:"members"`
	Channels  []Channel        `json:"channels"`
	Threads   []Channel        `json:"threads"`
	Presences []PresenceUpdate `json:"presences"`
}

type Ready struct {
	APIVersion       int     `json:"v"`
	User             User    `json:"user"`
	Guilds           []Guild `json:"guilds"`
	SessionID        string  `json:"session_id"`
	ResumeGatewayURL string  `json:"resume_gateway_url"`
	Shard            []int   `json:"shard,omitempty"`
	// Application is only filled with ID and flag.
	Application Application `json:"application"`
}
