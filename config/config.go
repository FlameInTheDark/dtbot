package config

import (
    "os"
    "fmt"
	"github.com/BurntSushi/toml"
)

type WeatherConfig struct {
    WeatherToken    string
    CityZIP         string
    Country         string
}

type GeneralConfig struct {
    Language string
}

type Config struct {
    Weather WeatherConfig
    General GeneralConfig
}

var (
    Weather         WeatherConfig
    BotToken        string
    Language        string
    Locales         LocalesMap
)

// Loading configs from file
func LoadConfig() {
    var cfg Config
    if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
        fmt.Printf("Config loading error. Please create a \"config.toml\"")
        os.Exit(1)
	}
    
    BotToken = "Bot " + os.Getenv("BOT_TOKEN")
    Weather = cfg.Weather
    Language = cfg.General.Language
    LoadLocales()
}