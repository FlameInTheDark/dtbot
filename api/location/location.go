package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"../../config"
)

type LocationResultData struct {
	ResultCount		int			`json:"totalResultsCount"`
	Geonames		[]Geoname	`json:"geonames"`
}

type Geoname struct {
	GeonameID		int			`json:"geonameId"`
	CountryID		string		`json:"countryId"`
	ToponymName		string		`json:"toponymName"`
	Population		int			`json:"population"`
	CountryCode		string		`json:"countryCode"`
	Name			string		`json:"name"`
	CountryName		string		`json:"countryName"`
	Lat				string		`json:"lat"`
	Lng				string		`json:"lng"`
}

func (l LocationResultData) GetCoordinates() (string, string) {
	return l.Geonames[0].Lat, l.Geonames[0].Lng
}

func New(locationName string) (LocationResultData, error) {
	var result LocationResultData
	resp, err := http.Get(fmt.Sprintf("http://api.geonames.org/searchJSON?q=%v&maxRows=1&username=%v", locationName, config.General.GeonamesUsername))
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(result.Geonames) > 0 {
		return result, nil
	} else {
		return result, errors.New("City not found!")
	}
}
