package cmd

import (
	"fmt"

	"../api/currency"
	"../bot"
)

// CurrencyCommand Translate handler
func CurrencyCommand(ctx bot.Context) {
	ctx.ReplyEmbed("", fmt.Sprintf("%v:", ctx.Loc("currency")), currency.GetCurrency(&ctx), "", false)
}
