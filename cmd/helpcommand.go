package cmd

import (
	"github.com/FlameInTheDark/dtbot/bot"
)

// HelpCommand shows help
func HelpCommand(ctx bot.Context) {
	ctx.MetricsCommand("help_command", "main")
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_reply"))
		return
	}

	commandMap := map[string]string{
		"v":         "help_command_!v",
		"b":         "help_command_!b",
		"y":         "help_command_!y",
		"r":         "help_command_!r",
		"w":         "help_command_!w",
		"n":         "help_command_!n",
		"t":         "help_command_!t",
		"c":         "help_command_!c",
		"p":         "help_command_!p",
		"geoip":     "help_command_!geoip",
		"twitch":    "help_command_!twitch",
		"greetings": "help_command_!greetings",
		"bot.admin": "admin_help",
	}

	adminCommandMap := map[string]string{
		"b": "help_command_!b_admin",
	}

	if _, ok := commandMap[ctx.Args[0]]; !ok {
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_reply"))
	}

	ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc(commandMap[ctx.Args[0]]))

	if ctx.IsAdmin() {
		if _, ok := adminCommandMap[ctx.Args[0]]; ok {
			ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc(adminCommandMap[ctx.Args[0]]))
		}
	}
}
