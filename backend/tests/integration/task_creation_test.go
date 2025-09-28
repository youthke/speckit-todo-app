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

// TestTaskCreationScenario tests the complete user scenario for creating tasks
func TestTaskCreationScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers and storage are implemented
	taskHandler := &handlers.TaskHandler{}
	router.POST("/api/v1/tasks", taskHandler.CreateTask)
	router.GET("/api/v1/tasks", taskHandler.GetTasks)

	t.Run("Complete task creation scenario", func(t *testing.T) {
		// Step 1: Verify empty list initially
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var initialResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &initialResponse)
		assert.Equal(t, float64(0), initialResponse["count"])

		// Step 2: Create a new task
		taskData := map[string]interface{}{
			"title": "Buy groceries",
		}
		bodyBytes, _ := json.Marshal(taskData)
		req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var createdTask map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createdTask)

		// Verify task properties
		assert.Equal(t, "Buy groceries", createdTask["title"])
		assert.Equal(t, false, createdTask["completed"])
		assert.NotZero(t, createdTask["id"])
		assert.NotEmpty(t, createdTask["created_at"])
		assert.NotEmpty(t, createdTask["updated_at"])

		// Step 3: Verify task appears in list
		req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var listResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, float64(1), listResponse["count"])

		tasks := listResponse["tasks"].([]interface{})
		assert.Len(t, tasks, 1)

		firstTask := tasks[0].(map[string]interface{})
		assert.Equal(t, createdTask["id"], firstTask["id"])
		assert.Equal(t, "Buy groceries", firstTask["title"])
		assert.Equal(t, false, firstTask["completed"])

		// Step 4: Create second task
		taskData2 := map[string]interface{}{
			"title": "Walk the dog",
		}
		bodyBytes, _ = json.Marshal(taskData2)
		req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Step 5: Verify both tasks in list
		req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, float64(2), listResponse["count"])

		tasks = listResponse["tasks"].([]interface{})
		assert.Len(t, tasks, 2)
	})
}