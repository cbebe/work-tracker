package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cbebe/worktracker"
)

const tokenVar = "DISCORD_TOKEN"

type BotService struct {
	*worktracker.WorkService
	userID string
}

func newBotService() BotService {
	store, err := worktracker.NewStore()
	if err != nil {
		log.Fatalf("error creating work store: %v\n", err)
	}
	service := worktracker.NewWorkService(store)

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

func (b *BotService) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	args := getArgs(m.Content)
	if args[0] != "task" || len(args) < 2 {
		return
	}
	switch args[1] {
	case "list":
		fallthrough
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
		return worktracker.DefaultType
	} else {
		return args[2]
	}
}

func sendAck(s *discordgo.Session, m *discordgo.MessageCreate, e error, t string, action string) {
	var message string
	if e != nil {
		message = e.Error()
	} else {
		message = fmt.Sprint(action, t)
	}
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, "send ack:", err)
	}
}

func (b *BotService) id(m *discordgo.MessageCreate) string {
	if m.Author.ID == b.userID {
		return "cli"
	}
	return m.Author.ID
}

func (b *BotService) stopTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StopLog(t, b.id(m)), t, "Stopped ")
}

func (b *BotService) startTask(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	t := getType(args)
	sendAck(s, m, b.StartLog(t, b.id(m)), t, "Started ")
}

func (b *BotService) sortLogs(works []worktracker.Work) ([]worktracker.Line, []worktracker.Work) {
	logs := make(map[string][]worktracker.Work)
	unfinished := make([]worktracker.Work, 0)
	lines := make([]worktracker.Line, 0)
	for _, w := range works {
		logs[w.Type] = append(logs[w.Type], w)
	}

	for k, v := range logs {
		for i := 0; i < len(v); i += 2 {
			if (i + 1) >= len(v) {
				unfinished = append(unfinished, v[i])
				continue
			}
			s := time.Unix(int64(v[i].Timestamp), 0)
			e := time.Unix(int64(v[i+1].Timestamp), 0)
			d := e.Sub(s)
			lines = append(lines, worktracker.Line{
				Start:    s,
				Type:     k,
				Message:  fmt.Sprintf("**%s:** %s to %s - %s", k, format(&s), format(&e), d),
				Duration: d,
			})
		}
	}

	worktracker.By(worktracker.StartDate).Sort(lines)
	return lines, unfinished
}

func format(t *time.Time) string {
	f := "3:04:05 PM"
	if t.Day() != time.Now().Day() {
		f = "Mon 3:04:05 PM"
	}
	return t.Format(f)
}

func (b *BotService) getTasks(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var reply string
	var works []worktracker.Work
	var err error

	if len(args) < 3 {
		works, err = b.GetWork(b.id(m))
	} else {
		works, err = b.GetWorkType(args[2], b.id(m))
	}

	var lines []worktracker.Line
	unfinished := make([]worktracker.Work, 0)
	total := make(map[string]time.Duration)
	if err != nil {
		reply = "Error getting work"
	} else {
		lines, unfinished = b.sortLogs(works)
		for i, l := range lines {
			reply += fmt.Sprintf("%d. %s\n", i+1, l.Message)
			total[l.Type] += l.Duration
		}
	}
	if reply == "" && len(unfinished) == 0 {
		reply = "No logs found"
	}

	embed := discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Logs", m.Author.Username),
		Description: reply,
	}

	if len(total) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Total Time",
			Value: "Per category:",
		})
		for k, v := range total {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   k,
				Value:  v.String(),
				Inline: true,
			})
		}
	}

	if len(unfinished) > 0 {
		v := ""
		for i, w := range unfinished {
			t := time.Unix(int64(w.Timestamp), 0)
			d := time.Unix(time.Now().Unix(), 0).Sub(t)
			v += fmt.Sprintf("%d. **%s**: Started %s - %s\n", i+1, w.Type, format(&t), d)
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Unfinished Logs",
			Value: v,
		})
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)

	if err != nil {
		fmt.Fprintln(os.Stderr, "get tasks:", err)
	}
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
