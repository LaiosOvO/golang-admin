package cron

import (
	"time"
)

// Config Cron配置结构
type Config struct {
	// 基础配置
	Enabled     bool   `yaml:"enabled" mapstructure:"enabled"`
	Timezone    string `yaml:"timezone" mapstructure:"timezone"`
	Concurrency int    `yaml:"concurrency" mapstructure:"concurrency"`

	// 任务配置
	Jobs []*JobConfig `yaml:"jobs" mapstructure:"jobs"`

	// 执行配置
	Timeout    time.Duration `yaml:"timeout" mapstructure:"timeout"`
	MaxRetries int           `yaml:"maxRetries" mapstructure:"maxRetries"`
	RetryDelay time.Duration `yaml:"retryDelay" mapstructure:"retryDelay"`

	// 监控配置
	LogEnabled bool   `yaml:"logEnabled" mapstructure:"logEnabled"`
	LogLevel   string `yaml:"logLevel" mapstructure:"logLevel"`

	// 管理配置
	ManagementEnabled bool `yaml:"managementEnabled" mapstructure:"managementEnabled"`
}

// JobConfig 任务配置
type JobConfig struct {
	// 基础信息
	ID          string `yaml:"id" mapstructure:"id"`
	Name        string `yaml:"name" mapstructure:"name"`
	Description string `yaml:"description" mapstructure:"description"`
	Enabled     bool   `yaml:"enabled" mapstructure:"enabled"`

	// 调度配置
	Cron     string `yaml:"cron" mapstructure:"cron"`
	Timezone string `yaml:"timezone" mapstructure:"timezone"`

	// 执行配置
	Handler    string        `yaml:"handler" mapstructure:"handler"`
	Timeout    time.Duration `yaml:"timeout" mapstructure:"timeout"`
	MaxRetries int           `yaml:"maxRetries" mapstructure:"maxRetries"`
	RetryDelay time.Duration `yaml:"retryDelay" mapstructure:"retryDelay"`

	// 参数配置
	Params map[string]interface{} `yaml:"params" mapstructure:"params"`

	// 依赖配置
	DependsOn []string `yaml:"dependsOn" mapstructure:"dependsOn"`

	// 执行限制
	MaxInstances int  `yaml:"maxInstances" mapstructure:"maxInstances"`
	Singleton    bool `yaml:"singleton" mapstructure:"singleton"`

	// 标签和元数据
	Tags     []string          `yaml:"tags" mapstructure:"tags"`
	Metadata map[string]string `yaml:"metadata" mapstructure:"metadata"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		Timezone:          "UTC",
		Concurrency:       100,
		Jobs:              []*JobConfig{},
		Timeout:           time.Minute * 30,
		MaxRetries:        3,
		RetryDelay:        time.Second * 5,
		LogEnabled:        true,
		LogLevel:          "info",
		ManagementEnabled: true,
	}
}

// DefaultJobConfig 返回默认任务配置
func DefaultJobConfig() *JobConfig {
	return &JobConfig{
		Enabled:      true,
		Cron:         "0 */1 * * *", // 每小时执行一次
		Timezone:     "UTC",
		Timeout:      time.Minute * 10,
		MaxRetries:   3,
		RetryDelay:   time.Second * 5,
		MaxInstances: 1,
		Singleton:    true,
		Params:       make(map[string]interface{}),
		DependsOn:    []string{},
		Tags:         []string{},
		Metadata:     make(map[string]string),
	}
}
