package service

import "account-manager/internal/models"

// IServerService defines the interface for server deployment business logic
type IServerService interface {
	GetConfig() (*models.ServerConfig, error)
	UpdateConfig(host string, port int, username, password, privateKey, deployPath string, isActive bool) error
	TestConnection() error
	DetectServerInfo() (*models.ServerInfo, error)
	DeployEmailService(emailConfig *models.EmailConfig) error
	GetServiceStatus() (string, error)
	StopService() error
	StartService() error
}
