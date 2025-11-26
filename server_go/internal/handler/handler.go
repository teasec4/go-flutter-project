// Package handler defines HTTP request handlers for the bank API
package handler

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"server/internal/middleware"
	"server/internal/storage"
	"github.com/go-chi/chi"
)

// Routes registers all account-related API routes
func Routes(r *chi.Mux, accountRepo *storage.AccountRepository, userRepo *storage.UserRepository, tokenRepo *storage.TokenRepository) {
	// Login route (no auth required)
	r.Post("/login", login(accountRepo, tokenRepo))

	r.Route("/account", func(router chi.Router) {
		// Apply auth middleware to all /account routes
		router.Use(middleware.AuthWithTokenRepo(tokenRepo))
		// Apply logging middleware to all /account routes
		router.Use(middleware.Logging)
		router.Get("/", getBalance(accountRepo))
		router.Post("/deposit", deposit(accountRepo))
		router.Post("/withdraw", withdraw(accountRepo))
	})
}

// ============= Handlers =============

// getBalance handles GET /account?id={accountId}
// Retrieves the current balance of a specified account
func getBalance(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract account ID from query parameters
		accountID := r.URL.Query().Get("id")
		if accountID == "" {
			sendError(w, http.StatusBadRequest, "id parameter required")
			return
		}

		// Retrieve account from database
		account, err := accountRepo.GetAccountByID(accountID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
				return
			}
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Write successful response with current balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balanceResponse{
			AccountId: accountID,
			Balance:   (*account).GetBalance(),
		})
	}
}

// deposit handles POST /account/deposit
// Deposits funds into an account
func deposit(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req depositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Lookup account in database
		account, err := accountRepo.GetAccountByID(req.AccountId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
				return
			}
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Execute deposit operation (validates amount)
		if err := account.Deposit(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save updated balance to database
		if err := accountRepo.UpdateBalance(req.AccountId, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(depositResponse{
			AccountId: req.AccountId,
			Balance:   account.GetBalance(),
		})
	}
}

// withdraw handles POST /account/withdraw
// Withdraws funds from an account
func withdraw(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req withdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Lookup account in database
		account, err := accountRepo.GetAccountByID(req.AccountId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
				return
			}
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Execute withdrawal operation
		if err := account.Withdraw(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save updated balance to database
		if err := accountRepo.UpdateBalance(req.AccountId, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(withdrawResponse{
			AccountId: req.AccountId,
			Balance:   account.GetBalance(),
		})
	}
}

// sendError is a helper function to send error responses in JSON format
func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

// login handles POST /login
// Authenticates user and returns a token
func login(accountRepo *storage.AccountRepository, tokenRepo *storage.TokenRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Check if account exists
		_, err := accountRepo.GetAccountByID(req.AccountId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
				return
			}
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Generate and save token to database
		token, err := middleware.GenerateToken(req.UserId, req.AccountId, tokenRepo)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to generate token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}
