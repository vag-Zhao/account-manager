package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Loader handles configuration loading
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// LoadOrCreate loads configuration from file or creates default config file
func (l *Loader) LoadOrCreate() (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(l.configPath); os.IsNotExist(err) {
		// Create default config file
		if err := l.CreateDefaultConfig(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
	}

	// Load configuration
	return Load(l.configPath)
}

// CreateDefaultConfig creates a default configuration file
func (l *Loader) CreateDefaultConfig() error {
	defaults := GetDefaults()

	data, err := yaml.Marshal(defaults)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := "."
	if l.configPath != "config.yaml" {
		dir = l.configPath[:len(l.configPath)-len("config.yaml")]
	}
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(l.configPath, data, 0644)
}

// Reload reloads configuration from file
func (l *Loader) Reload() (*Config, error) {
	return Load(l.configPath)
}

// Validate validates the configuration
func Validate(cfg *Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name cannot be empty")
	}
	if cfg.Database.Path == "" {
		return fmt.Errorf("database.path cannot be empty")
	}
	if cfg.Worker.PoolSize < 1 {
		return fmt.Errorf("worker.pool_size must be at least 1")
	}
	return nil
}
