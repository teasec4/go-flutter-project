package handler

import (
	"encoding/json"
	"net/http"
	"server/internal/bank"
)

// REQUEST
type depositRequest struct {
	AccountId string `json:"accountId"`
	Amount    int    `json:"amount"`
}

type withdrawRequest struct {
	AccountId string `json:"accountId"`
	Amount    int `json:"amount"`
}

type errorResponse struct{
	Error string `json:"error"`
}

// RESPONSE
type depositResponse struct{
	AccountId string `json:"accountId"`
	Balance int `json:"balance"`
}

type withdrawResponse struct{
	AccountId string `json:"accountId"`
	Balance int `json:"balance"`
}

type balanceResponse struct {
	AccountId string  `json:"accountId"`
	Balance   int `json:"balance"`
}

// Balance handles GET /account?id={id}
func Balance(b *bank.Bank) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		if r.Method != http.MethodGet{
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return 
		}

		accountIDStr := r.URL.Query().Get("id")
		if accountIDStr == ""{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: "account_id parameter required"})
			return
		}

		a, ok := b.GetAccount(accountIDStr)
		if !ok{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{Error: "account not found"})
			return
		} 

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balanceResponse{
			AccountId: accountIDStr, 
			Balance: a.GetBalance(),
		})
	}
}

// Deposit handles POST
func Deposit(b *bank.Bank) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost{
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req depositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: "invalid request"})
			return
		}
		a, ok := b.GetAccount(req.AccountId)
		if !ok{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{Error: "account not found"})
			return
		}

		if err := a.Deposit(req.Amount); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(depositResponse{
			AccountId: req.AccountId,
			Balance: a.GetBalance(),
		})
	}
}