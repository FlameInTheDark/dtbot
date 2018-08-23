package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Weather config struct
type WeatherConfig struct {
	WeatherToken string
	City         string
}

// General config struct
type GeneralConfig struct {
	Language         string
	Timezone         int
	GeonamesUsername string
}

// News config struct
type NewsConfig struct {
	ApiKey   string
	Country  string
	Articles int
}

// Yandex translate config struct
type TranslateConfig struct {
	ApiKey string
}

// Map with locales
type LocalesMap map[string]map[string]string

// Main config struct
type Config struct {
	Weather   WeatherConfig
	General   GeneralConfig
	News      NewsConfig
	Translate TranslateConfig
	Locales   LocalesMap
}

func (c Config) GetLocale(key string) string {
	return c.Locales[c.General.Language][key]
}

// Loading configs from file
func LoadConfig() *Config {
	var cfg Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		fmt.Printf("Config loading error: %v\n", err)
		os.Exit(1)
	}
	cfg.LoadLocales()
	return &cfg
}

// Loading locales from file
func (c Config) LoadLocales() {
	file, e := ioutil.ReadFile("./locales.json")
	if e != nil {
		fmt.Printf("Locale file error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &c.Locales)
	if err != nil {
		panic(err)
	}

	if _, ok := c.Locales[c.General.Language]; ok {
		return
	} else {
		fmt.Printf("Locale file not contain language \"%v\"\n", c.General.Language)
		os.Exit(1)
	}
}
