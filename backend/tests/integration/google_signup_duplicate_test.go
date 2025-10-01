package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/storage"
	"todo-app/internal/models"
)

// TestGoogleSignup_DuplicateUser_RedirectsToLogin tests duplicate signup prevention
func TestGoogleSignup_DuplicateUser_RedirectsToLogin(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// Create existing user with Google identity
	existingUser := models.User{
		Email:      "existing@example.com",
		Name:       "Existing User",
		AuthMethod: "google",
		IsActive:   true,
	}
	storage.DB.Create(&existingUser)

	existingIdentity := models.GoogleIdentity{
		UserID:        existingUser.ID,
		GoogleUserID:  "google_123456",
		Email:         "existing@example.com",
		EmailVerified: true,
	}
	storage.DB.Create(&existingIdentity)

	// Count users before attempt
	var userCountBefore int64
	storage.DB.Model(&models.User{}).Count(&userCountBefore)

	var identityCountBefore int64
	storage.DB.Model(&models.GoogleIdentity{}).Count(&identityCountBefore)

	// TODO: Setup mock OAuth server to return the same google_user_id
	// mockOAuthServer := setupMockGoogleOAuthServerWithID(t, "google_123456")
	// defer mockOAuthServer.Close()

	// Attempt signup with same Google account
	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=duplicate_code&state=test_state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test_state",
	})

	w := httptest.NewRecorder()
	// TODO: Execute request against actual router
	// router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusFound, w.Code, "Should redirect")

	location := w.Header().Get("Location")
	assert.Equal(t, "http://localhost:3000/login", location, "Should redirect to login page, not home")

	// Verify no new user created
	var userCountAfter int64
	storage.DB.Model(&models.User{}).Count(&userCountAfter)
	assert.Equal(t, userCountBefore, userCountAfter, "No new user should be created")

	// Verify no new GoogleIdentity created
	var identityCountAfter int64
	storage.DB.Model(&models.GoogleIdentity{}).Count(&identityCountAfter)
	assert.Equal(t, identityCountBefore, identityCountAfter, "No new GoogleIdentity should be created")
}
