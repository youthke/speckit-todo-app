package services

import (
	"fmt"
	"log"
	"time"

	"todo-app/internal/models"
	"todo-app/internal/storage"
)

// HealthService provides health checking functionality
type HealthService struct {
	startTime time.Time
	version   string
}

// NewHealthService creates a new health service instance
func NewHealthService() *HealthService {
	return &HealthService{
		startTime: time.Now(),
		version:   "1.0.0", // This could be injected from build info
	}
}

// GetHealthStatus performs comprehensive health checks and returns the current status
func (hs *HealthService) GetHealthStatus() (*models.HealthResponse, error) {
	// Check database connectivity
	dbStatus := hs.checkDatabaseConnectivity()

	// Determine overall health based on database status
	overallHealth := models.DetermineOverallHealth(dbStatus)

	// Calculate uptime
	uptime := int64(time.Since(hs.startTime).Seconds())

	// Create health response
	response := models.NewHealthResponse(
		overallHealth,
		dbStatus,
		hs.version,
		uptime,
	)

	// Validate response before returning
	if err := response.Validate(); err != nil {
		log.Printf("Health response validation failed: %v", err)
		return nil, fmt.Errorf("health check validation failed: %w", err)
	}

	return response, nil
}

// checkDatabaseConnectivity tests the database connection and returns status
func (hs *HealthService) checkDatabaseConnectivity() models.DatabaseStatus {
	// Get the database instance
	db := storage.GetDB()
	if db == nil {
		log.Printf("Database instance is nil")
		return models.DatabaseStatusDisconnected
	}

	// Get underlying sql.DB to test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get underlying database connection: %v", err)
		return models.DatabaseStatusError
	}

	// Test connection with ping
	if err := sqlDB.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)
		return models.DatabaseStatusDisconnected
	}

	// Additional checks for database health
	if err := hs.performDatabaseHealthChecks(sqlDB); err != nil {
		log.Printf("Database health check failed: %v", err)
		return models.DatabaseStatusError
	}

	return models.DatabaseStatusConnected
}

// performDatabaseHealthChecks performs additional database health validations
func (hs *HealthService) performDatabaseHealthChecks(sqlDB interface{}) error {
	// For SQLite, we can check if we can perform a simple query
	// This ensures the database is not only connected but also responsive

	// Note: In a real implementation, you might want to:
	// - Check disk space for SQLite files
	// - Verify database schema integrity
	// - Test transaction capabilities
	// - Check for any locks or connection pool issues

	// For now, we'll do a minimal check since the Ping() already validates basic connectivity
	// and we don't want to impact performance significantly

	return nil
}

// GetDatabaseStatus returns just the database connectivity status
func (hs *HealthService) GetDatabaseStatus() models.DatabaseStatus {
	return hs.checkDatabaseConnectivity()
}

// IsHealthy returns whether the service is currently healthy
func (hs *HealthService) IsHealthy() bool {
	dbStatus := hs.checkDatabaseConnectivity()
	overallHealth := models.DetermineOverallHealth(dbStatus)
	return overallHealth == models.HealthStatusHealthy
}

// GetUptime returns the service uptime in seconds
func (hs *HealthService) GetUptime() int64 {
	return int64(time.Since(hs.startTime).Seconds())
}

// GetVersion returns the service version
func (hs *HealthService) GetVersion() string {
	return hs.version
}

// SetVersion allows updating the service version (useful for testing or build info injection)
func (hs *HealthService) SetVersion(version string) {
	hs.version = version
}

// ValidateHealthResponse validates a health response structure
func (hs *HealthService) ValidateHealthResponse(response *models.HealthResponse) error {
	if response == nil {
		return fmt.Errorf("health response cannot be nil")
	}
	return response.Validate()
}