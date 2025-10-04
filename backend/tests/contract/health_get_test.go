package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todo-app/internal/dtos"
	"todo-app/internal/services"
	"todo-app/internal/storage"
)

// HealthResponse represents the expected structure of the health endpoint response
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version,omitempty"`
	Uptime    int64  `json:"uptime,omitempty"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func TestGetHealthContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize database for testing
	err := storage.InitDatabase()
	require.NoError(t, err, "Failed to initialize database for testing")

	// Create health service
	healthService := services.NewHealthService()

	// Setup router with enhanced health endpoint
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		healthResponse, err := healthService.GetHealthStatus()
		if err != nil {
			errorResponse := models.NewErrorResponse("internal_error", "Health check failed unexpectedly")
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}

		// Determine HTTP status code based on health status
		var statusCode int
		switch healthResponse.Status {
		case models.HealthStatusHealthy:
			statusCode = http.StatusOK
		case models.HealthStatusDegraded:
			statusCode = http.StatusServiceUnavailable
		case models.HealthStatusUnhealthy:
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, healthResponse)
	})

	tests := []struct {
		name               string
		expectedStatus     int
		expectedFields     []string
		requiredFields     []string
		optionalFields     []string
		validateResponse   func(t *testing.T, response map[string]interface{})
		description        string
	}{
		{
			name:           "Healthy service response",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"status", "database", "timestamp"},
			requiredFields: []string{"status", "database", "timestamp"},
			optionalFields: []string{"version", "uptime"},
			description:    "Should return 200 with complete health information when service is healthy",
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				// Validate status enum values
				status := response["status"].(string)
				validStatuses := []string{"healthy", "degraded", "unhealthy"}
				assert.Contains(t, validStatuses, status, "Status must be one of: healthy, degraded, unhealthy")

				// Validate database enum values
				if database, exists := response["database"]; exists {
					databaseStr := database.(string)
					validDatabaseStates := []string{"connected", "disconnected", "error"}
					assert.Contains(t, validDatabaseStates, databaseStr, "Database must be one of: connected, disconnected, error")
				}

				// Validate timestamp format (ISO 8601)
				if timestamp, exists := response["timestamp"]; exists {
					timestampStr := timestamp.(string)
					parsedTime, err := time.Parse(time.RFC3339, timestampStr)
					assert.NoError(t, err, "Timestamp must be valid ISO 8601 format")

					// Validate timestamp is recent (within last 5 seconds)
					timeDiff := time.Since(parsedTime)
					assert.True(t, timeDiff < 5*time.Second, "Timestamp should be recent")
				}

				// Validate version field if present
				if version, exists := response["version"]; exists {
					assert.IsType(t, "", version, "Version must be a string")
					assert.NotEmpty(t, version, "Version should not be empty if present")
				}

				// Validate uptime field if present
				if uptime, exists := response["uptime"]; exists {
					// JSON numbers are float64 by default
					switch v := uptime.(type) {
					case float64:
						assert.True(t, v >= 0, "Uptime must be non-negative")
					case int64:
						assert.True(t, v >= 0, "Uptime must be non-negative")
					default:
						t.Errorf("Uptime must be a number, got %T", v)
					}
				}
			},
		},
		{
			name:           "Degraded service response",
			expectedStatus: http.StatusServiceUnavailable,
			expectedFields: []string{"status", "database", "timestamp"},
			requiredFields: []string{"status", "database", "timestamp"},
			optionalFields: []string{"version", "uptime"},
			description:    "Should return 503 when service is degraded",
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				// For degraded state, status should be "degraded" or "unhealthy"
				if status, exists := response["status"]; exists {
					statusStr := status.(string)
					degradedStatuses := []string{"degraded", "unhealthy"}
					assert.Contains(t, degradedStatuses, statusStr, "Status should indicate degraded service")
				}

				// Database might be disconnected or error
				if database, exists := response["database"]; exists {
					databaseStr := database.(string)
					validDatabaseStates := []string{"connected", "disconnected", "error"}
					assert.Contains(t, validDatabaseStates, databaseStr, "Database must be one of valid states")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until enhanced health endpoint is implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "HTTP status code should match expected")

			// Validate content type
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			// Parse response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Verify required fields exist
			for _, field := range tt.requiredFields {
				assert.Contains(t, response, field, "Response must contain required field: %s", field)
			}

			// Verify response structure matches health-api.yaml schema
			t.Run("Schema validation", func(t *testing.T) {
				// Check that response contains only expected fields
				allowedFields := append(tt.requiredFields, tt.optionalFields...)
				for field := range response {
					assert.Contains(t, allowedFields, field, "Response contains unexpected field: %s", field)
				}

				// Validate field types
				if status, exists := response["status"]; exists {
					assert.IsType(t, "", status, "Status must be a string")
				}

				if database, exists := response["database"]; exists {
					assert.IsType(t, "", database, "Database must be a string")
				}

				if timestamp, exists := response["timestamp"]; exists {
					assert.IsType(t, "", timestamp, "Timestamp must be a string")
				}
			})

			// Run custom validation
			if tt.validateResponse != nil {
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestGetHealthContractFieldValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize database for testing
	err := storage.InitDatabase()
	require.NoError(t, err, "Failed to initialize database for testing")

	// Create health service
	healthService := services.NewHealthService()

	// Setup router with enhanced health endpoint
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		healthResponse, err := healthService.GetHealthStatus()
		if err != nil {
			errorResponse := models.NewErrorResponse("internal_error", "Health check failed unexpectedly")
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}

		// Determine HTTP status code based on health status
		var statusCode int
		switch healthResponse.Status {
		case models.HealthStatusHealthy:
			statusCode = http.StatusOK
		case models.HealthStatusDegraded:
			statusCode = http.StatusServiceUnavailable
		case models.HealthStatusUnhealthy:
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, healthResponse)
	})

	t.Run("Required fields validation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// These assertions will fail until enhanced endpoint is implemented
		requiredFields := []string{"status", "database", "timestamp"}
		for _, field := range requiredFields {
			assert.Contains(t, response, field, "Required field %s is missing", field)
		}
	})

	t.Run("Status enum validation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// This will fail until status field returns proper enum values
		if status, exists := response["status"]; exists {
			validStatuses := []string{"healthy", "degraded", "unhealthy"}
			assert.Contains(t, validStatuses, status, "Status must be one of: %v", validStatuses)
		}
	})

	t.Run("Database enum validation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// This will fail until database field is implemented
		if database, exists := response["database"]; exists {
			validDatabaseStates := []string{"connected", "disconnected", "error"}
			assert.Contains(t, validDatabaseStates, database, "Database must be one of: %v", validDatabaseStates)
		} else {
			t.Error("Database field is required but missing")
		}
	})

	t.Run("Timestamp format validation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// This will fail until timestamp field is implemented
		if timestamp, exists := response["timestamp"]; exists {
			timestampStr := timestamp.(string)
			_, err := time.Parse(time.RFC3339, timestampStr)
			assert.NoError(t, err, "Timestamp must be valid ISO 8601 format (RFC3339)")
		} else {
			t.Error("Timestamp field is required but missing")
		}
	})
}

func TestGetHealthContractErrorScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This simulates a health endpoint that might return errors
	router.GET("/health", func(c *gin.Context) {
		// This will be replaced with actual health check logic
		// For now, it returns the simple response that will make tests fail
		c.JSON(200, gin.H{"status": "ok"})
	})

	t.Run("JSON response structure", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Validate JSON structure
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response must be valid JSON")

		// Response should not be empty
		assert.NotEmpty(t, response, "Response should not be empty")
	})

	t.Run("Content type header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		contentType := w.Header().Get("Content-Type")
		assert.Equal(t, "application/json; charset=utf-8", contentType, "Content-Type must be application/json")
	})

	t.Run("HTTP method validation", func(t *testing.T) {
		// Test that only GET method is supported
		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		for _, method := range methods {
			req, err := http.NewRequest(method, "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 or 405 for unsupported methods
			assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMethodNotAllowed,
				"Method %s should not be allowed on /health endpoint", method)
		}
	})
}

// TestGetHealthContractIntegration tests the health endpoint in a more realistic scenario
func TestGetHealthContractIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This represents the current implementation that will fail the contract tests
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	t.Run("Response time validation", func(t *testing.T) {
		start := time.Now()

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		duration := time.Since(start)

		// Health checks should be fast (under 1 second)
		assert.True(t, duration < time.Second, "Health check should complete within 1 second")
	})

	t.Run("Multiple requests consistency", func(t *testing.T) {
		// Make multiple requests to ensure consistent response structure
		for i := 0; i < 3; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "All health checks should return consistent status")

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "All responses should be valid JSON")
		}
	})
}