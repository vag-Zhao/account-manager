package utils

import (
	"regexp"
	"strings"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateRequired checks if a string is not empty
func ValidateRequired(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fieldName, "不能为空")
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length requirement
func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return NewValidationErrorf(fieldName, "长度至少为%d个字符", minLength)
	}
	return nil
}

// ValidateMaxLength checks if a string doesn't exceed maximum length
func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return NewValidationErrorf(fieldName, "长度不能超过%d个字符", maxLength)
	}
	return nil
}

// ValidateRange checks if a number is within a range
func ValidateRange(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return NewValidationErrorf(fieldName, "必须在%d到%d之间", min, max)
	}
	return nil
}

// ValidatePort checks if a port number is valid
func ValidatePort(port int) error {
	return ValidateRange(port, 1, 65535, "端口")
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// NewValidationErrorf creates a new validation error with formatted message
func NewValidationErrorf(field, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: sprintf(format, args...),
	}
}

// sprintf is a simple format function to avoid importing fmt
func sprintf(format string, args ...interface{}) string {
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			result = strings.Replace(result, "%d", itoa(v), 1)
		case string:
			result = strings.Replace(result, "%s", v, 1)
		}
	}
	return result
}

// itoa converts int to string without fmt package
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
