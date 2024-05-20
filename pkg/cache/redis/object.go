package rediscache

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/redis/go-redis/v9"
)

const (
	objectPrefix = "object"
)

type ObjectStorage struct {
	client          *redis.Client
	namespacePrefix string
}

func (o *ObjectStorage) ObjectPacks() ([]plumbing.Hash, error) {
	fmt.Printf("ObjectPacks")
	return nil, nil
}
func (o *ObjectStorage) DeleteOldObjectPackAndIndex(plumbing.Hash, time.Time) error {
	print("DeleteOldObjectPackAndIndex")
	return nil
}

func (o *ObjectStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

// SetEncodedObject stores an encoded object in Redis.
func (s *ObjectStorage) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	r, err := obj.Reader()
	if err != nil {
		return obj.Hash(), err
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		return obj.Hash(), err
	}

	key := fmt.Sprintf("%s:%s", obj.Type(), obj.Hash())
	if err := s.client.Set(ctx, withNamespace(s.namespacePrefix, objectPrefix, key), bytes, 0).Err(); err != nil {
		return plumbing.ZeroHash, err
	}

	return obj.Hash(), nil
}

// EncodedObject retrieves an encoded object from Redis by its hash.
func (s *ObjectStorage) EncodedObject(t plumbing.ObjectType, h plumbing.Hash) (plumbing.EncodedObject, error) {
	// if the object type is AnyObject, we need to try and get the object for each type until we find it.
	if t == plumbing.AnyObject {
		for _, ot := range []plumbing.ObjectType{plumbing.CommitObject, plumbing.TreeObject, plumbing.BlobObject, plumbing.TagObject, plumbing.OFSDeltaObject, plumbing.REFDeltaObject} {
			obj, err := s.EncodedObject(ot, h)
			if err == nil {
				return obj, nil
			}
		}
		return nil, plumbing.ErrObjectNotFound
	}

	// Get object data by hash as key with prefix "object:", type and hash.
	key := fmt.Sprintf("%s:%s", t, h)
	data, err := s.client.Get(ctx, withNamespace(s.namespacePrefix, objectPrefix, key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, plumbing.ErrObjectNotFound
		}
		return nil, err
	}

	o := &plumbing.MemoryObject{}
	o.SetType(t)
	o.SetSize(int64(len(data)))
	_, err = o.Write(data)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// EncodedObjectSize implements storage.Storer.
func (s *Storage) EncodedObjectSize(h plumbing.Hash) (int64, error) {
	pattern := withNamespace(s.namespacePrefix, objectPrefix, fmt.Sprintf("*:%s", h.String()))
	var cursor uint64
	var err error
	for {
		// Scan for keys that match the hash pattern
		var keys []string
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 1).Result()
		if err != nil {
			return 0, err
		}
		if len(keys) > 0 {
			// Use STRLEN to get the size of the object stored at the found key
			size, err := s.client.StrLen(ctx, keys[0]).Result()
			if err != nil {
				return 0, err
			}
			return size, nil
		}
		if cursor == 0 { // No more keys to scan
			break
		}
	}
	return 0, plumbing.ErrObjectNotFound
}

func (s *Storage) HasEncodedObject(h plumbing.Hash) error {
	pattern := withNamespace(s.namespacePrefix, objectPrefix, fmt.Sprintf("*:%s", h.String()))
	var cursor uint64
	var err error
	for {
		var keys []string
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 1).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			return nil // Key exists
		}
		if cursor == 0 { // No more keys to scan
			break
		}
	}
	return plumbing.ErrObjectNotFound
}

type EncodedObjectIter struct {
	client          *redis.Client
	namespacePrefix string
	t               plumbing.ObjectType
	cursor          uint64
	keys            []string
	moreData        bool
}

// IterEncodedObjects returns an iterator for encoded objects stored in Redis.
func (s *ObjectStorage) IterEncodedObjects(t plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	iter := &EncodedObjectIter{
		client:          s.client,
		namespacePrefix: s.namespacePrefix,
		t:               t,
		cursor:          0,
		moreData:        true, // Initially assume there is data to fetch
	}
	iter.fetchNextBatch() // Initial fetch
	return iter, nil
}

// fetchNextBatch fetches the next batch of keys from Redis matching the object type pattern.
func (iter *EncodedObjectIter) fetchNextBatch() {
	if iter.moreData {
		pattern := withNamespace(iter.namespacePrefix, objectPrefix, fmt.Sprintf("%s:*", iter.t))
		var err error
		iter.keys, iter.cursor, err = iter.client.Scan(ctx, iter.cursor, pattern, 100).Result()
		if err != nil {
			iter.moreData = false // Stop further fetches on error
			return
		}
		// Update moreData based on whether the cursor returned to 0
		iter.moreData = (iter.cursor != 0)
	}
}

// Next moves the iterator to the next object and returns it.
func (iter *EncodedObjectIter) Next() (plumbing.EncodedObject, error) {
	if len(iter.keys) == 0 && iter.moreData {
		iter.fetchNextBatch()
	}
	if len(iter.keys) > 0 {
		key := iter.keys[0]
		iter.keys = iter.keys[1:]

		// Fetch the data corresponding to the key.
		data, err := iter.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		obj := &plumbing.MemoryObject{}
		obj.SetType(iter.t)
		obj.SetSize(int64(len(data)))
		if _, err := obj.Write(data); err != nil {
			return nil, err
		}

		return obj, nil
	}

	return nil, io.EOF
}

// ForEach iterates over all objects and applies a callback function to each.
func (iter *EncodedObjectIter) ForEach(cb func(plumbing.EncodedObject) error) error {
	for {
		obj, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		err = cb(obj)
		if err != nil {
			if err == storer.ErrStop {
				return nil
			}
			return err
		}
	}
}

// Close is a placeholder to satisfy the EncodedObjectIter interface.
// Redis connections are managed outside the iterator, so there's nothing to close here.
func (iter *EncodedObjectIter) Close() {}

func (s *ObjectStorage) ForEachObjectHash(fun func(plumbing.Hash) error) error {
	// Scan for keys that match the object pattern

	pattern := withNamespace(s.namespacePrefix, objectPrefix, "*")
	var cursor uint64
	var err error
	for {
		var keys []string
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			// Extract the hash from the key by splitting on ":"
			hash := plumbing.NewHash(key[strings.LastIndex(key, ":")+1:])
			if err := fun(hash); err != nil {
				return err
			}
		}
		if cursor == 0 { // No more keys to scan
			break
		}
	}
	return nil
}

// LooseObjectTime looks up the (m)time associated with the
// loose object (that is not in a pack file). Some
// implementations (e.g. without loose objects)
// always return an error.
func (s *ObjectStorage) LooseObjectTime(plumbing.Hash) (time.Time, error) {
	return time.Time{}, fmt.Errorf("not implemented")
}

// DeleteLooseObject deletes a loose object if it exists.
func (s *ObjectStorage) DeleteLooseObject(plumbing.Hash) error {
	fmt.Printf("DeleteLooseObject")
	return fmt.Errorf("not implemented")

}
