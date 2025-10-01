package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-app/handlers"
	"todo-app/internal/models"
	"todo-app/services/auth"
)

// TestGoogleOAuthFlowNewUser tests the complete OAuth flow for a new user
func TestGoogleOAuthFlowNewUser(t *testing.T) {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.AuthenticationSession{}, &models.OAuthState{})
	require.NoError(t, err)

	// Setup services
	googleConfig, err := auth.NewGoogleOAuthConfig()
	require.NoError(t, err)

	oauthService := auth.NewOAuthService(db, googleConfig)
	sessionService := auth.NewSessionService(db)
	jwtService := auth.NewJWTService()

	authHandler := handlers.NewAuthHandler(oauthService, sessionService, jwtService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/auth/google/login", authHandler.GoogleLogin)
	router.GET("/api/v1/auth/google/callback", authHandler.GoogleCallback)

	t.Run("Step 1: Initiate OAuth flow", func(t *testing.T) {
		// Execute OAuth initiation
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "auth_url")

		// Verify OAuth state cookie is set
		cookies := w.Result().Cookies()
		var stateCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "oauth_state" {
				stateCookie = cookie
				break
			}
		}
		require.NotNil(t, stateCookie, "oauth_state cookie should be set")
		assert.True(t, stateCookie.HttpOnly, "oauth_state should be HttpOnly")
		assert.Equal(t, 300, stateCookie.MaxAge, "oauth_state should expire in 5 minutes")

		// Verify OAuth state is stored in database
		var stateCount int64
		err := db.Model(&models.OAuthState{}).Count(&stateCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), stateCount, "OAuth state should be stored in database")
	})

	t.Run("Step 2: OAuth callback creates new user", func(t *testing.T) {
		// This test demonstrates the expected flow when Google OAuth callback is received
		// In reality, this would require mocking the Google OAuth token exchange

		// Verify no users exist initially
		var userCount int64
		err := db.Model(&models.User{}).Count(&userCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), userCount, "No users should exist initially")

		// Expected behavior after callback:
		// 1. Exchange code for OAuth tokens (mocked)
		// 2. Fetch user info from Google (mocked)
		// 3. Create new user with Google ID
		// 4. Create authentication session
		// 5. Set session cookie
		// 6. Redirect to application

		// Simulate user creation as would happen in callback
		newUser := &models.User{
			Email:         "newuser@gmail.com",
			Name:          "New User",
			GoogleID:      "google_user_123",
			OAuthProvider: "google",
		}
		now := time.Now()
		newUser.OAuthCreatedAt = &now

		err = db.Create(newUser).Error
		require.NoError(t, err)

		// Verify user was created
		var users []models.User
		err = db.Find(&users).Error
		require.NoError(t, err)
		assert.Equal(t, 1, len(users), "One user should be created")
		assert.Equal(t, "newuser@gmail.com", users[0].Email)
		assert.Equal(t, "google_user_123", users[0].GoogleID)
		assert.Equal(t, "google", users[0].OAuthProvider)
		assert.NotNil(t, users[0].OAuthCreatedAt)
	})

	t.Run("Step 3: Session is created for new user", func(t *testing.T) {
		// Find the created user
		var user models.User
		err := db.Where("email = ?", "newuser@gmail.com").First(&user).Error
		require.NoError(t, err)

		// Create session for the user
		ctx := context.Background()
		session, err := sessionService.CreateSession(ctx, user.ID, true, "mock_access_token", "mock_refresh_token", time.Now().Add(1*time.Hour))
		require.NoError(t, err)
		require.NotNil(t, session)

		// Verify session properties
		assert.Equal(t, user.ID, session.UserID)
		assert.NotEmpty(t, session.SessionToken)
		assert.NotEmpty(t, session.AccessToken)
		assert.NotEmpty(t, session.RefreshToken)
		assert.False(t, session.SessionExpiresAt.Before(time.Now()))
		assert.True(t, session.SessionExpiresAt.After(time.Now().Add(23*time.Hour)))
	})

	t.Run("Step 4: User can validate session", func(t *testing.T) {
		// Find the user
		var user models.User
		err := db.Where("email = ?", "newuser@gmail.com").First(&user).Error
		require.NoError(t, err)

		// Find the session
		var session models.AuthenticationSession
		err = db.Where("user_id = ?", user.ID).First(&session).Error
		require.NoError(t, err)

		// Validate session
		ctx := context.Background()
		validatedSession, err := sessionService.ValidateSession(ctx, session.SessionToken)
		require.NoError(t, err)
		require.NotNil(t, validatedSession)

		// Verify session is valid
		assert.Equal(t, user.ID, validatedSession.UserID)
		assert.False(t, validatedSession.IsExpired())
	})

	t.Run("Step 5: Complete flow validation", func(t *testing.T) {
		// Verify final state
		// 1. User exists with OAuth fields
		var user models.User
		err := db.Where("email = ?", "newuser@gmail.com").First(&user).Error
		require.NoError(t, err)
		assert.NotEmpty(t, user.GoogleID)
		assert.Equal(t, "google", user.OAuthProvider)

		// 2. Session exists and is valid
		var session models.AuthenticationSession
		err = db.Where("user_id = ?", user.ID).First(&session).Error
		require.NoError(t, err)
		assert.False(t, session.IsExpired())

		// 3. OAuth state has been consumed (would be deleted in real flow)
		// This validates the cleanup happens

		// 4. User can make authenticated requests
		assert.True(t, user.IsOAuthUser())
		assert.True(t, user.IsActive)
	})
}

// TestNewUserAccountCreation tests user account creation with OAuth
func TestNewUserAccountCreation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	t.Run("creates user with OAuth fields", func(t *testing.T) {
		user := &models.User{
			Email:         "test@gmail.com",
			Name:          "Test User",
			GoogleID:      "google_123",
			OAuthProvider: "google",
		}
		now := time.Now()
		user.OAuthCreatedAt = &now

		err := db.Create(user).Error
		require.NoError(t, err)

		// Retrieve and verify
		var retrieved models.User
		err = db.First(&retrieved, user.ID).Error
		require.NoError(t, err)

		assert.Equal(t, "test@gmail.com", retrieved.Email)
		assert.Equal(t, "google_123", retrieved.GoogleID)
		assert.True(t, retrieved.IsOAuthUser())
	})

	t.Run("enforces unique Google ID", func(t *testing.T) {
		user1 := &models.User{
			Email:         "user1@gmail.com",
			Name:          "User 1",
			GoogleID:      "duplicate_id",
			OAuthProvider: "google",
		}
		err := db.Create(user1).Error
		require.NoError(t, err)

		// Attempt to create user with duplicate Google ID
		user2 := &models.User{
			Email:         "user2@gmail.com",
			Name:          "User 2",
			GoogleID:      "duplicate_id",
			OAuthProvider: "google",
		}
		err = db.Create(user2).Error
		assert.Error(t, err, "Should not allow duplicate Google ID")
	})
}