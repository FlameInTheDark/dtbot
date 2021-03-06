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

// SearchResult contains search response
type AlbionSearchResult struct {
	Guilds  []AlbionGuildSearch  `json:"guilds"`
	Players []AlbionPlayerSearch `json:"players"`
}

// GuildSearch contains guild data from search response
type AlbionGuildSearch struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	KillFame     int    `json:"KillFame"`
	DeathFame    int    `json:"DeathFame"`
}

// PlayerSearch contains player data from search response
type AlbionPlayerSearch struct {
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
type AlbionPlayer struct {
	AverageItemPower float64         `json:"AverageItemPower"`
	Equipment        AlbionEquipment `json:"Equipment"`
	Inventory        []AlbionItem    `json:"Inventory"`
	Name             string          `json:"Name"`
	Id               string          `json:"Id"`
	GuildName        string          `json:"GuildName"`
	GuildId          string          `json:"GuildId"`
	AllianceName     string          `json:"AllianceName"`
	AllianceId       string          `json:"AllianceId"`
	AllianceTag      string          `json:"AllianceTag"`
	Avatar           string          `json:"Avatar"`
	AvatarRing       string          `json:"AvatarRing"`
	DeathFame        int             `json:"DeathFame"`
	KillFame         int             `json:"KillFame"`
	FameRatio        float64         `json:"FameRatio"`
	DamageDone       float64         `json:"DamageDone"`
}

// Equipment contains items in slots
type AlbionEquipment struct {
	MainHand AlbionItem `json:"MainHand"`
	OffHand  AlbionItem `json:"OffHand"`
	Head     AlbionItem `json:"Head"`
	Armor    AlbionItem `json:"Armor"`
	Shoes    AlbionItem `json:"Shoes"`
	Bag      AlbionItem `json:"Bag"`
	Cape     AlbionItem `json:"Cape"`
	Mount    AlbionItem `json:"Mount"`
	Potion   AlbionItem `json:"Potion"`
	Food     AlbionItem `json:"Food"`
}

// Item contains item data
type AlbionItem struct {
	Type    string `json:"Type"`
	Count   int    `json:"Count"`
	Quality int    `json:"Quality"`
}

// Kill contains kill data
type AlbionKill struct {
	GroupMemberCount     int            `json:"groupMemberCount"`
	NumberOfParticipants int            `json:"numberOfParticipants"`
	EventID              int            `json:"EventId"`
	TimeStamp            string         `json:"TimeStamp"`
	Version              int            `json:"Version"`
	Killer               AlbionPlayer   `json:"Killer"`
	Victim               AlbionPlayer   `json:"Victim"`
	TotalVictimKillFame  int            `json:"TotalVictimKillFame"`
	Location             string         `json:"Location"`
	Participants         []AlbionPlayer `json:"Participants"`
	GroupMembers         []AlbionPlayer `json:"GroupMembers"`
	BattleID             int            `json:"BattleId"`
	Type                 string         `json:"Type"`
}

type AlbionUpdater struct {
	Players map[string]*AlbionPlayerUpdater
}

type AlbionPlayerUpdater struct {
	PlayerID string
	UserID   string
	Language string
	LastKill int64
	StartAt  int64
}

// SearchPlayers returns player list by name
func AlbionSearchPlayers(name string) (result *AlbionSearchResult, err error) {
	var sResult AlbionSearchResult
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
func AlbionGetPlayerKills(id string) (result []AlbionKill, err error) {
	var kills []AlbionKill
	resp, err := http.Get(fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/players/%v/topkills?range=month&offset=0&limit=20", id))
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

// AlbionGetKillID returns kill by ID
func AlbionGetKillID(id string) (kill *AlbionKill, err error) {
	var result AlbionKill
	resp, err := http.Get(fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/events/%v", id))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ShowKills sends embed message in discord
func (ctx *Context) AlbionShowKills() {
	search, err := AlbionSearchPlayers(ctx.Args[1])
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}
	fmt.Println("Founded players")
	if len(search.Players) > 0 {
		fmt.Println("Players more then 0")
		kills, err := AlbionGetPlayerKills(search.Players[0].ID)
		fmt.Println("Searching kills of " + search.Players[0].Name + search.Players[0].ID)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		fmt.Println("Founded kills of " + search.Players[0].Name)
		if len(kills) > 0 {
			fmt.Println("Kills more then 0")
			embed := NewEmbed("Albion Killboard")
			embed.Desc(fmt.Sprintf("[%v](https://albiononline.com/ru/killboard/player/%v)", search.Players[0].Name, search.Players[0].ID)) // "https://assets.albiononline.com/assets/images/icons/favicon.ico")
			embed.Color(ctx.GuildConf().EmbedColor)
			for _, k := range kills {
				fmt.Println("Killed " + k.Victim.Name)
				var timeString string
				t, err := time.Parse(time.RFC3339Nano, k.TimeStamp)
				if err == nil {
					timeString = fmt.Sprintf("%v.%v.%v %v:%v", t.Day(), t.Month().String(), t.Year(), t.Hour(), t.Minute())
				} else {
					fmt.Println("Parse time: ", err.Error())
				}
				embed.Field(
					k.Victim.Name,
					fmt.Sprintf("%v [[%v](https://albiononline.com/ru/killboard/kill/%v)]",
						fmt.Sprintf(ctx.Loc("albion_kill_short"),
							k.Victim.DeathFame,
							k.Victim.AverageItemPower,
							timeString),
						k.EventID,
						k.EventID), false)
			}
			embed.Send(ctx)
		}

	}
}

// AlbionShowKill sends kill embed to user
func (ctx *Context) AlbionShowKill() {
	kill, err := AlbionGetKillID(ctx.Args[1])
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}

	embed := NewEmbed(fmt.Sprintf("Show on killboard #%v", kill.EventID))
	embed.Desc(fmt.Sprintf("%v :crossed_swords: %v", kill.Killer.Name, kill.Victim.Name))
	embed.Color(ctx.GuildConf().EmbedColor)
	embed.URL(fmt.Sprintf("https://albiononline.com/ru/killboard/kill/%v", kill.EventID))
	embed.AttachThumbURL("https://assets.albiononline.com/assets/images/header/logo.png")
	embed.Author("Albion Killboard", "https://albiononline.com/ru/killboard", "")
	embed.TimeStamp(kill.TimeStamp)
	embed.Field(ctx.Loc("albion_guild"), kill.Victim.GuildName, true)
	embed.Field(ctx.Loc("albion_fame"), fmt.Sprintf("%d", kill.Victim.DeathFame), true)
	embed.Field(ctx.Loc("albion_item_power"), fmt.Sprintf("%.3f", kill.Victim.AverageItemPower), true)
	embed.Field(ctx.Loc("albion_killer_item_power"), fmt.Sprintf("%.3f", kill.Killer.AverageItemPower), true)
	if len(kill.Participants) > 0 {
		var names []string
		for _, p := range kill.Participants {
			names = append(names, fmt.Sprintf("%v (%.0f)", p.Name, p.DamageDone))
		}
		embed.Field(ctx.Loc("albion_participants"), strings.Join(names, ", "), true)
	}
	embed.Send(ctx)
}

// SendKill sends kill to user
func SendKill(session *discordgo.Session, conf *Config, kill *AlbionKill, userID, lang string) {
	embed := NewEmbed(fmt.Sprintf("Show on killboard #%v", kill.EventID))
	embed.Desc(fmt.Sprintf("%v :crossed_swords: %v", kill.Killer.Name, kill.Victim.Name))
	embed.Color(4460547)
	embed.URL(fmt.Sprintf("https://albiononline.com/ru/killboard/kill/%v", kill.EventID))
	embed.AttachThumbURL("https://assets.albiononline.com/assets/images/header/logo.png")
	embed.Author("Albion Killboard", "https://albiononline.com/ru/killboard", "")
	embed.TimeStamp(kill.TimeStamp)
	embed.Field(conf.GetLocaleLang("albion_item_power", lang), fmt.Sprintf("%.2f", kill.Victim.AverageItemPower), true)
	embed.Field(conf.GetLocaleLang("albion_killer_item_power", lang), fmt.Sprintf("%.2f", kill.Killer.AverageItemPower), true)
	embed.Field(conf.GetLocaleLang("albion_fame", lang), fmt.Sprintf("%d", kill.Victim.DeathFame), true)
	if kill.Victim.GuildName != "" {
		embed.Field(conf.GetLocaleLang("albion_guild", lang), kill.Victim.GuildName, true)
	}
	if len(kill.Participants) > 0 {
		var names []string
		for _, p := range kill.Participants {
			names = append(names, fmt.Sprintf("%v (%.0f)", p.Name, p.DamageDone))
		}
		participants := strings.Join(names, ", ")
		if len(participants) < 1000 {
			embed.Field(conf.GetLocaleLang("albion_participants", lang), participants, true)
		}
	}
	ch, err := session.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("Error whilst creating private channel, ", err.Error())
		return
	}
	_, err = session.ChannelMessageSendEmbed(ch.ID, embed.GetEmbed())
	if err != nil {
		fmt.Println("Error whilst sending embed message, ", err.Error())
		return
	}
}

// GetPlayerByID returns player ID by player name
func GetPlayerByName(name string) string {
	search, err := AlbionSearchPlayers(name)
	if err == nil {
		if len(search.Players) > 0 {
			return search.Players[0].ID
		}
	}
	return ""
}

// AlbionGetUpdater creates and returns albion kills updater
func AlbionGetUpdater(db *DBWorker) *AlbionUpdater {
	var updater = &AlbionUpdater{Players: make(map[string]*AlbionPlayerUpdater)}
	var players []AlbionPlayerUpdater
	players = db.GetAlbionPlayers()
	for i, p := range players {
		updater.Players[p.UserID] = &players[i]
	}
	return updater
}

// SendPlayerKills sends player kills
func SendPlayerKills(session *discordgo.Session, worker *DBWorker, conf *Config, updater *AlbionUpdater, userID string) {
	startTime := time.Unix(updater.Players[userID].StartAt, 0)
	lastTime := time.Unix(updater.Players[userID].LastKill, 0)
	if startTime.Add(time.Hour * 24).Unix() < time.Now().Unix() {
		worker.RemoveAlbionPlayer(updater.Players[userID].UserID)
		delete(updater.Players, updater.Players[userID].UserID)
		return
	} else {
		kills, err := AlbionGetPlayerKills(updater.Players[userID].PlayerID)
		if err != nil {
			return
		}
		var newKillTime int64
		for i, k := range kills {
			killTime, err := time.Parse(time.RFC3339Nano, k.TimeStamp)
			if err != nil {
				fmt.Println("Kill time parse error: ", err.Error())
				continue
			}
			if killTime.Unix() > lastTime.Unix() {
				if killTime.Unix() > newKillTime {
					newKillTime = killTime.Unix()
				}
				SendKill(session, conf, &kills[i], userID, updater.Players[userID].Language)
			}
		}
		if newKillTime > lastTime.Unix() {
			worker.UpdateAlbionPlayerLast(userID, newKillTime)
			updater.Players[userID].LastKill = newKillTime
		}
	}
}

// Update updates players kills and sends to users
func (u *AlbionUpdater) Update(session *discordgo.Session, worker *DBWorker, conf *Config) {
	for _, p := range u.Players {
		startTime := time.Unix(p.StartAt, 0)
		lastTime := time.Unix(p.LastKill, 0)
		if startTime.Add(time.Hour * 24).Unix() < time.Now().Unix() {
			worker.RemoveAlbionPlayer(p.UserID)
			delete(u.Players, p.UserID)
			return
		} else {
			kills, err := AlbionGetPlayerKills(p.PlayerID)
			if err != nil {
				worker.Log("albion", "", fmt.Sprintf("Getting player kills error: %v", err.Error()))
				fmt.Println("Getting player kills error: ", err.Error())
				return
			}
			var newKillTime int64
			for i, k := range kills {
				killTime, err := time.Parse(time.RFC3339Nano, k.TimeStamp)
				if err != nil {
					fmt.Println("Kill time parse error: ", err.Error())
					worker.Log("albion", "", fmt.Sprintf("Parse time error: %v", err.Error()))
					continue
				}
				if killTime.Unix() > lastTime.Unix() {
					if killTime.Unix() > newKillTime {
						newKillTime = killTime.Unix()
					}
					go SendKill(session, conf, &kills[i], p.UserID, p.Language)
				}
			}
			if newKillTime > lastTime.Unix() {
				worker.UpdateAlbionPlayerLast(p.UserID, newKillTime)
				u.Players[p.UserID].LastKill = newKillTime
			}
		}
	}
}

// AlbionAddPlayer adds player to updater
func (ctx *Context) AlbionAddPlayer() error {
	if len(ctx.Args) > 1 {
		search, err := AlbionSearchPlayers(ctx.Args[1])
		if err != nil {
			ctx.Log("albion", "", fmt.Sprintf("Searching player error: %v", err.Error()))
			return errors.New("error searching Albion player")
		}
		if _, ok := ctx.Albion.Players[ctx.User.ID]; !ok {
			kills, err := AlbionGetPlayerKills(search.Players[0].ID)
			if err != nil {
				ctx.Log("albion", "", fmt.Sprintf("Getting kills error: %v", err.Error()))
				return errors.New("error getting Albion kills")
			}
			var lastKill int64
			for _, k := range kills {
				killTime, err := time.Parse(time.RFC3339Nano, k.TimeStamp)
				if err != nil {
					continue
				}
				if killTime.Unix() > lastKill {
					lastKill = killTime.Unix()
				}
			}
			player := &AlbionPlayerUpdater{search.Players[0].ID, ctx.User.ID, ctx.GuildConf().Language, lastKill, time.Now().Unix()}
			ctx.Albion.Players[ctx.User.ID] = player
			ctx.DB.AddAlbionPlayer(player)
			return nil
		}
	}
	return errors.New("error")
}
