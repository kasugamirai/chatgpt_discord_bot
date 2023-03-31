// Package main provides the main application to run the Discord bot
// which interacts with the GPT API
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"xy.com/discordbot/c2gptapi"
)

// main initializes the Discord bot, sets up the event handlers,
// and starts the bot
func main() {
	// Get the Discord bot token from the environment variable
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Error: DISCORD_BOT_TOKEN environment variable not set.")
		return
	}

	// Create a new Discord session with the given bot token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Add the message create event handler to the Discord session
	dg.AddHandler(messageCreate)

	// Open a WebSocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	// Notify the user that the bot is running and waiting for messages
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// Wait for a termination signal (CTRL-C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close the Discord session
	dg.Close()
}

// messageCreate is called every time a new message is created on any channel
// that the authenticated bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Process commands starting with "."
	if strings.HasPrefix(m.Content, ".") {
		search := m.Content[1:]
		output := make(chan string)
		var msg *discordgo.Message
		ans := ""

		// Send an initial message to indicate the bot is processing the command
		msg, _ = s.ChannelMessageSend(m.ChannelID, "Processing...")

		// Start a new Goroutine to chat with the GPT API
		go c2gptapi.ChatWithGPT(search, output)
		cnt := 0

		for value := range output {
			ans += value
			cnt++
			if cnt%12 == 0 {
				_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, ans)
			}
		}
	}
}
