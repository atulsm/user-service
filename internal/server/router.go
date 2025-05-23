package server

import (
	"github.com/atulsm/user-service/internal/config"
	"github.com/atulsm/user-service/internal/handlers"
	"github.com/atulsm/user-service/internal/middleware"
	"github.com/atulsm/user-service/internal/repository"
	"github.com/atulsm/user-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns the Gin router with all routes configured
func SetupRouter() (*gin.Engine, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize repository
	repo := repository.NewUserRepository(cfg.DatabaseURL)

	// Initialize router
	router := gin.Default()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers with all required dependencies
	userHandler := handlers.NewUserHandler(
		repo,
		middleware.NewTokenGenerator(cfg.JWTSecret),
		utils.NewPasswordHasher(),
	)

	// Public routes
	router.POST("/api/v1/auth/register", userHandler.Register)
	router.POST("/api/v1/auth/login", userHandler.Login)
	router.POST("/api/v1/auth/reset-password", userHandler.ResetPassword)
	router.GET("/users/:id", userHandler.GetUser)
	router.GET("/users", userHandler.ListUsers)

	// Protected routes
	authorized := router.Group("/api/v1")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/users/profile", userHandler.GetProfile)
		authorized.PUT("/users/profile", userHandler.UpdateProfile)
		authorized.GET("/users", userHandler.ListUsers)
		authorized.POST("/users", userHandler.CreateUser)
		authorized.PUT("/users/:id", userHandler.UpdateUser)
		authorized.DELETE("/users/:id", userHandler.DeleteUser)
		authorized.POST("/auth/logout", userHandler.Logout)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router, nil
}
