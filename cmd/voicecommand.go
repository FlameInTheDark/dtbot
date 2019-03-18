package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/bot"
)

// VoiceCommand voice handler
func VoiceCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	switch ctx.Args[0] {
	case "join":
		ctx.MetricsCommand("voice", "join")
		if ctx.Sessions.GetByGuild(ctx.Guild.ID) != nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_connected"))
			return
		}
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		sess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		})
		if err != nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
			return
		}
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID))
	case "leave":
		ctx.MetricsCommand("voice", "leave")
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		ctx.Sessions.Leave(ctx.Discord, *sess)
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_left"), sess.ChannelID))
	}
}
