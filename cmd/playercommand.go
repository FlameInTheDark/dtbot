package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/bot"
)

// PlayerCommand Player handler
func PlayerCommand(ctx bot.Context) {
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
	case "stop":
		ctx.MetricsCommand("player", "stop")
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
			return
		}
		sess.Stop()
	}
}

func playerPlay(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("player", "play")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_not_in_voice"))
		return
	}
	if len(ctx.Args) > 1 {
		go sess.Player.Start(sess, ctx.Args[1], func(msg string) {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), msg)
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)

	}
}

func playerList(ctx *bot.Context) {
	stations := ctx.DB.GetRadioStations()
	if len(stations) > 0 {
		var response string
		if len(stations) > 20 {
			for _, s := range stations[:20] {
				response += fmt.Sprintf("[%v] - %v\n", s.Key, s.Name)
			}
		} else {
			for _, s := range stations {
				response += fmt.Sprintf("[%v] - %v\n", s.Key, s.Name)
			}
		}
		ctx.ReplyEmbed(ctx.Loc("player"), response)
	} else {
		ctx.ReplyEmbed(ctx.Loc("player"), ctx.Loc("stations_not_found"))
	}
}

func playerStation(sess *bot.Session, ctx *bot.Context) {
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
