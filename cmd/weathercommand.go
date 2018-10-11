package cmd

import (
	"fmt"

	"../api/weather"
	"../bot"
)

// WeatherCommand weather handler
func WeatherCommand(ctx bot.Context) {
	buf, err := weather.GetWeatherImage(&ctx)
	if err != nil {
		ctx.Reply(fmt.Sprintf("%v: %v", ctx.Loc("weather_error"), err.Error()))
		return
	}
	ctx.ReplyFile("weather.png", buf)
}
