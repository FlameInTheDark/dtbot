package cmd

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"
	"strings"

	"github.com/FlameInTheDark/dtbot/bot"
)

// BotCommand special bot commands handler
func BotCommand(ctx bot.Context) {
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
		if len(ctx.Args)<2 {
			logs := ctx.DB.LogGet(10)
			var logString = []string{""}
			for _,log := range logs {
				logString = append(logString,fmt.Sprintf("[%v] %v: %v\n",log.Date,log.Module,log.Text))
			}
			ctx.ReplyEmbedPM("Logs",strings.Join(logString,""))
		} else {
			count, err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			logs := ctx.DB.LogGet(count)
			var logString = []string{""}
			for _,log := range logs {
				logString = append(logString,fmt.Sprintf("[%v] %v: %v\n",log.Date,log.Module,log.Text))
			}
			ctx.ReplyEmbedPM("Logs",strings.Join(logString,""))
		}
	case "setconf":
		if len(ctx.Args) > 2{
			switch ctx.Args[1] {
			case "language":
				ctx.Guilds[ctx.Guild.ID].Language = ctx.Args[2]
				ctx.DB.DBSession.DB(ctx.DB.DBName).C("guilds").Update(bson.M{"id":ctx.Guild.ID},bson.M{"$set":bson.M{"language":ctx.Args[2]}})
				ctx.ReplyEmbedPM("Config",fmt.Sprintf("Language set to: %v", ctx.Args[2]))
			case "timezone":
				tz,err := strconv.Atoi(ctx.Args[1])
				if err != nil {
					ctx.ReplyEmbedPM("Settings","Not a number")
				}
				ctx.Guilds[ctx.Guild.ID].Timezone = tz
				ctx.DB.DBSession.DB(ctx.DB.DBName).C("guilds").Update(bson.M{"id":ctx.Guild.ID},bson.M{"$set":bson.M{"timezone":tz}})
				ctx.ReplyEmbedPM("Config",fmt.Sprintf("Timezone set to: %v", ctx.Args[2]))
			case "weather.city":
				ctx.Guilds[ctx.Guild.ID].WeatherCity = ctx.Args[2]
				ctx.DB.DBSession.DB(ctx.DB.DBName).C("guilds").Update(bson.M{"id":ctx.Guild.ID},bson.M{"$set":bson.M{"weathercity":ctx.Args[2]}})
				ctx.ReplyEmbedPM("Config",fmt.Sprintf("Weather city set to: %v", ctx.Args[2]))
			case "news.country":
				ctx.Guilds[ctx.Guild.ID].NewsCounty = ctx.Args[2]
				ctx.DB.DBSession.DB(ctx.DB.DBName).C("guilds").Update(bson.M{"id":ctx.Guild.ID},bson.M{"$set":bson.M{"weathercountry":ctx.Args[2]}})
				ctx.ReplyEmbedPM("Config",fmt.Sprintf("News country set to: %v", ctx.Args[2]))
			}
		}

	}
}
