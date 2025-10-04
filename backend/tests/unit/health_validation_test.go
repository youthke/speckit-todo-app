package unit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/dtos"
)

func TestHealthStatusEnum(t *testing.T) {
	tests := []struct {
		name     string
		status   models.HealthStatus
		expected bool
	}{
		{
			name:     "Valid healthy status",
			status:   models.HealthStatusHealthy,
			expected: true,
		},
		{
			name:     "Valid degraded status",
			status:   models.HealthStatusDegraded,
			expected: true,
		},
		{
			name:     "Valid unhealthy status",
			status:   models.HealthStatusUnhealthy,
			expected: true,
		},
		{
			name:     "Invalid empty status",
			status:   "",
			expected: false,
		},
		{
			name:     "Invalid unknown status",
			status:   "unknown",
			expected: false,
		},
		{
			name:     "Invalid mixed case status",
			status:   "Healthy",
			expected: false,
		},
		{
			name:     "Invalid status with spaces",
			status:   " healthy ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			assert.Equal(t, tt.expected, result, "IsValid() result should match expected for %s", tt.name)
		})
	}
}

func TestDatabaseStatusEnum(t *testing.T) {
	tests := []struct {
		name     string
		status   models.DatabaseStatus
		expected bool
	}{
		{
			name:     "Valid connected status",
			status:   models.DatabaseStatusConnected,
			expected: true,
		},
		{
			name:     "Valid disconnected status",
			status:   models.DatabaseStatusDisconnected,
			expected: true,
		},
		{
			name:     "Valid error status",
			status:   models.DatabaseStatusError,
			expected: true,
		},
		{
			name:     "Invalid empty status",
			status:   "",
			expected: false,
		},
		{
			name:     "Invalid unknown status",
			status:   "unknown",
			expected: false,
		},
		{
			name:     "Invalid mixed case status",
			status:   "Connected",
			expected: false,
		},
		{
			name:     "Invalid status with spaces",
			status:   " connected ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			assert.Equal(t, tt.expected, result, "IsValid() result should match expected for %s", tt.name)
		})
	}
}

func TestHealthResponseValidation(t *testing.T) {
	validTimestamp := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name        string
		response    models.HealthResponse
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid health response",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Version:   "1.0.0",
				Uptime:    3600,
			},
			expectError: false,
		},
		{
			name: "Valid health response without optional fields",
			response: models.HealthResponse{
				Status:    models.HealthStatusDegraded,
				Database:  models.DatabaseStatusDisconnected,
				Timestamp: validTimestamp,
			},
			expectError: false,
		},
		{
			name: "Invalid health status",
			response: models.HealthResponse{
				Status:    "invalid",
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
			},
			expectError: true,
			errorMsg:    "invalid status",
		},
		{
			name: "Invalid database status",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  "invalid",
				Timestamp: validTimestamp,
			},
			expectError: true,
			errorMsg:    "invalid database status",
		},
		{
			name: "Empty timestamp",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: "",
			},
			expectError: true,
			errorMsg:    "timestamp cannot be empty",
		},
		{
			name: "Invalid timestamp format",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: "2023-01-01 12:00:00",
			},
			expectError: true,
			errorMsg:    "invalid timestamp format",
		},
		{
			name: "Negative uptime",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Uptime:    -1,
			},
			expectError: true,
			errorMsg:    "uptime must be non-negative",
		},
		{
			name: "Empty version string",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Version:   "",
			},
			expectError: false, // Empty version is allowed
		},
		{
			name: "Whitespace-only version",
			response: models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Version:   "   ",
			},
			expectError: true,
			errorMsg:    "version cannot be empty or whitespace-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Validate()

			if tt.expectError {
				assert.Error(t, err, "Expected validation to fail for %s", tt.name)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected validation to pass for %s", tt.name)
			}
		})
	}
}

func TestTimestampFormatValidation(t *testing.T) {
	tests := []struct {
		name        string
		timestamp   string
		expectError bool
	}{
		{
			name:        "Valid RFC3339 timestamp",
			timestamp:   "2023-12-01T15:30:45Z",
			expectError: false,
		},
		{
			name:        "Valid RFC3339 with timezone",
			timestamp:   "2023-12-01T15:30:45+05:00",
			expectError: false,
		},
		{
			name:        "Valid RFC3339 with nanoseconds",
			timestamp:   "2023-12-01T15:30:45.123456789Z",
			expectError: false,
		},
		{
			name:        "Invalid format - missing timezone",
			timestamp:   "2023-12-01T15:30:45",
			expectError: true,
		},
		{
			name:        "Invalid format - date only",
			timestamp:   "2023-12-01",
			expectError: true,
		},
		{
			name:        "Invalid format - time only",
			timestamp:   "15:30:45",
			expectError: true,
		},
		{
			name:        "Invalid format - non-ISO",
			timestamp:   "Dec 1, 2023 3:30:45 PM",
			expectError: true,
		},
		{
			name:        "Invalid format - unix timestamp",
			timestamp:   "1701439845",
			expectError: true,
		},
		{
			name:        "Empty timestamp",
			timestamp:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: tt.timestamp,
			}

			err := response.Validate()

			if tt.expectError {
				assert.Error(t, err, "Expected validation to fail for timestamp: %s", tt.timestamp)
			} else {
				assert.NoError(t, err, "Expected validation to pass for timestamp: %s", tt.timestamp)
			}
		})
	}
}

func TestVersionValidation(t *testing.T) {
	validTimestamp := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "Valid semantic version",
			version:     "1.0.0",
			expectError: false,
		},
		{
			name:        "Valid version with build",
			version:     "1.0.0-beta.1",
			expectError: false,
		},
		{
			name:        "Valid git hash version",
			version:     "abc123def",
			expectError: false,
		},
		{
			name:        "Valid complex version",
			version:     "v1.2.3-rc.1+build.20231201",
			expectError: false,
		},
		{
			name:        "Empty version (allowed)",
			version:     "",
			expectError: false,
		},
		{
			name:        "Whitespace-only version",
			version:     "   ",
			expectError: true,
		},
		{
			name:        "Tab-only version",
			version:     "\t",
			expectError: true,
		},
		{
			name:        "Newline-only version",
			version:     "\n",
			expectError: true,
		},
		{
			name:        "Mixed whitespace version",
			version:     " \t\n ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Version:   tt.version,
			}

			err := response.Validate()

			if tt.expectError {
				assert.Error(t, err, "Expected validation to fail for version: %q", tt.version)
				assert.Contains(t, err.Error(), "version cannot be empty or whitespace-only")
			} else {
				assert.NoError(t, err, "Expected validation to pass for version: %q", tt.version)
			}
		})
	}
}

func TestUptimeValidation(t *testing.T) {
	validTimestamp := time.Now().UTC().Format(time.RFC3339)

	tests := []struct {
		name        string
		uptime      int64
		expectError bool
	}{
		{
			name:        "Zero uptime",
			uptime:      0,
			expectError: false,
		},
		{
			name:        "Positive uptime",
			uptime:      3600,
			expectError: false,
		},
		{
			name:        "Large uptime",
			uptime:      86400 * 365, // 1 year in seconds
			expectError: false,
		},
		{
			name:        "Maximum int64 uptime",
			uptime:      9223372036854775807, // max int64
			expectError: false,
		},
		{
			name:        "Negative uptime",
			uptime:      -1,
			expectError: true,
		},
		{
			name:        "Large negative uptime",
			uptime:      -9223372036854775808, // min int64
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := models.HealthResponse{
				Status:    models.HealthStatusHealthy,
				Database:  models.DatabaseStatusConnected,
				Timestamp: validTimestamp,
				Uptime:    tt.uptime,
			}

			err := response.Validate()

			if tt.expectError {
				assert.Error(t, err, "Expected validation to fail for uptime: %d", tt.uptime)
				assert.Contains(t, err.Error(), "uptime must be non-negative")
			} else {
				assert.NoError(t, err, "Expected validation to pass for uptime: %d", tt.uptime)
			}
		})
	}
}

func TestHealthStatusJSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		status   models.HealthStatus
		expected string
	}{
		{
			name:     "Healthy status",
			status:   models.HealthStatusHealthy,
			expected: `"healthy"`,
		},
		{
			name:     "Degraded status",
			status:   models.HealthStatusDegraded,
			expected: `"degraded"`,
		},
		{
			name:     "Unhealthy status",
			status:   models.HealthStatusUnhealthy,
			expected: `"unhealthy"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.status)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestHealthStatusJSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expected    models.HealthStatus
		expectError bool
	}{
		{
			name:        "Valid healthy status",
			jsonData:    `"healthy"`,
			expected:    models.HealthStatusHealthy,
			expectError: false,
		},
		{
			name:        "Valid degraded status",
			jsonData:    `"degraded"`,
			expected:    models.HealthStatusDegraded,
			expectError: false,
		},
		{
			name:        "Valid unhealthy status",
			jsonData:    `"unhealthy"`,
			expected:    models.HealthStatusUnhealthy,
			expectError: false,
		},
		{
			name:        "Invalid status",
			jsonData:    `"invalid"`,
			expectError: true,
		},
		{
			name:        "Empty status",
			jsonData:    `""`,
			expectError: true,
		},
		{
			name:        "Null status",
			jsonData:    `null`,
			expectError: true,
		},
		{
			name:        "Invalid JSON format",
			jsonData:    `healthy`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status models.HealthStatus
			err := json.Unmarshal([]byte(tt.jsonData), &status)

			if tt.expectError {
				assert.Error(t, err, "Expected unmarshaling to fail for %s", tt.jsonData)
			} else {
				assert.NoError(t, err, "Expected unmarshaling to succeed for %s", tt.jsonData)
				assert.Equal(t, tt.expected, status)
			}
		})
	}
}

func TestDatabaseStatusJSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		status   models.DatabaseStatus
		expected string
	}{
		{
			name:     "Connected status",
			status:   models.DatabaseStatusConnected,
			expected: `"connected"`,
		},
		{
			name:     "Disconnected status",
			status:   models.DatabaseStatusDisconnected,
			expected: `"disconnected"`,
		},
		{
			name:     "Error status",
			status:   models.DatabaseStatusError,
			expected: `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.status)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestDatabaseStatusJSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expected    models.DatabaseStatus
		expectError bool
	}{
		{
			name:        "Valid connected status",
			jsonData:    `"connected"`,
			expected:    models.DatabaseStatusConnected,
			expectError: false,
		},
		{
			name:        "Valid disconnected status",
			jsonData:    `"disconnected"`,
			expected:    models.DatabaseStatusDisconnected,
			expectError: false,
		},
		{
			name:        "Valid error status",
			jsonData:    `"error"`,
			expected:    models.DatabaseStatusError,
			expectError: false,
		},
		{
			name:        "Invalid status",
			jsonData:    `"invalid"`,
			expectError: true,
		},
		{
			name:        "Empty status",
			jsonData:    `""`,
			expectError: true,
		},
		{
			name:        "Null status",
			jsonData:    `null`,
			expectError: true,
		},
		{
			name:        "Invalid JSON format",
			jsonData:    `connected`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status models.DatabaseStatus
			err := json.Unmarshal([]byte(tt.jsonData), &status)

			if tt.expectError {
				assert.Error(t, err, "Expected unmarshaling to fail for %s", tt.jsonData)
			} else {
				assert.NoError(t, err, "Expected unmarshaling to succeed for %s", tt.jsonData)
				assert.Equal(t, tt.expected, status)
			}
		})
	}
}

func TestNewHealthResponse(t *testing.T) {
	status := models.HealthStatusHealthy
	database := models.DatabaseStatusConnected
	version := "1.0.0"
	uptime := int64(3600)

	response := models.NewHealthResponse(status, database, version, uptime)

	assert.NotNil(t, response)
	assert.Equal(t, status, response.Status)
	assert.Equal(t, database, response.Database)
	assert.Equal(t, version, response.Version)
	assert.Equal(t, uptime, response.Uptime)
	assert.NotEmpty(t, response.Timestamp)

	// Validate that timestamp is in correct format
	_, err := time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err, "Timestamp should be in RFC3339 format")

	// Validate the response itself
	err = response.Validate()
	assert.NoError(t, err, "Generated response should be valid")
}

func TestDetermineOverallHealth(t *testing.T) {
	tests := []struct {
		name           string
		databaseStatus models.DatabaseStatus
		expected       models.HealthStatus
	}{
		{
			name:           "Connected database results in healthy",
			databaseStatus: models.DatabaseStatusConnected,
			expected:       models.HealthStatusHealthy,
		},
		{
			name:           "Disconnected database results in degraded",
			databaseStatus: models.DatabaseStatusDisconnected,
			expected:       models.HealthStatusDegraded,
		},
		{
			name:           "Error database results in unhealthy",
			databaseStatus: models.DatabaseStatusError,
			expected:       models.HealthStatusUnhealthy,
		},
		{
			name:           "Invalid database status results in unhealthy",
			databaseStatus: "invalid",
			expected:       models.HealthStatusUnhealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.DetermineOverallHealth(tt.databaseStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHealthStatusString(t *testing.T) {
	tests := []struct {
		name     string
		status   models.HealthStatus
		expected string
	}{
		{
			name:     "Healthy status string",
			status:   models.HealthStatusHealthy,
			expected: "healthy",
		},
		{
			name:     "Degraded status string",
			status:   models.HealthStatusDegraded,
			expected: "degraded",
		},
		{
			name:     "Unhealthy status string",
			status:   models.HealthStatusUnhealthy,
			expected: "unhealthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseStatusString(t *testing.T) {
	tests := []struct {
		name     string
		status   models.DatabaseStatus
		expected string
	}{
		{
			name:     "Connected status string",
			status:   models.DatabaseStatusConnected,
			expected: "connected",
		},
		{
			name:     "Disconnected status string",
			status:   models.DatabaseStatusDisconnected,
			expected: "disconnected",
		},
		{
			name:     "Error status string",
			status:   models.DatabaseStatusError,
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}