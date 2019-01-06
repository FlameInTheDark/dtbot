package yandexmap

import (
	"bytes"
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/location"
	"github.com/FlameInTheDark/dtbot/bot"
	"io"
	"net/http"
	"strings"
)

func GetMapImage(ctx *bot.Context) (buf *bytes.Buffer, err error) {
	var (
		city     = ctx.GetGuild().WeatherCity
	)

	if len(ctx.Args) > 0 {
		city = strings.Join(ctx.Args, "+")
	}

	loc, err := location.New(ctx.Conf.General.GeonamesUsername, city)
	if err != nil {
		fmt.Printf("Location API: %v", err)
		return
	}

	// Get coordinates and get weather data
	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%v&lon=%v&lang=%v&units=metric&appid=%v",
		newlat, newlng, ctx.Conf.General.Language, ctx.Conf.Weather.WeatherToken))
	if err != nil {
		fmt.Printf("Map API: %v", err)
		return
	}

	buf = new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)//png.Encode(buf, resp.Body)
	if err != nil {
		fmt.Printf("Image: %v", err.Error())
	}
	return buf, err
}