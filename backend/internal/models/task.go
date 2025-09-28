package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a single TODO item
type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"type:varchar(500);not null" validate:"required,max=500"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the Task model
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate hook to validate task before creation
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	return t.Validate()
}

// BeforeUpdate hook to validate task before update
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
}

// Validate performs validation on the Task model
func (t *Task) Validate() error {
	if t.Title == "" {
		return gorm.ErrInvalidValue
	}
	if len(t.Title) > 500 {
		return gorm.ErrInvalidValue
	}
	return nil
}

// CreateTaskRequest represents the request payload for creating a task
type CreateTaskRequest struct {
	Title string `json:"title" binding:"required,max=500"`
}

// UpdateTaskRequest represents the request payload for updating a task
type UpdateTaskRequest struct {
	Title     *string `json:"title,omitempty" binding:"omitempty,max=500"`
	Completed *bool   `json:"completed,omitempty"`
}

// TaskResponse represents the response format for task operations
type TaskResponse struct {
	Tasks []Task `json:"tasks"`
	Count int    `json:"count"`
}