package service

import (
	"errors"
	"sync"
	"time"

	"account-manager/internal/cache"
	"account-manager/internal/config"
	"account-manager/internal/logger"
	"account-manager/internal/models"
	"account-manager/internal/repository"
	"account-manager/internal/utils"
)

type AccountService struct {
	repo       *repository.AccountRepository
	emailRepo  *repository.EmailRepository
	auditLog   *AuditLogService
}

func NewAccountService() *AccountService {
	return &AccountService{
		repo:      repository.NewAccountRepository(),
		emailRepo: repository.NewEmailRepository(),
		auditLog:  NewAuditLogService(),
	}
}

func (s *AccountService) CreateAccount(account string, password string, accountType string, expireAt *time.Time, notes string, isSold bool) error {
	if account == "" {
		return errors.New("账号不能为空")
	}

	// Check if account already exists
	existing, _ := s.repo.FindByAccount(account)
	if existing != nil {
		return errors.New("账号已存在")
	}

	// Encrypt password
	encryptedPassword := ""
	if password != "" {
		var err error
		encryptedPassword, err = utils.Encrypt(password)
		if err != nil {
			return err
		}
	}

	// Calculate expire date
	var finalExpireAt *time.Time
	accType := models.AccountType(accountType)

	if accType != models.AccountTypeFREE {
		if expireAt != nil {
			finalExpireAt = expireAt
		} else {
			// Default: 30 days from now
			sysConfig, _ := s.emailRepo.GetSystemConfig()
			days := 30
			if sysConfig != nil {
				days = sysConfig.DefaultValidityDays
			}
			expire := time.Now().AddDate(0, 0, days)
			finalExpireAt = &expire
		}
	}

	var soldAt *time.Time
	if isSold {
		now := time.Now()
		soldAt = &now
	}

	newAccount := &models.Account{
		Account:     account,
		Password:    encryptedPassword,
		AccountType: accType,
		ExpireAt:    finalExpireAt,
		Notes:       notes,
		IsSold:      isSold,
		SoldAt:      soldAt,
	}

	err := s.repo.Create(newAccount)
	if err == nil {
		// Invalidate stats cache after creating account
		cache.InvalidateStats()

		// Audit log
		s.auditLog.LogAccountCreate(newAccount.ID, "user", account)
	}
	return err
}

func (s *AccountService) UpdateAccount(id uint, account string, password string, accountType string, expireAt *time.Time, notes string, isSold bool) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("账号不存在")
	}

	// Check if new account name conflicts with another account
	if account != existing.Account {
		conflict, _ := s.repo.FindByAccount(account)
		if conflict != nil {
			return errors.New("账号名已被使用")
		}
	}

	existing.Account = account
	existing.AccountType = models.AccountType(accountType)
	existing.Notes = notes

	// Update isSold status
	if isSold != existing.IsSold {
		existing.IsSold = isSold
		if isSold {
			now := time.Now()
			existing.SoldAt = &now
		} else {
			existing.SoldAt = nil
		}
	}

	// Update password if provided
	if password != "" {
		encryptedPassword, err := utils.Encrypt(password)
		if err != nil {
			return err
		}
		existing.Password = encryptedPassword
	}

	// Update expire date
	if models.AccountType(accountType) == models.AccountTypeFREE {
		existing.ExpireAt = nil
	} else if expireAt != nil {
		existing.ExpireAt = expireAt
	}

	err = s.repo.Update(existing)
	if err == nil {
		// Invalidate stats cache after updating account
		cache.InvalidateStats()

		// Audit log
		changes := map[string]interface{}{
			"account": account,
			"type":    accountType,
		}
		s.auditLog.LogAccountUpdate(id, "user", changes)
	}
	return err
}

func (s *AccountService) DeleteAccount(id uint) error {
	// Get account info before deletion for audit log
	account, _ := s.repo.FindByID(id)
	accountName := ""
	if account != nil {
		accountName = account.Account
	}

	err := s.repo.Delete(id)
	if err == nil {
		// Invalidate stats cache after deleting account
		cache.InvalidateStats()

		// Audit log
		s.auditLog.LogAccountDelete(id, "user", accountName)
	}
	return err
}

func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Decrypt password for display
	if account.Password != "" {
		decrypted, err := utils.Decrypt(account.Password)
		if err == nil {
			account.Password = decrypted
		}
	}

	return account, nil
}

func (s *AccountService) GetAccounts(filter models.AccountFilter) (*models.PaginatedAccounts, error) {
	result, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	// Batch decrypt passwords using goroutine pool
	s.batchDecrypt(result.Data)

	return result, nil
}

// batchDecrypt decrypts passwords in parallel using a goroutine pool
func (s *AccountService) batchDecrypt(accounts []models.Account) {
	cfg := config.Get()
	maxWorkers := cfg.Worker.DecryptionWorkers
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for i := range accounts {
		if accounts[i].Password == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			decrypted, err := utils.Decrypt(accounts[idx].Password)
			if err == nil {
				accounts[idx].Password = decrypted
			} else {
				logger.WithFields(map[string]interface{}{
					"account_id": accounts[idx].ID,
					"error":      err.Error(),
				}).Error("Failed to decrypt password")
			}
		}(i)
	}
	wg.Wait()
}

func (s *AccountService) GetStats() (*models.AccountStats, error) {
	// Try to get from cache first
	c := cache.GetCache()
	if cached, found := c.Get(cache.KeyStats); found {
		if stats, ok := cached.(*models.AccountStats); ok {
			return stats, nil
		}
	}

	// Cache miss, fetch from database
	stats, err := s.repo.GetStats()
	if err != nil {
		return nil, err
	}

	// Store in cache with 5 minute TTL
	c.Set(cache.KeyStats, stats, cache.TTLStats)

	return stats, nil
}

func (s *AccountService) MarkAsSold(id uint) error {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("账号不存在")
	}

	now := time.Now()
	account.IsSold = true
	account.SoldAt = &now

	err = s.repo.Update(account)
	if err == nil {
		// Invalidate stats cache after marking as sold
		cache.InvalidateStats()
	}
	return err
}

func (s *AccountService) MarkAsUnsold(id uint) error {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("账号不存在")
	}

	account.IsSold = false
	account.SoldAt = nil

	err = s.repo.Update(account)
	if err == nil {
		// Invalidate stats cache after marking as unsold
		cache.InvalidateStats()
	}
	return err
}

func (s *AccountService) BatchImport(accounts []map[string]interface{}) (int, []string) {
	successCount := 0
	var errors []string

	for _, acc := range accounts {
		account, _ := acc["account"].(string)
		password, _ := acc["password"].(string)
		accountType, _ := acc["accountType"].(string)
		isSold, _ := acc["isSold"].(bool)

		var expireAt *time.Time
		if expireStr, ok := acc["expireAt"].(string); ok && expireStr != "" {
			t, err := time.Parse("2006-01-02", expireStr)
			if err == nil {
				expireAt = &t
			}
		}

		err := s.CreateAccount(account, password, accountType, expireAt, "", isSold)
		if err != nil {
			errors = append(errors, account+": "+err.Error())
		} else {
			successCount++
		}
	}

	// Invalidate stats cache after batch import
	if successCount > 0 {
		cache.InvalidateStats()
	}

	return successCount, errors
}

// DecryptPassword decrypts a single password on-demand
func (s *AccountService) DecryptPassword(id uint) (string, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return "", errors.New("账号不存在")
	}

	if account.Password == "" {
		return "", nil
	}

	decrypted, err := utils.Decrypt(account.Password)
	if err != nil {
		return "", errors.New("解密失败")
	}

	// Audit log - password access
	s.auditLog.LogPasswordView(id, "user", "decrypt")

	return decrypted, nil
}
