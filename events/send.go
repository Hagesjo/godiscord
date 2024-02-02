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
	Presence       *Presence  `json:"presence,omitempty"`
	Intents        int        `json:"intents"`
}

type Properties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type Presence struct {
	Since      *int64      `json:"since,omitempty"`
	Activities []*Activity `json:"activities,omitempty"`
	Status     string      `json:"status"`
	AFK        bool        `json:"afk"`
}

type Activity struct {
	Name string       `json:"name"`
	Type ActivityType `json:"type"`
	URL  string       `json:"url,omitempty"`
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
