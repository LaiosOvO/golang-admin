package postgresql

import (
	"time"
)

// Config PostgreSQL配置结构
type Config struct {
	Host          string        `yaml:"host" mapstructure:"host"`
	Port          int           `yaml:"port" mapstructure:"port"`
	Database      string        `yaml:"database" mapstructure:"database"`
	Username      string        `yaml:"username" mapstructure:"username"`
	Password      string        `yaml:"password" mapstructure:"password"`
	SSLMode       string        `yaml:"sslMode" mapstructure:"sslMode"`
	Timezone      string        `yaml:"timezone" mapstructure:"timezone"`
	MaxIdleConns  int           `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxOpenConns  int           `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxLifetime   time.Duration `yaml:"maxLifetime" mapstructure:"maxLifetime"`
	LogLevel      string        `yaml:"logLevel" mapstructure:"logLevel"`
	SlowThreshold time.Duration `yaml:"slowThreshold" mapstructure:"slowThreshold"`

	// 扩展配置
	Extensions []ExtensionConfig `yaml:"extensions" mapstructure:"extensions"`
}

// ExtensionConfig 扩展配置
type ExtensionConfig struct {
	Name    string `yaml:"name" mapstructure:"name"`
	Version string `yaml:"version" mapstructure:"version"`
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:          "localhost",
		Port:          5432,
		Database:      "gin_admin",
		Username:      "postgres",
		Password:      "password",
		SSLMode:       "disable",
		Timezone:      "Asia/Shanghai",
		MaxIdleConns:  10,
		MaxOpenConns:  100,
		MaxLifetime:   time.Hour,
		LogLevel:      "info",
		SlowThreshold: 200 * time.Millisecond,
		Extensions: []ExtensionConfig{
			{Name: "postgis", Version: "3.3", Enabled: true},
			{Name: "vector", Version: "0.5.1", Enabled: true},
			{Name: "uuid-ossp", Version: "1.1", Enabled: true},
			{Name: "btree_gin", Version: "1.0", Enabled: true},
		},
	}
}
