package news

import (
	"encoding/json"
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/pkg/errors"
	"net/http"
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
	URL         string                 `json:"url"`
	PublishedAt string                 `json:"publishedAt"`
}

// NewsArticeleSourceData : Article source struct
type NewsArticeleSourceData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetNews returns news string
func GetNews(ctx *bot.Context) error{
	var (
		result   NewsResponseData
		category string
	)
	if len(ctx.Args) > 0 {
		category = ctx.Args[0]
	}
	resp, err := http.Get(fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%v&category=%v&apiKey=%v", ctx.Conf.News.Country, category, ctx.Conf.News.APIKey))
	if err != nil {
		ctx.Log("news", ctx.Guild.ID, fmt.Sprintf("Get news error: %v", err))
		return errors.New(fmt.Sprintf("Get news error: %v", err))
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		ctx.Log("news", ctx.Guild.ID, fmt.Sprintf("Parse news error: %v", err))
		return errors.New(fmt.Sprintf("Parse news error: %v", err))
	}

	if result.Status == "ok" {
		if len(result.Articles) > 0 {
			emb := bot.NewEmbed(ctx.Loc("news"))
			for i := 0; i < ctx.Conf.News.Articles; i++ {
				emb.Field(result.Articles[i].Title, result.Articles[i].Description + "\n" + result.Articles[i].URL, false)
			}
			emb.Desc(fmt.Sprintf("%v %v",ctx.Loc("requested_by"), ctx.Message.Author.Username))
			_,_=ctx.Discord.ChannelMessageSendEmbed(ctx.Message.ChannelID, emb.GetEmbed())
			return nil
		} else {
			return errors.New(ctx.Loc("news_404"))
		}
	} else {
		return errors.New(ctx.Loc("news_api_error"))
	}
}
