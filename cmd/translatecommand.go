package cmd

import (
	"fmt"

	"../api/translate"
	"../bot"
)

// TranslateCommand Translate handler
func TranslateCommand(ctx bot.Context) {
	ctx.ReplyEmbed(fmt.Sprintf("%v: ", ctx.Loc("translate")), translate.GetTranslation(&ctx))
}
