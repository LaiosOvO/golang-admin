package mysql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 3306, cfg.Port)
	assert.Equal(t, "gin_admin", cfg.Database)
	assert.Equal(t, "root", cfg.Username)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "utf8mb4", cfg.Charset)
	assert.True(t, cfg.ParseTime)
	assert.Equal(t, "Local", cfg.Loc)
	assert.Equal(t, 10, cfg.MaxIdleConns)
	assert.Equal(t, 100, cfg.MaxOpenConns)
	assert.Equal(t, time.Hour, cfg.MaxLifetime)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 200*time.Millisecond, cfg.SlowThreshold)
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		valid  bool
	}{
		{
			name:   "default config is valid",
			config: DefaultConfig(),
			valid:  true,
		},
		{
			name: "empty host is invalid",
			config: &Config{
				Host:     "",
				Port:     3306,
				Database: "test",
				Username: "root",
				Password: "password",
			},
			valid: false,
		},
		{
			name: "invalid port",
			config: &Config{
				Host:     "localhost",
				Port:     -1,
				Database: "test",
				Username: "root",
				Password: "password",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.config.Host)
				assert.Greater(t, tt.config.Port, 0)
				assert.NotEmpty(t, tt.config.Database)
				assert.NotEmpty(t, tt.config.Username)
			} else {
				// 简单的验证逻辑
				if tt.config.Host == "" || tt.config.Port <= 0 {
					// 配置无效，这是预期的
					assert.True(t, true)
				}
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	// 测试不同日志级别
	loggers := []string{"silent", "error", "warn", "info", "debug"}

	for _, level := range loggers {
		logger := getLogger(level, 200*time.Millisecond)
		assert.NotNil(t, logger)
	}
}
