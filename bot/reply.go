package bot

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

// NewEmbedStruct generated embed
type NewEmbedStruct struct {
	Emb discordgo.MessageSend
}

// NewEmbed creates new embed
func NewEmbed(title string) NewEmbedStruct {
	var emb = NewEmbedStruct{
		Emb: discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: title,
				Color: 0x00ff00,
			},
		},
	}
	return emb
}

// Field adds field to embed
func (emb NewEmbedStruct) Field(name, value string, inline bool) NewEmbedStruct {
	emb.Emb.Embed.Fields = append(emb.Emb.Embed.Fields, &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline})
	return emb
}

// Desc adds description to embed
func (emb NewEmbedStruct) Desc(desc string) NewEmbedStruct {
	emb.Emb.Embed.Description = desc
	return emb
}

// Color adds color to embed
func (emb NewEmbedStruct) Color(color int) NewEmbedStruct {
	emb.Emb.Embed.Color = color
	return emb
}

// AttachImg adds attached image to embed
func (emb NewEmbedStruct) AttachImg(name string, file io.Reader) NewEmbedStruct {
	fmt.Println(name)
	emb.Emb.Embed.Image = &discordgo.MessageEmbedImage{URL: "attachment://" + name}
	emb.Emb.Files = append(emb.Emb.Files, &discordgo.File{Name: name, Reader: file})
	return emb
}

// Send send embed message to Discord
func (emb NewEmbedStruct) Send(ctx Context) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSendComplex(ctx.TextChannel.ID, &emb.Emb)
	if err != nil {
		fmt.Println("Error whilst sending embed message, ", err)
		return nil
	}
	return msg
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
func (ctx Context) ReplyEmbed(name, content string) *discordgo.Message {
	return NewEmbed("").Field(name, content, false).Send(ctx)
}

// ReplyEmbedAttachment reply on message with embed message with attachment
func (ctx Context) ReplyEmbedAttachment(name, content, fileName string, file io.Reader) *discordgo.Message {
	fmt.Println(fileName)
	return NewEmbed("").Field(name, content, false).AttachImg(fileName, file).Send(ctx)
}
