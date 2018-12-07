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

	"github.com/FlameInTheDark/dtbot/api/location"
	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/fogleman/gg"
)

// Forecast Weather forecast struct
type Forecast struct {
	Cod     string        `json:"cod"`
	Weather []WeatherData `json:"list"`
	City    CityData      `json:"city"`
}

// WeatherData Weather data struct
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

// WDescData Weather description struct
type WDescData struct {
	Id   int64  `json:"id"`
	Main string `json:"main"`
	Desc string `json:"description"`
	Icon string `json:"icon"`
}

// MainData Weather main data struct
type MainData struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity int     `json:"humidity"`
}

// WindData Weather wind data struct
type WindData struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

// CloudsData Weather cloud data struct
type CloudsData struct {
	All int `json:"all"`
}

// CityData Weather city data struct
type CityData struct {
	Name string `json:"name"`
}

// DrawOne returns image with one day forecast
func DrawOne(temp, hum, clo int, time, icon string) image.Image {
	dpc := gg.NewContext(300, 400)
	dpc.SetRGBA(0, 0, 0, 0)
	dpc.Clear()
	// Drawing weather icon

	dpc.Push()
	if err := dpc.LoadFontFace("owfont-regular.ttf", 140); err != nil {
		fmt.Printf("Weather font: %v", err)
	}
	if temp < 0 {
		dpc.SetRGB255(51, 179, 216)
	} else {
		dpc.SetRGB255(231, 81, 56)
	}
	dpc.DrawStringAnchored(icon, 150, 145, 0.5, 0.5)
	dpc.Pop()

	// Drawing lines
	dpc.SetRGB255(76, 90, 109)
	dpc.SetLineWidth(6)
	dpc.DrawLine(299, 61, 299, 400)
	dpc.Stroke()

	// Drawing rectangle
	//dpc.DrawRectangle(0, 0, 300, 60)
	//dpc.SetRGBA(1, 1, 1, 0.6)
	//dpc.Fill()

	// Drawing hummidity and cloudnes
	if err := dpc.LoadFontFace("arial.ttf", 50); err != nil {
		fmt.Printf("Image font: %v", err)
	}
	dpc.SetRGB255(255, 255, 255)
	dpc.DrawStringAnchored(time, 150, 30, 0.5, 0.5)
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
func GetWeatherImage(ctx *bot.Context) (buf *bytes.Buffer, cityName string, err error) {
	var (
		forecast Forecast
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

	cityName = loc.Geonames[0].CountryName + ", " + loc.Geonames[0].Name

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
	//dc := gg.NewContext(1500, 400)
	//dc.SetRGBA255(36, 48, 64, 255)
	//dc.Clear()
	//for i := 0; i < 6; i++ {
	//	dc.DrawImage(DrawOne(int(forecast.Weather[i].Main.TempMin),
	//		forecast.Weather[i].Main.Humidity,
	//		int(forecast.Weather[i].Clouds.All),
	//		fmt.Sprintf("%.2v:00", forecast.Weather[i].TZTime(ctx.Conf.General.Timezone).Hour()),
	//		ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[i].WDesc[0].Id))), 300*i, 0)
	//}

	gc := gg.NewContext(400, 700)
	gc.Clear()

	// Template
	gc.SetRGB255(242, 97, 73)
	gc.DrawRectangle(0, 0, 400, 700)
	gc.Fill()

	gc.SetRGB255(234, 89, 65)
	gc.DrawRectangle(0, 250, 400, 100)
	gc.DrawRectangle(0, 450, 400, 100)
	gc.DrawRectangle(0, 650, 400, 100)
	gc.Fill()

	gc.SetLineWidth(2)
	gc.SetRGBA(0, 0, 0,0.05)
	gc.DrawLine(0, 250, 400, 250)
	gc.DrawLine(0, 348, 400, 348)
	gc.DrawLine(0, 450, 400, 450)
	gc.DrawLine(0, 548, 400, 548)
	gc.DrawLine(0, 650, 400, 650)
	gc.DrawLine(0, 748, 400, 748)
	gc.Stroke()

	// Text
	if err := gc.LoadFontFace("lato.ttf", 20); err != nil {
		panic(err)
	}
	// Header
	gc.SetRGBA(1, 1, 1, 0.7)
	gc.DrawStringAnchored(cityName, 10, 15, 0, 0.5)
	gc.SetRGBA(1, 1, 1, 0.4)
	gc.DrawStringAnchored(time.Now().Format("Jan 2, 2006"), 280, 15, 0, 0.5)

	// First weather data
	gc.SetRGBA(1, 1, 1, 0.5)
	if err := gc.LoadFontFace("lato.ttf", 30); err != nil {
		panic(err)
	}
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Weather[0].TZTime(ctx.Conf.General.Timezone).Hour()), 50, 200, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Weather[0].Main.Humidity), 200, 200, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Weather[0].Clouds.All)), 350, 200, 0.5, 0.5)

	gc.SetRGBA(1, 1, 1, 1)
	if err := gc.LoadFontFace("lato.ttf", 90); err != nil {
		panic(err)
	}

	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Weather[0].Main.TempMin)), 30, 120, 0, 0.5)

	if err := gc.LoadFontFace("owfont-regular.ttf", 90); err != nil {
		panic(err)
	}

	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[0].WDesc[0].Id)), 250, 120, 0, 0.7)

	if err := gc.LoadFontFace("lato.ttf", 30); err != nil {
		panic(err)
	}

	// Time
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Weather[1].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 300, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Weather[2].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 400, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Weather[3].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 500, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Weather[4].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 600, 0, 0.5)

	if err := gc.LoadFontFace("lato.ttf", 50); err != nil {
		panic(err)
	}

	// Temperature
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Weather[1].Main.TempMin)), 250, 300, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Weather[2].Main.TempMin)), 250, 400, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Weather[3].Main.TempMin)), 250, 500, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Weather[4].Main.TempMin)), 250, 600, 0, 0.5)

	if err := gc.LoadFontFace("owfont-regular.ttf", 60); err != nil {
		panic(err)
	}

	// Weather icon
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[1].WDesc[0].Id)), 20, 300, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[2].WDesc[0].Id)), 20, 400, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[3].WDesc[0].Id)), 20, 500, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Weather[4].WDesc[0].Id)), 20, 600, 0, 0.7)

	buf = new(bytes.Buffer)
	pngerr := png.Encode(buf, gc.Image())
	if pngerr != nil {
		fmt.Printf("Image: %v", pngerr)
	}
	return
}
