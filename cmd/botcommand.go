package cmd

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"
	"strings"

	"github.com/FlameInTheDark/dtbot/bot"
)

func showLogs(ctx *bot.Context, count int) {
	logs := ctx.DB.LogGet(count)
	var logString = []string{""}
	for _, log := range logs {
		logString = append(logString, fmt.Sprintf("[%v] %v: %v\n", log.Date, log.Module, log.Text))
	}
	ctx.ReplyEmbedPM("Logs", strings.Join(logString, ""))
}

// TODO: I should make it more tasty and remove all "switch/case"!
// BotCommand special bot commands handler
func BotCommand(ctx bot.Context) {
	if ctx.GetRoles().ExistsName("bot.admin") {
		if len(ctx.Args) == 0 {
			return
		}
		switch ctx.Args[0] {
		case "clear":
			if len(ctx.Args) < 2 {
				ctx.BotMsg.Clear(&ctx, 0)
				return
			}
			from, err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				return
			}
			ctx.BotMsg.Clear(&ctx, from)
		case "help":
			ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_reply"))
		case "logs":
			if len(ctx.Args) < 2 {
				showLogs(&ctx, 10)
			} else {
				count, err := strconv.Atoi(ctx.Args[1])
				if err != nil {
					fmt.Println(err)
					return
				}
				showLogs(&ctx, count)
			}
		case "setconf":
			if len(ctx.Args) > 2 {
				target := strings.Split(ctx.Args[1], ".")
				switch target[0] {
				case "general":
					switch target[1] {
					case "language":
						ctx.Guilds[ctx.Guild.ID].Language = ctx.Args[2]
						ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"language": ctx.Args[2]}})
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("Language set to: %v", ctx.Args[2]))
					case "timezone":
						tz, err := strconv.Atoi(ctx.Args[1])
						if err != nil {
							ctx.ReplyEmbedPM("Settings", "Not a number")
						}
						ctx.Guilds[ctx.Guild.ID].Timezone = tz
						ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"timezone": tz}})
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("Timezone set to: %v", ctx.Args[2]))
					case "nick":
						ctx.Discord.GuildMemberNickname(ctx.Guild.ID, "@me", ctx.Args[2])
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("Nickname changed to %v", ctx.Args[2]))
					}
				case "weather":
					switch target[1] {
					case "city":
						ctx.Guilds[ctx.Guild.ID].WeatherCity = ctx.Args[2]
						ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"weathercity": ctx.Args[2]}})
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("Weather city set to: %v", ctx.Args[2]))
					}
				case "news":
					switch target[1] {
					case "country":
						ctx.Guilds[ctx.Guild.ID].NewsCounty = ctx.Args[2]
						ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"weathercountry": ctx.Args[2]}})
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("News country set to: %v", ctx.Args[2]))
					}
				case "embed":
					switch target[1] {
					case "color":
						var color int64
						var err error
						if strings.HasPrefix(ctx.Args[2], "#") {
							color, err = strconv.ParseInt(ctx.Args[2][1:], 16, 32)
							if err != nil {
								ctx.Log("Config", ctx.Guild.ID, fmt.Sprintf("error setting parameter %v to value %v: %v", ctx.Args[1], target[2], err.Error()))
								return
							}
						} else {
							color, err = strconv.ParseInt(ctx.Args[2], 16, 32)
							if err != nil {
								ctx.Log("Config", ctx.Guild.ID, fmt.Sprintf("error setting parameter %v to value %v: %v", ctx.Args[1], ctx.Args[2], err.Error()))
								return
							}
						}
						ctx.Guilds[ctx.Guild.ID].EmbedColor = int(color)
						ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"embedcolor": int(color)}})
						ctx.ReplyEmbedPM("Config", fmt.Sprintf("Embed color set to: %v", ctx.Args[2]))
					}
				}
			}

		}
	}
	ctx.MetricsCommand("currency")
}
