package persistence

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"todo-app/domain/user/entities"
	"todo-app/domain/user/repositories"
	"todo-app/domain/user/valueobjects"
	taskvo "todo-app/domain/task/valueobjects"
)

// UserModel represents the GORM model for user persistence
type UserModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email       string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FirstName   string    `gorm:"not null;size:50" json:"first_name"`
	LastName    string    `gorm:"not null;size:50" json:"last_name"`
	Timezone    string    `gorm:"not null;size:50" json:"timezone"`
	Preferences string    `gorm:"type:text" json:"preferences"` // JSON-encoded preferences
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// UserPreferencesJSON represents the JSON structure for user preferences
type UserPreferencesJSON struct {
	DefaultTaskPriority string `json:"default_task_priority"`
	EmailNotifications  bool   `json:"email_notifications"`
	ThemePreference     string `json:"theme_preference"`
}

// gormUserRepository implements the UserRepository interface using GORM
type gormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

// Save persists a user entity
func (r *gormUserRepository) Save(user *entities.User) error {
	model, err := r.entityToModel(user)
	if err != nil {
		return err
	}

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the user entity with the generated ID would require entity reconstruction
	// In a real implementation, this would be handled differently
	return nil
}

// FindByID retrieves a user by their ID
func (r *gormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
	var model UserModel

	if err := r.db.First(&model, id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	return r.modelToEntity(model)
}

// FindByEmail retrieves a user by their email address
func (r *gormUserRepository) FindByEmail(email valueobjects.Email) (*entities.User, error) {
	var model UserModel

	if err := r.db.Where("email = ?", email.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, err
	}

	return r.modelToEntity(model)
}

// Update updates an existing user
func (r *gormUserRepository) Update(user *entities.User) error {
	model, err := r.entityToModel(user)
	if err != nil {
		return err
	}

	// Update specific fields
	result := r.db.Model(&model).Where("id = ?", model.ID).Updates(map[string]interface{}{
		"email":       model.Email,
		"first_name":  model.FirstName,
		"last_name":   model.LastName,
		"timezone":    model.Timezone,
		"preferences": model.Preferences,
		"updated_at":  time.Now(),
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
	result := r.db.Delete(&UserModel{}, id.Value())

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

	if err := r.db.Model(&UserModel{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEmail checks if a user exists by email address
func (r *gormUserRepository) ExistsByEmail(email valueobjects.Email) (bool, error) {
	var count int64

	if err := r.db.Model(&UserModel{}).Where("email = ?", email.Value()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindAll retrieves all users (for admin purposes)
func (r *gormUserRepository) FindAll() ([]*entities.User, error) {
	var models []UserModel

	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(models)
}

// Count returns the total number of users
func (r *gormUserRepository) Count() (int64, error) {
	var count int64

	if err := r.db.Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// entityToModel converts a domain entity to a GORM model
func (r *gormUserRepository) entityToModel(user *entities.User) (UserModel, error) {
	// Convert preferences to JSON
	prefsJSON := UserPreferencesJSON{
		DefaultTaskPriority: user.Preferences().DefaultTaskPriority().Value(),
		EmailNotifications:  user.Preferences().EmailNotifications(),
		ThemePreference:     user.Preferences().ThemePreference(),
	}

	preferencesBytes, err := json.Marshal(prefsJSON)
	if err != nil {
		return UserModel{}, err
	}

	return UserModel{
		ID:          user.ID().Value(),
		Email:       user.Email().Value(),
		FirstName:   user.Profile().FirstName(),
		LastName:    user.Profile().LastName(),
		Timezone:    user.Profile().Timezone(),
		Preferences: string(preferencesBytes),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
	}, nil
}

// modelToEntity converts a GORM model to a domain entity
func (r *gormUserRepository) modelToEntity(model UserModel) (*entities.User, error) {
	// Create value objects
	id, err := valueobjects.NewUserID(model.ID)
	if err != nil {
		return nil, err
	}

	email, err := valueobjects.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	profile, err := valueobjects.NewUserProfile(model.FirstName, model.LastName, model.Timezone)
	if err != nil {
		return nil, err
	}

	// Parse preferences from JSON
	var prefsJSON UserPreferencesJSON
	if err := json.Unmarshal([]byte(model.Preferences), &prefsJSON); err != nil {
		return nil, err
	}

	// Create preferences value object
	defaultPriority, err := taskvo.NewTaskPriority(prefsJSON.DefaultTaskPriority)
	if err != nil {
		return nil, err
	}

	preferences, err := valueobjects.NewUserPreferences(
		defaultPriority,
		prefsJSON.EmailNotifications,
		prefsJSON.ThemePreference,
	)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user, err := entities.NewUser(id, email, profile, preferences)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// modelsToEntities converts multiple GORM models to domain entities
func (r *gormUserRepository) modelsToEntities(models []UserModel) ([]*entities.User, error) {
	entities := make([]*entities.User, len(models))

	for i, model := range models {
		entity, err := r.modelToEntity(model)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}

	return entities, nil
}