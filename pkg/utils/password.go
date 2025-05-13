package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles password hashing operations
type PasswordHasher struct{}

// NewPasswordHasher creates a new PasswordHasher
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

// HashPassword creates a bcrypt hash of the password
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent
func (p *PasswordHasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// For backward compatibility, keep the package-level functions
func HashPassword(password string) (string, error) {
	return NewPasswordHasher().HashPassword(password)
}

func CheckPasswordHash(password, hash string) bool {
	return NewPasswordHasher().CheckPasswordHash(password, hash)
}
