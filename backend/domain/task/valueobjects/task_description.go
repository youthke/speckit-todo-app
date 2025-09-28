package valueobjects

import (
	"fmt"
	"strings"
)

// TaskDescription represents a task description with validation
type TaskDescription struct {
	value string
}

// NewTaskDescription creates a new TaskDescription with validation
func NewTaskDescription(description string) (TaskDescription, error) {
	// Trim whitespace
	description = strings.TrimSpace(description)

	// Validate description length (optional field, so empty is allowed)
	if len(description) > 2000 {
		return TaskDescription{}, fmt.Errorf("description too long: maximum 2000 characters, got %d", len(description))
	}

	return TaskDescription{value: description}, nil
}

// Value returns the underlying description value
func (t TaskDescription) Value() string {
	return t.value
}

// Equals checks if two TaskDescriptions are equal
func (t TaskDescription) Equals(other TaskDescription) bool {
	return t.value == other.value
}

// String returns the string representation of the TaskDescription
func (t TaskDescription) String() string {
	return t.value
}

// IsEmpty checks if the description is empty
func (t TaskDescription) IsEmpty() bool {
	return t.value == ""
}