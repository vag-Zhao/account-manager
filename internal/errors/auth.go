package errors

// Authentication-specific error constructors

func NewPasswordTooShort(minLength int) *AppError {
	return New(ErrCodePasswordTooShort, "密码长度不足")
}

func NewInvalidPassword() *AppError {
	return New(ErrCodeInvalidPassword, "密码错误")
}
