package errors

import "fmt"

// Email-specific error constructors

func NewEmailConfigFailed(err error) *AppError {
	return Wrap(err, ErrCodeEmailConfigFailed, "获取邮件配置失败")
}

func NewEmailNotEnabled() *AppError {
	return New(ErrCodeEmailNotEnabled, "邮件服务未启用")
}

func NewEmailSendFailed(err error) *AppError {
	return Wrap(err, ErrCodeEmailSendFailed, "发送邮件失败")
}

func NewEmailTimeout() *AppError {
	return New(ErrCodeEmailTimeout, "邮件发送超时")
}

// Server-specific error constructors

func NewServerConfigFailed(err error) *AppError {
	return Wrap(err, ErrCodeServerConfigFailed, "获取服务器配置失败")
}

func NewServerNotEnabled() *AppError {
	return New(ErrCodeServerNotEnabled, "服务器部署未启用")
}

func NewSSHConnectionFailed(err error) *AppError {
	return Wrap(err, ErrCodeSSHConnectionFailed, "SSH连接失败")
}

func NewEncryptionFailed(field string, err error) *AppError {
	return Wrap(err, ErrCodeEncryptionFailed, fmt.Sprintf("加密%s失败", field))
}
