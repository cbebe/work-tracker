package worktracker

import (
	"fmt"
	"io"
)

type RecordType int8

const (
	Start RecordType = iota
	Stop
)

func (r RecordType) String() string {
	switch r {
	case Start:
		return "start"
	case Stop:
		return "stop"
	}
	return "n/a"
}

func PrintWork(w io.Writer, works []Work) {
	for _, work := range works {
		fmt.Fprintf(w, "%s %d\n", work.GetRecordType().String(), work.GetTimestamp())
	}
}

type Work interface {
	GetRecordType() RecordType
	GetTimestamp() int64
}
