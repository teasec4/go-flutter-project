// Package middleware provides HTTP middleware functions for the API
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"server/internal/models"
	"server/internal/storage"
)

// Logging is middleware that logs HTTP request details
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", start.Format("15:04:05"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - %v", start.Format("15:04:05"), r.Method, r.RequestURI, time.Since(start))
	})
}

// JSONContent is middleware that validates Content-Type header for POST/PUT requests
func JSONContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Header.Get("Content-Type") != "application/json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Content-Type must be application/json"}`))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// StripSlashes is chi's built-in middleware that removes trailing slashes from request paths
var StripSlashes = chimiddleware.StripSlashes

// GenerateToken creates a token and saves it to database
// Token format: "accountId:uuid"
// Default expiration: 24 hours
func GenerateToken(userID, accountID string, tokenRepo *storage.TokenRepository) (string, error) {
	tokenUUID := uuid.New().String()
	tokenValue := fmt.Sprintf("%s:%s", accountID, tokenUUID)
	
	token := &models.Token{
		ID:        uuid.New().String(),
		UserID:    userID,
		AccountID: accountID,
		Token:     tokenValue,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	err := tokenRepo.CreateToken(token)
	if err != nil {
		return "", err
	}
	
	log.Printf("Generated token for user %s (account %s)", userID, accountID)
	return tokenValue, nil
}

// Auth verifies that requests contain a valid authorization token from database
// Token is expected in "Authorization: Bearer <token>" header
func AuthWithTokenRepo(tokenRepo *storage.TokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"missing authorization header"}`))
				return
			}

			// Extract token from "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid authorization format"}`))
				return
			}

			tokenValue := parts[1]
			
			// Check if token exists in database and not expired
			token, err := tokenRepo.GetTokenByValue(tokenValue)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error":"invalid or expired token"}`))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"database error"}`))
				return
			}

			log.Printf("Authorized request for account: %s", token.AccountID)
			next.ServeHTTP(w, r)
		})
	}
}
