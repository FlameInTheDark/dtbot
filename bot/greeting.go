package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
)

// Greetings sends greetings for user
func Greetings(discord *discordgo.Session, event *discordgo.GuildMemberAdd, guild *GuildData) {
	if guild.Greeting != "" {
		ch, cErr := discord.UserChannelCreate(event.User.ID)
		if cErr != nil {
			return
		}
		_, mErr := discord.ChannelMessageSend(ch.ID, guild.Greeting)
		if mErr != nil {
			return
		}
	}
}

// AddGreetings adds new greetings to guild
func (ctx *Context) AddGreetings(text string) {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = text
	_ = ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"greeting": text}})
}

// RemoveGreetings removes greetings from guild
func (ctx *Context) RemoveGreetings() {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = ""
	_ = ctx.DB.Guilds().Update(bson.M{"id": ctx.Guild.ID}, bson.M{"$set": bson.M{"greeting": ""}})

}
