package messages

import (
	"strings"
	
	"../api/weather"
	"../api/news"
	"../api/translate"
	"github.com/bwmarrin/discordgo"
)

// Bot messages reactions
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := strings.Split(m.Content, " ")
	switch args[0] {

	case "!n":
		go news.GetNews(s, m, args[1:])
		return
	case "!w":
		go weather.GetForecast(s, m, args[1:])
		return
	case "!t":
		go translate.GetTranslation(s, m, args[1:])
	}
	
}
