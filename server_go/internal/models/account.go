package models

import "time"

type Account struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"unique;not null"`
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
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
