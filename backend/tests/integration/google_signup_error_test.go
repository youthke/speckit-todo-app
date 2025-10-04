package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/storage"
	"todo-app/internal/dtos"
)

// TestGoogleSignup_OAuthDenied_ShowsError tests error handling when user denies OAuth
func TestGoogleSignup_OAuthDenied_ShowsError(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// Simulate OAuth denial (Google redirects with error parameter)
	req, _ := http.NewRequest("GET", "/api/auth/google/callback?error=access_denied&state=test_state", nil)
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
	assert.Contains(t, location, "/signup?error=authentication_failed", "Should show generic error message")

	// Verify no user created
	var userCount int64
	storage.DB.Model(&models.User{}).Count(&userCount)
	assert.Equal(t, int64(0), userCount, "No user should be created")
}

// TestGoogleSignup_NetworkError_ShowsError tests error handling for network issues
func TestGoogleSignup_NetworkError_ShowsError(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// TODO: Setup mock OAuth server that times out or returns error
	// mockOAuthServer := setupMockGoogleOAuthServerWithError(t)
	// defer mockOAuthServer.Close()

	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=network_error_code&state=test_state", nil)
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
	assert.Contains(t, location, "/signup?error=authentication_failed", "Should show generic error message")
}
