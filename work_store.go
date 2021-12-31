package worktracker

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type WorkStore interface {
	StartWork() error
	StopWork() error
	GetWork() ([]Work, error)
}

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

func (s *SqliteWorkStore) newWork(r RecordType) error {
	return s.db.Create(&WorkRecord{RecordType: r, Timestamp: time.Now().Unix()}).Error
}

func (s *SqliteWorkStore) StartWork() error {
	return s.newWork(Start)
}

func (s *SqliteWorkStore) StopWork() error {
	return s.newWork(Stop)
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
	RecordType RecordType
	Timestamp  int64
}

func (w WorkRecord) GetRecordType() RecordType {
	return w.RecordType
}

func (w WorkRecord) GetTimestamp() int64 {
	return w.Timestamp
}
