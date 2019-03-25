package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/FlameInTheDark/dtbot/bot"
)

// DebugCommand special bot commands handler
func DebugCommand(ctx bot.Context) {
	ctx.MetricsCommand("debug", "admin")
	if ctx.GetRoles().ExistsName("bot.admin") {
		if len(ctx.Args) == 0 {
			return
		}
		switch ctx.Args[0] {
		case "roles":
			var roles []string
			for _, val := range ctx.GetRoles().Roles {
				roles = append(roles, val.Name)
			}
			ctx.ReplyEmbedPM("Debug", strings.Join(roles, ", "))
		case "time":
			ctx.ReplyEmbedPM("Debug", time.Now().String())
		case "session":
			sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
			if sess != nil {
				ctx.ReplyEmbed("Debug", sess.ChannelID)
			} else {
				ctx.ReplyEmbed("Debug", "Session is nil")
			}
		case "voice":
			var resp string
			resp += fmt.Sprintf("Voice connections: %v\n", len(ctx.Discord.VoiceConnections))
			for i,c := range ctx.Discord.VoiceConnections {
				resp += i + " | G: " + c.GuildID + " | C: " + c.ChannelID + "\n"
			}
			ctx.ReplyEmbed("Debug", resp)
		case "leavevoice":
			if v,ok := ctx.Discord.VoiceConnections[ctx.Args[1]];ok {
				err := v.Disconnect()
				if err != nil {
					ctx.ReplyEmbed("Debug", "Voice: " + err.Error())
				}
			} else {
				ctx.ReplyEmbed("Debug", "Voice connection not found")
			}
		case "volume":
			ctx.ReplyEmbed("Debug", fmt.Sprintf("Voice volume is %.2f", ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume))
		}
	} else {
		ctx.ReplyEmbedPM("Debug", "Not a Admin")
	}
}
