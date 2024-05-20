// Package redisstore is a storage backend based on Redis.
package postgrescache

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
	"gorm.io/gorm"
)

var (
	ErrUnsupportedObjectType = fmt.Errorf("unsupported object type")
)

// Storage implements git.Storer interface with Redis as backend.
type Storage struct {
	db              *gorm.DB
	namespacePrefix string

	ObjectStorage
	ReferenceStorage
	ConfigStorage
	ShallowStorage
	IndexStorage
}

// CountLooseRefs implements storage.Storer.
func (s *Storage) CountLooseRefs() (int, error) {
	panic("unimplemented")
}

// Module implements storage.Storer.
func (s *Storage) Module(name string) (storage.Storer, error) {
	panic("unimplemented")
}

// PackRefs implements storage.Storer.
func (s *Storage) PackRefs() error {
	return nil
}

// RemoveReference implements storage.Storer.
func (s *Storage) RemoveReference(plumbing.ReferenceName) error {
	panic("unimplemented")
}

// NewStorage returns a new Redis-based storage.
func NewStorage(db *gorm.DB, ns string) *Storage {
	return &Storage{
		db:              db,
		namespacePrefix: ns,
		ObjectStorage: ObjectStorage{
			db:              db,
			namespacePrefix: ns,
		},
		ReferenceStorage: ReferenceStorage{
			db:              db,
			namespacePrefix: ns,
		},
		ConfigStorage: ConfigStorage{
			db:              db,
			namespacePrefix: ns,
		},
		ShallowStorage: ShallowStorage{
			db:              db,
			namespacePrefix: ns,
		},
		IndexStorage: IndexStorage{
			db:              db,
			namespacePrefix: ns,
		},
	}
}

var errNotSupported = fmt.Errorf("not supported")

// AddAlternate implements storage.Storer.
func (s *Storage) AddAlternate(remote string) error {
	return errNotSupported
}
