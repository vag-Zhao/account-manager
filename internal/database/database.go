package database

import (
	"os"
	"path/filepath"

	"account-manager/internal/cache"
	"account-manager/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize() error {
	// Get executable directory
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	// Create data directory
	dataDir := filepath.Join(execDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		// Fallback to current directory
		dataDir = "data"
		os.MkdirAll(dataDir, 0755)
	}

	dbPath := filepath.Join(dataDir, "account_manager.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Auto migrate
	err = db.AutoMigrate(
		&models.Account{},
		&models.EmailConfig{},
		&models.EmailLog{},
		&models.SystemConfig{},
		&models.ServerConfig{},
		&models.HostKey{},
		&models.AuditLog{},
	)
	if err != nil {
		return err
	}

	// Initialize default system config if not exists
	var sysConfig models.SystemConfig
	if db.First(&sysConfig).Error != nil {
		db.Create(&models.SystemConfig{
			DefaultValidityDays: 30,
			ReminderDaysBefore:  1,
		})
	}

	// Initialize default email config if not exists
	var emailConfig models.EmailConfig
	if db.First(&emailConfig).Error != nil {
		db.Create(&models.EmailConfig{
			SMTPHost:       "smtp.126.com",
			SMTPPort:       465,
			SenderEmail:    "vag3344@126.com",
			SenderPassword: "",
			RecipientEmail: "zgs3344@hunnu.edu.cn",
			IsActive:       false,
		})
	}

	DB = db

	// Initialize cache
	cache.Initialize()

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
