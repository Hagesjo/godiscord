package godiscord

// TODO: Fetcher is not really the name I'm looking for... Context? Taken by stdlib tho.

func newFetcher(guildEvent GuildCreate, restClient *restClient) *Fetcher {
	fetcher := Fetcher{
		guildID:         guildEvent.Guild.ID,
		membersByID:     make(map[string]GuildMember),
		channelsByID:    make(map[string]Channel),
		threadsByID:     make(map[string]Channel),
		voiceStatesByID: make(map[string]VoiceState),
		restClient:      restClient,
	}

	for _, voiceState := range guildEvent.VoiceStates {
		fetcher.voiceStatesByID[voiceState.UserID] = voiceState
	}

	for _, member := range guildEvent.Members {
		member := member
		if member.User == nil {
			continue
		}

		fetcher.membersByID[member.User.ID] = member
		// Discord does not provide member in voice states for GUILD_CREATE for whatever reason.
		// So we fill it from received members.
		if voiceState, ok := fetcher.voiceStatesByID[member.User.ID]; ok {
			voiceState.Member = &member
			fetcher.voiceStatesByID[member.User.ID] = voiceState
		}
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
	guildID         string
	membersByID     map[string]GuildMember
	channelsByID    map[string]Channel
	threadsByID     map[string]Channel
	voiceStatesByID map[string]VoiceState
	restClient      *restClient
}

func (f *Fetcher) SendContent(channelID, content string) error {
	return f.restClient.MessageSend(channelID, MessageCreateRequest{
		Content: content,
	})
}

func (f *Fetcher) SendEmbeds(channelID string, embeds []Embed) error {
	return f.restClient.MessageSend(channelID, MessageCreateRequest{
		Embeds: embeds,
	})
}

func (f *Fetcher) CreateThread(channelID, messageID string, req CreateThreadRequest) error {
	return f.restClient.CreateThread(channelID, messageID, req)
}

func (f *Fetcher) Do(path, method string, res any, resp any) error {
	return f.restClient.do(path, method, res, resp)
}

func (f *Fetcher) CreateChannel(req CreateChannelRequest) error {
	return f.restClient.CreateChannel(f.guildID, req)
}

func (f *Fetcher) GetVoiceStates() []VoiceState {
	return Values(f.voiceStatesByID)
}

func (f *Fetcher) GetMembers() []GuildMember {
	return Values(f.membersByID)
}

func (f *Fetcher) GetMembersByIDs(userIDs ...string) (members []GuildMember) {
	for _, id := range userIDs {
		member, ok := f.membersByID[id]
		if ok {
			members = append(members, member)
		}
	}

	return members
}

func (f *Fetcher) GetChannels() []Channel {
	return Values(f.channelsByID)
}

func (f *Fetcher) GetChannelByID(channelID string) (Channel, bool) {
	c, ok := f.channelsByID[channelID]
	return c, ok
}

func (f *Fetcher) GetChannelByName(name string) (Channel, bool) {
	// TODO: inefficient. If it becomes a problem, save channels by name as well in a dict (and update it on received event).
	cs := Filter(f.GetChannels(), func(channel Channel) bool {
		return channel.Name != nil && *channel.Name == name
	})

	if len(cs) == 0 {
		return Channel{}, false
	}

	return cs[0], true
}

func (f *Fetcher) GetChannelsByIDs(channelIDs ...string) (channels []Channel) {
	for _, id := range channelIDs {
		channel, ok := f.channelsByID[id]
		if ok {
			channels = append(channels, channel)
		}
	}

	return channels
}

func Filter[T any](ts []T, f func(T) bool) (ret []T) {
	for _, t := range ts {
		if f(t) {
			ret = append(ret, t)
		}
	}

	return ret
}

func Values[K comparable, V any](m map[K]V) (vs []V) {
	for _, v := range m {
		vs = append(vs, v)
	}

	return vs
}
