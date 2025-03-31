package storage

import (
	"fmt"
	"log"
	"order/configs"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	Order OrderInterface
	write *gorm.DB
	read  *gorm.DB
}

var (
	StorageInstance *Storage
	once            sync.Once
)

// Initialize initializes the storage singleton
func (s *Storage) Initialize() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	envConfig := configs.GetEnvConfig()

	var err error
	s.write, err = gorm.Open(postgres.Open(envConfig.PostgresConnStringMaster), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB, err := s.write.DB()
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB.SetMaxOpenConns(envConfig.PostgresMaxOpenConns)
	writeDB.SetMaxIdleConns(envConfig.PostgresMaxIdleConns)

	s.read, err = gorm.Open(postgres.Open(envConfig.PostgresConnStringSlave), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("status: ", err)
	}
	readDB, err := s.read.DB()
	if err != nil {
		fmt.Println("status: ", err)
	}
	readDB.SetMaxOpenConns(envConfig.PostgresMaxOpenConns)
	readDB.SetMaxIdleConns(envConfig.PostgresMaxIdleConns)
}

// GetInstance returns the storage singleton
func GetInstance() *Storage {
	once.Do(func() {
		StorageInstance = &Storage{}
		StorageInstance.Initialize()
		StorageInstance.Order = NewOrderTable(StorageInstance.read, StorageInstance.write)
	})
	return StorageInstance
}

// AutoMigrate automatically migrates the database schema
func (s *Storage) AutoMigrate(model interface{}) {
	s.write.AutoMigrate(model)
	s.read.AutoMigrate(model)
}

// Close closes the database connections
func (s *Storage) Close() error {
	sqlDB, err := s.read.DB()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
		return err
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
		return err
	}

	sqlDB, err = s.write.DB()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
		return err
	}
	return sqlDB.Close()
}
