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
}

// NewContext create new context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	return ctx
}

// Reply reply on massege
func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

// ReplyEmbed reply on message with embed message
func (ctx Context) ReplyEmbed(title, field, content, description string, inline bool) *discordgo.Message {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00,
		Description: description,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   field,
				Value:  content,
				Inline: inline,
			},
		},
		Title: title,
	}
	msg, err := ctx.Discord.ChannelMessageSendEmbed(ctx.TextChannel.ID, embed)
	if err != nil {
		fmt.Println("Error whilst sending embed message, ", err)
		return nil
	}
	return msg
}

func (ctx Context) ReplyEmbedAttachment(title, field, content, description, fileName string, file io.Reader, inline bool) *discordgo.Message {
	mess := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x00ff00,
			Description: description,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   field,
					Value:  content,
					Inline: inline,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + fileName,
			},
			Title: title,
		},
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   fileName,
				Reader: file,
			},
		},
	}

	msg, err := ctx.Discord.ChannelMessageSendComplex(ctx.TextChannel.ID, mess)
	if err != nil {
		fmt.Println("Error whilst sending embed message, ", err)
		return nil
	}
	return msg
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
