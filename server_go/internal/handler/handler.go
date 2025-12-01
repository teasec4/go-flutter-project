// Package handler defines HTTP request handlers for the bank API
package handler

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"server/internal/middleware"
	"server/internal/models"
	"server/internal/storage"
	"github.com/go-chi/chi"
)

// Routes registers all account-related API routes
func Routes(r *chi.Mux, accountRepo *storage.AccountRepository, userRepo *storage.UserRepository, tokenRepo *storage.TokenRepository) {
	// Login route (no auth required)
	r.Post("/login", login(userRepo, tokenRepo))
	r.Post("/register", register(userRepo, accountRepo))

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

// getBalance handles GET /account
// Retrieves the current balance of the authenticated user's account
func getBalance(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Get userID from header (set by auth middleware)
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			sendError(w, http.StatusUnauthorized, "missing user information")
			return
		}

		// Retrieve account from database by user ID
		accounts, err := accountRepo.GetAccountsByUserID(ctx, userID)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		if len(accounts) == 0 {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		account := accounts[0] // User has exactly one account
		// Write successful response with current balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balanceResponse{
			AccountId: account.ID,
			Balance:   account.GetBalance(),
		})
	}
}

// deposit handles POST /account/deposit
// Deposits funds into the authenticated user's account
func deposit(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse incoming JSON request body
		var req struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Get userID from header (set by auth middleware)
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			sendError(w, http.StatusUnauthorized, "missing user information")
			return
		}

		// Lookup account in database by user ID
		accounts, err := accountRepo.GetAccountsByUserID(ctx, userID)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		if len(accounts) == 0 {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		account := accounts[0] // User has exactly one account

		// Execute deposit operation (validates amount)
		if err := account.Deposit(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save updated balance to database
		if err := accountRepo.UpdateBalance(ctx, account.ID, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(depositResponse{
			AccountId: account.ID,
			Balance:   account.GetBalance(),
		})
	}
}

// withdraw handles POST /account/withdraw
// Withdraws funds from the authenticated user's account
func withdraw(accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse incoming JSON request body
		var req struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Get userID from header (set by auth middleware)
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			sendError(w, http.StatusUnauthorized, "missing user information")
			return
		}

		// Lookup account in database by user ID
		accounts, err := accountRepo.GetAccountsByUserID(ctx, userID)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		if len(accounts) == 0 {
			sendError(w, http.StatusNotFound, "account not found")
			return
		}

		account := accounts[0] // User has exactly one account

		// Execute withdrawal operation
		if err := account.Withdraw(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save updated balance to database
		if err := accountRepo.UpdateBalance(ctx, account.ID, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		// Write successful response with updated balance
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(withdrawResponse{
			AccountId: account.ID,
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

// register handles POST /register
// Creates a new user with userId and hashed password, and creates an associated account
func register(userRepo *storage.UserRepository, accountRepo *storage.AccountRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Validate userId and password
		if req.UserId == "" || req.Password == "" {
			sendError(w, http.StatusBadRequest, "userId and password are required")
			return
		}

		// Check if user already exists
		_, err := userRepo.GetUserByID(req.UserId)
		if err == nil {
			// User exists
			sendError(w, http.StatusBadRequest, "user already exists")
			return
		}
		if err != gorm.ErrRecordNotFound {
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Hash password
		hashedPassword, err := models.HashPassword(req.Password)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to process password")
			return
		}

		// Create user
		user := &models.User{
			ID:       req.UserId,
			Password: hashedPassword,
		}

		if err := userRepo.CreateUser(user); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to create user")
			return
		}

		// Create associated account
		account := &models.Account{
			UserID:  req.UserId,
			Balance: 0,
		}

		if err := accountRepo.CreateAccount(account); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to create account")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(registerResponse{
			UserId:  req.UserId,
			Message: "user registered successfully",
		})
	}
}

// login handles POST /login
// Authenticates user with userId and password, returns a token
func login(userRepo *storage.UserRepository, tokenRepo *storage.TokenRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Find user by ID
		user, err := userRepo.GetUserByID(req.UserId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusUnauthorized, "invalid userId or password")
				return
			}
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		// Verify password
		if err := models.CheckPassword(user.Password, req.Password); err != nil {
			sendError(w, http.StatusUnauthorized, "invalid userId or password")
			return
		}

		// Generate and save token to database
		token, err := middleware.GenerateToken(req.UserId, tokenRepo)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to generate token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}
