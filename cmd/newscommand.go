package cmd

import (
	"../api/news"
	"../bot"
)

// NewsCommand News handler
func NewsCommand(ctx bot.Context) {
	ctx.Reply(news.GetNews(&ctx))
}
