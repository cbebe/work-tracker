package worktracker

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteWorkStore struct {
	db *gorm.DB
}

func NewSqliteWorkStore(path string) (*SqliteWorkStore, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&WorkRecord{})
	return &SqliteWorkStore{db}, nil
}

func (s *SqliteWorkStore) newWork(r RecordType, t string) error {
	return s.db.Create(&WorkRecord{RecordType: r, Timestamp: time.Now().Unix(), Type: t}).Error
}

func (s *SqliteWorkStore) StartWork() error {
	return s.StartLog("work")
}

func (s *SqliteWorkStore) StopWork() error {
	return s.StopLog("work")
}

func (s *SqliteWorkStore) StartLog(t string) error {
	var workRecord WorkRecord
	result := s.db.Model(&WorkRecord{}).Where("type = ?", t).Order("id desc").Limit(1).Find(&workRecord)
	if result.Error == nil && workRecord.RecordType == Start {
		return &ExistingLogError{work: workRecord}
	}
	return s.newWork(Start, t)
}

func (s *SqliteWorkStore) StopLog(t string) error {
	var workRecord WorkRecord
	result := s.db.Model(&WorkRecord{}).Where("type = ?", t).Order("id desc").Limit(1).Find(&workRecord)
	if result.Error != nil {
		return result.Error
	}
	if workRecord.RecordType != Start {
		return &ExistingLogError{work: workRecord}
	}
	return s.newWork(Stop, t)
}

func (s *SqliteWorkStore) GetWorkType(t string) ([]Work, error) {
	var workRecords []WorkRecord
	result := s.db.Where("type = ?", t).Find(&workRecords)
	if result.Error != nil {
		return nil, result.Error
	}
	works := make([]Work, len(workRecords))
	for i, v := range workRecords {
		works[i] = v
	}
	return works, nil
}

func (s *SqliteWorkStore) GetWork() ([]Work, error) {
	var workRecords []WorkRecord
	result := s.db.Find(&workRecords)
	if result.Error != nil {
		return nil, result.Error
	}
	works := make([]Work, len(workRecords))
	for i, v := range workRecords {
		works[i] = v
	}
	return works, nil
}

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
