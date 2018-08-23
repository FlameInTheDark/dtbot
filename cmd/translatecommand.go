package cmd

import (
	"../api/translate"
	"../bot"
)

func TranslateCommand(ctx bot.Context) {
	ctx.Reply(translate.GetTranslation(&ctx))
}
