package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/fun"
	"github.com/FlameInTheDark/dtbot/bot"
)

// SlapCommand returns slap image
func SlapCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		url, err := fun.GetImageURL("slap")
		if err == nil {
			ctx.ReplyEmbedAttachmentImageURL(fmt.Sprintf("%v slaping %v", ctx.User.Mention(), ctx.Args[0]), url)
		} else {
			fmt.Printf(err.Error())
		}
	}
}
