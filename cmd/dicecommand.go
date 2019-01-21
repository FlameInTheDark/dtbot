package cmd

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/FlameInTheDark/dtbot/bot"
)

// DiceCommand handle dice
func DiceCommand(ctx bot.Context) {
	ctx.MetricsCommand("dice")
	if len(ctx.Args) > 0 {
		val, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			return
		}
		if val <= 0 {
			return
		}
		ctx.Reply(fmt.Sprintf("Dice: %v", rand.Intn(val)+1))
	} else {
		ctx.Reply(fmt.Sprintf("Dice: %v", rand.Intn(6)+1))
	}
}
