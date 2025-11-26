package storage

import (
	"server/internal/models"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) GetAccountByID(id string) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetAccountsByUserID(userID string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Find(&accounts, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(accountID string, newBalance int) error {
	return r.db.Model(&models.Account{}).Where("id = ?", accountID).Update("balance", newBalance).Error
}
