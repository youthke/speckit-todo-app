package persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"todo-app/domain/task/entities"
	"todo-app/domain/task/repositories"
	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
)

// TaskModel represents the GORM model for task persistence
type TaskModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"not null;size:500" json:"title"`
	Description string    `gorm:"size:2000" json:"description"`
	Status      string    `gorm:"not null;size:20" json:"status"`
	Priority    string    `gorm:"not null;size:10" json:"priority"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (TaskModel) TableName() string {
	return "tasks"
}

// gormTaskRepository implements the TaskRepository interface using GORM
type gormTaskRepository struct {
	db *gorm.DB
}

// NewGormTaskRepository creates a new GORM task repository
func NewGormTaskRepository(db *gorm.DB) repositories.TaskRepository {
	return &gormTaskRepository{
		db: db,
	}
}

// Save persists a task entity
func (r *gormTaskRepository) Save(task *entities.Task) error {
	model := r.entityToModel(task)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the task entity with the generated ID
	newID, err := valueobjects.NewTaskID(model.ID)
	if err != nil {
		return err
	}

	// Since task is immutable, we need to return a new task with the ID
	// In a real implementation, this would be handled differently
	// For now, we'll assume the caller handles this

	return nil
}

// FindByID retrieves a task by its ID
func (r *gormTaskRepository) FindByID(id valueobjects.TaskID) (*entities.Task, error) {
	var model TaskModel

	if err := r.db.First(&model, id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	return r.modelToEntity(model)
}

// FindByUserID retrieves all tasks for a specific user
func (r *gormTaskRepository) FindByUserID(userID uservo.UserID) ([]*entities.Task, error) {
	var models []TaskModel

	if err := r.db.Where("user_id = ?", userID.Value()).Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models)
}

// FindByUserIDAndStatus retrieves tasks by user and status
func (r *gormTaskRepository) FindByUserIDAndStatus(userID uservo.UserID, status valueobjects.TaskStatus) ([]*entities.Task, error) {
	var models []TaskModel

	if err := r.db.Where("user_id = ? AND status = ?", userID.Value(), status.Value()).Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models)
}

// FindByUserIDAndPriority retrieves tasks by user and priority
func (r *gormTaskRepository) FindByUserIDAndPriority(userID uservo.UserID, priority valueobjects.TaskPriority) ([]*entities.Task, error) {
	var models []TaskModel

	if err := r.db.Where("user_id = ? AND priority = ?", userID.Value(), priority.Value()).Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models)
}

// Update updates an existing task
func (r *gormTaskRepository) Update(task *entities.Task) error {
	model := r.entityToModel(task)

	// Update specific fields to avoid overwriting timestamps incorrectly
	result := r.db.Model(&model).Where("id = ?", model.ID).Updates(map[string]interface{}{
		"title":       model.Title,
		"description": model.Description,
		"status":      model.Status,
		"priority":    model.Priority,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found or no changes made")
	}

	return nil
}

// Delete removes a task by ID
func (r *gormTaskRepository) Delete(id valueobjects.TaskID) error {
	result := r.db.Delete(&TaskModel{}, id.Value())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// ExistsByID checks if a task exists by ID
func (r *gormTaskRepository) ExistsByID(id valueobjects.TaskID) (bool, error) {
	var count int64

	if err := r.db.Model(&TaskModel{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// entityToModel converts a domain entity to a GORM model
func (r *gormTaskRepository) entityToModel(task *entities.Task) TaskModel {
	return TaskModel{
		ID:          task.ID().Value(),
		Title:       task.Title().Value(),
		Description: task.Description().Value(),
		Status:      task.Status().Value(),
		Priority:    task.Priority().Value(),
		UserID:      task.UserID().Value(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}
}

// modelToEntity converts a GORM model to a domain entity
func (r *gormTaskRepository) modelToEntity(model TaskModel) (*entities.Task, error) {
	// Create value objects
	id, err := valueobjects.NewTaskID(model.ID)
	if err != nil {
		return nil, err
	}

	title, err := valueobjects.NewTaskTitle(model.Title)
	if err != nil {
		return nil, err
	}

	description, err := valueobjects.NewTaskDescription(model.Description)
	if err != nil {
		return nil, err
	}

	status, err := valueobjects.NewTaskStatus(model.Status)
	if err != nil {
		return nil, err
	}

	priority, err := valueobjects.NewTaskPriority(model.Priority)
	if err != nil {
		return nil, err
	}

	userID, err := uservo.NewUserID(model.UserID)
	if err != nil {
		return nil, err
	}

	// Create task entity
	task, err := entities.NewTask(id, title, description, status, priority, userID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// modelsToEntities converts multiple GORM models to domain entities
func (r *gormTaskRepository) modelsToEntities(models []TaskModel) ([]*entities.Task, error) {
	entities := make([]*entities.Task, len(models))

	for i, model := range models {
		entity, err := r.modelToEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}