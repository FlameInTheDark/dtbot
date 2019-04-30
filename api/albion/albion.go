package albion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"net/http"
	"time"
)

type SearchResult struct {
	Guilds  []GuildSearch  `json:"guilds"`
	Players []PlayerSearch `json:"players"`
}

type GuildSearch struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	KillFame     int    `json:"KillFame"`
	DeathFame    int    `json:"DeathFame"`
}

type PlayerSearch struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	KillFame     int    `json:"KillFame"`
	DeathFame    int    `json:"DeathFame"`
	Avatar       string `json:"Avatar"`
	AvatarRing   string `json:"AvatarRing"`
	FameRation   int    `json:"FameRatio"`
	TotalKills   int    `json:"totalKills"`
	GVGKills     int    `json:"gvgKills"`
	GVGWon       int    `json:"gvgWon"`
}

type Player struct {
	AverageItemPower int       `json:"AverageItemPower"`
	Equipment        Equipment `json:"Equipment"`
	Inventory        []Item    `json:"Inventory"`
	Name             string    `json:"Name"`
	Id               string    `json:"Id"`
	GuildName        string    `json:"GuildName"`
	GuildId          string    `json:"GuildId"`
	AllianceName     string    `json:"AllianceName"`
	AllianceId       string    `json:"AllianceId"`
	AllianceTag      string    `json:"AllianceTag"`
	Avatar           string    `json:"Avatar"`
	AvatarRing       string    `json:"AvatarRing"`
	DeathFame        int       `json:"DeathFame"`
	KillFame         int       `json:"KillFame"`
	FameRatio        int       `json:"FameRatio"`
}

type Equipment struct {
	MainHand Item `json:"MainHand"`
	OffHand  Item `json:"OffHand"`
	Head     Item `json:"Head"`
	Armor    Item `json:"Armor"`
	Shoes    Item `json:"Shoes"`
	Bag      Item `json:"Bag"`
	Cape     Item `json:"Cape"`
	Mount    Item `json:"Mount"`
	Potion   Item `json:"Potion"`
	Food     Item `json:"Food"`
}

type Item struct {
	Type    string `json:"Type"`
	Count   int    `json:"Count"`
	Quality int    `json:"Quality"`
}

type Kill struct {
	GroupMemberCount     int      `json:"groupMemberCount"`
	NumberOfParticipants int      `json:"numberOfParticipants"`
	EventID              int      `json:"EventId"`
	TimeStamp            string   `json:"TimeStamp"`
	Version              int      `json:"Version"`
	Killer               Player   `json:"Killer"`
	Victim               Player   `json:"Victim"`
	TotalVictimKillFame  int      `json:"TotalVictimKillFame"`
	Location             string   `json:"Location"`
	Participants         []Player `json:"Participants"`
	GroupMembers         []Player `json:"GroupMembers"`
	BattleID             int      `json:"BattleId"`
	Type                 string   `json:"Type"`
}

func SearchPlayers(name string) (result *SearchResult, err error) {
	var sresult SearchResult
	resp, err := http.Get(fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/search?q=%v", name))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&sresult)
	if err != nil {
		return nil, err
	}

	return &sresult, nil
}

func GetPlayerKills(id string) (result []Kill, err error) {
	var kills []Kill
	resp, err := http.Get(fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/players/%v/topkills?range=lastWeek&offset=0&limit=11", id))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&kills)
	if err != nil {
		return nil, err
	}

	return kills, nil
}

func ShowKills(ctx *bot.Context) {
	search, err := SearchPlayers(ctx.Args[1])
	if err != nil {
		return
	}
	if len(search.Players) > 0 {
		kills, err := GetPlayerKills(search.Players[0].ID)
		if err != nil {
			return
		}
		embed := bot.NewEmbed(ctx.Loc("Albion Killboard"))
		embed.Author("https://albiononline.com/ru/killboard/player/"+search.Players[0].ID, "", "https://assets.albiononline.com/assets/images/icons/favicon.ico")
		for _, k := range kills {
			t, err := time.Parse("2006-01-02T15:04:05.000+0000", k.TimeStamp)
			var timeString string
			if err == nil {
				timeString = fmt.Sprintf("%v.%v.%v %v:%v", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
			}
			embed.Field(k.Victim.Name, fmt.Sprintf(ctx.Loc("albion_kill_short"), k.Victim.FameRatio, k.Victim.AverageItemPower, timeString), false)
		}
		embed.Send(ctx)
	}
}
