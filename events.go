package godiscord

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
)

// Gateway close codes
const (
	CloseEventUnknownError         = 4000
	CloseEventUnknownOpCode        = 4001
	CloseEventDecodeError          = 4002
	CloseEventNotAuthenticated     = 4003
	CloseEventAuthenticationFailed = 4004
	CloseEventAlreadyAuthenticated = 4005
	CloseEventInvalidSeq           = 4007
	CloseEventRateLimited          = 4008
	CloseEventSessionTimeout       = 4009
	CloseEventInvalidShard         = 4010
	CloseEventShardingRequired     = 4011
	CloseEventInvalidAPIVersion    = 4012
	CloseEventInvalidIntents       = 4013
	CloseEventDisallowedIntents    = 4014
)

var ResumeableCloseEvents = map[int]bool{
	CloseEventUnknownError:         true,
	CloseEventUnknownOpCode:        true,
	CloseEventDecodeError:          true,
	CloseEventNotAuthenticated:     true,
	CloseEventAuthenticationFailed: false,
	CloseEventAlreadyAuthenticated: true,
	CloseEventInvalidSeq:           true,
	CloseEventRateLimited:          true,
	CloseEventSessionTimeout:       true,
	CloseEventInvalidShard:         false,
	CloseEventShardingRequired:     false,
	CloseEventInvalidAPIVersion:    false,
	CloseEventInvalidIntents:       false,
	CloseEventDisallowedIntents:    false,
}

type Event struct {
	Type           *string          `json:"t,omitempty"`
	OpCode         int              `json:"op,omitempty"`
	SequenceNumber *uint64          `json:"s,omitempty"`
	Data           *json.RawMessage `json:"d,omitempty"`
}

type DispatchEvent interface {
	Name() string
}
