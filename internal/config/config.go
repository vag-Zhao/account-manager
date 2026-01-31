package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"cache"`
	Worker   WorkerConfig   `yaml:"worker"`
	Server   ServerConfig   `yaml:"server"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path            string        `yaml:"path"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	DefaultExpiration time.Duration `yaml:"default_expiration"`
	CleanupInterval   time.Duration `yaml:"cleanup_interval"`
	StatsExpiration   time.Duration `yaml:"stats_expiration"`
}

// WorkerConfig holds worker pool configuration
type WorkerConfig struct {
	PoolSize           int `yaml:"pool_size"`
	DecryptionWorkers  int `yaml:"decryption_workers"`
	EmailQueueSize     int `yaml:"email_queue_size"`
	EmailWorkers       int `yaml:"email_workers"`
}

// ServerConfig holds server deployment configuration
type ServerConfig struct {
	DefaultPort     int    `yaml:"default_port"`
	SSHTimeout      int    `yaml:"ssh_timeout"`
	DeployTimeout   int    `yaml:"deploy_timeout"`
	BuildTarget     string `yaml:"build_target"`
}

// Global configuration instance
var globalConfig *Config

// Load loads configuration from a YAML file
func Load(configPath string) (*Config, error) {
	// If config file doesn't exist, use defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaults(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Merge with defaults for any missing values
	defaults := GetDefaults()
	mergeDefaults(&cfg, defaults)

	globalConfig = &cfg
	return &cfg, nil
}

// Get returns the global configuration instance
func Get() *Config {
	if globalConfig == nil {
		globalConfig = GetDefaults()
	}
	return globalConfig
}

// mergeDefaults fills in missing values with defaults
func mergeDefaults(cfg, defaults *Config) {
	if cfg.App.Name == "" {
		cfg.App = defaults.App
	}
	if cfg.Database.Path == "" {
		cfg.Database = defaults.Database
	}
	if cfg.Cache.DefaultExpiration == 0 {
		cfg.Cache = defaults.Cache
	}
	if cfg.Worker.PoolSize == 0 {
		cfg.Worker = defaults.Worker
	}
	if cfg.Server.DefaultPort == 0 {
		cfg.Server = defaults.Server
	}
}
