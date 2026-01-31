package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"account-manager/internal/models"
	"account-manager/internal/repository"

	"golang.org/x/crypto/ssh"
)

type HostKeyService struct {
	repo *repository.HostKeyRepository
}

func NewHostKeyService() *HostKeyService {
	return &HostKeyService{
		repo: repository.NewHostKeyRepository(),
	}
}

// VerifyHostKey creates a HostKeyCallback that verifies SSH host keys
func (s *HostKeyService) VerifyHostKey(host string, port int) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Calculate fingerprint
		fingerprint := ssh.FingerprintSHA256(key)
		keyType := key.Type()
		publicKeyBytes := key.Marshal()
		publicKeyB64 := base64.StdEncoding.EncodeToString(publicKeyBytes)

		// Check if we have this host key stored
		storedKey, err := s.repo.FindByHostAndPort(host, port)

		if err != nil {
			// First time seeing this host - need user verification
			// Store as untrusted for now
			newKey := &models.HostKey{
				Host:        host,
				Port:        port,
				KeyType:     keyType,
				Fingerprint: fingerprint,
				PublicKey:   publicKeyB64,
				FirstSeen:   time.Now(),
				LastUsed:    time.Now(),
				Trusted:     false,
			}
			s.repo.Create(newKey)

			return fmt.Errorf("UNTRUSTED_HOST_KEY|%d|%s|%s|%s", newKey.ID, host, keyType, fingerprint)
		}

		// Check if the key is trusted
		if !storedKey.Trusted {
			return fmt.Errorf("UNTRUSTED_HOST_KEY|%d|%s|%s|%s", storedKey.ID, host, keyType, storedKey.Fingerprint)
		}

		// Verify the key matches
		if storedKey.Fingerprint != fingerprint {
			return fmt.Errorf("警告: 主机密钥已更改！\n存储的指纹: %s\n当前指纹: %s\n这可能表示中间人攻击",
				storedKey.Fingerprint, fingerprint)
		}

		// Update last used time
		s.repo.UpdateLastUsed(storedKey.ID)

		return nil
	}
}

// TrustHostKey marks a host key as trusted
func (s *HostKeyService) TrustHostKey(id uint) error {
	hostKey, err := s.repo.FindByFingerprint("")
	if err != nil {
		// Find by ID instead
		keys, err := s.repo.FindAll()
		if err != nil {
			return err
		}

		for _, k := range keys {
			if k.ID == id {
				hostKey = &k
				break
			}
		}

		if hostKey == nil {
			return fmt.Errorf("host key not found")
		}
	}

	hostKey.Trusted = true
	return s.repo.Update(hostKey)
}

// GetAllHostKeys returns all stored host keys
func (s *HostKeyService) GetAllHostKeys() ([]models.HostKey, error) {
	return s.repo.FindAll()
}

// DeleteHostKey deletes a host key
func (s *HostKeyService) DeleteHostKey(id uint) error {
	return s.repo.Delete(id)
}

// GetHostKeyFingerprint returns the fingerprint for display
func GetHostKeyFingerprint(key ssh.PublicKey) string {
	hash := sha256.Sum256(key.Marshal())
	return base64.RawStdEncoding.EncodeToString(hash[:])
}
