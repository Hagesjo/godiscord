package discordgo

import "encoding/json"

const (
	opCodeDispatch = iota
	opCodeHeartbeat
	opCodeIdentity
	opCodePresenceUpdate
	opCodeVoiceStateUpdate
	opCodeResume
	opCodeReconnect
	opCodeRequestGuildMembers
	opCodeInvalidSession
	opCodeHello
	opCodeHeartbeatAck

	// Gateway close codes
	opCodeUnknownError = iota + 4000
	opCodeUnknownOpCode
	opCodeDecodeError
	opCodeNotAuthenticated
	opCodeAlreadyAuthenticated
	opCodeInvalidSeq
	opCodeRateLimited
	opCodeSessionTimeout
	opCodeInvalidShard
	opCodeShardingRequired
	opCodeInvalidAPIVersion
	opCodeInvalidIntents
	opCodeDisallowedIntents
)

type Event struct {
	OpCode         int              `json:"o,omitempty"`
	SequenceNumber int              `json:"s,omitempty"`
	Name           *string          `json:"name,omitempty"`
	Data           *json.RawMessage `json:"d,omitempty"`
}

type EventHello struct {
	Data struct {
		HeartbeatInterval int `json:"heartbeat_interval,omitempty"`
	} `json:"d,omitempty"`
}

type Timestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

type EventIdentify struct {
	Token      string `json:"token,omitempty"`
	Properties struct {
		OS      string `json:"os,omitempty"`
		Browser string `json:"browser,omitempty"`
		Device  string `json:"device,omitempty"`
	} `json:"properties,omitempty"`
	Compress       *bool     `json:"compress,omitempty"`
	LargeThreshold *int      `json:"large_threshold,omitempty"`
	Shard          []int     `json:"shard,omitempty"`
	Presence       *Presence `json:"presence,omitempty"`
	Intents        int       `json:"intents,omitempty"`
}
