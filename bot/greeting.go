package bot

import "github.com/bwmarrin/discordgo"

func Greeting(discord *discordgo.Session, event *discordgo.GuildMemberAdd, guild *GuildData, conf *Config) {
	if guild.Greeting != "" {
		NewEmbed(conf.GetLocaleLang("", guild.Language))
		_, _ = discord.ChannelMessageSend(event.User.ID, guild.Greeting)
	}
}

func (ctx *Context) AddGreeting(text string) {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = text
}

func (ctx *Context) RemoveGreeting() {
	ctx.Guilds.Guilds[ctx.Guild.ID].Greeting = ""
}
