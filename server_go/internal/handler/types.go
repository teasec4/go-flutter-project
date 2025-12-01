// Package handler defines HTTP request handlers for the bank API
package handler

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

// loginRequest represents the incoming JSON payload for login
type loginRequest struct {
	UserId    string `json:"userId"`
	Password  string `json:"password"`
}

// registerRequest represents the incoming JSON payload for user registration
type registerRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
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

// loginResponse represents the JSON response after successful login
type loginResponse struct {
	Token string `json:"token"`
}

// registerResponse represents the JSON response after successful registration
type registerResponse struct {
	UserId   string `json:"userId"`
	Message  string `json:"message"`
}
