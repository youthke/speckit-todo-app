package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-app/internal/dtos"
	"todo-app/internal/services"
)

// GoogleOAuthHandler handles Google OAuth signup/login requests
type GoogleOAuthHandler struct {
	oauthService   *services.GoogleOAuthService
	sessionService *services.SessionService
}

// NewGoogleOAuthHandler creates a new Google OAuth handler
func NewGoogleOAuthHandler(db *gorm.DB) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{
		oauthService:   services.NewGoogleOAuthService(db),
		sessionService: services.NewSessionService(),
	}
}

// GoogleLogin initiates the Google OAuth flow
// GET /api/v1/auth/google/login
func (h *GoogleOAuthHandler) GoogleLogin(c *gin.Context) {
	// Generate random state token for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		log.Printf("Failed to generate state token: %v", err)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Store state in session cookie (10 min expiration for the OAuth flow)
	c.SetCookie(
		"oauth_state",
		state,
		600, // 10 minutes
		"/",
		"",
		false, // Secure (set to true in production with HTTPS)
		true,  // HttpOnly
	)

	// Generate OAuth URL
	url := h.oauthService.GenerateAuthURL(state)

	// Redirect to Google OAuth
	c.Redirect(http.StatusFound, url)
}

// GoogleCallback handles the OAuth callback from Google
// GET /api/v1/auth/google/callback
func (h *GoogleOAuthHandler) GoogleCallback(c *gin.Context) {
	// Validate state parameter (CSRF protection)
	code := c.Query("code")
	state := c.Query("state")

	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		log.Printf("State validation failed: %v", err)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Handle OAuth error (user denied permission)
	if c.Query("error") != "" {
		log.Printf("OAuth error: %s", c.Query("error"))
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Exchange code for user info
	userInfo, err := h.oauthService.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Validate email is present
	if userInfo.Email == "" {
		log.Printf("No email provided by Google for user: %s", userInfo.GoogleUserID)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Validate email is verified
	if !userInfo.EmailVerified {
		log.Printf("Email not verified for user: %s", userInfo.Email)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Check for duplicate (existing Google account)
	existingUser, err := h.oauthService.FindUserByGoogleID(userInfo.GoogleUserID)
	if err != nil {
		log.Printf("Error checking for existing user: %v", err)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	var user *dtos.User
	// If user already exists, auto-login (create new session)
	if existingUser != nil {
		log.Printf("User already exists with Google ID: %s, auto-logging in", userInfo.GoogleUserID)
		user = existingUser
	} else {
		// Create new user from Google info
		var err error
		user, err = h.oauthService.CreateUserFromGoogle(userInfo)
		if err != nil {
			log.Printf("Failed to create user from Google: %v", err)
			c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
			return
		}
	}

	// Create session token
	token, err := h.sessionService.CreateSession(user.ID)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		c.Redirect(http.StatusFound, "http://localhost:3000/signup?error=authentication_failed")
		return
	}

	// Set session cookie with 7-day expiration
	c.SetCookie(
		"session_token",
		token,
		h.sessionService.GetSessionMaxAge(), // 7 days
		"/",
		"",
		false, // Secure (set to true in production with HTTPS)
		true,  // HttpOnly
	)

	// Redirect to frontend home page
	c.Redirect(http.StatusFound, "http://localhost:3000/")
}

// generateRandomState generates a random state token for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
