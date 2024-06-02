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

	for _, voiceState := range guildEvent.VoiceStates {
		fetcher.voiceStatesByID[voiceState.UserID] = voiceState
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

func (f *Fetcher) Send(channelID, content string) error {
	return f.restClient.MessageSend(channelID, MessageCreateRequest{
		Content: content,
	})
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

func (f *Fetcher) GetChannelsByIDs(userIDs ...string) (channels []Channel) {
	for _, id := range userIDs {
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
