package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/internal/models"
)

// TestAuthMe_GoogleUser tests that /api/auth/me returns correct data for Google OAuth user
func TestAuthMe_GoogleUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// router.GET("/api/auth/me", middleware.AuthRequired(), handlers.GetCurrentUser)

	// TODO: Create test user with auth_method="google" in test database
	// TODO: Create valid session token for the test user

	// For now, create a mock request with a session token
	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "mock_jwt_token",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response models.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Should parse JSON response")

	assert.Equal(t, "google", response.OAuthProvider, "User should have Google OAuth provider")
	assert.NotEmpty(t, response.Email, "User should have email")
	assert.True(t, response.IsActive, "User should be active")
}

// TestAuthMe_Unauthorized tests that /api/auth/me returns 401 without valid session
func TestAuthMe_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// router.GET("/api/auth/me", middleware.AuthRequired(), handlers.GetCurrentUser)

	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	// No session cookie provided

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized")
}
