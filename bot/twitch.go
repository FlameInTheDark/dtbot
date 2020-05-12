package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
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
	Login           string
	UserID          string
	Name            string
	Guild           string
	Channel         string
	ProfileImageURL string
	IsOnline        bool
	IsCustom        bool
	CustomMessage   string
	CustomImageURL  string
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

// TwitchGameData Twitch API response struct
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

func (t *Twitch) OAuthToken() string {
	dbt := t.DB.GetTwitchToken()
	if dbt.Expire.Unix() < time.Now().Unix() {
		var oauth2Config *clientcredentials.Config
		oauth2Config = &clientcredentials.Config{
			ClientID:     t.Conf.Twitch.ClientID,
			ClientSecret: t.Conf.Twitch.ClientSecret,
			TokenURL:     twitch.Endpoint.TokenURL,
		}

		token, err := oauth2Config.Token(context.Background())
		if err != nil {
			log.Println("[Twitch] Getting token error: ", err)
		}
		t.DB.UpdateTwitchToken(token.AccessToken, token.Expiry)
		return token.AccessToken
	}
	return dbt.Token
}

// Update updates status of streamers and notify
func (t *Twitch) Update() {
	var gameResult TwitchGameResult
	var streamResult TwitchStreamResult
	var streams = make(map[string]*TwitchStreamData)
	var games = make(map[string]*TwitchGameData)
	var oauthToken = t.OAuthToken()
	timeout := time.Duration(time.Duration(1) * time.Second)
	client := &http.Client{
		Timeout: time.Duration(timeout),
	}
	streamQuery := url.Values{}
	gameQuery := url.Values{}
	for _, g := range t.Guilds {
		for _, s := range g.Streams {
			streamQuery.Add("user_login", s.Login)
		}
	}
	// Streams
	tsreq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/streams?%v", streamQuery.Encode()), nil)
	tsreq.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
	tsreq.Header.Add("Authorization", "Bearer "+oauthToken)
	tsresp, tserr := client.Do(tsreq)
	if tserr == nil {
		jerr := json.NewDecoder(tsresp.Body).Decode(&streamResult)
		if jerr != nil {
			t.DB.Log("Twitch", "", "Parsing Twitch API stream error")
		}
	} else {
		t.DB.Log("Twitch", "", fmt.Sprintf("Getting Twitch API stream error: %v", tserr.Error()))
		return
	}
	for i, s := range streamResult.Data {
		gameQuery.Add("id", s.GameID)
		streams[s.UserID] = &streamResult.Data[i]
	}

	// Games
	tgreq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/games?%v", gameQuery.Encode()), nil)
	tgreq.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
	tgreq.Header.Add("Authorization", "Bearer "+oauthToken)
	tgresp, tgerr := client.Do(tgreq)
	if tgerr == nil {
		jerr := json.NewDecoder(tgresp.Body).Decode(&gameResult)
		if jerr != nil {
			t.DB.Log("Twitch", "", "Parsing Twitch API game error")
			return
		}
	} else {
		t.DB.Log("Twitch", "", fmt.Sprintf("Getting Twitch API game error: %v", tgerr.Error()))
		return
	}
	for i, g := range gameResult.Data {
		games[g.ID] = &gameResult.Data[i]
	}

	for _, g := range t.Guilds {
		for _, s := range g.Streams {
			if stream, ok := streams[s.UserID]; ok {
				if !s.IsOnline {
					gameName := "Unknown"
					if _, ok := games[stream.GameID]; ok {
						gameName = games[stream.GameID].Name
					}
					t.Guilds[s.Guild].Streams[s.Login].IsOnline = true
					t.DB.UpdateStream(s)
					imgURL := strings.Replace(stream.ThumbnailURL, "{width}", "320", -1)
					imgURL = strings.Replace(imgURL, "{height}", "180", -1)
					emb := NewEmbed(stream.Title).
						URL(fmt.Sprintf("http://www.twitch.tv/%v", s.Login)).
						Author(s.Name, "", s.ProfileImageURL).
						Field("Viewers", fmt.Sprintf("%v", stream.Viewers), true).
						Field("Game", gameName, true).
						AttachImgURL(imgURL).
						Color(t.Conf.General.EmbedColor)
					if s.CustomMessage != "" {
						emb.Content = s.CustomMessage
					} else {
						emb.Content = fmt.Sprintf(t.Conf.GetLocaleLang("twitch_online", stream.Language), s.Name, s.Login)
					}
					_, _ = t.Discord.ChannelMessageSendComplex(s.Channel, emb.MessageSend)
				}
			} else {
				if s.IsOnline == true {
					t.Guilds[s.Guild].Streams[s.Login].IsOnline = false
					t.DB.UpdateStream(s)
				}
			}
		}
	}
}

// AddStreamer adds new streamer to list
func (t *Twitch) AddStreamer(guild, channel, login, message string) (string, error) {
	if g, ok := t.Guilds[guild]; ok {
		if g.Streams == nil {
			t.Guilds[guild].Streams = make(map[string]*TwitchStream)
		}
		for _, s := range g.Streams {
			if s.Guild == guild && s.Login == login {
				return "", errors.New("streamer already exists")
			}
		}
		var oauthToken = t.OAuthToken()
		timeout := time.Duration(time.Duration(1) * time.Second)
		client := &http.Client{
			Timeout: time.Duration(timeout),
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?login=%v", login), nil)
		req.Header.Add("Client-ID", t.Conf.Twitch.ClientID)
		req.Header.Add("Authorization", "Bearer "+oauthToken)
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
				stream.UserID = result.Data[0].ID
				if result.Data[0].Name == "" {
					stream.Name = login
				} else {
					stream.Name = result.Data[0].Name
				}
				stream.ProfileImageURL = result.Data[0].ProfileImgURL
				stream.CustomMessage = message
				t.Guilds[guild].Streams[login] = &stream
				t.DB.AddStream(&stream)
			} else {
				return "", errors.New("streamer not found")
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
