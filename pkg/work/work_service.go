package work

type WorkService struct {
	*SqliteWorkStore
}

func (w *WorkService) StartWork() error {
	return w.StartLog("work")
}

func (w *WorkService) StopWork() error {
	return w.StopLog("work")
}

func (w *WorkService) StartLog(t string) error {
	work, err := w.GetLatestWork(t)
	if err == nil && work.RecordType == Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Start, t)
}

func (w *WorkService) StopLog(t string) error {
	work, err := w.GetLatestWork(t)
	if err != nil {
		return err
	}
	if work.RecordType != Start {
		return &ExistingLogError{work}
	}
	return w.NewWork(Stop, t)
}
