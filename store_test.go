package worktracker_test

import (
	"bytes"
	"os"
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
)

func TestStore_NewStore(t *testing.T) {
	t.Run("it creates an sqlite database", func(t *testing.T) {
		store, err := NewStore(":memory:")
		assert.NoError(t, err)
		assert.NotNil(t, store)
	})
	t.Run("it fails to create an sqlite database", func(t *testing.T) {
		store, err := NewStore("/")
		assert.Error(t, err)
		assert.Nil(t, store)
	})
}

func TestStore_GetPath(t *testing.T) {
	t.Run("it gets a path from the environment", func(t *testing.T) {
		expected := "test"
		stdout := &bytes.Buffer{}
		os.Setenv("DB_PATH", expected)
		actual := GetPath(stdout)
		assert.Equal(t, expected, actual)
	})
	t.Run("it uses the default path if the env is not set", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		os.Setenv("DB_PATH", "")
		actual := GetPath(stdout)
		assert.Equal(t, DefaultDBPath, actual)
	})
}
