package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"todo-app/internal/dtos"
	"todo-app/internal/services"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler creates a new TaskHandler instance
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: services.NewTaskService(),
	}
}

// GetTasks handles GET /api/v1/tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	// Parse query parameters
	var completed *bool
	if completedStr := c.Query("completed"); completedStr != "" {
		if completedBool, err := strconv.ParseBool(completedStr); err == nil {
			completed = &completedBool
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": "Invalid 'completed' parameter. Must be true or false.",
			})
			return
		}
	}

	// Get tasks from service
	tasks, err := h.taskService.GetTasks(completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve tasks",
		})
		return
	}

	// Get count
	count, err := h.taskService.GetTaskCount(completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to count tasks",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dtos.TaskResponse{
		Tasks: tasks,
		Count: int(count),
	})
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	// Parse task ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid task ID",
		})
		return
	}

	// Get task from service
	task, err := h.taskService.GetTaskByID(uint(id))
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Task with ID " + idStr + " not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dtos.CreateTaskRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Create task via service
	task, err := h.taskService.CreateTask(req)
	if err != nil {
		if err.Error() == "title cannot be empty" || err.Error() == "title must be 500 characters or less" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Parse task ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid task ID",
		})
		return
	}

	var req dtos.UpdateTaskRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Update task via service
	task, err := h.taskService.UpdateTask(uint(id), req)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Task with ID " + idStr + " not found",
			})
			return
		}
		if err.Error() == "title cannot be empty" || err.Error() == "title must be 500 characters or less" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Parse task ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid task ID",
		})
		return
	}

	// Delete task via service
	err = h.taskService.DeleteTask(uint(id))
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Task with ID " + idStr + " not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to delete task",
		})
		return
	}

	c.Status(http.StatusNoContent)
}