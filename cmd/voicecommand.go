package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/globalsign/mgo/bson"
	"log"
	"strconv"
)

// VoiceCommand voice handler
func VoiceCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) < 1 {
		return
	}
	switch ctx.Args[0] {
	case "join":
		voiceJoin(sess, &ctx)
	case "leave":
		voiceLeave(sess, &ctx)
	case "volume":
		voiceVolume(sess, &ctx)
	}
}

func voiceJoin(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("voice", "join")
	log.Println("used")
	if ctx.Sessions.GetByGuild(ctx.Guild.ID) != nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_connected"))
		return
	}
	log.Println("guild founded")
	vc := ctx.GetVoiceChannel()
	if vc == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
		return
	}
	log.Println("get voice channel")
	s := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if s != nil {
		if !s.IsOk() {
			log.Println("removed failed session")
			ctx.Sessions.Leave(ctx.Discord, *sess)
		}
	}

	sess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
		Muted:    false,
		Deafened: true,
	}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
	if err != nil {
		log.Println("Voice join error: ", err)
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
		if sess != nil {
			ctx.Sessions.Leave(ctx.Discord, *sess)
		}
		return
	}
	log.Println("join session")
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID))
}

func voiceLeave(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("voice", "leave")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
		return
	}
	ctx.Sessions.Leave(ctx.Discord, *sess)
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_left"), sess.ChannelID))
}

func voiceVolume(sess *bot.Session, ctx *bot.Context) {
	if len(ctx.Args) > 1 {
		vol, err := strconv.ParseFloat(ctx.Args[1], 32)
		if err != nil {
			ctx.ReplyEmbed(ctx.Loc("player"), ctx.Loc("player_wrong_volume"))
			return
		}
		ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume = float32(vol * 0.01)
		_ = ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"voicevolume": float32(vol * 0.01)}})
		ctx.ReplyEmbed(ctx.Loc("player"), fmt.Sprintf(ctx.Loc("player_volume_changed"), ctx.Args[1]))
		sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
		if sess != nil {
			sess.Volume = float32(vol * 0.01)
		}
	}
}
