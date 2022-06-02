package worktracker

import (
	"fmt"
	"time"
)

type WorkStore interface {
	StartWork() error
	StartLog(string) error
	StopWork() error
	StopLog(string) error
	GetWork() ([]Work, error)
	GetWorkType(string) ([]Work, error)
}

type ExistingLogError struct {
	work Work
}

func exists(r RecordType) string {
	switch r {
	case Start:
		return "started"
	case Stop:
		return "stopped"
	}
	return "n/a"
}

func (e *ExistingLogError) Error() string {
	return fmt.Sprintf("log already %s: %s at %v", exists(e.work.GetRecordType()), e.work.GetType(), time.Unix(e.work.GetTimestamp(), 0))
}
