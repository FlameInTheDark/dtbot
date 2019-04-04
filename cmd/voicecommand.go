package cmd

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"

	"github.com/FlameInTheDark/dtbot/bot"
)

// VoiceCommand voice handler
func VoiceCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) < 1 {
		return
	}
	switch ctx.Args[0] {
	case "join":
		ctx.MetricsCommand("voice", "join")
		if sess != nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_connected"))
			return
		}
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		sess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
		if err != nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
			return
		}
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID))
	case "leave":
		ctx.MetricsCommand("voice", "leave")
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		ctx.Sessions.Leave(ctx.Discord, *sess)
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_left"), sess.ChannelID))
	case "volume":
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
}
