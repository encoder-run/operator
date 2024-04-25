package rediscache

import (
	"github.com/go-git/go-git/v5/storage"
	"github.com/redis/go-redis/v9"
)

const (
	modulePrefix = "module"
)

type ModuleStorage struct {
	client          *redis.Client
	namespacePrefix string
}

// Module retrieves or initializes a new module-specific Storage instance.
func (s *ModuleStorage) Module(name string) (storage.Storer, error) {

	// Check if the module exists in Redis.
	exists, err := s.client.SIsMember(ctx, withNamespace(s.namespacePrefix, modulePrefix, ""), name).Result()
	if err != nil {
		return nil, err
	}

	if !exists {
		// If the module does not exist, add it to the set of modules.
		err = s.client.SAdd(ctx, withNamespace(s.namespacePrefix, modulePrefix, ""), name).Err()
		if err != nil {
			return nil, err
		}
	}

	// Initialize a new Storage instance for the module.
	// The module-specific namespace or prefix can be used for module-specific data.
	moduleStorage := NewStorage(s.client.Options(), withNamespace(s.namespacePrefix, modulePrefix, ""))

	return moduleStorage, nil
}
