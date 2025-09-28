package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/handlers"
)

func TestPutTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers are implemented
	taskHandler := &handlers.TaskHandler{}
	router.PUT("/api/v1/tasks/:id", taskHandler.UpdateTask)

	tests := []struct {
		name           string
		taskID         string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Update task title",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title": "Updated task title",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Mark task as completed",
			taskID: "1",
			requestBody: map[string]interface{}{
				"completed": true,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update both title and completion",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title":     "Updated and completed",
				"completed": true,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update with empty title - should fail",
			taskID: "1",
			requestBody: map[string]interface{}{
				"title": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Update non-existent task - should fail",
			taskID: "999",
			requestBody: map[string]interface{}{
				"title": "This should fail",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Update with invalid ID - should fail",
			taskID: "invalid",
			requestBody: map[string]interface{}{
				"title": "This should fail",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/v1/tasks/"+tt.taskID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify required fields
				requiredFields := []string{"id", "title", "completed", "created_at", "updated_at"}
				for _, field := range requiredFields {
					assert.Contains(t, response, field, "Response should contain field: %s", field)
				}

				// Verify updates were applied
				if title, exists := tt.requestBody["title"]; exists {
					assert.Equal(t, title, response["title"])
				}
				if completed, exists := tt.requestBody["completed"]; exists {
					assert.Equal(t, completed, response["completed"])
				}
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)

				// Verify error response structure
				assert.Contains(t, errorResponse, "error")
				assert.Contains(t, errorResponse, "message")
			}
		})
	}
}