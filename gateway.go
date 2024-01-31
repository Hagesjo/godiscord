package discordgo

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
