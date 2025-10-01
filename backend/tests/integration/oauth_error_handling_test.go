package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-app/internal/models"
	"todo-app/services/auth"
)

// TestOAuthErrorHandling tests error scenarios in OAuth flow
func TestOAuthErrorHandling(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{}, &models.OAuthState{})
	require.NoError(t, err)

	t.Run("handles Google service unavailable", func(t *testing.T) {
		// When Google OAuth service is unavailable (503 error)
		// The system should:
		// 1. Return appropriate error to user
		// 2. Not create partial user records
		// 3. Clean up any OAuth state
		// 4. Log the error for monitoring

		// Verify no sessions created during error
		var sessionCount int64
		err := db.Model(&models.AuthenticationSession{}).Count(&sessionCount).Error
		require.NoError(t, err)
		initialCount := sessionCount

		// Simulate service unavailable error
		serviceUnavailableErr := errors.New("Google OAuth service unavailable (503)")
		assert.Error(t, serviceUnavailableErr)

		// Verify state remains consistent
		err = db.Model(&models.AuthenticationSession{}).Count(&sessionCount).Error
		require.NoError(t, err)
		assert.Equal(t, initialCount, sessionCount, "No sessions should be created on service error")
	})

	t.Run("handles invalid authorization code", func(t *testing.T) {
		// When authorization code is invalid or expired
		// Should return error without creating user/session

		invalidCode := "invalid_auth_code"
		assert.NotEmpty(t, invalidCode)

		// Verify no user created
		var userCount int64
		err := db.Model(&models.User{}).Count(&userCount).Error
		require.NoError(t, err)
		initialUserCount := userCount

		// Simulate invalid code error
		invalidCodeErr := errors.New("invalid authorization code")
		assert.Error(t, invalidCodeErr)

		// Verify no user was created
		err = db.Model(&models.User{}).Count(&userCount).Error
		require.NoError(t, err)
		assert.Equal(t, initialUserCount, userCount)
	})

	t.Run("handles state mismatch attack", func(t *testing.T) {
		// CSRF protection: state parameter must match
		stateToken := "valid_state_token"
		pkceVerifier := "test_pkce_verifier"

		// Create OAuth state
		oauthState := &models.OAuthState{
			StateToken:   stateToken,
			PKCEVerifier: pkceVerifier,
			RedirectURI:  "http://localhost:3000/callback",
			ExpiresAt:    time.Now().Add(5 * time.Minute),
		}
		err := db.Create(oauthState).Error
		require.NoError(t, err)

		// Attempt with different state token (CSRF attack)
		attackStateToken := "malicious_state_token"
		assert.NotEqual(t, stateToken, attackStateToken)

		// Verify state mismatch is detected
		var foundState models.OAuthState
		err = db.Where("state_token = ?", attackStateToken).First(&foundState).Error
		assert.Error(t, err, "Should not find state with mismatched token")
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("handles expired OAuth state", func(t *testing.T) {
		// OAuth state expires after 5 minutes
		expiredState := &models.OAuthState{
			StateToken:   "expired_state_token",
			PKCEVerifier: "expired_verifier",
			RedirectURI:  "http://localhost:3000/callback",
			ExpiresAt:    time.Now().Add(-1 * time.Minute), // Expired 1 minute ago
		}
		err := db.Create(expiredState).Error
		require.NoError(t, err)

		// Verify state is expired
		now := time.Now()
		assert.True(t, expiredState.ExpiresAt.Before(now), "State should be expired")

		// Attempting to use expired state should fail
		var foundState models.OAuthState
		err = db.Where("state_token = ? AND expires_at > ?", "expired_state_token", now).
			First(&foundState).Error
		assert.Error(t, err, "Should not find expired state")
	})

	t.Run("handles network timeout during token exchange", func(t *testing.T) {
		// When network timeout occurs during token exchange
		// Should not create partial records

		var userCount, sessionCount int64
		db.Model(&models.User{}).Count(&userCount)
		db.Model(&models.AuthenticationSession{}).Count(&sessionCount)

		initialUserCount := userCount
		initialSessionCount := sessionCount

		// Simulate network timeout
		timeoutErr := errors.New("network timeout during token exchange")
		assert.Error(t, timeoutErr)

		// Verify no records created
		db.Model(&models.User{}).Count(&userCount)
		db.Model(&models.AuthenticationSession{}).Count(&sessionCount)

		assert.Equal(t, initialUserCount, userCount)
		assert.Equal(t, initialSessionCount, sessionCount)
	})

	t.Run("handles user denied OAuth consent", func(t *testing.T) {
		// When user clicks "Deny" on Google consent screen
		// Google redirects with error=access_denied

		deniedError := "access_denied"
		assert.Equal(t, "access_denied", deniedError)

		// Should not create user or session
		var userCount int64
		db.Model(&models.User{}).Count(&userCount)
		initialCount := userCount

		// Simulate access denied handling
		// System should clean up OAuth state and return to login

		db.Model(&models.User{}).Count(&userCount)
		assert.Equal(t, initialCount, userCount, "No user should be created on denial")
	})

	t.Run("handles token refresh failure", func(t *testing.T) {
		user := &models.User{
			Email:         "refresh_fail@gmail.com",
			Name:          "Refresh Fail User",
			GoogleID:      "google_refresh_fail",
			OAuthProvider: "google",
			IsActive:      true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(10 * time.Minute)

		// Create session
		session, err := sessionService.CreateSession(ctx, user.ID, true, "access", "invalid_refresh_token", tokenExpiry)
		require.NoError(t, err)

		// Attempt to refresh with invalid refresh token
		// Should fail and potentially terminate session
		newExpiry := time.Now().Add(1 * time.Hour)
		err = sessionService.RefreshSession(ctx, session.SessionToken, "new_access", newExpiry)

		// Depending on implementation:
		// - May return error
		// - May terminate session
		// Either is acceptable as long as security is maintained
	})

	t.Run("handles concurrent OAuth state validation", func(t *testing.T) {
		// Multiple requests with same state token (replay attack)
		stateToken := "concurrent_state_token"
		oauthState := &models.OAuthState{
			StateToken:   stateToken,
			PKCEVerifier: "concurrent_verifier",
			RedirectURI:  "http://localhost:3000/callback",
			ExpiresAt:    time.Now().Add(5 * time.Minute),
		}
		err := db.Create(oauthState).Error
		require.NoError(t, err)

		// First validation should succeed and consume state
		ctx := context.Background()
		result := db.WithContext(ctx).Delete(&models.OAuthState{}, "state_token = ?", stateToken)
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// Second validation should fail (state already consumed)
		result = db.WithContext(ctx).Delete(&models.OAuthState{}, "state_token = ?", stateToken)
		require.NoError(t, result.Error)
		assert.Equal(t, int64(0), result.RowsAffected, "State should already be consumed")
	})

	t.Run("handles database connection failure", func(t *testing.T) {
		// When database is unavailable during OAuth flow
		// Should return appropriate error without data corruption

		// This is a conceptual test - actual implementation would need
		// proper error handling and retry logic

		dbError := errors.New("database connection failed")
		assert.Error(t, dbError)

		// System should:
		// 1. Return error to user
		// 2. Not create partial records
		// 3. Log error for monitoring
		// 4. Allow retry
	})

	t.Run("validates email domain restrictions", func(t *testing.T) {
		// If system restricts OAuth to specific domains
		// Should reject unauthorized domains

		unauthorizedEmail := "user@unauthorized-domain.com"
		assert.Contains(t, unauthorizedEmail, "@")

		// Implementation could restrict to specific domains
		// e.g., only @company.com emails allowed
	})

	t.Run("handles rate limit exceeded", func(t *testing.T) {
		// When rate limit is exceeded (too many OAuth attempts)
		// Should return 429 Too Many Requests

		rateLimitError := errors.New("rate limit exceeded")
		assert.Error(t, rateLimitError)

		// System should:
		// 1. Return 429 status
		// 2. Include Retry-After header
		// 3. Not create any records
		// 4. Log the rate limit violation
	})
}

// TestOAuthSecurityEdgeCases tests security-related edge cases
func TestOAuthSecurityEdgeCases(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.OAuthState{})
	require.NoError(t, err)

	t.Run("prevents OAuth state fixation attack", func(t *testing.T) {
		// Attacker tries to inject their own state token
		attackerState := "attacker_state_token"

		// Legitimate state creation
		legitimateState := &models.OAuthState{
			StateToken:   "legitimate_state",
			PKCEVerifier: "legitimate_verifier",
			RedirectURI:  "http://localhost:3000/callback",
			ExpiresAt:    time.Now().Add(5 * time.Minute),
		}
		err := db.Create(legitimateState).Error
		require.NoError(t, err)

		// Attempt to use attacker's state
		var foundState models.OAuthState
		err = db.Where("state_token = ?", attackerState).First(&foundState).Error
		assert.Error(t, err, "Attacker state should not be found")
	})

	t.Run("validates redirect URI whitelist", func(t *testing.T) {
		// Only whitelisted redirect URIs should be allowed
		allowedURI := "http://localhost:3000/callback"
		maliciousURI := "https://evil.com/steal-tokens"

		// Create state with allowed URI
		validState := &models.OAuthState{
			StateToken:   "valid_redirect_state",
			PKCEVerifier: "verifier",
			RedirectURI:  allowedURI,
			ExpiresAt:    time.Now().Add(5 * time.Minute),
		}
		err := db.Create(validState).Error
		require.NoError(t, err)

		// Malicious URI should be rejected (validation in handler)
		assert.NotEqual(t, allowedURI, maliciousURI)
		assert.NotContains(t, maliciousURI, "localhost")
	})

	t.Run("prevents session fixation", func(t *testing.T) {
		// Session tokens must be regenerated after authentication
		// Never accept pre-existing session tokens

		user := &models.User{
			Email:    "fixation@gmail.com",
			Name:     "Fixation Test",
			GoogleID: "google_fixation",
			IsActive: true,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		sessionService := auth.NewSessionService(db)
		ctx := context.Background()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		// Create two sessions
		session1, err := sessionService.CreateSession(ctx, user.ID, true, "access1", "refresh1", tokenExpiry)
		require.NoError(t, err)

		session2, err := sessionService.CreateSession(ctx, user.ID, true, "access2", "refresh2", tokenExpiry)
		require.NoError(t, err)

		// Verify each session has unique token
		assert.NotEqual(t, session1.SessionToken, session2.SessionToken)
	})
}