package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"xy.com/discordbot/bard"
	"xy.com/discordbot/c2gptapi"
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
		handleGPTCommand(s, m, msg)
	} else if strings.HasPrefix(m.Content, "!") {
		handleTextResponseCommand(s, m, msg)
	} else if strings.HasPrefix(m.Content, "~") {
		handleChatResponseCommand(s, m, msg)
	}
}

func handleGPTCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output := make(chan string)
	builder := &strings.Builder{}

	go c2gptapi.ChatWithGPT(search, output)
	ticker := time.NewTicker(800 * time.Millisecond)

	for {
		select {
		case value, ok := <-output:
			if !ok {
				ticker.Stop()
				_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, builder.String())
				return
			}
			builder.WriteString(value)
		case <-ticker.C:
			_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, builder.String())
		}
	}
}

func handleTextResponseCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output, err := bard.GenerateTextResponse(search)
	if err != nil {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
	} else {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
	}
}

func handleChatResponseCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output, err := bard.GenerateChatResponse(search)
	if err != nil {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
	} else {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
	}
}
