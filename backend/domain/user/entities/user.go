package entities

import (
	"errors"
	"time"

	"domain/user/valueobjects"
)

// User represents a domain entity for user management
type User struct {
	id          valueobjects.UserID
	email       valueobjects.Email
	profile     valueobjects.UserProfile
	preferences valueobjects.UserPreferences
	createdAt   time.Time
	updatedAt   time.Time
}

// NewUser creates a new User entity
func NewUser(
	id valueobjects.UserID,
	email valueobjects.Email,
	profile valueobjects.UserProfile,
	preferences valueobjects.UserPreferences,
) (*User, error) {
	if id.IsZero() {
		return nil, errors.New("user ID cannot be zero")
	}

	if email.IsEmpty() {
		return nil, errors.New("user email cannot be empty")
	}

	now := time.Now()

	return &User{
		id:          id,
		email:       email,
		profile:     profile,
		preferences: preferences,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// NewUserWithDefaults creates a new User entity with default preferences
func NewUserWithDefaults(
	id valueobjects.UserID,
	email valueobjects.Email,
	profile valueobjects.UserProfile,
) (*User, error) {
	defaultPreferences := valueobjects.NewDefaultUserPreferences()
	return NewUser(id, email, profile, defaultPreferences)
}

// UpdateProfile updates the user profile
func (u *User) UpdateProfile(profile valueobjects.UserProfile) error {
	u.profile = profile
	u.updatedAt = time.Now()
	return nil
}

// UpdatePreferences updates the user preferences
func (u *User) UpdatePreferences(preferences valueobjects.UserPreferences) error {
	u.preferences = preferences
	u.updatedAt = time.Now()
	return nil
}

// ChangeEmail updates the user's email address
func (u *User) ChangeEmail(email valueobjects.Email) error {
	if email.IsEmpty() {
		return errors.New("email cannot be empty")
	}

	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// GetDisplayName returns the user's display name from profile
func (u *User) GetDisplayName() string {
	return u.profile.DisplayName()
}

// IsActive returns true if the user is active (for future use with user status)
func (u *User) IsActive() bool {
	// For now, all users are considered active
	// This can be extended when user status/deactivation is needed
	return true
}

// HasEmailDomain checks if the user's email belongs to a specific domain
func (u *User) HasEmailDomain(domain string) bool {
	return u.email.Domain() == domain
}

// UpdateDefaultTaskPriority updates the default task priority in preferences
// Note: This method will be called from the application layer with proper typing
func (u *User) UpdateDefaultTaskPriority(newPrefs valueobjects.UserPreferences) error {
	u.preferences = newPrefs
	u.updatedAt = time.Now()
	return nil
}

// EnableEmailNotifications enables email notifications
func (u *User) EnableEmailNotifications() error {
	u.preferences = u.preferences.WithEmailNotifications(true)
	u.updatedAt = time.Now()
	return nil
}

// DisableEmailNotifications disables email notifications
func (u *User) DisableEmailNotifications() error {
	u.preferences = u.preferences.WithEmailNotifications(false)
	u.updatedAt = time.Now()
	return nil
}

// UpdateThemePreference updates the theme preference
func (u *User) UpdateThemePreference(theme string) error {
	newPrefs, err := u.preferences.WithThemePreference(theme)
	if err != nil {
		return err
	}

	u.preferences = newPrefs
	u.updatedAt = time.Now()
	return nil
}

// Getters

// ID returns the user ID
func (u *User) ID() valueobjects.UserID {
	return u.id
}

// Email returns the user email
func (u *User) Email() valueobjects.Email {
	return u.email
}

// Profile returns the user profile
func (u *User) Profile() valueobjects.UserProfile {
	return u.profile
}

// Preferences returns the user preferences
func (u *User) Preferences() valueobjects.UserPreferences {
	return u.preferences
}

// CreatedAt returns the creation time
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update time
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}