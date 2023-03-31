package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
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

	if strings.HasPrefix(m.Content, ".") {
		search := m.Content[1:]
		output := make(chan string)
		var msg *discordgo.Message
		ans := ""

		msg, _ = s.ChannelMessageSend(m.ChannelID, "Processing...")

		go c2gptapi.ChatWithGPT(search, output)
		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case value, ok := <-output:
				if !ok {
					ticker.Stop()
					_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, ans)
					return
				}
				ans += value
			case <-ticker.C:
				_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, ans)
			}
		}
	}
}
