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

func TestPostTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers are implemented
	taskHandler := &handlers.TaskHandler{}
	router.POST("/api/v1/tasks", taskHandler.CreateTask)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		shouldHaveID   bool
	}{
		{
			name: "Create valid task",
			requestBody: map[string]interface{}{
				"title": "Buy groceries",
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
			name: "Create task without title - should fail",
			requestBody: map[string]interface{}{
				"description": "Missing title field",
			},
			expectedStatus: http.StatusBadRequest,
			shouldHaveID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify required fields for successful creation
				requiredFields := []string{"id", "title", "completed", "created_at", "updated_at"}
				for _, field := range requiredFields {
					assert.Contains(t, response, field, "Response should contain field: %s", field)
				}

				// Verify defaults
				assert.Equal(t, tt.requestBody["title"], response["title"])
				assert.Equal(t, false, response["completed"])
				assert.NotZero(t, response["id"])
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