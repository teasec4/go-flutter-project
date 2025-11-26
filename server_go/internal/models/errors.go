package models

import "errors"

var (
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
