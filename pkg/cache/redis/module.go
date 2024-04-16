package rediscache

import (
	"fmt"

	"github.com/go-git/go-git/v5/storage"
	"github.com/redis/go-redis/v9"
)

type ModuleStorage struct {
	client *redis.Client
}

// Module retrieves or initializes a new module-specific Storage instance.
func (manager *ModuleStorage) Module(name string) (storage.Storer, error) {
	key := fmt.Sprintf("modules:%s", name)

	// Check if the module exists in Redis.
	exists, err := manager.client.SIsMember(ctx, "modules", name).Result()
	if err != nil {
		return nil, err
	}

	if !exists {
		// If the module does not exist, add it to the set of modules.
		err = manager.client.SAdd(ctx, "modules", name).Err()
		if err != nil {
			return nil, err
		}
	}

	// Initialize a new Storage instance for the module.
	// The module-specific namespace or prefix can be used for module-specific data.
	moduleStorage := NewStorage(manager.client.Options(), key)

	return moduleStorage, nil
}
