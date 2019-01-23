package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
	ctx.MetricsCommand("bot")
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
			if ctx.IsAdmin() {
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
		case "guild":
			if len(ctx.Args) < 2 && !ctx.IsAdmin() {
				return
			}
			switch ctx.Args[1] {
			case "leave":
				if len(ctx.Args) < 3 {
					return
				}
				err := ctx.Discord.GuildLeave(ctx.Args[2])
				if err != nil {
					ctx.Log("Guild", ctx.Guild.ID, fmt.Sprintf("error leaving from guild [%v]: %v", ctx.Args[2], err.Error()))
					ctx.ReplyEmbedPM("Guild", fmt.Sprintf("Error leaving from guild [%v]: %v", ctx.Args[2], err.Error()))
					return
				}
				ctx.ReplyEmbedPM("Guild", fmt.Sprintf("Leave from guild: %v", ctx.Args[2]))
			case "list":
				var selected string
				if len(ctx.Args) > 2 && ctx.Args[2] == "id" {
					if len(ctx.Guilds) > 4 {
						selected = ctx.Args[3]
					} else {
						selected = "1"
					}
				} else {
					if len(ctx.Guilds) > 3 {
						fmt.Println(len(ctx.Args))
						selected = ctx.Args[2]
					} else {
						selected = "1"
					}
				}
				guilds := ctx.Discord.State.Guilds
				pages := int(len(guilds)/20) + 1
				if len(ctx.Args) > 3 {
					index := 0
					page, err := strconv.Atoi(selected)
					if err != nil {
						index = 0
					} else {
						index = page
						if index > 1 {
							index = index * 20
						} else {
							index = 0
						}
						if index > len(guilds) {
							index = len(guilds) - 1
							if index < 0 {
								index = 0
							}
						}
					}
					var indexEnd = index + 20
					if indexEnd > len(guilds) {
						indexEnd = len(guilds)
					}
					if len(ctx.Args) > 2 && ctx.Args[2] == "id" {
						ctx.ReplyEmbed("Guilds", guildsListID(guilds[index:indexEnd], page, pages))
					} else {
						ctx.ReplyEmbed("Guilds", guildsListName(guilds[index:indexEnd], page, pages))
					}

				} else {
					var indexEnd = 20
					if indexEnd > len(guilds) {
						indexEnd = len(guilds)
					}
					if len(ctx.Args) > 2 && ctx.Args[2] == "id" {
						ctx.ReplyEmbed("Guilds", guildsListID(guilds[:indexEnd], 1, 1))
					} else {
						ctx.ReplyEmbed("Guilds", guildsListName(guilds[:indexEnd], 1, 1))
					}

				}
			}
		}
	}
}

func guildsListID(guilds []*discordgo.Guild, current, pages int) string {
	var list string
	for _, g := range guilds {
		var gName string
		if len(g.Name) > 20 {
			gName = fmt.Sprintf("%v...", g.Name[:20])
		} else {
			gName = g.Name
		}
		list += fmt.Sprintf("[%v] - %v\n", g.ID, gName)
	}
	list += fmt.Sprintf("Pages: %v | Current: %v", pages, current)
	return list
}

func guildsListName(guilds []*discordgo.Guild, current, pages int) string {
	var list string
	for i, g := range guilds {
		var gName string
		if len(g.Name) > 20 {
			gName = fmt.Sprintf("%v...", g.Name[:20])
		} else {
			gName = g.Name
		}
		list += fmt.Sprintf("[%v] - %v | U: %v\n", i, gName, len(g.Members))
	}
	list += fmt.Sprintf("Pages: %v | Current: %v", pages, current)
	return list
}
