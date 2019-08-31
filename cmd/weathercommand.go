package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/darksky"
	"github.com/FlameInTheDark/dtbot/bot"
)

// WeatherCommand weather handler
func WeatherCommand(ctx bot.Context) {
	ctx.MetricsCommand("weather", "main")
	buf, err := darksky.GetWeatherImage(&ctx)
	if err != nil {
		ctx.Log("Weather", ctx.Guild.ID, err.Error())
		return
	}
	ctx.ReplyFile("weather.png", buf)
}

// WeatherCommand weather handler
func WeatherWeekCommand(ctx bot.Context) {
	ctx.MetricsCommand("weather", "week")
	buf, err := darksky.GetWeatherWeekImage(&ctx)
	if err != nil {
		ctx.Log("Weather", ctx.Guild.ID, err.Error())
		return
	}
	ctx.ReplyFile("weather.png", buf)
}