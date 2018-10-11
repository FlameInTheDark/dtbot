package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// LocationResultData : location struct
type LocationResultData struct {
	ResultCount int       `json:"totalResultsCount"`
	Geonames    []Geoname `json:"geonames"`
}

// Geoname : location geonames struct
type Geoname struct {
	GeonameID   int    `json:"geonameId"`
	CountryID   string `json:"countryId"`
	ToponymName string `json:"toponymName"`
	Population  int    `json:"population"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	CountryName string `json:"countryName"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
}

// GetCoordinates returns latitude and longitude
func (l LocationResultData) GetCoordinates() (string, string) {
	return l.Geonames[0].Lat, l.Geonames[0].Lng
}

// New creates and returns location struct
func New(user string, locationName string) (result LocationResultData, err error) {
	resp, err := http.Get(fmt.Sprintf("http://api.geonames.org/searchJSON?q=%v&maxRows=1&username=%v", locationName, user))
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(result.Geonames) == 0 {
		return result, errors.New("City not found")
	}

	return result, nil
}
