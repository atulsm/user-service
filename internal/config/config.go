package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		if env == "development" {
			dbURL = "postgres://postgres:postgres@localhost:5432/userservice?sslmode=disable"
		} else {
			return nil, errors.New("DATABASE_URL environment variable is required")
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if env == "development" {
			jwtSecret = "dev-jwt-secret-do-not-use-in-production"
		} else {
			return nil, errors.New("JWT_SECRET environment variable is required")
		}
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		Environment: env,
	}, nil
}
