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
			user := ctx.GetGuildUser(ctx.Args[0][2 : len(ctx.Args[0])-1])
			var mention = ctx.Args[0]
			if user != nil {
				mention = user.Username
			}
			ctx.ReplyEmbedAttachmentImageURL("", fmt.Sprintf("%v slaping %v", ctx.User.Username, mention), url)
		} else {
			fmt.Printf(err.Error())
		}
	}
}
