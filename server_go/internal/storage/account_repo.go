package storage

import (
	"context"
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

func (r *AccountRepository) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.WithContext(ctx).Find(&accounts, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID string, newBalance int) error {
	return r.db.WithContext(ctx).Model(&models.Account{}).Where("id = ?", accountID).Update("balance", newBalance).Error
}
