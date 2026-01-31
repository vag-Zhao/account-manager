package repository

import (
	"time"

	"account-manager/internal/database"
	"account-manager/internal/models"

	"gorm.io/gorm"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) Create(account *models.Account) error {
	return database.GetDB().Create(account).Error
}

func (r *AccountRepository) Update(account *models.Account) error {
	return database.GetDB().Save(account).Error
}

func (r *AccountRepository) Delete(id uint) error {
	return database.GetDB().Delete(&models.Account{}, id).Error
}

func (r *AccountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	err := database.GetDB().First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByAccount(accountName string) (*models.Account, error) {
	var account models.Account
	err := database.GetDB().Where("account = ?", accountName).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindAll(filter models.AccountFilter) (*models.PaginatedAccounts, error) {
	db := database.GetDB().Model(&models.Account{})

	// Apply filters
	if filter.AccountType != "" {
		db = db.Where("account_type = ?", filter.AccountType)
	}
	if filter.IsSold != nil {
		db = db.Where("is_sold = ?", *filter.IsSold)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		db = db.Where("account LIKE ? OR notes LIKE ?", search, search)
	}

	// Count total
	var total int64
	db.Count(&total)

	// Pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	var accounts []models.Account
	err := db.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedAccounts{
		Data:       accounts,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *AccountRepository) GetStats() (*models.AccountStats, error) {
	db := database.GetDB()
	var stats models.AccountStats

	// Single aggregated query instead of 7 separate queries
	type statsResult struct {
		Total           int64
		PlusCount       int64
		BusinessCount   int64
		FreeCount       int64
		SoldCount       int64
		ExpiredCount    int64
		ExpiringIn7Days int64
	}

	now := time.Now()
	sevenDaysLater := now.AddDate(0, 0, 7)

	var result statsResult
	err := db.Model(&models.Account{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN account_type = ? THEN 1 ELSE 0 END) as plus_count,
			SUM(CASE WHEN account_type = ? THEN 1 ELSE 0 END) as business_count,
			SUM(CASE WHEN account_type = ? THEN 1 ELSE 0 END) as free_count,
			SUM(CASE WHEN is_sold = 1 THEN 1 ELSE 0 END) as sold_count,
			SUM(CASE WHEN expire_at IS NOT NULL AND expire_at < ? THEN 1 ELSE 0 END) as expired_count,
			SUM(CASE WHEN expire_at IS NOT NULL AND expire_at > ? AND expire_at <= ? THEN 1 ELSE 0 END) as expiring_in7_days
		`, models.AccountTypePLUS, models.AccountTypeBUSINESS, models.AccountTypeFREE, now, now, sevenDaysLater).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.Total = result.Total
	stats.PlusCount = result.PlusCount
	stats.BusinessCount = result.BusinessCount
	stats.FreeCount = result.FreeCount
	stats.SoldCount = result.SoldCount
	stats.ExpiredCount = result.ExpiredCount
	stats.ExpiringIn7Days = result.ExpiringIn7Days

	return &stats, nil
}

func (r *AccountRepository) FindExpiringAccounts(daysBefore int) ([]models.Account, error) {
	now := time.Now()
	targetDate := now.AddDate(0, 0, daysBefore)
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var accounts []models.Account
	err := database.GetDB().Where(
		"expire_at >= ? AND expire_at < ? AND reminder_sent = ? AND is_sold = ?",
		startOfDay, endOfDay, false, false,
	).Find(&accounts).Error

	return accounts, err
}

func (r *AccountRepository) MarkReminderSent(ids []uint) error {
	return database.GetDB().Model(&models.Account{}).Where("id IN ?", ids).Update("reminder_sent", true).Error
}

func (r *AccountRepository) BatchCreate(accounts []models.Account) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, account := range accounts {
			if err := tx.Create(&account).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
