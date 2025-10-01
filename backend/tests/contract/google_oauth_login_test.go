package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGoogleOAuthLogin tests the /api/auth/google/login endpoint
func TestGoogleOAuthLogin(t *testing.T) {
	// Setup test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// TODO: Register the actual handler once implemented
	// For now, this test will fail
	// router.GET("/api/auth/google/login", handlers.GoogleLogin)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/auth/google/login", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusFound, w.Code, "Should return 302 redirect")

	location := w.Header().Get("Location")
	assert.Contains(t, location, "accounts.google.com", "Should redirect to Google OAuth")

	cookies := w.Result().Cookies()
	var oauthStateCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			oauthStateCookie = cookie
			break
		}
	}

	assert.NotNil(t, oauthStateCookie, "Should set oauth_state cookie")
	assert.True(t, oauthStateCookie.HttpOnly, "oauth_state cookie should be HttpOnly")
}
