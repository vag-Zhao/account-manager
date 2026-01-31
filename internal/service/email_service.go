package service

import (
	"crypto/tls"
	"fmt"
	"time"

	"account-manager/internal/config"
	"account-manager/internal/logger"
	"account-manager/internal/models"
	"account-manager/internal/queue"
	"account-manager/internal/repository"
	"account-manager/internal/utils"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	repo  *repository.EmailRepository
	queue *queue.EmailQueue
}

func NewEmailService() *EmailService {
	service := &EmailService{
		repo: repository.NewEmailRepository(),
	}

	// Initialize queue with the actual send function
	cfg := config.Get()
	service.queue = queue.NewEmailQueue(cfg.Worker.EmailWorkers, service.sendEmailSync)
	service.queue.Start()

	return service
}

func (s *EmailService) GetConfig() (*models.EmailConfig, error) {
	config, err := s.repo.GetConfig()
	if err != nil {
		return nil, err
	}

	// Decrypt password for display
	if config.SenderPassword != "" {
		decrypted, err := utils.Decrypt(config.SenderPassword)
		if err == nil {
			config.SenderPassword = decrypted
		}
	}

	return config, nil
}

func (s *EmailService) UpdateConfig(smtpHost string, smtpPort int, senderEmail, senderPassword, recipientEmail string, isActive bool) error {
	config, err := s.repo.GetConfig()
	if err != nil {
		config = &models.EmailConfig{}
	}

	config.SMTPHost = smtpHost
	config.SMTPPort = smtpPort
	config.SenderEmail = senderEmail
	config.RecipientEmail = recipientEmail
	config.IsActive = isActive

	// Encrypt password if provided
	if senderPassword != "" {
		encrypted, err := utils.Encrypt(senderPassword)
		if err != nil {
			return err
		}
		config.SenderPassword = encrypted
	}

	return s.repo.UpdateConfig(config)
}

// SendEmail sends an email asynchronously using the queue
func (s *EmailService) SendEmail(subject, content string) error {
	// Enqueue the email
	resultChan := s.queue.Enqueue(subject, content)

	// Wait for result with timeout
	select {
	case err := <-resultChan:
		return err
	case <-time.After(60 * time.Second):
		return fmt.Errorf("邮件发送超时")
	}
}

// SendEmailAsync sends an email asynchronously and returns immediately
func (s *EmailService) SendEmailAsync(subject, content string) <-chan error {
	return s.queue.Enqueue(subject, content)
}

// sendEmailSync is the synchronous email sending implementation
func (s *EmailService) sendEmailSync(subject, content string) error {
	config, err := s.repo.GetConfig()
	if err != nil {
		return fmt.Errorf("获取邮件配置失败: %v", err)
	}

	if !config.IsActive {
		return fmt.Errorf("邮件服务未启用")
	}

	// Decrypt password
	password := ""
	if config.SenderPassword != "" {
		decrypted, err := utils.Decrypt(config.SenderPassword)
		if err != nil {
			return fmt.Errorf("解密密码失败: %v", err)
		}
		password = decrypted
	}

	logger.WithFields(map[string]interface{}{
		"smtp_host":  config.SMTPHost,
		"smtp_port":  config.SMTPPort,
		"sender":     config.SenderEmail,
		"recipient":  config.RecipientEmail,
	}).Debug("Sending email via gomail")

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", config.SenderEmail)
	m.SetHeader("To", config.RecipientEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// Create dialer with SSL/TLS
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SenderEmail, password)

	// Configure TLS with proper certificate verification
	d.TLSConfig = &tls.Config{
		ServerName: config.SMTPHost,
		MinVersion: tls.VersionTLS12,
	}

	// For port 465, use SSL
	if config.SMTPPort == 465 {
		d.SSL = true
		logger.Debug("Using SSL connection (port 465)")
	} else {
		d.SSL = false
		logger.WithField("port", config.SMTPPort).Debug("Using STARTTLS connection")
	}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		logger.WithFields(map[string]interface{}{
			"error":     err.Error(),
			"recipient": config.RecipientEmail,
			"subject":   subject,
		}).Error("Email send failed")

		// Log the result
		log := &models.EmailLog{
			Subject:   subject,
			Content:   content,
			Recipient: config.RecipientEmail,
			Status:    "failed",
			Error:     err.Error(),
		}
		s.repo.CreateLog(log)

		return fmt.Errorf("发送邮件失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"recipient": config.RecipientEmail,
		"subject":   subject,
	}).Info("Email sent successfully")

	// Log success
	log := &models.EmailLog{
		Subject:   subject,
		Content:   content,
		Recipient: config.RecipientEmail,
		Status:    "success",
	}
	s.repo.CreateLog(log)

	return nil
}

// StopQueue stops the email queue (call on shutdown)
func (s *EmailService) StopQueue() {
	if s.queue != nil {
		s.queue.Stop()
	}
}

// GetQueueSize returns the current queue size
func (s *EmailService) GetQueueSize() int {
	if s.queue != nil {
		return s.queue.GetQueueSize()
	}
	return 0
}

func (s *EmailService) TestSend() error {
	subject := "账号管理系统 - 测试邮件"
	content := `
	<html>
	<body>
		<h2>测试邮件</h2>
		<p>这是一封来自账号管理系统的测试邮件。</p>
		<p>如果您收到此邮件，说明邮件配置正确。</p>
		<p>发送时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
	</body>
	</html>
	`
	return s.SendEmail(subject, content)
}

func (s *EmailService) GetLogs(page, pageSize int) ([]models.EmailLog, int64, error) {
	return s.repo.GetLogs(page, pageSize)
}

func (s *EmailService) GetSystemConfig() (*models.SystemConfig, error) {
	return s.repo.GetSystemConfig()
}

func (s *EmailService) UpdateSystemConfig(defaultValidityDays, reminderDaysBefore int, copyFormat, emailFormat, accountTypes, accountStatuses string) error {
	config, err := s.repo.GetSystemConfig()
	if err != nil {
		config = &models.SystemConfig{}
	}

	config.DefaultValidityDays = defaultValidityDays
	config.ReminderDaysBefore = reminderDaysBefore
	if copyFormat != "" {
		config.CopyFormat = copyFormat
	}
	if emailFormat != "" {
		config.EmailFormat = emailFormat
	}
	if accountTypes != "" {
		config.AccountTypes = accountTypes
	}
	if accountStatuses != "" {
		config.AccountStatuses = accountStatuses
	}

	return s.repo.UpdateSystemConfig(config)
}
