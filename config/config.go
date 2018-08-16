package config

import (
	"fmt"
	"os"
	"github.com/BurntSushi/toml"

)

type WeatherConfig struct {
	WeatherToken string
	City         string
}

type GeneralConfig struct {
	Language			string
	Timezone			int
	GeonamesUsername	string
}

type NewsConfig struct {
	ApiKey		string
	Country		string
	Articles	int
}

type TranslateConfig struct {
	ApiKey	string
}

type Config struct {
	Weather 	WeatherConfig
	General 	GeneralConfig
	News		NewsConfig
	Translate	TranslateConfig
}

var (
	BotToken	string
	Weather		WeatherConfig
	Locales		LocalesMap
	General		GeneralConfig
	News		NewsConfig
	Translate	TranslateConfig
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
	News = cfg.News
	Translate = cfg.Translate
	LoadLocales()
}
