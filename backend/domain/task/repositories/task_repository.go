package repositories

import (
	"domain/task/entities"
	"domain/task/valueobjects"
	uservo "domain/user/valueobjects"
)

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	// Save persists a task entity
	Save(task *entities.Task) error

	// FindByID retrieves a task by its ID
	FindByID(id valueobjects.TaskID) (*entities.Task, error)

	// FindByUserID retrieves all tasks for a specific user
	FindByUserID(userID uservo.UserID) ([]*entities.Task, error)

	// FindByUserIDAndStatus retrieves tasks by user and status
	FindByUserIDAndStatus(userID uservo.UserID, status valueobjects.TaskStatus) ([]*entities.Task, error)

	// FindByUserIDAndPriority retrieves tasks by user and priority
	FindByUserIDAndPriority(userID uservo.UserID, priority valueobjects.TaskPriority) ([]*entities.Task, error)

	// Update updates an existing task
	Update(task *entities.Task) error

	// Delete removes a task by ID
	Delete(id valueobjects.TaskID) error

	// ExistsByID checks if a task exists by ID
	ExistsByID(id valueobjects.TaskID) (bool, error)
}