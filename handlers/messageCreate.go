package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefixes := map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message){
		".": HandleGPTCommand,
		"!": HandleTextResponseCommand,
		"~": HandleChatResponseCommand,
	}

	for prefix, handler := range prefixes {
		if strings.HasPrefix(m.Content, prefix) {
			msg, err := s.ChannelMessageSend(m.ChannelID, "Processing...")
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
			handler(s, m, msg)
			return
		}
	}
}
