package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddTypingSimulatorParser(s *discordgo.Session) {
	messageCreate := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		msg := strings.ToLower(m.Content)
		if strings.HasPrefix(msg, "/type") {
			split := strings.Split(msg, " ")
			channel := m.ChannelID
			if len(split) > 2 {
				channel = split[2]
			}

			e := s.ChannelTyping(channel)
			if e != nil {
				fmt.Println(e)
			}
		}
	}
	s.AddHandlerOnce(messageCreate)
}
