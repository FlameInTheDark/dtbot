package cmd

import (
	"fmt"
	"math/rand"
	"strconv"

	"../bot"
)

// DiceCommand handle dice
func DiceCommand(ctx bot.Context) {
	if len(ctx.Args) > 0 {
		val, err := strconv.Atoi(ctx.Args[0])
		if err != nil {
			return
		}
		if val <= 0 {
			return
		}
		ctx.Reply(fmt.Sprintf("Dice: %v", rand.Intn(val)))
	} else {
		ctx.Reply(fmt.Sprintf("Dice: %v", rand.Int()))
	}
}
