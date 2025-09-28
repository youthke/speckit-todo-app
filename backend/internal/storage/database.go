package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"todo-app/internal/models"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase() error {
	var err error

	// Use SQLite for development
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "todo.db"
	}

	// Configure GORM logger
	gormLogger := logger.Default
	if os.Getenv("ENV") == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto migrations
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// ResetDatabase drops all tables and recreates them (for testing)
func ResetDatabase() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Drop existing tables
	err := DB.Migrator().DropTable(&models.Task{})
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	// Recreate tables
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		return fmt.Errorf("failed to recreate tables: %w", err)
	}

	return nil
}