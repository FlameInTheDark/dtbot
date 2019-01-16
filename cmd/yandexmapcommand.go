package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/yandexmap"
	"github.com/FlameInTheDark/dtbot/bot"
)

// WeatherCommand weather handler
func YandexmapCommand(ctx bot.Context) {
	buf, err := yandexmap.GetMapImage(&ctx)
	if err != nil {
		ctx.DB.Log("Map", ctx.Guild.ID, err.Error())
		return
	}
	ctx.ReplyFile("map.png",buf)
	ctx.MetricsCommand("yandexmap")
}