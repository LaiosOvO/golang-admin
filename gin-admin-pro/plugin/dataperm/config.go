package dataperm

import (
	"gorm.io/gorm"
)

// Config 数据权限配置
type Config struct {
	// 是否启用数据权限
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 默认数据权限范围
	DefaultDataScope int `yaml:"defaultDataScope" json:"defaultDataScope"`
	// 是否启用权限缓存
	EnableCache bool `yaml:"enableCache" json:"enableCache"`
	// 缓存过期时间（秒）
	CacheExpire int `yaml:"cacheExpire" json:"cacheExpire"`
	// 是否启用详细日志
	EnableLog bool `yaml:"enableLog" json:"enableLog"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:          true,
		DefaultDataScope: DataScopeDept, // 默认本部门数据权限
		EnableCache:      true,
		CacheExpire:      300, // 5分钟
		EnableLog:        false,
	}
}

// Plugin 数据权限插件
type Plugin struct {
	config *Config
	db     *gorm.DB
}

// NewPlugin 创建数据权限插件
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

// GetService 获取数据权限服务
func (p *Plugin) GetService() *Service {
	if !p.IsEnabled() {
		return nil
	}
	return NewService(p.db)
}

// AutoMigrate 自动迁移数据库表
func (p *Plugin) AutoMigrate() error {
	if !p.IsEnabled() {
		return nil
	}

	return p.db.AutoMigrate(&RoleDept{})
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

	// 可以在这里添加其他初始化逻辑
	// 比如初始化缓存、创建索引等

	return nil
}
