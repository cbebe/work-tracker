package work

import "gorm.io/gorm"

type WorkRecord struct {
	gorm.Model
	Timestamp  int64
	RecordType RecordType
	Type       string `gorm:"default:'work'"`
}

func (w WorkRecord) GetRecordType() RecordType {
	return w.RecordType
}

func (w WorkRecord) GetTimestamp() int64 {
	return w.Timestamp
}

func (w WorkRecord) GetType() string {
	return w.Type
}
