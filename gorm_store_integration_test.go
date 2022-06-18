package worktracker_test

import (
	"testing"

	. "github.com/cbebe/worktracker"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGORMStore_HappyPath(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&Work{})

	store := NewGORMWorkStore(db)

	var first Work

	t.Run("it creates new work", func(t *testing.T) {
		err := store.NewWork(Start, DefaultType, ID)
		assert.NoError(t, err)
		db.Model(&Work{}).Find(&first)
		assert.Equal(t, uint(1), first.ID)
	})

	t.Run("it gets the latest work", func(t *testing.T) {
		work, err := store.GetLatestWork(DefaultType, ID)
		assert.NoError(t, err)
		assert.EqualValues(t, first, work)
	})

	t.Run("it returns LogDoesNotExist if no log is found", func(t *testing.T) {
		work, err := store.GetLatestWork("does not exist", ID)
		assert.Error(t, err)
		assert.Equal(t, uint(0), work.ID)
		assert.IsType(t, &LogDoesNotExistError{}, err)
	})

	t.Run("it gets all work of the default type", func(t *testing.T) {
		err := store.NewWork(Stop, DefaultType, ID)
		assert.NoError(t, err)
		works, err := store.GetWork(ID)
		assert.NoError(t, err)
		assert.Len(t, works, 2)
	})

	t.Run("it returns an empty array if no work is found", func(t *testing.T) {
		works, err := store.GetWorkType("nothing", ID)
		assert.NoError(t, err)
		assert.Len(t, works, 0)
	})
}
