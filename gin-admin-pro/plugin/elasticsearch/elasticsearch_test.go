package elasticsearch

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if len(config.Addresses) == 0 {
		t.Error("Should have default addresses")
	}

	if config.Addresses[0] != "http://localhost:9200" {
		t.Errorf("Default address should be 'http://localhost:9200', got '%s'", config.Addresses[0])
	}

	if config.Timeout != 30 {
		t.Errorf("Default timeout should be 30, got %d", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Default max retries should be 3, got %d", config.MaxRetries)
	}

	if config.DefaultIndexPrefix != "gin_admin" {
		t.Errorf("Default index prefix should be 'gin_admin', got '%s'", config.DefaultIndexPrefix)
	}
}

func TestPlugin(t *testing.T) {
	// 测试默认配置插件
	plugin := NewPlugin(nil)

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled by default")
	}

	config := plugin.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	client := plugin.GetClient()
	if client == nil {
		t.Error("Client should not be nil")
	}

	// 测试禁用插件
	disabledConfig := &Config{Enabled: false}
	disabledPlugin := NewPlugin(disabledConfig)

	if disabledPlugin.IsEnabled() {
		t.Error("Disabled plugin should not be enabled")
	}
}

func TestClient(t *testing.T) {
	config := DefaultConfig()
	client, err := NewClient(config)
	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Error("Client should not be nil")
	}

	if !client.IsEnabled() {
		t.Error("Client should be enabled")
	}

	// 测试禁用客户端
	disabledConfig := &Config{Enabled: false}
	disabledClient, err := NewClient(disabledConfig)
	if err != nil {
		t.Errorf("Failed to create disabled client: %v", err)
	}

	if disabledClient.IsEnabled() {
		t.Error("Disabled client should not be enabled")
	}
}

func TestGetHealthStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"green", "健康"},
		{"yellow", "警告"},
		{"red", "异常"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetHealthStatus(tt.input)
			if result != tt.expected {
				t.Errorf("GetHealthStatus(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
