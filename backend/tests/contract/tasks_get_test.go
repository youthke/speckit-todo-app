package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTasksContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until DDD presentation layer handlers are implemented
	// Will need to update to use: todo-app/presentation/http.TaskHandlers
	// router.GET("/api/v1/tasks", dddTaskHandler.GetTasks)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Get all tasks - DDD contract",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Filter by status pending - DDD contract",
			queryParams:    "?status=pending",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Filter by status completed - DDD contract",
			queryParams:    "?status=completed",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Filter by priority high - DDD contract",
			queryParams:    "?priority=high",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Invalid status should fail - DDD contract",
			queryParams:    "?status=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/tasks"+tt.queryParams, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until DDD handlers are implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status code %d, got %d", tt.expectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Verify required response fields exist
				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field, "Response should contain field: %s", field)
				}

				// Verify DDD task structure according to new API contract
				if tasks, ok := response["tasks"].([]interface{}); ok {
					for _, task := range tasks {
						taskMap := task.(map[string]interface{})

						// DDD API contract requires these fields (from task-api.yaml)
						dddRequiredFields := []string{
							"id",          // Task ID
							"title",       // Task title
							"description", // Task description (new in DDD)
							"status",      // "pending", "completed", "archived" (replaces "completed" boolean)
							"priority",    // "low", "medium", "high" (new in DDD)
							"user_id",     // User who owns the task (new in DDD)
							"created_at",  // Creation timestamp
							"updated_at",  // Last update timestamp
						}

						for _, field := range dddRequiredFields {
							assert.Contains(t, taskMap, field, "DDD Task should contain field: %s", field)
						}

						// Validate enum values
						if status, ok := taskMap["status"].(string); ok {
							validStatuses := []string{"pending", "completed", "archived"}
							assert.Contains(t, validStatuses, status, "Status should be one of: %v", validStatuses)
						}

						if priority, ok := taskMap["priority"].(string); ok {
							validPriorities := []string{"low", "medium", "high"}
							assert.Contains(t, validPriorities, priority, "Priority should be one of: %v", validPriorities)
						}
					}
				}
			}
		})
	}
}