package models

import "time"

type Token struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    // прямая связь на User
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
