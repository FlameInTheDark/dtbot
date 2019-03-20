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
	switch ctx.Args[0] {
	case "!v":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!v"))
	case "!b":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!b"))
		if (ctx.IsAdmin()) {
			ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!b_admin"))
		}
	case "!y":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!y"))
	case "!r":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!r"))
	case "!w":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!w"))
	case "!n":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!n"))
	case "!t":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!t"))
	case "!c":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!c"))
	case "!p":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!p"))
	case "!geoip":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!geoip"))
	case "!twitch":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!twitch"))
	case "bot.admin":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("admin_help"))
	}
}
