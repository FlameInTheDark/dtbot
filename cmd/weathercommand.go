package cmd

import (
	"fmt"
	"strings"

	"../api/weather"
	"../bot"
)

// WeatherCommand weather handler
func WeatherCommand(ctx bot.Context) {
	buf, err := weather.GetWeatherImage(&ctx)
	if err != nil {
		bot.NewEmbed("").Color(0xff0000).Field(fmt.Sprintf("%v:", ctx.Loc("weather_error")), err.Error(), false).Send(ctx)
		return
	}
	var city string
	if len(ctx.Args) > 0 {
		city = strings.Join(ctx.Args, " ")
	} else {
		city = ctx.Conf.Weather.City
	}
	ctx.ReplyEmbedAttachment(fmt.Sprintf("%v:", ctx.Loc("weather")), city, "weather.png", buf)
}
