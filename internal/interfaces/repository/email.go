package repository

import "account-manager/internal/models"

// IEmailRepository defines the interface for email data access
type IEmailRepository interface {
	GetConfig() (*models.EmailConfig, error)
	UpdateConfig(config *models.EmailConfig) error
	CreateLog(log *models.EmailLog) error
	GetLogs(page, pageSize int) ([]models.EmailLog, int64, error)
	GetSystemConfig() (*models.SystemConfig, error)
	UpdateSystemConfig(config *models.SystemConfig) error
}
