package rediscache

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/redis/go-redis/v9"
)

type ShallowStorage struct {
	client *redis.Client
}

// SetShallow stores a list of shallow commit hashes in Redis.
func (s *ShallowStorage) SetShallow(commits []plumbing.Hash) error {
	key := "git:shallow"

	// Convert plumbing.Hash slices to string slices for Redis.
	var hashes []string
	for _, hash := range commits {
		hashes = append(hashes, hash.String())
	}

	// Use SAdd for adding items to a set. This operation is idempotent.
	_, err := s.client.SAdd(ctx, key, hashes).Result()
	return err
}

// Shallow retrieves the list of shallow commit hashes from Redis.
func (s *ShallowStorage) Shallow() ([]plumbing.Hash, error) {
	key := "git:shallow"

	// Fetch all members of the set.
	members, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var commits []plumbing.Hash
	for _, member := range members {
		commits = append(commits, plumbing.NewHash(member))
	}

	return commits, nil
}
