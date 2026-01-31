package main

import (
	"context"
	"os"
	"time"

	"account-manager/internal/config"
	"account-manager/internal/database"
	"account-manager/internal/logger"
	"account-manager/internal/migration"
	"account-manager/internal/models"
	"account-manager/internal/scheduler"
	"account-manager/internal/service"
)

// App struct
type App struct {
	ctx              context.Context
	accountService   *service.AccountService
	emailService     *service.EmailService
	serverService    *service.ServerService
	auditService     *service.AuditLogService
	hostKeyService   *service.HostKeyService
	migrationService *migration.MigrationService
	scheduler        *scheduler.Scheduler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize logger
	logLevel := logger.InfoLevel
	logger.Initialize(logLevel, logger.NewMultiWriter(os.Stdout))

	// Load configuration
	loader := config.NewLoader("config.yaml")
	cfg, err := loader.LoadOrCreate()
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Failed to load config, using defaults")
	} else {
		logger.Info("Configuration loaded successfully")
		if cfg.App.Debug {
			logger.SetLevel(logger.DebugLevel)
			logger.Info("Debug mode enabled")
		}
	}

	// Initialize database
	database.Initialize()

	// Initialize migration service
	a.migrationService = migration.NewMigrationService(database.GetDB())

	// Ensure migration table exists
	if err := a.migrationService.EnsureMigrationTableExists(); err != nil {
		// Log error but continue - this is not critical for startup
		logger.WithField("error", err.Error()).Warn("Failed to create migration table")
	}

	// Initialize services
	a.accountService = service.NewAccountService()
	a.emailService = service.NewEmailService()
	a.serverService = service.NewServerService()
	a.auditService = service.NewAuditLogService()
	a.hostKeyService = service.NewHostKeyService()

	// Initialize and start scheduler
	a.scheduler = scheduler.NewScheduler()
	a.scheduler.Start()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.scheduler != nil {
		a.scheduler.Stop()
	}
}

// ============ Account Methods ============

func (a *App) CreateAccount(account, password, accountType string, expireAt string, isSold bool) error {
	var expireTime *time.Time
	if expireAt != "" {
		t, err := time.Parse("2006-01-02", expireAt)
		if err == nil {
			expireTime = &t
		}
	}
	return a.accountService.CreateAccount(account, password, accountType, expireTime, "", isSold)
}

func (a *App) UpdateAccount(id uint, account, password, accountType string, expireAt string, isSold bool) error {
	var expireTime *time.Time
	if expireAt != "" {
		t, err := time.Parse("2006-01-02", expireAt)
		if err == nil {
			expireTime = &t
		}
	}
	return a.accountService.UpdateAccount(id, account, password, accountType, expireTime, "", isSold)
}

func (a *App) DeleteAccount(id uint) error {
	return a.accountService.DeleteAccount(id)
}

func (a *App) GetAccount(id uint) (*models.Account, error) {
	return a.accountService.GetAccount(id)
}

func (a *App) GetAccounts(accountType string, isSold *bool, search string, page, pageSize int) (*models.PaginatedAccounts, error) {
	filter := models.AccountFilter{
		AccountType: accountType,
		IsSold:      isSold,
		Search:      search,
		Page:        page,
		PageSize:    pageSize,
	}
	return a.accountService.GetAccounts(filter)
}

func (a *App) GetStats() (*models.AccountStats, error) {
	return a.accountService.GetStats()
}

func (a *App) MarkAsSold(id uint) error {
	return a.accountService.MarkAsSold(id)
}

func (a *App) MarkAsUnsold(id uint) error {
	return a.accountService.MarkAsUnsold(id)
}

func (a *App) BatchImport(accounts []map[string]interface{}) *models.BatchImportResult {
	success, errors := a.accountService.BatchImport(accounts)
	return &models.BatchImportResult{
		Success: success,
		Errors:  errors,
	}
}

// ============ Email Methods ============

func (a *App) GetEmailConfig() (*models.EmailConfig, error) {
	return a.emailService.GetConfig()
}

func (a *App) UpdateEmailConfig(smtpHost string, smtpPort int, senderEmail, senderPassword, recipientEmail string, isActive bool) error {
	return a.emailService.UpdateConfig(smtpHost, smtpPort, senderEmail, senderPassword, recipientEmail, isActive)
}

func (a *App) TestEmailSend() error {
	return a.emailService.TestSend()
}

func (a *App) GetEmailLogs(page, pageSize int) *models.EmailLogsResult {
	logs, total, _ := a.emailService.GetLogs(page, pageSize)
	return &models.EmailLogsResult{
		Logs:  logs,
		Total: total,
	}
}

func (a *App) GetSystemConfig() (*models.SystemConfig, error) {
	return a.emailService.GetSystemConfig()
}

func (a *App) UpdateSystemConfig(defaultValidityDays, reminderDaysBefore int, copyFormat, emailFormat, accountTypes, accountStatuses string) error {
	return a.emailService.UpdateSystemConfig(defaultValidityDays, reminderDaysBefore, copyFormat, emailFormat, accountTypes, accountStatuses)
}

func (a *App) ManualCheckExpiry() (int, error) {
	return a.scheduler.ManualCheck()
}

// ============ Server Methods ============

func (a *App) GetServerConfig() (*models.ServerConfig, error) {
	return a.serverService.GetConfig()
}

func (a *App) UpdateServerConfig(host string, port int, username, password, privateKey, deployPath string, isActive bool) error {
	return a.serverService.UpdateConfig(host, port, username, password, privateKey, deployPath, isActive)
}

func (a *App) TestServerConnection() error {
	return a.serverService.TestConnection()
}

func (a *App) DeployEmailService() error {
	emailConfig, err := a.emailService.GetConfig()
	if err != nil {
		return err
	}
	return a.serverService.DeployEmailService(emailConfig)
}

func (a *App) GetServiceStatus() (string, error) {
	return a.serverService.GetServiceStatus()
}

func (a *App) StopEmailService() error {
	return a.serverService.StopService()
}

func (a *App) StartEmailService() error {
	return a.serverService.StartService()
}

func (a *App) GetSMTPProviders() []models.SMTPProvider {
	return models.GetSMTPProviders()
}

func (a *App) DetectServerInfo() (*models.ServerInfo, error) {
	return a.serverService.DetectServerInfo()
}

// ============ Audit Log Methods ============

func (a *App) GetAuditLogs(action, resourceType string, startDate, endDate string, page, pageSize int) (*models.PaginatedAuditLogs, error) {
	filter := models.AuditLogFilter{
		Action:       action,
		ResourceType: resourceType,
		Page:         page,
		PageSize:     pageSize,
	}

	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			filter.StartDate = t
		}
	}

	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			filter.EndDate = t
		}
	}

	return a.auditService.GetLogs(filter)
}

func (a *App) GetAuditStats() (map[string]int64, error) {
	return a.auditService.GetStats()
}

func (a *App) ExportAuditLogs(action, resourceType string, startDate, endDate string) (string, error) {
	filter := models.AuditLogFilter{
		Action:       action,
		ResourceType: resourceType,
	}

	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			filter.StartDate = t
		}
	}

	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			filter.EndDate = t
		}
	}

	return a.auditService.ExportToCSV(filter)
}

func (a *App) CleanupOldAuditLogs(retentionDays int) error {
	return a.auditService.CleanupOldLogs(retentionDays)
}

// ============ Host Key Methods ============

func (a *App) GetAllHostKeys() ([]models.HostKey, error) {
	return a.hostKeyService.GetAllHostKeys()
}

func (a *App) TrustHostKey(id uint) error {
	return a.hostKeyService.TrustHostKey(id)
}

func (a *App) DeleteHostKey(id uint) error {
	return a.hostKeyService.DeleteHostKey(id)
}

// ============ Migration Methods ============

func (a *App) IsEncryptionMigrated() bool {
	migrated, _ := a.migrationService.IsEncryptionMigrated()
	return migrated
}

func (a *App) MigrateEncryption() error {
	return a.migrationService.MigrateEncryption()
}
