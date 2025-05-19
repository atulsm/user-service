package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	Email       string         `json:"email" db:"email"`
	Password    string         `json:"-" db:"password_hash"` // Never returned in responses
	FirstName   string         `json:"first_name" db:"first_name"`
	LastName    string         `json:"last_name" db:"last_name"`
	PhoneNumber sql.NullString `json:"phone_number,omitempty" db:"phone_number"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,e164"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,e164"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}
