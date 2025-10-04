package repositories

import (
	"domain/user/entities"
	"domain/user/valueobjects"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Save persists a user entity
	Save(user *entities.User) error

	// FindByID retrieves a user by their ID
	FindByID(id valueobjects.UserID) (*entities.User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(email valueobjects.Email) (*entities.User, error)

	// Update updates an existing user
	Update(user *entities.User) error

	// Delete removes a user by ID
	Delete(id valueobjects.UserID) error

	// ExistsByID checks if a user exists by ID
	ExistsByID(id valueobjects.UserID) (bool, error)

	// ExistsByEmail checks if a user exists by email address
	ExistsByEmail(email valueobjects.Email) (bool, error)

	// FindAll retrieves all users (for admin purposes)
	FindAll() ([]*entities.User, error)

	// Count returns the total number of users
	Count() (int64, error)
}