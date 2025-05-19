package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"atulsm/userservice/internal/models"
	"atulsm/userservice/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}

type PasswordHasher interface {
	CheckPasswordHash(password, hash string) bool
	HashPassword(password string) (string, error)
}

type UserHandler struct {
	repo     repository.UserRepository
	tokenGen TokenGenerator
	pwHasher PasswordHasher
}

func NewUserHandler(repo repository.UserRepository, tokenGen TokenGenerator, pwHasher PasswordHasher) *UserHandler {
	return &UserHandler{
		repo:     repo,
		tokenGen: tokenGen,
		pwHasher: pwHasher,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := h.tokenGen.GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber.String,
			CreatedAt:   user.CreatedAt,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug log: Log the email being checked
	log.Printf("Attempting login for email: %s", req.Email)

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		// Debug log: Log if user not found
		log.Printf("User not found for email: %s, error: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Debug log: Log if password check fails
	if !h.pwHasher.CheckPasswordHash(req.Password, user.Password) {
		log.Printf("Invalid password for user: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate token
	token, err := h.tokenGen.GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber.String,
			CreatedAt:   user.CreatedAt,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse UUID
	id, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get user
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber.String,
		CreatedAt:   user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse UUID
	id, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	user, err := h.repo.UpdateUser(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get user
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber.String,
		CreatedAt:   user.CreatedAt,
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Invalid limit parameter: %s", limitStr)
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("Invalid offset parameter: %s", offsetStr)
		offset = 0
	}

	log.Printf("Fetching users with limit: %d, offset: %d", limit, offset)

	// Get users
	users, err := h.repo.ListUsers(limit, offset)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully fetched %d users", len(users))

	// Convert to response objects
	response := make([]models.UserResponse, len(users))
	for i, user := range users {
		response[i] = models.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber.String,
			CreatedAt:   user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Delete user
	err = h.repo.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Hash the new password
	passwordHash, err := h.pwHasher.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update the user's password
	if err := h.repo.UpdatePassword(user.ID, passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	user, err := h.repo.UpdateUser(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber.String,
		CreatedAt:   user.CreatedAt,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header is required"})
		return
	}

	// Extract the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header format"})
		return
	}

	// In a real implementation, you might want to:
	// 1. Add the token to a blacklist in Redis
	// 2. Clear any server-side sessions
	// 3. Update the user's last logout timestamp
	// For now, we'll just return success as the client will remove the token

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
