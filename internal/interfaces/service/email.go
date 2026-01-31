package service

import "account-manager/internal/models"

// IEmailService defines the interface for email business logic
type IEmailService interface {
	GetConfig() (*models.EmailConfig, error)
	UpdateConfig(smtpHost string, smtpPort int, senderEmail, senderPassword, recipientEmail string, isActive bool) error
	SendEmail(subject, content string) error
	SendEmailAsync(subject, content string) <-chan error
	TestSend() error
	GetLogs(page, pageSize int) ([]models.EmailLog, int64, error)
	GetSystemConfig() (*models.SystemConfig, error)
	UpdateSystemConfig(defaultValidityDays, reminderDaysBefore int, copyFormat, emailFormat, accountTypes, accountStatuses string) error
	StopQueue()
	GetQueueSize() int
}
