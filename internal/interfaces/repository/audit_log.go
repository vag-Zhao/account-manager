package repository

import (
	"time"

	"account-manager/internal/models"
)

// IAuditLogRepository defines the interface for audit log data access
type IAuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindAll(filter models.AuditLogFilter) (*models.PaginatedAuditLogs, error)
	DeleteOlderThan(date time.Time) error
	GetStats() (map[string]int64, error)
}
