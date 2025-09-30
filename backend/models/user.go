package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system with OAuth support
type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	Name      string `json:"name" gorm:"type:varchar(255);not null" validate:"required"`

	// Traditional authentication
	PasswordHash string `json:"-" gorm:"type:varchar(255)"`

	// OAuth fields
	GoogleID       string     `json:"google_id,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	OAuthProvider  string     `json:"oauth_provider,omitempty" gorm:"type:varchar(50)"`
	OAuthCreatedAt *time.Time `json:"oauth_created_at,omitempty"`

	// Status and timestamps
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to validate user before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// BeforeUpdate hook to validate user before update
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}

// Validate performs validation on the User model
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if u.Name == "" {
		return errors.New("name cannot be empty")
	}

	// Either password_hash OR google_id must be present
	if u.PasswordHash == "" && u.GoogleID == "" {
		return errors.New("either password_hash or google_id must be present")
	}

	// If google_id is present, oauth_provider must be "google"
	if u.GoogleID != "" && u.OAuthProvider != "google" {
		return errors.New("oauth_provider must be 'google' when google_id is present")
	}

	// If oauth_provider is set, google_id must be present
	if u.OAuthProvider != "" && u.GoogleID == "" {
		return errors.New("google_id must be present when oauth_provider is set")
	}

	return nil
}

// IsOAuthUser returns true if the user was created via OAuth
func (u *User) IsOAuthUser() bool {
	return u.GoogleID != ""
}

// IsTraditionalUser returns true if the user has password authentication
func (u *User) IsTraditionalUser() bool {
	return u.PasswordHash != ""
}

// IsLinkedUser returns true if the user has both OAuth and password authentication
func (u *User) IsLinkedUser() bool {
	return u.IsOAuthUser() && u.IsTraditionalUser()
}

// LinkGoogleAccount links a Google OAuth account to an existing user
func (u *User) LinkGoogleAccount(googleID string, linkedAt time.Time) error {
	if googleID == "" {
		return errors.New("google_id cannot be empty")
	}

	u.GoogleID = googleID
	u.OAuthProvider = "google"
	u.OAuthCreatedAt = &linkedAt
	u.UpdatedAt = time.Now()

	return u.Validate()
}

// UnlinkGoogleAccount removes Google OAuth linking from the user
// Only allowed if user has password authentication
func (u *User) UnlinkGoogleAccount() error {
	if !u.IsTraditionalUser() {
		return errors.New("cannot unlink OAuth account without password authentication")
	}

	u.GoogleID = ""
	u.OAuthProvider = ""
	u.OAuthCreatedAt = nil
	u.UpdatedAt = time.Now()

	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates the user's display name
func (u *User) UpdateProfile(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	u.Name = name
	u.UpdatedAt = time.Now()

	return nil
}

// ChangeEmail updates the user's email address
func (u *User) ChangeEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	u.Email = email
	u.UpdatedAt = time.Now()

	return u.Validate()
}

// CreateOAuthUserRequest represents the request for creating an OAuth user
type CreateOAuthUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	GoogleID string `json:"google_id" binding:"required"`
}

// CreateTraditionalUserRequest represents the request for creating a traditional user
type CreateTraditionalUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserProfileRequest represents the request for updating user profile
type UpdateUserProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

// LinkOAuthAccountRequest represents the request for linking OAuth account
type LinkOAuthAccountRequest struct {
	GoogleID string `json:"google_id" binding:"required"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID             uint       `json:"id"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	OAuthProvider  string     `json:"oauth_provider,omitempty"`
	OAuthCreatedAt *time.Time `json:"oauth_created_at,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name,
		OAuthProvider:  u.OAuthProvider,
		OAuthCreatedAt: u.OAuthCreatedAt,
		IsActive:       u.IsActive,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}