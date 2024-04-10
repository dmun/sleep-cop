package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Starting sleep-cop..")

	token := os.Getenv("SLEEP_COP_TOKEN")
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Sleepcop is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("ready!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)

	if strings.HasPrefix(m.Content, "bruh") {
		c, err := s.State.Channel(m.ChannelID)

		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)

		if err != nil {
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := s.ChannelVoiceJoin(c.GuildID, vs.ChannelID, false, false)

				if err != nil {
					return
				}

				time.Sleep(5 * time.Second)

				err = vc.Disconnect()

				if err != nil {
					return
				}
			}
		}
	}
}
