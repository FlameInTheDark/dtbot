package cmd

import (
	"fmt"

	"github.com/FlameInTheDark/dtbot/api/currency"
	"github.com/FlameInTheDark/dtbot/bot"
)

// CurrencyCommand Translate handlergt
func CurrencyCommand(ctx bot.Context) {
	ctx.ReplyEmbed(fmt.Sprintf("%v:", ctx.Loc("currency")), currency.GetCurrency(&ctx))
}
