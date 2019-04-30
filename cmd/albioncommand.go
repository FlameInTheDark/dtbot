package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/albion"
	"github.com/FlameInTheDark/dtbot/bot"
)

// AlbionCommand handle dice
func AlbionCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		switch ctx.Args[0] {
		case "kills":
			if len(ctx.Args) > 1 {
				ctx.MetricsCommand("albion", "kills")
				albion.ShowKills(&ctx)
			}

		}
	}

}
