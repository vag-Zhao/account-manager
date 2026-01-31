package repository

import (
	"account-manager/internal/database"
	"account-manager/internal/models"
)

type ServerRepository struct{}

func NewServerRepository() *ServerRepository {
	return &ServerRepository{}
}

func (r *ServerRepository) GetConfig() (*models.ServerConfig, error) {
	var config models.ServerConfig
	result := database.DB.First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

func (r *ServerRepository) UpdateConfig(config *models.ServerConfig) error {
	if config.ID == 0 {
		return database.DB.Create(config).Error
	}
	return database.DB.Save(config).Error
}
