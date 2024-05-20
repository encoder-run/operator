package postgrescache

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/encoder-run/operator/pkg/database"
	"github.com/go-git/go-git/v5/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigStorage struct {
	db              *gorm.DB
	namespacePrefix string
}

// SetConfig stores the configuration in Redis.
func (s *ConfigStorage) SetConfig(cfg *config.Config) error {
	fmt.Printf("%v", cfg)
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	config := &database.Config{
		URL:  s.namespacePrefix,
		Blob: data,
	}

	// Upsert operation: Update if exists, else create
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},                           // conflict target
		DoUpdates: clause.Assignments(map[string]interface{}{"blob": data}), // update just the blob
	}).Create(config).Error

}

// Config retrieves the configuration from Redis.
func (s *ConfigStorage) Config() (*config.Config, error) {
	dbConfig := &database.Config{}
	if err := s.db.Where("url = ?", s.namespacePrefix).First(dbConfig).Error; err != nil {
		// If not found return a new configuration.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return config.NewConfig(), nil
		}
	}

	var cfg config.Config
	if err := json.Unmarshal(dbConfig.Blob, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
