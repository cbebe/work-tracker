package work

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
		return "started"
	case Stop:
		return "stopped"
	}
	return "n/a"
}

func (e *LogDoesNotExistError) Error() string {
	return fmt.Sprintf("log type does not exist: %s", e.Type)
}

func (e *ExistingLogError) Error() string {
	return fmt.Sprintf("log already %s: %s at %s", exists(e.work.RecordType), e.work.Type, e.work.Timestamp)
}
