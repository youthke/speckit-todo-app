package services

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"todo-app/internal/dtos"
	"todo-app/internal/storage"
)

// TaskService handles business logic for tasks
type TaskService struct {
	db *gorm.DB
}

// NewTaskService creates a new TaskService instance
func NewTaskService() *TaskService {
	return &TaskService{
		db: storage.GetDB(),
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(req dtos.CreateTaskRequest) (*dtos.Task, error) {
	// Trim whitespace from title
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if len(title) > 500 {
		return nil, errors.New("title must be 500 characters or less")
	}

	task := dtos.Task{
		Title:     title,
		Completed: false,
	}

	result := s.db.Create(&task)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create task: %w", result.Error)
	}

	return &task, nil
}

// GetTasks retrieves tasks with optional filtering
func (s *TaskService) GetTasks(completed *bool) ([]dtos.Task, error) {
	var tasks []dtos.Task
	query := s.db.Order("created_at DESC")

	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", result.Error)
	}

	return tasks, nil
}

// GetTaskByID retrieves a task by its ID
func (s *TaskService) GetTaskByID(id uint) (*dtos.Task, error) {
	var task dtos.Task
	result := s.db.First(&task, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("failed to retrieve task: %w", result.Error)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(id uint, req dtos.UpdateTaskRequest) (*dtos.Task, error) {
	// First, get the existing task
	task, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, errors.New("title cannot be empty")
		}
		if len(title) > 500 {
			return nil, errors.New("title must be 500 characters or less")
		}
		updates["title"] = title
	}

	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}

	// Perform update
	result := s.db.Model(task).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update task: %w", result.Error)
	}

	// Fetch updated task
	updatedTask, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

// DeleteTask removes a task by ID
func (s *TaskService) DeleteTask(id uint) error {
	// Check if task exists
	_, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Delete the task
	result := s.db.Delete(&dtos.Task{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}

	return nil
}

// GetTaskCount returns the total number of tasks
func (s *TaskService) GetTaskCount(completed *bool) (int64, error) {
	var count int64
	query := s.db.Model(&dtos.Task{})

	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", result.Error)
	}

	return count, nil
}