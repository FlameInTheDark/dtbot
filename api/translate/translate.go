package translate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../../config"
	"github.com/bwmarrin/discordgo"
)

type TranslateResponse struct {
	Code		int			`json:"code"`
	Language	string		`json:"lang"`
	Text		[]string	`json:"text"`
	Message		string		`json:"message"`
}

func GetTranslation(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var (
		result TranslateResponse
		translate string = ""
	)
	
	if len(args) > 1 {
		translate = strings.Join(args[1:], "+")
	} else {
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("translate_request_error"))
		return
	}
	
	resp, err := http.Get(fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?key=%v&text=%v&lang=%v&format=plain", config.Translate.ApiKey, translate, args[0]))
	if err != nil {
		fmt.Printf("Get translation error: %v", err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Parse translation error: %v", err)
		return
	}
	
	switch result.Code {
	case 502:
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("translate_request_error"))
		return
	case 200:
		s.ChannelMessageSend(m.ChannelID, strings.Join(result.Text, "\n"))
	default:
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("translate_api_error"))
	}
}