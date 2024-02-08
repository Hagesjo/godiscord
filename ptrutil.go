package discordgo

func strPtr(s string) *string {
	return &s
}

func valueOrZero[T any](t *T) T {
	if t != nil {
		return *t
	}

	var ret T
	return ret
}
