package models

import (
	"time"
)

type AccountType string

const (
	AccountTypePLUS     AccountType = "PLUS"
	AccountTypeBUSINESS AccountType = "BUSINESS"
	AccountTypeFREE     AccountType = "FREE"
)

type Account struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	Account      string      `json:"account" gorm:"uniqueIndex;not null"`
	Password     string      `json:"password"`
	AccountType  AccountType `json:"accountType" gorm:"type:varchar(20);not null;index:idx_type_sold"`
	IsSold       bool        `json:"isSold" gorm:"default:false;index:idx_type_sold"`
	SoldAt       *time.Time  `json:"soldAt"`
	ExpireAt     *time.Time  `json:"expireAt" gorm:"index:idx_expire"`
	ReminderSent bool        `json:"reminderSent" gorm:"default:false"`
	Notes        string      `json:"notes"`
	CreatedAt    time.Time   `json:"createdAt" gorm:"index:idx_created_desc"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

type AccountFilter struct {
	AccountType string `json:"accountType"`
	IsSold      *bool  `json:"isSold"`
	Search      string `json:"search"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

type AccountStats struct {
	Total           int64 `json:"total"`
	PlusCount       int64 `json:"plusCount"`
	BusinessCount   int64 `json:"businessCount"`
	FreeCount       int64 `json:"freeCount"`
	SoldCount       int64 `json:"soldCount"`
	ExpiredCount    int64 `json:"expiredCount"`
	ExpiringIn7Days int64 `json:"expiringIn7Days"`
}

type PaginatedAccounts struct {
	Data       []Account `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
	TotalPages int       `json:"totalPages"`
}
