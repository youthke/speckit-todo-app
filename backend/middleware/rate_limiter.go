package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiting per IP address
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit // requests per second
	b   int        // bucket size (burst)
}

// NewIPRateLimiter creates a new IP-based rate limiter
// r: rate limit (requests per second)
// b: burst size (maximum tokens in bucket)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}

	// Start cleanup goroutine to remove inactive IPs
	go limiter.cleanupInactive()

	return limiter
}

// GetLimiter returns the rate limiter for the given IP
// Creates a new limiter if one doesn't exist
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// cleanupInactive removes IP entries that haven't been used in 30 minutes
// Prevents memory leak from IP accumulation
func (i *IPRateLimiter) cleanupInactive() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		// Note: In production, you'd want to track last access time
		// For now, we'll just periodically clear the map
		// A more sophisticated approach would store access timestamps
		if len(i.ips) > 10000 {
			// Reset if too many IPs (memory protection)
			i.ips = make(map[string]*rate.Limiter)
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func (i *IPRateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := i.GetLimiter(ip)

		if !limiter.Allow() {
			// Calculate retry-after in seconds
			reservation := limiter.Reserve()
			delay := reservation.DelayFrom(time.Now())
			retryAfter := int(delay.Seconds())
			reservation.Cancel() // Cancel the reservation since we're rejecting

			c.Header("Retry-After", string(rune(retryAfter)))

			// Check if client wants JSON response
			if c.GetHeader("Accept") == "application/json" {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":       "rate_limit_exceeded",
					"message":     "Too many signup attempts. Please try again later.",
					"retry_after": retryAfter,
				})
				c.Abort()
				return
			}

			// For HTML clients (browser), redirect to signup page with error
			c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=rate_limit_exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
