package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"domain/auth/entities"
	"todo-app/internal/dtos"
)

func setupSessionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the tables
	err = db.AutoMigrate(&dtos.User{}, &entities.AuthenticationSession{})
	require.NoError(t, err)

	return db
}

func createTestUser(t *testing.T, db *gorm.DB) dtos.User {
	user := dtos.User{
		Email:        "session_test@example.com",
		Name:         "Session Test User",
		PasswordHash: "hashed_password",
		IsActive:     true,
	}

	result := db.Create(&user)
	require.NoError(t, result.Error)
	return user
}

func TestAuthenticationSession_Validation(t *testing.T) {
	db := setupSessionTestDB(t)
	user := createTestUser(t, db)

	tests := []struct {
		name        string
		session     entities.AuthenticationSession
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid session",
			session: entities.AuthenticationSession{
				UserID:           user.ID,
				SessionToken:     "valid.jwt.token",
				SessionExpiresAt: time.Now().Add(24 * time.Hour),
				LastActivity:     time.Now(),
				UserAgent:       "Mozilla/5.0 Test Browser",
				IPAddress:       "192.168.1.1",
			},
			shouldError: false,
		},
		{
			name: "valid OAuth session",
			session: func() entities.AuthenticationSession {
				tokenExp := time.Now().Add(1 * time.Hour)
				return entities.AuthenticationSession{
					UserID:           user.ID,
					SessionToken:     "valid.oauth.jwt.token",
					RefreshToken:     "encrypted_refresh_token",
					AccessToken:      "encrypted_access_token",
					TokenExpiresAt:   &tokenExp,
					SessionExpiresAt: time.Now().Add(24 * time.Hour),
					LastActivity:     time.Now(),
					UserAgent:       "Mozilla/5.0 OAuth Browser",
					IPAddress:       "192.168.1.2",
				}
			}(),
			shouldError: false,
		},
		{
			name: "invalid - missing user_id",
			session: entities.AuthenticationSession{
				SessionToken:     "valid.jwt.token",
				SessionExpiresAt: time.Now().Add(24 * time.Hour),
				LastActivity:     time.Now(),
			},
			shouldError: true,
			errorMsg:    "user_id is required",
		},
		{
			name: "invalid - empty session_token",
			session: entities.AuthenticationSession{
				UserID:           user.ID,
				SessionExpiresAt: time.Now().Add(24 * time.Hour),
				LastActivity:     time.Now(),
			},
			shouldError: true,
			errorMsg:    "session_token cannot be empty",
		},
		{
			name: "invalid - expired session",
			session: entities.AuthenticationSession{
				UserID:           user.ID,
				SessionToken:     "expired.jwt.token",
				SessionExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
				LastActivity:     time.Now(),
			},
			shouldError: true,
			errorMsg:    "session cannot be expired",
		},
		{
			name: "invalid - access_token without token_expires_at",
			session: entities.AuthenticationSession{
				UserID:           user.ID,
				SessionToken:     "oauth.jwt.token",
				AccessToken:      "encrypted_access_token",
				SessionExpiresAt: time.Now().Add(24 * time.Hour),
				LastActivity:     time.Now(),
			},
			shouldError: true,
			errorMsg:    "token_expires_at required when access_token present",
		},
		{
			name: "invalid - future session expiry beyond 24h",
			session: entities.AuthenticationSession{
				UserID:           user.ID,
				SessionToken:     "long.jwt.token",
				SessionExpiresAt: time.Now().Add(25 * time.Hour), // Too long
				LastActivity:     time.Now(),
			},
			shouldError: true,
			errorMsg:    "session_expires_at cannot exceed 24 hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()

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

func TestAuthenticationSession_CreateSession(t *testing.T) {
	db := setupSessionTestDB(t)
	user := createTestUser(t, db)

	session := entities.AuthenticationSession{
		UserID:           user.ID,
		SessionToken:     "test.jwt.token.12345",
		SessionExpiresAt: time.Now().Add(24 * time.Hour),
		LastActivity:     time.Now(),
		UserAgent:       "Test Browser",
		IPAddress:       "127.0.0.1",
	}

	result := db.Create(&session)
	require.NoError(t, result.Error)
	assert.NotEmpty(t, session.ID)
	assert.NotZero(t, session.CreatedAt)
}

func TestAuthenticationSession_CreateOAuthSession(t *testing.T) {
	db := setupSessionTestDB(t)
	user := createTestUser(t, db)

	session := entities.AuthenticationSession{
		UserID:           user.ID,
		SessionToken:     "oauth.jwt.token.67890",
		RefreshToken:     "encrypted_refresh_token_abc123",
		AccessToken:      "encrypted_access_token_def456",
		TokenExpiresAt:   time.Now().Add(1 * time.Hour),
		SessionExpiresAt: time.Now().Add(24 * time.Hour),
		LastActivity:     time.Now(),
		UserAgent:       "OAuth Test Browser",
		IPAddress:       "192.168.1.100",
	}

	result := db.Create(&session)
	require.NoError(t, result.Error)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, "encrypted_refresh_token_abc123", session.RefreshToken)
	assert.Equal(t, "encrypted_access_token_def456", session.AccessToken)
	assert.NotNil(t, session.TokenExpiresAt)
}

func TestAuthenticationSession_UniqueSessionToken(t *testing.T) {
	db := setupSessionTestDB(t)
	user := createTestUser(t, db)

	// Create first session
	session1 := entities.AuthenticationSession{
		UserID:           user.ID,
		SessionToken:     "unique.session.token",
		SessionExpiresAt: time.Now().Add(24 * time.Hour),
		LastActivity:     time.Now(),
	}

	result := db.Create(&session1)
	require.NoError(t, result.Error)

	// Try to create session with same token
	session2 := entities.AuthenticationSession{
		UserID:           user.ID,
		SessionToken:     "unique.session.token", // Same token
		SessionExpiresAt: time.Now().Add(24 * time.Hour),
		LastActivity:     time.Now(),
	}

	result = db.Create(&session2)
	assert.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")
}

func TestAuthenticationSession_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		session  entities.AuthenticationSession
		expected bool
	}{
		{
			name: "not expired session",
			session: entities.AuthenticationSession{
				SessionExpiresAt: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "expired session",
			session: entities.AuthenticationSession{
				SessionExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "session expiring now",
			session: entities.AuthenticationSession{
				SessionExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthenticationSession_IsTokenExpired(t *testing.T) {
	tests := []struct {
		name     string
		session  entities.AuthenticationSession
		expected bool
	}{
		{
			name: "no token expiry set",
			session: entities.AuthenticationSession{
				AccessToken: "some_token",
			},
			expected: false,
		},
		{
			name: "token not expired",
			session: entities.AuthenticationSession{
				AccessToken:    "some_token",
				TokenExpiresAt: time.Now().Add(30 * time.Minute),
			},
			expected: false,
		},
		{
			name: "token expired",
			session: entities.AuthenticationSession{
				AccessToken:    "some_token",
				TokenExpiresAt: time.Now().Add(-10 * time.Minute),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsTokenExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthenticationSession_NeedsRefresh(t *testing.T) {
	tests := []struct {
		name     string
		session  entities.AuthenticationSession
		expected bool
	}{
		{
			name: "no OAuth tokens - no refresh needed",
			session: entities.AuthenticationSession{
				SessionExpiresAt: time.Now().Add(12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "OAuth tokens not expiring soon",
			session: entities.AuthenticationSession{
				AccessToken:      "token",
				TokenExpiresAt:   time.Now().Add(30 * time.Minute),
				SessionExpiresAt: time.Now().Add(12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "OAuth tokens expiring soon",
			session: entities.AuthenticationSession{
				AccessToken:      "token",
				TokenExpiresAt:   time.Now().Add(2 * time.Minute),
				SessionExpiresAt: time.Now().Add(12 * time.Hour),
			},
			expected: true,
		},
		{
			name: "OAuth tokens expired",
			session: entities.AuthenticationSession{
				AccessToken:      "token",
				TokenExpiresAt:   time.Now().Add(-5 * time.Minute),
				SessionExpiresAt: time.Now().Add(12 * time.Hour),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.NeedsRefresh()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthenticationSession_UpdateActivity(t *testing.T) {
	session := entities.AuthenticationSession{
		LastActivity: time.Now().Add(-1 * time.Hour),
	}

	oldActivity := session.LastActivity

	session.UpdateActivity()

	assert.True(t, session.LastActivity.After(oldActivity))
	assert.WithinDuration(t, time.Now(), session.LastActivity, 1*time.Second)
}