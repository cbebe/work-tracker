package worktracker

import (
	"fmt"
	"io"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store interface {
	NewWork(r RecordType, t, u string) error
	GetLatestWork(t, u string) (Work, error)
	GetWorkType(t, u string) ([]Work, error)
	GetWork(u string) ([]Work, error)
}

func GetPath(w io.Writer) string {
	p := os.Getenv("DB_PATH")
	if p == "" {
		fmt.Fprintln(w, "Using default db path:", DefaultDBPath)
		return DefaultDBPath
	}
	return p
}

func NewStore(p string) (Store, error) {
	t := os.Getenv("DB_TYPE")
	switch t {
	default:
		db, err := gorm.Open(sqlite.Open(p), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
		db.AutoMigrate(&Work{})
		return NewGORMWorkStore(db), nil
	}
}
