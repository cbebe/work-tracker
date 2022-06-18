package worktracker

import (
	"time"

	"gorm.io/gorm"
)

type GORMWorkStore struct {
	db *gorm.DB
}

func NewGORMWorkStore(db *gorm.DB) *GORMWorkStore {
	return &GORMWorkStore{db}
}

func (s *GORMWorkStore) NewWork(r RecordType, t, u string) error {
	return s.db.Create(&Work{
		RecordType: r,
		Timestamp:  Timestamp(time.Now().Unix()),
		Type:       t,
		UserID:     u,
	}).Error
}

func (s *GORMWorkStore) GetLatestWork(t, u string) (Work, error) {
	var work Work
	r := s.db.Where(&Work{Type: t, UserID: u}).Order("id desc").Limit(1).Find(&work)
	if work.ID == 0 {
		return work, NewLogDoesNotExistError(t)
	}
	return work, r.Error
}

func (s *GORMWorkStore) GetWorkType(t, u string) ([]Work, error) {
	var works []Work
	r := s.db.Where(&Work{Type: t, UserID: u}).Find(&works)
	return works, r.Error
}

func (s *GORMWorkStore) GetWork(u string) ([]Work, error) {
	var works []Work
	r := s.db.Where(&Work{UserID: u}).Find(&works)
	return works, r.Error
}
