package server

import (
	"fmt"
	"strings"

	"account-manager/internal/models"
	"account-manager/internal/repository"

	"golang.org/x/crypto/ssh"
)

// ServiceControl handles service start/stop/status operations
type ServiceControl struct {
	sshClient *SSHClient
	repo      *repository.ServerRepository
}

// NewServiceControl creates a new service control helper
func NewServiceControl(sshClient *SSHClient, repo *repository.ServerRepository) *ServiceControl {
	return &ServiceControl{
		sshClient: sshClient,
		repo:      repo,
	}
}

// GetStatus checks the status of the email service on remote server
func (sc *ServiceControl) GetStatus(client *ssh.Client, config *models.ServerConfig) (string, error) {
	if !config.IsActive {
		return "disabled", nil
	}

	session, err := client.NewSession()
	if err != nil {
		return "unknown", err
	}
	defer session.Close()

	output, err := session.CombinedOutput("sudo systemctl is-active account-manager-email")
	status := strings.TrimSpace(string(output))

	if status == "active" {
		config.ServiceStatus = "running"
	} else {
		config.ServiceStatus = "stopped"
	}
	sc.repo.UpdateConfig(config)

	return config.ServiceStatus, nil
}

// Stop stops the email service on remote server
func (sc *ServiceControl) Stop(client *ssh.Client, config *models.ServerConfig) error {
	if err := sc.sshClient.RunCommand(client, "sudo systemctl stop account-manager-email"); err != nil {
		return fmt.Errorf("停止服务失败: %v", err)
	}

	config.ServiceStatus = "stopped"
	sc.repo.UpdateConfig(config)

	return nil
}

// Start starts the email service on remote server
func (sc *ServiceControl) Start(client *ssh.Client, config *models.ServerConfig) error {
	if err := sc.sshClient.RunCommand(client, "sudo systemctl start account-manager-email"); err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	config.ServiceStatus = "running"
	sc.repo.UpdateConfig(config)

	return nil
}
