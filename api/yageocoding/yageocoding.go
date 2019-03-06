package yageocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// YaGeoResponse contains response data
type YaGeoResponse struct {
	Response struct {
		ObjectCollection YaGeoObjectCollection `json:"GeoObjectCollection"`
	} `json:"response"`
}

// YaGeoObjectCollection contains geo objects
type YaGeoObjectCollection struct {
	MetaData YaGeoMetaData `json:"metaDataProperty"`
	Member   []YaGeoMember `json:"featureMember"`
}

// YaGeoMetaData contains metadata
type YaGeoMetaData struct {
	ResponseMetaData struct {
		Request string `json:"request"`
		Found   string `json:"found"`
		Results string `json:"results"`
	} `json:"GeocoderResponseMetaData"`
}

// YaGeoMember contains object data
type YaGeoMember struct {
	GeoObject struct {
		MetaData    YaGeoMemberMetaData `json:"metaDataProperty"`
		Description string              `json:"description"`
		Name        string              `json:"name"`
		Point       struct {
			Pos string `json:"pos"`
		} `json:"Point"`
	} `json:"GeoObject"`
}

// YaGeoMemberMetaData contains member metadata
type YaGeoMemberMetaData struct {
	Meta struct {
		Kind      string `json:"kind"`
		Text      string `json:"text"`
		Precision string `json:"precision"`
	} `json:"GeocoderMetaData"`
}

// YaGeoAddress contains address structure
type YaGeoAddress struct {
	CountryCode string                  `json:"country_code"`
	PostalCode  string                  `json:"postal_code"`
	Formatted   string                  `json:"formatted"`
	Components  []YaGeoAddressComponent `json:"Components"`
}

// YaGeoAddressComponent is address component
type YaGeoAddressComponent struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// GetCoordinates returns coordinates
func (loc *YaGeoResponse) GetCoordinates() (string, string) {
	if len(loc.Response.ObjectCollection.Member) == 0 {
		return "0", "0"
	}

	str := strings.Split(loc.Response.ObjectCollection.Member[0].GeoObject.Point.Pos, " ")
	return str[0], str[1]
}

// GetData creates request to API and returns result
func GetData(key, location string) (result YaGeoResponse, err error) {
	resp, err := http.Get(fmt.Sprintf("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=%v&apikey=%v", location, key))
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	if len(result.Response.ObjectCollection.Member) == 0 {
		return result, errors.New("location not fount")
	}
	return
}
