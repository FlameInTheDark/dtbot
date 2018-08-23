package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../../bot"
)

// NewsResponseData : News main struct
type NewsResponseData struct {
	Status       string            `json:"status"`
	TotalResults int               `json:"totalResults"`
	Articles     []NewsArticleData `json:"articles"`
}

// NewsArticleData : News article struct
type NewsArticleData struct {
	Source      NewsArticeleSourceData `json:"source"`
	Author      string                 `json:"author"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Url         string                 `json:"url"`
	PublishedAt string                 `json:"publishedAt"`
}

// Article source struct
type NewsArticeleSourceData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetNews returns news string
func GetNews(ctx *bot.Context) string {
	var (
		result   NewsResponseData
		category string = ""
	)
	if len(ctx.Args) > 0 {
		category = ctx.Args[0]
	}
	resp, err := http.Get(fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%v&category=%v&apiKey=%v", ctx.Conf.News.Country, category, ctx.Conf.News.ApiKey))
	if err != nil {
		return fmt.Sprintf("Get news error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Sprintf("Parse news error: %v", err)
	}

	if result.Status == "ok" {
		if len(result.Articles) > 0 {
			var news []string
			for i := 0; i < ctx.Conf.News.Articles; i++ {
				news = append(news, fmt.Sprintf("```%v\n\n%v\nLink: %v```", result.Articles[i].Title, result.Articles[i].Description, result.Articles[i].Url))
			}

			return strings.Join(news, "\n")
		} else {
			return ctx.Conf.GetLocale("news_404")
		}
	} else {
		return ctx.Conf.GetLocale("news_api_error")
	}
}
