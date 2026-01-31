package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	SMTPHost       string `json:"smtp_host"`
	SMTPPort       int    `json:"smtp_port"`
	SenderEmail    string `json:"sender_email"`
	SenderPassword string `json:"sender_password"`
	RecipientEmail string `json:"recipient_email"`
	CheckInterval  string `json:"check_interval"`
	DBPath         string `json:"db_path"`
}

type EmailService struct {
	config *Config
}

func main() {
	log.Println("账号管理系统 - 邮件服务启动中...")

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	service := &EmailService{config: config}

	// Parse check interval
	interval, err := time.ParseDuration(config.CheckInterval)
	if err != nil {
		log.Printf("解析检查间隔失败，使用默认值1小时: %v", err)
		interval = time.Hour
	}

	// Create ticker for periodic checks
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("邮件服务已启动，检查间隔: %v", interval)
	log.Printf("SMTP服务器: %s:%d", config.SMTPHost, config.SMTPPort)
	log.Printf("发件人: %s", config.SenderEmail)
	log.Printf("收件人: %s", config.RecipientEmail)

	// Send startup notification
	service.sendStartupNotification()

	// Main loop
	for {
		select {
		case <-ticker.C:
			log.Println("执行定时检查...")
			if err := service.checkAndSendReminders(); err != nil {
				log.Printf("检查失败: %v", err)
			}
		case <-sigChan:
			log.Println("收到停止信号，正在关闭服务...")
			service.sendShutdownNotification()
			return
		}
	}
}

func loadConfig() (*Config, error) {
	// Get executable directory
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("获取可执行文件路径失败: %v", err)
	}
	execDir := filepath.Dir(execPath)

	configPath := filepath.Join(execDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// Set default DB path if not specified
	if config.DBPath == "" {
		config.DBPath = filepath.Join(execDir, "email_service.db")
	}

	return &config, nil
}

func (s *EmailService) sendStartupNotification() {
	subject := "邮件服务启动通知"
	content := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2 style="color: #52c41a;">✓ 邮件服务已启动</h2>
		<p>账号管理系统的邮件服务已在远程服务器上成功启动。</p>
		<p><strong>服务信息：</strong></p>
		<ul>
			<li>启动时间: %s</li>
			<li>SMTP服务器: %s:%d</li>
			<li>发件邮箱: %s</li>
			<li>收件邮箱: %s</li>
		</ul>
		<p style="color: #666; margin-top: 20px;">
			此邮件由账号管理系统自动发送
		</p>
	</body>
	</html>
	`, time.Now().Format("2006-01-02 15:04:05"), s.config.SMTPHost, s.config.SMTPPort, s.config.SenderEmail, s.config.RecipientEmail)

	if err := s.sendEmail(subject, content); err != nil {
		log.Printf("发送启动通知失败: %v", err)
	} else {
		log.Println("启动通知已发送")
	}
}

func (s *EmailService) sendShutdownNotification() {
	subject := "邮件服务停止通知"
	content := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2 style="color: #ff4d4f;">⚠ 邮件服务已停止</h2>
		<p>账号管理系统的邮件服务已停止运行。</p>
		<p><strong>停止时间：</strong> %s</p>
		<p style="color: #666; margin-top: 20px;">
			此邮件由账号管理系统自动发送
		</p>
	</body>
	</html>
	`, time.Now().Format("2006-01-02 15:04:05"))

	if err := s.sendEmail(subject, content); err != nil {
		log.Printf("发送停止通知失败: %v", err)
	}
}

func (s *EmailService) checkAndSendReminders() error {
	// This is a placeholder for the actual check logic
	// In a real implementation, this would connect to the database
	// and check for expiring accounts
	log.Println("检查即将过期的账号...")

	// For now, just log that the check was performed
	log.Println("检查完成")
	return nil
}

func (s *EmailService) sendEmail(subject, content string) error {
	message := s.buildMessage(s.config.SenderEmail, s.config.RecipientEmail, subject, content)
	return s.sendMailSSL(s.config.SMTPHost, s.config.SMTPPort, s.config.SenderEmail, s.config.SenderPassword, s.config.RecipientEmail, message)
}

func (s *EmailService) sendMailSSL(host string, port int, from, password, to string, message []byte) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	tlsConfig := &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", from, password, host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %v", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭写入器失败: %v", err)
	}

	return client.Quit()
}

func (s *EmailService) buildMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", s.base64Encode(subject))
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	return []byte(message.String())
}

func (s *EmailService) base64Encode(str string) string {
	encoded := make([]byte, len(str)*4/3+4)
	n := 0
	for i := 0; i < len(str); i += 3 {
		b1, b2, b3 := byte(0), byte(0), byte(0)
		if i < len(str) {
			b1 = str[i]
		}
		if i+1 < len(str) {
			b2 = str[i+1]
		}
		if i+2 < len(str) {
			b3 = str[i+2]
		}

		encoded[n] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"[b1>>2]
		encoded[n+1] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"[((b1&0x03)<<4)|((b2&0xf0)>>4)]
		encoded[n+2] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"[((b2&0x0f)<<2)|((b3&0xc0)>>6)]
		encoded[n+3] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"[b3&0x3f]
		n += 4
	}
	return string(encoded[:n])
}
