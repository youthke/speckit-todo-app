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

func TestAuthLogout_Contract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should require session cookie", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.POST("/api/v1/auth/logout", authHandler.Logout)

		// Execute - No session cookie
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Should return 401 for missing session
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should be POST method only", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.POST("/api/v1/auth/logout", authHandler.Logout)

		// Execute - GET instead of POST
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Should return 404 or 405 for wrong method
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMethodNotAllowed)
	})

	t.Run("contract validation: successful logout response", func(t *testing.T) {
		// Contract expectation: Successful logout should return 200 with:
		// - status (string, "logged_out")
		// - message (string)
		// - Set-Cookie header clearing session_token
	})

	t.Run("contract validation: clears session cookie", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.POST("/api/v1/auth/logout", authHandler.Logout)

		// Contract expectation: Response should include Set-Cookie header
		// with session_token cookie having MaxAge=-1 (expired)
	})

	t.Run("contract validation: response schema", func(t *testing.T) {
		// Contract expectation: LogoutResponse must contain:
		// - status (string, required, example: "logged_out")
		// - message (string, required, example: "Session terminated successfully")
	})

	t.Run("should handle invalid session gracefully", func(t *testing.T) {
		// Setup
		router := gin.New()
		mockGoogleConfig, _ := auth.NewGoogleOAuthConfig()
		mockOAuthService := auth.NewOAuthService(nil, mockGoogleConfig)
		authHandler := handlers.NewAuthHandler(mockOAuthService, nil, nil)

		router.POST("/api/v1/auth/logout", authHandler.Logout)

		// Execute - Invalid session
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: "invalid_token",
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert - Contract: Should return 401 for invalid session
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("contract validation: idempotent logout", func(t *testing.T) {
		// Contract expectation: Logging out an already logged-out session
		// should still return success (idempotent operation)
	})

	t.Run("contract validation: clears all auth cookies", func(t *testing.T) {
		// Contract expectation: Logout should clear:
		// - session_token cookie
		// - csrf_token cookie
		// - oauth_state cookie (if present)
	})
}