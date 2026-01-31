package config

import "time"

// GetDefaults returns the default configuration
func GetDefaults() *Config {
	return &Config{
		App: AppConfig{
			Name:    "Account Manager",
			Version: "2.0.0",
			Debug:   false,
		},
		Database: DatabaseConfig{
			Path:            "data/account_manager.db",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Cache: CacheConfig{
			DefaultExpiration: 5 * time.Minute,
			CleanupInterval:   10 * time.Minute,
			StatsExpiration:   5 * time.Minute,
		},
		Worker: WorkerConfig{
			PoolSize:          10,
			DecryptionWorkers: 10,
			EmailQueueSize:    100,
			EmailWorkers:      3,
		},
		Server: ServerConfig{
			DefaultPort:   22,
			SSHTimeout:    10,
			DeployTimeout: 300,
			BuildTarget:   "linux/amd64",
		},
	}
}
