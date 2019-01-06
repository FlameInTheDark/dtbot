package yandexmap

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/yageocoding"
	"github.com/FlameInTheDark/dtbot/bot"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetMapImage(ctx *bot.Context) (buf *bytes.Buffer, err error) {
	var (
		mapType = "map"
		mapSize = 14
		city    = ctx.GetGuild().WeatherCity
	)

	if len(ctx.Args) > 2 {
		mapType = ctx.Args[0]
		mapSize, _ = strconv.Atoi(ctx.Args[1])
		city = strings.Join(ctx.Args[2:], "+")
	}

	if mapSize > 17 {
		return buf, errors.New("map size out of range")
	}

	loc, err := yageocoding.GetData(ctx.Conf.General.GeocodingApiKey, city) //location.New(ctx.Conf.General.GeonamesUsername, city)
	if err != nil {
		fmt.Printf("Location API: %v", err)
		return
	}

	// Get coordinates and get weather data
	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%v,%v&size=450,450&z=%v&l=%v&pt=%v,%v,vkbkm",
		newlat, newlng, mapSize, mapType, newlat, newlng))
	if err != nil {
		fmt.Printf("Map API: %v", err)
		return
	}

	buf = new(bytes.Buffer)

	_, err = io.Copy(buf, resp.Body) //png.Encode(buf, resp.Body)
	if err != nil {
		fmt.Printf("Map image: %v", err.Error())
	}
	return buf, err
}
