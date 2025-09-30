package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-app/handlers"
	"todo-app/services/auth"
)

func TestAuthGoogleLogin_Contract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return auth URL with 200 status", func(t *testing.T) {
		// Setup
		router := gin.New()

		// Mock OAuth service (in real implementation, would use actual service)
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.GET("/api/v1/auth/google/login", authHandler.GoogleLogin)

		// Execute
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Should return 200 with auth_url field
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
		assert.Contains(t, w.Body.String(), "auth_url")

		// Assert - Contract: Should set oauth_state cookie
		cookies := w.Result().Cookies()
		foundStateCookie := false
		for _, cookie := range cookies {
			if cookie.Name == "oauth_state" {
				foundStateCookie = true
				assert.True(t, cookie.HttpOnly, "oauth_state cookie should be HttpOnly")
				assert.Equal(t, 300, cookie.MaxAge, "oauth_state cookie should expire in 5 minutes")
			}
		}
		assert.True(t, foundStateCookie, "oauth_state cookie should be set")
	})

	t.Run("should handle optional redirect_uri parameter", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.GET("/api/v1/auth/google/login", authHandler.GoogleLogin)

		// Execute
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login?redirect_uri=http://localhost:3000/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "auth_url")
	})

	t.Run("should reject invalid redirect_uri", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.GET("/api/v1/auth/google/login", authHandler.GoogleLogin)

		// Execute with malicious redirect
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login?redirect_uri=https://evil.com", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Should return 400 for invalid redirect_uri
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "message")
	})

	t.Run("contract validation: response schema", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.GET("/api/v1/auth/google/login", authHandler.GoogleLogin)

		// Execute
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Response must match OpenAPI schema
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify JSON structure matches contract
		body := w.Body.String()
		assert.Contains(t, body, "auth_url", "Response must contain auth_url field")

		// Verify auth_url is a valid URL format
		assert.Contains(t, body, "https://", "auth_url should be HTTPS URL")
		assert.Contains(t, body, "accounts.google.com", "auth_url should point to Google OAuth")
	})
}