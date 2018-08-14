package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type WeatherConfig struct {
	WeatherToken string
	City         string
}

type GeneralConfig struct {
	Language        string
    Timezone        int
    GeocodingKey    string
}

type Config struct {
	Weather WeatherConfig
	General GeneralConfig
}

var (
	Weather  WeatherConfig
	BotToken    string
	Locales     LocalesMap
    General     GeneralConfig
)

// Loading configs from file
func LoadConfig() {
	var cfg Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		fmt.Printf("Config loading error: %v\n", err)
		os.Exit(1)
	}

	BotToken = "Bot " + os.Getenv("BOT_TOKEN")
	Weather = cfg.Weather
    General = cfg.General
	LoadLocales()
}
