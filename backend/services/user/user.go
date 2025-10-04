package user

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"todo-app/internal/dtos"
)

// UserService handles user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID uint) (*dtos.User, error) {
	var user dtos.User

	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*dtos.User, error) {
	var user dtos.User

	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByGoogleID retrieves a user by Google ID
func (s *UserService) GetUserByGoogleID(googleID string) (*dtos.User, error) {
	var user dtos.User

	result := s.db.Where("google_id = ?", googleID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// CreateOAuthUser creates a new user from OAuth information
func (s *UserService) CreateOAuthUser(email, name, googleID string) (*dtos.User, error) {
	now := time.Now()
	user := dtos.User{
		Email:          email,
		Name:           name,
		GoogleID:       googleID,
		OAuthProvider:  "google",
		OAuthCreatedAt: &now,
		IsActive:       true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// LinkGoogleAccount links a Google account to an existing user
func (s *UserService) LinkGoogleAccount(userID uint, googleID string) (*dtos.User, error) {
	var user dtos.User

	// Find the user
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if user already has Google linked
	if user.GoogleID != "" {
		return nil, errors.New("user already has a Google account linked")
	}

	// Check if this Google ID is already used by another user
	var existingUser dtos.User
	result = s.db.Where("google_id = ?", googleID).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("this Google account is already linked to another user")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// Link the account
	now := time.Now()
	if err := user.LinkGoogleAccount(googleID, now); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UnlinkGoogleAccount removes Google OAuth linking from a user
func (s *UserService) UnlinkGoogleAccount(userID uint) (*dtos.User, error) {
	var user dtos.User

	// Find the user
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if user has Google linked
	if user.GoogleID == "" {
		return nil, errors.New("user does not have a Google account linked")
	}

	// Unlink the account
	if err := user.UnlinkGoogleAccount(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindOrCreateOAuthUser finds an existing user or creates a new one from OAuth data
// This implements automatic account linking based on email
func (s *UserService) FindOrCreateOAuthUser(email, name, googleID string) (*dtos.User, bool, error) {
	var user dtos.User
	isNewUser := false

	// Try to find user by Google ID
	result := s.db.Where("google_id = ?", googleID).First(&user)
	if result.Error == nil {
		// User exists with this Google ID
		return &user, false, nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false, result.Error
	}

	// Try to find user by email (for automatic account linking)
	result = s.db.Where("email = ?", email).First(&user)
	if result.Error == nil {
		// User exists with this email - link Google account automatically
		now := time.Now()
		err := user.LinkGoogleAccount(googleID, now)
		if err != nil {
			return nil, false, err
		}

		// Save the linked account
		if err := s.db.Save(&user).Error; err != nil {
			return nil, false, err
		}

		return &user, false, nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false, result.Error
	}

	// Create new user
	newUser, err := s.CreateOAuthUser(email, name, googleID)
	if err != nil {
		return nil, false, err
	}

	isNewUser = true
	return newUser, isNewUser, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(userID uint, name string) (*dtos.User, error) {
	var user dtos.User

	// Find the user
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update profile
	if err := user.UpdateProfile(name); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(userID uint) (*dtos.User, error) {
	var user dtos.User

	// Find the user
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Deactivate
	user.Deactivate()

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(userID uint) (*dtos.User, error) {
	var user dtos.User

	// Find the user
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Activate
	user.Activate()

	// Save changes
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(limit, offset int) ([]dtos.User, int64, error) {
	var users []dtos.User
	var total int64

	// Get total count
	if err := s.db.Model(&dtos.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	result := s.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// SearchUsersByEmail searches for users by email pattern
func (s *UserService) SearchUsersByEmail(emailPattern string) ([]dtos.User, error) {
	var users []dtos.User

	result := s.db.Where("email LIKE ?", "%"+emailPattern+"%").
		Order("email").
		Limit(50).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetOAuthUsers retrieves all users who signed up via OAuth
func (s *UserService) GetOAuthUsers() ([]dtos.User, error) {
	var users []dtos.User

	result := s.db.Where("google_id IS NOT NULL AND google_id != ''").
		Order("oauth_created_at DESC").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetLinkedUsers retrieves all users who have both password and OAuth authentication
func (s *UserService) GetLinkedUsers() ([]dtos.User, error) {
	var users []dtos.User

	result := s.db.Where("google_id IS NOT NULL AND google_id != '' AND password_hash IS NOT NULL AND password_hash != ''").
		Order("created_at DESC").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}