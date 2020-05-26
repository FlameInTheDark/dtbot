package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"os"
	"time"
)

// DBWorker MongoDB instance
type DBWorker struct {
	DBSession *mgo.Session
	DBName    string
}

type dbLog struct {
	Date   time.Time
	Text   string
	Module string
	Guild  string
}

// GuildData contains data about guild settings
type GuildData struct {
	ID          string
	WeatherCity string
	NewsCounty  string
	Language    string
	Timezone    int
	EmbedColor  int
	VoiceVolume float32
	Greeting    string
}

// GuildsMap contains guilds settings
type GuildsMap struct {
	Guilds map[string]*GuildData
}

// RadioStation contains info about radio station
type RadioStation struct {
	Name     string
	URL      string
	Key      string
	Category string
}

type TwitchDBConfig struct {
	Type   string
	Token  string
	Expire time.Time
}

type BlackListElement struct {
	ID string
}

// NewDBSession creates new MongoDB instance
func NewDBSession(dbname string) *DBWorker {
	session, err := mgo.Dial(os.Getenv("MONGO_CONN"))
	if err != nil {
		log.Printf("[Mongo] Mongo connection error: %v", err)
	}
	count, err := session.DB("dtbot").C("logs").Count()
	if err != nil {
		log.Printf("[Mongo] DB_ERR: ", err)
	}
	log.Printf("[Mongo] connected\nLogs in base: %v\n", count)
	return &DBWorker{DBSession: session, DBName: dbname}
}

// InitGuilds initialize guilds in database
func (db *DBWorker) InitGuilds(sess *discordgo.Session, conf *Config) *GuildsMap {
	var data = &GuildsMap{Guilds: make(map[string]*GuildData)}
	var loaded, initialized = 0, 0
	for _, guild := range sess.State.Guilds {
		count, err := db.DBSession.DB(db.DBName).C("guilds").Find(bson.M{"id": guild.ID}).Count()
		if err != nil {
			log.Printf("[Mongo] guilds, DB: %s, Guild: %s, Error: %v\n", db.DBName, guild.ID, err)
		}
		if count == 0 {
			newData := &GuildData{
				ID:          guild.ID,
				WeatherCity: conf.Weather.City,
				NewsCounty:  conf.News.Country,
				Language:    conf.General.Language,
				Timezone:    conf.General.Timezone,
				EmbedColor:  conf.General.EmbedColor,
				VoiceVolume: conf.Voice.Volume,
				Greeting:    "",
			}
			_ = db.DBSession.DB(db.DBName).C("guilds").Insert(newData)
			data.Guilds[guild.ID] = newData
			initialized++
		} else {
			var newData = &GuildData{}
			_ = db.DBSession.DB(db.DBName).C("guilds").Find(bson.M{"id": guild.ID}).One(newData)
			if err != nil {
				log.Printf("[Mongo] guilds, DB: %s, Guild: %s, Error: %v\n", db.DBName, guild.ID, err)
				continue
			}
			data.Guilds[guild.ID] = newData
			loaded++
		}
	}
	log.Printf("[Mongo] Guilds loaded [%v], initialized [%v]\n", loaded, initialized)
	return data
}

// InitNewGuild creates new guild in database
func (db *DBWorker) InitNewGuild(guildID string, conf *Config, data *GuildsMap) {
	newData := &GuildData{
		ID:          guildID,
		WeatherCity: conf.Weather.City,
		NewsCounty:  conf.News.Country,
		Language:    conf.General.Language,
		Timezone:    conf.General.Timezone,
		EmbedColor:  conf.General.EmbedColor,
		VoiceVolume: conf.Voice.Volume,
		Greeting:    "",
	}
	_ = db.DBSession.DB(db.DBName).C("guilds").Insert(newData)
	data.Guilds[guildID] = newData
}

// Log saves log in database
func (db *DBWorker) Log(module, guildID, text string) {
	_ = db.DBSession.DB(db.DBName).C("logs").Insert(dbLog{Date: time.Now(), Text: text, Module: module, Guild: guildID})
}

// LogGet returns last N log rows
func (db *DBWorker) LogGet(count int) []dbLog {
	var log = make([]dbLog, count)
	_ = db.DBSession.DB(db.DBName).C("logs").Find(nil).Sort("-$natural").Limit(count).All(&log)
	return log
}

// Guilds returns guilds collection from mongodb
func (db *DBWorker) Guilds() *mgo.Collection {
	return db.DBSession.DB(db.DBName).C("guilds")
}

func (db *DBWorker) GetTwitchToken() *TwitchDBConfig {
	var token TwitchDBConfig
	err := db.DBSession.DB(db.DBName).C("config").Find(bson.M{"type": "twitch"}).One(&token)
	if err != nil {
		return &token
	}
	return &token
}

func (db *DBWorker) UpdateTwitchToken(token string, expire time.Time) {
	err := db.DBSession.DB(db.DBName).C("config").Update(bson.M{"type": "twitch"}, bson.M{"$set": bson.M{"token": token, "expire": expire}})
	if err != nil {
		log.Println("[Mongo] Update twitch token error: ", err)
		err = db.DBSession.DB(db.DBName).C("config").Insert(TwitchDBConfig{
			Type:   "twitch",
			Token:  token,
			Expire: expire,
		})
		if err != nil {
			log.Println("[Mongo] Update twitch token error: ", err)
		}
	}
}

// GetTwitchStreams returns twitch streams from mongodb
func (db *DBWorker) GetTwitchStreams(guildID string) map[string]*TwitchStream {
	streams := []TwitchStream{}
	err := db.DBSession.DB(db.DBName).C("streams").Find(bson.M{"guild": guildID}).All(&streams)
	if err != nil {
		log.Printf("[Mongo] streams, streams DB: %s, Guild: %s, Error: %v\n", db.DBName, guildID, err)
	}
	var newMap = make(map[string]*TwitchStream)
	for i, s := range streams {
		newMap[s.Login] = &streams[i]
	}
	return newMap
}

// UpdateStream updates stream in mongodb
func (db *DBWorker) UpdateStream(stream *TwitchStream) {
	err := db.DBSession.DB(db.DBName).C("streams").
		Update(
			bson.M{"guild": stream.Guild, "login": stream.Login},
			bson.M{"$set": bson.M{"isonline": stream.IsOnline}})
	if err != nil {
		log.Println("[Mongo] ", err)
	}
}

// AddStream adds stream in mongodb
func (db *DBWorker) AddStream(stream *TwitchStream) {
	_ = db.DBSession.DB(db.DBName).C("streams").Insert(stream)
}

// RemoveStream removes stream from mongodb
func (db *DBWorker) RemoveStream(stream *TwitchStream) {
	_ = db.DBSession.DB(db.DBName).C("streams").Remove(bson.M{"login": stream.Login, "guild": stream.Guild})
}

// GetRadioStations gets stations from database and returns slice of them
func (db *DBWorker) GetRadioStations(category string) []RadioStation {
	stations := []RadioStation{}
	var request = bson.M{}
	if category != "" {
		request = bson.M{"category": category}
	}
	err := db.DBSession.DB(db.DBName).C("stations").Find(request).All(&stations)
	if err != nil {
		log.Printf("[Mongo] stations, DB: %s, Error: %v\n", db.DBName, err)

	}
	return stations
}

// GetRadioStationByKey returns one station by key
func (db *DBWorker) GetRadioStationByKey(key string) (*RadioStation, error) {
	station := RadioStation{}
	err := db.DBSession.DB(db.DBName).C("stations").Find(bson.M{"key": key}).One(&station)
	if err != nil {
		log.Printf("[Mongo] stations, DB: %s, Key: %s, Error: %v\n", db.DBName, key, err)
		return nil, fmt.Errorf("station not found")
	}
	return &station, nil
}

// RemoveRadioStation removes radio station by key
func (db *DBWorker) RemoveRadioStation(key string) error {
	err := db.DBSession.DB(db.DBName).C("stations").Remove(bson.M{"key": key})
	return err
}

// AddRadioStation adds new radio station
func (db *DBWorker) AddRadioStation(name, url, key, category string) error {
	station := RadioStation{Name: name, URL: url, Key: key, Category: category}
	err := db.DBSession.DB(db.DBName).C("stations").Insert(&station)
	return err
}

// GetAlbionPlayers gets players from database
func (db *DBWorker) GetAlbionPlayers() []AlbionPlayerUpdater {
	var players []AlbionPlayerUpdater
	_ = db.DBSession.DB(db.DBName).C("albion").Find(nil).All(&players)
	fmt.Println(len(players))
	return players
}

// AddAlbionPlayer adds new player in database
func (db *DBWorker) AddAlbionPlayer(player *AlbionPlayerUpdater) {
	err := db.DBSession.DB(db.DBName).C("albion").Insert(player)
	if err != nil {
		log.Printf("[Mongo] Error adding Albion player: ", err.Error())
	}
}

// RemoveAlbionPlayer removes player from database
func (db *DBWorker) RemoveAlbionPlayer(id string) {
	err := db.DBSession.DB(db.DBName).C("albion").Remove(bson.M{"userid": id})
	if err != nil {
		log.Printf("[Mongo] Error removing Albion player: ", err.Error())
	}
}

// UpdateAlbionPlayerLast updates last kill of albion player
func (db *DBWorker) UpdateAlbionPlayerLast(userID string, lastKill int64) {
	err := db.DBSession.DB(db.DBName).C("albion").
		Update(
			bson.M{"userid": userID},
			bson.M{"$set": bson.M{"lastkill": lastKill}})
	if err != nil {
		log.Printf("[Mongo] ", err)
	}
}

// GetBlackList gets blacklist from database
func (db *DBWorker) GetBlacklist() *BlackListStruct {
	var (
		blacklist BlackListStruct
		Guilds    []BlackListElement
		Users     []BlackListElement
	)
	_ = db.DBSession.DB(db.DBName).C("blusers").Find(nil).All(&Users)
	_ = db.DBSession.DB(db.DBName).C("blguilds").Find(nil).All(&Guilds)

	for _, g := range Guilds {
		blacklist.Guilds = append(blacklist.Guilds, g.ID)
	}
	for _, u := range Users {
		blacklist.Users = append(blacklist.Users, u.ID)
	}

	return &blacklist
}

// AddBlacklistGuild adds guild in database blacklist
func (db *DBWorker) AddBlacklistGuild(id string) {
	err := db.DBSession.DB(db.DBName).C("blguilds").Insert(BlackListElement{ID: id})
	if err != nil {
		log.Printf("[Mongo] Error adding guild in blacklist: ", err.Error())
	}
}

// AddBlacklistUser adds user in database blacklist
func (db *DBWorker) AddBlacklistUser(id string) {
	err := db.DBSession.DB(db.DBName).C("blusers").Insert(BlackListElement{ID: id})
	if err != nil {
		log.Printf("[Mongo] Error adding user in blacklist: ", err.Error())
	}
}

// RemoveBlacklistGuild removes guild from database blacklist
func (db *DBWorker) RemoveBlacklistGuild(id string) {
	err := db.DBSession.DB(db.DBName).C("blguilds").Remove(bson.M{"id": id})
	if err != nil {
		log.Printf("[Mongo] Error removing guild from blacklist: ", err.Error())
	}
}

// RemoveBlacklistUser removes user from database blacklist
func (db *DBWorker) RemoveBlacklistUser(id string) {
	err := db.DBSession.DB(db.DBName).C("blusers").Remove(bson.M{"id": id})
	if err != nil {
		log.Printf("[Mongo] Error removing user from blacklist: ", err.Error())
	}
}

// GetNewsCountry returns news country string
func (db *DBWorker) GetNewsCountry(guild string) string {
	var dbGuild GuildData
	err := db.DBSession.DB(db.DBName).C("guilds").Find(bson.M{"id": guild}).One(&dbGuild)
	if err != nil {
		log.Printf("[Mongo] Error getting news country: ", err)
	}
	return dbGuild.NewsCounty
}
