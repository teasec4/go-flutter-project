// Package auth provides JWT token generation and verification
package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims stored in a JWT token
type JWTClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// init initializes the JWT secret from environment variable or uses default
func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "max03"
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT creates a new JWT token for a user with 24 hour expiration
// Token is self-contained and does not require database storage
func GenerateJWT(userID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyJWT validates a JWT token and returns the claims if valid
// Does not require database lookup (stateless authentication)
func VerifyJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
