package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware provides security-related middleware functions
type SecurityMiddleware struct {
	allowedOrigins []string
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware(allowedOrigins []string) *SecurityMiddleware {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:3000"}
	}
	return &SecurityMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func (m *SecurityMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range m.allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Max-Age", "86400") // 24 hours
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CSRFProtection middleware validates CSRF tokens for OAuth state
func (m *SecurityMiddleware) CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for GET and OPTIONS requests
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Check for CSRF token in header
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			// Try cookie
			csrfToken, _ = c.Cookie("csrf_token")
		}

		// Validate CSRF token
		expectedToken, exists := c.Get("csrf_token")
		if exists && expectedToken != csrfToken {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "csrf_token_invalid",
				"message": "Invalid CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OAuthStateValidation middleware validates OAuth state parameter
func (m *SecurityMiddleware) OAuthStateValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is specifically for OAuth callback
		if c.Request.URL.Path == "/api/v1/auth/google/callback" {
			state := c.Query("state")
			if state == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "missing_state",
					"message": "OAuth state parameter is required",
				})
				c.Abort()
				return
			}

			// Verify state matches cookie
			stateCookie, err := c.Cookie("oauth_state")
			if err != nil || stateCookie == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_state",
					"message": "OAuth state cookie not found",
				})
				c.Abort()
				return
			}

			if state != stateCookie {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "state_mismatch",
					"message": "OAuth state does not match",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// SecurityHeaders middleware adds security-related HTTP headers
func (m *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// Strict Transport Security (HSTS) - only enable in production with HTTPS
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// RateLimitByIP middleware implements basic rate limiting by IP address
func (m *SecurityMiddleware) RateLimitByIP() gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, use Redis or similar for distributed rate limiting
	type rateLimitEntry struct {
		count     int
		resetTime int64
	}

	// This is a simplified implementation
	// For production, use a proper rate limiting library
	return func(c *gin.Context) {
		// Rate limiting logic would go here
		// For now, this is a placeholder
		c.Next()
	}
}

// ValidateOrigin middleware validates the request origin
func (m *SecurityMiddleware) ValidateOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")

		// Allow requests without Origin/Referer (e.g., server-to-server)
		if origin == "" && referer == "" {
			c.Next()
			return
		}

		// Check if origin is allowed
		if origin != "" {
			allowed := false
			for _, allowedOrigin := range m.allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "origin_not_allowed",
					"message": "Request origin not allowed",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// PreventCSRF middleware prevents CSRF attacks by validating state tokens
func (m *SecurityMiddleware) PreventCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Check Origin header
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			allowed := false
			for _, allowedOrigin := range m.allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "csrf_origin_check_failed",
					"message": "Request origin validation failed",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}