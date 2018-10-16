package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

const (
	// ERROR_TYPE returns if error
	ERROR_TYPE = -1
	// VIDEO_TYPE returns if video
	VIDEO_TYPE = 0
	// PLAYLIST_TYPE returns if playlist
	PLAYLIST_TYPE = 1
)

type (
	videoResponse struct {
		Formats []struct {
			Url string `json:"url"`
		} `json:"formats"`
		Title string `json:"title"`
	}

	// VideoResult contains information about video
	VideoResult struct {
		Media string
		Title string
	}

	// PlaylistVideo contains playlist ID
	PlaylistVideo struct {
		Id string `json:"id"`
	}

	// YTSearchContent contains Youtube search result
	YTSearchContent struct {
		Id           string `json:"id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ChannelTitle string `json:"channel_title"`
		Duration     string `json:"duration"`
	}

	ytApiResponse struct {
		Error   bool              `json:"error"`
		Content []YTSearchContent `json:"content"`
	}

	// Youtube contains pointer to bot configuration struct
	Youtube struct {
		Conf *Config
	}
)

func (youtube Youtube) getType(input string) int {
	if strings.Contains(input, "upload_date") {
		return VIDEO_TYPE
	}
	if strings.Contains(input, "_type") {
		return PLAYLIST_TYPE
	}
	return ERROR_TYPE
}

// Get returns data grabbed from youtube
func (youtube Youtube) Get(link string) (int, *string, error) {
	cmd := exec.Command("./youtube-dl", "--skip-download", "--print-json", "--flat-playlist", link)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ERROR_TYPE, nil, err
	}
	str := out.String()
	return youtube.getType(str), &str, nil
}

// Video returns unmarshaled data from Youtube
func (youtube Youtube) Video(input string) (*VideoResult, error) {
	var resp videoResponse
	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		return nil, err
	}
	return &VideoResult{resp.Formats[0].Url, resp.Title}, nil
}

// Playlist returns Playlist
func (youtube Youtube) Playlist(input string) (*[]PlaylistVideo, error) {
	lines := strings.Split(input, "\n")
	videos := make([]PlaylistVideo, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var video PlaylistVideo
		fmt.Println("line,", line)
		err := json.Unmarshal([]byte(line), &video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return &videos, nil
}

func (youtube Youtube) buildUrl(query string) (*string, error) {
	base := youtube.Conf.General.ServiceURL + "/v1/youtube/search"
	address, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("search", query)
	address.RawQuery = params.Encode()
	str := address.String()
	return &str, nil
}

// Search returns array of search results
func (youtube Youtube) Search(query string) ([]YTSearchContent, error) {
	addr, err := youtube.buildUrl(query)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(*addr)
	if err != nil {
		return nil, err
	}
	var apiResp ytApiResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	return apiResp.Content, nil
}
