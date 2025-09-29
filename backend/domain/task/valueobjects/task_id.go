package valueobjects

import "fmt"

// TaskID represents a unique identifier for a task
type TaskID struct {
	value uint
}

// NewTaskID creates a new TaskID
func NewTaskID(id uint) TaskID {
	return TaskID{value: id}
}

// Value returns the underlying ID value
func (t TaskID) Value() uint {
	return t.value
}

// Equals checks if two TaskIDs are equal
func (t TaskID) Equals(other TaskID) bool {
	return t.value == other.value
}

// String returns the string representation of the TaskID
func (t TaskID) String() string {
	return fmt.Sprintf("TaskID(%d)", t.value)
}

// IsZero checks if the TaskID is zero (uninitialized)
func (t TaskID) IsZero() bool {
	return t.value == 0
}