// Package redisstore is a storage backend based on Redis.
package rediscache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	ErrUnsupportedObjectType = fmt.Errorf("unsupported object type")
	ctx                      = context.Background()
)

// Storage implements git.Storer interface with Redis as backend.
type Storage struct {
	client          *redis.Client
	namespacePrefix string

	ConfigStorage
	IndexStorage
	ObjectStorage
	ModuleStorage
	ShallowStorage
	ReferenceStorage
}

// NewStorage returns a new Redis-based storage.
func NewStorage(redisOptions *redis.Options, ns string) *Storage {
	c := redis.NewClient(redisOptions)
	return &Storage{
		client:          c,
		namespacePrefix: ns,

		ConfigStorage: ConfigStorage{
			client:          c,
			namespacePrefix: ns,
		},
		IndexStorage: IndexStorage{
			client:          c,
			namespacePrefix: ns,
		},
		ObjectStorage: ObjectStorage{
			client:          c,
			namespacePrefix: ns,
		},
		ModuleStorage: ModuleStorage{
			client:          c,
			namespacePrefix: ns,
		},
		ShallowStorage: ShallowStorage{
			client:          c,
			namespacePrefix: ns,
		},
		ReferenceStorage: ReferenceStorage{
			client:          c,
			namespacePrefix: ns,
		},
	}
}

var errNotSupported = fmt.Errorf("not supported")

// AddAlternate implements storage.Storer.
func (s *Storage) AddAlternate(remote string) error {
	return errNotSupported
}

func withNamespace(ns string, keyType string, key string) string {
	k := ns
	if keyType != "" {
		k = fmt.Sprintf("%s:%s", k, keyType)
	}
	if key != "" {
		k = fmt.Sprintf("%s:%s", k, key)
	}
	return k
}
