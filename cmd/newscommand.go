package cmd

import (
	"../api/news"
	"../bot"
)

// News command handler
func NewsCommand(ctx bot.Context) {
	ctx.Reply(news.GetNews(&ctx))
}
