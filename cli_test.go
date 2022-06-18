package worktracker_test

import (
	"bytes"
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

type WorkServiceSpy struct {
	*WorkService
	CalledStartWork bool
	CalledStopWork  bool
	Type            string
	ID              string
}

func (w *WorkServiceSpy) StartWork() error {
	w.CalledStartWork = true
	return nil
}

func (w *WorkServiceSpy) StopWork() error {
	w.CalledStopWork = true
	return nil
}

func TestCLI(t *testing.T) {
	t.Run("it prints the usage when no arguments are given", func(t *testing.T) {
		spy := StoreSpy{}
		s := NewWorkService(&spy)
		stdout := &bytes.Buffer{}
		args := []string{"test"}
		err := HandleCommand(stdout, args, s)
		assert.Error(t, err)
	})

	t.Run("it prints the usage when an unknown argument is given", func(t *testing.T) {
		spy := StoreSpy{}
		s := NewWorkService(&spy)
		stdout := &bytes.Buffer{}
		args := []string{"test", "unknown"}
		err := HandleCommand(stdout, args, s)
		assert.Error(t, err)
	})

	t.Run("it calls StartWork on `start`", func(t *testing.T) {
		s := &WorkServiceSpy{}
		stdout := &bytes.Buffer{}
		args := []string{"test", "start"}
		HandleCommand(stdout, args, s)
		assert.True(t, s.CalledStartWork)
	})

	t.Run("it calls StopWork on `start`", func(t *testing.T) {
		s := &WorkServiceSpy{}
		stdout := &bytes.Buffer{}
		args := []string{"test", "stop"}
		HandleCommand(stdout, args, s)
		assert.True(t, s.CalledStopWork)
	})
}
