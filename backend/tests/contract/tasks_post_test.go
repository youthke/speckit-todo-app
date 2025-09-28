package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until DDD presentation layer handlers are implemented
	// Will need to update to use: todo-app/presentation/http.TaskHandlers
	// router.POST("/api/v1/tasks", dddTaskHandler.CreateTask)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldHaveID   bool
	}{
		{
			name: "Create valid task - DDD contract",
			requestBody: map[string]interface{}{
				"title":       "Buy groceries",
				"description": "Get milk, bread, and eggs from the store",
				"priority":    "medium",
			},
			expectedStatus: http.StatusCreated,
			shouldHaveID:   true,
		},
		{
			name: "Create task with high priority - DDD contract",
			requestBody: map[string]interface{}{
				"title":    "Fix critical bug",
				"priority": "high",
			},
			expectedStatus: http.StatusCreated,
			shouldHaveID:   true,
		},
		{
			name: "Create task with minimal fields - DDD contract",
			requestBody: map[string]interface{}{
				"title": "Simple task",
			},
			expectedStatus: http.StatusCreated,
			shouldHaveID:   true,
		},
		{
			name: "Create task with empty title - should fail",
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "Create task with long title - should fail",
			requestBody: map[string]interface{}{
				"title": string(make([]byte, 501)),
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "Create task with invalid priority - should fail",
			requestBody: map[string]interface{}{
				"title":    "Valid title",
				"priority": "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
		{
			name: "Create task without title - should fail",
			requestBody: map[string]interface{}{
				"description": "Missing title field",
				"priority":    "low",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until DDD handlers are implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status code %d, got %d", tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Verify DDD task creation response fields
				dddRequiredFields := []string{
					"id",          // Task ID
					"title",       // Task title
					"description", // Task description
					"status",      // Should default to "pending"
					"priority",    // Task priority
					"user_id",     // User who owns the task
					"created_at",  // Creation timestamp
					"updated_at",  // Last update timestamp
				}

				for _, field := range dddRequiredFields {
					assert.Contains(t, response, field, "DDD response should contain field: %s", field)
				}

				// Verify request data mapping
				assert.Equal(t, tt.requestBody["title"], response["title"])
				if desc, exists := tt.requestBody["description"]; exists {
					assert.Equal(t, desc, response["description"])
				}
				if priority, exists := tt.requestBody["priority"]; exists {
					assert.Equal(t, priority, response["priority"])
				} else {
					assert.Equal(t, "medium", response["priority"], "Priority should default to medium")
				}

				// Verify DDD defaults
				assert.Equal(t, "pending", response["status"], "New task should have pending status")
				assert.NotZero(t, response["id"])
				assert.NotZero(t, response["user_id"], "Task should be assigned to a user")

			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)

				// Verify DDD error response structure
				assert.Contains(t, errorResponse, "error")
				assert.Contains(t, errorResponse, "message")
			}
		})
	}
}