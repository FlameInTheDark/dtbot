package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/fun"
	"github.com/FlameInTheDark/dtbot/bot"
)

// SlapCommand returns slap image
func SlapCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		ctx.MetricsCommand("fun", "slap")
		url, err := fun.GetImageURL("slap")
		if err == nil {
			if len(ctx.Args) > 0 {
				var userID string
				if len(ctx.Args[0]) > 3 {
					userID = ctx.Args[0][2 : len(ctx.Args[0])-1]
				}
				user := ctx.GetGuildUser(userID)
				var mention= ctx.Args[0]
				if user != nil {
					mention = user.Username
				}
				ctx.ReplyEmbedAttachmentImageURL("", fmt.Sprintf(ctx.Loc("slapping"), ctx.User.Username, mention), url)
			}
		} else {
			fmt.Printf(err.Error())
		}
	}
}

// FUCommand returns FU image
func FUCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		ctx.MetricsCommand("fun", "fu")
		url, err := fun.GetImageURL("fu")
		if err == nil {
			if len(ctx.Args) > 0 {
				var userID string
				if len(ctx.Args[0]) > 3 {
					userID = ctx.Args[0][2 : len(ctx.Args[0])-1]
				}
				user := ctx.GetGuildUser(userID)
				var mention = ctx.Args[0]
				if user != nil {
					mention = user.Username
				}
				ctx.ReplyEmbedAttachmentImageURL("", fmt.Sprintf(ctx.Loc("send_fu"), ctx.User.Username, mention), url)
			}
		} else {
			fmt.Printf(err.Error())
		}
	}
}