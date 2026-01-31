package repository

import "account-manager/internal/models"

// IAccountRepository defines the interface for account data access
type IAccountRepository interface {
	Create(account *models.Account) error
	Update(account *models.Account) error
	Delete(id uint) error
	FindByID(id uint) (*models.Account, error)
	FindByAccount(accountName string) (*models.Account, error)
	FindAll(filter models.AccountFilter) (*models.PaginatedAccounts, error)
	GetStats() (*models.AccountStats, error)
	FindExpiringAccounts(daysBefore int) ([]models.Account, error)
	MarkReminderSent(ids []uint) error
	BatchCreate(accounts []models.Account) error
}
