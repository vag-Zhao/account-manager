package utils

import (
	"time"
)

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ParseDateTime parses a datetime string in YYYY-MM-DD HH:MM:SS format
func ParseDateTime(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FormatDate formats a time pointer to YYYY-MM-DD string
func FormatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateTime formats a time pointer to YYYY-MM-DD HH:MM:SS string
func FormatDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// IsExpired checks if a date is in the past
func IsExpired(t *time.Time) bool {
	if t == nil {
		return false
	}
	return t.Before(time.Now())
}

// IsExpiringSoon checks if a date is within the specified number of days
func IsExpiringSoon(t *time.Time, days int) bool {
	if t == nil {
		return false
	}
	now := time.Now()
	targetDate := now.AddDate(0, 0, days)
	return t.After(now) && t.Before(targetDate)
}

// DaysUntil returns the number of days until the given date
func DaysUntil(t *time.Time) int {
	if t == nil {
		return 0
	}
	duration := t.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// StartOfDay returns the start of the day for the given time
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for the given time
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// ParseDateRange parses start and end date strings
func ParseDateRange(startStr, endStr string) (start, end time.Time, err error) {
	if startStr != "" {
		startPtr, err := ParseDate(startStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		if startPtr != nil {
			start = StartOfDay(*startPtr)
		}
	}

	if endStr != "" {
		endPtr, err := ParseDate(endStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		if endPtr != nil {
			end = EndOfDay(*endPtr)
		}
	}

	return start, end, nil
}

// AddDays adds the specified number of days to a time
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// TruncateToDate truncates a time to just the date (removes time component)
func TruncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
