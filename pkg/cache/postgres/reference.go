package postgrescache

import (
	"errors"
	"io"

	"github.com/encoder-run/operator/pkg/database"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReferenceStorage struct {
	db              *gorm.DB
	namespacePrefix string
}

// SetReference stores a reference in Postgres.
func (r *ReferenceStorage) SetReference(ref *plumbing.Reference) error {
	reference := &database.Reference{
		Name:   ref.Name().String(),
		Type:   ref.Type().String(),
		Target: ref.Target().String(),
		Hash:   ref.Hash().String(),
		URL:    r.namespacePrefix,
	}

	// Upsert operation: Update if exists, else create
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}, {Name: "name"}},                          // conflict target
		DoUpdates: clause.Assignments(map[string]interface{}{"hash": ref.Hash().String()}), // update just the hash
	}).Create(reference).Error
}

// Reference retrieves a reference from Postgres by its name.
func (r *ReferenceStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	var reference database.Reference
	if err := r.db.Where("url = ? AND name = ?", r.namespacePrefix, name.String()).First(&reference).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Check if the error is a not-found error
			return nil, plumbing.ErrReferenceNotFound // Return the Git-specific not-found error
		}
		return nil, err // Return other types of errors directly
	}

	switch reference.Type {
	case plumbing.HashReference.String():
		return plumbing.NewHashReference(name, plumbing.NewHash(reference.Hash)), nil
	case plumbing.SymbolicReference.String():
		return plumbing.NewSymbolicReference(name, plumbing.ReferenceName(reference.Target)), nil
	default:
		return nil, errors.New("unknown reference type")
	}
}

// type ReferenceIter struct {
// 	refs  []*plumbing.Reference
// 	index int // Track the current index in the refs slice
// }

// // IterReferences initializes a ReferenceIter with all references loaded upfront.
// func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
// 	var results []database.Reference
// 	if err := r.db.Where("url = ?", r.namespacePrefix).Find(&results).Error; err != nil {
// 		return nil, err
// 	}

// 	// Convert database references to plumbing.Reference and store them in the iter
// 	refs := make([]*plumbing.Reference, len(results))
// 	for i, dbRef := range results {
// 		ref := plumbing.NewReferenceFromStrings(dbRef.Name, dbRef.Target)
// 		refs[i] = ref
// 	}

// 	return &ReferenceIter{
// 		refs:  refs,
// 		index: 0,
// 	}, nil
// }

// // Next retrieves the next reference, advancing the iterator.
// func (iter *ReferenceIter) Next() (*plumbing.Reference, error) {
// 	if iter.index >= len(iter.refs) {
// 		return nil, io.EOF
// 	}
// 	ref := iter.refs[iter.index]
// 	iter.index++
// 	return ref, nil
// }

// // ForEach implements the required method to iterate over each reference.
// func (iter *ReferenceIter) ForEach(cb func(obj *plumbing.Reference) error) error {
// 	for {
// 		ref, err := iter.Next()
// 		if err != nil {
// 			if err == io.EOF {
// 				return nil
// 			}
// 			return err
// 		}
// 		if err := cb(ref); err != nil {
// 			if err == storer.ErrStop {
// 				return nil
// 			}
// 			return err
// 		}
// 	}
// }

// // Close is a placeholder to satisfy the ReferenceIter interface.
// func (iter *ReferenceIter) Close() {
// 	// Nothing specific to close in this context, since DB connection is managed outside.
// }

type ReferenceIter struct {
	db              *gorm.DB
	namespacePrefix string
	offset          int
	limit           int
	refs            []*plumbing.Reference
	moreData        bool
}

// IterReferences returns an iterator for references stored in Postgres.
func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
	return &ReferenceIter{
		db:              r.db,
		namespacePrefix: r.namespacePrefix,
		offset:          0,
		limit:           100, // Define your batch size
		moreData:        true,
	}, nil
}

// fetchNextBatch fetches the next batch of references from the database.
func (iter *ReferenceIter) fetchNextBatch() error {
	var results []database.Reference
	result := iter.db.Where("url = ?", iter.namespacePrefix).Offset(iter.offset).Limit(iter.limit).Find(&results)
	if result.Error != nil {
		iter.moreData = false
		return result.Error
	}
	iter.offset += len(results)
	iter.moreData = len(results) == iter.limit

	// Convert database references to plumbing.Reference
	iter.refs = make([]*plumbing.Reference, len(results))
	for i, dbRef := range results {
		var ref *plumbing.Reference
		switch dbRef.Type {
		case plumbing.HashReference.String():
			ref = plumbing.NewHashReference(plumbing.ReferenceName(dbRef.Name), plumbing.NewHash(dbRef.Target))
		case plumbing.SymbolicReference.String():
			ref = plumbing.NewSymbolicReference(plumbing.ReferenceName(dbRef.Name), plumbing.ReferenceName(dbRef.Target))
		default:
			return errors.New("unknown reference type")
		}
		iter.refs[i] = ref
	}
	return nil
}

// Next retrieves the next reference, advancing the iterator.
func (iter *ReferenceIter) Next() (*plumbing.Reference, error) {
	if len(iter.refs) == 0 && iter.moreData {
		if err := iter.fetchNextBatch(); err != nil {
			return nil, err
		}
	}
	if len(iter.refs) > 0 {
		ref := iter.refs[0]
		iter.refs = iter.refs[1:]
		return ref, nil
	}
	return nil, io.EOF
}

// ForEach implements the required method to iterate over each reference.
func (iter *ReferenceIter) ForEach(cb func(obj *plumbing.Reference) error) error {
	for {
		ref, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := cb(ref); err != nil {
			if err == storer.ErrStop {
				return nil
			}
			return err
		}
	}
}

// Close is a placeholder to satisfy the ReferenceIter interface.
func (iter *ReferenceIter) Close() {
	// Nothing specific to close in this context, since DB connection is managed outside.
}

// CheckAndSetReference atomically checks if the existing reference has changed.
func (r *ReferenceStorage) CheckAndSetReference(new, old *plumbing.Reference) error {
	if new == nil {
		return errors.New("new reference cannot be nil")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Set the transaction isolation level to Serializable
		if err := tx.Clauses(clause.Locking{Strength: "SERIALIZABLE"}).Error; err != nil {
			return err
		}

		var current database.Reference
		err := tx.Where("url = ? AND name = ?", r.namespacePrefix, new.Name().String()).First(&current).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Handle not found error
		if old != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return plumbing.ErrReferenceNotFound
			}

			currentRef := plumbing.NewReferenceFromStrings(current.Name, current.Target)
			if currentRef.Hash() != old.Hash() {
				return storage.ErrReferenceHasChanged // Hash mismatch indicates concurrent update.
			}
		}

		// If the reference was found or no old reference was provided, proceed to set the new reference.
		return r.setReference(tx, new)
	})
}

// setReference is a helper function to upsert a reference within a transaction
func (r *ReferenceStorage) setReference(tx *gorm.DB, ref *plumbing.Reference) error {
	reference := &database.Reference{
		Name:   ref.Name().String(),
		Type:   ref.Type().String(),
		Target: ref.Target().String(),
		Hash:   ref.Hash().String(),
		URL:    r.namespacePrefix,
	}

	// Upsert operation: Update if exists, else create
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}, {Name: "name"}},                                                                                        // conflict target
		DoUpdates: clause.Assignments(map[string]interface{}{"hash": ref.Hash().String(), "type": ref.Type().String(), "target": ref.Target().String()}), // update fields
	}).Create(reference).Error
}
