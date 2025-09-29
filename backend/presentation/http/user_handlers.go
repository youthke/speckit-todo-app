package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"todo-app/application/user"
)

// UserResponse represents the HTTP response format for a user
type UserResponse struct {
	ID          uint                `json:"id"`
	Email       string              `json:"email"`
	Profile     UserProfileResponse `json:"profile"`
	Preferences UserPreferencesResponse `json:"preferences"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// UserProfileResponse represents the HTTP response format for user profile
type UserProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Timezone  string `json:"timezone"`
}

// UserPreferencesResponse represents the HTTP response format for user preferences
type UserPreferencesResponse struct {
	DefaultTaskPriority string `json:"default_task_priority"`
	EmailNotifications  bool   `json:"email_notifications"`
	ThemePreference     string `json:"theme_preference"`
}

// RegisterUserRequest represents the HTTP request format for user registration
type RegisterUserRequest struct {
	Email     string                      `json:"email" binding:"required,email,max=255"`
	Profile   RegisterUserProfileRequest  `json:"profile" binding:"required"`
	Preferences *RegisterUserPreferencesRequest `json:"preferences,omitempty"`
}

// RegisterUserProfileRequest represents the profile part of user registration
type RegisterUserProfileRequest struct {
	FirstName string `json:"first_name" binding:"required,max=50"`
	LastName  string `json:"last_name" binding:"required,max=50"`
	Timezone  string `json:"timezone" binding:"required"`
}

// RegisterUserPreferencesRequest represents the preferences part of user registration
type RegisterUserPreferencesRequest struct {
	DefaultTaskPriority *string `json:"default_task_priority,omitempty" binding:"omitempty,oneof=low medium high"`
	EmailNotifications  *bool   `json:"email_notifications,omitempty"`
	ThemePreference     *string `json:"theme_preference,omitempty" binding:"omitempty,oneof=light dark auto"`
}

// UpdateUserProfileRequest represents the HTTP request format for updating user profile
type UpdateUserProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=50"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=50"`
	Timezone  *string `json:"timezone,omitempty"`
}

// UpdateUserPreferencesRequest represents the HTTP request format for updating user preferences
type UpdateUserPreferencesRequest struct {
	DefaultTaskPriority *string `json:"default_task_priority,omitempty" binding:"omitempty,oneof=low medium high"`
	EmailNotifications  *bool   `json:"email_notifications,omitempty"`
	ThemePreference     *string `json:"theme_preference,omitempty" binding:"omitempty,oneof=light dark auto"`
}

// UserHandlers contains HTTP handlers for user-related endpoints
type UserHandlers struct {
	userService user.UserApplicationService
}

// NewUserHandlers creates a new user handlers instance
func NewUserHandlers(userService user.UserApplicationService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

// RegisterRoutes registers all user-related routes
func (h *UserHandlers) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", h.RegisterUser)
		userRoutes.GET("/profile", h.GetUserProfile)
		userRoutes.PUT("/profile", h.UpdateUserProfile)
		userRoutes.GET("/preferences", h.GetUserPreferences)
		userRoutes.PUT("/preferences", h.UpdateUserPreferences)
	}
}

// RegisterUser handles POST /api/v1/users/register
func (h *UserHandlers) RegisterUser(c *gin.Context) {
	// Parse request body
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Create command
	cmd := user.RegisterUserCommand{
		Email:     req.Email,
		FirstName: req.Profile.FirstName,
		LastName:  req.Profile.LastName,
		Timezone:  req.Profile.Timezone,
	}

	// Add optional preferences
	if req.Preferences != nil {
		cmd.DefaultTaskPriority = req.Preferences.DefaultTaskPriority
		cmd.EmailNotifications = req.Preferences.EmailNotifications
		cmd.ThemePreference = req.Preferences.ThemePreference
	}

	// Register user using application service
	registeredUser, err := h.userService.RegisterUser(cmd)
	if err != nil {
		if isValidationError(err) || isEmailConflictError(err) {
			statusCode := http.StatusUnprocessableEntity
			if isEmailConflictError(err) {
				statusCode = http.StatusConflict
			}
			c.JSON(statusCode, ErrorResponse{
				Error:   "registration_failed",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "registration_failed",
				Message: "Failed to register user",
				Details: err.Error(),
			})
		}
		return
	}

	// Convert to response format
	response := h.convertUserToResponse(registeredUser)
	c.JSON(http.StatusCreated, response)
}

// GetUserProfile handles GET /api/v1/users/profile
func (h *UserHandlers) GetUserProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get user profile from application service
	userEntity, err := h.userService.GetUserProfile(userIDUint)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "retrieval_failed",
				Message: "Failed to retrieve user profile",
			})
		}
		return
	}

	// Convert to response format
	response := h.convertUserToResponse(userEntity)
	c.JSON(http.StatusOK, response)
}

// UpdateUserProfile handles PUT /api/v1/users/profile
func (h *UserHandlers) UpdateUserProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Create command
	cmd := user.UpdateUserProfileCommand{
		UserID:    userIDUint,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Timezone:  req.Timezone,
	}

	// Update user profile using application service
	updatedUser, err := h.userService.UpdateUserProfile(cmd)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
		} else if isValidationError(err) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update user profile",
				Details: err.Error(),
			})
		}
		return
	}

	// Convert to response format
	response := h.convertUserToResponse(updatedUser)
	c.JSON(http.StatusOK, response)
}

// GetUserPreferences handles GET /api/v1/users/preferences
func (h *UserHandlers) GetUserPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get user preferences from application service
	preferences, err := h.userService.GetUserPreferences(userIDUint)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "retrieval_failed",
				Message: "Failed to retrieve user preferences",
			})
		}
		return
	}

	// Convert to response format
	response := h.convertPreferencesToResponse(preferences)
	c.JSON(http.StatusOK, response)
}

// UpdateUserPreferences handles PUT /api/v1/users/preferences
func (h *UserHandlers) UpdateUserPreferences(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req UpdateUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Create command
	cmd := user.UpdateUserPreferencesCommand{
		UserID:              userIDUint,
		DefaultTaskPriority: req.DefaultTaskPriority,
		EmailNotifications:  req.EmailNotifications,
		ThemePreference:     req.ThemePreference,
	}

	// Update user preferences using application service
	updatedPreferences, err := h.userService.UpdateUserPreferences(cmd)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
		} else if isValidationError(err) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update user preferences",
				Details: err.Error(),
			})
		}
		return
	}

	// Convert to response format
	response := h.convertPreferencesToResponse(updatedPreferences)
	c.JSON(http.StatusOK, response)
}

// Helper functions

// convertUserToResponse converts a domain user entity to HTTP response format
func (h *UserHandlers) convertUserToResponse(userEntity interface{}) UserResponse {
	// This would use proper type assertion in a real implementation
	// For now, we'll assume the conversion logic
	return UserResponse{
		// Conversion logic would be implemented here
		// This is a placeholder
	}
}

// convertPreferencesToResponse converts user preferences to HTTP response format
func (h *UserHandlers) convertPreferencesToResponse(preferences interface{}) UserPreferencesResponse {
	// This would use proper type assertion in a real implementation
	// For now, we'll assume the conversion logic
	return UserPreferencesResponse{
		// Conversion logic would be implemented here
		// This is a placeholder
	}
}

// Error checking helper functions specific to user operations
func isEmailConflictError(err error) bool {
	// Check if error indicates email conflict (already exists)
	// This would be implemented based on your error handling strategy
	return false
}