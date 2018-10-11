package cmd

import (
	"fmt"

	"../bot"
)

// PlayerCommand Player handler
func PlayerCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_no_args"), "", false)
		return
	}
	switch ctx.Args[0] {
	case "play":
		if sess == nil {
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"), "", false)
			return
		}
		player := sess.Player
		go player.Start(sess, ctx.Args[1], func(msg string) {
			ctx.Reply(msg)
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), msg, "", false)
		})
	case "stop":
		sess.Stop()
	case "join":
		if ctx.Sessions.GetByGuild(ctx.Guild.ID) != nil {
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_connected"), "", false)
			return
		}
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"), "", false)
			return
		}
		sess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		})
		if err != nil {
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"), "", false)
			return
		}
		ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID), "", false)
	case "leave":
		if sess == nil {
			ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"), "", false)
			return
		}
		ctx.Sessions.Leave(ctx.Discord, *sess)
		ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_left"), sess.ChannelID), "", false)
	}
}
