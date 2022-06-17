package worktracker

import "fmt"

type WorkService struct {
	Store
}

func NewWorkService(path string) (*WorkService, error) {
	if path == "" {
		fmt.Println("Using default db path:", DefaultDBPath)
		path = DefaultDBPath
	}
	store, err := NewSqliteWorkStore(path)

	if err != nil {
		return nil, fmt.Errorf("error creating work service: %v", err)
	}
	return &WorkService{store}, nil
}

func (w *WorkService) StartWork() error {
	return w.StartLog(DefaultType, ID)
}

func (w *WorkService) StopWork() error {
	return w.StopLog(DefaultType, ID)
}

func (w *WorkService) StartLog(t, u string) error {
	work, err := w.GetLatestWork(t, u)
	if err == nil && work.RecordType == Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Start, t, u)
}

func (w *WorkService) StopLog(t, u string) error {
	work, err := w.GetLatestWork(t, u)
	if err != nil {
		return err
	}
	if work.RecordType != Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Stop, t, u)
}