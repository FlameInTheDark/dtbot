package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// WeatherConfig Weather config struct
type WeatherConfig struct {
	WeatherToken string
	City         string
}

// GeneralConfig General config struct
type GeneralConfig struct {
	Language         string
	Timezone         int
	GeonamesUsername string
	Game             string
	EmbedColor       int
	ServiceURL       string
	MessagePool      int
	DatabaseName     string
	GeocodingApiKey  string
	AdminID          string
}

// NewsConfig News config struct
type NewsConfig struct {
	APIKey   string
	Country  string
	Articles int
}

// MetricsConfig InfluxDB connection settings
type MetricsConfig struct {
	Address  string
	Database string
	User     string
	Password string
}

type DBLConfig struct {
	Token string
}

// TranslateConfig Yandex translate config struct
type TranslateConfig struct {
	APIKey string
}

// CurrencyConfig Currency config struct
type CurrencyConfig struct {
	Default []string
}

// LocalesMap Map with locales
type LocalesMap map[string]map[string]string

// WeatherCodesMap symbols for font
type WeatherCodesMap map[string]string

// Config Main config struct. Contains all another config structs data.
type Config struct {
	Weather      WeatherConfig
	General      GeneralConfig
	News         NewsConfig
	Translate    TranslateConfig
	Locales      LocalesMap
	Currency     CurrencyConfig
	WeatherCodes WeatherCodesMap
	Metrics      MetricsConfig
	DBL          DBLConfig
}

// GetLocale returns locale string by key
func (c *Config) GetLocale(key string) string {
	return c.Locales[c.General.Language][key]
}

// LoadConfig loads configs from file 'config.toml'. Terminate program if error.
func LoadConfig() *Config {
	var cfg Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		fmt.Printf("Config loading error: %v\n", err)
		os.Exit(1)
	}
	cfg.LoadLocales()
	cfg.LoadWeatherCodes()
	return &cfg
}

// LoadLocales loads locales from file 'locales.json'. Terminate program if error.
func (c *Config) LoadLocales() {
	file, e := ioutil.ReadFile("./locales.json")
	if e != nil {
		fmt.Printf("Locale file error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &c.Locales)
	if err != nil {
		panic(err)
	}

	if _, ok := c.Locales[c.General.Language]; !ok {
		fmt.Printf("Locale file not contain language \"%v\"\n", c.General.Language)
		os.Exit(1)
	}

	fmt.Printf("Loaded %v translations for '%v' language\n", len(c.Locales[c.General.Language]), c.General.Language)
	return
}

// LoadWeatherCodes loads weather font codes from file 'codes.json' in map. Terminate program if error.
func (c *Config) LoadWeatherCodes() {
	file, e := ioutil.ReadFile("./codes.json")
	if e != nil {
		fmt.Printf("Codes file error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &c.WeatherCodes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %v weather codes\n", len(c.WeatherCodes))
	return
}
