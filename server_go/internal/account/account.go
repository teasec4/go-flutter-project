// Package account provides the Account interface and its implementation
package account

import "errors"

// Account defines the interface for bank account operations
// All account implementations must satisfy this interface
type Account interface {
	// Deposit adds funds to the account
	// Returns error if amount is not positive
	Deposit(amount int) error

	// Withdraw removes funds from the account
	// Returns error if amount is not positive or exceeds balance
	Withdraw(amount int) error

	// GetBalance returns the current account balance
	GetBalance() int
}

// impl is the concrete implementation of the Account interface
// This is the actual account storage with balance tracking
type impl struct {
	balance int // current account balance in currency units
}

// New creates and returns a new Account with the specified initial balance
// Parameters:
//   - initialBalance: int, the starting balance for the new account
// Returns: Account interface pointing to the new account implementation
func New(initialBalance int) Account {
	return &impl{balance: initialBalance}
}

// Deposit adds the specified amount to the account balance
// Validates that the deposit amount is positive before updating balance
//
// Parameters:
//   - amount: int, the amount to deposit (must be positive)
//
// Returns:
//   - nil: if deposit was successful
//   - error: if amount is not positive (amount <= 0)
func (a *impl) Deposit(amount int) error {
	if amount <= 0 {
		return errors.New("deposit amount must be greater than 0")
	}
	a.balance += amount
	return nil
}

// Withdraw removes the specified amount from the account balance
// Validates that:
//   1. The withdrawal amount is positive
//   2. The account has sufficient balance
//
// Parameters:
//   - amount: int, the amount to withdraw (must be positive)
//
// Returns:
//   - nil: if withdrawal was successful
//   - error: if amount is not positive OR if balance is insufficient
func (a *impl) Withdraw(amount int) error {
	if amount <= 0 {
		return errors.New("withdraw amount should be positive")
	}
	if a.balance < amount {
		return errors.New("insufficient balance")
	}
	a.balance -= amount
	return nil
}

// GetBalance returns the current balance of the account
// This method is safe to call concurrently (read-only operation)
// Note: For guaranteed consistency with concurrent deposits/withdrawals,
//       the caller should use Bank.GetAccount() which holds the mutex
//
// Returns: int, the current account balance
func (a *impl) GetBalance() int {
	return a.balance
}

