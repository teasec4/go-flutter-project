// Package handler defines HTTP request handlers for the bank API
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/middleware"
	"server/internal/models"
	"server/internal/store"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

// Routes registers all account-related API routes
func Routes(r *chi.Mux, db *store.DB) {
	// Apply CORS middleware globally
	r.Use(middleware.CORS)
	
	// Login route (no auth required)
	r.Post("/login", login(db))
	r.Post("/register", register(db))

	r.Route("/account", func(router chi.Router) {
		// Apply auth middleware to all /account routes
		router.Use(middleware.Auth)
		// Apply logging middleware to all /account routes
		router.Use(middleware.Logging)
		router.Get("/", getBalance(db))
		router.Post("/deposit", deposit(db))
		router.Post("/withdraw", withdraw(db))
	})
}

// ============= Handlers =============

// getBalance handles GET /account
func getBalance(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		account, err := getAccountForUser(r.Context(), db, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
			} else {
				sendError(w, http.StatusInternalServerError, "database error")
			}
			return
		}

		sendSuccess(w, http.StatusOK, balanceResponse{
			AccountId: account.ID,
			Balance:   account.GetBalance(),
		})
	}
}

// deposit handles POST /account/deposit
func deposit(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userID := r.Header.Get("X-User-ID")
		account, err := getAccountForUser(r.Context(), db, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
			} else {
				sendError(w, http.StatusInternalServerError, "database error")
			}
			return
		}

		if err := account.Deposit(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := db.UpdateBalance(r.Context(), account.ID, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		sendSuccess(w, http.StatusOK, depositResponse{
			AccountId: account.ID,
			Balance:   account.GetBalance(),
		})
	}
}

// withdraw handles POST /account/withdraw
func withdraw(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userID := r.Header.Get("X-User-ID")
		account, err := getAccountForUser(r.Context(), db, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				sendError(w, http.StatusNotFound, "account not found")
			} else {
				sendError(w, http.StatusInternalServerError, "database error")
			}
			return
		}

		if err := account.Withdraw(req.Amount); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := db.UpdateBalance(r.Context(), account.ID, account.GetBalance()); err != nil {
			sendError(w, http.StatusInternalServerError, "failed to update balance")
			return
		}

		sendSuccess(w, http.StatusOK, withdrawResponse{
			AccountId: account.ID,
			Balance:   account.GetBalance(),
		})
	}
}

// sendError sends an error response in JSON format
func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

// sendSuccess sends a successful response in JSON format
func sendSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// getAccountForUser retrieves the account for an authenticated user
func getAccountForUser(ctx context.Context, db *store.DB, userID string) (*models.Account, error) {
	accounts, err := db.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &accounts[0], nil
}

// register handles POST /register
func register(db *store.DB) http.HandlerFunc {
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

		_, err := db.GetUserByID(req.UserId)
		if err == nil {
			sendError(w, http.StatusBadRequest, "user already exists")
			return
		}
		if err != gorm.ErrRecordNotFound {
			sendError(w, http.StatusInternalServerError, "database error")
			return
		}

		hashedPassword, err := models.HashPassword(req.Password)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to process password")
			return
		}

		// Use transaction to ensure atomicity: both User and Account created together
		err = db.WithTx(func(txDB *store.DB) error {
			user := &models.User{
				ID:       req.UserId,
				Password: hashedPassword,
			}
			if err := txDB.CreateUser(user); err != nil {
				return err
			}

			account := &models.Account{
				UserID:  req.UserId,
				Balance: 0,
			}
			if err := txDB.CreateAccount(account); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to create user and account")
			return
		}

		sendSuccess(w, http.StatusCreated, registerResponse{
			UserId:  req.UserId,
			Message: "user registered successfully",
		})
	}
}

// login handles POST /login
// Authenticates user with userId and password, returns a JWT token
func login(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse incoming JSON request body
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Find user by ID
		user, err := db.GetUserByID(req.UserId)
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

		// Generate JWT token (no database storage required - stateless)
		token, err := auth.GenerateJWT(req.UserId)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to generate token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{Token: token})
	}
}
