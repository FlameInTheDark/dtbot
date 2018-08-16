package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings"

	"../../config"
	"../location"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

type Forecast struct {
	Cod     string        `json:"cod"`
	Weather []WeatherData `json:"list"`
	City    CityData      `json:"city"`
}

type WeatherData struct {
	Time   int64       `json:"dt"`
	Main   MainData    `json:"main"`
	Wind   WindData    `json:"wind"`
	Clouds CloudsData  `json:"clouds"`
	WDesc  []WDescData `json:"weather"`
}

func (w WeatherData) TZTime() time.Time {
	return time.Unix(w.Time, 0).UTC().Add(time.Hour * time.Duration(config.General.Timezone))
}

type WDescData struct {
	Id   int64  `json:"id"`
	Main string `json:"main"`
	Desc string `json:"description"`
}

type MainData struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity int     `json:"humidity"`
}

type WindData struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type CloudsData struct {
	All int `json:"all"`
}

type CityData struct {
	Name string `json:"name"`
}

func GetForecast(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var (
		forecast      Forecast
		city          string = config.Weather.City
		forecastData  [][]string
		forecastTable bytes.Buffer
	)

	if len(args) > 0 {
		city = strings.Join(args, "+")
	}

	loc, err := location.New(city)
	if err != nil {
		fmt.Printf("Location API: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("location_404"))
		return
	}

	newlat, newlng := loc.GetCoordinates()
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%v&lon=%v&lang=%v&units=metric&appid=%v",
		newlat, newlng, config.General.Language, config.Weather.WeatherToken))
	if err != nil {
		fmt.Printf("Weather API: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_api_error"))
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		fmt.Printf("Weather Decode: %v", err)
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_parse_error"))
		return
	}

	switch forecast.Cod {
	case "404":
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_404"))
		return
	case "200":
		// Generate forecast table
		var tempStr = []string{"°C"}
		var pressureStr = []string{"hPa"}
		var humidityStr = []string{"Hum %"}
		var windStr = []string{"Wind"}
		var cloudsStr = []string{"Clouds"}
		var timeStr = []string{fmt.Sprintf("UTC %v", config.General.Timezone)}
		for i := 0; i < 5; i++ {
			tempStr = append(tempStr, fmt.Sprintf("%v|%v", int(forecast.Weather[i].Main.TempMin), int(forecast.Weather[i].Main.TempMax)))
			pressureStr = append(pressureStr, fmt.Sprintf("%v", int(forecast.Weather[i].Main.Pressure)))
			humidityStr = append(humidityStr, fmt.Sprintf("%.1v", forecast.Weather[i].Main.Humidity))
			windStr = append(windStr, fmt.Sprintf("%.1v", int(forecast.Weather[i].Wind.Speed)))
			cloudsStr = append(cloudsStr, fmt.Sprintf("%.1v", int(forecast.Weather[i].Clouds.All)))
			timeStr = append(timeStr, fmt.Sprintf("%.2v:00", forecast.Weather[i].TZTime().Hour()))
		}
		forecastData = append(forecastData, tempStr)
		forecastData = append(forecastData, pressureStr)
		forecastData = append(forecastData, humidityStr)
		forecastData = append(forecastData, windStr)
		forecastData = append(forecastData, cloudsStr)

		table := tablewriter.NewWriter(&forecastTable)
		table.SetHeader(timeStr)
		table.SetCaption(true, forecast.City.Name)

		for _, v := range forecastData {
			table.Append(v)
		}
		table.Render()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%v```", forecastTable.String()))
		return
	default:
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_error"))
		return
	}
}
