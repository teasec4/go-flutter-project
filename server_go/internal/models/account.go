package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"unique;not null"`
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate автоматически генерирует ID для нового Account
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
