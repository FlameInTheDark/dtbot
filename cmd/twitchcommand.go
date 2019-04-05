package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"strings"
)

// TwitchCommand manipulates twitch announcer
func TwitchCommand(ctx bot.Context) {
	if ctx.IsServerAdmin() {
		if len(ctx.Args) == 0 {
			return
		}
		switch ctx.Args[0] {
		case "add":
			twitchAdd(&ctx)
		case "remove":
			twitchRemove(&ctx)
		case "list":
			twitchList(&ctx)
		case "count":
			twitchCount(&ctx)
		}
	} else {
		ctx.ReplyEmbed("Twitch", ctx.Loc("admin_require"))
		ctx.MetricsCommand("twitch", "error")
	}
}

func twitchAdd(ctx *bot.Context) {
	ctx.MetricsCommand("twitch", "add")
	if len(ctx.Args) > 2 {
		username, err := ctx.Twitch.AddStreamer(ctx.Guild.ID, ctx.Message.ChannelID, ctx.Args[1], strings.Join(ctx.Args[2:], " "))
		if err != nil {
			ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_add_error"))
		} else {
			ctx.ReplyEmbed("Twitch", fmt.Sprintf(ctx.Loc("twitch_added"), username))
		}
	} else if len(ctx.Args) > 1 {
		username, err := ctx.Twitch.AddStreamer(ctx.Guild.ID, ctx.Message.ChannelID, ctx.Args[1], "")
		if err != nil {
			ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_add_error"))
		} else {
			ctx.ReplyEmbed("Twitch", fmt.Sprintf(ctx.Loc("twitch_added"), username))
		}
	}
}

func twitchRemove(ctx *bot.Context) {
	ctx.MetricsCommand("twitch", "remove")
	if len(ctx.Args) > 1 {
		err := ctx.Twitch.RemoveStreamer(ctx.Args[1], ctx.Guild.ID)
		if err != nil {
			ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_remove_error"))
		} else {
			ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_removed"))
		}
	}
}

func twitchList(ctx *bot.Context) {
	ctx.MetricsCommand("twitch", "list")
	if g, ok := ctx.Twitch.Guilds[ctx.Guild.ID]; ok {
		if len(g.Streams) > 0 {
			list := ""
			var counter int
			if g.Streams != nil {
				for _, s := range g.Streams {
					list += fmt.Sprintf("%v. %v\n", counter, s.Login)
					counter++
				}
				ctx.ReplyEmbed("Twitch", fmt.Sprintf(ctx.Loc("twitch_list"), list))
			} else {
				ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_list_empty"))
			}

		} else {
			ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_list_empty"))
		}
	} else {
		ctx.ReplyEmbed("Twitch", ctx.Loc("twitch_list_empty"))
	}
}

func twitchCount(ctx *bot.Context) {
	ctx.MetricsCommand("twitch", "count")
	if ctx.IsAdmin() {
		count := 0
		for _, g := range ctx.Twitch.Guilds {
			count += len(g.Streams)
		}
		ctx.ReplyEmbed("Twitch", fmt.Sprintf("Streamers: %v", count))
	}
}
