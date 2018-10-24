package cmd

import (
	"../bot"
)

// HelpCommand shows help
func HelpCommand(ctx bot.Context) {
	if len(ctx.Args) == 0 {
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_reply"))
		return
	}
	switch ctx.Args[0] {
	case "!v":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!v"))
	case "!b":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_command_!b"))
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
	}
}
