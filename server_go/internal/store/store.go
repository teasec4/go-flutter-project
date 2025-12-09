// Package store provides database operations with a single unified interface
package store

import (
	"context"
	"server/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps gorm.DB and provides all database operations for the application
type DB struct {
	conn *gorm.DB
}

// InitDB initializes the database connection and runs migrations
func InitDB(dbPath string) (*DB, error) {
	conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate creates tables automatically if they do not exist
	err = conn.AutoMigrate(&models.User{}, &models.Account{})
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// New creates a new DB instance from gorm.DB connection
func New(conn *gorm.DB) *DB {
	return &DB{conn: conn}
}

// ==================== USER OPERATIONS ====================

// CreateUser creates a new user in the database
func (db *DB) CreateUser(user *models.User) error {
	return db.conn.Create(user).Error
}

// CreateUserWithTx creates a user within an existing transaction
func (db *DB) CreateUserWithTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := db.conn.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ==================== ACCOUNT OPERATIONS ====================

// CreateAccount creates a new account in the database
func (db *DB) CreateAccount(account *models.Account) error {
	return db.conn.Create(account).Error
}

// CreateAccountWithTx creates an account within an existing transaction
func (db *DB) CreateAccountWithTx(tx *gorm.DB, account *models.Account) error {
	return tx.Create(account).Error
}

// GetAccountsByUserID retrieves all accounts for a user
func (db *DB) GetAccountsByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	var accounts []models.Account
	err := db.conn.WithContext(ctx).Find(&accounts, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// UpdateBalance updates the balance of an account
func (db *DB) UpdateBalance(ctx context.Context, accountID string, newBalance int) error {
	return db.conn.WithContext(ctx).Model(&models.Account{}).Where("id = ?", accountID).Update("balance", newBalance).Error
}

// ==================== TRANSACTION OPERATIONS ====================

// WithTx executes a function within a database transaction
// If fn returns an error, the transaction is automatically rolled back
// Otherwise, the transaction is committed
func (db *DB) WithTx(fn func(*DB) error) error {
	tx := db.conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := fn(&DB{conn: tx})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
