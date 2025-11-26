// Package handler defines HTTP request handlers for the bank API
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"server/internal/bank"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// tokenStore holds valid tokens in memory (token -> accountId mapping)
var (
	tokenStore = &sync.Map{}
)

// Routes registers all account-related API routes
// Routes: GET /account, POST /account/deposit, POST /account/withdraw
// Applies logging middleware to track all account operations
func Routes(r *chi.Mux, b *bank.Bank) {
	// Import logging middleware
	// The middleware is applied at the route group level to log all account operations
	
	// Login route (no auth required)
	r.Post("/login", login(b))
	
	r.Route("/account", func(router chi.Router) {
		// Apply auth middleware to all /account routes
		router.Use(authMiddleware)
		// Apply logging middleware to all /account routes
		router.Use(loggingMiddleware)
		
		router.Get("/", getBalance(b))
		router.Post("/deposit", deposit(b))
		router.Post("/withdraw", withdraw(b))
	})
}

// ============= Request Types =============

// depositRequest represents the incoming JSON payload for deposit operations
// Fields:
//   - AccountId: the account identifier to deposit funds into
//   - Amount: the amount of money to deposit (must be positive)
type depositRequest struct {
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}

// withdrawRequest represents the incoming JSON payload for withdrawal operations
// Fields:
//   - AccountId: the account identifier to withdraw funds from
//   - Amount: the amount of money to withdraw (must be positive and not exceed balance)
type withdrawRequest struct {
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}

// ============= Response Types =============

// balanceResponse represents the JSON response when checking account balance
type balanceResponse struct {
	AccountId string `json:"accountId"`
	Balance   int    `json:"balance"`
}

// depositResponse represents the JSON response after a successful deposit
// Returns the updated balance of the account
type depositResponse struct {
	AccountId string `json:"accountId"`
	Balance   int    `json:"balance"`
}

// withdrawResponse represents the JSON response after a successful withdrawal
// Returns the updated balance of the account
type withdrawResponse struct {
	AccountId string `json:"accountId"`
	Balance   int    `json:"balance"`
}

// errorResponse represents a generic error response sent when operations fail
type errorResponse struct {
	Error string `json:"error"`
}

// loginRequest represents the incoming JSON payload for login
type loginRequest struct {
	AccountId string `json:"accountId"`
}

// loginResponse represents the JSON response after successful login
type loginResponse struct {
	Token string `json:"token"`
}

// ============= Handlers =============

// getBalance handles GET /account?id={accountId}
// Retrieves the current balance of a specified account
//
// Query Parameters:
//   - id: required, the account identifier
//
// Responses:
//   - 200 OK: with balanceResponse containing account balance
//   - 400 Bad Request: if id parameter is missing
//   - 404 Not Found: if account does not exist
func getBalance(b *bank.Bank) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract account ID from query parameters
		accountID := r.URL.Query().Get("id")
		if accountID == "" {
			sendError(w, http.StatusBadRequest, "id parameter required")
			return
		}

		// Retrieve account from bank using account ID
		a, ok := b.GetAccount(accountID)
		if !ok {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		// Write successful response with current balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balanceResponse{
			AccountId: accountID,
			Balance:   a.GetBalance(),
		})
	}
}

// deposit handles POST /account/deposit
// Deposits funds into an account
//
// Request body (JSON):
//   - accountId: string, the target account
//   - amount: integer, positive amount to deposit
//
// Responses:
//   - 200 OK: with updated account balance
//   - 400 Bad Request: invalid request body or invalid amount (must be positive)
//   - 404 Not Found: account does not exist
func deposit(b *bank.Bank) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body into depositRequest struct
		var req depositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Lookup account in bank to verify it exists
		a, ok := b.GetAccount(req.AccountId)
		if !ok {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		// Execute deposit operation on the account (validates amount)
		if err := a.Deposit(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(depositResponse{
			AccountId: req.AccountId,
			Balance:   a.GetBalance(),
		})
	}
}

// withdraw handles POST /account/withdraw
// Withdraws funds from an account
//
// Request body (JSON):
//   - accountId: string, the source account
//   - amount: integer, positive amount to withdraw
//
// Responses:
//   - 200 OK: with updated account balance
//   - 400 Bad Request: invalid request, invalid amount, or insufficient balance
//   - 404 Not Found: account does not exist
func withdraw(b *bank.Bank) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body into withdrawRequest struct
		var req withdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Lookup account in bank to verify it exists
		a, ok := b.GetAccount(req.AccountId)
		if !ok {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		// Execute withdrawal operation on the account
		// Validates: amount must be positive and not exceed balance
		if err := a.Withdraw(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(withdrawResponse{
			AccountId: req.AccountId,
			Balance:   a.GetBalance(),
		})
	}
}

// sendError is a helper function to send error responses in JSON format
// Parameters:
//   - w: ResponseWriter to write the response to
//   - statusCode: HTTP status code (e.g., 400, 404, 500)
//   - message: error message to include in response
func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

// ============= Middleware =============

// loggingMiddleware logs all incoming HTTP requests and their processing time
// Format: [HH:MM:SS] METHOD /path - duration
// Example: [14:32:15] POST /account/deposit - 125.4ms
// This middleware is applied to all /account routes to track API usage
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()
		
		// Log request start with timestamp, HTTP method, and URI
		log.Printf("[%s] %s %s", start.Format("15:04:05"), r.Method, r.RequestURI)
		
		// Call the next handler/middleware
		next.ServeHTTP(w, r)
		
		// Log request completion with duration
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %v", start.Format("15:04:05"), r.Method, r.RequestURI, duration)
	})
}

// login handles POST /login
// Authenticates user and returns a token
//
// Request body (JSON):
//   - accountId: string, the account ID to login
//
// Responses:
//   - 200 OK: with auth token
//   - 400 Bad Request: invalid request body
//   - 404 Not Found: account does not exist
func login(b *bank.Bank) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Check if account exists
		_, ok := b.GetAccount(req.AccountId)
		if !ok {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		// Generate and return token
		token := generateToken(req.AccountId)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}

// generateToken creates a token from accountId and stores it in memory
// Format: "accountId:uuid"
// Token is stored in tokenStore for validation
func generateToken(accountId string) string {
	tokenUUID := uuid.New().String()
	token := fmt.Sprintf("%s:%s", accountId, tokenUUID)
	
	// Store token in memory (token -> accountId)
	tokenStore.Store(token, accountId)
	
	log.Printf("Generated token for account %s: %s", accountId, token)
	return token
}

// authMiddleware verifies that requests contain a valid authorization token
// Token is expected in "Authorization: Bearer <token>" header
// Token must be previously issued by /login endpoint
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := parts[1]
		
		// Check if token exists in tokenStore
		accountId, found := tokenStore.Load(token)
		if !found {
			sendError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		log.Printf("Authorized request for account: %s", accountId)
		next.ServeHTTP(w, r)
	})
}