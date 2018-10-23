package cmd

import (
	"fmt"

	"../api/translate"
	"../bot"
)

// TranslateCommand Translate handler
func TranslateCommand(ctx bot.Context) {
	resp, err := translate.GetTranslation(&ctx)
	if err != nil {
		bot.NewEmbed("").
			Color(0xff0000).
			Field(fmt.Sprintf("%v:", ctx.Loc("translate_error")), err.Error(), false).
			Footer(ctx.Loc("requested_by") + ": " + ctx.User.Username).
			Send(&ctx)
	} else {
		ctx.ReplyEmbed(fmt.Sprintf("%v: ", ctx.Loc("translate")), resp)
	}
}
