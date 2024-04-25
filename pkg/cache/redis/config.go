package rediscache

import (
	"encoding/json"

	"github.com/go-git/go-git/v5/config"
	"github.com/redis/go-redis/v9"
)

const (
	configPrefix = "config"
)

type ConfigStorage struct {
	client          *redis.Client
	namespacePrefix string
}

// SetConfig stores the configuration in Redis.
func (s *ConfigStorage) SetConfig(cfg *config.Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, withNamespace(s.namespacePrefix, configPrefix, ""), data, 0).Err()
}

// Config retrieves the configuration from Redis.
func (s *ConfigStorage) Config() (*config.Config, error) {
	data, err := s.client.Get(ctx, withNamespace(s.namespacePrefix, configPrefix, "")).Bytes()
	if err != nil {
		// If the key does not exist, return a new configuration.
		if err == redis.Nil {
			return config.NewConfig(), nil
		}
	}

	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
