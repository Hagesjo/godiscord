package discordgo

const (
	IntentGuilds = 1 << iota
	IntentGuildMembers
	IntentGuildModeration
	IntentGuildEmojisAndStickers
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping
	IntentMessageContent
	IntentGuildScheduledEvents
	IntentAutoModerationConfiguration
	IntentAutoModerationExecution
)

// TODO: make some sensible default intents.
