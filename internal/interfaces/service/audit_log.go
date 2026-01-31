package service

import "account-manager/internal/models"

// IAuditLogService defines the interface for audit logging business logic
type IAuditLogService interface {
	Log(action, resourceType string, resourceID uint, user string, details map[string]interface{}, success bool, errorMsg string) error
	LogAccountView(accountID uint, user string) error
	LogAccountCreate(accountID uint, user string, account string) error
	LogAccountUpdate(accountID uint, user string, changes map[string]interface{}) error
	LogAccountDelete(accountID uint, user string, account string) error
	LogPasswordView(accountID uint, user string, action string) error
	LogLogin(user string, success bool, errorMsg string) error
	LogLogout(user string) error
	LogConfigChange(configType string, user string, changes map[string]interface{}) error
	GetLogs(filter models.AuditLogFilter) (*models.PaginatedAuditLogs, error)
	CleanupOldLogs(retentionDays int) error
	GetStats() (map[string]int64, error)
	ExportToCSV(filter models.AuditLogFilter) (string, error)
}
