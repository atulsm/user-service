package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header format is correct
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid Authorization header format: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Debug log: Log the token being validated
		log.Printf("Validating token: %s", parts[1])

		// Validate token
		userID, err := ValidateToken(parts[1])
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Debug log: Log successful token validation
		log.Printf("Token validated successfully for user: %s", userID)

		// Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
}

// TokenGenerator generates JWT tokens
type TokenGenerator struct {
	secret string
}

// NewTokenGenerator creates a new TokenGenerator with the given secret
func NewTokenGenerator(secret string) *TokenGenerator {
	return &TokenGenerator{secret: secret}
}

// GenerateToken generates a new JWT token for the given user ID
func (t *TokenGenerator) GenerateToken(userID string) (string, error) {
	if t.secret == "" {
		return "", errors.New("JWT secret is not set")
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
		"iat": time.Now().Unix(),
	})

	// Sign token
	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf("JWT_SECRET environment variable is not set")
		return "", errors.New("JWT_SECRET environment variable is not set")
	}
	log.Printf("Using JWT secret: %s", secret)

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("Token parsing failed: %v", err)
		return "", err
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				log.Printf("Token expired at %v, current time: %v", exp, time.Now().Unix())
				return "", errors.New("token expired")
			}
		} else {
			log.Printf("Invalid expiration claim in token")
			return "", errors.New("invalid token claims")
		}

		// Get user ID
		if sub, ok := claims["sub"].(string); ok {
			log.Printf("Token validated successfully for user: %s", sub)
			return sub, nil
		}
		log.Printf("Invalid user ID claim in token")
		return "", errors.New("invalid user ID in token")
	}
	log.Printf("Token validation failed: invalid token")
	return "", errors.New("invalid token")
}
