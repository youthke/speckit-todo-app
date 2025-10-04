package mappers

import (
	"fmt"

	"domain/task/entities"
	"domain/task/valueobjects"
	uservo "domain/user/valueobjects"
	"todo-app/internal/dtos"
)

// TaskMapper handles conversion between Task DTOs and Task entities
type TaskMapper struct{}

// ToEntity converts a TaskDTO to a Task entity
func (m *TaskMapper) ToEntity(dto *dtos.Task) (*entities.Task, error) {
	// Validate and create TaskID
	taskID := valueobjects.NewTaskID(dto.ID)
	if taskID.IsZero() && dto.ID != 0 {
		return nil, fmt.Errorf("task ID cannot be zero")
	}

	// Validate and create TaskTitle
	title, err := valueobjects.NewTaskTitle(dto.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	// Create empty TaskDescription (not in DTO)
	description, err := valueobjects.NewTaskDescription("")
	if err != nil {
		return nil, fmt.Errorf("failed to create description: %w", err)
	}

	// Convert completed boolean to TaskStatus
	var status valueobjects.TaskStatus
	if dto.Completed {
		status = valueobjects.NewCompletedStatus()
	} else {
		status = valueobjects.NewPendingStatus()
	}

	// Create default TaskPriority (medium)
	priority := valueobjects.NewMediumPriority()

	// Create UserID value object from DTO
	ownerID := uservo.NewUserID(dto.UserID)
	if ownerID.IsZero() {
		return nil, fmt.Errorf("user ID cannot be zero")
	}

	// Create the Task entity
	task, err := entities.NewTask(taskID, title, description, status, priority, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create task entity: %w", err)
	}

	return task, nil
}

// ToDTO converts a Task entity to a TaskDTO
func (m *TaskMapper) ToDTO(entity *entities.Task) *dtos.Task {
	return &dtos.Task{
		ID:        entity.ID().Value(),
		Title:     entity.Title().Value(),
		Completed: entity.Status().IsCompleted(), // Convert TaskStatus to boolean
		UserID:    entity.UserID().Value(),       // Include UserID for database
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
	}
}
