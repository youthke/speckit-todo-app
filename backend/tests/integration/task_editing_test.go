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

// TestTaskEditingScenario tests editing task titles
func TestTaskEditingScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers and storage are implemented
	taskHandler := &handlers.TaskHandler{}
	router.POST("/api/v1/tasks", taskHandler.CreateTask)
	router.PUT("/api/v1/tasks/:id", taskHandler.UpdateTask)
	router.GET("/api/v1/tasks/:id", taskHandler.GetTask)

	t.Run("Complete task editing scenario", func(t *testing.T) {
		// Step 1: Create a task
		taskData := map[string]interface{}{
			"title": "Buy groceries",
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
		originalTitle := createdTask["title"]
		originalUpdatedAt := createdTask["updated_at"]

		// Step 2: Edit the task title
		updateData := map[string]interface{}{
			"title": "Buy groceries and cook dinner",
		}
		bodyBytes, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/v1/tasks/"+taskID.(string), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &updatedTask)

		// Verify title was updated
		assert.Equal(t, "Buy groceries and cook dinner", updatedTask["title"])
		assert.NotEqual(t, originalTitle, updatedTask["title"])

		// Verify other properties remain unchanged
		assert.Equal(t, createdTask["id"], updatedTask["id"])
		assert.Equal(t, createdTask["completed"], updatedTask["completed"])
		assert.Equal(t, createdTask["created_at"], updatedTask["created_at"])

		// Verify updated_at timestamp changed
		assert.NotEqual(t, originalUpdatedAt, updatedTask["updated_at"])

		// Step 3: Edit title and completion status together
		updateData = map[string]interface{}{
			"title":     "Buy groceries, cook dinner, and clean",
			"completed": true,
		}
		bodyBytes, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/v1/tasks/"+taskID.(string), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var finalTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &finalTask)

		// Verify both updates were applied
		assert.Equal(t, "Buy groceries, cook dinner, and clean", finalTask["title"])
		assert.Equal(t, true, finalTask["completed"])

		// Step 4: Verify changes persist by fetching the task
		req, _ = http.NewRequest("GET", "/api/v1/tasks/"+taskID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &fetchedTask)

		// Verify persistence
		assert.Equal(t, finalTask["title"], fetchedTask["title"])
		assert.Equal(t, finalTask["completed"], fetchedTask["completed"])
		assert.Equal(t, finalTask["id"], fetchedTask["id"])

		// Step 5: Try to edit with invalid data (empty title)
		updateData = map[string]interface{}{
			"title": "",
		}
		bodyBytes, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/v1/tasks/"+taskID.(string), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should fail with bad request
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Verify original data was not changed
		req, _ = http.NewRequest("GET", "/api/v1/tasks/"+taskID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &fetchedTask)
		assert.Equal(t, "Buy groceries, cook dinner, and clean", fetchedTask["title"])
	})
}