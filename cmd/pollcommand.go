package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FlameInTheDark/dtbot/bot"
)

// PollCommand handle polls commands
func PollCommand(ctx bot.Context) {
	if len(ctx.Args) == 0 {
		return
	}
	switch ctx.Args[0] {
	case "new":
		pollNew(&ctx)
	case "vote":
		pollVote(&ctx)
	case "end":
		pollEnd(&ctx)
	}
}

func pollNew(ctx *bot.Context) {
	ctx.MetricsCommand("poll", "new")
	err := ctx.Data.CreatePoll(ctx, strings.Split(strings.Join(ctx.Args[1:], " "), "|"))
	if err != nil {
		ctx.ReplyEmbed(ctx.Loc("polls"), err.Error())
		return
	}
	fields := strings.Split(strings.Join(ctx.Args[1:], " "), "|")
	for key, val := range fields {
		fields[key] = fmt.Sprintf("%v: %v", key+1, val)
	}
	ctx.ReplyEmbed(ctx.Loc("polls"), fmt.Sprintf("%v:\n%v", ctx.Loc("polls_created"), strings.Join(fields, "\n")))
}

func pollVote(ctx *bot.Context) {
	ctx.MetricsCommand("poll", "vote")
	val, err := strconv.Atoi(ctx.Args[1])
	if err != nil {
		ctx.ReplyEmbed(ctx.Loc("polls"), ctx.Loc("polls_wrong_field"))
		return
	}
	verr := ctx.Data.AddPollVote(ctx, val)
	if verr != nil {
		ctx.ReplyEmbed(ctx.Loc("polls"), verr.Error())
		return
	}
}

func pollEnd(ctx *bot.Context) {
	ctx.MetricsCommand("poll", "end")
	result, err := ctx.Data.EndPoll(ctx)
	if err != nil {
		ctx.ReplyEmbed(ctx.Loc("polls"), err.Error())
		return
	}
	var newResults []string
	for name, count := range result {
		newResults = append(newResults, fmt.Sprintf("[%v] %v", count, name))
	}
	ctx.ReplyEmbed(ctx.Loc("polls"), fmt.Sprintf("%v:\n%v", ctx.Loc("polls_ends"), strings.Join(newResults, "\n")))
}