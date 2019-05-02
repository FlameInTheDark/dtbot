package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"gopkg.in/robfig/cron.v2"
	"strconv"
	"strings"
)

// CronCommand manipulates cron functions
func CronCommand(ctx bot.Context) {
	if ctx.IsServerAdmin() {
		// !cron add 0 0 7 * * * !w Chelyabinsk
		switch ctx.Arg(0) {
		case "add":
			ctx.MetricsCommand("cron", "add")
			if len(ctx.Args) > 7 {
				if ctx.Args[0] != "*" && ctx.Args[1] != "*" && ctx.Args[2] != "*" {
					if len(ctx.Args) > 7 {
						if !ctx.Data.CronIsFull(&ctx) {
							cmd := strings.Join(ctx.Args[1:], " ")
							cronTime := strings.Join(ctx.Args[1:7], " ")
							trigger := ctx.Args[7]
							ctx.Args = ctx.Args[8:]
							id, _ := ctx.Cron.AddFunc(cronTime, func() {
								switch trigger {
								case "!w":
									WeatherCommand(ctx)
								case "!c":
									CurrencyCommand(ctx)
								case "!p":
									PollCommand(ctx)
								case "!v":
									VoiceCommand(ctx)
								case "!y":
									YoutubeCommand(ctx)
								case "!play":
									YoutubeShortCommand(ctx)
								case "!b":
									BotCommand(ctx)
								case "!n":
									NewsCommand(ctx)
								}
							})
							_ = ctx.Data.AddCronJob(&ctx, id, cmd)
							ctx.ReplyEmbedPM("Cron", fmt.Sprintf("Job added: [%v] [%v]", cmd, id))
						} else {
							ctx.ReplyEmbedPM("Cron", "Schedule is full")
						}
					}
				}
			}
		case "remove":
			ctx.MetricsCommand("cron", "remove")
			val, err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				ctx.ReplyEmbedPM("Cron", err.Error())
				return
			}
			cErr := ctx.Data.CronRemove(&ctx, cron.EntryID(val))
			if cErr != nil {
				ctx.ReplyEmbedPM("Cron", "Error adding job")
				fmt.Println("Error adding job: ", cErr.Error())
				return
			}
			ctx.ReplyEmbedPM("Cron", "Job removed")
		case "list":
			ctx.MetricsCommand("cron", "list")
			s, err := ctx.Data.CronList(&ctx)
			if err != nil {
				ctx.ReplyEmbedPM("Cron", err.Error())
				return
			}
			var reply = []string{"Jobs:"}
			for key, val := range s.CronJobs {
				reply = append(reply, fmt.Sprintf("[%v] - [%v]", key, val))
			}
			ctx.ReplyEmbedPM("Cron", strings.Join(reply, "\n"))
		}
	} else {
		ctx.MetricsCommand("cron", "error")
	}
}
