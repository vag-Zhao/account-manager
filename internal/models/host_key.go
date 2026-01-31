package models

import "time"

// HostKey stores SSH host key fingerprints for verification
type HostKey struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Host        string    `json:"host" gorm:"type:varchar(255);not null"`
	Port        int       `json:"port" gorm:"not null"`
	KeyType     string    `json:"keyType" gorm:"type:varchar(50);not null"` // ssh-rsa, ecdsa-sha2-nistp256, etc.
	Fingerprint string    `json:"fingerprint" gorm:"type:varchar(255);not null;uniqueIndex:idx_host_port_fingerprint"`
	PublicKey   string    `json:"publicKey" gorm:"type:text;not null"` // Base64 encoded public key
	FirstSeen   time.Time `json:"firstSeen"`
	LastUsed    time.Time `json:"lastUsed"`
	Trusted     bool      `json:"trusted" gorm:"default:false"` // User has verified this key
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
