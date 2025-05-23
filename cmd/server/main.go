package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atulsm/user-service/internal/config"
	"github.com/atulsm/user-service/internal/grpc"
	"github.com/atulsm/user-service/internal/handlers"
	"github.com/atulsm/user-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// Dummy implementations for demonstration
// Replace with your actual implementations

type dummyTokenGen struct{}

func (d *dummyTokenGen) GenerateToken(userID string) (string, error) { return "dummy-token", nil }

type dummyPwHasher struct{}

func (d *dummyPwHasher) CheckPasswordHash(password, hash string) bool { return password == hash }
func (d *dummyPwHasher) HashPassword(password string) (string, error) { return password, nil }

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	userRepo := repository.NewUserRepository(cfg.DatabaseURL)

	grpcServer := grpc.NewServer(userRepo)
	go func() {
		if err := grpcServer.Start(50051); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	router := gin.Default()

	// router.Use(middleware.CORS())
	// router.Use(middleware.RequestID())
	// router.Use(middleware.Logger())

	tokenGen := &dummyTokenGen{}
	pwHasher := &dummyPwHasher{}
	userHandler := handlers.NewUserHandler(userRepo, tokenGen, pwHasher)

	router.POST("/api/users/register", userHandler.Register)
	router.POST("/api/users/login", userHandler.Login)

	protected := router.Group("/api")
	// protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.GET("/users", userHandler.ListUsers)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	grpcServer.Stop()

	log.Println("Servers shutdown complete")
}
