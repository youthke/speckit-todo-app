package utils

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CookieConfig holds configuration for secure cookies
type CookieConfig struct {
	Secure   bool
	HttpOnly bool
	SameSite string
	Domain   string
	Path     string
}

// GetDefaultCookieConfig returns the default cookie configuration from environment
func GetDefaultCookieConfig() CookieConfig {
	secure := os.Getenv("SESSION_COOKIE_SECURE") == "true"
	httpOnly := os.Getenv("SESSION_COOKIE_HTTPONLY") != "false" // Default true
	sameSite := os.Getenv("SESSION_COOKIE_SAMESITE")
	if sameSite == "" {
		sameSite = "Lax"
	}

	return CookieConfig{
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		Domain:   "",
		Path:     "/",
	}
}

// SetSessionCookie sets a secure session cookie
func SetSessionCookie(c *gin.Context, name, value string, maxAge int) {
	config := GetDefaultCookieConfig()

	sameSite := http.SameSiteLaxMode
	switch config.SameSite {
	case "Strict":
		sameSite = http.SameSiteStrictMode
	case "None":
		sameSite = http.SameSiteNoneMode
	case "Lax":
		sameSite = http.SameSiteLaxMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		name,
		value,
		maxAge,
		config.Path,
		config.Domain,
		config.Secure,
		config.HttpOnly,
	)
}

// SetOAuthStateCookie sets a secure cookie for OAuth state
func SetOAuthStateCookie(c *gin.Context, stateToken string) {
	config := GetDefaultCookieConfig()

	// OAuth state cookies are short-lived (5 minutes)
	maxAge := 300

	sameSite := http.SameSiteLaxMode
	if config.SameSite == "Strict" {
		sameSite = http.SameSiteStrictMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		"oauth_state",
		stateToken,
		maxAge,
		config.Path,
		config.Domain,
		config.Secure,
		config.HttpOnly,
	)
}

// SetCSRFCookie sets a CSRF token cookie
func SetCSRFCookie(c *gin.Context, csrfToken string) {
	config := GetDefaultCookieConfig()

	// CSRF cookies should not be HttpOnly so JavaScript can read them
	// but should still be secure
	sameSite := http.SameSiteStrictMode // CSRF tokens should use Strict

	c.SetSameSite(sameSite)
	c.SetCookie(
		"csrf_token",
		csrfToken,
		3600, // 1 hour
		config.Path,
		config.Domain,
		config.Secure,
		false, // Not HttpOnly for CSRF tokens
	)
}

// ClearSessionCookie clears the session cookie
func ClearSessionCookie(c *gin.Context) {
	config := GetDefaultCookieConfig()

	c.SetCookie(
		"session_token",
		"",
		-1, // Expire immediately
		config.Path,
		config.Domain,
		config.Secure,
		config.HttpOnly,
	)
}

// ClearOAuthStateCookie clears the OAuth state cookie
func ClearOAuthStateCookie(c *gin.Context) {
	config := GetDefaultCookieConfig()

	c.SetCookie(
		"oauth_state",
		"",
		-1,
		config.Path,
		config.Domain,
		config.Secure,
		config.HttpOnly,
	)
}

// ClearAllAuthCookies clears all authentication-related cookies
func ClearAllAuthCookies(c *gin.Context) {
	ClearSessionCookie(c)
	ClearOAuthStateCookie(c)

	config := GetDefaultCookieConfig()

	// Clear CSRF cookie
	c.SetCookie(
		"csrf_token",
		"",
		-1,
		config.Path,
		config.Domain,
		config.Secure,
		false,
	)
}

// GetSessionExpiry returns the session expiry duration from environment
func GetSessionExpiry() time.Duration {
	hoursStr := os.Getenv("JWT_EXPIRES_HOURS")
	if hoursStr == "" {
		return 24 * time.Hour // Default 24 hours
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return 24 * time.Hour
	}

	return time.Duration(hours) * time.Hour
}

// GetOAuthStateExpiry returns the OAuth state expiry duration
func GetOAuthStateExpiry() time.Duration {
	minutesStr := os.Getenv("OAUTH_STATE_EXPIRES_MINUTES")
	if minutesStr == "" {
		return 5 * time.Minute // Default 5 minutes
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return 5 * time.Minute
	}

	return time.Duration(minutes) * time.Minute
}

// ValidateCookieDomain validates a cookie domain
func ValidateCookieDomain(domain string) bool {
	// Basic validation - in production, use more sophisticated validation
	if domain == "" {
		return true // Empty domain is valid (current domain)
	}

	// Domain should not start with a dot for explicit domains
	// but can start with dot for subdomains
	return len(domain) > 0
}

// GetSecureCookieSettings returns secure cookie settings based on environment
func GetSecureCookieSettings() map[string]interface{} {
	config := GetDefaultCookieConfig()

	return map[string]interface{}{
		"secure":    config.Secure,
		"http_only": config.HttpOnly,
		"same_site": config.SameSite,
		"domain":    config.Domain,
		"path":      config.Path,
		"max_age":   int(GetSessionExpiry().Seconds()),
	}
}