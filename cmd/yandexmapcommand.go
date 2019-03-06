package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/yandexmap"
	"github.com/FlameInTheDark/dtbot/bot"
)

// WeatherCommand weather handler
func YandexmapCommand(ctx bot.Context) {
	ctx.MetricsCommand("yandexmap")
	buf, err := yandexmap.GetMapImage(&ctx)
	if err != nil {
		ctx.Log("Map", ctx.Guild.ID, err.Error())
		return
	}
	ctx.ReplyFile("map.png", buf)
}
