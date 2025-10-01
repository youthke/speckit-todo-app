package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"todo-app/internal/handlers"
	"todo-app/internal/models"
	"todo-app/internal/services"
	"todo-app/internal/storage"
	"todo-app/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Initialize database
	if err := storage.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func() {
		if err := storage.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Add middleware
	router.Use(handlers.ErrorHandler())
	router.Use(handlers.RequestLogger())
	router.Use(handlers.SecurityHeaders())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Allow requests from the frontend development server
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler()
	healthService := services.NewHealthService()
	googleOAuthHandler := handlers.NewGoogleOAuthHandler(storage.DB)

	// Initialize rate limiter for signup/OAuth endpoints
	// 10 requests per 15 minutes = 10 / (15 * 60) = 0.0111 requests per second
	signupRateLimiter := middleware.NewIPRateLimiter(rate.Every(15*time.Minute)/10, 10)

	// Setup routes
	setupRoutes(router, taskHandler, healthService, googleOAuthHandler, signupRateLimiter)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler, healthService *services.HealthService, googleOAuthHandler *handlers.GoogleOAuthHandler, signupRateLimiter *middleware.IPRateLimiter) {
	// Health check handler function
	healthHandler := func(c *gin.Context) {
		healthResponse, err := healthService.GetHealthStatus()
		if err != nil {
			log.Printf("Health check failed: %v", err)
			errorResponse := models.NewErrorResponse("internal_error", "Health check failed unexpectedly")
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}

		// Determine HTTP status code based on health status
		var statusCode int
		switch healthResponse.Status {
		case models.HealthStatusHealthy:
			statusCode = http.StatusOK
		case models.HealthStatusDegraded:
			statusCode = http.StatusServiceUnavailable
		case models.HealthStatusUnhealthy:
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, healthResponse)
	}

	// API group
	api := router.Group("/api")
	{
		// Health endpoint in API group
		api.GET("/health", healthHandler)

		// API v1 routes
		v1 := api.Group("/v1")
		{
			// Google OAuth routes
			auth := v1.Group("/auth")
			{
				// Apply rate limiter to signup/login endpoint
				auth.GET("/google/login", signupRateLimiter.RateLimitMiddleware(), googleOAuthHandler.GoogleLogin)
				auth.GET("/google/callback", googleOAuthHandler.GoogleCallback)
			}

			// Task routes
			tasks := v1.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetTasks)
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}
		}
	}

	// Enhanced health check endpoint (also available at root level)
	router.GET("/health", healthHandler)
}