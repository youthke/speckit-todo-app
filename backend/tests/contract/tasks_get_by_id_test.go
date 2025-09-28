package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTaskByIDContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until DDD presentation layer handlers are implemented
	// Will need to update to use: todo-app/presentation/http.TaskHandlers
	// router.GET("/api/v1/tasks/:id", dddTaskHandler.GetTaskByID)

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "Get existing task by ID - DDD contract",
			taskID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get non-existent task - should return 404",
			taskID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get task with invalid ID format - should return 400",
			taskID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/tasks/%s", tt.taskID), nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until DDD handlers are implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status code %d, got %d", tt.expectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Verify DDD task structure according to API contract
				dddRequiredFields := []string{
					"id", "title", "description", "status", "priority",
					"user_id", "created_at", "updated_at",
				}

				for _, field := range dddRequiredFields {
					assert.Contains(t, response, field, "DDD Task should contain field: %s", field)
				}

				// Validate enum values
				if status, ok := response["status"].(string); ok {
					validStatuses := []string{"pending", "completed", "archived"}
					assert.Contains(t, validStatuses, status, "Status should be one of: %v", validStatuses)
				}

				if priority, ok := response["priority"].(string); ok {
					validPriorities := []string{"low", "medium", "high"}
					assert.Contains(t, validPriorities, priority, "Priority should be one of: %v", validPriorities)
				}

			} else if w.Code == http.StatusNotFound || w.Code == http.StatusBadRequest {
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