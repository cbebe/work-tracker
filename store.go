package worktracker

import (
	"fmt"
	"os"
)

type Store interface {
	NewWork(r RecordType, t, u string) error
	GetLatestWork(t, u string) (Work, error)
	GetWorkType(t, u string) ([]Work, error)
	GetWork(u string) ([]Work, error)
}

func NewStore() (Store, error) {
	t := os.Getenv("DB_TYPE")

	switch t {
	default:
		p := os.Getenv("DB_PATH")
		if p == "" {
			if p == "" {
				fmt.Println("Using default db path:", DefaultDBPath)
				p = DefaultDBPath
			}
		}
		return NewSqliteWorkStore(p)
	}
}
