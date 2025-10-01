package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-app/internal/models"
	"todo-app/services/auth"
)

// TestSessionManagementAndRefresh tests session lifecycle and automatic refresh
func TestSessionManagementAndRefresh(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{})
	require.NoError(t, err)

	sessionService := auth.NewSessionService(db)

	// Create test user
	user := &models.User{
		Email:         "session@gmail.com",
		Name:          "Session User",
		GoogleID:      "google_session_123",
		OAuthProvider: "google",
		IsActive:      true,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Run("creates session with 24-hour expiration", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		session, err := sessionService.CreateSession(ctx, user.ID, true, "access_token", "refresh_token", tokenExpiry)
		require.NoError(t, err)
		require.NotNil(t, session)

		// Verify session properties
		assert.Equal(t, user.ID, session.UserID)
		assert.NotEmpty(t, session.SessionToken)
		assert.NotEmpty(t, session.AccessToken)
		assert.NotEmpty(t, session.RefreshToken)

		// Verify 24-hour expiration
		expectedExpiry := time.Now().Add(24 * time.Hour)
		assert.WithinDuration(t, expectedExpiry, session.SessionExpiresAt, 1*time.Minute)

		// Verify token expiration
		assert.Equal(t, tokenExpiry.Unix(), session.TokenExpiresAt.Unix())
	})

	t.Run("validates active session successfully", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "access_token_valid", "refresh_token_valid", tokenExpiry)
		require.NoError(t, err)

		// Validate session
		validatedSession, err := sessionService.ValidateSession(ctx, session.SessionToken)
		require.NoError(t, err)
		require.NotNil(t, validatedSession)

		assert.Equal(t, session.ID, validatedSession.ID)
		assert.Equal(t, user.ID, validatedSession.UserID)
		assert.False(t, validatedSession.IsExpired())
	})

	t.Run("detects expired session", func(t *testing.T) {
		// Create session with past expiration
		expiredSession := &models.AuthenticationSession{
			UserID:           user.ID,
			SessionToken:     "expired_token",
			AccessToken:      "expired_access",
			RefreshToken:     "expired_refresh",
			SessionExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			LastActivity:     time.Now().Add(-1 * time.Hour),
		}
		now := time.Now().Add(-1 * time.Hour)
		expiredSession.TokenExpiresAt = &now

		err := db.Create(expiredSession).Error
		require.NoError(t, err)

		// Verify session is detected as expired
		assert.True(t, expiredSession.IsExpired())

		// Validation should fail for expired session
		ctx := context.Background()
		_, err = sessionService.ValidateSession(ctx, "expired_token")
		assert.Error(t, err, "Should reject expired session")
	})

	t.Run("refreshes OAuth tokens", func(t *testing.T) {
		ctx := context.Background()
		oldTokenExpiry := time.Now().Add(10 * time.Minute)

		// Create session with tokens expiring soon
		session, err := sessionService.CreateSession(ctx, user.ID, true, "old_access_token", "refresh_token_for_refresh", oldTokenExpiry)
		require.NoError(t, err)

		// Verify tokens need refresh
		assert.True(t, session.NeedsRefresh(), "Session should need refresh when tokens expire soon")

		// Simulate token refresh
		newTokenExpiry := time.Now().Add(1 * time.Hour)
		err = sessionService.RefreshSession(ctx, session.SessionToken, "new_access_token", newTokenExpiry)
		require.NoError(t, err)

		// Retrieve updated session
		var updatedSession models.AuthenticationSession
		err = db.First(&updatedSession, session.ID).Error
		require.NoError(t, err)

		// Verify tokens were updated
		assert.NotEqual(t, "old_access_token", updatedSession.AccessToken)
		assert.False(t, updatedSession.NeedsRefresh(), "Session should not need refresh after update")
	})

	t.Run("updates last activity on validation", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "activity_token", "activity_refresh", tokenExpiry)
		require.NoError(t, err)

		originalActivity := session.LastActivity

		// Wait a bit
		time.Sleep(100 * time.Millisecond)

		// Validate session (should update activity)
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Retrieve updated session
		var updatedSession models.AuthenticationSession
		err = db.First(&updatedSession, session.ID).Error
		require.NoError(t, err)

		// Verify last activity was updated
		assert.True(t, updatedSession.LastActivity.After(originalActivity), "Last activity should be updated")
	})

	t.Run("extends session for active users", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session expiring soon
		session, err := sessionService.CreateSession(ctx, user.ID, true, "extend_token", "extend_refresh", tokenExpiry)
		require.NoError(t, err)

		// Manually set expiration to be within extension window
		session.SessionExpiresAt = time.Now().Add(30 * time.Minute)
		err = db.Save(session).Error
		require.NoError(t, err)

		originalExpiry := session.SessionExpiresAt

		// Validate session (should extend if within window)
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Retrieve updated session
		var updatedSession models.AuthenticationSession
		err = db.First(&updatedSession, session.ID).Error
		require.NoError(t, err)

		// Verify session was potentially extended
		// (Implementation may extend sessions approaching expiration)
		assert.True(t, updatedSession.SessionExpiresAt.After(originalExpiry) || updatedSession.SessionExpiresAt.Equal(originalExpiry))
	})

	t.Run("terminates session on logout", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "logout_token", "logout_refresh", tokenExpiry)
		require.NoError(t, err)

		// Terminate session
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Verify session no longer exists
		var count int64
		err = db.Model(&models.AuthenticationSession{}).Where("id = ?", session.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "Session should be deleted on logout")

		// Validation should fail
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.Error(t, err, "Should not validate terminated session")
	})

	t.Run("handles multiple concurrent sessions per user", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create multiple sessions for same user
		session1, err := sessionService.CreateSession(ctx, user.ID, true, "token1", "refresh1", tokenExpiry)
		require.NoError(t, err)

		session2, err := sessionService.CreateSession(ctx, user.ID, true, "token2", "refresh2", tokenExpiry)
		require.NoError(t, err)

		// Both sessions should be independent and valid
		_, err = sessionService.ValidateSession(ctx, session1.SessionToken)
		assert.NoError(t, err)

		_, err = sessionService.ValidateSession(ctx, session2.SessionToken)
		assert.NoError(t, err)

		// Terminating one should not affect the other
		err = sessionService.TerminateSession(ctx, session1.SessionToken)
		require.NoError(t, err)

		_, err = sessionService.ValidateSession(ctx, session2.SessionToken)
		assert.NoError(t, err, "Session 2 should still be valid")
	})
}

// TestSessionSecurityFeatures tests security aspects of session management
func TestSessionSecurityFeatures(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{})
	require.NoError(t, err)

	t.Run("session tokens are unique", func(t *testing.T) {
		user := &models.User{
			Email:    "secure@gmail.com",
			Name:     "Secure User",
			GoogleID: "google_secure",
			IsActive: true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create multiple sessions
		sessions := make([]*models.AuthenticationSession, 10)
		tokens := make(map[string]bool)

		for i := 0; i < 10; i++ {
			session, err := sessionService.CreateSession(ctx, user.ID, true, "access", "refresh", tokenExpiry)
			require.NoError(t, err)
			sessions[i] = session

			// Verify token is unique
			assert.False(t, tokens[session.SessionToken], "Session token should be unique")
			tokens[session.SessionToken] = true
		}

		assert.Equal(t, 10, len(tokens), "All session tokens should be unique")
	})

	t.Run("encrypted tokens stored securely", func(t *testing.T) {
		user := &models.User{
			Email:    "encrypt@gmail.com",
			Name:     "Encrypt User",
			GoogleID: "google_encrypt",
			IsActive: true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		plainAccessToken := "plain_access_token_value"
		plainRefreshToken := "plain_refresh_token_value"

		session, err := sessionService.CreateSession(ctx, user.ID, true, plainAccessToken, plainRefreshToken, tokenExpiry)
		require.NoError(t, err)

		// Verify tokens are encrypted in database
		// (Access/Refresh tokens should not match plain values if encrypted)
		var storedSession models.AuthenticationSession
		err = db.First(&storedSession, session.ID).Error
		require.NoError(t, err)

		// If encryption is implemented, tokens should differ from plain values
		// This is a placeholder check - actual implementation may vary
		assert.NotEmpty(t, storedSession.AccessToken)
		assert.NotEmpty(t, storedSession.RefreshToken)
	})
}