package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"atulsm/userservice/internal/models"
	"atulsm/userservice/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	router     *gin.Engine
	authToken  string
	adminToken string
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	// Set test environment
	os.Setenv("ENVIRONMENT", "test")

	// Ensure we're using the test database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		s.T().Fatal("DATABASE_URL environment variable is required")
	}

	// Set a consistent JWT secret for tests
	os.Setenv("JWT_SECRET", "test-jwt-secret")

	// Initialize the router
	router, err := server.SetupRouter()
	if err != nil {
		s.T().Fatalf("Failed to setup router: %v", err)
	}
	s.router = router

	// Login as regular user and admin to get tokens
	s.authToken = s.loginAndGetToken("user@example.com", "User123!")
	s.adminToken = s.loginAndGetToken("admin@example.com", "Admin123!")
}

func (s *APITestSuite) loginAndGetToken(email, password string) string {
	loginReq := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	var response models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Token
}

// Test Registration
func (s *APITestSuite) TestRegisterUser() {
	tests := []struct {
		name       string
		input      models.RegisterRequest
		wantStatus int
		wantError  bool
	}{
		{
			name: "valid registration",
			input: models.RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "Password123!",
				FirstName: "New",
				LastName:  "User",
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "existing email",
			input: models.RegisterRequest{
				Email:     "user@example.com",
				Password:  "Password123!",
				FirstName: "Existing",
				LastName:  "User",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "invalid email",
			input: models.RegisterRequest{
				Email:     "invalid-email",
				Password:  "Password123!",
				FirstName: "Invalid",
				LastName:  "Email",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
			if !tt.wantError {
				var response models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tt.input.Email, response.Email)
			}
		})
	}
}

// Test Login
func (s *APITestSuite) TestLogin() {
	tests := []struct {
		name       string
		input      models.LoginRequest
		wantStatus int
		wantError  bool
	}{
		{
			name: "valid login",
			input: models.LoginRequest{
				Email:    "user@example.com",
				Password: "User123!",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "invalid password",
			input: models.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
		{
			name: "non-existent user",
			input: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "Password123!",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
			if !tt.wantError {
				var response models.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), response.Token)
				assert.Equal(s.T(), tt.input.Email, response.User.Email)
			}
		})
	}
}

// Test Get Profile
func (s *APITestSuite) TestGetProfile() {
	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid token",
			token:      s.authToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "no token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/users/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), "user@example.com", response.Email)
			}
		})
	}
}

// Test Update Profile
func (s *APITestSuite) TestUpdateProfile() {
	tests := []struct {
		name       string
		token      string
		input      models.UpdateProfileRequest
		wantStatus int
	}{
		{
			name:  "valid update",
			token: s.authToken,
			input: models.UpdateProfileRequest{
				FirstName: "UpdatedFirst",
				LastName:  "UpdatedLast",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "update email",
			token: s.authToken,
			input: models.UpdateProfileRequest{
				Email: "updated@example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid email format",
			token: s.authToken,
			input: models.UpdateProfileRequest{
				Email: "invalid-email",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			body, _ := json.Marshal(tt.input)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/users/profile", bytes.NewBuffer(body))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			req.Header.Set("Content-Type", "application/json")
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				if tt.input.FirstName != "" {
					assert.Equal(s.T(), tt.input.FirstName, response.FirstName)
				}
				if tt.input.LastName != "" {
					assert.Equal(s.T(), tt.input.LastName, response.LastName)
				}
				if tt.input.Email != "" {
					assert.Equal(s.T(), tt.input.Email, response.Email)
				}
			}
		})
	}
}

// Test List Users
func (s *APITestSuite) TestListUsers() {
	tests := []struct {
		name       string
		token      string
		query      string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "list all users",
			token:      s.adminToken,
			query:      "",
			wantStatus: http.StatusOK,
			wantCount:  4, // Total number of seed users
		},
		{
			name:       "list with pagination",
			token:      s.adminToken,
			query:      "?limit=2&offset=0",
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "unauthorized",
			token:      "",
			query:      "",
			wantStatus: http.StatusUnauthorized,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/users"+tt.query, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response []*models.UserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tt.wantCount, len(response))
			}
		})
	}
}

// Test Delete User
func (s *APITestSuite) TestDeleteUser() {
	// First create a user to delete
	registerReq := models.RegisterRequest{
		Email:     "todelete@example.com",
		Password:  "Delete123!",
		FirstName: "To",
		LastName:  "Delete",
	}
	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	var newUser models.UserResponse
	json.Unmarshal(w.Body.Bytes(), &newUser)

	tests := []struct {
		name       string
		token      string
		userID     string
		wantStatus int
	}{
		{
			name:       "delete existing user",
			token:      s.adminToken,
			userID:     newUser.ID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "delete non-existent user",
			token:      s.adminToken,
			userID:     "non-existent-id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unauthorized deletion",
			token:      s.authToken, // Non-admin user
			userID:     newUser.ID.String(),
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/users/"+tt.userID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.wantStatus, w.Code)
		})
	}
}
