package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPutTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until DDD presentation layer handlers are implemented
	// Will need to update to use: todo-app/presentation/http.TaskHandlers
	// router.PUT("/api/v1/tasks/:id", dddTaskHandler.UpdateTask)

	tests := []struct {
		name           string
		taskID         string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Update task title - DDD contract",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title": "Updated task title",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update task status to completed - DDD contract",
			taskID: "1",
			requestBody: map[string]interface{}{
				"status": "completed",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update task priority - DDD contract",
			taskID: "1",
			requestBody: map[string]interface{}{
				"priority": "high",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update multiple fields - DDD contract",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title":       "Complete project",
				"description": "Finish all remaining tasks for the project",
				"status":      "completed",
				"priority":    "high",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update with invalid status - should fail",
			taskID: "1",
			requestBody: map[string]interface{}{
				"status": "invalid_status",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Update with invalid priority - should fail",
			taskID: "1",
			requestBody: map[string]interface{}{
				"priority": "invalid_priority",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Update non-existent task - should return 404",
			taskID: "999999",
			requestBody: map[string]interface{}{
				"title": "Updated title",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Update with empty title - should fail",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/tasks/%s", tt.taskID), bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until DDD handlers are implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status code %d, got %d", tt.expectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Verify DDD task update response fields
				dddRequiredFields := []string{
					"id", "title", "description", "status", "priority",
					"user_id", "created_at", "updated_at",
				}

				for _, field := range dddRequiredFields {
					assert.Contains(t, response, field, "DDD response should contain field: %s", field)
				}

				// Verify updated fields
				for key, expectedValue := range tt.requestBody {
					assert.Equal(t, expectedValue, response[key], "Field %s should be updated", key)
				}

				// Verify updated_at was modified
				assert.NotZero(t, response["updated_at"])

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