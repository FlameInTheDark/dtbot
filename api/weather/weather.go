package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"strings"
	"time"

	"../../bot"
	"../location"
	"github.com/fogleman/gg"
)

// Forecast : Weather forecast struct
type Forecast struct {
	Cod     string        `json:"cod"`
	Weather []WeatherData `json:"list"`
	City    CityData      `json:"city"`
}

// WeatherData : Weather data struct
type WeatherData struct {
	Time   int64       `json:"dt"`
	Main   MainData    `json:"main"`
	Wind   WindData    `json:"wind"`
	Clouds CloudsData  `json:"clouds"`
	WDesc  []WDescData `json:"weather"`
}

// Returns time in specified timezone
func (w WeatherData) TZTime(tz int) time.Time {
	return time.Unix(w.Time, 0).UTC().Add(time.Hour * time.Duration(tz))
}

// WDescData : Weather description struct
type WDescData struct {
	Id   int64  `json:"id"`
	Main string `json:"main"`
	Desc string `json:"description"`
	Icon string `json:"icon"`
}

// MainData : Weather main data struct
type MainData struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity int     `json:"humidity"`
}

// WindData : Weather wind data struct
type WindData struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

// CloudsData : Weather cloud data struct
type CloudsData struct {
	All int `json:"all"`
}

// CityData : Weather city data struct
type CityData struct {
	Name string `json:"name"`
}

// DrawOne draw one day
func DrawOne(temp, hum, clo int, time, icon string) image.Image {
	dpc := gg.NewContext(300, 400)
	dpc.SetRGBA(0, 0, 0, 0)
	dpc.Clear()
	dpc.SetRGB(1, 1, 1)

	res, err := http.Get(fmt.Sprintf("http://openweathermap.org/img/w/%v.png", icon))
	if err != nil || res.StatusCode != 200 {
		fmt.Println(err)
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	dpc.Push()
	dpc.Scale(3, 3)
	dpc.DrawImage(m, 25, 12)
	dpc.Pop()

	if err := dpc.LoadFontFace("arial.ttf", 50); err != nil {
		fmt.Printf("Image font: %v", err)
	}
	dpc.DrawStringAnchored(time, 150, 30, 0.5, 0.5)
	dpc.DrawStringAnchored(fmt.Sprintf("H: %v%%", hum), 150, 280, 0.5, 0.5)
	dpc.DrawStringAnchored(fmt.Sprintf("C: %v%%", clo), 150, 330, 0.5, 0.5)

	if err := dpc.LoadFontFace("arial.ttf", 80); err != nil {
		fmt.Printf("Image font: %v", err)
	}
	dpc.DrawStringAnchored(fmt.Sprintf("%v°", temp), 150, 200, 0.5, 0.5)

	return dpc.Image()
}

// GetWeatherImage returns buffer with weather image
func GetWeatherImage(ctx *bot.Context) (buf *bytes.Buffer, err error) {
	var (
		forecast Forecast
		city     string = ctx.Conf.Weather.City
	)

	if len(ctx.Args) > 0 {
		city = strings.Join(ctx.Args, "+")
	}

	loc, err := location.New(ctx.Conf.General.GeonamesUsername, city)
	if err != nil {
		fmt.Printf("Location API: %v", err)
		return
	}

	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%v&lon=%v&lang=%v&units=metric&appid=%v",
		newlat, newlng, ctx.Conf.General.Language, ctx.Conf.Weather.WeatherToken))
	if err != nil {
		fmt.Printf("Weather API: %v", err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		fmt.Printf("Weather Decode: %v", err)
		return
	}

	dc := gg.NewContext(1500, 400)
	dc.SetRGBA(0, 0, 0, 0.7)
	dc.Clear()
	for i := 0; i < 6; i++ {
		dc.DrawImage(DrawOne(int(forecast.Weather[i].Main.TempMin),
			forecast.Weather[i].Main.Humidity,
			int(forecast.Weather[i].Clouds.All),
			fmt.Sprintf("%.2v:00", forecast.Weather[i].TZTime(ctx.Conf.General.Timezone).Hour()),
			forecast.Weather[i].WDesc[0].Icon), 300*i, 0)
	}

	buf = new(bytes.Buffer)
	pngerr := png.Encode(buf, dc.Image())
	if pngerr != nil {
		fmt.Printf("Image: %v", pngerr)
	}
	return
}
