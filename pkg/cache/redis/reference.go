package rediscache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"github.com/redis/go-redis/v9"
)

const (
	referencePrefix = "ref"
)

type ReferenceStorage struct {
	client          *redis.Client
	namespacePrefix string
}

// SetReference stores a reference in Redis.
func (r *ReferenceStorage) SetReference(ref *plumbing.Reference) error {
	refMap, err := referenceToMap(ref)
	if err != nil {
		return err
	}

	data, err := json.Marshal(refMap)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, withNamespace(r.namespacePrefix, referencePrefix, ref.Name().String()), data, 0).Err()
}

// Reference retrieves a reference from Redis by its name.
func (r *ReferenceStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	data, err := r.client.Get(ctx, withNamespace(r.namespacePrefix, referencePrefix, name.String())).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, plumbing.ErrReferenceNotFound
		}
		return nil, err
	}

	return jsonToReference(data)
}

func (r *ReferenceStorage) CheckAndSetReference(new, old *plumbing.Reference) error {
	if new == nil {
		return errors.New("new reference cannot be nil")
	}

	key := withNamespace(r.namespacePrefix, referencePrefix, new.Name().String())

	// The transaction function. It will be retried if the watched keys change.
	txnFunc := func(tx *redis.Tx) error {
		// Get the current reference data.
		data, err := tx.Get(ctx, key).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}

		// If old reference is provided, ensure the current reference matches.
		if old != nil {
			if err == redis.Nil {
				return storage.ErrReferenceHasChanged // Reference was expected but not found.
			}

			currentRef, err := jsonToReference(data)
			if err != nil {
				return err
			}

			if currentRef.Hash() != old.Hash() {
				return storage.ErrReferenceHasChanged // Hash mismatch indicates concurrent update.
			}
		}

		// If the reference was found or no old reference was provided, proceed to set the new reference.
		newRefMap, err := referenceToMap(new)
		if err != nil {
			return err
		}

		newData, err := json.Marshal(newRefMap)
		if err != nil {
			return err
		}

		// Set the new reference data inside the transaction.
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, newData, 0)
			return nil
		})

		return err
	}

	// Execute the transaction.
	return r.client.Watch(ctx, txnFunc, key)
}

// PackRefs implements storage.Storer.
func (s *ReferenceStorage) PackRefs() error {
	return nil
}

// CountLooseRefs counts the number of references stored in Redis.
func (s *ReferenceStorage) CountLooseRefs() (int, error) {
	var count int
	var cursor uint64
	var err error

	// Use the SCAN command to iterate through all keys that match the pattern "ref:*".
	// We're only interested in counting keys, so we don't need to retrieve their values.
	for {
		var keys []string
		keys, cursor, err = s.client.Scan(ctx, cursor, withNamespace(s.namespacePrefix, referencePrefix, "*"), 0).Result()
		if err != nil {
			return 0, err // Return the error if the SCAN command fails.
		}

		count += len(keys) // Increment the count by the number of keys returned in this batch.

		if cursor == 0 {
			break // If the cursor returned by SCAN is 0, we've iterated through all keys.
		}
	}

	return count, nil
}

// RemoveReference implements storage.Storer.
func (s *ReferenceStorage) RemoveReference(ref plumbing.ReferenceName) error {
	fmt.Printf("RemoveReference")
	panic("unimplemented")
}

type ReferenceIter struct {
	client          *redis.Client
	namespacePrefix string
	cursor          uint64
	keys            []string
	moreData        bool
}

// IterReferences returns an iterator for references stored in Redis.
func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
	// Initialize the iterator
	iter := &ReferenceIter{
		client:          r.client,
		namespacePrefix: r.namespacePrefix,
		cursor:          0,
		keys:            nil,
		moreData:        true, // Assume there's more data until proven otherwise
	}
	// Perform the initial scan
	iter.fetchNextBatch()
	return iter, nil
}

// fetchNextBatch fetches the next batch of keys from Redis matching the pattern "ref:*".
func (iter *ReferenceIter) fetchNextBatch() {
	if iter.moreData {
		var err error
		iter.keys, iter.cursor, err = iter.client.Scan(ctx, iter.cursor, withNamespace(iter.namespacePrefix, referencePrefix, "*"), 100).Result()
		if err != nil {
			iter.moreData = false // In case of error, stop further fetching
			return
		}
		// Update moreData based on whether the cursor returned to 0
		iter.moreData = (iter.cursor != 0)
	}
}

// Next retrieves the next reference, advancing the iterator.
func (iter *ReferenceIter) Next() (*plumbing.Reference, error) {
	if len(iter.keys) == 0 && iter.moreData {
		iter.fetchNextBatch()
	}
	if len(iter.keys) > 0 {
		// Assume key to be the key itself and value to be something sensible
		key := iter.keys[0]
		iter.keys = iter.keys[1:]

		data, err := iter.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		ref, err := jsonToReference(data)
		if err != nil {
			return nil, err
		}

		return ref, nil
	}
	return nil, io.EOF
}

func (iter *ReferenceIter) ForEach(cb func(obj *plumbing.Reference) error) error {
	for {
		obj, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err := cb(obj); err != nil {
			if err == storer.ErrStop {
				return nil
			}

			return err
		}
	}
}

// Close is a placeholder to satisfy the ReferenceIter interface.
// Redis connections are managed outside the iterator, so there's nothing to close here.
func (iter *ReferenceIter) Close() {}

func referenceToMap(ref *plumbing.Reference) (map[string]interface{}, error) {
	if ref == nil {
		return nil, fmt.Errorf("nil reference")
	}
	refMap := map[string]interface{}{
		"Name":   ref.Name().String(),
		"Type":   ref.Type().String(),
		"Target": ref.Target().String(),
		// Include other fields as necessary.
	}
	if ref.Type() == plumbing.HashReference {
		refMap["Hash"] = ref.Hash().String()
	}
	return refMap, nil
}

func jsonToReference(data []byte) (*plumbing.Reference, error) {
	var refMap map[string]interface{}
	if err := json.Unmarshal(data, &refMap); err != nil {
		return nil, err
	}

	name, ok := refMap["Name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reference name")
	}

	target, ok := refMap["Target"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reference target")
	}

	refType, ok := refMap["Type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reference type")
	}

	hash, ok := refMap["Hash"].(string)
	if refType == plumbing.HashReference.String() && !ok {
		return nil, fmt.Errorf("invalid reference hash")
	}

	var ref plumbing.Reference
	switch refType {
	case plumbing.SymbolicReference.String():
		ref = *plumbing.NewSymbolicReference(plumbing.ReferenceName(name), plumbing.ReferenceName(target))
	case plumbing.HashReference.String():
		ref = *plumbing.NewHashReference(plumbing.ReferenceName(name), plumbing.NewHash(hash))
	default:
		return nil, fmt.Errorf("unsupported reference type")
	}

	return &ref, nil
}
