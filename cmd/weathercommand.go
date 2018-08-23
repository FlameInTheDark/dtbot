package cmd

import (
	"../api/weather"
	"../bot"
)

func WeatherCommand(ctx bot.Context) {
	buf, err := weather.GetWeatherImage(&ctx)
	if err != nil {
		ctx.Reply(err.Error())
		return
	}
	ctx.ReplyFile("weather.png", buf)
}
