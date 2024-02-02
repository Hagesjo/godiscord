package events

import "encoding/json"

const (
	OpCodeDispatch            = 0
	OpCodeHeartbeat           = 1
	OpCodeIdentity            = 2
	OpCodePresenceUpdate      = 3
	OpCodeVoiceStateUpdate    = 5
	OpCodeResume              = 6
	OpCodeReconnect           = 7
	OpCodeRequestGuildMembers = 8
	OpCodeInvalidSession      = 9
	OpCodeHello               = 10
	OpCodeHeartbeatAck        = 11

	// Gateway close codes
	OpCodeUnknownError = iota + 4000
	OpCodeUnknownOpCode
	OpCodeDecodeError
	OpCodeNotAuthenticated
	OpCodeAlreadyAuthenticated
	OpCodeInvalidSeq
	OpCodeRateLimited
	OpCodeSessionTimeout
	OpCodeInvalidShard
	OpCodeShardingRequired
	OpCodeInvalidAPIVersion
	OpCodeInvalidIntents
	OpCodeDisallowedIntents
)

type Event struct {
	OpCode         int              `json:"op,omitempty"`
	SequenceNumber int              `json:"s,omitempty"`
	Name           *string          `json:"name,omitempty"`
	Data           *json.RawMessage `json:"d,omitempty"`
}
