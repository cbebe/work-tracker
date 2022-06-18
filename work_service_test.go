package worktracker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StoreSpy struct {
	Works []Work
}

func (s StoreSpy) GetLatestWork(t, u string) (Work, error) {
	for i := len(s.Works) - 1; i >= 0; i-- {
		if s.Works[i].Type == t && s.Works[i].UserID == u {
			return s.Works[i], nil
		}
	}

	return Work{}, fmt.Errorf("Not found")
}

func (s *StoreSpy) NewWork(r RecordType, t, u string) error {
	s.Works = append(s.Works, Work{
		UserID:     u,
		RecordType: r,
		Type:       t,
	})
	return nil
}

func (s *StoreSpy) GetWorkType(t, u string) ([]Work, error) {
	works := make([]Work, 0)
	for _, w := range s.Works {
		if w.Type == t && w.UserID == u {
			works = append(works, w)
		}
	}

	return works, nil
}

func (s *StoreSpy) GetWork(u string) ([]Work, error) {
	works := make([]Work, 0)
	for _, w := range s.Works {
		if w.UserID == u {
			works = append(works, w)
		}
	}

	return works, nil
}

func TestWorkService_Start(t *testing.T) {
	t.Run("create start", func(t *testing.T) {
		spy := StoreSpy{}
		service := NewWorkService(&spy)
		expected := Work{
			RecordType: Start,
			Type:       "test",
			UserID:     "1",
		}
		err := service.StartLog(expected.Type, expected.UserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, spy.Works[0])
	})

	t.Run("start error when a log already exists", func(t *testing.T) {
		spy := StoreSpy{}
		service := NewWorkService(&spy)
		expected := Work{
			RecordType: Start,
			Type:       DefaultType,
			UserID:     ID,
		}
		err := service.StartWork()
		assert.NoError(t, err)
		assert.Equal(t, expected, spy.Works[0])
		err = service.StartWork()
		assert.Error(t, err)
	})

	t.Run("create stop", func(t *testing.T) {
		spy := StoreSpy{}
		service := NewWorkService(&spy)
		expected := Work{
			RecordType: Stop,
			Type:       "test",
			UserID:     "1",
		}
		err := service.StartLog(expected.Type, expected.UserID)
		assert.NoError(t, err)
		err = service.StopLog(expected.Type, expected.UserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, spy.Works[1])
	})

	t.Run("stop error when a log does not exist", func(t *testing.T) {
		spy := StoreSpy{}
		service := NewWorkService(&spy)
		err := service.StopWork()
		assert.Error(t, err)
	})
}
