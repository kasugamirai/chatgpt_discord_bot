package main

import (
	"fmt"
	"os"
	"os/signal"
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
	/*
		commands := []*discordgo.ApplicationCommand{
			{
				Name:        "chatGpt",
				Description: "Interact with GPT-3",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "input",
						Description: "Input for GPT-3",
						Required:    true,
					},
				},
			},
			{
				Name:        "bardText",
				Description: "Generate a text response",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "input",
						Description: "Input for bard text response",
						Required:    true,
					},
				},
			},
			{
				Name:        "bardChat",
				Description: "Generate a chat response",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "input",
						Description: "Input for bard chat response",
						Required:    true,
					},
				},
			},
		}

		for _, command := range commands {
			_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "your guild ID", command)
			if err != nil {
				fmt.Println("Cannot create command: ", err)
				return
			}
		}
	*/

	dg.AddHandler(handlers.MessageCreate)

	if err = dg.Open(); err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
