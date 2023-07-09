package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"xy.com/discordbot/bard"

	"github.com/bwmarrin/discordgo"
	"xy.com/discordbot/chatGPT"
)

// main initializes the Discord bot, sets up the event handlers, and starts the bot
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
		builder := &strings.Builder{}

		// Send an initial message to indicate the bot is processing the command
		msg, _ = s.ChannelMessageSend(m.ChannelID, "Processing...")

		// Start a new Goroutine to chat with the GPT API
		go chatGPT.ChatWithGPT(search, output)
		ticker := time.NewTicker(800 * time.Millisecond)

		// Update the bot's message with the GPT API's response
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

	if strings.HasPrefix(m.Content, "!") {
		msg, _ := s.ChannelMessageSend(m.ChannelID, "Processing...")
		search := m.Content[1:]
		output, err := bard.GenerateTextResponse(search)
		if err != nil {
			_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
		} else {
			_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
		}
	}

	if strings.HasPrefix(m.Content, "~") {
		msg, _ := s.ChannelMessageSend(m.ChannelID, "Processing...")
		search := m.Content[1:]
		output, err := bard.GenerateChatResponse(search)
		if err != nil {
			_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
		} else {
			_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
		}

	}
}
