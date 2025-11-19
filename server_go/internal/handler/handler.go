package handler

import (
	"encoding/json"
	"net/http"
	"server/internal/account"
)

type balanceResponse struct {
	AccountId string  `json:"accountId"`
	Balance   float64 `json:"balance"`
}

// Balance handles GET /account?id={id}
func Balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get query param ?id=1
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	a, ok := account.GetByID(id)
	if !ok {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	resp := balanceResponse{AccountId: id, Balance: a.GetBalance()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
