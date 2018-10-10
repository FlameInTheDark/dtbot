package cmd

import (
	"../api/currency"
	"../bot"
)

// CurrencyCommand Translate handler
func CurrencyCommand(ctx bot.Context) {
	ctx.Reply(currency.GetCurrency(&ctx))
}
