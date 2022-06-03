package work

import (
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"
)

type Work struct {
	gorm.Model
	Timestamp  Timestamp
	RecordType RecordType
	Type       string `gorm:"default:'work'"`
}

type Timestamp int64

func (t Timestamp) String() string {
	return time.Unix(int64(t), 0).Format(time.UnixDate)
}

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
		fmt.Fprintf(w, "%s %s %s\n", work.RecordType.String(), work.Type, work.Timestamp)
	}
}
