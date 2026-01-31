package server

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"account-manager/internal/logger"
	"account-manager/internal/models"
	"account-manager/internal/utils"

	"golang.org/x/crypto/ssh"
)

// Deployment handles deployment operations
type Deployment struct {
	sshClient *SSHClient
}

// NewDeployment creates a new deployment helper
func NewDeployment(sshClient *SSHClient) *Deployment {
	return &Deployment{
		sshClient: sshClient,
	}
}

// UploadFile uploads a local file to the remote server
func (d *Deployment) UploadFile(client *ssh.Client, localPath, remotePath string) error {
	// First, ensure the remote directory exists
	remoteDir := path.Dir(remotePath)
	logger.WithField("remote_dir", remoteDir).Debug("Creating remote directory")
	if err := d.sshClient.RunCommand(client, fmt.Sprintf("mkdir -p %s", remoteDir)); err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}

	// Check if scp is available
	logger.Debug("Checking if SCP command is available")
	if err := d.sshClient.RunCommand(client, "which scp"); err != nil {
		return fmt.Errorf("远程服务器未安装SCP命令: %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"local_path":  localPath,
		"remote_path": remotePath,
		"size":        stat.Size(),
	}).Debug("Uploading file")

	// Capture stderr for debugging
	var stderr bytes.Buffer
	session.Stderr = &stderr

	// Get stdin pipe
	w, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("获取stdin管道失败: %v", err)
	}

	// Start the SCP command
	logger.WithField("remote_path", remotePath).Debug("Starting SCP command")
	if err := session.Start(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return fmt.Errorf("启动SCP命令失败: %v", err)
	}

	// Send file header
	logger.WithFields(map[string]interface{}{
		"size":     stat.Size(),
		"filename": path.Base(remotePath),
	}).Debug("Sending file header")
	fmt.Fprintf(w, "C0755 %d %s\n", stat.Size(), path.Base(remotePath))

	// Copy file content
	logger.Debug("Copying file content")
	copied, err := io.Copy(w, file)
	if err != nil {
		w.Close()
		session.Wait()
		return fmt.Errorf("复制文件内容失败: %v", err)
	}
	logger.WithField("bytes", copied).Debug("File content copied")

	// Send termination byte
	logger.Debug("Sending termination byte")
	fmt.Fprint(w, "\x00")
	w.Close()

	// Wait for command to complete
	logger.Debug("Waiting for SCP command to complete")
	if err := session.Wait(); err != nil {
		stderrStr := stderr.String()
		if stderrStr != "" {
			return fmt.Errorf("SCP上传失败: %v\nstderr: %s", err, stderrStr)
		}
		return fmt.Errorf("SCP上传失败: %v", err)
	}

	logger.WithField("remote_path", remotePath).Info("File uploaded successfully")

	// Verify file exists
	logger.Debug("Verifying file exists")
	if err := d.sshClient.RunCommand(client, fmt.Sprintf("test -f %s", remotePath)); err != nil {
		return fmt.Errorf("文件上传后验证失败，文件不存在: %v", err)
	}

	logger.WithField("remote_path", remotePath).Info("File upload verified successfully")
	return nil
}

// UploadContent uploads content as a file to the remote server
func (d *Deployment) UploadContent(client *ssh.Client, content, remotePath string) error {
	// Ensure remote directory exists
	remoteDir := path.Dir(remotePath)
	if err := d.sshClient.RunCommand(client, fmt.Sprintf("mkdir -p %s", remoteDir)); err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C0644 %d %s\n", len(content), path.Base(remotePath))
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()

	if err := session.Run(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return err
	}

	return nil
}

// BuildEmailService builds the email service binary for Linux
func (d *Deployment) BuildEmailService() (string, error) {
	// Create temporary directory for build
	tmpDir, err := os.MkdirTemp("", "email-service-build")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "email-service")

	// Check if email service source exists
	emailServicePath := "cmd/email-service"
	if _, err := os.Stat(emailServicePath); os.IsNotExist(err) {
		return "", fmt.Errorf("邮件服务源代码不存在: %s", emailServicePath)
	}

	// Build for Linux AMD64
	logger.Info("Building email service (Linux AMD64)")
	cmd := exec.Command("go", "build", "-o", outputPath, "./cmd/email-service")
	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("构建失败: %v\n输出: %s", err, string(output))
	}

	logger.WithField("output_path", outputPath).Info("Build successful")
	return outputPath, nil
}

// GenerateConfigFile generates the email service configuration file content
func (d *Deployment) GenerateConfigFile(emailConfig *models.EmailConfig) string {
	// Decrypt password for config file
	password := emailConfig.SenderPassword
	if password != "" {
		decrypted, err := utils.Decrypt(password)
		if err == nil {
			password = decrypted
		}
	}

	return fmt.Sprintf(`{
  "smtp_host": "%s",
  "smtp_port": %d,
  "sender_email": "%s",
  "sender_password": "%s",
  "recipient_email": "%s",
  "check_interval": "1h"
}`, emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.SenderEmail, password, emailConfig.RecipientEmail)
}

// GenerateSystemdService generates the systemd service file content
func (d *Deployment) GenerateSystemdService(deployPath string) string {
	return fmt.Sprintf(`[Unit]
Description=Account Manager Email Service
After=network.target

[Service]
Type=simple
User=emailservice
Group=emailservice
WorkingDirectory=%s
ExecStart=%s/email-service
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
`, deployPath, deployPath)
}
