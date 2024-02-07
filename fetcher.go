package discordgo

import "github.com/hagesjo/discordgo/events"

func newFetcher(guildEvent events.GuildCreate) *Fetcher {
	fetcher := Fetcher{
		membersByGuild:  make(map[string]events.GuildMember),
		channelsByGuild: make(map[string]events.Channel),
		threadsByGuild:  make(map[string]events.Channel),
	}
	for _, member := range guildEvent.Members {
		if member.User == nil {
			continue
		}

		fetcher.membersByGuild[member.User.ID] = member
	}

	for _, channel := range guildEvent.Channels {
		fetcher.channelsByGuild[channel.ID] = channel
	}

	for _, thread := range guildEvent.Threads {
		fetcher.channelsByGuild[thread.ID] = thread
	}

	return &fetcher
}

type Fetcher struct {
	guildID         string
	membersByGuild  map[string]events.GuildMember
	channelsByGuild map[string]events.Channel
	threadsByGuild  map[string]events.Channel
}
