package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"xy.com/discordbot/c2gptapi"
)

func main() {
	// Load bot token from environment variable
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Error: DISCORD_BOT_TOKEN environment variable not set.")
		return
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Add event handler for message create event
	dg.AddHandler(messageCreate)

	// Open the Discord session
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	// Wait until CTRL-C or other termination signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clean up before exiting
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message starts with the "!roll" command
	if strings.HasPrefix(m.Content, "!roll") {
		// Roll a six-sided die and send the result as a message
		result := rand.Intn(6) + 1
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You rolled a %d!", result))
	}

	if strings.HasPrefix(m.Content, ".") {
		// send the result as a message
		search := m.Content[1:]
		ret := ""
		res, err := c2gptapi.ChatWithGPT(search)
		if err != nil {
			ret = fmt.Sprintf("Error: %v", err)
		} else {
			ret = res
		}
		s.ChannelMessageSend(m.ChannelID, ret)
	}
}
