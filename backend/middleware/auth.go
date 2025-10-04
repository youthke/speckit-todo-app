package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"todo-app/internal/dtos"
	"todo-app/services/auth"
)

// AuthMiddleware creates a middleware for OAuth session validation
type AuthMiddleware struct {
	sessionService *auth.SessionService
	jwtService     *auth.JWTService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(sessionService *auth.SessionService, jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
		jwtService:     jwtService,
	}
}

// RequireAuth middleware requires valid authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := m.extractToken(c)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "no_auth_token",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Validate session
		result, err := m.sessionService.ValidateSession(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "validation_error",
				"message": "Failed to validate session",
			})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_session",
				"message": result.Error,
			})
			c.Abort()
			return
		}

		// Set user and session in context
		c.Set("user", result.User)
		c.Set("session", result.Session)

		// Extract user ID from interface
		if user, ok := result.User.(*dtos.User); ok {
			c.Set("user_id", user.ID)
		}
		c.Set("session_id", result.Session.ID)

		c.Next()
	}
}

// OptionalAuth middleware validates authentication if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := m.extractToken(c)

		if tokenString != "" {
			// Validate session
			result, err := m.sessionService.ValidateSession(tokenString)
			if err == nil && result.Valid {
				// Set user and session in context
				c.Set("user", result.User)
				c.Set("session", result.Session)

				// Extract user ID from interface
				if user, ok := result.User.(*dtos.User); ok {
					c.Set("user_id", user.ID)
				}
				c.Set("session_id", result.Session.ID)
			}
		}

		c.Next()
	}
}

// RequireOAuth middleware requires OAuth authentication specifically
func (m *AuthMiddleware) RequireOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := m.extractToken(c)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "no_auth_token",
				"message": "OAuth authentication required",
			})
			c.Abort()
			return
		}

		// Validate JWT token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid authentication token",
			})
			c.Abort()
			return
		}

		// Check if it's an OAuth session
		if !claims.IsOAuth {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "oauth_required",
				"message": "This endpoint requires OAuth authentication",
			})
			c.Abort()
			return
		}

		// Validate full session
		result, err := m.sessionService.ValidateSession(tokenString)
		if err != nil || !result.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_session",
				"message": "Session validation failed",
			})
			c.Abort()
			return
		}

		// Set user and session in context
		c.Set("user", result.User)
		c.Set("session", result.Session)

		// Extract user ID from interface
		if user, ok := result.User.(*dtos.User); ok {
			c.Set("user_id", user.ID)
		}
		c.Set("session_id", result.Session.ID)

		c.Next()
	}
}

// RefreshIfNeeded middleware automatically refreshes OAuth tokens if needed
func (m *AuthMiddleware) RefreshIfNeeded() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := m.extractToken(c)

		if tokenString == "" {
			c.Next()
			return
		}

		// Validate session
		result, err := m.sessionService.ValidateSession(tokenString)
		if err != nil || !result.Valid {
			c.Next()
			return
		}

		// Check if refresh is needed
		if result.NeedsRefresh && result.Session.IsOAuthSession() {
			// Attempt to refresh tokens
			// Note: This is async and won't block the request
			go func() {
				// Refresh is handled by a separate endpoint
				// This just sets a flag in the response
			}()

			c.Header("X-Token-Refresh-Needed", "true")
		}

		// Set user and session in context
		c.Set("user", result.User)
		c.Set("session", result.Session)

		// Extract user ID from interface
		if user, ok := result.User.(*dtos.User); ok {
			c.Set("user_id", user.ID)
		}
		c.Set("session_id", result.Session.ID)

		c.Next()
	}
}

// extractToken extracts the authentication token from cookie or Authorization header
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// Try cookie first
	token, err := c.Cookie("session_token")
	if err == nil && token != "" {
		return token
	}

	// Try Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Extract Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(c *gin.Context) interface{} {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user
}

// GetCurrentUserID retrieves the current user ID from context
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	if id, ok := userID.(uint); ok {
		return id, true
	}
	return 0, false
}

// GetCurrentSession retrieves the current session from context
func GetCurrentSession(c *gin.Context) interface{} {
	session, exists := c.Get("session")
	if !exists {
		return nil
	}
	return session
}

// GetCurrentSessionID retrieves the current session ID from context
func GetCurrentSessionID(c *gin.Context) (string, bool) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		return "", false
	}
	if id, ok := sessionID.(string); ok {
		return id, true
	}
	return "", false
}