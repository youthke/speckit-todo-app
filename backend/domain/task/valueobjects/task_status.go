package valueobjects

import (
	"errors"
	"fmt"
)

// TaskStatus represents the status of a task
type TaskStatus struct {
	value string
}

// Valid task statuses
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusArchived  = "archived"
)

// NewTaskStatus creates a new TaskStatus with validation
func NewTaskStatus(status string) (TaskStatus, error) {
	switch status {
	case StatusPending, StatusCompleted, StatusArchived:
		return TaskStatus{value: status}, nil
	default:
		return TaskStatus{}, fmt.Errorf("invalid task status: %s, must be one of: %s, %s, %s",
			status, StatusPending, StatusCompleted, StatusArchived)
	}
}

// NewPendingStatus creates a new pending status
func NewPendingStatus() TaskStatus {
	return TaskStatus{value: StatusPending}
}

// NewCompletedStatus creates a new completed status
func NewCompletedStatus() TaskStatus {
	return TaskStatus{value: StatusCompleted}
}

// NewArchivedStatus creates a new archived status
func NewArchivedStatus() TaskStatus {
	return TaskStatus{value: StatusArchived}
}

// Value returns the underlying status value
func (t TaskStatus) Value() string {
	return t.value
}

// Equals checks if two TaskStatuses are equal
func (t TaskStatus) Equals(other TaskStatus) bool {
	return t.value == other.value
}

// String returns the string representation of the TaskStatus
func (t TaskStatus) String() string {
	return t.value
}

// IsPending checks if the status is pending
func (t TaskStatus) IsPending() bool {
	return t.value == StatusPending
}

// IsCompleted checks if the status is completed
func (t TaskStatus) IsCompleted() bool {
	return t.value == StatusCompleted
}

// IsArchived checks if the status is archived
func (t TaskStatus) IsArchived() bool {
	return t.value == StatusArchived
}

// CanBeModified checks if a task with this status can be modified
func (t TaskStatus) CanBeModified() bool {
	return !t.IsArchived()
}

// CanChangePriority checks if priority can be changed for this status
func (t TaskStatus) CanChangePriority() bool {
	return t.IsPending() // Only pending tasks can change priority
}