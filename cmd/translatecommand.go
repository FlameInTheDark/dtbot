package cmd

import (
	"../api/translate"
	"../bot"
)

// Translate command handler
func TranslateCommand(ctx bot.Context) {
	ctx.Reply(translate.GetTranslation(&ctx))
}
