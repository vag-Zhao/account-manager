package container

import (
	"account-manager/internal/cache"
	cacheInterface "account-manager/internal/interfaces/cache"
	repoInterface "account-manager/internal/interfaces/repository"
	serviceInterface "account-manager/internal/interfaces/service"
	"account-manager/internal/migration"
	"account-manager/internal/repository"
	"account-manager/internal/service"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	// Repositories
	AccountRepo   repoInterface.IAccountRepository
	EmailRepo     repoInterface.IEmailRepository
	ServerRepo    repoInterface.IServerRepository
	AuditLogRepo  repoInterface.IAuditLogRepository
	HostKeyRepo   repoInterface.IHostKeyRepository

	// Services
	AccountService  serviceInterface.IAccountService
	EmailService    serviceInterface.IEmailService
	ServerService   serviceInterface.IServerService
	AuditLogService serviceInterface.IAuditLogService
	HostKeyService  serviceInterface.IHostKeyService

	// Infrastructure
	MigrationService *migration.MigrationService
	Cache            cacheInterface.ICache
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer(db *gorm.DB) *Container {
	c := &Container{}

	// Initialize repositories
	c.AccountRepo = repository.NewAccountRepository()
	c.EmailRepo = repository.NewEmailRepository()
	c.ServerRepo = repository.NewServerRepository()
	c.AuditLogRepo = repository.NewAuditLogRepository()
	c.HostKeyRepo = repository.NewHostKeyRepository()

	// Initialize services
	c.AccountService = service.NewAccountService()
	c.EmailService = service.NewEmailService()
	c.ServerService = service.NewServerService()
	c.AuditLogService = service.NewAuditLogService()
	c.HostKeyService = service.NewHostKeyService()

	// Initialize infrastructure
	c.MigrationService = migration.NewMigrationService(db)
	c.Cache = cache.GetCache()

	return c
}

// Cleanup performs cleanup operations on container shutdown
func (c *Container) Cleanup() {
	// Stop email queue
	if emailService, ok := c.EmailService.(*service.EmailService); ok {
		emailService.StopQueue()
	}
}
