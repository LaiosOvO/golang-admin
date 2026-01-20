package errorcode

import (
	"gorm.io/gorm"
)

// Config 错误码配置
type Config struct {
	// 是否启用错误码管理
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 是否启用缓存
	EnableCache bool `yaml:"enableCache" json:"enableCache"`
	// 缓存过期时间（秒）
	CacheExpire int `yaml:"cacheExpire" json:"cacheExpire"`
	// 缓存前缀
	CachePrefix string `yaml:"cachePrefix" json:"cachePrefix"`
	// 是否自动加载预定义错误码
	AutoLoadPredefined bool `yaml:"autoLoadPredefined" json:"autoLoadPredefined"`
	// 默认分页大小
	DefaultPageSize int `yaml:"defaultPageSize" json:"defaultPageSize"`
	// 最大分页大小
	MaxPageSize int `yaml:"maxPageSize" json:"maxPageSize"`
	// 是否启用错误码验证
	EnableValidation bool `yaml:"enableValidation" json:"enableValidation"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:            true,
		EnableCache:        true,
		CacheExpire:        3600, // 1小时
		CachePrefix:        "errorcode:",
		AutoLoadPredefined: true,
		DefaultPageSize:    20,
		MaxPageSize:        100,
		EnableValidation:   true,
	}
}

// Plugin 错误码插件
type Plugin struct {
	config *Config
	db     *gorm.DB
}

// NewPlugin 创建错误码插件
func NewPlugin(db *gorm.DB, cfg *Config) *Plugin {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Plugin{
		config: cfg,
		db:     db,
	}
}

// GetConfig 获取配置
func (p *Plugin) GetConfig() *Config {
	return p.config
}

// GetDB 获取数据库连接
func (p *Plugin) GetDB() *gorm.DB {
	return p.db
}

// IsEnabled 是否启用
func (p *Plugin) IsEnabled() bool {
	return p.config.Enabled
}

// GetService 获取错误码服务
func (p *Plugin) GetService() *Service {
	if !p.IsEnabled() {
		return nil
	}
	return NewService(p.db, p.config)
}

// AutoMigrate 自动迁移数据库表
func (p *Plugin) AutoMigrate() error {
	if !p.IsEnabled() {
		return nil
	}

	return p.db.AutoMigrate(&ErrorCode{})
}

// Init 初始化插件
func (p *Plugin) Init() error {
	if !p.IsEnabled() {
		return nil
	}

	// 自动迁移数据库表
	if err := p.AutoMigrate(); err != nil {
		return err
	}

	// 自动加载预定义错误码
	if p.config.AutoLoadPredefined {
		if err := p.loadPredefinedErrorCodes(); err != nil {
			return err
		}
	}

	return nil
}

// loadPredefinedErrorCodes 加载预定义错误码
func (p *Plugin) loadPredefinedErrorCodes() error {
	service := p.GetService()
	if service == nil {
		return nil
	}

	predefinedCodes := GetPredefinedErrorCodes()
	for _, code := range predefinedCodes {
		// 检查是否已存在
		var existing ErrorCode
		err := p.db.Where("code = ?", code.Code).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 不存在，创建新记录
			if err := p.db.Create(&code).Error; err != nil {
				return err
			}
		} else if err == nil {
			// 已存在，更新记录（只更新消息和解决方案）
			existing.Message = code.Message
			existing.Solution = code.Solution
			existing.Status = code.Status
			if err := p.db.Save(&existing).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
