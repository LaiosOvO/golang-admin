package mysql

import (
	"time"
)

// Config MySQL配置结构
type Config struct {
	Host          string        `yaml:"host" mapstructure:"host"`
	Port          int           `yaml:"port" mapstructure:"port"`
	Database      string        `yaml:"database" mapstructure:"database"`
	Username      string        `yaml:"username" mapstructure:"username"`
	Password      string        `yaml:"password" mapstructure:"password"`
	Charset       string        `yaml:"charset" mapstructure:"charset"`
	ParseTime     bool          `yaml:"parseTime" mapstructure:"parseTime"`
	Loc           string        `yaml:"loc" mapstructure:"loc"`
	MaxIdleConns  int           `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxOpenConns  int           `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxLifetime   time.Duration `yaml:"maxLifetime" mapstructure:"maxLifetime"`
	LogLevel      string        `yaml:"logLevel" mapstructure:"logLevel"`
	SlowThreshold time.Duration `yaml:"slowThreshold" mapstructure:"slowThreshold"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:          "localhost",
		Port:          3306,
		Database:      "gin_admin",
		Username:      "root",
		Password:      "password",
		Charset:       "utf8mb4",
		ParseTime:     true,
		Loc:           "Local",
		MaxIdleConns:  10,
		MaxOpenConns:  100,
		MaxLifetime:   time.Hour,
		LogLevel:      "info",
		SlowThreshold: 200 * time.Millisecond,
	}
}
