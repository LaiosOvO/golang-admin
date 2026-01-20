package dataperm

import (
	"testing"
)

func TestGetDataScopeName(t *testing.T) {
	tests := []struct {
		dataScope int
		expected  string
	}{
		{DataScopeAll, "全部数据权限"},
		{DataScopeCustom, "自定义数据权限"},
		{DataScopeDept, "本部门数据权限"},
		{DataScopeDeptChild, "本部门及以下数据权限"},
		{DataScopeSelf, "仅本人数据权限"},
		{999, "未知权限"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetDataScopeName(tt.dataScope)
			if result != tt.expected {
				t.Errorf("GetDataScopeName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDataScopeOptions(t *testing.T) {
	options := GetDataScopeOptions()

	if len(options) != 5 {
		t.Errorf("GetDataScopeOptions() returned %d options, want 5", len(options))
	}

	for i, option := range options {
		if _, ok := option["value"]; !ok {
			t.Errorf("Option %d missing 'value' field", i)
		}
		if _, ok := option["label"]; !ok {
			t.Errorf("Option %d missing 'label' field", i)
		}
	}
}

func TestPlugin(t *testing.T) {
	// 测试默认配置
	plugin := NewPlugin(nil, nil)
	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled by default")
	}

	if plugin.GetConfig().DefaultDataScope != DataScopeDept {
		t.Errorf("Default data scope should be %d, got %d", DataScopeDept, plugin.GetConfig().DefaultDataScope)
	}

	// 测试自定义配置
	customConfig := &Config{
		Enabled:          false,
		DefaultDataScope: DataScopeSelf,
		EnableCache:      false,
		CacheExpire:      600,
		EnableLog:        true,
	}

	plugin2 := NewPlugin(nil, customConfig)
	if plugin2.IsEnabled() {
		t.Error("Plugin2 should be disabled")
	}

	if plugin2.GetConfig().DefaultDataScope != DataScopeSelf {
		t.Errorf("Default data scope should be %d, got %d", DataScopeSelf, plugin2.GetConfig().DefaultDataScope)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if config.DefaultDataScope != DataScopeDept {
		t.Errorf("Default data scope should be %d, got %d", DataScopeDept, config.DefaultDataScope)
	}

	if !config.EnableCache {
		t.Error("Cache should be enabled by default")
	}

	if config.CacheExpire != 300 {
		t.Errorf("Default cache expire should be 300, got %d", config.CacheExpire)
	}

	if config.EnableLog {
		t.Error("Log should be disabled by default")
	}
}
