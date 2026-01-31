package service

import (
	"encoding/json"
	"time"

	"account-manager/internal/models"
	"account-manager/internal/repository"
)

type AuditLogService struct {
	repo *repository.AuditLogRepository
}

func NewAuditLogService() *AuditLogService {
	return &AuditLogService{
		repo: repository.NewAuditLogRepository(),
	}
}

// Log creates a new audit log entry
func (s *AuditLogService) Log(action, resourceType string, resourceID uint, user string, details map[string]interface{}, success bool, errorMsg string) error {
	detailsJSON := ""
	if details != nil {
		bytes, _ := json.Marshal(details)
		detailsJSON = string(bytes)
	}

	log := &models.AuditLog{
		Timestamp:    time.Now(),
		User:         user,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    "127.0.0.1", // Desktop app - always localhost
		Details:      detailsJSON,
		Success:      success,
		ErrorMessage: errorMsg,
	}

	return s.repo.Create(log)
}

// LogAccountView logs account viewing
func (s *AuditLogService) LogAccountView(accountID uint, user string) error {
	return s.Log("view", "account", accountID, user, nil, true, "")
}

// LogAccountCreate logs account creation
func (s *AuditLogService) LogAccountCreate(accountID uint, user string, account string) error {
	details := map[string]interface{}{
		"account": account,
	}
	return s.Log("create", "account", accountID, user, details, true, "")
}

// LogAccountUpdate logs account updates
func (s *AuditLogService) LogAccountUpdate(accountID uint, user string, changes map[string]interface{}) error {
	return s.Log("update", "account", accountID, user, changes, true, "")
}

// LogAccountDelete logs account deletion
func (s *AuditLogService) LogAccountDelete(accountID uint, user string, account string) error {
	details := map[string]interface{}{
		"account": account,
	}
	return s.Log("delete", "account", accountID, user, details, true, "")
}

// LogPasswordView logs password viewing/copying
func (s *AuditLogService) LogPasswordView(accountID uint, user string, action string) error {
	details := map[string]interface{}{
		"action": action, // "view" or "copy"
	}
	return s.Log("password_access", "account", accountID, user, details, true, "")
}

// LogLogin logs user login
func (s *AuditLogService) LogLogin(user string, success bool, errorMsg string) error {
	return s.Log("login", "auth", 0, user, nil, success, errorMsg)
}

// LogLogout logs user logout
func (s *AuditLogService) LogLogout(user string) error {
	return s.Log("logout", "auth", 0, user, nil, true, "")
}

// LogConfigChange logs configuration changes
func (s *AuditLogService) LogConfigChange(configType string, user string, changes map[string]interface{}) error {
	return s.Log("config_change", configType, 0, user, changes, true, "")
}

// GetLogs returns paginated audit logs
func (s *AuditLogService) GetLogs(filter models.AuditLogFilter) (*models.PaginatedAuditLogs, error) {
	return s.repo.FindAll(filter)
}

// CleanupOldLogs deletes logs older than the retention period (default 90 days)
func (s *AuditLogService) CleanupOldLogs(retentionDays int) error {
	if retentionDays <= 0 {
		retentionDays = 90
	}
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.repo.DeleteOlderThan(cutoffDate)
}

// GetStats returns audit log statistics
func (s *AuditLogService) GetStats() (map[string]int64, error) {
	return s.repo.GetStats()
}

// ExportToCSV exports audit logs to CSV format
func (s *AuditLogService) ExportToCSV(filter models.AuditLogFilter) (string, error) {
	// Set a large page size to get all logs
	filter.PageSize = 10000
	filter.Page = 1

	logs, err := s.repo.FindAll(filter)
	if err != nil {
		return "", err
	}

	csv := "ID,Timestamp,User,Action,ResourceType,ResourceID,IPAddress,Details,Success,ErrorMessage\n"
	for _, log := range logs.Data {
		csv += formatCSVRow(
			log.ID,
			log.Timestamp.Format("2006-01-02 15:04:05"),
			log.User,
			log.Action,
			log.ResourceType,
			log.ResourceID,
			log.IPAddress,
			log.Details,
			log.Success,
			log.ErrorMessage,
		)
	}

	return csv, nil
}

func formatCSVRow(values ...interface{}) string {
	row := ""
	for i, v := range values {
		if i > 0 {
			row += ","
		}
		row += escapeCSV(v)
	}
	return row + "\n"
}

func escapeCSV(value interface{}) string {
	str := ""
	switch v := value.(type) {
	case string:
		str = v
	case uint:
		str = string(rune(v))
	case bool:
		if v {
			str = "true"
		} else {
			str = "false"
		}
	default:
		str = ""
	}

	// Escape quotes and wrap in quotes if contains comma or newline
	if len(str) > 0 && (str[0] == '"' || containsSpecialChar(str)) {
		str = `"` + replaceAll(str, `"`, `""`) + `"`
	}
	return str
}

func containsSpecialChar(s string) bool {
	for _, c := range s {
		if c == ',' || c == '\n' || c == '\r' {
			return true
		}
	}
	return false
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
