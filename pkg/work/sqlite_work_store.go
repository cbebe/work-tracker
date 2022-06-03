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

func (s *SqliteWorkStore) NewWork(r RecordType, t string) error {
	return s.db.Create(&Work{RecordType: r, Timestamp: Timestamp(time.Now().Unix()), Type: t}).Error
}

func (s *SqliteWorkStore) GetLatestWork(t string) (Work, error) {
	var work Work
	result := s.db.Model(&Work{}).Where("type = ?", t).Order("id desc").Limit(1).Find(&work)
	if work.ID == 0 {
		return work, &LogDoesNotExistError{t}
	}
	return work, result.Error
}

func (s *SqliteWorkStore) GetWorkType(t string) ([]Work, error) {
	var works []Work
	result := s.db.Where("type = ?", t).Find(&works)
	if result.Error != nil {
		return nil, result.Error
	}
	return works, nil
}

func (s *SqliteWorkStore) GetWork() ([]Work, error) {
	var works []Work
	result := s.db.Find(&works)
	if result.Error != nil {
		return nil, result.Error
	}
	return works, nil
}
