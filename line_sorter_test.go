package worktracker_test

import (
	"testing"
	"time"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

func TestLineSort(t *testing.T) {
	t.Run("it sorts by start date", func(t *testing.T) {
		lines := []Line{
			{Start: time.Unix(3, 0)},
			{Start: time.Unix(2, 0)},
			{Start: time.Unix(1, 0)},
			{Start: time.Unix(4, 0)},
		}
		By(StartDate).Sort(lines)
		for i := 0; i < len(lines)-1; i++ {
			assert.LessOrEqual(t, lines[i].Start, lines[i+1].Start)
		}
	})

}
