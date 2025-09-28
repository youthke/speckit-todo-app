package integration

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

// TestTaskCompletionScenario tests marking tasks as complete/incomplete
func TestTaskCompletionScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers and storage are implemented
	taskHandler := &handlers.TaskHandler{}
	router.POST("/api/v1/tasks", taskHandler.CreateTask)
	router.PUT("/api/v1/tasks/:id", taskHandler.UpdateTask)
	router.GET("/api/v1/tasks", taskHandler.GetTasks)

	t.Run("Complete task completion scenario", func(t *testing.T) {
		// Step 1: Create a task
		taskData := map[string]interface{}{
			"title": "Complete project",
		}
		bodyBytes, _ := json.Marshal(taskData)
		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var createdTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createdTask)
		taskID := createdTask["id"]

		// Verify initial state is not completed
		assert.Equal(t, false, createdTask["completed"])

		// Step 2: Mark task as completed
		updateData := map[string]interface{}{
			"completed": true,
		}
		bodyBytes, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/v1/tasks/"+taskID.(string), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &updatedTask)

		// Verify completion status changed
		assert.Equal(t, true, updatedTask["completed"])
		assert.Equal(t, createdTask["title"], updatedTask["title"])
		assert.Equal(t, createdTask["id"], updatedTask["id"])

		// Verify updated_at timestamp changed
		assert.NotEqual(t, createdTask["updated_at"], updatedTask["updated_at"])

		// Step 3: Verify completed task shows in filtered list
		req, _ = http.NewRequest("GET", "/api/v1/tasks?completed=true", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var completedResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &completedResponse)
		assert.Equal(t, float64(1), completedResponse["count"])

		// Step 4: Verify task does not show in pending list
		req, _ = http.NewRequest("GET", "/api/v1/tasks?completed=false", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var pendingResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &pendingResponse)
		assert.Equal(t, float64(0), pendingResponse["count"])

		// Step 5: Mark task as incomplete again
		updateData = map[string]interface{}{
			"completed": false,
		}
		bodyBytes, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/v1/tasks/"+taskID.(string), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var revertedTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &revertedTask)

		// Verify task is now incomplete
		assert.Equal(t, false, revertedTask["completed"])

		// Step 6: Verify task shows in pending list again
		req, _ = http.NewRequest("GET", "/api/v1/tasks?completed=false", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &pendingResponse)
		assert.Equal(t, float64(1), pendingResponse["count"])
	})
}