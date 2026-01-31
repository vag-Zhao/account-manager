package migration

import (
	"errors"

	"gorm.io/gorm"
)

// SystemConfig stores system-wide configuration flags
type SystemConfig struct {
	ID                uint   `gorm:"primaryKey"`
	Key               string `gorm:"uniqueIndex;not null"`
	Value             string `gorm:"type:text"`
	EncryptionMigrated bool  `gorm:"default:false"`
}

// MigrationService handles data migration operations
type MigrationService struct {
	db *gorm.DB
}

func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{db: db}
}

// IsEncryptionMigrated checks if encryption migration has been completed
func (s *MigrationService) IsEncryptionMigrated() (bool, error) {
	var config SystemConfig
	err := s.db.Where("key = ?", "encryption_migrated").First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return config.EncryptionMigrated, nil
}

// MigrateEncryption marks encryption migration as complete
// No longer needed with fixed encryption key
func (s *MigrationService) MigrateEncryption() error {
	return s.MarkEncryptionMigrated()
}

// MarkEncryptionMigrated marks the encryption migration as complete
func (s *MigrationService) MarkEncryptionMigrated() error {
	config := SystemConfig{
		Key:                "encryption_migrated",
		Value:              "true",
		EncryptionMigrated: true,
	}
	return s.db.Save(&config).Error
}

// EnsureMigrationTableExists creates the system_config table if it doesn't exist
func (s *MigrationService) EnsureMigrationTableExists() error {
	return s.db.AutoMigrate(&SystemConfig{})
}
