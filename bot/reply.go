package bot

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

// NewEmbedStruct generated embed
type NewEmbedStruct struct {
	*discordgo.MessageSend
}

// NewEmbed creates new embed
func NewEmbed(title string) *NewEmbedStruct {
	return &NewEmbedStruct{&discordgo.MessageSend{Embed: &discordgo.MessageEmbed{Title: title}}}
}

// Field adds field to embed
func (emb *NewEmbedStruct) Field(name, value string, inline bool) *NewEmbedStruct {
	emb.Embed.Fields = append(emb.Embed.Fields, &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline})
	return emb
}

// Desc adds description to embed
func (emb *NewEmbedStruct) Desc(desc string) *NewEmbedStruct {
	emb.Embed.Description = desc
	return emb
}

// Footer adds footer text
func (emb *NewEmbedStruct) Footer(text string) *NewEmbedStruct {
	emb.Embed.Footer = &discordgo.MessageEmbedFooter{Text: text}
	return emb
}

// Color adds color to embed
func (emb *NewEmbedStruct) Color(color int) *NewEmbedStruct {
	emb.Embed.Color = color
	return emb
}

// AttachImg adds attached image to embed frim io.Reader
func (emb *NewEmbedStruct) AttachImg(name string, file io.Reader) *NewEmbedStruct {
	emb.Embed.Image = &discordgo.MessageEmbedImage{URL: "attachment://" + name}
	emb.Files = append(emb.Files, &discordgo.File{Name: name, Reader: file})
	return emb
}

// AttachImgURL adds attached image to embed from url
func (emb *NewEmbedStruct) AttachImgURL(url string) *NewEmbedStruct {
	emb.Embed.Image = &discordgo.MessageEmbedImage{URL: url}
	return emb
}

// AttachThumbURL adds attached thumbnail to embed from url
func (emb *NewEmbedStruct) AttachThumbURL(url string) *NewEmbedStruct {
	emb.Embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
	return emb
}

// Send send embed message to Discord
func (emb *NewEmbedStruct) Send(ctx *Context) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSendComplex(ctx.TextChannel.ID, emb.MessageSend)
	if err != nil {
		fmt.Println("Error whilst sending embed message, ", err)
		return nil
	}
	ctx.BotMsg.Add(ctx, msg.ID)
	return msg
}

// GetEmbed returns discords embed
func (emb *NewEmbedStruct) GetEmbed() *discordgo.MessageEmbed {
	return emb.Embed
}

// Reply reply on massege
func (ctx *Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	ctx.BotMsg.Add(ctx, msg.ID)
	return msg
}

// EditEmbed edits embed by id
func (ctx *Context) EditEmbed(ID, name, value string, inline bool) {
	ctx.Discord.ChannelMessageEditEmbed(ctx.TextChannel.ID, ID, NewEmbed("").Color(ctx.Conf.General.EmbedColor).Footer(fmt.Sprintf("%v %v", ctx.Loc("requested_by"), ctx.User.Username)).Field(name, value, inline).GetEmbed())
}

// ReplyEmbed reply on message with embed message
func (ctx *Context) ReplyEmbed(name, content string) *discordgo.Message {
	return NewEmbed("").
		Field(name, content, false).
		Footer(ctx.Loc("requested_by") + ": " + ctx.User.Username).
		Color(ctx.Conf.General.EmbedColor).
		Send(ctx)
}

// ReplyEmbedAttachment reply on message with embed message with attachment
func (ctx *Context) ReplyEmbedAttachment(name, content, fileName string, file io.Reader) *discordgo.Message {
	return NewEmbed("").
		Field(name, content, false).
		AttachImg(fileName, file).
		Footer(ctx.Loc("requested_by") + ": " + ctx.User.Username).
		Color(ctx.Conf.General.EmbedColor).
		Send(ctx)
}
