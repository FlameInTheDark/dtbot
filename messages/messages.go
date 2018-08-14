package messages

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Bot messages reactions
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := strings.Split(m.Content, " ")
    if args[0] == "!w" {
		if len(args) == 1 {
			go getForecast(s, m)
		} else if len(args) == 2 {
			go getForecast(s, m, args[1])
		} else {
			go getForecast(s, m, args[1], args[2])
		}
	}
}
