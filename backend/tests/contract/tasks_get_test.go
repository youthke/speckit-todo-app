package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/handlers"
)

func TestGetTasksContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers are implemented
	taskHandler := &handlers.TaskHandler{}
	router.GET("/api/v1/tasks", taskHandler.GetTasks)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Get all tasks - empty response",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Get completed tasks only",
			queryParams:    "?completed=true",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
		{
			name:           "Get pending tasks only",
			queryParams:    "?completed=false",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"tasks", "count"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/tasks"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify required fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Verify tasks array structure
			if tasks, ok := response["tasks"].([]interface{}); ok {
				for _, task := range tasks {
					taskMap := task.(map[string]interface{})
					requiredFields := []string{"id", "title", "completed", "created_at", "updated_at"}
					for _, field := range requiredFields {
						assert.Contains(t, taskMap, field, "Task should contain field: %s", field)
					}
				}
			}
		})
	}
}