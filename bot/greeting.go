package bot

import "github.com/bwmarrin/discordgo"

func Greetings(discord *discordgo.Session, event *discordgo.GuildMemberAdd, guild *GuildData, conf *Config) {
	if guild.Greeting != "" {
		_, _ = discord.ChannelMessageSend(event.User.ID, guild.Greeting)
	}
}

func (ctx *Context) AddGreetings(text string) {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = text
}

func (ctx *Context) RemoveGreetings() {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = ""
}
