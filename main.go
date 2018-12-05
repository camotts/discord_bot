package main

import (
	"bots/discord_bot/reputation"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creationg Discord Session", err)
		os.Exit(1)
	}
	//dg.AddHandler(messageCreate)
	//coding.AddCodingCommandParser(dg)
	reputation.AddReputationHandler(dg)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection", err)
		os.Exit(1)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	msg := strings.ToLower(m.Content)
	if strings.HasPrefix(msg, "ping") {
		for i := 0; i < strings.Count(msg, "s")+1; i++ {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
	}

	if strings.HasPrefix(msg, "pong") {
		for i := 0; i < strings.Count(msg, "s")+1; i++ {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}
	}
}
