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
		go getForecast(s, m, args[1:])
	}
}
