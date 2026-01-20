package dict

import (
	"gorm.io/gorm"
)

// Config 字典管理配置
type Config struct {
	// 是否启用字典管理
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 是否启用缓存
	EnableCache bool `yaml:"enableCache" json:"enableCache"`
	// 缓存过期时间（秒）
	CacheExpire int `yaml:"cacheExpire" json:"cacheExpire"`
	// 缓存前缀
	CachePrefix string `yaml:"cachePrefix" json:"cachePrefix"`
	// 是否启用自动刷新
	AutoRefresh bool `yaml:"autoRefresh" json:"autoRefresh"`
	// 刷新间隔（秒）
	RefreshInterval int `yaml:"refreshInterval" json:"refreshInterval"`
	// 默认分页大小
	DefaultPageSize int `yaml:"defaultPageSize" json:"defaultPageSize"`
	// 最大分页大小
	MaxPageSize int `yaml:"maxPageSize" json:"maxPageSize"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:         true,
		EnableCache:     true,
		CacheExpire:     3600, // 1小时
		CachePrefix:     "dict:",
		AutoRefresh:     false,
		RefreshInterval: 300, // 5分钟
		DefaultPageSize: 20,
		MaxPageSize:     100,
	}
}

// Plugin 字典插件
type Plugin struct {
	config *Config
	db     *gorm.DB
}

// NewPlugin 创建字典插件
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

// GetService 获取字典服务
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

	return p.db.AutoMigrate(&DictType{}, &DictData{})
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

	// 初始化默认字典数据
	if err := p.initDefaultData(); err != nil {
		return err
	}

	return nil
}

// initDefaultData 初始化默认字典数据
func (p *Plugin) initDefaultData() error {
	service := p.GetService()
	if service == nil {
		return nil
	}

	// 初始化系统状态字典
	statusType := &DictType{
		Name:   "系统状态",
		Type:   "system_status",
		Status: 1,
		Remark: "系统状态字典",
	}

	// 检查是否已存在
	_, err := service.GetDictTypeByType("system_status")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建字典类型
			if err := service.CreateDictType(statusType); err != nil {
				return err
			}

			// 创建字典数据
			statusData := []*DictData{
				{DictSort: 1, Label: "启用", Value: "1", DictType: "system_status", Status: 1},
				{DictSort: 2, Label: "禁用", Value: "0", DictType: "system_status", Status: 1},
			}
			for _, data := range statusData {
				if err := service.CreateDictData(data); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	// 初始化用户性别字典
	genderType := &DictType{
		Name:   "用户性别",
		Type:   "user_gender",
		Status: 1,
		Remark: "用户性别字典",
	}

	_, err = service.GetDictTypeByType("user_gender")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := service.CreateDictType(genderType); err != nil {
				return err
			}

			genderData := []*DictData{
				{DictSort: 1, Label: "男", Value: "1", DictType: "user_gender", Status: 1},
				{DictSort: 2, Label: "女", Value: "2", DictType: "user_gender", Status: 1},
				{DictSort: 3, Label: "保密", Value: "3", DictType: "user_gender", Status: 1},
			}
			for _, data := range genderData {
				if err := service.CreateDictData(data); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	// 初始化数据权限字典
	dataScopeType := &DictType{
		Name:   "数据权限",
		Type:   "data_scope",
		Status: 1,
		Remark: "数据权限范围字典",
	}

	_, err = service.GetDictTypeByType("data_scope")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := service.CreateDictType(dataScopeType); err != nil {
				return err
			}

			dataScopeData := []*DictData{
				{DictSort: 1, Label: "全部数据权限", Value: "1", DictType: "data_scope", Status: 1},
				{DictSort: 2, Label: "自定义数据权限", Value: "2", DictType: "data_scope", Status: 1},
				{DictSort: 3, Label: "本部门数据权限", Value: "3", DictType: "data_scope", Status: 1},
				{DictSort: 4, Label: "本部门及以下数据权限", Value: "4", DictType: "data_scope", Status: 1},
				{DictSort: 5, Label: "仅本人数据权限", Value: "5", DictType: "data_scope", Status: 1},
			}
			for _, data := range dataScopeData {
				if err := service.CreateDictData(data); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	return nil
}
