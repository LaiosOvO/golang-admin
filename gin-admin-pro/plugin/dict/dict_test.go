package dict

import (
	"testing"
)

func TestGetDataScopeName(t *testing.T) {
	// 测试字典插件的默认配置
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if !config.EnableCache {
		t.Error("Cache should be enabled by default")
	}

	if config.CacheExpire != 3600 {
		t.Errorf("Default cache expire should be 3600, got %d", config.CacheExpire)
	}

	if config.CachePrefix != "dict:" {
		t.Errorf("Default cache prefix should be 'dict:', got '%s'", config.CachePrefix)
	}

	if config.AutoRefresh {
		t.Error("Auto refresh should be disabled by default")
	}

	if config.RefreshInterval != 300 {
		t.Errorf("Default refresh interval should be 300, got %d", config.RefreshInterval)
	}

	if config.DefaultPageSize != 20 {
		t.Errorf("Default page size should be 20, got %d", config.DefaultPageSize)
	}

	if config.MaxPageSize != 100 {
		t.Errorf("Max page size should be 100, got %d", config.MaxPageSize)
	}
}

func TestPlugin(t *testing.T) {
	// 测试默认配置插件
	plugin := NewPlugin(nil, nil)

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled by default")
	}

	if plugin.GetDB() != nil {
		t.Error("DB should be nil when not provided")
	}

	config := plugin.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	// 测试自定义配置插件
	customConfig := &Config{
		Enabled:         false,
		EnableCache:     false,
		CacheExpire:     7200,
		CachePrefix:     "test:",
		AutoRefresh:     true,
		RefreshInterval: 600,
		DefaultPageSize: 10,
		MaxPageSize:     50,
	}

	plugin2 := NewPlugin(nil, customConfig)

	if plugin2.IsEnabled() {
		t.Error("Plugin2 should be disabled")
	}

	config2 := plugin2.GetConfig()
	if !config2.AutoRefresh {
		t.Error("Auto refresh should be enabled")
	}

	if config2.CachePrefix != "test:" {
		t.Errorf("Cache prefix should be 'test:', got '%s'", config2.CachePrefix)
	}
}

func TestTableName(t *testing.T) {
	dictType := DictType{}
	if dictType.TableName() != "system_dict_type" {
		t.Errorf("DictType table name should be 'system_dict_type', got '%s'", dictType.TableName())
	}

	dictData := DictData{}
	if dictData.TableName() != "system_dict_data" {
		t.Errorf("DictData table name should be 'system_dict_data', got '%s'", dictData.TableName())
	}
}

func TestService(t *testing.T) {
	// 测试服务创建（不使用真实数据库）
	service := NewService(nil)
	if service == nil {
		t.Error("Service should not be nil")
	}

	// 测试插件服务获取
	plugin := NewPlugin(nil, nil)
	service2 := plugin.GetService()
	if service2 == nil {
		t.Error("Plugin service should not be nil when enabled")
	}

	// 测试禁用状态的服务获取
	disabledPlugin := NewPlugin(nil, &Config{Enabled: false})
	service3 := disabledPlugin.GetService()
	if service3 != nil {
		t.Error("Disabled plugin service should be nil")
	}
}
