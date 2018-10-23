package cmd

import (
	"../api/news"
	"../bot"
)

// NewsCommand News handler
func NewsCommand(ctx bot.Context) {
	ctx.ReplyEmbed(ctx.Loc("news"), news.GetNews(&ctx))
}
