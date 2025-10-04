package persistence

import (
	"errors"

	"gorm.io/gorm"

	"todo-app/application/mappers"
	"todo-app/domain/user/entities"
	"todo-app/domain/user/repositories"
	"todo-app/domain/user/valueobjects"
	"todo-app/internal/dtos"
)

// gormUserRepository implements the UserRepository interface using GORM
type gormUserRepository struct {
	db     *gorm.DB
	mapper *mappers.UserMapper
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB, mapper *mappers.UserMapper) repositories.UserRepository {
	return &gormUserRepository{
		db:     db,
		mapper: mapper,
	}
}

// Save persists a user entity
func (r *gormUserRepository) Save(user *entities.User) error {
	// Convert entity to DTO using mapper
	dto := r.mapper.ToDTO(user)

	if err := r.db.Create(dto).Error; err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *gormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
	var dto dtos.User

	if err := r.db.First(&dto, id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	// Convert DTO to entity using mapper
	return r.mapper.ToEntity(&dto)
}

// FindByEmail retrieves a user by their email address
func (r *gormUserRepository) FindByEmail(email valueobjects.Email) (*entities.User, error) {
	var dto dtos.User

	if err := r.db.Where("email = ?", email.Value()).First(&dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	// Convert DTO to entity using mapper
	return r.mapper.ToEntity(&dto)
}

// Update updates an existing user
func (r *gormUserRepository) Update(user *entities.User) error {
	// Convert entity to DTO using mapper
	dto := r.mapper.ToDTO(user)

	// Update specific fields
	result := r.db.Model(&dtos.User{}).Where("id = ?", dto.ID).Updates(map[string]interface{}{
		"email": dto.Email,
		"name":  dto.Name,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or no changes made")
	}

	return nil
}

// Delete removes a user by ID
func (r *gormUserRepository) Delete(id valueobjects.UserID) error {
	result := r.db.Delete(&dtos.User{}, id.Value())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ExistsByID checks if a user exists by ID
func (r *gormUserRepository) ExistsByID(id valueobjects.UserID) (bool, error) {
	var count int64

	if err := r.db.Model(&dtos.User{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEmail checks if a user exists by email address
func (r *gormUserRepository) ExistsByEmail(email valueobjects.Email) (bool, error) {
	var count int64

	if err := r.db.Model(&dtos.User{}).Where("email = ?", email.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindAll retrieves all users (for admin purposes)
func (r *gormUserRepository) FindAll() ([]*entities.User, error) {
	var dtoList []dtos.User

	if err := r.db.Find(&dtoList).Error; err != nil {
		return nil, err
	}

	// Convert DTOs to entities using mapper
	entities := make([]*entities.User, len(dtoList))
	for i, dto := range dtoList {
		entity, err := r.mapper.ToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}

// Count returns the total number of users
func (r *gormUserRepository) Count() (int64, error) {
	var count int64

	if err := r.db.Model(&dtos.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}