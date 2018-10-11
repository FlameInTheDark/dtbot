package cmd

import (
	"fmt"

	"../api/currency"
	"../bot"
)

// CurrencyCommand Translate handlergt
func CurrencyCommand(ctx bot.Context) {
	ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("currency")), currency.GetCurrency(&ctx), "", false)
}
