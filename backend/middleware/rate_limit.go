package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu          sync.Mutex
	buckets     map[string]*bucket
	rate        int           // requests per window
	window      time.Duration // time window
	cleanup     time.Duration // cleanup interval
	maxBuckets  int           // maximum number of buckets to track
}

// bucket represents a rate limit bucket for a specific key
type bucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		buckets:    make(map[string]*bucket),
		rate:       rate,
		window:     window,
		cleanup:    time.Hour,
		maxBuckets: 10000,
	}

	// Start cleanup goroutine
	go limiter.cleanupExpiredBuckets()

	return limiter
}

// RateLimit middleware limits requests per IP address
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if !rl.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Header("Retry-After", rl.window.String())
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser middleware limits requests per authenticated user
func (rl *RateLimiter) RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// Fall back to IP-based rate limiting
			key := c.ClientIP()
			if !rl.allow(key) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "rate_limit_exceeded",
					"message": "Too many requests. Please try again later.",
				})
				c.Abort()
				return
			}
		} else {
			key := "user_" + userID.(string)
			if !rl.allow(key) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "rate_limit_exceeded",
					"message": "Too many requests. Please try again later.",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// allow checks if a request is allowed for the given key
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	b, exists := rl.buckets[key]
	if !exists {
		// Create new bucket
		b = &bucket{
			tokens:     rl.rate,
			lastRefill: time.Now(),
		}
		rl.buckets[key] = b
	}
	rl.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill tokens if window has passed
	now := time.Now()
	if now.Sub(b.lastRefill) >= rl.window {
		b.tokens = rl.rate
		b.lastRefill = now
	}

	// Check if tokens are available
	if b.tokens <= 0 {
		return false
	}

	// Consume a token
	b.tokens--
	return true
}

// cleanupExpiredBuckets periodically removes old buckets
func (rl *RateLimiter) cleanupExpiredBuckets() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		// Remove buckets that haven't been used for 2x the window duration
		expiry := rl.window * 2
		for key, b := range rl.buckets {
			b.mu.Lock()
			if now.Sub(b.lastRefill) > expiry {
				delete(rl.buckets, key)
			}
			b.mu.Unlock()
		}

		// If we have too many buckets, remove the oldest ones
		if len(rl.buckets) > rl.maxBuckets {
			// Simple strategy: clear half of them
			count := 0
			target := len(rl.buckets) / 2
			for key := range rl.buckets {
				delete(rl.buckets, key)
				count++
				if count >= target {
					break
				}
			}
		}

		rl.mu.Unlock()
	}
}

// OAuthRateLimiter creates a rate limiter specifically for OAuth endpoints
func OAuthRateLimiter() gin.HandlerFunc {
	// More restrictive rate limiting for OAuth endpoints
	// 10 requests per minute
	limiter := NewRateLimiter(10, time.Minute)
	return limiter.RateLimit()
}

// APIRateLimiter creates a standard rate limiter for API endpoints
func APIRateLimiter() gin.HandlerFunc {
	// Standard rate limiting for API endpoints
	// 100 requests per minute
	limiter := NewRateLimiter(100, time.Minute)
	return limiter.RateLimit()
}

// StrictRateLimiter creates a very restrictive rate limiter
func StrictRateLimiter(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)
	return limiter.RateLimit()
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	Enabled       bool
	RequestsPerWindow int
	WindowDuration    time.Duration
	ByUser        bool
}

// NewRateLimitMiddleware creates rate limiting middleware from config
func NewRateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	if !config.Enabled {
		// Return no-op middleware
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewRateLimiter(config.RequestsPerWindow, config.WindowDuration)

	if config.ByUser {
		return limiter.RateLimitByUser()
	}

	return limiter.RateLimit()
}