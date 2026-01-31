package cache

import (
	"time"

	"account-manager/internal/config"
	cacheInterface "account-manager/internal/interfaces/cache"

	"github.com/patrickmn/go-cache"
)

var (
	// Global cache instance
	Cache *CacheWrapper
)

// CacheWrapper wraps go-cache to implement ICache interface
type CacheWrapper struct {
	cache *cache.Cache
}

// Get retrieves a value from cache
func (c *CacheWrapper) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set stores a value in cache with expiration
func (c *CacheWrapper) Set(key string, value interface{}, duration time.Duration) {
	c.cache.Set(key, value, duration)
}

// Delete removes a value from cache
func (c *CacheWrapper) Delete(key string) {
	c.cache.Delete(key)
}

// Flush clears all cache entries
func (c *CacheWrapper) Flush() {
	c.cache.Flush()
}

// Initialize creates the cache instance with configuration
func Initialize() {
	cfg := config.Get()
	Cache = &CacheWrapper{
		cache: cache.New(cfg.Cache.DefaultExpiration, cfg.Cache.CleanupInterval),
	}
}

// GetCache returns the global cache instance
func GetCache() cacheInterface.ICache {
	if Cache == nil {
		Initialize()
	}
	return Cache
}

// Cache TTL constants
func GetTTLStats() time.Duration {
	return config.Get().Cache.StatsExpiration
}

func GetTTLSystemConfig() time.Duration {
	return 1 * time.Hour
}

func GetTTLEmailConfig() time.Duration {
	return 1 * time.Hour
}

func GetTTLPassword() time.Duration {
	return 30 * time.Second
}

// Legacy constants for backward compatibility
const (
	KeyStats        = "stats"
	KeySystemConfig = "system_config"
	KeyEmailConfig  = "email_config"
)

var (
	TTLStats        = GetTTLStats()
	TTLSystemConfig = GetTTLSystemConfig()
	TTLEmailConfig  = GetTTLEmailConfig()
	TTLPassword     = GetTTLPassword()
)

// InvalidateStats clears the stats cache
func InvalidateStats() {
	if Cache != nil {
		Cache.Delete(KeyStats)
	}
}

// InvalidateSystemConfig clears the system config cache
func InvalidateSystemConfig() {
	if Cache != nil {
		Cache.Delete(KeySystemConfig)
	}
}

// InvalidateEmailConfig clears the email config cache
func InvalidateEmailConfig() {
	if Cache != nil {
		Cache.Delete(KeyEmailConfig)
	}
}

// GetPasswordKey returns the cache key for a decrypted password
func GetPasswordKey(accountID uint) string {
	return "password_" + string(rune(accountID))
}
