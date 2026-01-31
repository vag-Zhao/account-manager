package repository

import (
	"account-manager/internal/database"
	"account-manager/internal/models"
	"time"
)

type HostKeyRepository struct{}

func NewHostKeyRepository() *HostKeyRepository {
	return &HostKeyRepository{}
}

// FindByHostAndPort finds a trusted host key for the given host and port
func (r *HostKeyRepository) FindByHostAndPort(host string, port int) (*models.HostKey, error) {
	db := database.GetDB()
	var hostKey models.HostKey

	err := db.Where("host = ? AND port = ? AND trusted = ?", host, port, true).First(&hostKey).Error
	if err != nil {
		return nil, err
	}

	return &hostKey, nil
}

// FindByFingerprint finds a host key by its fingerprint
func (r *HostKeyRepository) FindByFingerprint(fingerprint string) (*models.HostKey, error) {
	db := database.GetDB()
	var hostKey models.HostKey

	err := db.Where("fingerprint = ?", fingerprint).First(&hostKey).Error
	if err != nil {
		return nil, err
	}

	return &hostKey, nil
}

// Create creates a new host key record
func (r *HostKeyRepository) Create(hostKey *models.HostKey) error {
	db := database.GetDB()
	return db.Create(hostKey).Error
}

// Update updates a host key record
func (r *HostKeyRepository) Update(hostKey *models.HostKey) error {
	db := database.GetDB()
	return db.Save(hostKey).Error
}

// UpdateLastUsed updates the last used timestamp
func (r *HostKeyRepository) UpdateLastUsed(id uint) error {
	db := database.GetDB()
	return db.Model(&models.HostKey{}).Where("id = ?", id).Update("last_used", time.Now()).Error
}

// FindAll returns all host keys
func (r *HostKeyRepository) FindAll() ([]models.HostKey, error) {
	db := database.GetDB()
	var hostKeys []models.HostKey

	err := db.Order("last_used DESC").Find(&hostKeys).Error
	if err != nil {
		return nil, err
	}

	return hostKeys, nil
}

// Delete deletes a host key
func (r *HostKeyRepository) Delete(id uint) error {
	db := database.GetDB()
	return db.Delete(&models.HostKey{}, id).Error
}
