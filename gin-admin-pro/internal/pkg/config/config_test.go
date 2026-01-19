package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						Database: "test",
					},
					PostgreSQL: PostgreSQLConfig{
						Host: "localhost",
						Port: 5432,
					},
					Redis: RedisConfig{
						Host: "localhost",
						Port: 6379,
					},
				},
				JWT: JWTConfig{Secret: "test-secret"},
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "invalid port",
			config: &Config{
				Server: ServerConfig{Port: -1},
			},
			wantErr: true,
		},
		{
			name: "empty jwt secret",
			config: &Config{
				Server: ServerConfig{Port: 8080},
				JWT:    JWTConfig{Secret: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GlobalConfig = tt.config
			err := validateConfig()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsDebug(t *testing.T) {
	GlobalConfig = &Config{Server: ServerConfig{Mode: "debug"}}
	assert.True(t, IsDebug())
	assert.False(t, IsRelease())
	assert.False(t, IsTest())
}

func TestIsRelease(t *testing.T) {
	GlobalConfig = &Config{Server: ServerConfig{Mode: "release"}}
	assert.False(t, IsDebug())
	assert.True(t, IsRelease())
	assert.False(t, IsTest())
}

func TestIsTest(t *testing.T) {
	GlobalConfig = &Config{Server: ServerConfig{Mode: "test"}}
	assert.False(t, IsDebug())
	assert.False(t, IsRelease())
	assert.True(t, IsTest())
}

func TestGetConfig(t *testing.T) {
	// 测试未初始化的情况
	GlobalConfig = nil
	assert.Panics(t, func() {
		GetConfig()
	})

	// 测试正常情况
	GlobalConfig = &Config{Server: ServerConfig{Port: 8080}}
	config := GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestLoadNonExistentFile(t *testing.T) {
	err := Load("/non/existent/path")
	assert.Error(t, err)
}

// 集成测试 - 测试实际配置文件加载（从项目根目录运行）
func TestIntegrationLoad(t *testing.T) {
	// 检查是否在正确的目录中
	if _, err := os.Stat("../../config/config.yaml"); os.IsNotExist(err) {
		t.Skip("配置文件不存在，跳过集成测试")
		return
	}

	// 测试加载默认配置
	err := Load("../../config")
	require.NoError(t, err)
	assert.NotNil(t, GlobalConfig)
	assert.Equal(t, 8080, GlobalConfig.Server.Port)
	assert.Equal(t, "debug", GlobalConfig.Server.Mode)
}
