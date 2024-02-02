package discordgo

const (
	UserStatusOnline    = "online"
	UserStatusDND       = "dnd"
	UserStatusIdle      = "idle"
	UserStatusInvisible = "invisible"
	UserStatusOffline   = "offline"
)

type Emoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

type Party struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}

type Assets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

type Button struct {
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Activity struct {
	// Bot users are only able to set these four below.

	Name  string  `json:"name,omitempty"`
	Type  int     `json:"type,omitempty"`
	URL   *string `json:"url,omitempty"`
	State *string `json:"state,omitempty"`

	CreatedAt     int64       `json:"created_at,omitempty"`
	Timestamps    *Timestamps `json:"timestamps,omitempty"`
	ApplicationID string      `json:"application_id,omitempty"`
	Details       *string     `json:"details,omitempty"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	Party         *Party      `json:"party,omitempty"`
	Assets        *Assets     `json:"assets,omitempty"`
	Secrets       *Secrets    `json:"secrets,omitempty"`
	Instance      bool        `json:"instance,omitempty"`
	Flags         int         `json:"flags,omitempty"`
	Buttons       []Button    `json:"buttons,omitempty"`
}

type Presence struct {
	Since      *int       `json:"since,omitempty"`
	Activities []Activity `json:"activities,omitempty"`
	Status     string     `json:"status,omitempty"`
	AFK        bool       `json:"afk,omitempty"`
}
