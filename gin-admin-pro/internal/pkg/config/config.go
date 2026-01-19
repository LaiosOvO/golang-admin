package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var GlobalConfig *Config

// Load 加载配置文件
func Load(configPath string) error {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configPath == "" {
		// 默认配置文件路径
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../config")
		viper.AddConfigPath("../../config")
	} else {
		viper.AddConfigPath(configPath)
	}

	// 支持环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("配置文件未找到: %w", err)
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// LoadWithEnv 根据环境加载配置
func LoadWithEnv(env string) error {
	configDir := "./config"

	// 首先加载默认配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// 读取默认配置
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取默认配置文件失败: %w", err)
	}

	// 环境特定配置文件
	envConfigFile := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))

	// 检查环境配置文件是否存在
	if _, err := os.Stat(envConfigFile); err == nil {
		// 设置环境配置文件并合并
		viper.SetConfigFile(envConfigFile)
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("合并环境配置文件失败: %w", err)
		}
	}

	// 支持环境变量覆盖
	viper.AutomaticEnv()

	// 解析配置
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// validateConfig 验证配置
func validateConfig() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置为空")
	}

	// 验证服务器配置
	if GlobalConfig.Server.Port <= 0 || GlobalConfig.Server.Port > 65535 {
		return fmt.Errorf("服务器端口配置无效: %d", GlobalConfig.Server.Port)
	}

	// 验证数据库配置
	if err := validateDatabaseConfig(); err != nil {
		return fmt.Errorf("数据库配置验证失败: %w", err)
	}

	// 验证JWT配置
	if GlobalConfig.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// validateDatabaseConfig 验证数据库配置
func validateDatabaseConfig() error {
	// MySQL配置验证
	if GlobalConfig.Database.MySQL.Host == "" {
		return fmt.Errorf("MySQL主机地址不能为空")
	}
	if GlobalConfig.Database.MySQL.Port <= 0 || GlobalConfig.Database.MySQL.Port > 65535 {
		return fmt.Errorf("MySQL端口配置无效: %d", GlobalConfig.Database.MySQL.Port)
	}
	if GlobalConfig.Database.MySQL.Database == "" {
		return fmt.Errorf("MySQL数据库名不能为空")
	}

	// PostgreSQL配置验证
	if GlobalConfig.Database.PostgreSQL.Host == "" {
		return fmt.Errorf("PostgreSQL主机地址不能为空")
	}
	if GlobalConfig.Database.PostgreSQL.Port <= 0 || GlobalConfig.Database.PostgreSQL.Port > 65535 {
		return fmt.Errorf("PostgreSQL端口配置无效: %d", GlobalConfig.Database.PostgreSQL.Port)
	}

	// Redis配置验证
	if GlobalConfig.Database.Redis.Host == "" {
		return fmt.Errorf("Redis主机地址不能为空")
	}
	if GlobalConfig.Database.Redis.Port <= 0 || GlobalConfig.Database.Redis.Port > 65535 {
		return fmt.Errorf("Redis端口配置无效: %d", GlobalConfig.Database.Redis.Port)
	}

	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if GlobalConfig == nil {
		panic("配置未初始化，请先调用 Load() 函数")
	}
	return GlobalConfig
}

// IsDebug 是否为调试模式
func IsDebug() bool {
	return GetConfig().Server.Mode == "debug"
}

// IsRelease 是否为生产模式
func IsRelease() bool {
	return GetConfig().Server.Mode == "release"
}

// IsTest 是否为测试模式
func IsTest() bool {
	return GetConfig().Server.Mode == "test"
}

// WatchConfig 监听配置文件变化（热加载）
func WatchConfig(callback func(*Config)) error {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已更改: %s\n", e.Name)

		// 重新解析配置
		newConfig := &Config{}
		if err := viper.Unmarshal(newConfig); err != nil {
			fmt.Printf("解析新配置失败: %v\n", err)
			return
		}

		// 验证新配置
		oldConfig := GlobalConfig
		GlobalConfig = newConfig
		if err := validateConfig(); err != nil {
			fmt.Printf("新配置验证失败: %v\n", err)
			GlobalConfig = oldConfig // 恢复原配置
			return
		}

		// 调用回调函数
		if callback != nil {
			callback(GlobalConfig)
		}
	})

	return nil
}
