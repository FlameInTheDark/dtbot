package cmd

import (
	"github.com/FlameInTheDark/dtbot/bot"
)

// TwitchCommand manipulates twitch announcer
func TwitchCommand(ctx bot.Context) {
	ctx.MetricsCommand("bot")
	if ctx.GetRoles().ExistsName("bot.admin") || ctx.IsAdmin() {
		if len(ctx.Args) == 0 {
			return
		}
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) > 1 {
				err := ctx.Twitch.AddStreamer(ctx.Guild.ID, ctx.Message.ChannelID, ctx.Args[1])
				if err != nil {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_add_error"))
				} else {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_added"))
				}
			}
		case "remove":
			if len(ctx.Args) > 1 {
				err := ctx.Twitch.RemoveStreamer(ctx.Guild.ID, ctx.Args[1])
				if err != nil {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_remove_error"))
				} else {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_removed"))
				}
			}
		}
	} else {
		ctx.ReplyEmbed("Twitch", ctx.Loc("admin_require"))
	}
}