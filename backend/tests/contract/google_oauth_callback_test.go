package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGoogleOAuthCallback_Success tests successful OAuth callback
func TestGoogleOAuthCallback_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// router.GET("/api/auth/google/callback", handlers.GoogleCallback)

	// Mock request with code and state parameters
	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=test_code&state=test_state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test_state",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusFound, w.Code, "Should return 302 redirect")

	location := w.Header().Get("Location")
	assert.Contains(t, location, "http://localhost:3000", "Should redirect to frontend")

	// Check for session cookie
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	assert.NotNil(t, sessionCookie, "Should set session_token cookie")
	assert.Equal(t, 604800, sessionCookie.MaxAge, "Session should last 7 days (604800 seconds)")
	assert.True(t, sessionCookie.HttpOnly, "Session cookie should be HttpOnly")
}

// TestGoogleOAuthCallback_InvalidState tests callback with invalid state parameter
func TestGoogleOAuthCallback_InvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// router.GET("/api/auth/google/callback", handlers.GoogleCallback)

	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=test_code&state=wrong_state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "correct_state",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code, "Should redirect on error")

	location := w.Header().Get("Location")
	assert.Contains(t, location, "/signup?error=authentication_failed", "Should redirect to signup with error")
}

// TestGoogleOAuthCallback_UnverifiedEmail tests callback with unverified email
func TestGoogleOAuthCallback_UnverifiedEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// This test requires mocking the OAuth service to return email_verified=false
	// router.GET("/api/auth/google/callback", handlers.GoogleCallback)

	req, _ := http.NewRequest("GET", "/api/auth/google/callback?code=unverified_email_code&state=test_state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test_state",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code, "Should redirect on error")

	location := w.Header().Get("Location")
	assert.Contains(t, location, "/signup?error=authentication_failed", "Should redirect to signup with error")
}
