package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/handlers"
)

func TestDeleteTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers are implemented
	taskHandler := &handlers.TaskHandler{}
	router.DELETE("/api/v1/tasks/:id", taskHandler.DeleteTask)

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "Delete existing task",
			taskID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent task - should fail",
			taskID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete with invalid ID - should fail",
			taskID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/api/v1/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusNoContent {
				// Successful deletion should have no content
				assert.Empty(t, w.Body.String())
			} else {
				// Error responses should have error structure
				if w.Body.Len() > 0 {
					var errorResponse map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
					assert.NoError(t, err)

					assert.Contains(t, errorResponse, "error")
					assert.Contains(t, errorResponse, "message")
				}
			}
		})
	}
}