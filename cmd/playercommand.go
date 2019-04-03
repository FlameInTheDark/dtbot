package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/bot"
)

// PlayerCommand Player handler
func PlayerCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_no_args"))
		return
	}
	switch ctx.Args[0] {
	case "play":
		ctx.MetricsCommand("player", "play")
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
			return
		}
		if len(ctx.Args) > 1 {
			go sess.Player.Start(sess, ctx.Args[1], func(msg string) {
				ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
			}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
			fmt.Printf("volume=%.3f", ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
		}
	case "stop":
		ctx.MetricsCommand("player", "stop")
		sess.Stop()
	}
}
