package work

import "fmt"

type WorkService struct {
	*SqliteWorkStore
}

func NewWorkService(path string) (*WorkService, error) {
	store, err := NewSqliteWorkStore(path)

	if err != nil {
		return nil, fmt.Errorf("Error creating work service: %v", err)
	}
	return &WorkService{SqliteWorkStore: store}, nil
}

func (w WorkService) StartWork() error {
	return w.StartLog(DefaultType, ID)
}

func (w WorkService) StopWork() error {
	return w.StopLog(DefaultType, ID)
}

func (w WorkService) StartLog(t, u string) error {
	work, err := w.GetLatestWork(t, u)
	if err == nil && work.RecordType == Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Start, t, u)
}

func (w WorkService) StopLog(t, u string) error {
	work, err := w.GetLatestWork(t, u)
	if err != nil {
		return err
	}
	if work.RecordType != Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Stop, t, u)
}
