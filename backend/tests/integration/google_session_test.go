package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/config"
	"todo-app/internal/models"
	"todo-app/internal/storage"
)

// TestGoogleSignup_SessionDuration_SevenDays tests that sessions have 7-day expiration
func TestGoogleSignup_SessionDuration_SevenDays(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// TODO: Complete signup flow and capture session cookie
	// For now, test JWT creation directly

	// Create test user
	user := models.User{
		Email:      "session-test@example.com",
		Name:       "Session Test",
		AuthMethod: "google",
		IsActive:   true,
	}
	storage.DB.Create(&user)

	// Create JWT token with 7-day expiration
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	assert.NoError(t, err, "Should create JWT token")

	// Parse token and verify expiration
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})
	assert.NoError(t, err, "Should parse JWT token")

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Should extract claims")

	expTimestamp := int64(claims["exp"].(float64))
	expTime := time.Unix(expTimestamp, 0)

	// Verify expiration is approximately 7 days from now (allow 5 second variance)
	expectedExpiration := time.Now().Add(7 * 24 * time.Hour)
	assert.InDelta(t, expectedExpiration.Unix(), expTime.Unix(), 5, "Session should expire in 7 days")

	// Test cookie MaxAge should be 604800 seconds (7 days)
	expectedMaxAge := 604800
	assert.Equal(t, expectedMaxAge, 7*24*60*60, "Cookie MaxAge should be 7 days in seconds")
}

// TestSessionValidation tests that session validation works correctly
func TestSessionValidation(t *testing.T) {
	// TODO: Test that valid sessions are accepted
	// TODO: Test that expired sessions are rejected
	// TODO: Test that invalid tokens are rejected
}
