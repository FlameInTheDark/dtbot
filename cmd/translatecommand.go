package cmd

import (
	"../api/translate"
	"../bot"
)

// TranslateCommand Translate handler
func TranslateCommand(ctx bot.Context) {
	ctx.Reply(translate.GetTranslation(&ctx))
}
