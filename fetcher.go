package discordgo

import "github.com/hagesjo/discordgo/events"

func newFetcher(guildEvent events.GuildCreate) *Fetcher {
	fetcher := Fetcher{
		membersByID:  make(map[string]events.GuildMember),
		channelsByID: make(map[string]events.Channel),
		threadsByID:  make(map[string]events.Channel),
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
	membersByID  map[string]events.GuildMember
	channelsByID map[string]events.Channel
	threadsByID  map[string]events.Channel
}
