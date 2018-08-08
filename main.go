package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "./messages"
    "./config"
    "github.com/bwmarrin/discordgo"
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
        fmt.Println("Connection open error")
        return
    }
    fmt.Println("Bot is now running.")

    sc := make(chan os.Signal, 1)

    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    defer dg.Close()
}