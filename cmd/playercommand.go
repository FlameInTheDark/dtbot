package cmd

import (
	"../bot"
)

// Player command handler
func PlayerCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.Reply("Already connected! Use `music leave` for the bot to disconnect.")
		return
	}
	switch ctx.Args[0] {
	case "play":
		if sess == nil {
			ctx.Reply("Not in a voice channel! To make the bot join one, use `!radio join`.")
			return
		}
		player := sess.Player
		go player.Start(sess, ctx.Args[1], func(msg string) {
			ctx.Reply(msg)
		})
	case "stop":
		sess.Stop()
	case "join":
		if ctx.Sessions.GetByGuild(ctx.Guild.ID) != nil {
			ctx.Reply("Already connected! Use `music leave` for the bot to disconnect.")
			return
		}
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.Reply("You must be in a voice channel to use the bot!")
			return
		}
		sess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		})
		if err != nil {
			ctx.Reply("An error occured!")
			return
		}
		ctx.Reply("Joined <#" + sess.ChannelId + ">!")
	case "leave":
		if sess == nil {
			ctx.Reply("Not in a voice channel! To make the bot join one, use `music join`.")
			return
		}
		ctx.Sessions.Leave(ctx.Discord, *sess)
		ctx.Reply("Left <#" + sess.ChannelId + ">!")
	}
}
