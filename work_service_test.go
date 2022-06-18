package worktracker_test

import (
	"testing"

	. "github.com/cbebe/worktracker"
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

	return Work{}, NewLogDoesNotExistError(t)
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

	t.Run("start error when a log has already started", func(t *testing.T) {
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
		expectedErr := NewExistingLogError(expected)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})
}

func TestWorkService_Stop(t *testing.T) {
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
		expected := NewLogDoesNotExistError(DefaultType)
		assert.Contains(t, err.Error(), expected.Error())
	})

	t.Run("stop error when a log has already been stopped", func(t *testing.T) {
		spy := StoreSpy{}
		service := NewWorkService(&spy)
		service.StartWork()
		service.StopWork()
		err := service.StopWork()

		expected := NewExistingLogError(Work{
			RecordType: Stop,
			Type:       DefaultType,
			UserID:     ID,
		})
		assert.Contains(t, err.Error(), expected.Error())
	})
}
