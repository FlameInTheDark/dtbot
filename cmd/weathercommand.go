package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/weather"
	"github.com/FlameInTheDark/dtbot/bot"
)

// WeatherCommand weather handler
func WeatherCommand(ctx bot.Context) {
	ctx.MetricsCommand("weather", "main")
	buf, err := weather.GetWeatherImage(&ctx)
	if err != nil {
		ctx.Log("Weather", ctx.Guild.ID, err.Error())
		return
	}
	ctx.ReplyFile("weather.png", buf)
}
