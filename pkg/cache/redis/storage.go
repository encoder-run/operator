// Package redisstore is a storage backend based on Redis.
package rediscache

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	ErrUnsupportedObjectType = fmt.Errorf("unsupported object type")
	ctx                      = context.Background()
)

// Storage implements git.Storer interface with Redis as backend.
type Storage struct {
	client   *redis.Client
	moduleNS string

	ConfigStorage
	IndexStorage
	ObjectStorage
	ModuleStorage
	ShallowStorage
	ReferenceStorage
}

// NewStorage returns a new Redis-based storage.
func NewStorage(redisOptions *redis.Options, moduleNS string) *Storage {
	c := redis.NewClient(redisOptions)
	return &Storage{
		client:   c,
		moduleNS: moduleNS,

		ConfigStorage:    ConfigStorage{client: c},
		IndexStorage:     IndexStorage{client: c},
		ObjectStorage:    ObjectStorage{client: c},
		ModuleStorage:    ModuleStorage{client: c},
		ShallowStorage:   ShallowStorage{client: c},
		ReferenceStorage: ReferenceStorage{client: c},
	}
}

var errNotSupported = fmt.Errorf("not supported")

func (o *ObjectStorage) ObjectPacks() ([]plumbing.Hash, error) {
	fmt.Printf("ObjectPacks")
	return nil, nil
}
func (o *ObjectStorage) DeleteOldObjectPackAndIndex(plumbing.Hash, time.Time) error {
	print("DeleteOldObjectPackAndIndex")
	return nil
}

// AddAlternate implements storage.Storer.
func (s *Storage) AddAlternate(remote string) error {
	return errNotSupported
}
