package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cbebe/worktracker/pkg/work"
	"github.com/joho/godotenv"
)

const tokenVar = "DISCORD_TOKEN"

type BotService struct {
	*work.WorkService
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	service, err := work.NewWorkService("work.db")

	if err != nil {
		log.Fatalln("Error connecting to database")
	}

	bot := BotService{&service}

	token := os.Getenv(tokenVar)
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session, ", err)
		return
	}

	dg.AddHandler(bot.messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection, ", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func (b *BotService) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := getArgs(m.Content)
	if args[0] != "task" || len(args) < 2 {
		return
	}
	switch args[1] {
	case "get":
		b.getTasks(s, m, args)
	case "start":
		b.startTask(s, m, args)
	case "stop":
		b.stopTask(s, m, args)
	}

}

func getType(args []string) string {
	if len(args) < 3 {
		return "work"
	} else {
		return args[2]
	}
}

func sendAck(s *discordgo.Session, m *discordgo.MessageCreate, e error, t string, action string) {
	if e != nil {
		s.ChannelMessageSend(m.ChannelID, e.Error())
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint(action, t))
	}
}

func (b *BotService) stopTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StopLog(t), t, "Stopped ")
}

func (b *BotService) startTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StartLog(t), t, "Started ")
}

func (b *BotService) getTasks(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var reply string
	var works []work.Work
	var err error
	if len(args) < 3 {
		works, err = b.GetWork()
	} else {
		works, err = b.GetWorkType(args[2])
	}
	if err != nil {
		reply = "Error getting work"
	} else {
		fmt.Println(works)
		for _, work := range works {
			reply += fmt.Sprintln(work)
		}
	}
	s.ChannelMessageSend(m.ChannelID, reply)
}

func getArgs(s string) []string {
	var args []string
	for _, item := range strings.Split(s, " ") {
		str := strings.TrimSpace(item)
		if str != "" {
			args = append(args, str)
		}
	}
	return args
}
