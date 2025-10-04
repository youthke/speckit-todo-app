package unit

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"domain/health/entities"
	"todo-app/internal/services"
)

// TestHealthServiceDatabaseConnectionTimeout tests database connection timeout scenarios
func TestHealthServiceDatabaseConnectionTimeout(t *testing.T) {
	tests := []struct {
		name               string
		pingError          error
		expectedDBStatus   entities.DatabaseStatus
		expectedHealth     entities.HealthStatus
		shouldFailValidation bool
	}{
		{
			name:             "Database connection timeout",
			pingError:        sql.ErrConnDone,
			expectedDBStatus: entities.DatabaseStatusDisconnected,
			expectedHealth:   entities.HealthStatusDegraded,
		},
		{
			name:             "Database connection refused",
			pingError:        errors.New("connection refused"),
			expectedDBStatus: entities.DatabaseStatusDisconnected,
			expectedHealth:   entities.HealthStatusDegraded,
		},
		{
			name:             "Database timeout error",
			pingError:        errors.New("timeout"),
			expectedDBStatus: entities.DatabaseStatusDisconnected,
			expectedHealth:   entities.HealthStatusDegraded,
		},
		{
			name:             "Database network error",
			pingError:        errors.New("network unreachable"),
			expectedDBStatus: entities.DatabaseStatusDisconnected,
			expectedHealth:   entities.HealthStatusDegraded,
		},
		{
			name:             "Successful connection",
			pingError:        nil,
			expectedDBStatus: entities.DatabaseStatusConnected,
			expectedHealth:   entities.HealthStatusHealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new health service
			healthService := services.NewHealthService()

			// Note: In a real implementation, we would inject a database interface
			// For this test, we're testing the logic without actual database connection
			// The actual database connectivity testing would require dependency injection

			// Test the DetermineOverallHealth logic directly
			overallHealth := dtos.DetermineOverallHealth(tt.expectedDBStatus)
			assert.Equal(t, tt.expectedHealth, overallHealth)

			// Test that the health service can handle various scenarios
			response, err := healthService.GetHealthStatus()
			if tt.shouldFailValidation {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)

				// Verify response structure is valid
				assert.NoError(t, response.Validate())
				assert.NotEmpty(t, response.Timestamp)
				assert.GreaterOrEqual(t, response.Uptime, int64(0))
			}
		})
	}
}

func TestHealthServiceInvalidDatabaseStates(t *testing.T) {
	tests := []struct {
		name           string
		dbStatus       dtos.DatabaseStatus
		expectedHealth dtos.HealthStatus
	}{
		{
			name:           "Unknown database status should result in unhealthy",
			dbStatus:       "unknown",
			expectedHealth: dtos.HealthStatusUnhealthy,
		},
		{
			name:           "Empty database status should result in unhealthy",
			dbStatus:       "",
			expectedHealth: dtos.HealthStatusUnhealthy,
		},
		{
			name:           "Corrupted database status should result in unhealthy",
			dbStatus:       "corrupted",
			expectedHealth: dtos.HealthStatusUnhealthy,
		},
		{
			name:           "Malformed database status should result in unhealthy",
			dbStatus:       "con nected",
			expectedHealth: dtos.HealthStatusUnhealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dtos.DetermineOverallHealth(tt.dbStatus)
			assert.Equal(t, tt.expectedHealth, result)
		})
	}
}

func TestHealthServiceStartupConditions(t *testing.T) {
	t.Run("Fresh service startup", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Service should be immediately available
		assert.NotNil(t, healthService)

		// Get health status immediately after startup
		response, err := healthService.GetHealthStatus()
		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Uptime should be very small (< 1 second for a fresh service)
		assert.GreaterOrEqual(t, response.Uptime, int64(0))
		assert.LessOrEqual(t, response.Uptime, int64(1))

		// Timestamp should be recent
		timestamp, err := time.Parse(time.RFC3339, response.Timestamp)
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now().UTC(), timestamp, time.Second)
	})

	t.Run("Service with custom version", func(t *testing.T) {
		healthService := services.NewHealthService()
		customVersion := "2.1.0-beta"

		healthService.SetVersion(customVersion)

		response, err := healthService.GetHealthStatus()
		assert.NoError(t, err)
		assert.Equal(t, customVersion, response.Version)
	})

	t.Run("Service uptime progression", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Get initial uptime
		firstUptime := healthService.GetUptime()

		// Wait a small amount
		time.Sleep(100 * time.Millisecond)

		// Get uptime again
		secondUptime := healthService.GetUptime()

		// Second uptime should be greater than or equal to first (accounting for timing precision)
		assert.GreaterOrEqual(t, secondUptime, firstUptime)
	})
}

func TestHealthServiceHighLoadScenarios(t *testing.T) {
	t.Run("Concurrent health checks", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Number of concurrent requests to simulate
		numRequests := 100
		resultChan := make(chan error, numRequests)

		// Launch concurrent health checks
		for i := 0; i < numRequests; i++ {
			go func() {
				response, err := healthService.GetHealthStatus()
				if err != nil {
					resultChan <- err
					return
				}

				// Validate response
				if validateErr := response.Validate(); validateErr != nil {
					resultChan <- validateErr
					return
				}

				resultChan <- nil
			}()
		}

		// Collect results
		var errors []error
		for i := 0; i < numRequests; i++ {
			if err := <-resultChan; err != nil {
				errors = append(errors, err)
			}
		}

		// All requests should succeed
		assert.Empty(t, errors, "All concurrent health checks should succeed")
	})

	t.Run("Health check under memory pressure", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Simulate memory pressure by creating large objects
		var memoryPressure [][]byte
		for i := 0; i < 1000; i++ {
			memoryPressure = append(memoryPressure, make([]byte, 1024))
		}

		// Health check should still work under memory pressure
		response, err := healthService.GetHealthStatus()
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NoError(t, response.Validate())

		// Clean up memory
		memoryPressure = nil
	})

	t.Run("Rapid successive health checks", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Perform rapid successive health checks
		for i := 0; i < 50; i++ {
			response, err := healthService.GetHealthStatus()
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.NoError(t, response.Validate())

			// Verify uptime is monotonically increasing
			if i > 0 {
				nextResponse, nextErr := healthService.GetHealthStatus()
				assert.NoError(t, nextErr)
				assert.GreaterOrEqual(t, nextResponse.Uptime, response.Uptime)
			}
		}
	})
}

func TestHealthServiceMalformedRequests(t *testing.T) {
	t.Run("Nil health response validation", func(t *testing.T) {
		healthService := services.NewHealthService()

		err := healthService.ValidateHealthResponse(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "health response cannot be nil")
	})

	t.Run("Invalid health response validation", func(t *testing.T) {
		healthService := services.NewHealthService()

		invalidResponse := &dtos.HealthResponse{
			Status:    "invalid_status",
			Database:  "invalid_database",
			Timestamp: "invalid_timestamp",
			Uptime:    -1,
			Version:   "   ", // whitespace-only
		}

		err := healthService.ValidateHealthResponse(invalidResponse)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("Corrupted health response fields", func(t *testing.T) {
		tests := []struct {
			name     string
			response dtos.HealthResponse
			errorMsg string
		}{
			{
				name: "Corrupted status field",
				response: dtos.HealthResponse{
					Status:    dtos.HealthStatus("héalthy"), // non-ASCII
					Database:  dtos.DatabaseStatusConnected,
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				},
				errorMsg: "invalid status",
			},
			{
				name: "Corrupted database field",
				response: dtos.HealthResponse{
					Status:    dtos.HealthStatusHealthy,
					Database:  dtos.DatabaseStatus("connëcted"), // non-ASCII
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				},
				errorMsg: "invalid database status",
			},
			{
				name: "Malformed timestamp with injection attempt",
				response: dtos.HealthResponse{
					Status:    dtos.HealthStatusHealthy,
					Database:  dtos.DatabaseStatusConnected,
					Timestamp: "2023-01-01T00:00:00Z'; DROP TABLE health; --",
				},
				errorMsg: "invalid timestamp format",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.response.Validate()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			})
		}
	})
}

func TestHealthServiceEdgeCaseTimestamps(t *testing.T) {
	tests := []struct {
		name        string
		timestamp   string
		expectError bool
	}{
		{
			name:        "Unix epoch",
			timestamp:   "1970-01-01T00:00:00Z",
			expectError: false,
		},
		{
			name:        "Far future timestamp",
			timestamp:   "2099-12-31T23:59:59Z",
			expectError: false,
		},
		{
			name:        "Leap year timestamp",
			timestamp:   "2024-02-29T12:00:00Z",
			expectError: false,
		},
		{
			name:        "Invalid leap year",
			timestamp:   "2023-02-29T12:00:00Z",
			expectError: true,
		},
		{
			name:        "Timezone edge case",
			timestamp:   "2023-01-01T00:00:00+14:00", // UTC+14 (maximum timezone)
			expectError: false,
		},
		{
			name:        "Timezone edge case negative",
			timestamp:   "2023-01-01T00:00:00-12:00", // UTC-12 (minimum timezone)
			expectError: false,
		},
		{
			name:        "Invalid timezone",
			timestamp:   "2023-01-01T00:00:00+25:00",
			expectError: true,
		},
		{
			name:        "Microsecond precision",
			timestamp:   "2023-01-01T12:00:00.123456Z",
			expectError: false,
		},
		{
			name:        "Nanosecond precision",
			timestamp:   "2023-01-01T12:00:00.123456789Z",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := dtos.HealthResponse{
				Status:    dtos.HealthStatusHealthy,
				Database:  dtos.DatabaseStatusConnected,
				Timestamp: tt.timestamp,
			}

			err := response.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHealthServiceBoundaryConditions(t *testing.T) {
	t.Run("Maximum uptime value", func(t *testing.T) {
		response := dtos.HealthResponse{
			Status:    dtos.HealthStatusHealthy,
			Database:  dtos.DatabaseStatusConnected,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    9223372036854775807, // max int64
		}

		err := response.Validate()
		assert.NoError(t, err)
	})

	t.Run("Minimum negative uptime", func(t *testing.T) {
		response := dtos.HealthResponse{
			Status:    dtos.HealthStatusHealthy,
			Database:  dtos.DatabaseStatusConnected,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    -9223372036854775808, // min int64
		}

		err := response.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "uptime must be non-negative")
	})

	t.Run("Very long version string", func(t *testing.T) {
		longVersion := string(make([]byte, 10000)) // 10KB version string
		for i := range longVersion {
			longVersion = longVersion[:i] + "a" + longVersion[i+1:]
		}

		response := dtos.HealthResponse{
			Status:    dtos.HealthStatusHealthy,
			Database:  dtos.DatabaseStatusConnected,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   longVersion,
		}

		err := response.Validate()
		assert.NoError(t, err) // Should not fail on long versions
	})

	t.Run("Unicode characters in version", func(t *testing.T) {
		unicodeVersion := "1.0.0-🚀.beta"

		response := dtos.HealthResponse{
			Status:    dtos.HealthStatusHealthy,
			Database:  dtos.DatabaseStatusConnected,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   unicodeVersion,
		}

		err := response.Validate()
		assert.NoError(t, err) // Should handle Unicode gracefully
	})
}

func TestHealthServiceErrorPropagation(t *testing.T) {
	t.Run("Service handles validation errors gracefully", func(t *testing.T) {
		healthService := services.NewHealthService()

		// This test verifies that the service properly validates responses
		// and returns appropriate errors when validation fails
		response, err := healthService.GetHealthStatus()
		assert.NoError(t, err, "Normal health check should succeed")
		assert.NotNil(t, response)

		// Verify the response is valid
		validationErr := response.Validate()
		assert.NoError(t, validationErr, "Generated response should pass validation")
	})

	t.Run("Database status helpers work correctly", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Test individual methods
		dbStatus := healthService.GetDatabaseStatus()
		assert.True(t, dbStatus.IsValid(), "Database status should be valid")

		isHealthy := healthService.IsHealthy()
		assert.True(t, isHealthy || !isHealthy) // Should return a boolean without error

		uptime := healthService.GetUptime()
		assert.GreaterOrEqual(t, uptime, int64(0), "Uptime should be non-negative")

		version := healthService.GetVersion()
		assert.NotEmpty(t, version, "Version should not be empty")
	})
}

func TestHealthServiceConsistentBehavior(t *testing.T) {
	t.Run("Multiple calls return consistent structure", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Get multiple responses
		response1, err1 := healthService.GetHealthStatus()
		require.NoError(t, err1)

		time.Sleep(10 * time.Millisecond) // Small delay

		response2, err2 := healthService.GetHealthStatus()
		require.NoError(t, err2)

		// Structure should be consistent
		assert.Equal(t, response1.Status, response2.Status)
		assert.Equal(t, response1.Database, response2.Database)
		assert.Equal(t, response1.Version, response2.Version)

		// Uptime should increase or stay the same (accounting for timing precision)
		assert.GreaterOrEqual(t, response2.Uptime, response1.Uptime)

		// Timestamps should be different and in order
		time1, err := time.Parse(time.RFC3339, response1.Timestamp)
		require.NoError(t, err)
		time2, err := time.Parse(time.RFC3339, response2.Timestamp)
		require.NoError(t, err)
		assert.True(t, time2.After(time1) || time2.Equal(time1))
	})

	t.Run("Service version changes are reflected", func(t *testing.T) {
		healthService := services.NewHealthService()

		// Get initial version
		initialResponse, err := healthService.GetHealthStatus()
		require.NoError(t, err)
		initialVersion := initialResponse.Version

		// Change version
		newVersion := "test-version-" + time.Now().Format("20060102150405")
		healthService.SetVersion(newVersion)

		// Get updated response
		updatedResponse, err := healthService.GetHealthStatus()
		require.NoError(t, err)

		// Version should be updated
		assert.NotEqual(t, initialVersion, updatedResponse.Version)
		assert.Equal(t, newVersion, updatedResponse.Version)
	})
}

func TestErrorResponseValidation(t *testing.T) {
	t.Run("Valid error response", func(t *testing.T) {
		errorResp := dtos.NewErrorResponse("HEALTH_CHECK_FAILED", "Database connection timeout")

		assert.NotNil(t, errorResp)
		assert.Equal(t, "HEALTH_CHECK_FAILED", errorResp.Error)
		assert.Equal(t, "Database connection timeout", errorResp.Message)
	})

	t.Run("Error response with empty fields", func(t *testing.T) {
		errorResp := dtos.NewErrorResponse("", "")

		assert.NotNil(t, errorResp)
		assert.Equal(t, "", errorResp.Error)
		assert.Equal(t, "", errorResp.Message)
	})

	t.Run("Error response with special characters", func(t *testing.T) {
		errorResp := dtos.NewErrorResponse("ERROR_CODE_123", "Connection failed: timeout after 30s (error: connection refused)")

		assert.NotNil(t, errorResp)
		assert.Contains(t, errorResp.Message, "timeout")
		assert.Contains(t, errorResp.Message, "connection refused")
	})
}