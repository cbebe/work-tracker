package worktracker

import (
	"fmt"
)

type ExistingLogError struct {
	work Work
}

type LogDoesNotExistError struct {
	Type string
}

func exists(r RecordType) string {
	switch r {
	case Start:
		return "started:"
	case Stop:
		return "stopped:"
	}
	return "n/a:"
}

func NewExistingLogError(work Work) error {
	return &ExistingLogError{work}
}

func NewLogDoesNotExistError(t string) error {
	return &LogDoesNotExistError{t}
}

func (e *LogDoesNotExistError) Error() string {
	return fmt.Sprintln("log type does not exist:", e.Type)
}

func (e *ExistingLogError) Error() string {
	return fmt.Sprintln("log already", exists(e.work.RecordType), e.work.Type, "at", e.work.Timestamp)
}
