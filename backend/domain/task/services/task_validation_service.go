package services

import (
	"errors"

	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
)

// TaskValidationService provides domain validation logic for tasks
type TaskValidationService interface {
	// ValidateTaskCreation validates task creation rules
	ValidateTaskCreation(title valueobjects.TaskTitle, userID uservo.UserID) error

	// ValidateTaskUpdate validates task update rules
	ValidateTaskUpdate(currentStatus valueobjects.TaskStatus, updates TaskUpdates) error
}

// TaskUpdates represents the fields that can be updated on a task
type TaskUpdates struct {
	Title       *valueobjects.TaskTitle
	Description *valueobjects.TaskDescription
	Status      *valueobjects.TaskStatus
	Priority    *valueobjects.TaskPriority
}

// taskValidationService implements TaskValidationService
type taskValidationService struct{}

// NewTaskValidationService creates a new task validation service
func NewTaskValidationService() TaskValidationService {
	return &taskValidationService{}
}

// ValidateTaskCreation validates task creation business rules
func (s *taskValidationService) ValidateTaskCreation(title valueobjects.TaskTitle, userID uservo.UserID) error {
	if title.IsEmpty() {
		return errors.New("task title cannot be empty")
	}

	if userID.IsZero() {
		return errors.New("user ID is required for task creation")
	}

	return nil
}

// ValidateTaskUpdate validates task update business rules
func (s *taskValidationService) ValidateTaskUpdate(currentStatus valueobjects.TaskStatus, updates TaskUpdates) error {
	// Cannot modify archived tasks (except to unarchive)
	if currentStatus.IsArchived() {
		if updates.Status == nil || !updates.Status.IsPending() {
			return errors.New("archived tasks can only be updated to pending status")
		}

		// For archived tasks being unarchived, only status change is allowed
		if updates.Title != nil || updates.Description != nil || updates.Priority != nil {
			return errors.New("archived tasks can only have status updated to pending")
		}
	}

	// Priority can only be changed on pending tasks
	if updates.Priority != nil && !currentStatus.IsPending() {
		return errors.New("task priority can only be changed on pending tasks")
	}

	// Title cannot be empty if provided
	if updates.Title != nil && updates.Title.IsEmpty() {
		return errors.New("task title cannot be empty")
	}

	return nil
}