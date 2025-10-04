package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/storage"
	"todo-app/internal/dtos"
)

// TestGoogleSignup_NewUser_Success tests the complete OAuth flow for a new user
func TestGoogleSignup_NewUser_Success(t *testing.T) {
	// Setup test database
	if err := storage.InitDatabase(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer storage.CloseDatabase()

	// TODO: Setup mock Google OAuth server
	// mockOAuthServer := setupMockGoogleOAuthServer(t)
	// defer mockOAuthServer.Close()

	// Step 1: Initiate login flow
	loginReq, _ := http.NewRequest("GET", "/api/auth/google/login", nil)
	loginW := httptest.NewRecorder()

	// TODO: Execute request against actual router
	// router.ServeHTTP(loginW, loginReq)

	// Verify redirect to Google
	assert.Equal(t, http.StatusFound, loginW.Code)

	// Extract state cookie
	var stateCookie *http.Cookie
	for _, cookie := range loginW.Result().Cookies() {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
			break
		}
	}
	assert.NotNil(t, stateCookie)

	// Step 2: Simulate Google callback with valid verified email
	callbackURL := "/api/auth/google/callback?code=mock_valid_code&state=" + stateCookie.Value
	callbackReq, _ := http.NewRequest("GET", callbackURL, nil)
	callbackReq.AddCookie(stateCookie)

	callbackW := httptest.NewRecorder()
	// TODO: Execute request against actual router
	// router.ServeHTTP(callbackW, callbackReq)

	// Verify redirect to frontend home
	assert.Equal(t, http.StatusFound, callbackW.Code)
	location := callbackW.Header().Get("Location")
	assert.Equal(t, "http://localhost:3000/", location)

	// Verify session cookie set with 7-day expiration
	var sessionCookie *http.Cookie
	for _, cookie := range callbackW.Result().Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, 604800, sessionCookie.MaxAge, "Session should last 7 days")

	// Verify user created in database
	var user models.User
	result := storage.DB.Where("email = ?", "test@example.com").First(&user)
	assert.NoError(t, result.Error, "User should be created")
	assert.Equal(t, "google", user.AuthMethod, "User should have google auth method")

	// Verify GoogleIdentity created
	var googleIdentity models.GoogleIdentity
	result = storage.DB.Where("user_id = ?", user.ID).First(&googleIdentity)
	assert.NoError(t, result.Error, "GoogleIdentity should be created")
	assert.True(t, googleIdentity.EmailVerified, "Email should be verified")

	// Verify session created with 7-day expiration
	// Note: This assumes a sessions table exists - may need to adjust based on actual implementation
	// For now, we verify via the JWT token in the cookie
}
