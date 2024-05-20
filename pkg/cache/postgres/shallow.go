package postgrescache

import (
	"errors"

	"github.com/encoder-run/operator/pkg/database"
	"github.com/go-git/go-git/v5/plumbing"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShallowStorage struct {
	db              *gorm.DB
	namespacePrefix string
}

// SetShallow stores a list of shallow commit hashes in Postgres.
func (s *ShallowStorage) SetShallow(commits []plumbing.Hash) error {
	var hashes []string
	for _, hash := range commits {
		hashes = append(hashes, hash.String())
	}

	// Upsert operation to update or insert the shallow entry
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"hashes"}),
	}).Create(&database.Shallow{
		URL:    s.namespacePrefix,
		Hashes: hashes,
	}).Error
}

// Shallow retrieves the list of shallow commit hashes from Postgres.
func (s *ShallowStorage) Shallow() ([]plumbing.Hash, error) {
	var entry database.Shallow
	err := s.db.Where("url = ?", s.namespacePrefix).First(&entry).Error
	// Check if the error is because no records were found
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found, return an empty list of hashes
			return []plumbing.Hash{}, nil
		}
		// Some other error occurred, return the error
		return nil, err
	}

	// Convert the strings from the entry to plumbing.Hash
	var hashes []plumbing.Hash
	for _, h := range entry.Hashes {
		hashes = append(hashes, plumbing.NewHash(h))
	}

	return hashes, nil
}
