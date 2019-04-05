package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// YoutubeCommand youtube handler
func YoutubeCommand(ctx bot.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_no_args"))
		return
	}
	switch ctx.Args[0] {
	case "play":
		youtubePlay(sess, &ctx)
	case "stop":
		youtubeStop(sess, &ctx)
	case "skip":
		youtubeSkip(sess, &ctx)
	case "add":
		youtubeAdd(sess, &ctx)
	case "list":
		youtubeList(sess, &ctx)
	case "clear":
		youtubeClear(sess, &ctx)
	}
}

func youtubePlay(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "play")
	if sess == nil {
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		nsess, serr := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
		if serr != nil {
			//ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error") + " : " + serr.Error())
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("player error: %v", serr.Error()))
			return
		}
		sess = nsess
	}
	queue := sess.Queue
	if !queue.HasNext() {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
		return
	}
	msg := ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_starting"))
	shortPlay(ctx, sess, msg)
}

func youtubeStop(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "stop")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
		return
	}
	if sess.Queue.HasNext() {
		sess.Queue.Clear()
	}
	sess.Stop()
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_stopped"))
}

func youtubeSkip(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "skip")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
		return
	}
	sess.Stop()
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_skipped"))
}

func youtubeAdd(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "add")
	newargs := ctx.Args[1:]
	if len(newargs) == 0 {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_args_missing"))
		return
	}
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
		return
	}
	msg := ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_adding_song"))
	for _, arg := range newargs {
		t, inp, err := ctx.Youtube.Get(arg)
		if err != nil {
			//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting input: %v", err.Error()))
			return
		}

		switch t {
		case bot.ERROR_TYPE:
			//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error type: %v", t))
			return
		case bot.VIDEO_TYPE:
			{
				video, err := ctx.Youtube.Video(*inp)
				if err != nil {
					//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
					ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video1 (VIDEO_TYPE): %v", err.Error()))
					return
				}
				song := bot.NewSong(video.Media, video.Title, arg)
				sess.Queue.Add(song)
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_added_format"), song.Title), true)
				break
			}
		case bot.PLAYLIST_TYPE:
			{
				videos, err := ctx.Youtube.Playlist(*inp)
				if err != nil {
					//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
					ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting playlist: %v", err.Error()))
					return
				}
				for _, v := range *videos {
					id := v.Id
					_, i, err := ctx.Youtube.Get(id)
					if err != nil {
						//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video2: %v", err.Error()))
						continue
					}
					video, err := ctx.Youtube.Video(*i)
					if err != nil {
						//ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video3: %v", err.Error()))
						return
					}
					song := bot.NewSong(video.Media, video.Title, arg)
					sess.Queue.Add(song)
				}
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_added"), true)
				break
			}
		}
	}
}

func youtubeList(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "list")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
		return
	}
	if !sess.Queue.HasNext() {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
		return
	}
	var songsNames []string
	var count = 0
	var countMore = 0
	for _, val := range sess.Queue.Get() {
		count++
		if count <= 10 {
			songsNames = append(songsNames, fmt.Sprintf("[%v]: %v", count, val.Title))
		} else {
			countMore++
		}
	}
	if countMore > 0 {
		songsNames = append(songsNames, fmt.Sprintf(ctx.Loc("youtube_list_more_format"), countMore))
	}
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_list_format"), strings.Join(songsNames, "\n")))
}

func youtubeClear(sess *bot.Session, ctx *bot.Context) {
	ctx.MetricsCommand("youtube_command", "clear")
	if sess == nil {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
		return
	}
	if !sess.Queue.HasNext() {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
		return
	}
	sess.Queue.Clear()
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_cleared"))
}

func shortPlay(ctx *bot.Context, sess *bot.Session, msg *discordgo.Message) (isPlaying bool) {
	queue := sess.Queue
	if !queue.HasNext() {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
		return
	}
	go queue.Start(sess, func(relp string) {
		switch relp {
		case "stop":
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_stopped"), true)
		case "finish":
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_finished"), true)
		default:
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf("%v: %v", ctx.Loc("youtube_now_playing"), relp), true)
			isPlaying = true
		}
	})
	return false
}

// YoutubeShortCommand handle short command for playing song from youtube
func YoutubeShortCommand(ctx bot.Context) {
	ctx.MetricsCommand("youtube_command", "short")
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	newargs := ctx.Args
	if len(newargs) == 0 {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_args_missing"))
		return
	}
	if sess == nil {
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
			return
		}
		nsess, serr := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		}, ctx.Guilds.Guilds[ctx.Guild.ID].VoiceVolume)
		sess = nsess
		if serr != nil {
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("session error: %v", serr.Error()))
			//ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
			return
		}
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID))
	}
	msg := ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_adding_song"))
	for _, arg := range newargs {
		t, inp, err := ctx.Youtube.Get(arg)

		if err != nil {
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting input: %v", err.Error()))
			return
		}

		switch t {
		case bot.ERROR_TYPE:
			ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error type: %v", err.Error()))
			fmt.Println("error type", t)
			return
		case bot.VIDEO_TYPE:
			{
				video, err := ctx.Youtube.Video(*inp)
				if err != nil {
					ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video1: %v", err.Error()))
					return
				}
				song := bot.NewSong(video.Media, video.Title, arg)
				sess.Queue.Add(song)
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_added_format"), song.Title), true)
				shortPlay(&ctx, sess, msg)
				break
			}
		case bot.PLAYLIST_TYPE:
			{
				videos, err := ctx.Youtube.Playlist(*inp)
				if err != nil {
					ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting playlist: %v", err.Error()))
					return
				}
				var isPlaying bool
				for _, v := range *videos {
					id := v.Id
					_, i, err := ctx.Youtube.Get(id)
					if err != nil {
						ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video2: %v", err.Error()))
						continue
					}
					video, err := ctx.Youtube.Video(*i)
					if err != nil {
						ctx.Log("Youtube", ctx.Guild.ID, fmt.Sprintf("error getting video3: %v", err.Error()))
						return
					}
					song := bot.NewSong(video.Media, video.Title, arg)
					sess.Queue.Add(song)
					ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_added_format"), song.Title), true)
					if !isPlaying {
						shortPlay(&ctx, sess, msg)
						isPlaying = true
					}
				}
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_added"), true)
				break
			}
		}
	}
}
