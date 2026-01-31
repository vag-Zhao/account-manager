package server

import (
	"bytes"
	"fmt"
	"time"

	"account-manager/internal/models"

	"golang.org/x/crypto/ssh"
)

// HostKeyVerifier is a function type for host key verification
type HostKeyVerifier func(host string, port int) ssh.HostKeyCallback

// SSHClient handles SSH connection and command execution
type SSHClient struct {
	hostKeyVerifier HostKeyVerifier
}

// NewSSHClient creates a new SSH client helper
func NewSSHClient(hostKeyVerifier HostKeyVerifier) *SSHClient {
	return &SSHClient{
		hostKeyVerifier: hostKeyVerifier,
	}
}

// Connect establishes an SSH connection to the server
func (c *SSHClient) Connect(config *models.ServerConfig) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// Try password authentication
	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	// Try private key authentication
	if config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("未提供有效的认证方式")
	}

	clientConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethods,
		HostKeyCallback: c.hostKeyVerifier(config.Host, config.Port),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	return client, nil
}

// RunCommand executes a command on the remote server
func (c *SSHClient) RunCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var stderr bytes.Buffer
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}

	return nil
}

// RunCommandWithOutput executes a command and returns its output
func (c *SSHClient) RunCommandWithOutput(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
