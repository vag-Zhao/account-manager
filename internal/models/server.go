package models

import "time"

// ServerConfig stores remote server configuration for email service deployment
type ServerConfig struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Host           string    `json:"host" gorm:"not null"`           // Server IP or domain
	Port           int       `json:"port" gorm:"default:22"`         // SSH port
	Username       string    `json:"username" gorm:"not null"`       // SSH username
	Password       string    `json:"password"`                       // SSH password (encrypted)
	PrivateKey     string    `json:"privateKey" gorm:"type:text"`    // SSH private key (encrypted)
	DeployPath     string    `json:"deployPath" gorm:"default:'/opt/account-manager-email'"` // Deployment directory
	IsActive       bool      `json:"isActive" gorm:"default:false"`  // Whether server deployment is active
	LastDeployedAt *time.Time `json:"lastDeployedAt"`                // Last deployment time
	ServiceStatus  string    `json:"serviceStatus" gorm:"default:'unknown'"` // Service status: running, stopped, unknown
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// SMTPProvider represents common SMTP service providers
type SMTPProvider struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	HelpText string `json:"helpText"` // Instructions for getting authorization code
}

// GetSMTPProviders returns a list of common SMTP providers
func GetSMTPProviders() []SMTPProvider {
	return []SMTPProvider{
		{
			Name:     "QQ邮箱",
			Host:     "smtp.qq.com",
			Port:     465,
			HelpText: "登录QQ邮箱 -> 设置 -> 账户 -> 开启SMTP服务 -> 生成授权码（使用SSL端口465或587）",
		},
		{
			Name:     "163邮箱",
			Host:     "smtp.163.com",
			Port:     465,
			HelpText: "登录163邮箱 -> 设置 -> POP3/SMTP/IMAP -> 开启SMTP服务 -> 设置授权密码（使用SSL端口465）",
		},
		{
			Name:     "Gmail",
			Host:     "smtp.gmail.com",
			Port:     465,
			HelpText: "Enable 2FA -> Generate App Password -> Use app password (SSL port 465 or TLS port 587)",
		},
		{
			Name:     "Outlook/Hotmail",
			Host:     "smtp.office365.com",
			Port:     587,
			HelpText: "Use your account password or generate app password (STARTTLS port 587)",
		},
		{
			Name:     "126邮箱",
			Host:     "smtp.126.com",
			Port:     465,
			HelpText: "登录126邮箱 -> 设置 -> POP3/SMTP/IMAP -> 开启SMTP服务 -> 设置授权密码（使用SSL端口465）",
		},
		{
			Name:     "新浪邮箱",
			Host:     "smtp.sina.com",
			Port:     465,
			HelpText: "登录新浪邮箱 -> 设置 -> 账户 -> 开启SMTP服务（使用SSL端口465）",
		},
		{
			Name:     "企业微信邮箱",
			Host:     "smtp.exmail.qq.com",
			Port:     465,
			HelpText: "登录企业邮箱 -> 设置 -> 客户端设置 -> 开启SMTP服务 -> 生成授权码（使用SSL端口465）",
		},
	}
}

// ServerInfo contains detected server information
type ServerInfo struct {
	OSName         string `json:"osName"`         // e.g., "ubuntu", "centos", "debian"
	OSVersion      string `json:"osVersion"`      // e.g., "22.04", "8"
	OSPrettyName   string `json:"osPrettyName"`   // Full OS name from os-release
	PackageManager string `json:"packageManager"` // "apt", "yum", "dnf", "zypper", "pacman"
	HasSystemd     bool   `json:"hasSystemd"`     // Whether systemd is available
	SystemdVersion string `json:"systemdVersion"` // Systemd version if available
}
