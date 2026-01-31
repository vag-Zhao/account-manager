package service

import (
	"time"

	"account-manager/internal/models"
)

// IAccountService defines the interface for account business logic
type IAccountService interface {
	CreateAccount(account string, password string, accountType string, expireAt *time.Time, notes string, isSold bool) error
	UpdateAccount(id uint, account string, password string, accountType string, expireAt *time.Time, notes string, isSold bool) error
	DeleteAccount(id uint) error
	GetAccount(id uint) (*models.Account, error)
	GetAccounts(filter models.AccountFilter) (*models.PaginatedAccounts, error)
	GetStats() (*models.AccountStats, error)
	MarkAsSold(id uint) error
	MarkAsUnsold(id uint) error
	BatchImport(accounts []map[string]interface{}) (int, []string)
	DecryptPassword(id uint) (string, error)
}
