package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

// NewDBSession creates new MongoDB instance
func NewDBSession(dbname string) *DBWorker {
	session, err := mgo.Dial(os.Getenv("MONGO_CONN"))
	if err != nil {
		fmt.Printf("Mongo connection error: %v", err)
	}
	count, err := session.DB("dtbot").C("logs").Count()
	if err != nil {
		fmt.Println("DB_ERR: ", err)
	}
	fmt.Printf("Mongo connected\nLogs in base: %v\n", count)
	return &DBWorker{DBSession: session, DBName: dbname}
}

// InitGuilds initialize guilds in database
func (db *DBWorker) InitGuilds(sess *discordgo.Session, conf *Config) *GuildsMap {
	var data = &GuildsMap{Guilds: make(map[string]*GuildData)}
	var loaded, initialized = 0, 0
	for _, guild := range sess.State.Guilds {
		count, err := db.DBSession.DB(db.DBName).C("guilds").Find(bson.M{"id": guild.ID}).Count()
		if err != nil {
			fmt.Println("Mongo: ", err)
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
				fmt.Println("Mongo: ", err)
				continue
			}
			data.Guilds[guild.ID] = newData
			loaded++
		}
	}
	fmt.Printf("Guilds loaded [%v], initialized [%v]\n", loaded, initialized)
	return data
}

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

// GetTwitchStreams returns twitch streams from mongodb
func (db *DBWorker) GetTwitchStreams(guildID string) map[string]*TwitchStream {
	streams := []TwitchStream{}
	err := db.DBSession.DB(db.DBName).C("streams").Find(bson.M{"guild": guildID}).All(&streams)
	if err != nil {
		fmt.Println("Mongo: ", err)
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
		fmt.Println(err.Error())
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
