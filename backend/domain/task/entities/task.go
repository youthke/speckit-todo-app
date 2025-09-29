package entities

import (
	"errors"
	"time"

	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
)

// Task represents a domain entity for task management
type Task struct {
	id          valueobjects.TaskID
	title       valueobjects.TaskTitle
	description valueobjects.TaskDescription
	status      valueobjects.TaskStatus
	priority    valueobjects.TaskPriority
	userID      uservo.UserID
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTask creates a new Task entity
func NewTask(
	id valueobjects.TaskID,
	title valueobjects.TaskTitle,
	description valueobjects.TaskDescription,
	status valueobjects.TaskStatus,
	priority valueobjects.TaskPriority,
	userID uservo.UserID,
) (*Task, error) {
	if id.IsZero() {
		return nil, errors.New("task ID cannot be zero")
	}

	if userID.IsZero() {
		return nil, errors.New("user ID cannot be zero")
	}

	now := time.Now()

	return &Task{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		userID:      userID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// MarkAsCompleted marks the task as completed
func (t *Task) MarkAsCompleted() error {
	if t.status.IsArchived() {
		return errors.New("cannot complete archived task")
	}

	t.status = valueobjects.NewCompletedStatus()
	t.updatedAt = time.Now()
	return nil
}

// UpdateTitle updates the task title
func (t *Task) UpdateTitle(title valueobjects.TaskTitle) error {
	if !t.status.CanBeModified() {
		return errors.New("cannot modify archived task")
	}

	t.title = title
	t.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the task description
func (t *Task) UpdateDescription(description valueobjects.TaskDescription) error {
	if !t.status.CanBeModified() {
		return errors.New("cannot modify archived task")
	}

	t.description = description
	t.updatedAt = time.Now()
	return nil
}

// ChangePriority changes the task priority
func (t *Task) ChangePriority(priority valueobjects.TaskPriority) error {
	if !t.status.CanChangePriority() {
		return errors.New("can only change priority on pending tasks")
	}

	t.priority = priority
	t.updatedAt = time.Now()
	return nil
}

// Archive archives the task
func (t *Task) Archive() error {
	t.status = valueobjects.NewArchivedStatus()
	t.updatedAt = time.Now()
	return nil
}

// IsOwnedBy checks if the task is owned by the given user
func (t *Task) IsOwnedBy(userID uservo.UserID) bool {
	return t.userID.Equals(userID)
}

// Getters

// ID returns the task ID
func (t *Task) ID() valueobjects.TaskID {
	return t.id
}

// Title returns the task title
func (t *Task) Title() valueobjects.TaskTitle {
	return t.title
}

// Description returns the task description
func (t *Task) Description() valueobjects.TaskDescription {
	return t.description
}

// Status returns the task status
func (t *Task) Status() valueobjects.TaskStatus {
	return t.status
}

// Priority returns the task priority
func (t *Task) Priority() valueobjects.TaskPriority {
	return t.priority
}

// UserID returns the user ID that owns this task
func (t *Task) UserID() uservo.UserID {
	return t.userID
}

// CreatedAt returns the creation time
func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update time
func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}