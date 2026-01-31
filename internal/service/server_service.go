package service

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"account-manager/internal/models"
	"account-manager/internal/repository"
	"account-manager/internal/service/server"
	"account-manager/internal/utils"

	"golang.org/x/crypto/ssh"
)

type ServerService struct {
	repo           *repository.ServerRepository
	hostKeyService *HostKeyService
	sshClient      *server.SSHClient
	deployment     *server.Deployment
	serviceControl *server.ServiceControl
}

func NewServerService() *ServerService {
	hostKeyService := NewHostKeyService()
	sshClient := server.NewSSHClient(hostKeyService.VerifyHostKey)
	repo := repository.NewServerRepository()

	return &ServerService{
		repo:           repo,
		hostKeyService: hostKeyService,
		sshClient:      sshClient,
		deployment:     server.NewDeployment(sshClient),
		serviceControl: server.NewServiceControl(sshClient, repo),
	}
}

// GetConfig retrieves server configuration
func (s *ServerService) GetConfig() (*models.ServerConfig, error) {
	config, err := s.repo.GetConfig()
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive data
	if config.Password != "" {
		decrypted, err := utils.Decrypt(config.Password)
		if err == nil {
			config.Password = decrypted
		}
	}
	if config.PrivateKey != "" {
		decrypted, err := utils.Decrypt(config.PrivateKey)
		if err == nil {
			config.PrivateKey = decrypted
		}
	}

	return config, nil
}

// UpdateConfig updates server configuration
func (s *ServerService) UpdateConfig(host string, port int, username, password, privateKey, deployPath string, isActive bool) error {
	config, err := s.repo.GetConfig()
	if err != nil {
		config = &models.ServerConfig{}
	}

	config.Host = host
	config.Port = port
	config.Username = username
	config.DeployPath = deployPath
	config.IsActive = isActive

	// Encrypt password if provided
	if password != "" {
		encrypted, err := utils.Encrypt(password)
		if err != nil {
			return fmt.Errorf("加密密码失败: %v", err)
		}
		config.Password = encrypted
	}

	// Encrypt private key if provided
	if privateKey != "" {
		encrypted, err := utils.Encrypt(privateKey)
		if err != nil {
			return fmt.Errorf("加密私钥失败: %v", err)
		}
		config.PrivateKey = encrypted
	}

	return s.repo.UpdateConfig(config)
}

// TestConnection tests SSH connection to the server
func (s *ServerService) TestConnection() error {
	config, err := s.GetConfig()
	if err != nil {
		return fmt.Errorf("获取服务器配置失败: %v", err)
	}

	client, err := s.sshClient.Connect(config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Test by running a simple command
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput("echo 'Connection successful'")
	if err != nil {
		return fmt.Errorf("执行测试命令失败: %v", err)
	}

	if !strings.Contains(string(output), "Connection successful") {
		return fmt.Errorf("连接测试失败")
	}

	return nil
}

// DetectServerInfo detects server OS type, version, and systemd availability
func (s *ServerService) DetectServerInfo() (*models.ServerInfo, error) {
	config, err := s.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("获取服务器配置失败: %v", err)
	}

	client, err := s.sshClient.Connect(config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	info := &models.ServerInfo{}

	// Detect OS information from /etc/os-release
	osInfo, err := s.sshClient.RunCommandWithOutput(client, "cat /etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("读取系统信息失败: %v", err)
	}

	// Parse os-release file
	s.parseOSRelease(osInfo, info)

	// Detect package manager
	info.PackageManager = s.detectPackageManager(client)

	// Check systemd availability and version
	systemdOutput, err := s.sshClient.RunCommandWithOutput(client, "systemctl --version")
	if err == nil && strings.Contains(systemdOutput, "systemd") {
		info.HasSystemd = true
		// Parse systemd version from first line
		lines := strings.Split(systemdOutput, "\n")
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) >= 2 {
				info.SystemdVersion = parts[1]
			}
		}
	} else {
		info.HasSystemd = false
		info.SystemdVersion = ""
	}

	return info, nil
}

// parseOSRelease parses /etc/os-release content
func (s *ServerService) parseOSRelease(content string, info *models.ServerInfo) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID=") {
			info.OSName = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			info.OSVersion = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		} else if strings.HasPrefix(line, "PRETTY_NAME=") {
			info.OSPrettyName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		}
	}
}

// detectPackageManager detects the package manager available on the system
func (s *ServerService) detectPackageManager(client *ssh.Client) string {
	managers := []struct {
		name    string
		command string
	}{
		{"apt", "which apt"},
		{"yum", "which yum"},
		{"dnf", "which dnf"},
		{"zypper", "which zypper"},
		{"pacman", "which pacman"},
	}

	for _, mgr := range managers {
		if err := s.sshClient.RunCommand(client, mgr.command); err == nil {
			return mgr.name
		}
	}

	return "unknown"
}

// DeployEmailService deploys the email service to remote server
func (s *ServerService) DeployEmailService(emailConfig *models.EmailConfig) error {
	serverConfig, err := s.GetConfig()
	if err != nil {
		return fmt.Errorf("获取服务器配置失败: %v", err)
	}

	if !serverConfig.IsActive {
		return fmt.Errorf("服务器部署未启用")
	}

	// Connect to server
	client, err := s.sshClient.Connect(serverConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create deployment directory
	if err := s.sshClient.RunCommand(client, fmt.Sprintf("mkdir -p %s", serverConfig.DeployPath)); err != nil {
		return fmt.Errorf("创建部署目录失败: %v", err)
	}

	// Build the email service binary
	binaryPath, err := s.deployment.BuildEmailService()
	if err != nil {
		return fmt.Errorf("构建邮件服务失败: %v", err)
	}
	defer os.Remove(binaryPath)

	// Upload binary to server
	if err := s.deployment.UploadFile(client, binaryPath, path.Join(serverConfig.DeployPath, "email-service")); err != nil {
		return fmt.Errorf("上传服务文件失败: %v", err)
	}

	// Create configuration file
	configContent := s.deployment.GenerateConfigFile(emailConfig)
	if err := s.deployment.UploadContent(client, configContent, path.Join(serverConfig.DeployPath, "config.json")); err != nil {
		return fmt.Errorf("上传配置文件失败: %v", err)
	}

	// Create systemd service file
	serviceContent := s.deployment.GenerateSystemdService(serverConfig.DeployPath)
	if err := s.deployment.UploadContent(client, serviceContent, "/tmp/account-manager-email.service"); err != nil {
		return fmt.Errorf("上传服务配置失败: %v", err)
	}

	// Install and enable systemd service
	commands := []string{
		fmt.Sprintf("chmod +x %s/email-service", serverConfig.DeployPath),
		"sudo mv /tmp/account-manager-email.service /etc/systemd/system/",
		"sudo systemctl daemon-reload",
		"sudo systemctl enable account-manager-email",
		"sudo systemctl restart account-manager-email",
	}

	for _, cmd := range commands {
		if err := s.sshClient.RunCommand(client, cmd); err != nil {
			return fmt.Errorf("执行命令失败 [%s]: %v", cmd, err)
		}
	}

	// Update deployment time
	now := time.Now()
	serverConfig.LastDeployedAt = &now
	serverConfig.ServiceStatus = "running"
	s.repo.UpdateConfig(serverConfig)

	return nil
}

// GetServiceStatus checks the status of the email service on remote server
func (s *ServerService) GetServiceStatus() (string, error) {
	config, err := s.GetConfig()
	if err != nil {
		return "unknown", err
	}

	if !config.IsActive {
		return "disabled", nil
	}

	client, err := s.sshClient.Connect(config)
	if err != nil {
		return "unknown", err
	}
	defer client.Close()

	return s.serviceControl.GetStatus(client, config)
}

// StopService stops the email service on remote server
func (s *ServerService) StopService() error {
	config, err := s.GetConfig()
	if err != nil {
		return err
	}

	client, err := s.sshClient.Connect(config)
	if err != nil {
		return err
	}
	defer client.Close()

	return s.serviceControl.Stop(client, config)
}

// StartService starts the email service on remote server
func (s *ServerService) StartService() error {
	config, err := s.GetConfig()
	if err != nil {
		return err
	}

	client, err := s.sshClient.Connect(config)
	if err != nil {
		return err
	}
	defer client.Close()

	return s.serviceControl.Start(client, config)
}
