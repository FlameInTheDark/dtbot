package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/bot"
)

// RadioCommand Player handler
func RadioCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_no_args"))
		return
	}
	switch ctx.Args[0] {
	case "play":
		playerPlay(sess, &ctx)
	case "list":
		playerList(&ctx)
	case "station":
		playerStation(sess, &ctx)
	case "genres":
		playerCategories(&ctx)
	case "stop":
		ctx.MetricsCommand("radio", "stop")
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
			return
		}
		sess.Stop()
	}
}

func playerPlay(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("radio", "play")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
		return
	}
	if ctx.Arg(1) == "attachment" {
		go sess.Player.Start(sess, ctx.Message.Attachments[0].URL, func(msg string) {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
	} else if len(ctx.Args) > 1 {
		go sess.Player.Start(sess, ctx.Args[1], func(msg string) {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
	}
}

func playerList(ctx *bot.Context) {
	ctx.MetricsCommand("radio", "list")
	var stations []bot.RadioStation
	if len(ctx.Args) > 1 {
		stations = ctx.DB.GetRadioStations(ctx.Args[1])
	} else {
		stations = ctx.DB.GetRadioStations("")
	}

	var category = make(map[string][]*bot.RadioStation)

	for i, s := range stations {
		category[s.Category] = append(category[s.Category], &stations[i])
	}

	if len(stations) > 0 {
		var embed = bot.NewEmbed(ctx.Loc("player"))
		for c, st := range category {
			var response string
			if len(st) > 20 {
				for _, s := range st[:20] {
					response += fmt.Sprintf("[%v] - %v\n", s.Key, s.Name)
				}
			} else {
				for _, s := range st {
					response += fmt.Sprintf("[%v] - %v\n", s.Key, s.Name)
				}
			}
			embed.Field(c, response, false).Color(ctx.GuildConf().EmbedColor)
		}
		embed.Send(ctx)
	} else {
		ctx.ReplyEmbed(ctx.Loc("player"), ctx.Loc("stations_not_found"))
	}
}

func playerCategories(ctx *bot.Context) {
	ctx.MetricsCommand("radio", "categories")
	stations := ctx.DB.GetRadioStations("")
	var categories = make(map[string]bool)
	var reply = ctx.Loc("stations_categories")
	for _, st := range stations {
		categories[st.Category] = true
	}

	for c, _ := range categories {
		reply += "`" + c + "`, "
	}
	ctx.ReplyEmbed(ctx.Loc("player"), reply)
}

func playerStation(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("radio", "station")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
		return
	}
	if len(ctx.Args) > 1 {
		station, err := ctx.DB.GetRadioStationByKey(ctx.Args[1])
		if err != nil {
			ctx.ReplyEmbed(ctx.Loc("player"), ctx.Loc("stations_not_found"))
			return
		}
		go sess.Player.Start(sess, station.URL, func(msg string) {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
	}
}
