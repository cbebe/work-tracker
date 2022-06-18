package worktracker_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

func TestBotService_GetTasks(t *testing.T) {
	t.Run("it is defined", func(t *testing.T) {
		service := NewWorkService(&StoreSpy{})
		bot := NewBotService(service, nil)
		assert.NotNil(t, bot)
	})

	t.Run("it displays a finished log", func(t *testing.T) {
		spy := &StoreSpy{
			[]Work{
				{Type: "work", RecordType: Start, Timestamp: 1655234110},
				{Type: "work", RecordType: Stop, Timestamp: 1655234111},
			},
		}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		bot.GetTasks([]string{"task", "get", "work"}, "", "test")
	})

	t.Run("it displays an unfinished log", func(t *testing.T) {
		d := time.Now().Unix()
		spy := &StoreSpy{
			[]Work{
				{Type: "work", RecordType: Start, Timestamp: Timestamp(d)},
			},
		}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		bot.GetTasks([]string{"task", "get"}, "", "test")
	})

	t.Run("it replies if there is an error", func(t *testing.T) {
		spy := &ErrorStoreSpy{}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		message := bot.GetTasks([]string{"task", "get"}, "", "test")
		assert.Equal(t, message.Description, ERROR_GETTING_WORK)

	})

	t.Run("it replies if no logs found", func(t *testing.T) {
		spy := &StoreSpy{}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		message := bot.GetTasks([]string{"task", "get"}, "", "test")
		assert.Equal(t, message.Description, NO_LOGS_FOUND)
	})
}

type ErrorStoreSpy struct {
	*StoreSpy
}

func (m *ErrorStoreSpy) GetWork(u string) ([]Work, error) {
	return []Work{}, fmt.Errorf("oops")
}

func (m *ErrorStoreSpy) GetLatestWork(t, u string) (Work, error) {
	return Work{}, fmt.Errorf("oops")
}

func (m *ErrorStoreSpy) NewWork(r RecordType, t, u string) error {
	return fmt.Errorf("oops")
}

func TestBotService_NewLog(t *testing.T) {
	t.Run("it starts a new log of type work", func(t *testing.T) {
		spy := &StoreSpy{}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		bot.NewLog([]string{"task", "start", "work"}, "")
		assert.Len(t, spy.Works, 1)
		assert.Equal(t, spy.Works[0].RecordType, Start)
		assert.Equal(t, spy.Works[0].UserID, "")
		assert.Equal(t, spy.Works[0].Type, "work")
	})

	t.Run("it stop an existing log of the default type", func(t *testing.T) {
		spy := &StoreSpy{
			[]Work{
				{Type: DefaultType, RecordType: Start, Timestamp: 1655234110},
			},
		}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		bot.NewLog([]string{"task", "stop"}, "")
		assert.Len(t, spy.Works, 2)
		assert.Equal(t, spy.Works[1].RecordType, Stop)
		assert.Equal(t, spy.Works[1].UserID, "")
		assert.Equal(t, spy.Works[1].Type, DefaultType)
	})

	t.Run("it prints an error", func(t *testing.T) {
		spy := &ErrorStoreSpy{}
		service := NewWorkService(spy)
		bot := NewBotService(service, nil)
		message := bot.NewLog([]string{"task", "start"}, "")
		assert.Contains(t, message, "oops")
	})
}
