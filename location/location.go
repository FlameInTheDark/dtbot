package location

import (
	"encoding/json"
    "net/http"
    "errors"
    "fmt"
    "../config"
)

type LocationResultData struct {
    Status      string          `json:"status"`
    Locations   []LocationData  `json:"results"`
}

type LocationData struct {
    AddrComponents  []AddrComponentData `json:"address_components"`
    FormatedAddress string              `json:"formatted_address"`
    Geometry        GeometryData        `json:"geometry"`
}

type CoordinateData struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

type GeometryData struct {
    Location CoordinateData `json:"location"`
}

type AddrComponentData struct {
    LongName    string      `json:"long_name"`
    ShortName   string      `json:"short_name"`
    Types       []string    `json:"types"`
}

func (l LocationResultData) GetCoordinates() (float64, float64) {
    return l.Locations[0].Geometry.Location.Lat, l.Locations[0].Geometry.Location.Lng
}

func New(locationName string) (LocationResultData, error) {
    var result LocationResultData
    resp, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%v&sensor=true&key=%v", locationName, config.General.GeocodingKey))
	if err != nil {
		return result, err
	}
    
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
    }
    
    if result.Status == "OK" {
        return result, nil
    } else {
        return result, errors.New("City not found!")
    }
}