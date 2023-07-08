package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
	"xy.com/discordbot/bard"
	"xy.com/discordbot/chatGPT"
)

func HandleGPTCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output := make(chan string)
	builder := &strings.Builder{}

	go chatGPT.ChatWithGPT(search, output)
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

func HandleTextResponseCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output, err := bard.GenerateTextResponse(search)
	if err != nil {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
	} else {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
	}
}

func HandleChatResponseCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	search := m.Content[1:]
	output, err := bard.GenerateChatResponse(search)
	if err != nil {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprint(err))
	} else {
		_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
	}
}
