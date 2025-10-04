package persistence

import (
	"errors"

	"gorm.io/gorm"

	"todo-app/application/mappers"
	"todo-app/domain/task/entities"
	"todo-app/domain/task/repositories"
	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
	"todo-app/internal/dtos"
)

// gormTaskRepository implements the TaskRepository interface using GORM
type gormTaskRepository struct {
	db     *gorm.DB
	mapper *mappers.TaskMapper
}

// NewGormTaskRepository creates a new GORM task repository
func NewGormTaskRepository(db *gorm.DB, mapper *mappers.TaskMapper) repositories.TaskRepository {
	return &gormTaskRepository{
		db:     db,
		mapper: mapper,
	}
}

// Save persists a task entity
func (r *gormTaskRepository) Save(task *entities.Task) error {
	// Convert entity to DTO using mapper
	dto := r.mapper.ToDTO(task)

	if err := r.db.Create(dto).Error; err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a task by its ID
func (r *gormTaskRepository) FindByID(id valueobjects.TaskID) (*entities.Task, error) {
	var dto dtos.Task

	if err := r.db.First(&dto, id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	// Convert DTO to entity using mapper
	return r.mapper.ToEntity(&dto)
}

// FindByUserID retrieves all tasks for a specific user
func (r *gormTaskRepository) FindByUserID(userID uservo.UserID) ([]*entities.Task, error) {
	var dtoList []dtos.Task

	if err := r.db.Where("user_id = ?", userID.Value()).Find(&dtoList).Error; err != nil {
		return nil, err
	}

	// Convert DTOs to entities using mapper
	entities := make([]*entities.Task, len(dtoList))
	for i, dto := range dtoList {
		entity, err := r.mapper.ToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// FindByUserIDAndStatus retrieves tasks by user and status
func (r *gormTaskRepository) FindByUserIDAndStatus(userID uservo.UserID, status valueobjects.TaskStatus) ([]*entities.Task, error) {
	var dtoList []dtos.Task

	// Map status to completed boolean for DTO query
	completed := status.IsCompleted()

	if err := r.db.Where("user_id = ? AND completed = ?", userID.Value(), completed).Find(&dtoList).Error; err != nil {
		return nil, err
	}

	// Convert DTOs to entities using mapper
	entities := make([]*entities.Task, len(dtoList))
	for i, dto := range dtoList {
		entity, err := r.mapper.ToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// FindByUserIDAndPriority retrieves tasks by user and priority
func (r *gormTaskRepository) FindByUserIDAndPriority(userID uservo.UserID, priority valueobjects.TaskPriority) ([]*entities.Task, error) {
	var dtoList []dtos.Task

	// Note: Priority is not stored in DTO, so this query will return all tasks for the user
	// and we filter by priority in memory (not ideal, but maintains compatibility)
	if err := r.db.Where("user_id = ?", userID.Value()).Find(&dtoList).Error; err != nil {
		return nil, err
	}

	// Convert DTOs to entities and filter by priority
	var entities []*entities.Task
	for _, dto := range dtoList {
		entity, err := r.mapper.ToEntity(&dto)
		if err != nil {
			return nil, err
		}
		// Filter by priority (all tasks have medium priority from mapper)
		if entity.Priority().Value() == priority.Value() {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

// Update updates an existing task
func (r *gormTaskRepository) Update(task *entities.Task) error {
	// Convert entity to DTO using mapper
	dto := r.mapper.ToDTO(task)

	// Update specific fields
	result := r.db.Model(&dtos.Task{}).Where("id = ?", dto.ID).Updates(map[string]interface{}{
		"title":     dto.Title,
		"completed": dto.Completed,
		"user_id":   dto.UserID,
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
	result := r.db.Delete(&dtos.Task{}, id.Value())

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

	if err := r.db.Model(&dtos.Task{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}