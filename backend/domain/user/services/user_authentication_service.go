package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"todo-app/domain/user/repositories"
	"todo-app/domain/user/valueobjects"
)

// UserCredentials represents user authentication credentials
type UserCredentials struct {
	UserID         valueobjects.UserID
	Email          valueobjects.Email
	HashedPassword string
	Salt           string
}

// UserAuthenticationService provides domain authentication logic for users
type UserAuthenticationService interface {
	// ValidateEmailUniqueness ensures email is unique across the system
	ValidateEmailUniqueness(email valueobjects.Email) error

	// GenerateUserCredentials creates authentication credentials for a new user
	GenerateUserCredentials(email valueobjects.Email) (*UserCredentials, error)

	// ValidateRegistrationData validates all data required for user registration
	ValidateRegistrationData(email valueobjects.Email, profile valueobjects.UserProfile) error
}

// userAuthenticationService implements UserAuthenticationService
type userAuthenticationService struct {
	userRepo repositories.UserRepository
}

// NewUserAuthenticationService creates a new user authentication service
func NewUserAuthenticationService(userRepo repositories.UserRepository) UserAuthenticationService {
	return &userAuthenticationService{
		userRepo: userRepo,
	}
}

// ValidateEmailUniqueness ensures the email is not already taken
func (s *userAuthenticationService) ValidateEmailUniqueness(email valueobjects.Email) error {
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email address is already registered")
	}

	return nil
}

// GenerateUserCredentials creates new authentication credentials
func (s *userAuthenticationService) GenerateUserCredentials(email valueobjects.Email) (*UserCredentials, error) {
	// Generate a random salt
	salt, err := generateRandomSalt(32)
	if err != nil {
		return nil, err
	}

	// For now, we'll generate a placeholder hashed password
	// In a real implementation, this would hash a temporary password or work with OAuth
	hashedPassword, err := generateRandomSalt(64)
	if err != nil {
		return nil, err
	}

	// Create a new UserID (will be set by the repository when saving)
	userID := valueobjects.UserID{}

	return &UserCredentials{
		UserID:         userID,
		Email:          email,
		HashedPassword: hashedPassword,
		Salt:           salt,
	}, nil
}

// ValidateRegistrationData validates all data required for user registration
func (s *userAuthenticationService) ValidateRegistrationData(email valueobjects.Email, profile valueobjects.UserProfile) error {
	// Validate email uniqueness
	if err := s.ValidateEmailUniqueness(email); err != nil {
		return err
	}

	// Validate email format (should already be validated in value object, but double-check)
	if email.IsEmpty() {
		return errors.New("email is required for registration")
	}

	// Validate profile data (should already be validated in value object, but ensure completeness)
	if profile.FirstName() == "" {
		return errors.New("first name is required for registration")
	}

	if profile.LastName() == "" {
		return errors.New("last name is required for registration")
	}

	if profile.Timezone() == "" {
		return errors.New("timezone is required for registration")
	}

	return nil
}

// generateRandomSalt generates a random salt of the specified byte length
func generateRandomSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}