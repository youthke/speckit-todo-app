package config

import (
	"fmt"
	"os"
	"time"

	"domain/auth/entities"
	"domain/auth/valueobjects"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"todo-app/internal/models"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel
}

// GetDefaultDatabaseConfig returns default database configuration
func GetDefaultDatabaseConfig() DatabaseConfig {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "./todo.db" // Default SQLite database
	}

	logLevel := logger.Info
	if os.Getenv("APP_ENV") == "production" {
		logLevel = logger.Error
	}

	return DatabaseConfig{
		DSN:             dsn,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogLevel:        logLevel,
	}
}

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(config DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database for connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&valueobjects.GoogleIdentity{},
		&entities.AuthenticationSession{},
		&entities.OAuthState{},
	)
}

// InitializeDatabase initializes database connection and runs migrations
func InitializeDatabase() (*gorm.DB, error) {
	config := GetDefaultDatabaseConfig()
	db, err := NewDatabaseConnection(config)
	if err != nil {
		return nil, err
	}

	// Run auto migrations
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// HealthCheck performs a database health check
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// CloseDatabase gracefully closes database connection
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}