package cmd

import (
	"github.com/FlameInTheDark/dtbot/bot"
	"strings"
)

func GreetingCommand(ctx *bot.Context) {
	if len(ctx.Args) > 0 {
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) > 1 {
				ctx.MetricsCommand("greetings", "add")
				ctx.AddGreeting(strings.Join(ctx.Args[1:], " "))
				ctx.ReplyEmbed(ctx.Loc("greetings"), ctx.Loc("greetings_add"))
			} else {
				ctx.MetricsCommand("greetings", "add_no_text")
				ctx.ReplyEmbed(ctx.Loc("greetings"), ctx.Loc("greetings_no_text"))
			}
		case "remove":
			ctx.MetricsCommand("greetings", "remove")
			ctx.RemoveGreeting()
		case "test":
			ctx.ReplyEmbed(ctx.Loc("greetings"), ctx.Loc("greetings_removed"))
		}
	}
}
