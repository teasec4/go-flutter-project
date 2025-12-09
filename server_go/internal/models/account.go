package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Error definitions for account operations
var (
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Account struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"unique;not null"`
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate automatically generates a UUID for new Account records
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (a *Account) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.Balance < amount {
		return ErrInsufficientBalance
	}
	a.Balance -= amount
	return nil
}

func (a *Account) GetBalance() int {
	return a.Balance
}
