// Package main API types
package main

// BalanceResponse represents the response returned when querying account balance
// Used by: GET /account?id={accountId}
type BalanceResponse struct {
	AccountId string `json:"accountId"` // Account identifier
	Balance   int    `json:"balance"`   // Current account balance
}

// DepositResponse represents the response after successfully depositing funds
// Used by: POST /account/deposit
// Shows the updated balance after deposit operation
type DepositResponse struct {
	AccountId string `json:"accountId"` // Account identifier
	Balance   int    `json:"balance"`   // Updated account balance after deposit
}

// WithdrawResponse represents the response after successfully withdrawing funds
// Used by: POST /account/withdraw
// Shows the updated balance after withdrawal operation
type WithdrawResponse struct {
	AccountId string `json:"accountId"` // Account identifier
	Balance   int    `json:"balance"`   // Updated account balance after withdrawal
}

// ErrorResponse represents an error response returned when operations fail
// Used by: all endpoints when errors occur
type ErrorResponse struct {
	Error string `json:"error"` // Error message describing what went wrong
}
