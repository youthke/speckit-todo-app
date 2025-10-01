package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/storage"
	"todo-app/internal/models"
)

// TestGoogleSignup_UnverifiedEmail_Rejected tests that unverified emails are rejected
func TestGoogleSignup_UnverifiedEmail_Rejected(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// Count users before attempt
	var userCountBefore int64
	storage.DB.Model(&models.User{}).Count(&userCountBefore)

	// TODO: Setup mock OAuth server to return email_verified=false
	// mockOAuthServer := setupMockGoogleOAuthServerUnverified(t)
	// defer mockOAuthServer.Close()

	// Attempt signup with unverified email
	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=unverified_code&state=test_state", nil)
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
	assert.Contains(t, location, "/signup?error=authentication_failed", "Should redirect to signup with error")

	// Verify no user created
	var userCountAfter int64
	storage.DB.Model(&models.User{}).Count(&userCountAfter)
	assert.Equal(t, userCountBefore, userCountAfter, "No user should be created")

	// Verify no GoogleIdentity created
	var identityCount int64
	storage.DB.Model(&models.GoogleIdentity{}).Count(&identityCount)
	assert.Equal(t, int64(0), identityCount, "No GoogleIdentity should be created")
}
