package repository

import (
	"time"

	"account-manager/internal/database"
	"account-manager/internal/models"
)

type AuditLogRepository struct{}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	db := database.GetDB()
	return db.Create(log).Error
}

// FindAll returns paginated audit logs with filters
func (r *AuditLogRepository) FindAll(filter models.AuditLogFilter) (*models.PaginatedAuditLogs, error) {
	db := database.GetDB()

	query := db.Model(&models.AuditLog{})

	// Apply filters
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}
	if !filter.StartDate.IsZero() {
		query = query.Where("timestamp >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("timestamp <= ?", filter.EndDate)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	// Fetch data
	var logs []models.AuditLog
	err := query.Order("timestamp DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return &models.PaginatedAuditLogs{
		Data:       logs,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// DeleteOlderThan deletes audit logs older than the specified date
func (r *AuditLogRepository) DeleteOlderThan(date time.Time) error {
	db := database.GetDB()
	return db.Where("timestamp < ?", date).Delete(&models.AuditLog{}).Error
}

// GetStats returns audit log statistics
func (r *AuditLogRepository) GetStats() (map[string]int64, error) {
	db := database.GetDB()
	stats := make(map[string]int64)

	// Total logs
	var total int64
	db.Model(&models.AuditLog{}).Count(&total)
	stats["total"] = total

	// Logs by action
	type ActionCount struct {
		Action string
		Count  int64
	}
	var actionCounts []ActionCount
	db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&actionCounts)

	for _, ac := range actionCounts {
		stats[ac.Action] = ac.Count
	}

	return stats, nil
}
