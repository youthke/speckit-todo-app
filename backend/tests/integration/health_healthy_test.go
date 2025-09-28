package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/storage"
)

// HealthResponse represents the expected enhanced health endpoint response
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version,omitempty"`
	Uptime    int64  `json:"uptime,omitempty"`
}

// TestHealthyServiceScenario tests the complete end-to-end healthy service scenario
func TestHealthyServiceScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize database for healthy scenario
	db, err := storage.NewDatabase()
	require.NoError(t, err, "Database initialization should succeed for healthy scenario")
	defer db.Close()

	// This will fail until the enhanced health handler is implemented
	// Currently the endpoint only returns {"status": "ok"} but we expect full health information
	router.GET("/health", func(c *gin.Context) {
		// This is the current simple implementation that will make our tests fail
		// The enhanced version should check database connectivity, return proper status enum values, etc.
		c.JSON(200, gin.H{"status": "ok"})
	})

	t.Run("End-to-end healthy service scenario", func(t *testing.T) {
		// Record start time for response time validation
		start := time.Now()

		// Step 1: Call /health endpoint
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err, "HTTP request creation should succeed")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Step 2: Verify response time < 200ms
		responseTime := time.Since(start)
		assert.True(t, responseTime < 200*time.Millisecond,
			"Response time should be less than 200ms, got %v", responseTime)

		// Step 3: Verify 200 status code
		assert.Equal(t, http.StatusOK, w.Code,
			"Health endpoint should return 200 OK for healthy service")

		// Step 4: Verify content type
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"),
			"Content-Type should be application/json")

		// Step 5: Parse and validate response structure
		var response HealthResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON matching HealthResponse structure")

		// Step 6: Validate all required fields are present
		var responseMap map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &responseMap)
		require.NoError(t, err, "Response should be valid JSON")

		requiredFields := []string{"status", "database", "timestamp"}
		for _, field := range requiredFields {
			assert.Contains(t, responseMap, field,
				"Response must contain required field: %s", field)
		}

		// Step 7: Validate field values for healthy scenario
		assert.Equal(t, "healthy", response.Status,
			"Status should be 'healthy' for a healthy service")
		assert.Equal(t, "connected", response.Database,
			"Database status should be 'connected' when database is accessible")

		// Step 8: Validate timestamp format and recency
		parsedTime, err := time.Parse(time.RFC3339, response.Timestamp)
		assert.NoError(t, err, "Timestamp should be valid ISO 8601/RFC3339 format")

		timeDiff := time.Since(parsedTime)
		assert.True(t, timeDiff < 5*time.Second,
			"Timestamp should be recent (within last 5 seconds), got %v ago", timeDiff)

		// Step 9: Validate optional fields if present
		if response.Version != "" {
			assert.NotEmpty(t, response.Version, "Version should not be empty if present")
		}

		if response.Uptime > 0 {
			assert.True(t, response.Uptime >= 0, "Uptime should be non-negative if present")
		}
	})

	t.Run("Consistent healthy response over multiple calls", func(t *testing.T) {
		// Make multiple requests to ensure consistent healthy behavior
		for i := 0; i < 3; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// All requests should return 200 OK
			assert.Equal(t, http.StatusOK, w.Code,
				"Request %d should return 200 OK for healthy service", i+1)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response %d should be valid JSON", i+1)

			// Status should consistently be "healthy"
			assert.Equal(t, "healthy", response.Status,
				"Request %d should return 'healthy' status", i+1)

			// Database should consistently be "connected"
			assert.Equal(t, "connected", response.Database,
				"Request %d should return 'connected' database status", i+1)
		}
	})

	t.Run("Response structure validation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var responseMap map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &responseMap)
		require.NoError(t, err)

		// Validate field types
		if status, exists := responseMap["status"]; exists {
			assert.IsType(t, "", status, "Status field must be a string")
		}

		if database, exists := responseMap["database"]; exists {
			assert.IsType(t, "", database, "Database field must be a string")
		}

		if timestamp, exists := responseMap["timestamp"]; exists {
			assert.IsType(t, "", timestamp, "Timestamp field must be a string")
		}

		if version, exists := responseMap["version"]; exists {
			assert.IsType(t, "", version, "Version field must be a string")
		}

		if uptime, exists := responseMap["uptime"]; exists {
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
	})

	t.Run("Database connectivity verification", func(t *testing.T) {
		// Verify database is actually connected before health check
		assert.NotNil(t, db, "Database connection should be established")

		// Perform a simple database operation to confirm connectivity
		err := db.Ping()
		require.NoError(t, err, "Database should be reachable before health check")

		// Now call health endpoint
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response HealthResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Health endpoint should reflect the actual database connectivity
		assert.Equal(t, "connected", response.Database,
			"Health endpoint should report 'connected' when database is accessible")
	})
}