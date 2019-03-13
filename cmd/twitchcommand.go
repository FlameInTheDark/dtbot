package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
)

// TwitchCommand manipulates twitch announcer
func TwitchCommand(ctx bot.Context) {
	ctx.MetricsCommand("twitch")
	if ctx.GetRoles().ExistsName("bot.admin") || ctx.IsAdmin() {
		if len(ctx.Args) == 0 {
			return
		}
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) > 1 {
				username, err := ctx.Twitch.AddStreamer(ctx.Guild.ID, ctx.Message.ChannelID, ctx.Args[1])
				if err != nil {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_add_error"))
				} else {
					ctx.ReplyEmbed("Twitch", fmt.Sprintf(ctx.Loc("twitch_added"), username))
				}
			}
		case "remove":
			if len(ctx.Args) > 1 {
				err := ctx.Twitch.RemoveStreamer(ctx.Args[1], ctx.Guild.ID)
				if err != nil {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_remove_error"))
				} else {
					ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_removed"))
				}
			}
		case "count":
			if ctx.IsAdmin() {
				count := 0
				for _,g := range ctx.Twitch.Guilds {
					count += len(g.Streams)
				}
				ctx.ReplyEmbed("Twitch", fmt.Sprintf("Streamers: %v", count))
			}
		}
	} else {
		ctx.ReplyEmbed("Twitch", ctx.Loc("admin_require"))
	}
}