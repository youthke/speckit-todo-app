package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todo-app/internal/models"
	"todo-app/services/auth"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	googleConfig   *auth.GoogleOAuthConfig
	oauthService   *auth.OAuthService
	sessionService *auth.SessionService
	jwtService     *auth.JWTService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(
	googleConfig *auth.GoogleOAuthConfig,
	oauthService *auth.OAuthService,
	sessionService *auth.SessionService,
	jwtService *auth.JWTService,
) *AuthHandler {
	return &AuthHandler{
		googleConfig:   googleConfig,
		oauthService:   oauthService,
		sessionService: sessionService,
		jwtService:     jwtService,
	}
}

// GoogleLogin initiates the Google OAuth flow
// GET /auth/google/login
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Get optional redirect URI from query parameter
	redirectURI := c.DefaultQuery("redirect_uri", "http://localhost:3000/dashboard")

	// Initiate OAuth flow
	result, err := h.oauthService.InitiateOAuthFlow(c.Request.Context(), redirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "oauth_init_failed",
			"message": "Failed to initiate OAuth flow",
			"details": err.Error(),
		})
		return
	}

	// Set state token as secure cookie
	c.SetCookie(
		"oauth_state",
		result.StateToken,
		300, // 5 minutes
		"/",
		"",
		false, // Secure (should be true in production with HTTPS)
		true,  // HttpOnly
	)

	// Return authorization URL
	c.JSON(http.StatusOK, gin.H{
		"auth_url":    result.AuthURL,
		"state_token": result.StateToken,
	})
}

// GoogleCallback handles the OAuth callback from Google
// GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Get authorization code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Check for OAuth errors
	if errorParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "oauth_error",
			"message": "OAuth authorization failed",
			"details": errorParam,
		})
		return
	}

	// Validate required parameters
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Missing required parameters",
		})
		return
	}

	// Verify state token from cookie
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie != state {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_state",
			"message": "Invalid or missing OAuth state",
		})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Process OAuth callback
	result, err := h.oauthService.ProcessOAuthCallback(c.Request.Context(), code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "oauth_callback_failed",
			"message": "Failed to process OAuth callback",
			"details": err.Error(),
		})
		return
	}

	// Generate JWT token
	jwtToken, err := h.jwtService.GenerateToken(
		result.User.ID,
		result.User.Email,
		result.Session.ID,
		true, // isOAuth
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_failed",
			"message": "Failed to generate session token",
		})
		return
	}

	// Update session with JWT token
	result.Session.SessionToken = jwtToken

	// Set session cookie
	c.SetCookie(
		"session_token",
		jwtToken,
		86400, // 24 hours
		"/",
		"",
		false, // Secure (should be true in production)
		true,  // HttpOnly
	)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"user":         result.User.ToResponse(),
		"session":      result.Session.ToResponse(),
		"is_new_user":  result.IsNewUser,
		"redirect_uri": result.RedirectURI,
	})
}

// ValidateSession validates the current session
// GET /auth/session/validate
func (h *AuthHandler) ValidateSession(c *gin.Context) {
	// Get session token from cookie or Authorization header
	tokenString, err := c.Cookie("session_token")
	if err != nil {
		// Try Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "no_token",
				"message": "No session token provided",
			})
			return
		}

		// Extract token from "Bearer <token>"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token_format",
				"message": "Invalid authorization header format",
			})
			return
		}
	}

	// Validate session
	result, err := h.sessionService.ValidateSession(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "validation_failed",
			"message": "Failed to validate session",
			"details": err.Error(),
		})
		return
	}

	if !result.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_session",
			"message": result.Error,
		})
		return
	}

	// Return session and user information
	var userResponse interface{}
	if user, ok := result.User.(*models.User); ok {
		userResponse = user.ToResponse()
	} else {
		userResponse = result.User
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":         true,
		"user":          userResponse,
		"session":       result.Session.ToResponse(),
		"needs_refresh": result.NeedsRefresh,
	})
}

// RefreshSession refreshes the OAuth tokens and extends session
// POST /auth/session/refresh
func (h *AuthHandler) RefreshSession(c *gin.Context) {
	// Get session token from cookie or Authorization header
	tokenString, err := c.Cookie("session_token")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "no_token",
				"message": "No session token provided",
			})
			return
		}
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
	}

	// Extract session ID from token
	sessionID, err := h.jwtService.ExtractSessionID(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_token",
			"message": "Failed to extract session ID from token",
		})
		return
	}

	// Refresh OAuth tokens if needed
	session, err := h.oauthService.RefreshOAuthToken(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "refresh_failed",
			"message": "Failed to refresh tokens",
			"details": err.Error(),
		})
		return
	}

	// Refresh session (extend expiration)
	refreshedSession, newJWT, err := h.sessionService.RefreshSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "session_refresh_failed",
			"message": "Failed to refresh session",
			"details": err.Error(),
		})
		return
	}

	// Update session cookie
	c.SetCookie(
		"session_token",
		newJWT,
		86400, // 24 hours
		"/",
		"",
		false, // Secure
		true,  // HttpOnly
	)

	// Return refreshed session
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"session":    refreshedSession.ToResponse(),
		"expires_at": session.TokenExpiresAt,
	})
}

// Logout terminates the current session
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get session token
	tokenString, err := c.Cookie("session_token")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
	}

	if tokenString != "" {
		// Extract session ID
		sessionID, err := h.jwtService.ExtractSessionID(tokenString)
		if err == nil {
			// Terminate session
			h.sessionService.TerminateSession(sessionID)
		}
	}

	// Clear session cookie
	c.SetCookie(
		"session_token",
		"",
		-1, // Expire immediately
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// RevokeWebhook handles OAuth revocation webhook from Google
// POST /auth/revoke-webhook
func (h *AuthHandler) RevokeWebhook(c *gin.Context) {
	// Parse form data
	token := c.PostForm("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_token",
			"message": "Token parameter is required",
		})
		return
	}

	// Handle revocation
	err := h.oauthService.HandleRevocationWebhook(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "revocation_failed",
			"message": "Failed to process revocation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "revoked",
	})
}

// RegisterRoutes registers all authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		// OAuth routes
		auth.GET("/google/login", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)

		// Session management routes
		auth.GET("/session/validate", h.ValidateSession)
		auth.POST("/session/refresh", h.RefreshSession)
		auth.POST("/logout", h.Logout)

		// Webhook routes
		auth.POST("/revoke-webhook", h.RevokeWebhook)
	}
}