package persistence

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// MigrationRunner handles database schema migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{
		db: db,
	}
}

// RunMigrations executes all pending migrations
func (m *MigrationRunner) RunMigrations(migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := m.getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get already applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Execute pending migrations
	for _, file := range migrationFiles {
		migrationName := filepath.Base(file)

		// Skip if already applied
		if contains(appliedMigrations, migrationName) {
			fmt.Printf("Migration %s already applied, skipping\n", migrationName)
			continue
		}

		fmt.Printf("Applying migration: %s\n", migrationName)

		if err := m.executeMigration(file, migrationName); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
		}

		fmt.Printf("Migration %s applied successfully\n", migrationName)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := m.db.Exec(query)
	return err
}

// getMigrationFiles returns sorted list of migration files
func (m *MigrationRunner) getMigrationFiles(migrationsPath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure consistent execution order
	sort.Strings(files)

	return files, nil
}

// getAppliedMigrations returns list of already applied migration names
func (m *MigrationRunner) getAppliedMigrations() ([]string, error) {
	query := "SELECT migration_name FROM schema_migrations ORDER BY applied_at"

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations = append(migrations, name)
	}

	return migrations, rows.Err()
}

// executeMigration executes a single migration file
func (m *MigrationRunner) executeMigration(filePath, migrationName string) error {
	// Read migration file
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute migration SQL
	if _, err = tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	insertQuery := "INSERT INTO schema_migrations (migration_name) VALUES (?)"
	if _, err = tx.Exec(insertQuery, migrationName); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

// RollbackLastMigration rolls back the most recent migration (if supported)
func (m *MigrationRunner) RollbackLastMigration() error {
	// Get the last applied migration
	query := "SELECT migration_name FROM schema_migrations ORDER BY applied_at DESC LIMIT 1"

	var lastMigration string
	if err := m.db.QueryRow(query).Scan(&lastMigration); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Note: SQLite doesn't support complex schema rollbacks easily
	// This is a placeholder for rollback functionality
	// In a production system, you'd want to maintain rollback scripts
	fmt.Printf("Warning: Rollback of migration %s not implemented\n", lastMigration)
	fmt.Printf("SQLite does not support complex schema rollbacks\n")
	fmt.Printf("Consider backing up database before migrations\n")

	return fmt.Errorf("rollback not implemented for SQLite")
}

// GetMigrationStatus returns the current migration status
func (m *MigrationRunner) GetMigrationStatus(migrationsPath string) ([]MigrationStatus, error) {
	// Get all migration files
	migrationFiles, err := m.getMigrationFiles(migrationsPath)
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var status []MigrationStatus

	for _, file := range migrationFiles {
		migrationName := filepath.Base(file)
		isApplied := contains(appliedMigrations, migrationName)

		status = append(status, MigrationStatus{
			Name:     migrationName,
			Applied:  isApplied,
			FilePath: file,
		})
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Name     string
	Applied  bool
	FilePath string
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// fsys is a placeholder for the filesystem interface
// In real implementation, this would be properly initialized
var fsys fs.FS