package cmd

import (
	"fmt"
	"github.com/FlameInTheDark/dtbot/bot"
	"gopkg.in/robfig/cron.v2"
	"strconv"
	"strings"
)

func CronCommand(ctx bot.Context) {
	if ctx.GetRoles().ExistsName("bot.admin") {
		// !cron add 0 0 7 * * * !w Chelyabinsk
		switch ctx.Args[0] {
		case "add":
			if len(ctx.Args) > 7 {
				if ctx.Args[0] != "*" && ctx.Args[1] != "*" && ctx.Args[2] != "*" {
					switch ctx.Args[7] {
					case "!w":
						if len(ctx.Args) > 7 {
							if !ctx.Data.CronIsFull(&ctx) {
								cmd := strings.Join(ctx.Args[1:], " ")
								cronTime := strings.Join(ctx.Args[1:7], " ")
								ctx.Args = ctx.Args[8:]
								id, _ := ctx.Cron.AddFunc(cronTime, func() { WeatherCommand(ctx) })
								ctx.Data.AddCronJob(&ctx, id, cmd)
								ctx.ReplyEmbedPM("Cron", fmt.Sprintf("Job added: [%v] [%v]", cmd, id))
							} else {
								ctx.ReplyEmbedPM("Cron", "Schedule is full")
							}
						}
					}
				}
			}
		case "remove":
			val,err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				ctx.ReplyEmbedPM("Cron", err.Error())
				return
			}
			cerr := ctx.Data.CronRemove(&ctx, cron.EntryID(val))
			if cerr != nil {
				ctx.ReplyEmbedPM("Cron", err.Error())
				return
			}
			ctx.ReplyEmbedPM("Cron", "Job removed")
		case "list":
			s, err := ctx.Data.CronList(&ctx)
			if err != nil {
				ctx.ReplyEmbedPM("Cron", err.Error())
				return
			}
			var reply = []string{"Jobs:"}
			for key,val := range s.CronJobs {
				reply = append(reply, fmt.Sprintf("Job [%v] ID [%v]", val, key))
			}
			ctx.ReplyEmbedPM("Cron", strings.Join(reply,"\n"))
		}
	}
}