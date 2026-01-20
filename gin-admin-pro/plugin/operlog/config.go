package operlog

import (
	"gorm.io/gorm"
)

// Config 操作日志配置
type Config struct {
	// 是否启用操作日志
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 是否启用数据脱敏
	EnableMask bool `yaml:"enableMask" json:"enableMask"`
	// 敏感字段列表
	SensitiveFields []string `yaml:"sensitiveFields" json:"sensitiveFields"`
	// 请求参数最大长度
	MaxParamLength int `yaml:"maxParamLength" json:"maxParamLength"`
	// 返回结果最大长度
	MaxResultLength int `yaml:"maxResultLength" json:"maxResultLength"`
	// 是否记录GET请求
	RecordGet bool `yaml:"recordGet" json:"recordGet"`
	// 是否记录POST请求
	RecordPost bool `yaml:"recordPost" json:"recordPost"`
	// 是否记录PUT请求
	RecordPut bool `yaml:"recordPut" json:"recordPut"`
	// 是否记录DELETE请求
	RecordDelete bool `yaml:"recordDelete" json:"recordDelete"`
	// 日志保留天数
	RetentionDays int `yaml:"retentionDays" json:"retentionDays"`
	// 默认分页大小
	DefaultPageSize int `yaml:"defaultPageSize" json:"defaultPageSize"`
	// 最大分页大小
	MaxPageSize int `yaml:"maxPageSize" json:"maxPageSize"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:         true,
		EnableMask:      true,
		SensitiveFields: []string{"password", "pwd", "token", "secret", "key", "access_key", "secret_key"},
		MaxParamLength:  2000,
		MaxResultLength: 2000,
		RecordGet:       false, // 默认不记录GET请求
		RecordPost:      true,
		RecordPut:       true,
		RecordDelete:    true,
		RetentionDays:   90, // 保留90天
		DefaultPageSize: 20,
		MaxPageSize:     100,
	}
}

// Plugin 操作日志插件
type Plugin struct {
	config *Config
	db     *gorm.DB
}

// NewPlugin 创建操作日志插件
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

// GetService 获取操作日志服务
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

	return p.db.AutoMigrate(&OperLog{})
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

	return nil
}

// ShouldRecordMethod 判断是否记录指定方法的日志
func (p *Plugin) ShouldRecordMethod(method string) bool {
	switch method {
	case "GET":
		return p.config.RecordGet
	case "POST":
		return p.config.RecordPost
	case "PUT":
		return p.config.RecordPut
	case "DELETE":
		return p.config.RecordDelete
	default:
		return true // 其他方法默认记录
	}
}
