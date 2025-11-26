// Package middleware provides HTTP middleware functions for the API
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// tokenStore holds valid tokens in memory (token -> accountId mapping)
// Uses sync.Map for thread-safe concurrent access
var tokenStore = &sync.Map{}

// Logging is middleware that logs HTTP request details
// Logs the timestamp, HTTP method, URI, and request duration
// Useful for monitoring and debugging API usage
//
// Information logged:
//   - [HH:MM:SS] METHOD /path
//   - [HH:MM:SS] METHOD /path - duration (after request completes)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()
		// Log request start
		log.Printf("[%s] %s %s", start.Format("15:04:05"), r.Method, r.RequestURI)
		// Process the request
		next.ServeHTTP(w, r)
		// Log request completion with duration
		log.Printf("[%s] %s %s - %v", start.Format("15:04:05"), r.Method, r.RequestURI, time.Since(start))
	})
}

// JSONContent is middleware that validates Content-Type header for POST/PUT requests
// Ensures that all POST and PUT requests have "application/json" Content-Type
// Returns 400 Bad Request if Content-Type is missing or incorrect for these methods
//
// This middleware:
//   - Allows GET/DELETE requests without Content-Type validation
//   - Enforces application/json for POST/PUT requests
//   - Blocks requests with incorrect Content-Type on POST/PUT
func JSONContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a POST or PUT request (write operations)
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Verify Content-Type is application/json
			if r.Header.Get("Content-Type") != "application/json" {
				// Return error if Content-Type is incorrect
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Content-Type must be application/json"}`))
				return
			}
		}
		// Continue to next middleware/handler
		next.ServeHTTP(w, r)
	})
}

// StripSlashes is chi's built-in middleware that removes trailing slashes from request paths
// Example: /account/ -> /account
// This ensures consistent URL handling and prevents route mismatches
var StripSlashes = chimiddleware.StripSlashes

// GenerateToken creates a token from accountId and stores it in memory
// Format: "accountId:uuid"
// Token is stored in tokenStore for validation on subsequent requests
func GenerateToken(accountId string) string {
	tokenUUID := uuid.New().String()
	token := fmt.Sprintf("%s:%s", accountId, tokenUUID)
	
	// Store token in memory (token -> accountId)
	tokenStore.Store(token, accountId)
	
	log.Printf("Generated token for account %s: %s", accountId, token)
	return token
}

// Auth verifies that requests contain a valid authorization token
// Token is expected in "Authorization: Bearer <token>" header
// Token must be previously issued by /login endpoint
// Returns 401 Unauthorized if token is missing or invalid
func Auth(next http.Handler) http.Handler {
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

		token := parts[1]
		
		// Check if token exists in tokenStore
		accountId, found := tokenStore.Load(token)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid or expired token"}`))
			return
		}

		log.Printf("Authorized request for account: %s", accountId)
		next.ServeHTTP(w, r)
	})
}

