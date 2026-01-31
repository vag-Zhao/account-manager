package repository

import (
	"account-manager/internal/database"
	"account-manager/internal/models"
)

type EmailRepository struct{}

func NewEmailRepository() *EmailRepository {
	return &EmailRepository{}
}

func (r *EmailRepository) GetConfig() (*models.EmailConfig, error) {
	var config models.EmailConfig
	err := database.GetDB().First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *EmailRepository) UpdateConfig(config *models.EmailConfig) error {
	return database.GetDB().Save(config).Error
}

func (r *EmailRepository) CreateLog(log *models.EmailLog) error {
	return database.GetDB().Create(log).Error
}

func (r *EmailRepository) GetLogs(page, pageSize int) ([]models.EmailLog, int64, error) {
	var logs []models.EmailLog
	var total int64

	db := database.GetDB().Model(&models.EmailLog{})
	db.Count(&total)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (r *EmailRepository) GetSystemConfig() (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := database.GetDB().First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *EmailRepository) UpdateSystemConfig(config *models.SystemConfig) error {
	return database.GetDB().Save(config).Error
}
