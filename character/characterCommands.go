package character

import (
	"discord_bot/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	players        = make(map[string]player)
	namesToIds     = make(map[string]string)
	creationStates = make(map[string]CreationState)
)

func AddCharacterCommandParser(s *discordgo.Session) {
	messageCreate := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		msg := strings.ToLower(m.Content)
		if strings.HasPrefix(msg, "/character") {
			split := strings.Split(msg, " ")
			characterCommandParser(s, m, split[1], split[2:]...)
		}
		if ok, _ := utils.ComesFromDM(s, m); ok {
			fmt.Println("DM Recieved from ", m.Author.Username)
			directCreate(s, m)
		}
	}
	s.AddHandlerOnce(messageCreate)
}

func characterCommandParser(s *discordgo.Session, m *discordgo.MessageCreate, command string, args ...string) {
	switch command {
	case "create":
		createCharacter(s, m)
	case "stats":
		s.ChannelMessageSend(m.ChannelID, getPlayerStats(m.Author.ID))
	default:
	}
}

func directCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var state CreationState
	if cs, ok := creationStates[m.Author.ID]; ok {
		state = cs
	} else {
		creationStates[m.Author.ID] = CreationStateStart
		state = CreationStateStart
	}
	switch state {
	case CreationStateStart:
		s.ChannelMessageSend(m.ChannelID, "You now have a base character!")
		players[m.Author.ID] = player{
			Name:   "TBD",
			Health: 100,
			Damage: 1,
		}
		creationStates[m.Author.ID] = creationStates[m.Author.ID].Next()
	case CreationStateEnd:
		s.ChannelMessageSend(m.ChannelID, "You have already complete character creation")
	}
	fmt.Println("Finished DM handling")
}

func getPlayerStats(id string) string {
	if p, ok := players[id]; ok {
		return p.String()
	}
	return "You have not initialized a character. Please run the '/character create' command!"
}

func createCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.ChannelMessageSend(ch.ID, "Welcome to the Character Creator! Please type a name.")
	if err != nil {
		fmt.Println(err)
	}
}
