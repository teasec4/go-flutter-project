package bank

import (
	"sync"
	"server/internal/account"
)

// Bank manage accounts
type Bank struct{
	accounts map[string]account.Account
	mu sync.Mutex
}

// New create a new bank
func New() *Bank{
	b:= &Bank{
		accounts: make(map[string]account.Account),
	}
	// create a default account with id 
	b.accounts["1"] = account.New(1000)
	b.accounts["2"] = account.New(2000)
	return b
}

// get account an account by ID
func(b *Bank) GetAccount(id string) (account.Account, bool){
	b.mu.Lock()
	defer b.mu.Unlock()
	a, ok := b.accounts[id]
	return a, ok
}

// lock account for write operations
func(b *Bank) LockAccount(id string) (account.Account, bool){
	b.mu.Lock()
	defer b.mu.Unlock()
	a, ok := b.accounts[id]
	return a, ok
}
