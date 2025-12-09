// Package middleware provides HTTP middleware functions for the API
package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"server/internal/auth"

	chimiddleware "github.com/go-chi/chi/middleware"
)

// CORS middleware allows cross-origin requests from any origin
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging is middleware that logs HTTP request details
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", start.Format("15:04:05"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - %v", start.Format("15:04:05"), r.Method, r.RequestURI, time.Since(start))
	})
}



// StripSlashes is chi's built-in middleware that removes trailing slashes from request paths
var StripSlashes = chimiddleware.StripSlashes

// sendUnauthorized is a helper that sends standard unauthorized response
func sendUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

// Auth middleware verifies JWT tokens without database lookup (stateless)
// Token is expected in "Authorization: Bearer <token>" header
// Extracted userID is passed via X-User-ID header to handlers
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendUnauthorized(w, "missing authorization header")
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendUnauthorized(w, "invalid authorization format")
			return
		}

		// Verify JWT token (no database lookup required)
		claims, err := auth.VerifyJWT(parts[1])
		if err != nil {
			sendUnauthorized(w, "token expired or invalid, please login again")
			return
		}

		log.Printf("Authorized request for user: %s", claims.UserID)

		// Add UserID to request header for handlers to use
		r.Header.Set("X-User-ID", claims.UserID)
		next.ServeHTTP(w, r)
	})
}
