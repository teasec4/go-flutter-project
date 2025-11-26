package storage

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// CreateToken сохраняет новый токен в БД
func (r *TokenRepository) CreateToken(token *models.Token) error {
	return r.db.Create(token).Error
}

// GetTokenByValue находит токен по значению и проверяет что он не истёк
func (r *TokenRepository) GetTokenByValue(tokenValue string) (*models.Token, error) {
	var token models.Token
	err := r.db.Where("token = ? AND expires_at > ?", tokenValue, gorm.Expr("datetime('now')")).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteToken удаляет токен по значению
func (r *TokenRepository) DeleteToken(tokenValue string) error {
	return r.db.Where("token = ?", tokenValue).Delete(&models.Token{}).Error
}

// DeleteExpiredTokens удаляет все истёкшие токены
func (r *TokenRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at <= ?", gorm.Expr("datetime('now')")).
		Delete(&models.Token{}).Error
}
