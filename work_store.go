package worktracker

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type WorkStore interface {
	StartWork()
	StopWork()
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

func (s *SqliteWorkStore) newWork(r RecordType) {
	s.db.Create(&WorkRecord{RecordType: r, Timestamp: time.Now().Unix()})
}

func (s *SqliteWorkStore) StartWork() {
	s.newWork(Start)
}

func (s *SqliteWorkStore) StopWork() {
	s.newWork(Stop)
}

type WorkRecord struct {
	gorm.Model
	RecordType RecordType
	Timestamp  int64
}

func (w WorkRecord) GetRecordType() RecordType {
	return w.RecordType
}

func (w WorkRecord) GetTimestamp() int64 {
	return w.Timestamp
}
