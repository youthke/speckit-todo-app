package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/services"
	"todo-app/internal/models"
)

func TestTaskServiceEdgeCases(t *testing.T) {
	// Note: These tests demonstrate edge case validation logic
	// In practice, they would require database setup for full integration

	t.Run("Empty title validation", func(t *testing.T) {
		service := services.NewTaskService()

		// Test empty title
		req := models.CreateTaskRequest{Title: ""}
		_, err := service.CreateTask(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")

		// Test whitespace-only title
		req = models.CreateTaskRequest{Title: "   "}
		_, err = service.CreateTask(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")

		// Test tab and newline only
		req = models.CreateTaskRequest{Title: "\t\n\r "}
		_, err = service.CreateTask(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")
	})

	t.Run("Long title validation", func(t *testing.T) {
		service := services.NewTaskService()

		// Test exactly 500 characters (should pass)
		longTitle := string(make([]rune, 500))
		for i := range longTitle {
			longTitle = string(append([]rune(longTitle)[:i], 'a'))
		}
		req := models.CreateTaskRequest{Title: string(make([]byte, 500))}

		// This would fail due to no database, but validation should pass
		_, err := service.CreateTask(req)
		// We expect a database error, not a validation error
		assert.NotContains(t, err.Error(), "title must be 500 characters or less")

		// Test 501 characters (should fail validation)
		req = models.CreateTaskRequest{Title: string(make([]byte, 501))}
		_, err = service.CreateTask(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title must be 500 characters or less")
	})

	t.Run("Update task with empty title", func(t *testing.T) {
		service := services.NewTaskService()

		// Test empty title update
		updateReq := models.UpdateTaskRequest{Title: stringPtr("")}
		_, err := service.UpdateTask(1, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")

		// Test whitespace-only title update
		updateReq = models.UpdateTaskRequest{Title: stringPtr("   ")}
		_, err = service.UpdateTask(1, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")
	})

	t.Run("Update task with long title", func(t *testing.T) {
		service := services.NewTaskService()

		// Test 501 characters update (should fail validation)
		updateReq := models.UpdateTaskRequest{Title: stringPtr(string(make([]byte, 501)))}
		_, err := service.UpdateTask(1, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title must be 500 characters or less")
	})

	t.Run("Special characters in title", func(t *testing.T) {
		service := services.NewTaskService()

		// Test various special characters (should all be valid)
		specialTitles := []string{
			"Task with émojis 🚀",
			"Task with Unicode: こんにちは",
			"Task with symbols: !@#$%^&*()",
			"Task with\ttabs\tand\nlines",
			"Task with 'quotes' and \"double quotes\"",
			"Task with <HTML> tags & entities",
		}

		for _, title := range specialTitles {
			req := models.CreateTaskRequest{Title: title}
			_, err := service.CreateTask(req)
			// We expect database error, not validation error
			if err != nil {
				assert.NotContains(t, err.Error(), "title cannot be empty")
				assert.NotContains(t, err.Error(), "title must be 500 characters or less")
			}
		}
	})

	t.Run("Nil pointer handling in updates", func(t *testing.T) {
		service := services.NewTaskService()

		// Test update with nil title and completion
		updateReq := models.UpdateTaskRequest{}
		_, err := service.UpdateTask(1, updateReq)
		// Should get "task not found" error, not validation error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")

		// Test update with only completion
		updateReq = models.UpdateTaskRequest{Completed: boolPtr(true)}
		_, err = service.UpdateTask(1, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")

		// Test update with only title
		updateReq = models.UpdateTaskRequest{Title: stringPtr("Valid title")}
		_, err = service.UpdateTask(1, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("Non-existent task operations", func(t *testing.T) {
		service := services.NewTaskService()

		// Test getting non-existent task
		_, err := service.GetTaskByID(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")

		// Test updating non-existent task
		updateReq := models.UpdateTaskRequest{Title: stringPtr("Updated")}
		_, err = service.UpdateTask(999, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")

		// Test deleting non-existent task
		err = service.DeleteTask(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("Boundary value testing", func(t *testing.T) {
		service := services.NewTaskService()

		// Test task ID boundaries
		testIDs := []uint{0, 1, 999999999}
		for _, id := range testIDs {
			_, err := service.GetTaskByID(id)
			assert.Error(t, err) // Should all fail with "not found"
			assert.Contains(t, err.Error(), "task not found")
		}

		// Test title length boundaries
		titleLengths := []int{1, 499, 500}
		for _, length := range titleLengths {
			title := string(make([]byte, length))
			req := models.CreateTaskRequest{Title: title}
			_, err := service.CreateTask(req)
			// These should pass validation (though fail on DB)
			if err != nil {
				assert.NotContains(t, err.Error(), "title cannot be empty")
				assert.NotContains(t, err.Error(), "title must be 500 characters or less")
			}
		}
	})
}