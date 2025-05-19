package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"atulsm/userservice/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// MockTokenGenerator mocks the TokenGenerator interface
type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(userID string) (string, error) {
	return "test-jwt-token", nil
}

// MockPasswordHasher mocks the PasswordHasher interface
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) CheckPasswordHash(password, hash string) bool {
	return password == "password123" // Only return true for our test password
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(req *models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(id uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ListUsers(limit, offset int) ([]*models.User, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(id uuid.UUID, newPassword string) error {
	args := m.Called(id, newPassword)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockTokenGen := new(MockTokenGenerator)
	mockPwHasher := new(MockPasswordHasher)
	handler := NewUserHandler(mockRepo, mockTokenGen, mockPwHasher)

	testUser := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: sql.NullString{String: "+1234567890", Valid: true},
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"email":        "test@example.com",
				"password":     "password123",
				"first_name":   "John",
				"last_name":    "Doe",
				"phone_number": "+1234567890",
			},
			mockSetup: func() {
				mockRepo.On("CreateUser", mock.AnythingOfType("*models.RegisterRequest")).Return(testUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"token": "test-jwt-token",
				"user": map[string]interface{}{
					"id":           testUser.ID.String(),
					"email":        testUser.Email,
					"first_name":   testUser.FirstName,
					"last_name":    testUser.LastName,
					"phone_number": testUser.PhoneNumber,
					"created_at":   testUser.CreatedAt.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name: "missing required field",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid phone number",
			requestBody: map[string]interface{}{
				"email":        "test@example.com",
				"password":     "password123",
				"first_name":   "John",
				"last_name":    "Doe",
				"phone_number": "invalid-phone",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.New()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedBody, response)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockTokenGen := new(MockTokenGenerator)
	mockPwHasher := new(MockPasswordHasher)
	handler := NewUserHandler(mockRepo, mockTokenGen, mockPwHasher)

	testUser := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Password:    "hashed_password",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: sql.NullString{String: "+1234567890", Valid: true},
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", "test@example.com").Return(testUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "test-jwt-token",
				"user": map[string]interface{}{
					"id":           testUser.ID.String(),
					"email":        testUser.Email,
					"first_name":   testUser.FirstName,
					"last_name":    testUser.LastName,
					"phone_number": testUser.PhoneNumber,
					"created_at":   testUser.CreatedAt.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name: "invalid credentials",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpass",
			},
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", "test@example.com").Return(testUser, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.New()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedBody, response)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockTokenGen := new(MockTokenGenerator)
	mockPwHasher := new(MockPasswordHasher)
	handler := NewUserHandler(mockRepo, mockTokenGen, mockPwHasher)

	userID := uuid.New()
	testUser := &models.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful profile retrieval",
			setupContext: func(c *gin.Context) {
				c.Set("userID", userID.String())
			},
			mockSetup: func() {
				mockRepo.On("GetUserByID", userID).Return(testUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user ID in context",
			setupContext:   func(c *gin.Context) {},
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.New()
			router.GET("/profile", func(c *gin.Context) {
				tt.setupContext(c)
				handler.GetProfile(c)
			})

			req := httptest.NewRequest("GET", "/profile", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockTokenGen := new(MockTokenGenerator)
	mockPwHasher := new(MockPasswordHasher)
	handler := NewUserHandler(mockRepo, mockTokenGen, mockPwHasher)

	testUsers := []*models.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			FirstName: "User",
			LastName:  "One",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			FirstName: "User",
			LastName:  "Two",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful users list",
			queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
			},
			mockSetup: func() {
				mockRepo.On("ListUsers", 10, 0).Return(testUsers, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.New()
			router.GET("/users", handler.ListUsers)

			url := "/users"
			if len(tt.queryParams) > 0 {
				url += "?"
				for k, v := range tt.queryParams {
					url += k + "=" + v + "&"
				}
			}

			req := httptest.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	mockTokenGen := new(MockTokenGenerator)
	mockPwHasher := new(MockPasswordHasher)
	handler := NewUserHandler(mockRepo, mockTokenGen, mockPwHasher)

	userID := uuid.New()

	tests := []struct {
		name           string
		userID         string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "successful deletion",
			userID: userID.String(),
			mockSetup: func() {
				mockRepo.On("DeleteUser", userID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid UUID",
			userID:         "invalid-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			router := gin.New()
			router.DELETE("/users/:id", handler.DeleteUser)

			req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
