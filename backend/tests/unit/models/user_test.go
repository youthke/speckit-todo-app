package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-app/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the User table
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func TestUser_OAuthFieldsValidation(t *testing.T) {
	tests := []struct {
		name        string
		user        models.User
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid traditional user",
			user: models.User{
				Email:        "test@example.com",
				Name:         "Test User",
				PasswordHash: "hashed_password",
				IsActive:     true,
			},
			shouldError: false,
		},
		{
			name: "valid OAuth user",
			user: models.User{
				Email:           "oauth@gmail.com",
				Name:            "OAuth User",
				GoogleID:        "google_123456",
				OAuthProvider:   "google",
				OAuthCreatedAt:  &[]time.Time{time.Now()}[0],
				IsActive:        true,
			},
			shouldError: false,
		},
		{
			name: "invalid - empty email",
			user: models.User{
				Name:         "Test User",
				PasswordHash: "hashed_password",
				IsActive:     true,
			},
			shouldError: true,
			errorMsg:    "email cannot be empty",
		},
		{
			name: "invalid - empty name",
			user: models.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				IsActive:     true,
			},
			shouldError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name: "invalid - both password and google_id empty",
			user: models.User{
				Email:    "test@example.com",
				Name:     "Test User",
				IsActive: true,
			},
			shouldError: true,
			errorMsg:    "either password_hash or google_id must be present",
		},
		{
			name: "invalid - google_id without oauth_provider",
			user: models.User{
				Email:    "test@gmail.com",
				Name:     "Test User",
				GoogleID: "google_123456",
				IsActive: true,
			},
			shouldError: true,
			errorMsg:    "oauth_provider must be 'google' when google_id is present",
		},
		{
			name: "invalid - oauth_provider without google_id",
			user: models.User{
				Email:         "test@gmail.com",
				Name:          "Test User",
				OAuthProvider: "google",
				IsActive:      true,
			},
			shouldError: true,
			errorMsg:    "google_id must be present when oauth_provider is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_CreateOAuthUser(t *testing.T) {
	db := setupTestDB(t)

	now := time.Now()
	user := models.User{
		Email:          "oauth@gmail.com",
		Name:           "OAuth User",
		GoogleID:       "google_123456789",
		OAuthProvider:  "google",
		OAuthCreatedAt: &now,
		IsActive:       true,
	}

	result := db.Create(&user)
	require.NoError(t, result.Error)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUser_CreateTraditionalUser(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{
		Email:        "traditional@example.com",
		Name:         "Traditional User",
		PasswordHash: "hashed_password_123",
		IsActive:     true,
	}

	result := db.Create(&user)
	require.NoError(t, result.Error)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUser_UniqueConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Create first user
	now1 := time.Now()
	user1 := models.User{
		Email:        "unique@example.com",
		Name:         "User One",
		GoogleID:     "google_unique_123",
		OAuthProvider: "google",
		OAuthCreatedAt: &now1,
		IsActive:     true,
	}

	result := db.Create(&user1)
	require.NoError(t, result.Error)

	// Try to create user with same email
	user2 := models.User{
		Email:        "unique@example.com", // Same email
		Name:         "User Two",
		PasswordHash: "different_password",
		IsActive:     true,
	}

	result = db.Create(&user2)
	assert.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")

	// Try to create user with same google_id
	now3 := time.Now()
	user3 := models.User{
		Email:         "different@example.com",
		Name:          "User Three",
		GoogleID:      "google_unique_123", // Same GoogleID
		OAuthProvider: "google",
		OAuthCreatedAt: &now3,
		IsActive:      true,
	}

	result = db.Create(&user3)
	assert.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")
}

func TestUser_AccountLinking(t *testing.T) {
	db := setupTestDB(t)

	// Create traditional user
	user := models.User{
		Email:        "linking@example.com",
		Name:         "Linking User",
		PasswordHash: "original_password",
		IsActive:     true,
	}

	result := db.Create(&user)
	require.NoError(t, result.Error)
	originalID := user.ID

	// Link OAuth account
	err := user.LinkGoogleAccount("google_linked_123", time.Now())
	require.NoError(t, err)

	// Save changes
	result = db.Save(&user)
	require.NoError(t, result.Error)

	// Verify linking
	assert.Equal(t, originalID, user.ID)
	assert.Equal(t, "google_linked_123", user.GoogleID)
	assert.Equal(t, "google", user.OAuthProvider)
	assert.NotNil(t, user.OAuthCreatedAt)
	assert.NotEmpty(t, user.PasswordHash) // Should preserve original password
}

func TestUser_IsOAuthUser(t *testing.T) {
	tests := []struct {
		name     string
		user     models.User
		expected bool
	}{
		{
			name: "OAuth user",
			user: models.User{
				GoogleID:      "google_123",
				OAuthProvider: "google",
			},
			expected: true,
		},
		{
			name: "traditional user",
			user: models.User{
				PasswordHash: "hashed_password",
			},
			expected: false,
		},
		{
			name: "linked user (both OAuth and password)",
			user: models.User{
				GoogleID:      "google_123",
				OAuthProvider: "google",
				PasswordHash:  "hashed_password",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsOAuthUser()
			assert.Equal(t, tt.expected, result)
		})
	}
}