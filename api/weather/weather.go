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

	".github.com/FlameInTheDark/dtbot/bot"
	"github.com/FlameInTheDark/dtbot/api/location"
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

// TZTime returns time in specified timezone
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

// DrawOne returns image with one day forecast
func DrawOne(temp, hum, clo int, time, icon string) image.Image {
	dpc := gg.NewContext(300, 400)
	dpc.SetRGBA(0, 0, 0, 0)
	dpc.Clear()
	dpc.SetRGB(1, 1, 1)

	// Drawing weather icon
	dpc.Push()
	if err := dpc.LoadFontFace("owfont-regular.ttf", 140); err != nil {
		fmt.Printf("Weather font: %v", err)
	}
	dpc.DrawStringAnchored(icon, 150, 145, 0.5, 0.5)
	dpc.Pop()

	// Drawing lines
	dpc.SetLineWidth(1)
	dpc.DrawLine(299, 61, 299, 400)
	dpc.Stroke()
	dpc.DrawLine(0, 400, 0, 61)
	dpc.Stroke()

	// Drawing rectangle
	dpc.DrawRectangle(0, 0, 300, 60)
	dpc.SetRGBA(1, 1, 1, 0.6)
	dpc.Fill()

	// Drawing hummidity and cloudnes
	if err := dpc.LoadFontFace("arial.ttf", 50); err != nil {
		fmt.Printf("Image font: %v", err)
	}
	dpc.SetRGB(256, 256, 256)
	dpc.DrawStringAnchored(time, 150, 30, 0.5, 0.5)
	dpc.SetRGB(1, 1, 1)
	dpc.DrawStringAnchored(fmt.Sprintf("H: %v%%", hum), 150, 305, 0.5, 0.5)
	dpc.DrawStringAnchored(fmt.Sprintf("C: %v%%", clo), 150, 355, 0.5, 0.5)

	// Drawing temperature
	if err := dpc.LoadFontFace("arial.ttf", 80); err != nil {
		fmt.Printf("Image font: %v", err)
	}
	dpc.DrawStringAnchored(fmt.Sprintf("%v°", temp), 150, 230, 0.5, 0.5)

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

	// Get coordinates and get weather data
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

	// Drawing forecast
	dc := gg.NewContext(1500, 400)
	dc.SetRGBA(0, 0, 0, 0.7)
	dc.Clear()
	for i := 0; i < 6; i++ {
		dc.DrawImage(DrawOne(int(forecast.Weather[i].Main.TempMin),
			forecast.Weather[i].Main.Humidity,
			int(forecast.Weather[i].Clouds.All),
			fmt.Sprintf("%.2v:00", forecast.Weather[i].TZTime(ctx.Conf.General.Timezone).Hour()),
			ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[i].WDesc[0].Id))), 300*i, 0)
	}

	buf = new(bytes.Buffer)
	pngerr := png.Encode(buf, dc.Image())
	if pngerr != nil {
		fmt.Printf("Image: %v", pngerr)
	}
	return
}
