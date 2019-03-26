package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/news"
	"github.com/FlameInTheDark/dtbot/bot"
)

// NewsCommand News handler
func NewsCommand(ctx bot.Context) {
	ctx.MetricsCommand("news", "main")
	_ = news.GetNews(&ctx)
}
