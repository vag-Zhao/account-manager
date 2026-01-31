package repository

import "account-manager/internal/models"

// IServerRepository defines the interface for server config data access
type IServerRepository interface {
	GetConfig() (*models.ServerConfig, error)
	UpdateConfig(config *models.ServerConfig) error
}
