package bot

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

// Context : Bot context structure
type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string

	Conf       *Config
	CmdHandler *CommandHandler
	Sessions   *SessionManager
	Youtube    *Youtube
}

// NewContext create new context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager, youtube *Youtube) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	ctx.Youtube = youtube
	return ctx
}

// Loc Returns translated key string
func (ctx *Context) Loc(key string) string {
	// Check if translation exist
	if len(ctx.Conf.Locales[ctx.Conf.General.Language][key]) == 0 {
		return ctx.Conf.Locales["en"][key]
	}
	return ctx.Conf.Locales[ctx.Conf.General.Language][key]
}

// ReplyFile reply on massege with file
func (ctx Context) ReplyFile(name string, r io.Reader) *discordgo.Message {
	msg, err := ctx.Discord.ChannelFileSend(ctx.TextChannel.ID, name, r)
	if err != nil {
		fmt.Println("Error whilst sending file,", err)
		return nil
	}
	return msg
}

// GetVoiceChannel returns user voice channel
func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}
