package valueobjects

import (
	"fmt"
)

// TaskPriority represents the priority level of a task
type TaskPriority struct {
	value string
}

// Valid task priorities
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// NewTaskPriority creates a new TaskPriority with validation
func NewTaskPriority(priority string) (TaskPriority, error) {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return TaskPriority{value: priority}, nil
	default:
		return TaskPriority{}, fmt.Errorf("invalid task priority: %s, must be one of: %s, %s, %s",
			priority, PriorityLow, PriorityMedium, PriorityHigh)
	}
}

// NewMediumPriority creates a new medium priority (default)
func NewMediumPriority() TaskPriority {
	return TaskPriority{value: PriorityMedium}
}

// NewLowPriority creates a new low priority
func NewLowPriority() TaskPriority {
	return TaskPriority{value: PriorityLow}
}

// NewHighPriority creates a new high priority
func NewHighPriority() TaskPriority {
	return TaskPriority{value: PriorityHigh}
}

// Value returns the underlying priority value
func (t TaskPriority) Value() string {
	return t.value
}

// Equals checks if two TaskPriorities are equal
func (t TaskPriority) Equals(other TaskPriority) bool {
	return t.value == other.value
}

// String returns the string representation of the TaskPriority
func (t TaskPriority) String() string {
	return t.value
}

// IsLow checks if the priority is low
func (t TaskPriority) IsLow() bool {
	return t.value == PriorityLow
}

// IsMedium checks if the priority is medium
func (t TaskPriority) IsMedium() bool {
	return t.value == PriorityMedium
}

// IsHigh checks if the priority is high
func (t TaskPriority) IsHigh() bool {
	return t.value == PriorityHigh
}

// NumericValue returns a numeric representation for comparison
func (t TaskPriority) NumericValue() int {
	switch t.value {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	default:
		return 0
	}
}