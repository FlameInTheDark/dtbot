package cmd

import (
	"github.com/FlameInTheDark/dtbot/bot"
)

// AlbionCommand handle dice
func AlbionCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		switch ctx.Args[0] {
		case "kills":
			if len(ctx.Args) > 1 {
				ctx.MetricsCommand("albion", "kills")
				ctx.AlbionShowKills()
			}
		case "kill":
			if len(ctx.Args) > 1 {
				ctx.MetricsCommand("albion", "kill")
				ctx.AlbionShowKill()
			}
		case "watch":
			if len(ctx.Args) > 1 {
				err := ctx.Albion.Add(&ctx)
				if err != nil {
					ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_add_error"))
				} else {
					ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_added"))
				}
			}
		}
	}

}
