package fun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ImageResponse contains image API response data
type ImageResponse struct {
	Error    string `json:"error"`
	Success  bool   `json:"success"`
	ImageURL string `json:"image"`
}

// GetImageURL returns image url
func GetImageURL(category string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://botimages.realpha.org/?category=%v", category))
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Getting image url error: %v", err)
		return "", errors.New("getting image url error")
	}

	var result ImageResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}

	if result.Success {
		return fmt.Sprintf("https://botimages.realpha.ru/%v", result.ImageURL), nil
	}
	return "", errors.New("wrong data")
}

// GetImage return image bytes buffer
func GetImage(category string) (*bytes.Buffer, error) {
	resp, err := http.Get(fmt.Sprintf("https://botimages.realpha.ru/?category=%v", category))
	if err != nil {
		fmt.Printf("Getting image url error: %v", err)
		return nil, errors.New("getting image url error")
	}

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, resp.Body) //png.Encode(buf, resp.Body)
	if err != nil {
		fmt.Printf("Map image: %v", err.Error())
	}
	return buf, err
}
