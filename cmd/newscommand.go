package cmd

import (
	"../api/news"
	"../bot"
)

func NewsCommand(ctx bot.Context) {
	ctx.Reply(news.GetNews(&ctx))
}
