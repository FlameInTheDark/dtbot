package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"
)

// Context : Bot context structure
type Context struct {
	BotID string

	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string

	DB   *DBWorker
	Cron *cron.Cron

	Conf       *Config
	CmdHandler *CommandHandler
	Sessions   *SessionManager
	Youtube    *Youtube
	BotMsg     *BotMessages
	Data       *DataType
	Guilds     GuildsMap
}

// NewContext create new context
func NewContext(botID string, discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager, youtube *Youtube, botMsg *BotMessages, dataType *DataType, dbWorker *DBWorker, guilds GuildsMap, botCron *cron.Cron) *Context {
	ctx := new(Context)
	ctx.BotID = botID
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	ctx.Youtube = youtube
	ctx.BotMsg = botMsg
	ctx.Data = dataType
	ctx.DB = dbWorker
	ctx.Guilds = guilds
	ctx.Cron = botCron
	return ctx
}

// Loc returns translated string by key
func (ctx *Context) Loc(key string) string {
	// Check if translation exist

	if len(ctx.Conf.Locales[ctx.GetGuild().Language][key]) == 0 {
		return ctx.Conf.Locales["en"][key]
	}
	return ctx.Conf.Locales[ctx.GetGuild().Language][key]
}

// WeatherCode returns unicode symbol of weather font icon
func (ctx *Context) WeatherCode(code string) string {
	return ctx.Conf.WeatherCodes[code]
}

// GetVoiceChannel returns user voice channel
func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}

	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			// Check voice permissions
			perm, err := ctx.Discord.State.UserChannelPermissions(ctx.BotID, state.ChannelID)
			if err != nil {
				ctx.DB.Log("Voice", ctx.Guild.ID, fmt.Sprintf("Error whilst getting bot permissions on guild \"%v\": %v", ctx.Guild.ID, err))
				return nil
			}

			if perm&discordgo.PermissionVoiceConnect != discordgo.PermissionVoiceConnect ||
				perm&discordgo.PermissionVoiceSpeak != discordgo.PermissionVoiceSpeak ||
				perm&0x00000400 != 0x00000400 {
				ctx.DB.Log("Voice", ctx.Guild.ID, fmt.Sprintf("Voice permissions denied on guild \"%v\"", ctx.Guild.ID))
				return nil
			}

			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}

func (ctx *Context) GetGuild() *GuildData {
	if _, ok := ctx.Guilds[ctx.Guild.ID]; !ok {
		newData := &GuildData{
			ID:          ctx.Guild.ID,
			WeatherCity: ctx.Conf.Weather.City,
			NewsCounty:  ctx.Conf.News.Country,
			Language:    ctx.Conf.General.Language,
			Timezone:    ctx.Conf.General.Timezone,
			EmbedColor:  ctx.Conf.General.EmbedColor,
		}
		_ = ctx.DB.DBSession.DB(ctx.DB.DBName).C("guilds").Insert(newData)
		ctx.Guilds[ctx.Guild.ID] = newData
		return ctx.Guilds[ctx.Guild.ID]
	}
	return ctx.Guilds[ctx.Guild.ID]
}

func (ctx *Context) Log(module, guildID, text string) {
	ctx.DB.Log(module, guildID, text)
	ctx.MetricsLog(module)
}
