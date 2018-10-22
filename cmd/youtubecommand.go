package cmd

import (
	"fmt"
	"strings"

	"../bot"
	"github.com/bwmarrin/discordgo"
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
		if sess == nil {
			vc := ctx.GetVoiceChannel()
			if vc == nil {
				ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_must_be_in_voice"))
				return
			}
			nsess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
				Muted:    false,
				Deafened: true,
			})
			if err != nil {
				ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
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
		shortPlay(&ctx, sess, msg)
	case "stop":
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
			return
		}
		if sess.Queue.HasNext() {
			sess.Queue.Clear()
		}
		sess.Stop()
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_stopped"))
	case "skip":
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
			return
		}
		sess.Stop()
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_skipped"))
	case "add":
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
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
				fmt.Println("error getting input,", err)
				fmt.Println(arg)
				return
			}

			switch t {
			case bot.ERROR_TYPE:
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
				fmt.Println("error type", t)
				return
			case bot.VIDEO_TYPE:
				{
					video, err := ctx.Youtube.Video(*inp)
					if err != nil {
						ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						fmt.Println("error getting video1,", err)
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
						ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						fmt.Println("error getting playlist,", err)
						return
					}
					for _, v := range *videos {
						id := v.Id
						_, i, err := ctx.Youtube.Get(id)
						if err != nil {
							ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
							fmt.Println("error getting video2,", err)
							continue
						}
						video, err := ctx.Youtube.Video(*i)
						if err != nil {
							ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
							fmt.Println("error getting video3,", err)
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
	case "list":
		if sess == nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("player_not_in_voice"))
			return
		}
		if !sess.Queue.HasNext() {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
			return
		}
		var songsNames []string
		for _, val := range sess.Queue.Get() {
			songsNames = append(songsNames, val.Title)
		}
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_list_format"), strings.Join(songsNames, "\n")))
	case "clear":

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
}

func shortPlay(ctx *bot.Context, sess *bot.Session, msg *discordgo.Message) (isPlaying bool) {
	queue := sess.Queue
	if !queue.HasNext() {
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_queue_is_empty"))
		return
	}
	//cmsg := ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_starting"))
	go queue.Start(sess, func(relp string) {
		switch relp {
		case "stop":
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_stopped"), true)
			break
		case "finish":
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_finished"), true)
			break
		default:
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf("%v: %v", ctx.Loc("youtube_now_playing"), relp), true)
			isPlaying = true
		}
	})
	return false
}

func YoutubeShortCommand(ctx bot.Context) {
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
		nsess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, bot.JoinProperties{
			Muted:    false,
			Deafened: true,
		})
		sess = nsess
		if err != nil {
			ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), ctx.Loc("player_error"))
			return
		}
		ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("player")), fmt.Sprintf("%v <#%v>!", ctx.Loc("player_joined"), sess.ChannelID))
	}
	msg := ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_adding_song"))
	for _, arg := range newargs {
		t, inp, err := ctx.Youtube.Get(arg)

		if err != nil {
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
			fmt.Println("error getting input,", err)
			return
		}

		switch t {
		case bot.ERROR_TYPE:
			ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
			fmt.Println("error type", t)
			return
		case bot.VIDEO_TYPE:
			{
				video, err := ctx.Youtube.Video(*inp)
				if err != nil {
					ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
					fmt.Println("error getting video1,", err)
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
					ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
					fmt.Println("error getting playlist,", err)
					return
				}
				var is_playing bool
				for _, v := range *videos {
					id := v.Id
					_, i, err := ctx.Youtube.Get(id)
					if err != nil {
						ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						fmt.Println("error getting video2,", err)
						continue
					}
					video, err := ctx.Youtube.Video(*i)
					if err != nil {
						ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("error"), true)
						fmt.Println("error getting video3,", err)
						return
					}
					song := bot.NewSong(video.Media, video.Title, arg)
					sess.Queue.Add(song)
					ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), fmt.Sprintf(ctx.Loc("youtube_added_format"), song.Title), true)
					if !is_playing {
						shortPlay(&ctx, sess, msg)
						is_playing = true
					}
				}
				ctx.EditEmbed(msg.ID, fmt.Sprintf("%v:", ctx.Loc("youtube")), ctx.Loc("youtube_added"), true)
				break
			}
		}
	}
}
