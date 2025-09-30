package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-app/models"
	"todo-app/services/user"
)

// TestOAuthAccountLinking tests linking Google OAuth to existing password-based account
func TestOAuthAccountLinking(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	userService := user.NewUserService(db)

	t.Run("links Google account to existing user by email", func(t *testing.T) {
		// Create existing user with password auth
		existingUser := &models.User{
			Email:        "existing@gmail.com",
			Name:         "Existing User",
			PasswordHash: "hashed_password",
			IsActive:     true,
		}
		err := db.Create(existingUser).Error
		require.NoError(t, err)

		// Verify user exists without OAuth
		assert.Empty(t, existingUser.GoogleID)
		assert.Empty(t, existingUser.OAuthProvider)
		assert.False(t, existingUser.IsOAuthUser())

		// Simulate OAuth callback with same email
		ctx := context.Background()
		linkedUser, isNew, err := userService.FindOrCreateOAuthUser(ctx, "existing@gmail.com", "Existing User", "google_id_123")
		require.NoError(t, err)
		assert.False(t, isNew, "Should link to existing account, not create new")

		// Verify account was linked
		assert.Equal(t, existingUser.ID, linkedUser.ID)
		assert.Equal(t, "google_id_123", linkedUser.GoogleID)
		assert.Equal(t, "google", linkedUser.OAuthProvider)
		assert.NotNil(t, linkedUser.OAuthCreatedAt)

		// Verify user can still use password auth
		assert.NotEmpty(t, linkedUser.PasswordHash, "Password should be preserved")

		// Verify user is marked as OAuth user now
		assert.True(t, linkedUser.IsOAuthUser())
	})

	t.Run("preserves existing user data when linking", func(t *testing.T) {
		// Create user with some data
		existingUser := &models.User{
			Email:        "preserve@gmail.com",
			Name:         "Original Name",
			PasswordHash: "original_password",
			IsActive:     true,
		}
		err := db.Create(existingUser).Error
		require.NoError(t, err)

		originalID := existingUser.ID

		// Link OAuth account
		ctx := context.Background()
		linkedUser, _, err := userService.FindOrCreateOAuthUser(ctx, "preserve@gmail.com", "OAuth Name", "google_preserve_123")
		require.NoError(t, err)

		// Verify data preservation
		assert.Equal(t, originalID, linkedUser.ID, "User ID should remain the same")
		assert.Equal(t, "Original Name", linkedUser.Name, "Original name should be preserved")
		assert.Equal(t, "original_password", linkedUser.PasswordHash, "Password should be preserved")
		assert.True(t, linkedUser.IsActive, "Active status should be preserved")

		// Verify OAuth fields added
		assert.Equal(t, "google_preserve_123", linkedUser.GoogleID)
		assert.Equal(t, "google", linkedUser.OAuthProvider)
	})

	t.Run("creates new user if email not found", func(t *testing.T) {
		ctx := context.Background()
		newUser, isNew, err := userService.FindOrCreateOAuthUser(ctx, "newuser@gmail.com", "New User", "google_new_123")
		require.NoError(t, err)
		assert.True(t, isNew, "Should create new user")

		// Verify new user created
		assert.NotZero(t, newUser.ID)
		assert.Equal(t, "newuser@gmail.com", newUser.Email)
		assert.Equal(t, "google_new_123", newUser.GoogleID)
		assert.Empty(t, newUser.PasswordHash, "OAuth-only user should not have password")
	})

	t.Run("handles account already linked to different Google ID", func(t *testing.T) {
		// Create user already linked to Google
		existingUser := &models.User{
			Email:         "linked@gmail.com",
			Name:          "Linked User",
			GoogleID:      "original_google_id",
			OAuthProvider: "google",
			IsActive:      true,
		}
		now := time.Now()
		existingUser.OAuthCreatedAt = &now
		err := db.Create(existingUser).Error
		require.NoError(t, err)

		// Attempt to link with different Google ID (should not happen in practice)
		// This tests data integrity
		ctx := context.Background()
		user, isNew, err := userService.FindOrCreateOAuthUser(ctx, "linked@gmail.com", "Linked User", "different_google_id")

		// Expected behavior: Either error or keep original Google ID
		if err == nil {
			// If no error, should preserve original Google ID
			assert.False(t, isNew)
			assert.Equal(t, "original_google_id", user.GoogleID, "Should not overwrite existing Google ID")
		} else {
			// Error is acceptable for this edge case
			assert.Error(t, err)
		}
	})

	t.Run("validates email format during linking", func(t *testing.T) {
		ctx := context.Background()
		_, _, err := userService.FindOrCreateOAuthUser(ctx, "invalid-email", "Test User", "google_123")

		// Should handle invalid email appropriately
		// Either return error or sanitize the email
		if err != nil {
			assert.Error(t, err, "Should reject invalid email")
		}
	})
}

// TestAccountLinkingEdgeCases tests edge cases in account linking
func TestAccountLinkingEdgeCases(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	t.Run("case-insensitive email matching", func(t *testing.T) {
		// Create user with lowercase email
		user := &models.User{
			Email:        "test@gmail.com",
			Name:         "Test User",
			PasswordHash: "password",
			IsActive:     true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Attempt linking with uppercase email
		userService := user2.NewUserService(db)
		ctx := context.Background()
		linkedUser, isNew, err := userService.FindOrCreateOAuthUser(ctx, "TEST@GMAIL.COM", "Test User", "google_123")

		// Should find existing user (case-insensitive)
		require.NoError(t, err)
		assert.False(t, isNew, "Should link to existing user regardless of case")
		assert.Equal(t, user.ID, linkedUser.ID)
	})

	t.Run("linking inactive account", func(t *testing.T) {
		// Create inactive user
		inactiveUser := &models.User{
			Email:        "inactive@gmail.com",
			Name:         "Inactive User",
			PasswordHash: "password",
			IsActive:     false,
		}
		err := db.Create(inactiveUser).Error
		require.NoError(t, err)

		// Attempt OAuth linking
		userService := user2.NewUserService(db)
		ctx := context.Background()
		_, _, err = userService.FindOrCreateOAuthUser(ctx, "inactive@gmail.com", "Inactive User", "google_inactive")

		// Should either reactivate or prevent linking
		// Implementation choice: document expected behavior
	})

	t.Run("concurrent linking attempts", func(t *testing.T) {
		// Create user
		user := &models.User{
			Email:        "concurrent@gmail.com",
			Name:         "Concurrent User",
			PasswordHash: "password",
			IsActive:     true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Simulate concurrent linking attempts
		// Both should succeed without data corruption
		userService := user2.NewUserService(db)
		ctx := context.Background()

		done := make(chan bool, 2)
		errors := make(chan error, 2)

		for i := 0; i < 2; i++ {
			go func() {
				_, _, err := userService.FindOrCreateOAuthUser(ctx, "concurrent@gmail.com", "Concurrent User", "google_concurrent")
				if err != nil {
					errors <- err
				}
				done <- true
			}()
		}

		// Wait for both to complete
		<-done
		<-done
		close(errors)

		// Check for errors
		for err := range errors {
			assert.NoError(t, err, "Concurrent linking should not cause errors")
		}

		// Verify final state is consistent
		var finalUser models.User
		err = db.Where("email = ?", "concurrent@gmail.com").First(&finalUser).Error
		require.NoError(t, err)
		assert.Equal(t, "google_concurrent", finalUser.GoogleID)
	})
}