package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mu        sync.Mutex
	instances = make(map[string]*gorm.DB)
)

// GetPostgresClient initializes or returns a database connection for a given DSN
func GetPostgresClient(dsn string) (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if an instance already exists for this DSN
	if db, exists := instances[dsn]; exists {
		return db, nil
	}

	// Create a new instance since one doesn't exist
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set log mode to `Info` to log all SQL queries
	})
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return nil, err
	}

	// Execute initial setup queries
	if result := db.Exec("CREATE EXTENSION IF NOT EXISTS vector"); result.Error != nil {
		log.Printf("Failed to execute initial setup query: %v", result.Error)
		return nil, result.Error
	}

	// Auto-migrate the schema
	db.AutoMigrate(&Object{}, &Reference{}, &Config{}, &Shallow{}, &Index{}, &CodeEmbedding{})

	// Store the instance in the map
	instances[dsn] = db

	return db, nil
}

// ConstructPostgresDSN constructs the Data Source Name for a PostgreSQL connection
func ConstructPostgresDSN(host, user, password, dbname, port, sslmode, timezone string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)
}

type Object struct {
	Hash string `gorm:"primaryKey;type:varchar(255)"`
	Type string `gorm:"type:varchar(255)"`
	Blob []byte
	Size int64  `gorm:"type:bigint"`
	URL  string `gorm:"primaryKey;type:varchar(255)"`
}

type Reference struct {
	Name   string `gorm:"primaryKey;type:varchar(255)"`
	Type   string `gorm:"type:varchar(255)"`
	Target string `gorm:"type:varchar(255)"`
	Hash   string `gorm:"type:varchar(255)"`
	URL    string `gorm:"primaryKey;type:varchar(255)"`
}

type Config struct {
	URL  string `gorm:"primaryKey;type:varchar(255)"`
	Blob []byte
}

type Shallow struct {
	URL    string         `gorm:"primaryKey;type:varchar(255)"`
	Hashes pq.StringArray `gorm:"type:text[]"`
}

type Index struct {
	URL  string `gorm:"primaryKey;type:varchar(255)"`
	Blob []byte
}

type CodeEmbedding struct {
	URL        string `gorm:"primaryKey;type:varchar(255)"`
	FileHash   string `gorm:"primaryKey;type:varchar(255)"`
	FilePath   string `gorm:"primaryKey;type:varchar(255)"`
	ChunkID    int    `gorm:"primaryKey"`
	StartIndex int
	EndIndex   int
	Embedding  pgvector.Vector `gorm:"type:vector(768)"`
}
