package postgrescache

import (
	"encoding/json"
	"errors"

	"github.com/encoder-run/operator/pkg/database"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IndexStorage struct {
	db              *gorm.DB
	namespacePrefix string
}

// SetIndex serializes and stores the index in Postgres.
func (s *IndexStorage) SetIndex(idx *index.Index) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	index := &database.Index{
		URL:  s.namespacePrefix,
		Blob: data,
	}

	// Upsert operation: Update if exists, else create
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},                           // conflict target
		DoUpdates: clause.Assignments(map[string]interface{}{"blob": data}), // update just the blob
	}).Create(index).Error
}

// Index retrieves and deserializes the index from Postgres.
func (s *IndexStorage) Index() (*index.Index, error) {
	dbIndex := &database.Index{}
	if err := s.db.Where("url = ?", s.namespacePrefix).First(dbIndex).Error; err != nil {
		// If not found return a new index.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &index.Index{Version: 2}, nil
		}
	}

	var idx index.Index
	if err := json.Unmarshal(dbIndex.Blob, &idx); err != nil {
		return nil, err
	}

	return &idx, nil
}
