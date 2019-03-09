package main

import (
	"bytes"
	"fmt"
	"github.com/DiscordBotList/dblgo"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/robfig/cron.v2"

	"github.com/FlameInTheDark/dtbot/bot"
	"github.com/FlameInTheDark/dtbot/cmd"
	"github.com/bwmarrin/discordgo"
)

var (
	conf *bot.Config
	// CmdHandler bot command handler
	CmdHandler *bot.CommandHandler
	// Sessions bot session manager
	Sessions        *bot.SessionManager
	botId           string
	youtube         *bot.Youtube
	botMsg          *bot.BotMessages
	dataType        *bot.DataType
	dbWorker        *bot.DBWorker
	guilds          bot.GuildsMap
	botCron         *cron.Cron
	messagesCounter int
)

func main() {
	botCron = cron.New()
	conf = bot.LoadConfig()
	CmdHandler = bot.NewCommandHandler()
	registerCommands()
	Sessions = bot.NewSessionManager()
	youtube = &bot.Youtube{Conf: conf}
	botMsg = bot.NewMessagesMap()
	dataType = bot.NewDataType()
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Create session error, ", err)
		return
	}
	usr, err := discord.User("@me")
	if err != nil {
		fmt.Println("Error obtaining account details,", err)
		return
	}
	botId = usr.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		_ = discord.UpdateStatus(0, conf.General.Game)
		guilds := discord.State.Guilds
		fmt.Println("Ready with", len(guilds), "guilds.")
	})

	err = discord.Open()
	if err != nil {
		fmt.Printf("Connection open error: %v", err)
		return
	}
	defer discord.Close()
	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	dbWorker = bot.NewDBSession(conf.General.DatabaseName)
	guilds = dbWorker.InitGuilds(discord, conf)
	botCron.Start()
	go MetricsSender(discord)
	defer botCron.Stop()
	defer dbWorker.DBSession.Close()
	<-sc
}

// Handle discord messages
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	messagesCounter++
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}
	args := strings.Split(message.Content, " ")
	name := strings.ToLower(args[0])
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}

	var permission = true
	var msg string
	// Checking permissions
	perm, err := discord.State.UserChannelPermissions(botId, message.ChannelID)
	if err != nil {
		msg = fmt.Sprintf("Error whilst getting bot permissions %v\n", err)
		permission = false
	} else {
		if perm&discordgo.PermissionSendMessages != discordgo.PermissionSendMessages ||
			perm&discordgo.PermissionAttachFiles != discordgo.PermissionAttachFiles {
			msg = "Permissions denied"
			permission = false
		}
	}

	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}

	if permission {
		ctx := bot.NewContext(
			botId,
			discord,
			guild,
			channel,
			user,
			message,
			conf,
			CmdHandler,
			Sessions,
			youtube,
			botMsg,
			dataType,
			dbWorker,
			guilds,
			botCron)
		ctx.Args = args[1:]
		c := *command
		c(*ctx)
	} else {
		dbWorker.Log("Message", guild.ID, msg)
		query := []byte(fmt.Sprintf("logs,server=%v module=\"%v\"", guild.ID, "message"))
		addr := fmt.Sprintf("%v/write?db=%v", conf.Metrics.Address, conf.Metrics.Database)
		r := bytes.NewReader(query)
		_, _ = http.Post(addr, "", r)
	}
}

// Adds bot commands
func registerCommands() {
	CmdHandler.Register("!r", cmd.PlayerCommand)
	CmdHandler.Register("!w", cmd.WeatherCommand)
	CmdHandler.Register("!t", cmd.TranslateCommand)
	CmdHandler.Register("!n", cmd.NewsCommand)
	CmdHandler.Register("!c", cmd.CurrencyCommand)
	CmdHandler.Register("!y", cmd.YoutubeCommand)
	CmdHandler.Register("!v", cmd.VoiceCommand)
	CmdHandler.Register("!b", cmd.BotCommand)
	CmdHandler.Register("!play", cmd.YoutubeShortCommand)
	CmdHandler.Register("!d", cmd.DebugCommand)
	CmdHandler.Register("!p", cmd.PollCommand)
	CmdHandler.Register("!m", cmd.YandexmapCommand)
	CmdHandler.Register("!dice", cmd.DiceCommand)
	CmdHandler.Register("!help", cmd.HelpCommand)
	CmdHandler.Register("!cron", cmd.CronCommand)
	CmdHandler.Register("!geoip", cmd.GeoIPCommand)
}

// MetricsSender sends metrics to InfluxDB
func MetricsSender(d *discordgo.Session) {
	for {
		query := []byte(fmt.Sprintf("messages count=%v", messagesCounter))
		addr := fmt.Sprintf("%v/write?db=%v&u=%v&p=%v",
			conf.Metrics.Address, conf.Metrics.Database, conf.Metrics.User, conf.Metrics.Password)
		r := bytes.NewReader(query)
		_, _ = http.Post(addr, "", r)
		if conf.DBL.Token != "" {
			s := dblgo.NewDBL(conf.DBL.Token, d.State.User.ID)
			_ = s.PostStats(len(d.State.Guilds))
		}
		if conf.DBL.TokenDBL != "" {
			users := 0
			for _, g := range d.State.Guilds {
				users += len(g.Members)
			}
			data := url.Values{}
			data.Set("shard_id", "0")
			data.Set("guilds", string(len(d.State.Guilds)))
			fmt.Println("Guilds: ", string(len(d.State.Guilds)))
			data.Set("users", string(users))
			fmt.Println("Users: ", string(users))
			data.Set("voice_connections", string(Sessions.Count()))
			fmt.Println("Voice connections: ", string(Sessions.Count()))
			client := &http.Client{}
			fmt.Println(data.Encode())
			req, err := http.NewRequest("POST", fmt.Sprintf("https://discordbotlist.com/api/bots/%v/stats",conf.DBL.DBLID), bytes.NewBuffer([]byte(data.Encode())))
			if err != nil {
				fmt.Println(err.Error())
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bot %v", conf.DBL.TokenDBL))
			defer req.Body.Close()
			res, _ := client.Do(req)
			fmt.Println(res.Status)

		}
		messagesCounter = 0
		time.Sleep(time.Minute)
	}
}
