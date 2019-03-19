package darksky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FlameInTheDark/dtbot/api/location"
	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/fogleman/gg"
	"image/png"
	"net/http"
	"strings"
	"time"
)

type DarkSkyResponse struct {
	Latitude  float32       `json:"latitude"`
	Longitude float32       `json:"longitude"`
	Timezone  string        `json:"timezone"`
	Currently DarkSkyData   `json:"currently"`
	Hourly    DarkSkyHourly `json:"hourly"`
	Daily     DarkSkyDaily  `json:"daily"`
	Flags     DarkSkyFlags  `json:"flags"`
	Offset    int64           `json:"offset"`
}

type DarkSkyData struct {
	Time                int64     `json:"time"`
	Summary             string  `json:"summary"`
	Icon                string  `json:"icon"`
	PrecipIntensity     int64     `json:"precipIntensity"`
	PrecipProbability   int64     `json:"precipProbability"`
	Temperature         float32 `json:"temperature"`
	ApparentTemperature float32 `json:"apparentTemperature"`
	DewPoint            float32 `json:"dewPoint"`
	Humidity            float32 `json:"humidity"`
	Pressure            float32 `json:"pressure"`
	WindSpeed           float32 `json:"windSpeed"`
	WindGust            float32 `json:"windGust"`
	WindBearing         int64     `json:"windBearing"`
	CloudCover          float32 `json:"cloudCover"`
	UVIndex             int64     `json:"uvIndex"`
	Visibility          float32 `json:"visibility"`
	Ozone               float32 `json:"ozone"`
}

type DarkSkyHourly struct {
	Summary string        `json:"summary"`
	Icon    string        `json:"icon"`
	Data    []DarkSkyData `json:"data"`
}

type DarkSkyDaily struct {
	Summary string           `json:"summary"`
	Icon    string           `json:"icon"`
	Data    []DarkSkyDayData `json:"data"`
}

type DarkSkyDayData struct {
	Time                        int64     `json:"time"`
	Summary                     string  `json:"summary"`
	Icon                        string  `json:"icon"`
	SunriseTime                 int64     `json:"sunriseTime"`
	SunsetTime                  int64     `json:"sunsetTime"`
	MoonPhase                   float32 `json:"moonPhase"`
	PrecipIntensity             float32 `json:"precipIntensity"`
	PrecipIntensityMax          float32 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime      int64     `json:"precipIntensityMaxTime"`
	PrecipProbability           float32 `json:"precipProbability"`
	PrecipAccumulation          float32 `json:"precipAccumulation"`
	PrecipType                  string  `json:"precipType"`
	TemperatureHigh             float32 `json:"temperatureHigh"`
	TemperatureHighTime         int64     `json:"temperatureHighTime"`
	TemperatureLow              float32 `json:"temperatureLow"`
	TemperatureLowTime          int64     `json:"temperatureLowTime"`
	ApparentTemperatureHigh     float32 `json:"apparentTemperatureHigh"`
	ApparentTemperatureHighTime int64     `json:"apparentTemperatureHighTime"`
	ApparentTemperatureLow      float32 `json:"apparentTemperatureLow"`
	ApparentTemperatureLowTime  int64     `json:"apparentTemperatureLowTime"`
	DewPoint                    float32 `json:"dewPoint"`
	Humidity                    float32 `json:"humidity"`
	Pressure                    float32 `json:"pressure"`
	WindSpeed                   float32 `json:"windSpeed"`
	WindGust                    float32 `json:"windGust"`
	WindGustTime                int64     `json:"windGustTime"`
	WindBearing                 int64     `json:"windBearing"`
	CloudCover                  float32 `json:"cloudCover"`
	UVIndex                     int64     `json:"uvIndex"`
	UVIndexTime                 int64     `json:"uvIndexTime"`
	Visibility                  float32 `json:"visibility"`
	Ozone                       float32 `json:"ozone"`
	TemperatureMin              float32 `json:"temperatureMin"`
	TemperatureMinTime          int64     `json:"temperatureMinTime"`
	TemperatureMax              float32 `json:"temperatureMax"`
	TemperatureMaxTime          int64     `json:"temperatureMaxTime"`
	ApparentTemperatureMin      float32 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime  int64     `json:"apparentTemperatureMinTime"`
	ApparentTemperatureMax      float32 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime  int64     `json:"apparentTemperatureMaxTime"`
}

type DarkSkyFlags struct {
	Sources        []string `json:"sources"`
	NearestStation float32  `json:"nearest-station"`
	Units          string   `json:"units"`
}

func (d DarkSkyData) TZTime(tz int) time.Time {
	return time.Unix(d.Time, 0).UTC().Add(time.Hour * time.Duration(tz))
}

func GetWeatherImage(ctx *bot.Context) (buf *bytes.Buffer, err error) {
	var (
		forecast DarkSkyResponse
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

	cityName := loc.Geonames[0].CountryName + ", " + loc.Geonames[0].Name

	// Get coordinates and get weather data
	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://api.darksky.net/forecast/%v/%v,%v?units=ca&lang=%v",
		ctx.Conf.DarkSky.Token, newlat, newlng, ctx.Conf.General.Language))
	if err != nil {
		fmt.Printf("Weather API: %v", err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		fmt.Printf("Weather Decode: %v", err)
		return
	}

	gc := gg.NewContext(400, 650)
	gc.SetRGBA(0, 0, 0, 0)
	gc.Clear()

	// Template
	gc.SetRGB255(242, 97, 73)
	gc.DrawRoundedRectangle(0, 0, 400, 650, 10)
	gc.Fill()

	// Weather lines
	gc.SetRGB255(234, 89, 65)
	gc.DrawRectangle(0, 250, 400, 100)
	gc.DrawRectangle(0, 450, 400, 100)
	gc.Fill()

	gc.SetLineWidth(2)
	gc.SetRGBA(0, 0, 0, 0.05)
	gc.DrawLine(0, 250, 400, 250)
	gc.DrawLine(0, 349, 400, 348)
	gc.DrawLine(0, 450, 400, 450)
	gc.DrawLine(0, 549, 400, 548)
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
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Currently.TZTime(ctx.Conf.General.Timezone).Hour()), 50, 200, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Currently.Humidity), 200, 200, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Currently.CloudCover)), 350, 200, 0.5, 0.5)

	gc.SetRGBA(1, 1, 1, 1)
	if err := gc.LoadFontFace("lato.ttf", 90); err != nil {
		panic(err)
	}

	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Currently.Temperature)), 100, 120, 0.5, 0.5)

	if err := gc.LoadFontFace("owfont-regular.ttf", 90); err != nil {
		panic(err)
	}

	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Currently.Icon)), 250, 120, 0, 0.7)

	if err := gc.LoadFontFace("lato.ttf", 30); err != nil {
		panic(err)
	}

	// Time
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Hourly.Data[1].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 285, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Hourly.Data[2].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 385, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Hourly.Data[3].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 485, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%.2v:00", forecast.Hourly.Data[4].TZTime(ctx.Conf.General.Timezone).Hour()), 100, 585, 0, 0.5)

	// Humidity and cloudiness
	if err := gc.LoadFontFace("lato.ttf", 20); err != nil {
		panic(err)
	}
	gc.SetRGBA(1, 1, 1, 0.5)

	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Hourly.Data[1].Humidity), 100, 315, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Hourly.Data[2].Humidity), 100, 415, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Hourly.Data[3].Humidity), 100, 515, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("H:%v%%", forecast.Hourly.Data[4].Humidity), 100, 615, 0, 0.5)

	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Hourly.Data[1].CloudCover)), 170, 315, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Hourly.Data[2].CloudCover)), 170, 415, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Hourly.Data[3].CloudCover)), 170, 515, 0, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("C:%v%%", int(forecast.Hourly.Data[4].CloudCover)), 170, 615, 0, 0.5)

	gc.SetRGBA(1, 1, 1, 1)
	if err := gc.LoadFontFace("lato.ttf", 50); err != nil {
		panic(err)
	}

	// Temperature
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Hourly.Data[1].Temperature)), 320, 300, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Hourly.Data[2].Temperature)), 320, 400, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Hourly.Data[3].Temperature)), 320, 500, 0.5, 0.5)
	gc.DrawStringAnchored(fmt.Sprintf("%v°", int(forecast.Hourly.Data[4].Temperature)), 320, 600, 0.5, 0.5)

	if err := gc.LoadFontFace("owfont-regular.ttf", 60); err != nil {
		panic(err)
	}

	// Weather icon
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Hourly.Data[1].Icon)), 20, 300, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Hourly.Data[2].Icon)), 20, 400, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Hourly.Data[3].Icon)), 20, 500, 0, 0.7)
	gc.DrawStringAnchored(ctx.WeatherCode(fmt.Sprintf("%v", forecast.Hourly.Data[4].Icon)), 20, 600, 0, 0.7)

	buf = new(bytes.Buffer)
	pngerr := png.Encode(buf, gc.Image())
	if pngerr != nil {
		fmt.Printf("Image: %v", pngerr)
	}
	return
}