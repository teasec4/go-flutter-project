package storage

import (
	"server/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция - создаёт таблицы если их нет
	err = db.AutoMigrate(&models.User{}, &models.Account{}, &models.Token{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
