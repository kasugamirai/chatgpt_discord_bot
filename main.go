package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"xy.com/discordbot/handlers"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Error: DISCORD_BOT_TOKEN environment variable not set.")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	msg, _ := s.ChannelMessageSend(m.ChannelID, "Processing...")

	if strings.HasPrefix(m.Content, ".") {
		handlers.HandleGPTCommand(s, m, msg)
	} else if strings.HasPrefix(m.Content, "!") {
		handlers.HandleTextResponseCommand(s, m, msg)
	} else if strings.HasPrefix(m.Content, "~") {
		handlers.HandleChatResponseCommand(s, m, msg)
	}
}
