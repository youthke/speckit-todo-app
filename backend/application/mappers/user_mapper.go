package mappers

import (
	"fmt"
	"strings"

	"domain/user/entities"
	"domain/user/valueobjects"
	"todo-app/internal/dtos"
)

// UserMapper handles conversion between User DTOs and User entities
type UserMapper struct{}

// ToEntity converts a UserDTO to a User entity
func (m *UserMapper) ToEntity(dto *dtos.User) (*entities.User, error) {
	// Validate and create UserID
	userID := valueobjects.NewUserID(dto.ID)
	if userID.IsZero() {
		return nil, fmt.Errorf("user ID cannot be zero")
	}

	// Validate and create Email
	email, err := valueobjects.NewEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Create UserProfile from Name field
	// The DTO has a single Name field, but UserProfile expects firstName, lastName, timezone
	// We'll split the name and use a default timezone
	profile, err := m.createUserProfileFromName(dto.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid profile: %w", err)
	}

	// Create default UserPreferences
	preferences := valueobjects.NewDefaultUserPreferences()

	// Create the User entity
	user, err := entities.NewUser(userID, email, profile, preferences)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	return user, nil
}

// ToDTO converts a User entity to a UserDTO
func (m *UserMapper) ToDTO(entity *entities.User) *dtos.User {
	return &dtos.User{
		ID:         entity.ID().Value(),
		Email:      entity.Email().Value(),
		Name:       entity.Profile().DisplayName(), // Get full name from profile
		AuthMethod: "password",                      // Default value (auth-related fields managed by Auth domain)
		IsActive:   true,                            // Default value
		CreatedAt:  entity.CreatedAt(),
		UpdatedAt:  entity.UpdatedAt(),
	}
}

// createUserProfileFromName creates a UserProfile from a single name string
// This is a helper method to handle the mismatch between DTO (single Name field)
// and UserProfile (firstName, lastName, timezone)
func (m *UserMapper) createUserProfileFromName(name string) (valueobjects.UserProfile, error) {
	if name == "" {
		return valueobjects.UserProfile{}, fmt.Errorf("name cannot be empty")
	}

	// Split name into first and last name
	// Simple approach: split on first space
	parts := strings.Fields(name)
	var firstName, lastName string

	if len(parts) == 0 {
		return valueobjects.UserProfile{}, fmt.Errorf("name cannot be empty")
	} else if len(parts) == 1 {
		// Single word name: use as both first and last name
		firstName = parts[0]
		lastName = parts[0]
	} else {
		// Multiple words: first word is firstName, rest is lastName
		firstName = parts[0]
		lastName = strings.Join(parts[1:], " ")
	}

	// Use UTC as default timezone
	timezone := "UTC"

	return valueobjects.NewUserProfile(firstName, lastName, timezone)
}
