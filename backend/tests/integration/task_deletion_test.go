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

// TestTaskDeletionScenario tests deleting tasks
func TestTaskDeletionScenario(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// This will fail until handlers and storage are implemented
	taskHandler := &handlers.TaskHandler{}
	router.POST("/api/v1/tasks", taskHandler.CreateTask)
	router.DELETE("/api/v1/tasks/:id", taskHandler.DeleteTask)
	router.GET("/api/v1/tasks", taskHandler.GetTasks)
	router.GET("/api/v1/tasks/:id", taskHandler.GetTask)

	t.Run("Complete task deletion scenario", func(t *testing.T) {
		// Step 1: Create multiple tasks
		task1Data := map[string]interface{}{"title": "Task 1"}
		task2Data := map[string]interface{}{"title": "Task 2"}
		task3Data := map[string]interface{}{"title": "Task 3"}

		var task1ID, task2ID, task3ID interface{}

		// Create task 1
		bodyBytes, _ := json.Marshal(task1Data)
		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var task1 map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &task1)
		task1ID = task1["id"]

		// Create task 2
		bodyBytes, _ = json.Marshal(task2Data)
		req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var task2 map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &task2)
		task2ID = task2["id"]

		// Create task 3
		bodyBytes, _ = json.Marshal(task3Data)
		req, _ = http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var task3 map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &task3)
		task3ID = task3["id"]

		// Step 2: Verify all tasks exist
		req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var listResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, float64(3), listResponse["count"])

		// Step 3: Delete task 2 (middle task)
		req, _ = http.NewRequest("DELETE", "/api/v1/tasks/"+task2ID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())

		// Step 4: Verify task 2 no longer exists
		req, _ = http.NewRequest("GET", "/api/v1/tasks/"+task2ID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Step 5: Verify remaining tasks still exist and count is correct
		req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, float64(2), listResponse["count"])

		tasks := listResponse["tasks"].([]interface{})
		assert.Len(t, tasks, 2)

		// Verify tasks 1 and 3 still exist
		remainingTasks := make(map[string]bool)
		for _, task := range tasks {
			taskMap := task.(map[string]interface{})
			remainingTasks[taskMap["title"].(string)] = true
		}
		assert.True(t, remainingTasks["Task 1"])
		assert.True(t, remainingTasks["Task 3"])
		assert.False(t, remainingTasks["Task 2"])

		// Step 6: Try to delete non-existent task
		req, _ = http.NewRequest("DELETE", "/api/v1/tasks/999", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Step 7: Try to delete with invalid ID
		req, _ = http.NewRequest("DELETE", "/api/v1/tasks/invalid", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Step 8: Delete remaining tasks
		req, _ = http.NewRequest("DELETE", "/api/v1/tasks/"+task1ID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		req, _ = http.NewRequest("DELETE", "/api/v1/tasks/"+task3ID.(string), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Step 9: Verify list is empty
		req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &listResponse)
		assert.Equal(t, float64(0), listResponse["count"])
	})
}