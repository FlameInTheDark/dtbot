package main

import (
	"./config"
	"./messages"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.LoadConfig()
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Create session error, ", err)
		return
	}

	dg.AddHandler(messages.MessageCreate)

	err = dg.Open()
	if err != nil {
        fmt.Printf("Connection open error: %v", err)
		return
	}
	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer dg.Close()
}
