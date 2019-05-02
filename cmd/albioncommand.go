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
				ctx.MetricsCommand("albion", "watch")
				err := ctx.AlbionAddPlayer()
				if err != nil {
					ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_add_error"))
				} else {
					ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_added"))
				}
			}
		case "unwatch":
			if _,ok := ctx.Albion.Players[ctx.User.ID]; ok {
				ctx.MetricsCommand("albion", "unwatch")
				delete(ctx.Albion.Players, ctx.User.ID)
				ctx.DB.RemoveAlbionPlayer(ctx.User.ID)
				ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_removed"))
			} else {
				ctx.ReplyEmbed("Albion Killboard", ctx.Loc("albion_not_watching"))
			}
		}
	}

}
