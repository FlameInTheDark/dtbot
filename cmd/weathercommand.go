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
		bot.NewEmbed("").Color(0xff0000).Field(fmt.Sprintf("%v:", ctx.Loc("weather_error")), err.Error(), false).Footer(ctx.Loc("requested_by") + ": " + ctx.User.Username).Send(&ctx)
		return
	}
	ctx.ReplyEmbedAttachment(fmt.Sprintf("%v:", ctx.Loc("weather")), city, "weather.png", buf)
}
