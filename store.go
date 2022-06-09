package worktracker

type Store interface {
	NewWork(r RecordType, t, u string) error
	GetLatestWork(t, u string) (Work, error)
	GetWorkType(t, u string) ([]Work, error)
	GetWork(u string) ([]Work, error)
}
