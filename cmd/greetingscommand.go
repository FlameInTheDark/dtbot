package cmd

import (
	"github.com/FlameInTheDark/dtbot/bot"
	"strings"
)

func GreetingsCommand(ctx bot.Context) {
	if ctx.GetRoles().ExistsName("bot.admin") || ctx.IsAdmin() {
		if len(ctx.Args) > 0 {
			switch ctx.Args[0] {
			case "add":
				if len(ctx.Args) > 1 {
					ctx.MetricsCommand("greetings", "add")
					ctx.AddGreetings(strings.Join(ctx.Args[1:], " "))
					ctx.ReplyEmbed(ctx.Loc("greetings"), ctx.Loc("greetings_add"))
				} else {
					ctx.MetricsCommand("greetings", "add_no_text")
					ctx.ReplyEmbed(ctx.Loc("greetings"), ctx.Loc("greetings_no_text"))
				}
			case "remove":
				ctx.MetricsCommand("greetings", ctx.Loc("greetings_removed"))
				ctx.RemoveGreetings()
			case "test":
				_ = ctx.ReplyPM(ctx.Guilds.Guilds[ctx.Guild.ID].Greeting)
			}
		}
	} else {
		ctx.ReplyEmbed("Greetings", ctx.Loc("admin_require"))
		ctx.MetricsCommand("greetings", "error")
	}
}
