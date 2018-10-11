package translate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../../bot"
)

// TranslateResponse : Translate API struct
type TranslateResponse struct {
	Code     int      `json:"code"`
	Language string   `json:"lang"`
	Text     []string `json:"text"`
	Message  string   `json:"message"`
}

// GetTranslation returns translated text
func GetTranslation(ctx *bot.Context) string {
	var (
		result    TranslateResponse
		translate string
	)

	if len(ctx.Args) > 1 {
		translate = strings.Join(ctx.Args[1:], "+")
	} else {
		return ctx.Loc("translate_request_error")
	}

	resp, err := http.Get(fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?key=%v&text=%v&lang=%v&format=plain", ctx.Conf.Translate.ApiKey, translate, ctx.Args[0]))
	if err != nil {
		return fmt.Sprintf("%v: %v", ctx.Loc("translate_get_error"), err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Sprintf("%v: %v", ctx.Loc("translate_parse_error"), err)
	}

	// Checking request status
	switch result.Code {
	case 502:
		return ctx.Loc("translate_request_error")
	case 200:
		return strings.Join(result.Text, "\n")
	default:
		return ctx.Loc("translate_api_error")
	}
}
