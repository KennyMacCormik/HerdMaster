package storage

type DB interface {
	Close() error
}
