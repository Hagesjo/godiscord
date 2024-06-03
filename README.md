# GoDiscord

GoDiscord is a go wrapper around discord's API. It started out as a fun project for me to challenge myself to implement a comprehensive wrapper using no external dependencies at all. The only dependency is (webgockets)[https://github.com/Hagesjo/webgockets] which I also built from scratch.

You can listen to every guild event (an event with a GuildID field) available (which is the vast majority of all events).

An example:

```go
	bot, err := godiscord.NewBot("<your bot oauth token>", "!")
	if err != nil {
		panic(err)
	}

	bot.RegisterTextCommand("test", func(f *godiscord.Fetcher, args []string, channel godiscord.Channel) error {
		return f.Send(channel.ID, fmt.Sprintf("You sent `test` command with arguments '%v'", args))
	})

	bot.RegisterEventListener(func(f *godiscord.Fetcher, de godiscord.TypingStart) error {
        return f.Send(de.Channel.ID, "You started to type")
	})

	bot.Run()
```

This wrapper is by no means complete, as there's simply too much to cover with the restricted time I have.

Most of the REST api is not covered, but there's a .Do for you to call whatever you want.
