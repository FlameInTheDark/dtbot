package cmd

import (
	"fmt"

	"../bot"
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
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
			return
		}
		player := sess.Player
		go player.Start(sess, ctx.Args[1], func(msg string) {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
		})
	case "stop":
		sess.Stop()
	}
}
