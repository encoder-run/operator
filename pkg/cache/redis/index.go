package rediscache

import (
	"encoding/json"

	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/redis/go-redis/v9"
)

type IndexStorage struct {
	client *redis.Client
}

// SetIndex serializes and stores the index in Redis.
func (s *IndexStorage) SetIndex(idx *index.Index) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	// Assuming "index" is the key where the index is stored.
	// You might want to use a more specific key based on your application's needs.
	return s.client.Set(ctx, "index", data, 0).Err()
}

// Index retrieves and deserializes the index from Redis.
func (s *IndexStorage) Index() (*index.Index, error) {
	data, err := s.client.Get(ctx, "index").Bytes()
	if err != nil {
		if err == redis.Nil {
			return &index.Index{Version: 2}, nil
		}
	}

	var idx index.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}

	return &idx, nil
}
