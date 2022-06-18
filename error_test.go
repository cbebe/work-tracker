package worktracker_test

import (
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("should print n/a on an invalid Record type", func(t *testing.T) {
		err := NewExistingLogError(Work{
			RecordType: 2,
		})
		assert.Contains(t, err.Error(), "n/a")
	})
}
