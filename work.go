package worktracker

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

type Work interface {
	GetRecordType() RecordType
	GetTimestamp() int64
}
