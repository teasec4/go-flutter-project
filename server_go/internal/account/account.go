package account

import "errors"

type Account interface {
	Deposit(amount int) error
	Withdraw(amount int) error
	GetBalance() int
}

type impl struct {
	balance int
}

// new acc
func New(initialBalance int) Account {
	return &impl{balance: initialBalance}
}

// Deposit implements Account.
func (a *impl) Deposit(amount int) error {
	if amount <= 0 {
		return errors.New("deposit amount must be greater than 0")
	}
	a.balance += amount
	return nil
}

// Withdraw implements Account.
func (a *impl) Withdraw(amount int) error {
	if amount <= 0 {
		return  errors.New("withdraw amount should be positive")
	}
	if a.balance < amount {
		return errors.New("insufficient balance")
	}
	a.balance -= amount
	return nil
}


func (a *impl) GetBalance() int {
	return a.balance
}
