package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/api/currency"
	"github.com/FlameInTheDark/dtbot/bot"
)

// CurrencyCommand Translate handler
func CurrencyCommand(ctx bot.Context) {
	ctx.MetricsCommand("currency")
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("currency")), currency.GetCurrency(&ctx))
}
