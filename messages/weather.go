package messages

import (
    "fmt"
	"net/http"
    "github.com/bwmarrin/discordgo"
    "encoding/json"
    "../config"
)
// https://github.com/rafamds/discord-bot
type Forecast struct {
    Cod     string          `json:"cod"`
    Weather []WeatherData   `json:"list"`
    City    CityData        `json:"city"`
}

type WeatherData struct {
    Time        int64       `json:"dt"`
    Main        MainData    `json:"main"`
    Wind        WindData    `json:"wind"`
    Clouds      CloudsData  `json:"clouds"`
    WDesc       []WDescData `json:"weather"`
}

type WDescData struct {
    Id      int64   `json:"id"`
    Main    string  `json:"main"`
    Desc    string  `json:"description"`
}

type MainData struct {
    Temp        float64     `json:"temp"`
    Pressure    float64     `json:"pressure"`
    TempMin     float64     `json:"temp_min"`
    TempMax     float64     `json:"temp_max"`
    Humidity    int         `json:"humidity"`
}

type WindData struct {
    Speed       float64     `json:"speed"`
    Deg         float64     `json:"deg"`
}

type CloudsData struct {
    All int     `json:"all"`
}

type CityData struct {
    Name string `json:"name"`
}

// Getting weather from API
func getWeather(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) {
    var (
        forecast Forecast
        zip     string = config.Weather.CityZIP
        country string = config.Weather.Country
    )
    
    if len(args) == 1 {
        zip = args[0]
    } else if len(args) == 2 {
        zip = args[0]
        country = args[1]
    }
    
    resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?zip=%v,%v&lang=%v&units=metric&appid=%v", 
                                      zip, country, config.Language, config.Weather.WeatherToken))
    if err != nil {
        fmt.Println(err)
        s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_api_error"))
        return
    }

    err = json.NewDecoder(resp.Body).Decode(&forecast)
    if err != nil {
        fmt.Println(err)
        s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_parse_error"))
        return
    }
    
    switch forecast.Cod {
    case "404":
        s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_404"))
        return
    case "200":
        response := fmt.Sprintf(config.Locales.Get("weather_format"), 
                                forecast.City.Name, 
                                int(forecast.Weather[0].Main.Temp), 
                                forecast.Weather[0].Main.Pressure, 
                                forecast.Weather[0].Clouds.All, 
                                int(forecast.Weather[0].Wind.Speed), 
                                forecast.Weather[0].Main.Humidity, 
                                forecast.Weather[0].WDesc[0].Desc)
        s.ChannelMessageSend(m.ChannelID, response)
        return
    default:
        s.ChannelMessageSend(m.ChannelID, config.Locales.Get("weather_error"))
        return
    }
}