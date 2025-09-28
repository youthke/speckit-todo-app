package handlers

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Printf("Panic occurred: %v\n%s", err, debug.Stack())

				// Return generic error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "internal_error",
					"message": "An internal server error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RequestLogger middleware logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log request details
		log.Printf(
			"[%s] %s %s | %d | %v | %s",
			c.Request.Method,
			path,
			func() string {
				if raw != "" {
					return "?" + raw
				}
				return ""
			}(),
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}

// ValidationErrorHandler handles validation errors
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			// Log the error
			log.Printf("Request validation error: %v", c.Errors)

			// Return validation error response
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": "Request validation failed",
				"details": c.Errors.Errors(),
			})
		}
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}