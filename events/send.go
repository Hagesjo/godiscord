package events

type Heartbeat struct {
	OpCode       int  `json:"op"`
	LastSequence *int `json:"d"`
}

const (
	UserStatusOnline    = "online"
	UserStatusDND       = "dnd"
	UserStatusIdle      = "idle"
	UserStatusInvisible = "invisible"
	UserStatusOffline   = "offline"
)

type Identify struct {
	OpCode  int             `json:"op"`
	Payload IdentifyPayload `json:"d"`
}

type IdentifyPayload struct {
	Token          string     `json:"token"`
	Properties     Properties `json:"properties"`
	Compress       bool       `json:"compress,omitempty"`
	LargeThreshold int        `json:"large_threshold,omitempty"`
	Shard          []int      `json:"shard,omitempty"`
	Presence       *Presence  `json:"presence"`
	Intents        int        `json:"intents"`
}

type Resume struct {
	OpCode  int           `json:"op"`
	Payload ResumePayload `json:"d"`
}

type ResumePayload struct {
	Token        string `json:"token"`
	SessionID    string `json:"session_id"`
	LastSequence uint64 `json:"seq"`
}
