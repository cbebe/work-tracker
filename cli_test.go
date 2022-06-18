package worktracker_test

import (
	"bytes"
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

type WorkServiceSpy struct {
	*WorkService
	CalledStartWork   bool
	CalledStopWork    bool
	CalledStartLog    string
	CalledStopLog     string
	CalledGetWork     bool
	CalledGetWorkType string
	Type              string
	ID                string
}

func (w *WorkServiceSpy) StartWork() error {
	w.CalledStartWork = true
	return nil
}

func (w *WorkServiceSpy) StopWork() error {
	w.CalledStopWork = true
	return nil
}

func (w *WorkServiceSpy) StartLog(t, u string) error {
	w.CalledStartLog = t
	return nil
}

func (w *WorkServiceSpy) StopLog(t, u string) error {
	w.CalledStopLog = t
	return nil
}

func (w *WorkServiceSpy) GetWork(u string) ([]Work, error) {
	w.CalledGetWork = true
	return []Work{}, nil
}

func (w *WorkServiceSpy) GetWorkType(t, u string) ([]Work, error) {
	w.CalledGetWorkType = t
	return []Work{}, nil
}

func TestHandleCommand_Usage(t *testing.T) {
	t.Run("it prints the usage when no arguments are given", func(t *testing.T) {
		_, err := executeCommand(t)
		assert.Error(t, err)
	})

	t.Run("it prints the usage when an unknown argument is given", func(t *testing.T) {
		_, err := executeCommand(t, "unknown")
		assert.Error(t, err)
	})
}

func TestHandleCommand_Start(t *testing.T) {
	t.Run("it calls StartWork on `start`", func(t *testing.T) {
		s, _ := executeCommand(t, "start")
		assert.True(t, s.CalledStartWork)
	})

	t.Run("it calls StartLog on `start <type>`", func(t *testing.T) {
		wt := "test"
		s, _ := executeCommand(t, "start", wt)
		assert.Equal(t, wt, s.CalledStartLog)
	})
}

func TestHandleCommand_Stop(t *testing.T) {
	t.Run("it calls StopWork on `stop`", func(t *testing.T) {
		s, _ := executeCommand(t, "stop")
		assert.True(t, s.CalledStopWork)
	})

	t.Run("it calls StopWork on `stop <type>`", func(t *testing.T) {
		wt := "test"
		s, _ := executeCommand(t, "stop", wt)
		assert.Equal(t, wt, s.CalledStopLog)
	})
}

func TestHandleCommand_Get(t *testing.T) {
	t.Run("it calls GetWork on `get`", func(t *testing.T) {
		s, _ := executeCommand(t, "get")
		assert.True(t, s.CalledGetWork)
	})

	t.Run("it calls GetWorkType on `get <type>`", func(t *testing.T) {
		wt := "test"
		s, _ := executeCommand(t, "get", wt)
		assert.Equal(t, wt, s.CalledGetWorkType)
	})

	t.Run("it falls through on `list`", func(t *testing.T) {
		s, _ := executeCommand(t, "list")
		assert.True(t, s.CalledGetWork)
	})
}

func executeCommand(t testing.TB, args ...string) (*WorkServiceSpy, error) {
	t.Helper()
	s := &WorkServiceSpy{}
	stdout := &bytes.Buffer{}
	oa := append([]string{"cli"}, args...)
	err := HandleCommand(stdout, oa, s)
	return s, err
}
