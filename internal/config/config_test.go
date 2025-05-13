package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Helper function to clean environment variables
	cleanup := func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENVIRONMENT")
	}

	tests := []struct {
		name        string
		envVars     map[string]string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name:    "development environment with no env vars should use defaults",
			envVars: map[string]string{},
			wantConfig: &Config{
				Port:        "8080",
				DatabaseURL: "postgres://postgres:postgres@localhost:5432/userservice?sslmode=disable",
				JWTSecret:   "dev-jwt-secret-do-not-use-in-production",
				Environment: "development",
			},
			wantErr: false,
		},
		{
			name: "production environment requires DATABASE_URL and JWT_SECRET",
			envVars: map[string]string{
				"ENVIRONMENT": "production",
			},
			wantConfig:  nil,
			wantErr:     true,
			errContains: "DATABASE_URL environment variable is required",
		},
		{
			name: "production with DATABASE_URL but no JWT_SECRET should error",
			envVars: map[string]string{
				"ENVIRONMENT":  "production",
				"DATABASE_URL": "postgres://user:pass@host:5432/db",
			},
			wantConfig:  nil,
			wantErr:     true,
			errContains: "JWT_SECRET environment variable is required",
		},
		{
			name: "custom values should override defaults",
			envVars: map[string]string{
				"PORT":         "3000",
				"DATABASE_URL": "postgres://custom:custom@host:5432/customdb",
				"JWT_SECRET":   "custom-secret",
				"ENVIRONMENT":  "staging",
			},
			wantConfig: &Config{
				Port:        "3000",
				DatabaseURL: "postgres://custom:custom@host:5432/customdb",
				JWTSecret:   "custom-secret",
				Environment: "staging",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment before each test
			cleanup()

			// Set environment variables for the test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Run the test
			gotConfig, err := Load()

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Error("Load() expected error but got nil")
					return
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("Load() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			// Check no error when not expected
			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}

			// Compare configs
			if gotConfig.Port != tt.wantConfig.Port {
				t.Errorf("Load() Port = %v, want %v", gotConfig.Port, tt.wantConfig.Port)
			}
			if gotConfig.DatabaseURL != tt.wantConfig.DatabaseURL {
				t.Errorf("Load() DatabaseURL = %v, want %v", gotConfig.DatabaseURL, tt.wantConfig.DatabaseURL)
			}
			if gotConfig.JWTSecret != tt.wantConfig.JWTSecret {
				t.Errorf("Load() JWTSecret = %v, want %v", gotConfig.JWTSecret, tt.wantConfig.JWTSecret)
			}
			if gotConfig.Environment != tt.wantConfig.Environment {
				t.Errorf("Load() Environment = %v, want %v", gotConfig.Environment, tt.wantConfig.Environment)
			}
		})
	}
}
