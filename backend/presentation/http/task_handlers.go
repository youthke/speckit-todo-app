package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"todo-app/application/task"
)

// TaskResponse represents the HTTP response format for a task
type TaskResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskListResponse represents the HTTP response format for task lists
type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Count int            `json:"count"`
}

// CreateTaskRequest represents the HTTP request format for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,max=500"`
	Description string `json:"description" binding:"max=2000"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high"`
}

// UpdateTaskRequest represents the HTTP request format for updating a task
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,max=500"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=2000"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=pending completed archived"`
	Priority    *string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
}

// ErrorResponse represents the HTTP error response format
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// TaskHandlers contains HTTP handlers for task-related endpoints
type TaskHandlers struct {
	taskService task.TaskApplicationService
}

// NewTaskHandlers creates a new task handlers instance
func NewTaskHandlers(taskService task.TaskApplicationService) *TaskHandlers {
	return &TaskHandlers{
		taskService: taskService,
	}
}

// RegisterRoutes registers all task-related routes
func (h *TaskHandlers) RegisterRoutes(router *gin.RouterGroup) {
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", h.GetTasks)
		taskRoutes.POST("", h.CreateTask)
		taskRoutes.GET("/:id", h.GetTask)
		taskRoutes.PUT("/:id", h.UpdateTask)
		taskRoutes.DELETE("/:id", h.DeleteTask)
	}
}

// GetTasks handles GET /api/v1/tasks
func (h *TaskHandlers) GetTasks(c *gin.Context) {
	// Get user ID from context (would be set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Build query from request parameters
	query := task.TaskQuery{
		UserID: userIDUint,
	}

	// Parse optional status filter
	if statusParam := c.Query("status"); statusParam != "" {
		query.Status = &statusParam
	}

	// Parse optional priority filter
	if priorityParam := c.Query("priority"); priorityParam != "" {
		query.Priority = &priorityParam
	}

	// Get tasks from application service
	tasks, err := h.taskService.GetUserTasks(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_query",
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	response := TaskListResponse{
		Tasks: h.convertTasksToResponse(tasks),
		Count: len(tasks),
	}

	c.JSON(http.StatusOK, response)
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandlers) CreateTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Create command
	cmd := task.CreateTaskCommand{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		UserID:      userIDUint,
	}

	// Create task using application service
	createdTask, err := h.taskService.CreateTask(cmd)
	if err != nil {
		// Determine appropriate HTTP status based on error
		if isValidationError(err) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "creation_failed",
				Message: "Failed to create task",
				Details: err.Error(),
			})
		}
		return
	}

	// Convert to response format
	response := h.convertTaskToResponse(createdTask)
	c.JSON(http.StatusCreated, response)
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandlers) GetTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse task ID from path
	taskIDParam := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid task ID format",
		})
		return
	}

	// Get task from application service
	taskEntity, err := h.taskService.GetTask(uint(taskID), userIDUint)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else if isAccessDeniedError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{ // Return 404 instead of 403 for security
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "retrieval_failed",
				Message: "Failed to retrieve task",
			})
		}
		return
	}

	// Convert to response format
	response := h.convertTaskToResponse(taskEntity)
	c.JSON(http.StatusOK, response)
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandlers) UpdateTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse task ID from path
	taskIDParam := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid task ID format",
		})
		return
	}

	// Parse request body
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Create command
	cmd := task.UpdateTaskCommand{
		TaskID:      uint(taskID),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		UserID:      userIDUint,
	}

	// Update task using application service
	updatedTask, err := h.taskService.UpdateTask(cmd)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else if isAccessDeniedError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{ // Return 404 instead of 403 for security
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else if isValidationError(err) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update task",
				Details: err.Error(),
			})
		}
		return
	}

	// Convert to response format
	response := h.convertTaskToResponse(updatedTask)
	c.JSON(http.StatusOK, response)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandlers) DeleteTask(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse task ID from path
	taskIDParam := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid task ID format",
		})
		return
	}

	// Delete task using application service
	err = h.taskService.DeleteTask(uint(taskID), userIDUint)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else if isAccessDeniedError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{ // Return 404 instead of 403 for security
				Error:   "task_not_found",
				Message: "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "deletion_failed",
				Message: "Failed to delete task",
			})
		}
		return
	}

	// Return 204 No Content for successful deletion
	c.Status(http.StatusNoContent)
}

// Helper functions

// convertTaskToResponse converts a domain task entity to HTTP response format
func (h *TaskHandlers) convertTaskToResponse(taskEntity interface{}) TaskResponse {
	// This would use proper type assertion in a real implementation
	// For now, we'll assume the conversion logic
	return TaskResponse{
		// Conversion logic would be implemented here
		// This is a placeholder
	}
}

// convertTasksToResponse converts multiple domain task entities to HTTP response format
func (h *TaskHandlers) convertTasksToResponse(tasks interface{}) []TaskResponse {
	// This would use proper type assertion and conversion in a real implementation
	// For now, we'll return an empty slice
	return []TaskResponse{}
}

// Error checking helper functions
func isValidationError(err error) bool {
	// Check if error is a domain validation error
	// This would be implemented based on your error handling strategy
	return false
}

func isNotFoundError(err error) bool {
	// Check if error indicates resource not found
	// This would be implemented based on your error handling strategy
	return false
}

func isAccessDeniedError(err error) bool {
	// Check if error indicates access denied
	// This would be implemented based on your error handling strategy
	return false
}