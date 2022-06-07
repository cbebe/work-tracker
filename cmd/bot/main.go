package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cbebe/work-tracker/pkg/work"
)

const tokenVar = "DISCORD_TOKEN"

type BotService struct {
	*work.WorkService
	userID string
}

func newBotService() BotService {
	path := os.Getenv("DB_PATH")
	service, err := work.NewWorkService(path)
	if err != nil {
		log.Fatalln("Error connecting to database")
	}

	id := os.Getenv("USER_ID")
	return BotService{service, id}
}

func newDiscordGo(setup func(*discordgo.Session)) (*discordgo.Session, error) {
	token := os.Getenv(tokenVar)
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, fmt.Errorf("error creating Discord session, %v", err)
	}
	setup(dg)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening connection, %v", err)
	}

	return dg, nil
}

func main() {
	bot := newBotService()
	dg, err := newDiscordGo(func(d *discordgo.Session) {
		d.AddHandler(bot.messageCreate)
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func (b BotService) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
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
		return work.DefaultType
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

func (b BotService) id(m *discordgo.MessageCreate) string {
	if m.Author.ID == b.userID {
		return "cli"
	}
	return m.Author.ID
}

func (b BotService) stopTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StopLog(t, b.id(m)), t, "Stopped ")
}

func (b BotService) startTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StartLog(t, b.id(m)), t, "Started ")
}

func (b BotService) getTasks(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var reply string
	var works []work.Work
	var err error
	if len(args) < 3 {
		works, err = b.GetWork(b.id(m))
	} else {
		works, err = b.GetWorkType(args[2], b.id(m))
	}
	if err != nil {
		reply = "Error getting work"
	} else {
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
