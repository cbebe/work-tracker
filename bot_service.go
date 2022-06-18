package worktracker

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type BotServiceConfig struct {
	ErrLog io.Writer
	UserID string
}

type BotService struct {
	service BotWorkService
	Config  *BotServiceConfig
}

type BotWorkService interface {
	StartWork() error
	StopWork() error
	StartLog(t, u string) error
	StopLog(t, u string) error
	Store
}

func NewBotService(service BotWorkService, config *BotServiceConfig) *BotService {
	if config == nil {
		config = &BotServiceConfig{}
	}
	if config.ErrLog == nil {
		config.ErrLog = os.Stderr
	}
	return &BotService{service, config}
}

func (b *BotService) GetTasks(args []string, userID string, username string) *discordgo.MessageEmbed {
	var reply string
	var works []Work
	var err error

	if len(args) < 3 {
		works, err = b.service.GetWork(userID)
	} else {
		works, err = b.service.GetWorkType(args[2], userID)
	}

	var lines []Line
	unfinished := make([]Work, 0)
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
		Title:       fmt.Sprintf("%s's Logs", username),
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

	return &embed
}

func (b *BotService) NewLog(args []string, id string) string {
	var err error
	var action string
	t := getType(args)
	if args[1] == "start" {
		err = b.service.StartLog(t, id)
		action = "Started"
	} else if args[1] == "stop" {
		err = b.service.StopLog(t, id)
		action = "Stopped"
	}
	if err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("%s %s", action, t)
	}
}

func getType(args []string) string {
	if len(args) < 3 {
		return DefaultType
	} else {
		return args[2]
	}
}

func (b *BotService) sortLogs(works []Work) ([]Line, []Work) {
	logs := make(map[string][]Work)
	unfinished := make([]Work, 0)
	lines := make([]Line, 0)
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
			lines = append(lines, Line{
				Start:    s,
				Type:     k,
				Message:  fmt.Sprintf("**%s:** %s to %s - %s", k, format(&s), format(&e), d),
				Duration: d,
			})
		}
	}

	By(StartDate).Sort(lines)
	return lines, unfinished
}

func format(t *time.Time) string {
	f := "3:04:05 PM"
	if t.Day() != time.Now().Day() {
		f = "Mon 3:04:05 PM"
	}
	return t.Format(f)
}
