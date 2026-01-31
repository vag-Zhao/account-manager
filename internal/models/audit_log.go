package models

import "time"

// AuditLog records all sensitive operations
type AuditLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Timestamp    time.Time `json:"timestamp" gorm:"index:idx_timestamp"`
	User         string    `json:"user" gorm:"type:varchar(255)"` // Username or "system"
	Action       string    `json:"action" gorm:"type:varchar(100);not null;index:idx_action"` // view, create, update, delete, login, logout
	ResourceType string    `json:"resourceType" gorm:"type:varchar(50);not null"` // account, email_config, system_config
	ResourceID   uint      `json:"resourceId" gorm:"index:idx_resource"`
	IPAddress    string    `json:"ipAddress" gorm:"type:varchar(50)"`
	Details      string    `json:"details" gorm:"type:text"` // JSON string with additional details
	Success      bool      `json:"success" gorm:"default:true"`
	ErrorMessage string    `json:"errorMessage" gorm:"type:text"`
	CreatedAt    time.Time `json:"createdAt"`
}

// AuditLogFilter for querying audit logs
type AuditLogFilter struct {
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	Page         int       `json:"page"`
	PageSize     int       `json:"pageSize"`
}

// PaginatedAuditLogs represents paginated audit log results
type PaginatedAuditLogs struct {
	Data       []AuditLog `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	TotalPages int        `json:"totalPages"`
}
