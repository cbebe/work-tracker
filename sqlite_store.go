package worktracker

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteStore struct {
	db *gorm.DB
}

func NewSqliteWorkStore(path string) (*SqliteStore, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Work{})
	return &SqliteStore{db}, nil
}

func (s *SqliteStore) NewWork(r RecordType, t, u string) error {
	return s.db.Create(&Work{
		RecordType: r,
		Timestamp:  Timestamp(time.Now().Unix()),
		Type:       t,
		UserID:     u,
	}).Error
}

func (s *SqliteStore) GetLatestWork(t, u string) (Work, error) {
	var work Work
	r := s.db.Where(&Work{Type: t, UserID: u}).Order("id desc").Limit(1).Find(&work)
	if work.ID == 0 {
		return work, &LogDoesNotExistError{t}
	}
	return work, r.Error
}

func (s *SqliteStore) GetWorkType(t, u string) ([]Work, error) {
	var works []Work
	r := s.db.Where(&Work{Type: t, UserID: u}).Find(&works)
	if r.Error != nil {
		return nil, r.Error
	}
	return works, nil
}

func (s *SqliteStore) GetWork(u string) ([]Work, error) {
	var works []Work
	r := s.db.Where(&Work{UserID: u}).Find(&works)
	if r.Error != nil {
		return nil, r.Error
	}
	return works, nil
}
