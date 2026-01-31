package repository

import "account-manager/internal/models"

// IHostKeyRepository defines the interface for host key data access
type IHostKeyRepository interface {
	FindByHostAndPort(host string, port int) (*models.HostKey, error)
	FindByFingerprint(fingerprint string) (*models.HostKey, error)
	Create(hostKey *models.HostKey) error
	Update(hostKey *models.HostKey) error
	UpdateLastUsed(id uint) error
	FindAll() ([]models.HostKey, error)
	Delete(id uint) error
}
