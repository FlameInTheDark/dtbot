package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"strings"
	"time"
)

// Twitch contains streams
type Twitch struct {
	Guilds  map[string]*TwitchGuild
	DB      *DBWorker
	Conf    *Config
	Discord *discordgo.Session
}

// TwitchGuild contains streams from specified guild
type TwitchGuild struct {
	ID      string
	Streams map[string]*TwitchStream
}

// TwitchStream contains stream data
type TwitchStream struct {
	Login          string
	Guild          string
	Channel        string
	IsOnline       bool
	IsCustom       bool
	CustomMessage  string
	CustomImageURI string
}

// TwitchStreamResult contains response of Twitch API for streams
type TwitchStreamResult struct {
	Data []TwitchStreamData `json:"data"`
}

// TwitchStreamData Twitch API response struct
type TwitchStreamData struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	GameID       string `json:"game_id"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	Viewers      int    `json:"viewer_count"`
	Language     string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// TwitchUserResult contains response of Twitch API for users
type TwitchUserResult struct {
	Data []TwitchUserData `json:"data"`
}

// TwitchUserData Twitch API response struct
type TwitchUserData struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Name            string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImgURL   string `json:"profile_image_url"`
	OfflineImgURL   string `json:"offline_image_url"`
	Views           int    `json:"view_count"`
}

// TwitchGameResult contains response of Twitch API for games
type TwitchGameResult struct {
	Data []TwitchGameData `json:"data"`
}

// TwitchUserData Twitch API response struct
type TwitchGameData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ArtURL string `json:"box_art_url"`
}

// TwitchInit makes new instance of twitch api worker
func TwitchInit(session *discordgo.Session, conf *Config, db *DBWorker) *Twitch {
	guilds := make(map[string]*TwitchGuild)
	var counter int
	for _, g := range session.State.Guilds {
		guildStreams := db.GetTwitchStreams(g.ID)
		counter += len(guildStreams)
		guilds[g.ID] = &TwitchGuild{g.ID, guildStreams}
	}
	fmt.Printf("Loaded [%v] streamers\n", counter)
	return &Twitch{guilds, db, conf, session}
}

// Update updates status of streamers and notify
func (t *Twitch) Update() {
	for _,g := range t.Guilds {
		for _, s := range g.Streams {
			timeout := time.Duration(time.Duration(1) * time.Second)
			client := &http.Client{
				Timeout: time.Duration(timeout),
			}
			req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%v", s.Login), nil)
			req.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
			resp, err := client.Do(req)
			var result TwitchStreamResult
			var gameResult TwitchGameResult
			if err == nil {
				err = json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.DB.Log("Twitch", "", "Parsing Twitch API stream error")
					continue
				}
				if len(result.Data) > 0 {
					greq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/games?id=%v", result.Data[0].GameID), nil)
					greq.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
					gresp, gerr := client.Do(greq)
					err = json.NewDecoder(gresp.Body).Decode(&gameResult)
					if gerr != nil {
						t.DB.Log("Twitch", "", "Parsing Twitch API game error")
					}
					if s.IsOnline == false {
						s.IsOnline = true
						t.DB.UpdateStream(s)
						imgUrl := strings.Replace(result.Data[0].ThumbnailURL, "{width}", "320", -1)
						imgUrl = strings.Replace(imgUrl, "{height}", "180", -1)
						emb := NewEmbed(result.Data[0].UserName).
							Field("Title", result.Data[0].Title, false).
							Field("Viewers", fmt.Sprintf("%v", result.Data[0].Viewers), true).
							Field("Game", gameResult.Data[0].Name, true).
							AttachImgURL(imgUrl).
							Color(t.Conf.General.EmbedColor)
						_, _ = t.Discord.ChannelMessageSend(s.Channel, fmt.Sprintf(t.Conf.GetLocaleLang("twitch_online", result.Data[0].Language), result.Data[0].UserName, s.Login))
						_, _ = t.Discord.ChannelMessageSendEmbed(s.Channel, emb.GetEmbed())
					}
				} else {
					if s.IsOnline == true {
						s.IsOnline = false
						t.DB.UpdateStream(s)
					}
				}

			}
		}
	}
}

// AddStreamer adds new streamer to list
func (t *Twitch) AddStreamer(guild, channel, login string) (string, error) {
	if g, ok := t.Guilds[guild]; ok {
		if g.Streams == nil {
			t.Guilds[guild].Streams = make(map[string]*TwitchStream)
		}
		for _, s := range g.Streams {
			if s.Guild == guild && s.Login == login {
				return "", errors.New("streamer already exists")
			}
		}
		timeout := time.Duration(time.Duration(1) * time.Second)
		client := &http.Client{
			Timeout: time.Duration(timeout),
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?login=%v", login), nil)
		req.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
		resp, err := client.Do(req)
		var result TwitchUserResult
		if err == nil {
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				return "", errors.New("parsing streamer error")
			}
			if len(result.Data) > 0 {
				stream := TwitchStream{}
				stream.Login = login
				stream.Channel = channel
				stream.Guild = guild
				t.Guilds[guild].Streams[login] = &stream
				t.DB.AddStream(&stream)
			}
		} else {
			return "", errors.New("getting streamer error")
		}
		return result.Data[0].Name, nil
	}
	return "", errors.New("guild not found")
}

// RemoveStreamer removes streamer from list
func (t *Twitch) RemoveStreamer(login, guild string) error {
	complete := false
	if g, ok := t.Guilds[guild]; ok {
		if g.Streams != nil {
			if t.Guilds[guild].Streams[login] != nil {
				if g.Streams[login].Login == login && g.Streams[login].Guild == guild {
					t.DB.RemoveStream(g.Streams[login])
					delete(t.Guilds[guild].Streams, login)
					complete = true
				}
			}
		}
	} else {
		return errors.New("guild not found")
	}
	if !complete {
		return errors.New("streamer not found")
	}
	return nil
}
