package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atulsm/userservice/internal/config"
	"atulsm/userservice/internal/handlers"
	"atulsm/userservice/internal/middleware"
	"atulsm/userservice/internal/repository"
	"atulsm/userservice/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository
	repo := repository.NewUserRepository(cfg.DatabaseURL)
	defer repo.Close()

	// Initialize router
	router := gin.Default()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Initialize handlers with all required dependencies
	userHandler := handlers.NewUserHandler(
		repo,
		middleware.NewTokenGenerator(cfg.JWTSecret),
		utils.NewPasswordHasher(),
	)

	// Public routes
	router.POST("/api/users/register", userHandler.Register)
	router.POST("/api/users/login", userHandler.Login)

	// Protected routes
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/users/profile", userHandler.GetProfile)
		authorized.PUT("/users/profile", userHandler.UpdateProfile)
		authorized.GET("/users", userHandler.ListUsers)
		authorized.GET("/users/:id", userHandler.GetUser)
		authorized.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
