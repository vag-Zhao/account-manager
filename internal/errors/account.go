package errors

// Account-specific error constructors

func NewAccountEmpty() *AppError {
	return New(ErrCodeAccountEmpty, "账号不能为空")
}

func NewAccountExists() *AppError {
	return New(ErrCodeAccountExists, "账号已存在")
}

func NewAccountNotFound() *AppError {
	return New(ErrCodeAccountNotFound, "账号不存在")
}

func NewAccountNameInUse() *AppError {
	return New(ErrCodeAccountNameInUse, "账号名已被使用")
}

func NewDecryptionFailed(err error) *AppError {
	return Wrap(err, ErrCodeDecryptionFailed, "解密失败")
}
