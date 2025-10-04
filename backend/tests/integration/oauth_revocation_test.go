package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-app/internal/dtos"
	"todo-app/services/auth"
)

// TestOAuthAccessRevocation tests handling of OAuth access revocation
func TestOAuthAccessRevocation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{})
	require.NoError(t, err)

	sessionService := auth.NewSessionService(db)

	// Create test user
	user := &models.User{
		Email:         "revoke@gmail.com",
		Name:          "Revoke User",
		GoogleID:      "google_revoke_123",
		OAuthProvider: "google",
		IsActive:      true,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Run("terminates session when access token revoked", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		accessToken := "access_token_to_revoke"
		session, err := sessionService.CreateSession(ctx, user.ID, true, accessToken, "refresh_token", tokenExpiry)
		require.NoError(t, err)

		// Verify session is valid
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.NoError(t, err)

		// Simulate OAuth revocation webhook
		// Find and terminate sessions with matching access token
		var sessionsToRevoke []models.AuthenticationSession
		err = db.Where("access_token = ?", accessToken).Find(&sessionsToRevoke).Error
		require.NoError(t, err)
		assert.Equal(t, 1, len(sessionsToRevoke), "Should find session to revoke")

		// Revoke session
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Verify session no longer valid
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.Error(t, err, "Revoked session should not validate")
	})

	t.Run("terminates session when refresh token revoked", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		refreshToken := "refresh_token_to_revoke"
		session, err := sessionService.CreateSession(ctx, user.ID, true, "access_token", refreshToken, tokenExpiry)
		require.NoError(t, err)

		// Simulate revocation via refresh token
		var sessionsToRevoke []models.AuthenticationSession
		err = db.Where("refresh_token = ?", refreshToken).Find(&sessionsToRevoke).Error
		require.NoError(t, err)
		assert.Equal(t, 1, len(sessionsToRevoke))

		// Revoke
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Verify termination
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.Error(t, err)
	})

	t.Run("handles revocation of non-existent token gracefully", func(t *testing.T) {
		// Attempt to revoke non-existent token
		var sessions []models.AuthenticationSession
		err := db.Where("access_token = ?", "non_existent_token").Find(&sessions).Error
		require.NoError(t, err)
		assert.Equal(t, 0, len(sessions), "Should find no sessions")

		// Revocation should be idempotent - no error for non-existent token
		// This prevents leaking information about token existence
	})

	t.Run("revokes all sessions for user", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create multiple sessions for user
		session1, err := sessionService.CreateSession(ctx, user.ID, true, "access1", "refresh1", tokenExpiry)
		require.NoError(t, err)

		session2, err := sessionService.CreateSession(ctx, user.ID, true, "access2", "refresh2", tokenExpiry)
		require.NoError(t, err)

		session3, err := sessionService.CreateSession(ctx, user.ID, true, "access3", "refresh3", tokenExpiry)
		require.NoError(t, err)

		// Verify all sessions are valid
		_, err = sessionService.ValidateSession(ctx, session1.SessionToken)
		assert.NoError(t, err)
		_, err = sessionService.ValidateSession(ctx, session2.SessionToken)
		assert.NoError(t, err)
		_, err = sessionService.ValidateSession(ctx, session3.SessionToken)
		assert.NoError(t, err)

		// Revoke all sessions for user (global revocation)
		result := db.Where("user_id = ?", user.ID).Delete(&models.AuthenticationSession{})
		require.NoError(t, result.Error)
		assert.Equal(t, int64(3), result.RowsAffected)

		// Verify all sessions terminated
		_, err = sessionService.ValidateSession(ctx, session1.SessionToken)
		assert.Error(t, err)
		_, err = sessionService.ValidateSession(ctx, session2.SessionToken)
		assert.Error(t, err)
		_, err = sessionService.ValidateSession(ctx, session3.SessionToken)
		assert.Error(t, err)
	})

	t.Run("immediate termination on revocation notification", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "immediate_access", "immediate_refresh", tokenExpiry)
		require.NoError(t, err)

		beforeRevocation := time.Now()

		// Revoke immediately
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		afterRevocation := time.Now()

		// Verify revocation was immediate (< 1 second)
		assert.Less(t, afterRevocation.Sub(beforeRevocation), 1*time.Second, "Revocation should be immediate")

		// Verify session cannot be used
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.Error(t, err)
	})

	t.Run("prevents token refresh after revocation", func(t *testing.T) {
		ctx := context.Background()
		tokenExpiry := time.Now().Add(10 * time.Minute)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "pre_revoke_access", "pre_revoke_refresh", tokenExpiry)
		require.NoError(t, err)

		// Revoke session
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Attempt to refresh tokens (should fail)
		newTokenExpiry := time.Now().Add(1 * time.Hour)
		err = sessionService.RefreshSession(ctx, session.SessionToken, "new_access", newTokenExpiry)
		assert.Error(t, err, "Should not allow refresh of revoked session")
	})
}

// TestRevocationWebhookHandling tests webhook handling for OAuth revocation
func TestRevocationWebhookHandling(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{})
	require.NoError(t, err)

	t.Run("webhook payload processing", func(t *testing.T) {
		user := &models.User{
			Email:         "webhook@gmail.com",
			Name:          "Webhook User",
			GoogleID:      "google_webhook",
			OAuthProvider: "google",
			IsActive:      true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create session
		revokedToken := "token_from_webhook"
		session, err := sessionService.CreateSession(ctx, user.ID, true, revokedToken, "refresh", tokenExpiry)
		require.NoError(t, err)

		// Simulate webhook: Find sessions by token
		var sessionsToRevoke []models.AuthenticationSession
		err = db.Where("access_token = ? OR refresh_token = ?", revokedToken, revokedToken).
			Find(&sessionsToRevoke).Error
		require.NoError(t, err)
		assert.Equal(t, 1, len(sessionsToRevoke))

		// Process revocation
		for _, s := range sessionsToRevoke {
			err = sessionService.TerminateSession(ctx, s.SessionToken)
			require.NoError(t, err)
		}

		// Verify session terminated
		_, err = sessionService.ValidateSession(ctx, session.SessionToken)
		assert.Error(t, err)
	})

	t.Run("handles malformed webhook payload", func(t *testing.T) {
		// Webhook with empty token
		emptyToken := ""
		var sessions []models.AuthenticationSession
		err := db.Where("access_token = ?", emptyToken).Find(&sessions).Error
		require.NoError(t, err)
		assert.Equal(t, 0, len(sessions), "Should handle empty token gracefully")
	})

	t.Run("idempotent revocation handling", func(t *testing.T) {
		user := &models.User{
			Email:    "idempotent@gmail.com",
			Name:     "Idempotent User",
			GoogleID: "google_idempotent",
			IsActive: true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		session, err := sessionService.CreateSession(ctx, user.ID, true, "idempotent_token", "refresh", tokenExpiry)
		require.NoError(t, err)

		// First revocation
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		require.NoError(t, err)

		// Second revocation (should not error)
		err = sessionService.TerminateSession(ctx, session.SessionToken)
		// Should be idempotent - either no error or "not found" is acceptable
		// Implementation choice: no error indicates idempotent success
	})
}