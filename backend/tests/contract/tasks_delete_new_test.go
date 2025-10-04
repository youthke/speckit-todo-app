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

func TestDeleteNewTaskContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until DDD presentation layer handlers are implemented
	// Will need to update to use: todo-app/presentation/http.TaskHandlers
	// router.DELETE("/api/v1/tasks/:id", dddTaskHandler.DeleteTask)

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "Delete existing task - DDD contract",
			taskID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent task - should return 404",
			taskID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete with invalid ID format - should return 400",
			taskID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/tasks/%s", tt.taskID), nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// This assertion will fail until DDD handlers are implemented
			assert.Equal(t, tt.expectedStatus, w.Code, "Expected status code %d, got %d", tt.expectedStatus, w.Code)

			if w.Code == http.StatusNoContent {
				// Successful deletion should return empty body
				assert.Empty(t, w.Body.String(), "DELETE should return empty body on success")

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