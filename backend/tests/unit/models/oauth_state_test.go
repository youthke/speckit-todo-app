package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"domain/auth/entities"
)

func setupOAuthStateTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the OAuthState table
	err = db.AutoMigrate(&entities.OAuthState{})
	require.NoError(t, err)

	return db
}

func TestOAuthState_Validation(t *testing.T) {
	tests := []struct {
		name        string
		state       entities.OAuthState
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid OAuth state",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12", // 40 chars
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "http://localhost:3000/dashboard",
				ExpiresAt:     time.Now().Add(5 * time.Minute),
			},
			shouldError: false,
		},
		{
			name: "invalid - short state token",
			state: entities.OAuthState{
				StateToken:    "short", // Too short
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "http://localhost:3000/dashboard",
				ExpiresAt:     time.Now().Add(5 * time.Minute),
			},
			shouldError: true,
			errorMsg:    "state_token must be at least 32 characters",
		},
		{
			name: "invalid - empty PKCE verifier",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12",
				PKCEVerifier:  "", // Empty
				RedirectURI:   "http://localhost:3000/dashboard",
				ExpiresAt:     time.Now().Add(5 * time.Minute),
			},
			shouldError: true,
			errorMsg:    "pkce_verifier cannot be empty",
		},
		{
			name: "invalid - invalid redirect URI",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12",
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "not-a-valid-url", // Invalid URL
				ExpiresAt:     time.Now().Add(5 * time.Minute),
			},
			shouldError: true,
			errorMsg:    "redirect_uri must be a valid URL",
		},
		{
			name: "invalid - redirect URI not whitelisted",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12",
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "http://malicious-site.com/steal-tokens", // Not whitelisted
				ExpiresAt:     time.Now().Add(5 * time.Minute),
			},
			shouldError: true,
			errorMsg:    "redirect_uri not in whitelist",
		},
		{
			name: "invalid - expired state",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12",
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "http://localhost:3000/dashboard",
				ExpiresAt:     time.Now().Add(-1 * time.Minute), // Expired
			},
			shouldError: true,
			errorMsg:    "state cannot be expired",
		},
		{
			name: "invalid - expires too far in future",
			state: entities.OAuthState{
				StateToken:    "abcdef1234567890abcdef1234567890abcdef12",
				PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
				RedirectURI:   "http://localhost:3000/dashboard",
				ExpiresAt:     time.Now().Add(10 * time.Minute), // Too far
			},
			shouldError: true,
			errorMsg:    "expires_at cannot exceed 5 minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.state.Validate()

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

func TestOAuthState_Create(t *testing.T) {
	db := setupOAuthStateTestDB(t)

	state := entities.OAuthState{
		StateToken:    "test_state_token_1234567890_abcdef_secure",
		PKCEVerifier:  "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
		RedirectURI:   "http://localhost:3000/dashboard",
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}

	result := db.Create(&state)
	require.NoError(t, result.Error)
	assert.NotZero(t, state.CreatedAt)
	assert.Equal(t, "test_state_token_1234567890_abcdef_secure", state.StateToken)
}

func TestOAuthState_UniqueStateToken(t *testing.T) {
	db := setupOAuthStateTestDB(t)

	// Create first state
	state1 := entities.OAuthState{
		StateToken:    "unique_state_token_abcdef1234567890_test",
		PKCEVerifier:  "first_verifier_code_challenge_test_123",
		RedirectURI:   "http://localhost:3000/dashboard",
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}

	result := db.Create(&state1)
	require.NoError(t, result.Error)

	// Try to create state with same token
	state2 := entities.OAuthState{
		StateToken:    "unique_state_token_abcdef1234567890_test", // Same token
		PKCEVerifier:  "second_verifier_code_challenge_test_456",
		RedirectURI:   "http://localhost:3000/auth/callback",
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}

	result = db.Create(&state2)
	assert.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")
}

func TestOAuthState_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		state    entities.OAuthState
		expected bool
	}{
		{
			name: "not expired state",
			state: entities.OAuthState{
				ExpiresAt: time.Now().Add(2 * time.Minute),
			},
			expected: false,
		},
		{
			name: "expired state",
			state: entities.OAuthState{
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: true,
		},
		{
			name: "state expiring now",
			state: entities.OAuthState{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOAuthState_GenerateState(t *testing.T) {
	state, err := dtos.GenerateOAuthState("http://localhost:3000/dashboard")
	require.NoError(t, err)

	// Validate generated state
	assert.NotEmpty(t, state.StateToken)
	assert.GreaterOrEqual(t, len(state.StateToken), 32)
	assert.NotEmpty(t, state.PKCEVerifier)
	assert.Equal(t, "http://localhost:3000/dashboard", state.RedirectURI)
	assert.True(t, state.ExpiresAt.After(time.Now()))
	assert.True(t, state.ExpiresAt.Before(time.Now().Add(6*time.Minute)))
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), state.ExpiresAt, 10*time.Second)
}

func TestOAuthState_GenerateState_InvalidRedirectURI(t *testing.T) {
	invalidURIs := []string{
		"not-a-url",
		"ftp://invalid-scheme.com",
		"http://malicious-site.com/attack",
		"",
	}

	for _, uri := range invalidURIs {
		t.Run("invalid_uri_"+uri, func(t *testing.T) {
			_, err := dtos.GenerateOAuthState(uri)
			assert.Error(t, err)
		})
	}
}

func TestOAuthState_ValidateRedirectURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected bool
	}{
		{
			name:     "localhost frontend",
			uri:      "http://localhost:3000/dashboard",
			expected: true,
		},
		{
			name:     "localhost root",
			uri:      "http://localhost:3000/",
			expected: true,
		},
		{
			name:     "localhost auth callback",
			uri:      "http://localhost:3000/auth/callback",
			expected: true,
		},
		{
			name:     "production domain",
			uri:      "https://todo-app.example.com/dashboard",
			expected: false, // Would be true in production with proper config
		},
		{
			name:     "malicious site",
			uri:      "http://malicious.com/steal-tokens",
			expected: false,
		},
		{
			name:     "invalid scheme",
			uri:      "ftp://localhost:3000/dashboard",
			expected: false,
		},
		{
			name:     "empty uri",
			uri:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dtos.ValidateRedirectURI(tt.uri)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOAuthState_GeneratePKCEVerifier(t *testing.T) {
	verifier1 := dtos.GeneratePKCEVerifier()
	verifier2 := dtos.GeneratePKCEVerifier()

	// Should be different each time
	assert.NotEqual(t, verifier1, verifier2)

	// Should meet PKCE requirements (43-128 characters, URL-safe)
	assert.GreaterOrEqual(t, len(verifier1), 43)
	assert.LessOrEqual(t, len(verifier1), 128)
	assert.GreaterOrEqual(t, len(verifier2), 43)
	assert.LessOrEqual(t, len(verifier2), 128)

	// Should only contain URL-safe characters
	urlSafeRegex := `^[A-Za-z0-9\-._~]+$`
	assert.Regexp(t, urlSafeRegex, verifier1)
	assert.Regexp(t, urlSafeRegex, verifier2)
}

func TestOAuthState_CleanupExpired(t *testing.T) {
	db := setupOAuthStateTestDB(t)

	// Create expired state
	expiredState := entities.OAuthState{
		StateToken:    "expired_state_token_1234567890_old",
		PKCEVerifier:  "expired_verifier_code",
		RedirectURI:   "http://localhost:3000/dashboard",
		ExpiresAt:     time.Now().Add(-10 * time.Minute), // Expired
	}
	result := db.Create(&expiredState)
	require.NoError(t, result.Error)

	// Create valid state
	validState := entities.OAuthState{
		StateToken:    "valid_state_token_1234567890_current",
		PKCEVerifier:  "valid_verifier_code",
		RedirectURI:   "http://localhost:3000/dashboard",
		ExpiresAt:     time.Now().Add(3 * time.Minute), // Valid
	}
	result = db.Create(&validState)
	require.NoError(t, result.Error)

	// Count before cleanup
	var beforeCount int64
	db.Model(&entities.OAuthState{}).Count(&beforeCount)
	assert.Equal(t, int64(2), beforeCount)

	// Cleanup expired states
	deletedCount := entities.CleanupExpiredOAuthStates(db)
	assert.Equal(t, int64(1), deletedCount)

	// Count after cleanup
	var afterCount int64
	db.Model(&entities.OAuthState{}).Count(&afterCount)
	assert.Equal(t, int64(1), afterCount)

	// Verify only valid state remains
	var remainingState entities.OAuthState
	result = db.First(&remainingState)
	require.NoError(t, result.Error)
	assert.Equal(t, "valid_state_token_1234567890_current", remainingState.StateToken)
}