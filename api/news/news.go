package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../../config"
	"github.com/bwmarrin/discordgo"
)

type NewsResponseData struct {
	Status			string				`json:"status"`
	TotalResults	int					`json:"totalResults"`
	Articles		[]NewsArticleData	`json:"articles"`
}

type NewsArticleData struct {
	Source		NewsArticeleSourceData	`json:"source"`
	Author		string					`json:"author"`
	Title		string					`json:"title"`
	Description	string					`json:"description"`
	Url			string					`json:"url"`
	PublishedAt	string					`json:"publishedAt"`
}

type NewsArticeleSourceData	struct {
	Id		string	`json:"id"`
	Name	string	`json:"name"`
}

func GetNews(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var (
		result NewsResponseData
		category string = ""
	)
	if len(args) > 0 {
		category = args[0]
	}
	resp, err := http.Get(fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%v&category=%v&apiKey=%v", config.News.Country, category, config.News.ApiKey))
	if err != nil {
		fmt.Printf("Get news error: %v", err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Parse news error: %v", err)
		return
	}

	if result.Status == "ok" {
		if len(result.Articles) > 0 {
			var news []string
			for i := 0; i < config.News.Articles; i++ {
				news = append(news, fmt.Sprintf("```%v\n\n%v\nLink: %v```", result.Articles[i].Title, result.Articles[i].Description, result.Articles[i].Url))
			}
			
			s.ChannelMessageSend(m.ChannelID, strings.Join(news, "\n"))
			return
		} else {
			s.ChannelMessageSend(m.ChannelID, config.Locales.Get("news_404"))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, config.Locales.Get("news_api_error"))
	}
}