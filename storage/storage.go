package storage

import (
	"log"
	"order/configs"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	Order   OrderInterface
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

var (
	StorageInstance *Storage
	once            sync.Once
)

// Initialize initializes the storage singleton
func Initialize() {
	once.Do(func() {
		config := configs.GetEnvConfig()

		dsn := config.PostgresConnString

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		StorageInstance = &Storage{
			ReadDB:  db,
			WriteDB: db,
		}

		// Initialize order storage
		StorageInstance.Order = NewOrderTable(db, db)
	})
}

// GetInstance returns the storage singleton
func GetInstance() *Storage {
	if StorageInstance == nil {
		Initialize()
	}
	return StorageInstance
}

// AutoMigrate automatically migrates the database schema
func (s *Storage) AutoMigrate(models ...interface{}) error {
	return s.WriteDB.AutoMigrate(models...)
}

// Close closes the database connections
func (s *Storage) Close() error {
	sqlDB, err := s.ReadDB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	sqlDB, err = s.WriteDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
