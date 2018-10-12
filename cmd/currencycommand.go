package cmd

import (
	"fmt"

	"../api/currency"
	"../bot"
)

// CurrencyCommand Translate handlergt
func CurrencyCommand(ctx bot.Context) {
	bot.NewEmbed("").
		Field(fmt.Sprintf("%v:", ctx.Loc("currency")), currency.GetCurrency(&ctx), false).
		Color(0x00ff00).
		Footer(ctx.Loc("requested_by") + ": " + ctx.User.Username).
		Send(ctx)
}
