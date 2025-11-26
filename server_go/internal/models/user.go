package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey"`
	Password  string
	Account   Account `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
