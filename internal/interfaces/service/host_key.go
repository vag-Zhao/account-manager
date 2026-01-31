package service

import (
	"net"

	"account-manager/internal/models"
	"golang.org/x/crypto/ssh"
)

// IHostKeyService defines the interface for SSH host key management
type IHostKeyService interface {
	VerifyHostKey(host string, port int) ssh.HostKeyCallback
	TrustHostKey(id uint) error
	GetAllHostKeys() ([]models.HostKey, error)
	DeleteHostKey(id uint) error
}

// HostKeyCallback is the function signature for SSH host key verification
type HostKeyCallback func(hostname string, remote net.Addr, key ssh.PublicKey) error
