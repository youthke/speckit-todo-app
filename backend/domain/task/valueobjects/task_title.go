package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

// TaskTitle represents a task title with validation
type TaskTitle struct {
	value string
}

// NewTaskTitle creates a new TaskTitle with validation
func NewTaskTitle(title string) (TaskTitle, error) {
	// Trim whitespace
	title = strings.TrimSpace(title)

	// Validate title
	if title == "" {
		return TaskTitle{}, errors.New("title cannot be empty")
	}

	if len(title) > 500 {
		return TaskTitle{}, fmt.Errorf("title too long: maximum 500 characters, got %d", len(title))
	}

	return TaskTitle{value: title}, nil
}

// Value returns the underlying title value
func (t TaskTitle) Value() string {
	return t.value
}

// Equals checks if two TaskTitles are equal
func (t TaskTitle) Equals(other TaskTitle) bool {
	return t.value == other.value
}

// String returns the string representation of the TaskTitle
func (t TaskTitle) String() string {
	return t.value
}

// IsEmpty checks if the title is empty
func (t TaskTitle) IsEmpty() bool {
	return t.value == ""
}