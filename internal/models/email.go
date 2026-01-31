package models

import (
	"time"
)

type EmailConfig struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	SMTPHost       string `json:"smtpHost" gorm:"not null"`
	SMTPPort       int    `json:"smtpPort" gorm:"not null"`
	SenderEmail    string `json:"senderEmail" gorm:"not null"`
	SenderPassword string `json:"senderPassword"`
	RecipientEmail string `json:"recipientEmail" gorm:"not null"`
	IsActive       bool   `json:"isActive" gorm:"default:true"`
	UseRemoteServer bool  `json:"useRemoteServer" gorm:"default:false"` // Whether to use remote server for email service
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type EmailLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Subject   string    `json:"subject"`
	Content   string    `json:"content"`
	Recipient string    `json:"recipient"`
	Status    string    `json:"status"` // success, failed
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"createdAt"`
}

type SystemConfig struct {
	ID                  uint   `json:"id" gorm:"primaryKey"`
	DefaultValidityDays int    `json:"defaultValidityDays" gorm:"default:30"`
	ReminderDaysBefore  int    `json:"reminderDaysBefore" gorm:"default:1"`
	CopyFormat          string `json:"copyFormat" gorm:"default:'账号：{account}\n密码：{password}'"`
	EmailFormat         string `json:"emailFormat" gorm:"default:'您的账号 {account} 将在 {expireAt} 过期，请及时处理。'"`
	AccountTypes        string `json:"accountTypes" gorm:"type:text"`
	AccountStatuses     string `json:"accountStatuses" gorm:"type:text"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type EmailLogsResult struct {
	Logs  []EmailLog `json:"logs"`
	Total int64      `json:"total"`
}

type BatchImportResult struct {
	Success int      `json:"success"`
	Errors  []string `json:"errors"`
}
