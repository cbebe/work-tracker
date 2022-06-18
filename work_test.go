package worktracker_test

import (
	"bytes"
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

func TestWork(t *testing.T) {
	t.Run("it prints n/a on an invalid Record type", func(t *testing.T) {
		work := Work{
			RecordType: 2,
		}
		assert.Contains(t, work.String(), "n/a")
	})

	t.Run("it prints both start and stop", func(t *testing.T) {
		works := []Work{
			{RecordType: Start},
			{RecordType: Stop},
		}
		stdout := &bytes.Buffer{}
		PrintWorks(stdout, works)
		assert.Contains(t, stdout.String(), "start")
		assert.Contains(t, stdout.String(), "stop")
	})
}
