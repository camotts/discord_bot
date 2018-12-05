package coding

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bots/discord_bot/executor"

	"github.com/bwmarrin/discordgo"
)

func AddCodingCommandParser(s *discordgo.Session) {
	messageCreate := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		fmt.Println("Got message!")
		msgLower := strings.ToLower(m.Content)
		if strings.HasPrefix(msgLower, "!coding") {
			r := csv.NewReader(strings.NewReader(m.Content))
			r.Comma = ' '
			fields, err := r.Read()
			if err != nil {
				fmt.Println("Unable to parse command", err)
				return
			}
			codeCommandParser(s, m, fields[1], fields[2:]...)
		}
	}
	s.AddHandler(messageCreate)
}

func codeCommandParser(s *discordgo.Session, m *discordgo.MessageCreate, command string, args ...string) {
	switch command {
	case "random":
		fmt.Println("Generating random coding challenge...")
	case "submit":
		id, err := strconv.Atoi(args[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Not a problem ID")
			return
		}
		if len(m.Attachments) <= 0 {
			s.ChannelMessageSend(m.ChannelID, "No source file attached")
			return
		}
		err = submitChallenge(s, m, id, m.Attachments[0].URL)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Command not found")
	}
}

func submitChallenge(s *discordgo.Session, m *discordgo.MessageCreate, problemID int, attachURL string) error {
	fmt.Println("Generating temp file")
	f, err := ioutil.TempFile("", "golangProblem")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	fmt.Println("Getting code attachment")
	resp, err := http.Get(attachURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Copying file to temp file")
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Executing the go script")
	stdout, errout, err := executor.ExecuteGoScript(f.Name())
	if err != nil {
		return err
	}
	errbuf := new(bytes.Buffer)
	errbuf.ReadFrom(errout)
	strerr := errbuf.String()
	if strerr != "" {
		s.ChannelMessageSend(m.ChannelID, strerr)
		return nil
	}
	outbuf := new(bytes.Buffer)
	outbuf.ReadFrom(stdout)
	strout := outbuf.String()
	fmt.Println("Print restult to channel")
	fmt.Println(strout)
	fmt.Println("Printed =D")
	s.ChannelMessageSend(m.ChannelID, strout)
	return nil
}
