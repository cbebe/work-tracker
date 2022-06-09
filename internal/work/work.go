package work

import (
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"
)

const (
	ID            = "cli"
	DefaultType   = "work"
	DefaultDBPath = "db/work.sqlite"
)

type Work struct {
	gorm.Model
	Timestamp  Timestamp
	RecordType RecordType
	Type       string `gorm:"default:'work'"`
	UserID     string `gorm:"default:'cli'"`
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

func (w Work) String() string {
	return fmt.Sprintf("%s %s %s", w.RecordType, w.Type, w.Timestamp)
}

func PrintWorks(w io.Writer, works []Work) {
	for _, work := range works {
		fmt.Fprintln(w, work)
	}
}
