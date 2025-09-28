package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/storage"
)

// TestServiceStatusMonitoring tests monitoring scenarios that simulate real-world usage
func TestServiceStatusMonitoring(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Repeated requests simulate monitoring polling", func(t *testing.T) {
		router := gin.New()

		// Initialize database for monitoring scenario
		db, err := storage.NewDatabase()
		require.NoError(t, err, "Database initialization should succeed")
		defer db.Close()

		// This will fail until the enhanced health handler is implemented
		// Current implementation doesn't provide monitoring-friendly response structure
		router.GET("/health", func(c *gin.Context) {
			// Current simple implementation that will make monitoring tests fail
			// Enhanced version should provide consistent structure for monitoring tools
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Simulate monitoring tool polling every 30 seconds
		pollInterval := 100 * time.Millisecond // Faster for testing
		pollCount := 5

		var responses []HealthResponse
		var responseTimes []time.Duration

		for i := 0; i < pollCount; i++ {
			start := time.Now()

			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err, "Request %d should be created successfully", i+1)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseTime := time.Since(start)
			responseTimes = append(responseTimes, responseTime)

			// All polling requests should succeed
			assert.Equal(t, http.StatusOK, w.Code,
				"Polling request %d should return 200 OK", i+1)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Polling request %d should return valid JSON", i+1)

			responses = append(responses, response)

			// Wait for next poll (except last iteration)
			if i < pollCount-1 {
				time.Sleep(pollInterval)
			}
		}

		// Verify consistent response structure across all polls
		for i, response := range responses {
			assert.Equal(t, "healthy", response.Status,
				"Poll %d should return consistent 'healthy' status", i+1)
			assert.Equal(t, "connected", response.Database,
				"Poll %d should return consistent 'connected' database status", i+1)
			assert.NotEmpty(t, response.Timestamp,
				"Poll %d should include timestamp", i+1)
		}

		// Verify response times are consistently fast
		for i, responseTime := range responseTimes {
			assert.True(t, responseTime < 200*time.Millisecond,
				"Poll %d response time should be < 200ms, got %v", i+1, responseTime)
		}
	})

	t.Run("Consistent response structure across monitoring calls", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Make multiple requests and verify structure consistency
		requestCount := 10
		var responseStructures []map[string]interface{}

		for i := 0; i < requestCount; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var responseMap map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseMap)
			require.NoError(t, err, "Request %d should return valid JSON", i+1)

			responseStructures = append(responseStructures, responseMap)
		}

		// All responses should have the same field structure
		firstStructure := responseStructures[0]
		requiredFields := []string{"status", "database", "timestamp"}

		for i, structure := range responseStructures {
			// Verify all required fields are present
			for _, field := range requiredFields {
				assert.Contains(t, structure, field,
					"Response %d should contain required field: %s", i+1, field)
			}

			// Verify field count consistency (allowing for optional fields)
			assert.True(t, len(structure) >= len(requiredFields),
				"Response %d should have at least %d fields", i+1, len(requiredFields))

			// Verify field types are consistent
			if status, exists := structure["status"]; exists {
				assert.IsType(t, "", status,
					"Response %d status field should be string", i+1)
			}

			if database, exists := structure["database"]; exists {
				assert.IsType(t, "", database,
					"Response %d database field should be string", i+1)
			}

			if timestamp, exists := structure["timestamp"]; exists {
				assert.IsType(t, "", timestamp,
					"Response %d timestamp field should be string", i+1)
			}
		}

		// Verify all responses have the same set of fields (consistency)
		firstFields := make([]string, 0, len(firstStructure))
		for field := range firstStructure {
			firstFields = append(firstFields, field)
		}

		for i, structure := range responseStructures[1:] {
			currentFields := make([]string, 0, len(structure))
			for field := range structure {
				currentFields = append(currentFields, field)
			}

			assert.ElementsMatch(t, firstFields, currentFields,
				"Response %d should have same fields as first response", i+2)
		}
	})

	t.Run("Timestamp updates correctly between requests", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Make requests with delays to verify timestamp updates
		var timestamps []time.Time
		requestCount := 3
		delay := 100 * time.Millisecond

		for i := 0; i < requestCount; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Request %d should return valid JSON", i+1)

			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, response.Timestamp)
			require.NoError(t, err, "Request %d should have valid timestamp", i+1)

			timestamps = append(timestamps, timestamp)

			// Wait before next request (except last iteration)
			if i < requestCount-1 {
				time.Sleep(delay)
			}
		}

		// Verify timestamps are updated and in chronological order
		for i := 1; i < len(timestamps); i++ {
			assert.True(t, timestamps[i].After(timestamps[i-1]) || timestamps[i].Equal(timestamps[i-1]),
				"Timestamp %d should be after or equal to timestamp %d", i+1, i)

			// Timestamps should be reasonably spaced (accounting for test execution time)
			timeDiff := timestamps[i].Sub(timestamps[i-1])
			assert.True(t, timeDiff >= 0,
				"Timestamp %d should not be before timestamp %d", i+1, i)
		}

		// All timestamps should be recent
		now := time.Now()
		for i, timestamp := range timestamps {
			timeDiff := now.Sub(timestamp)
			assert.True(t, timeDiff < 10*time.Second,
				"Timestamp %d should be recent (within 10 seconds), got %v ago", i+1, timeDiff)
		}
	})

	t.Run("No performance degradation under monitoring load", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Simulate intensive monitoring (high frequency polling)
		requestCount := 50
		var responseTimes []time.Duration
		var wg sync.WaitGroup
		responseTimesChan := make(chan time.Duration, requestCount)

		// Make concurrent requests to simulate high monitoring load
		for i := 0; i < requestCount; i++ {
			wg.Add(1)
			go func(requestNum int) {
				defer wg.Done()

				start := time.Now()

				req, err := http.NewRequest("GET", "/health", nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				responseTime := time.Since(start)
				responseTimesChan <- responseTime

				// Verify response is valid
				assert.Equal(t, http.StatusOK, w.Code,
					"Concurrent request %d should return 200 OK", requestNum)

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err,
					"Concurrent request %d should return valid JSON", requestNum)
			}(i + 1)
		}

		wg.Wait()
		close(responseTimesChan)

		// Collect response times
		for responseTime := range responseTimesChan {
			responseTimes = append(responseTimes, responseTime)
		}

		// Verify no significant performance degradation
		assert.Len(t, responseTimes, requestCount,
			"Should receive response times for all requests")

		// Calculate average response time
		var totalTime time.Duration
		for _, responseTime := range responseTimes {
			totalTime += responseTime
		}
		avgResponseTime := totalTime / time.Duration(len(responseTimes))

		// Average response time should be reasonable
		assert.True(t, avgResponseTime < 100*time.Millisecond,
			"Average response time should be < 100ms under load, got %v", avgResponseTime)

		// No individual request should be excessively slow
		for i, responseTime := range responseTimes {
			assert.True(t, responseTime < 500*time.Millisecond,
				"Request %d response time should be < 500ms under load, got %v", i+1, responseTime)
		}
	})

	t.Run("Monitoring compatibility with different HTTP clients", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test different User-Agent headers to simulate various monitoring tools
		userAgents := []string{
			"Prometheus/2.0",
			"Nagios/4.0",
			"DataDog/7.0",
			"NewRelic-Monitor/1.0",
			"curl/7.68.0",
			"",
		}

		for _, userAgent := range userAgents {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			if userAgent != "" {
				req.Header.Set("User-Agent", userAgent)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should work with any monitoring tool
			assert.Equal(t, http.StatusOK, w.Code,
				"Should return 200 OK for User-Agent: %s", userAgent)

			// Response structure should be consistent regardless of client
			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err,
				"Should return valid JSON for User-Agent: %s", userAgent)

			// Required fields should be present for all monitoring tools
			assert.NotEmpty(t, response.Status,
				"Status should be present for User-Agent: %s", userAgent)
			assert.NotEmpty(t, response.Database,
				"Database status should be present for User-Agent: %s", userAgent)
			assert.NotEmpty(t, response.Timestamp,
				"Timestamp should be present for User-Agent: %s", userAgent)
		}
	})

	t.Run("Health check caching behavior for monitoring", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Make rapid successive requests to test caching behavior
		rapidRequestCount := 5
		var responses []HealthResponse

		for i := 0; i < rapidRequestCount; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			responses = append(responses, response)

			// Small delay between requests
			if i < rapidRequestCount-1 {
				time.Sleep(10 * time.Millisecond)
			}
		}

		// Verify responses are consistent (but timestamps may vary)
		for i, response := range responses {
			assert.Equal(t, "healthy", response.Status,
				"Rapid request %d should return consistent status", i+1)
			assert.Equal(t, "connected", response.Database,
				"Rapid request %d should return consistent database status", i+1)
			assert.NotEmpty(t, response.Timestamp,
				"Rapid request %d should have timestamp", i+1)
		}

		// All requests should complete quickly (no blocking)
		start := time.Now()
		for i := 0; i < 10; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
		totalTime := time.Since(start)

		// 10 requests should complete very quickly
		assert.True(t, totalTime < 1*time.Second,
			"10 rapid requests should complete in < 1 second, got %v", totalTime)
	})
}