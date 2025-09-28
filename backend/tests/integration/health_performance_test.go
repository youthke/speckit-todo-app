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

// TestResponseTimeValidation tests performance characteristics of the health endpoint
func TestResponseTimeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Response time under normal load", func(t *testing.T) {
		router := gin.New()

		// Initialize database for performance testing
		db, err := storage.NewDatabase()
		require.NoError(t, err, "Database initialization should succeed")
		defer db.Close()

		// This will fail until the enhanced health handler is implemented
		// Current implementation doesn't include database checks that might affect performance
		router.GET("/health", func(c *gin.Context) {
			// Current simple implementation that will make performance tests fail
			// Enhanced version should maintain performance while checking database
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test single request performance
		testCount := 10
		var responseTimes []time.Duration

		for i := 0; i < testCount; i++ {
			start := time.Now()

			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err, "Request %d should be created successfully", i+1)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseTime := time.Since(start)
			responseTimes = append(responseTimes, responseTime)

			// Verify successful response
			assert.Equal(t, http.StatusOK, w.Code,
				"Request %d should return 200 OK", i+1)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Request %d should return valid JSON", i+1)

			// Each individual request should be fast
			assert.True(t, responseTime < 200*time.Millisecond,
				"Request %d should complete in < 200ms, got %v", i+1, responseTime)
		}

		// Calculate statistics
		var totalTime time.Duration
		var maxTime time.Duration
		var minTime time.Duration = time.Hour // Start with large value

		for _, responseTime := range responseTimes {
			totalTime += responseTime
			if responseTime > maxTime {
				maxTime = responseTime
			}
			if responseTime < minTime {
				minTime = responseTime
			}
		}

		avgTime := totalTime / time.Duration(len(responseTimes))

		// Performance assertions
		assert.True(t, avgTime < 50*time.Millisecond,
			"Average response time should be < 50ms, got %v", avgTime)
		assert.True(t, maxTime < 200*time.Millisecond,
			"Maximum response time should be < 200ms, got %v", maxTime)
		assert.True(t, minTime < 100*time.Millisecond,
			"Minimum response time should be < 100ms, got %v", minTime)

		t.Logf("Performance stats - Avg: %v, Min: %v, Max: %v", avgTime, minTime, maxTime)
	})

	t.Run("Consistent response time under normal load", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test response time consistency
		testCount := 20
		var responseTimes []time.Duration

		for i := 0; i < testCount; i++ {
			start := time.Now()

			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseTime := time.Since(start)
			responseTimes = append(responseTimes, responseTime)

			assert.Equal(t, http.StatusOK, w.Code)

			// Add small delay to avoid overwhelming the handler
			time.Sleep(10 * time.Millisecond)
		}

		// Calculate variance to ensure consistency
		var totalTime time.Duration
		for _, responseTime := range responseTimes {
			totalTime += responseTime
		}
		avgTime := totalTime / time.Duration(len(responseTimes))

		// Check variance (all response times should be reasonably close to average)
		var varianceSum time.Duration
		for _, responseTime := range responseTimes {
			diff := responseTime - avgTime
			if diff < 0 {
				diff = -diff
			}
			varianceSum += diff
		}
		avgVariance := varianceSum / time.Duration(len(responseTimes))

		// Response times should be consistent (low variance)
		assert.True(t, avgVariance < 50*time.Millisecond,
			"Response time variance should be < 50ms, got %v", avgVariance)

		// No single request should be significantly slower than others
		for i, responseTime := range responseTimes {
			assert.True(t, responseTime < avgTime+100*time.Millisecond,
				"Request %d response time (%v) should not be much slower than average (%v)", i+1, responseTime, avgTime)
		}
	})

	t.Run("Performance with concurrent requests", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test concurrent request performance
		concurrentRequests := 20
		var wg sync.WaitGroup
		responseTimesChan := make(chan time.Duration, concurrentRequests)
		statusCodesChan := make(chan int, concurrentRequests)

		overallStart := time.Now()

		// Launch concurrent requests
		for i := 0; i < concurrentRequests; i++ {
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
				statusCodesChan <- w.Code

				// Verify response structure
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Concurrent request %d should return valid JSON", requestNum)
			}(i + 1)
		}

		wg.Wait()
		overallTime := time.Since(overallStart)

		close(responseTimesChan)
		close(statusCodesChan)

		// Collect results
		var responseTimes []time.Duration
		var statusCodes []int

		for responseTime := range responseTimesChan {
			responseTimes = append(responseTimes, responseTime)
		}

		for statusCode := range statusCodesChan {
			statusCodes = append(statusCodes, statusCode)
		}

		// Verify all requests completed successfully
		assert.Len(t, responseTimes, concurrentRequests, "Should receive all response times")
		assert.Len(t, statusCodes, concurrentRequests, "Should receive all status codes")

		// All requests should return 200 OK
		for i, statusCode := range statusCodes {
			assert.Equal(t, http.StatusOK, statusCode,
				"Concurrent request %d should return 200 OK", i+1)
		}

		// Individual response times should still be fast under concurrency
		for i, responseTime := range responseTimes {
			assert.True(t, responseTime < 500*time.Millisecond,
				"Concurrent request %d should complete in < 500ms, got %v", i+1, responseTime)
		}

		// Overall time should be reasonable (should handle concurrency well)
		assert.True(t, overallTime < 2*time.Second,
			"All %d concurrent requests should complete in < 2 seconds, got %v", concurrentRequests, overallTime)

		// Calculate average response time under concurrency
		var totalTime time.Duration
		for _, responseTime := range responseTimes {
			totalTime += responseTime
		}
		avgTime := totalTime / time.Duration(len(responseTimes))

		assert.True(t, avgTime < 200*time.Millisecond,
			"Average response time under concurrency should be < 200ms, got %v", avgTime)

		t.Logf("Concurrent performance - Requests: %d, Overall time: %v, Avg response: %v",
			concurrentRequests, overallTime, avgTime)
	})

	t.Run("No blocking operations in health check", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test that health check doesn't block on long-running operations
		// This simulates what happens when database checks are non-blocking

		maxResponseTime := 1 * time.Second
		testRuns := 5

		for run := 0; run < testRuns; run++ {
			start := time.Now()

			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseTime := time.Since(start)

			// Health check should never block for extended periods
			assert.True(t, responseTime < maxResponseTime,
				"Run %d: Health check should not block, got %v", run+1, responseTime)

			assert.Equal(t, http.StatusOK, w.Code,
				"Run %d should return 200 OK", run+1)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Run %d should return valid JSON", run+1)
		}
	})

	t.Run("Database check performance optimization", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		// The enhanced version should optimize database connectivity checks
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test database connectivity check doesn't significantly slow down response
		testCount := 15
		var responseTimes []time.Duration

		for i := 0; i < testCount; i++ {
			start := time.Now()

			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseTime := time.Since(start)
			responseTimes = append(responseTimes, responseTime)

			assert.Equal(t, http.StatusOK, w.Code,
				"Request %d should return 200 OK", i+1)

			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Request %d should return valid JSON", i+1)

			// Even with database checks, response should be fast
			assert.True(t, responseTime < 150*time.Millisecond,
				"Request %d with DB check should complete in < 150ms, got %v", i+1, responseTime)

			// Database status should be checked and reported
			validDatabaseStates := []string{"connected", "disconnected", "error"}
			assert.Contains(t, validDatabaseStates, response.Database,
				"Request %d should report valid database state", i+1)
		}

		// Calculate performance impact of database checks
		var totalTime time.Duration
		for _, responseTime := range responseTimes {
			totalTime += responseTime
		}
		avgTime := totalTime / time.Duration(len(responseTimes))

		// Average time should still be very fast even with database checks
		assert.True(t, avgTime < 100*time.Millisecond,
			"Average response time with DB checks should be < 100ms, got %v", avgTime)
	})

	t.Run("Performance under stress conditions", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test performance under high load (stress test)
		stressRequests := 100
		concurrency := 10
		var wg sync.WaitGroup
		responseTimesChan := make(chan time.Duration, stressRequests)
		errorsChan := make(chan error, stressRequests)

		// Rate limiting: launch requests in batches
		requestsPerBatch := stressRequests / concurrency
		overallStart := time.Now()

		for batch := 0; batch < concurrency; batch++ {
			wg.Add(1)
			go func(batchNum int) {
				defer wg.Done()

				for i := 0; i < requestsPerBatch; i++ {
					start := time.Now()

					req, err := http.NewRequest("GET", "/health", nil)
					if err != nil {
						errorsChan <- err
						continue
					}

					w := httptest.NewRecorder()
					router.ServeHTTP(w, req)

					responseTime := time.Since(start)
					responseTimesChan <- responseTime

					if w.Code != http.StatusOK {
						errorsChan <- assert.AnError
					}

					// Small delay to avoid overwhelming
					time.Sleep(5 * time.Millisecond)
				}
			}(batch)
		}

		wg.Wait()
		overallStressTime := time.Since(overallStart)

		close(responseTimesChan)
		close(errorsChan)

		// Collect results
		var responseTimes []time.Duration
		for responseTime := range responseTimesChan {
			responseTimes = append(responseTimes, responseTime)
		}

		var errors []error
		for err := range errorsChan {
			errors = append(errors, err)
		}

		// Verify stress test results
		assert.Empty(t, errors, "No errors should occur during stress test")
		assert.Len(t, responseTimes, stressRequests, "Should complete all stress test requests")

		// Calculate stress test statistics
		var totalTime time.Duration
		var maxTime time.Duration
		slowRequests := 0

		for _, responseTime := range responseTimes {
			totalTime += responseTime
			if responseTime > maxTime {
				maxTime = responseTime
			}
			if responseTime > 500*time.Millisecond {
				slowRequests++
			}
		}

		avgTime := totalTime / time.Duration(len(responseTimes))

		// Stress test performance assertions
		assert.True(t, avgTime < 300*time.Millisecond,
			"Average response time under stress should be < 300ms, got %v", avgTime)
		assert.True(t, maxTime < 1*time.Second,
			"Maximum response time under stress should be < 1s, got %v", maxTime)
		assert.True(t, slowRequests < stressRequests/10,
			"Less than 10%% of requests should be slow (>500ms), got %d/%d", slowRequests, stressRequests)
		assert.True(t, overallStressTime < 30*time.Second,
			"Stress test should complete in < 30 seconds, got %v", overallStressTime)

		t.Logf("Stress test results - Requests: %d, Overall: %v, Avg: %v, Max: %v, Slow: %d",
			stressRequests, overallStressTime, avgTime, maxTime, slowRequests)
	})

	t.Run("Memory and resource efficiency", func(t *testing.T) {
		router := gin.New()

		// This will fail until the enhanced health handler is implemented
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Test that repeated health checks don't cause memory leaks or resource issues
		repeatCount := 200

		for i := 0; i < repeatCount; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code,
				"Request %d should return 200 OK", i+1)

			// Verify response without keeping references to prevent memory accumulation
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Request %d should return valid JSON", i+1)

			// Periodically check that response times don't degrade
			if i%50 == 49 {
				start := time.Now()

				req, err := http.NewRequest("GET", "/health", nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				responseTime := time.Since(start)

				assert.True(t, responseTime < 200*time.Millisecond,
					"Performance should not degrade after %d requests, got %v", i+1, responseTime)
			}
		}

		t.Logf("Completed %d requests without performance degradation", repeatCount)
	})
}