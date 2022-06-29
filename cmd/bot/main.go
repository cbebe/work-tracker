package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cbebe/worktracker"
)

const tokenVar = "DISCORD_TOKEN"

type DiscordBotService struct {
	*worktracker.BotService
}

func initBotService() *DiscordBotService {
	store, err := worktracker.NewStore(worktracker.GetPath(os.Stdout))
	if err != nil {
		log.Fatalf("error creating work store: %v\n", err)
	}
	service := worktracker.NewWorkService(store)
	return &DiscordBotService{worktracker.NewBotService(service, &worktracker.BotServiceConfig{
		UserID: os.Getenv("USER_ID"),
	})}
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
	bot := initBotService()
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

func (b *DiscordBotService) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	args := getArgs(m.Content)
	if strings.ToLower(args[0]) != "task" || len(args) < 2 {
		return
	}
	authorID := b.id(m.Author.ID)
	cmd := strings.ToLower(args[1])
	if cmd == "list" || cmd == "get" {
		embed := b.GetTasks(args, authorID, m.Author.Username)
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Fprintln(os.Stderr, "get tasks:", err)
		}
	} else {
		message := b.NewLog(args, authorID)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Fprintln(b.Config.ErrLog, "send ack:", err)
		}
	}
}

func (b *DiscordBotService) id(authorID string) string {
	if authorID == b.Config.UserID {
		return worktracker.ID
	}
	return authorID
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
