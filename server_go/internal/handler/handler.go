// Package handler defines HTTP request handlers for the bank API
package handler

import (
	"encoding/json"
	"net/http"
	"server/internal/bank"
	"server/internal/middleware"
	"github.com/go-chi/chi"
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
		router.Use(middleware.Auth)
		// Apply logging middleware to all /account routes
		router.Use(middleware.Logging)
		router.Get("/", getBalance(b))
		router.Post("/deposit", deposit(b))
		router.Post("/withdraw", withdraw(b))
	})
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
		token := middleware.GenerateToken(req.AccountId)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}

