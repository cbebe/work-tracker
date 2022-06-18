package worktracker

type WorkService struct {
	Store
}

func NewWorkService(store Store) *WorkService {
	return &WorkService{store}
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
		return NewExistingLogError(work)
	}
	return w.NewWork(Stop, t, u)
}
