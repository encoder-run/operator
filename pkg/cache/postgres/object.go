package postgrescache

import (
	"errors"
	"io"

	"github.com/encoder-run/operator/pkg/database"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ObjectStorage struct {
	db              *gorm.DB
	namespacePrefix string
}

func (o *ObjectStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

// SetEncodedObject stores an encoded object in Postgres.
func (s *ObjectStorage) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	r, err := obj.Reader()
	if err != nil {
		return obj.Hash(), err
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		return obj.Hash(), err
	}

	// Create the object in the database.
	object := &database.Object{
		Hash: obj.Hash().String(),
		Type: obj.Type().String(),
		Size: int64(len(bytes)),
		Blob: bytes,
		URL:  s.namespacePrefix,
	}

	// Upsert operation: Update if exists, else create
	return obj.Hash(), s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hash"}, {Name: "url"}}, // conflict target
		DoUpdates: clause.Assignments(map[string]interface{}{
			"type": obj.Type().String(),
			"blob": bytes,
			"size": int64(len(bytes)),
		}), // update fields
	}).Create(object).Error
}

// EncodedObject retrieves an encoded object from the database by its hash, handling any type.
func (s *ObjectStorage) EncodedObject(t plumbing.ObjectType, h plumbing.Hash) (plumbing.EncodedObject, error) {
	var object database.Object
	var err error

	if t == plumbing.AnyObject {
		// Retrieve the first object that matches the hash across any object type
		err = s.db.Where("hash = ? AND url = ?", h.String(), s.namespacePrefix).First(&object).Error
	} else {
		// Specific object type requested
		err = s.db.Where("hash = ? AND type = ? AND url = ?", h.String(), t.String(), s.namespacePrefix).First(&object).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, plumbing.ErrObjectNotFound
		}
		return nil, err
	}

	o := &plumbing.MemoryObject{}
	oType, err := plumbing.ParseObjectType(object.Type)
	if err != nil {
		return nil, err
	}
	o.SetType(oType)
	o.SetSize(object.Size)
	_, err = o.Write(object.Blob)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// HasEncodedObject checks if an encoded object exists in the database by its hash.
func (s *ObjectStorage) HasEncodedObject(hash plumbing.Hash) error {
	// Using Count to check existence, more efficient than fetching the whole object
	var count int64
	err := s.db.Model(&database.Object{}).Where("hash = ? AND url = ?", hash.String(), s.namespacePrefix).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return plumbing.ErrObjectNotFound
	}
	return nil
}

type EncodedObjectIter struct {
	db              *gorm.DB
	namespacePrefix string
	objectType      plumbing.ObjectType
	offset          int
	limit           int
	objects         []*plumbing.MemoryObject
	moreData        bool
}

// IterEncodedObjects returns an iterator for encoded objects stored in Postgres.
func (s *ObjectStorage) IterEncodedObjects(t plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	return &EncodedObjectIter{
		db:              s.db,
		namespacePrefix: s.namespacePrefix,
		objectType:      t,
		offset:          0,
		limit:           100, // Define your batch size
		moreData:        true,
	}, nil
}

// fetchNextBatch fetches the next batch of encoded objects from the database.
func (iter *EncodedObjectIter) fetchNextBatch() error {
	var results []database.Object
	result := iter.db.Where("url = ? AND type = ?", iter.namespacePrefix, iter.objectType.String()).Offset(iter.offset).Limit(iter.limit).Find(&results)
	if result.Error != nil {
		iter.moreData = false
		return result.Error
	}
	iter.offset += len(results)
	iter.moreData = len(results) == iter.limit

	// Convert database objects to MemoryObject
	iter.objects = make([]*plumbing.MemoryObject, len(results))
	for i, dbObj := range results {
		o := plumbing.MemoryObject{}
		o.SetType(iter.objectType)
		o.SetSize(dbObj.Size)
		if _, err := o.Write(dbObj.Blob); err != nil {
			return err
		}
		iter.objects[i] = &o
	}
	return nil
}

// Next retrieves the next encoded object, advancing the iterator.
func (iter *EncodedObjectIter) Next() (plumbing.EncodedObject, error) {
	if len(iter.objects) == 0 && iter.moreData {
		if err := iter.fetchNextBatch(); err != nil {
			return nil, err
		}
	}
	if len(iter.objects) > 0 {
		obj := iter.objects[0]
		iter.objects = iter.objects[1:]
		return obj, nil
	}
	return nil, io.EOF
}

// ForEach implements the required method to iterate over each object.
func (iter *EncodedObjectIter) ForEach(cb func(obj plumbing.EncodedObject) error) error {
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

// Close is a placeholder to satisfy the EncodedObjectIter interface.
func (iter *EncodedObjectIter) Close() {
	// Nothing specific to close in this context, since DB connection is managed outside.
}

// EncodedObjectSize retrieves the size of an encoded object from Postgres by its hash.
func (s *ObjectStorage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	var size int64
	if err := s.db.Model(&database.Object{}).Select("size").Where("hash = ? AND url = ?", hash.String(), s.namespacePrefix).Scan(&size).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, plumbing.ErrObjectNotFound
		}
		return 0, err
	}

	return size, nil
}
