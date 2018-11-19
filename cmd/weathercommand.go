package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/weather"
	"github.com/FlameInTheDark/dtbot/bot"
)

// WeatherCommand weather handler
func WeatherCommand(ctx bot.Context) {
	buf, city, err := weather.GetWeatherImage(&ctx)
	if err != nil {
		ctx.DB.Log("Weather",err.Error())
		return
	}
	ctx.ReplyEmbedAttachment(fmt.Sprintf("%v:", ctx.Loc("weather")), city, "weather.png", buf)
}
