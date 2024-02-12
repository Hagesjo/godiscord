package discordgo

// TODO: Fetcher is not really the name I'm looking for... Context? Taken by stdlib tho.

func newFetcher(guildEvent GuildCreate, restClient *restClient) *Fetcher {
	fetcher := Fetcher{
		membersByID:  make(map[string]GuildMember),
		channelsByID: make(map[string]Channel),
		threadsByID:  make(map[string]Channel),
		restClient:   restClient,
	}
	for _, member := range guildEvent.Members {
		if member.User == nil {
			continue
		}

		fetcher.membersByID[member.User.ID] = member
	}

	for _, channel := range guildEvent.Channels {
		fetcher.channelsByID[channel.ID] = channel
	}

	for _, thread := range guildEvent.Threads {
		fetcher.threadsByID[thread.ID] = thread
	}

	return &fetcher
}

type Fetcher struct {
	guildID      string
	membersByID  map[string]GuildMember
	channelsByID map[string]Channel
	threadsByID  map[string]Channel
	restClient   *restClient
}

func (f *Fetcher) Send(channelID, content string) error {
	return f.restClient.MessageSend(channelID, MessageCreateRequest{
		Content: content,
	})
}
