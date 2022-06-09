package work

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
	db.AutoMigrate(&Work{})
	return &SqliteWorkStore{db}, nil
}

func (s *SqliteWorkStore) NewWork(r RecordType, t, u string) error {
	return s.db.Create(&Work{
		RecordType: r,
		Timestamp:  Timestamp(time.Now().Unix()),
		Type:       t,
		UserID:     u,
	}).Error
}

func (s *SqliteWorkStore) GetLatestWork(t, u string) (Work, error) {
	var work Work
	result := s.db.Where(&Work{Type: t, UserID: u}).Order("id desc").Limit(1).Find(&work)
	if work.ID == 0 {
		return work, &LogDoesNotExistError{t}
	}
	return work, result.Error
}

func (s *SqliteWorkStore) GetWorkType(t, u string) ([]Work, error) {
	var works []Work
	result := s.db.Where(&Work{Type: t, UserID: u}).Find(&works)
	if result.Error != nil {
		return nil, result.Error
	}
	return works, nil
}

func (s *SqliteWorkStore) GetWork(u string) ([]Work, error) {
	var works []Work
	result := s.db.Where(&Work{UserID: u}).Find(&works)
	if result.Error != nil {
		return nil, result.Error
	}
	return works, nil
}
