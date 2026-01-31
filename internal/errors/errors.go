package errors

import "fmt"

// ErrorCode represents a unique error code
type ErrorCode string

const (
	// Account errors
	ErrCodeAccountEmpty        ErrorCode = "ACCOUNT_EMPTY"
	ErrCodeAccountExists       ErrorCode = "ACCOUNT_EXISTS"
	ErrCodeAccountNotFound     ErrorCode = "ACCOUNT_NOT_FOUND"
	ErrCodeAccountNameInUse    ErrorCode = "ACCOUNT_NAME_IN_USE"

	// Authentication errors
	ErrCodePasswordTooShort    ErrorCode = "PASSWORD_TOO_SHORT"
	ErrCodeInvalidPassword     ErrorCode = "INVALID_PASSWORD"
	ErrCodeDecryptionFailed    ErrorCode = "DECRYPTION_FAILED"

	// Email errors
	ErrCodeEmailConfigFailed   ErrorCode = "EMAIL_CONFIG_FAILED"
	ErrCodeEmailNotEnabled     ErrorCode = "EMAIL_NOT_ENABLED"
	ErrCodeEmailSendFailed     ErrorCode = "EMAIL_SEND_FAILED"
	ErrCodeEmailTimeout        ErrorCode = "EMAIL_TIMEOUT"

	// Server errors
	ErrCodeServerConfigFailed  ErrorCode = "SERVER_CONFIG_FAILED"
	ErrCodeServerNotEnabled    ErrorCode = "SERVER_NOT_ENABLED"
	ErrCodeSSHConnectionFailed ErrorCode = "SSH_CONNECTION_FAILED"
	ErrCodeEncryptionFailed    ErrorCode = "ENCRYPTION_FAILED"

	// Validation errors
	ErrCodeValidationFailed    ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidInput        ErrorCode = "INVALID_INPUT"
)

// AppError represents an application error with code and message
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with code and message
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is checks if an error matches a specific error code
func Is(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}
