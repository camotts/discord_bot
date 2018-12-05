package reputation

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	userPoints = make(map[string]reputablePerson)
	mtx        = sync.Mutex{}
)

type reputablePerson struct {
	Reputation int
	NextSend   time.Time
	Recipients map[string]time.Time
}

func AddReputationHandler(s *discordgo.Session) {
	messageCreate := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		msgLower := strings.ToLower(m.Content)
		mtx.Lock()
		defer mtx.Unlock()

		if strings.HasPrefix(msgLower, "thanks") {
			if _, ok := userPoints[m.Author.ID]; ok {

			}
			for _, user := range m.Mentions {
				if m.Author.ID == user.ID {
					continue
				}
				if ur, ok := userPoints[user.ID]; ok {
					ur.Reputation++
					userPoints[user.ID] = ur
				} else {

					userPoints[user.ID] = reputablePerson{
						Reputation: 1,
					}
				}
			}
		}

		if strings.HasPrefix(msgLower, "bad") {
			for _, user := range m.Mentions {
				if m.Author.ID == user.ID {
					continue
				}
				if ur, ok := userPoints[user.ID]; ok {
					ur.Reputation--
					userPoints[user.ID] = ur
				} else {
					userPoints[user.ID] = reputablePerson{
						Reputation: -1,
					}
				}
			}
		}

		if strings.HasPrefix(msgLower, "myrep") {
			g, _ := s.State.Channel(m.ChannelID)
			mem, _ := s.GuildMember(g.GuildID, m.Author.ID)
			if rp, ok := userPoints[m.Author.ID]; ok {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d Rep Points", mem.Nick, rp.Reputation))
			} else {
				userPoints[m.Author.ID] = reputablePerson{
					Reputation: 0,
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has 0 Rep Points", mem.Nick))
			}
		}
	}
	s.AddHandler(messageCreate)
}
