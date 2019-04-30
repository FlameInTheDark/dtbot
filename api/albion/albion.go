package albion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"net/http"
	"time"
)

// SearchResult contains search response
type SearchResult struct {
	Guilds  []GuildSearch  `json:"guilds"`
	Players []PlayerSearch `json:"players"`
}

// GuildSearch contains guild data from search response
type GuildSearch struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	KillFame     int    `json:"KillFame"`
	DeathFame    int    `json:"DeathFame"`
}

// PlayerSearch contains player data from search response
type PlayerSearch struct {
	ID           string  `json:"Id"`
	Name         string  `json:"Name"`
	AllianceID   string  `json:"AllianceId"`
	AllianceName string  `json:"AllianceName"`
	KillFame     int     `json:"KillFame"`
	DeathFame    int     `json:"DeathFame"`
	Avatar       string  `json:"Avatar"`
	AvatarRing   string  `json:"AvatarRing"`
	FameRation   float64 `json:"FameRatio"`
	TotalKills   int     `json:"totalKills"`
	GVGKills     int     `json:"gvgKills"`
	GVGWon       int     `json:"gvgWon"`
}

// Player data
type Player struct {
	AverageItemPower float64   `json:"AverageItemPower"`
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
	FameRatio        float64   `json:"FameRatio"`
}

// Equipment contains items in slots
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

// Item contains item data
type Item struct {
	Type    string `json:"Type"`
	Count   int    `json:"Count"`
	Quality int    `json:"Quality"`
}

// Kill contains kill data
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

// SearchPlayers returns player list by name
func SearchPlayers(name string) (result *SearchResult, err error) {
	var sResult SearchResult
	resp, err := http.Get(fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/search?q=%v", name))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&sResult)
	if err != nil {
		return nil, err
	}

	return &sResult, nil
}

// GetPlayerKills returns array of kills by player id
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

// ShowKills sends embed message in discord
func ShowKills(ctx *bot.Context) {
	search, err := SearchPlayers(ctx.Args[1])
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}
	fmt.Println("Founded players")
	if len(search.Players) > 0 {
		fmt.Println("Players more then 0")
		kills, err := GetPlayerKills(search.Players[0].ID)
		fmt.Println("Searching kills of " + search.Players[0].Name + search.Players[0].ID)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		fmt.Println("Founded kills of " + search.Players[0].Name)
		if len(kills) > 0 {
			fmt.Println("Kills more then 0")
			embed := bot.NewEmbed("Albion Killboard")
			embed.Desc(fmt.Sprintf("[%v](https://albiononline.com/ru/killboard/player/%v)", search.Players[0].Name, search.Players[0].ID)) // "https://assets.albiononline.com/assets/images/icons/favicon.ico")
			embed.Color(ctx.GuildConf().EmbedColor)
			for _, k := range kills {
				fmt.Println("Killed " + k.Victim.Name)
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
}
