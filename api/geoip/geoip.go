package geoip

import (
	"encoding/json"
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"net/http"
)

type GeoIP struct {
	IP        string       `json:"ip"`
	City      GeoIPCity    `json:"city"`
	Region    GeoIPRegion  `json:"region"`
	Country   GeoIPCountry `json:"country"`
	Error     string       `json:"error"`
	Requests  int          `json:"requests"`
	Created   string       `json:"created"`
	TimeStamp int          `json:"timestamp"`
}

type GeoIPCity struct {
	Id         int     `json:"id"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
	NameRU     string  `json:"name_ru"`
	NameEN     string  `json:"name_en"`
	OKATO      string  `json:"okato"`
	VK         int     `json:"vk"`
	Population int     `json:"population"`
	Tel        string  `json:"tel"`
	PostalCode string  `json:"post"`
}

type GeoIPRegion struct {
	Id        int     `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	NameRU    string  `json:"name_ru"`
	NameEN    string  `json:"name_en"`
	OKATO     string  `json:"okato"`
	VK        int     `json:"vk"`
	ISO       string  `json:"iso"`
	TimeZone  string  `json:"timezone"`
	Auto      string  `json:"auto"`
	UTC       int     `json:"utc"`
}

type GeoIPCountry struct {
	ID            int     `json:"id"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lon"`
	NameRU        string  `json:"name_ru"`
	NameEN        string  `json:"name_en"`
	ISO           string  `json:"iso"`
	Continent     string  `json:"continent"`
	TimeZone      string  `json:"timezone"`
	Area          int     `json:"area"`
	Population    int     `json:"population"`
	CapitalID     int     `json:"capital_id"`
	CapitalNameRU string  `json:"capital_ru"`
	CapitalNameEN string  `json:"capital_en"`
	CurrencyCode  string  `json:"cur_code"`
	PhonePrefix   string  `json:"phone"`
	Neighbours    string  `json:"neighbours"`
	VK            int     `json:"vk"`
	UTC           int     `json:"utc"`
}

func GetGeoIP(ctx *bot.Context) string {
	resp, err := http.Get(fmt.Sprintf("http://api.sypexgeo.net/json/%v", ctx.Args[0]))
	if err != nil {
		return ctx.Loc("geoip_no_data")
	}

	var result GeoIP
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ctx.Loc("error")
	}

	if result.City.NameRU == "" || result.Region.NameRU == "" || result.Country.NameRU == "" {
		return ctx.Loc("geoip_no_data")
	}

	return fmt.Sprintf("IP [%v]\nCity: %v\nRegion: %v\nCountry: %v",
		result.IP, result.City.NameRU, result.Region.NameRU, result.Country.NameRU)

}
