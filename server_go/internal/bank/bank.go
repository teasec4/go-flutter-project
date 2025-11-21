// Package bank provides the Bank type that manages multiple accounts
package bank

import (
	"sync"
	"server/internal/account"
)

// Bank manages a collection of user accounts with thread-safe access
// Uses mutex (mu) to ensure concurrent access to accounts map is safe
type Bank struct {
	// accounts: map of account ID to Account interface
	// Protected by mutex for thread-safety
	accounts map[string]account.Account
	// mu: mutex lock for synchronizing access to accounts map
	mu sync.Mutex
}

// New creates and returns a new Bank instance
// Initializes the bank with two default accounts:
//   - Account "1" with balance 1000
//   - Account "2" with balance 2000
// Returns: pointer to newly created Bank
func New() *Bank {
	b := &Bank{
		accounts: make(map[string]account.Account),
	}
	// Initialize default accounts
	b.accounts["1"] = account.New(1000)
	b.accounts["2"] = account.New(2000)
	return b
}

// GetAccount retrieves an account by its ID in a thread-safe manner
// Uses mutex to ensure safe concurrent access
//
// Parameters:
//   - id: string, the account identifier to retrieve
//
// Returns:
//   - account.Account: the requested account (nil if not found)
//   - bool: true if account exists, false otherwise
func (b *Bank) GetAccount(id string) (account.Account, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	a, ok := b.accounts[id]
	return a, ok
}

// LockAccount retrieves an account with lock held for exclusive access
// This is designed for write operations that need atomic account state
// Currently works the same as GetAccount but reserved for future exclusive locking
//
// Parameters:
//   - id: string, the account identifier to retrieve
//
// Returns:
//   - account.Account: the requested account
//   - bool: true if account exists, false otherwise
func (b *Bank) LockAccount(id string) (account.Account, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	a, ok := b.accounts[id]
	return a, ok
}

