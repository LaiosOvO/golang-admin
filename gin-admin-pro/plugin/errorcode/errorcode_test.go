package errorcode

import (
	"fmt"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
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

	if config.CachePrefix != "errorcode:" {
		t.Errorf("Default cache prefix should be 'errorcode:', got '%s'", config.CachePrefix)
	}

	if !config.AutoLoadPredefined {
		t.Error("Auto load predefined should be enabled by default")
	}

	if config.DefaultPageSize != 20 {
		t.Errorf("Default page size should be 20, got %d", config.DefaultPageSize)
	}

	if config.MaxPageSize != 100 {
		t.Errorf("Max page size should be 100, got %d", config.MaxPageSize)
	}

	if !config.EnableValidation {
		t.Error("Validation should be enabled by default")
	}
}

func TestPlugin(t *testing.T) {
	// 测试默认配置插件
	plugin := NewPlugin(nil, nil)

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled by default")
	}

	config := plugin.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	// 测试自定义配置插件
	customConfig := &Config{
		Enabled:            false,
		EnableCache:        false,
		CacheExpire:        7200,
		CachePrefix:        "test:",
		AutoLoadPredefined: false,
		DefaultPageSize:    10,
		MaxPageSize:        50,
		EnableValidation:   false,
	}

	plugin2 := NewPlugin(nil, customConfig)

	if plugin2.IsEnabled() {
		t.Error("Plugin2 should be disabled")
	}

	config2 := plugin2.GetConfig()
	if config2.Enabled {
		t.Error("Config should be disabled")
	}

	if config2.CachePrefix != "test:" {
		t.Errorf("Cache prefix should be 'test:', got '%s'", config2.CachePrefix)
	}
}

func TestTableName(t *testing.T) {
	errorCode := ErrorCode{}
	if errorCode.TableName() != "system_error_code" {
		t.Errorf("ErrorCode table name should be 'system_error_code', got '%s'", errorCode.TableName())
	}
}

func TestGetPredefinedErrorCodes(t *testing.T) {
	codes := GetPredefinedErrorCodes()

	if len(codes) == 0 {
		t.Error("Should have predefined error codes")
	}

	// 检查是否包含常见的错误码
	codeMap := make(map[int]ErrorCode)
	for _, code := range codes {
		codeMap[code.Code] = code
	}

	if _, exists := codeMap[CodeSuccess]; !exists {
		t.Error("Should contain success code")
	}

	if _, exists := codeMap[CodeUnknownError]; !exists {
		t.Error("Should contain unknown error code")
	}

	if _, exists := codeMap[CodeUserNotFound]; !exists {
		t.Error("Should contain user not found code")
	}

	if _, exists := codeMap[CodePermissionDenied]; !exists {
		t.Error("Should contain permission denied code")
	}
}

func TestGetErrorTypeName(t *testing.T) {
	tests := []struct {
		errorType string
		expected  string
	}{
		{"common", "通用错误"},
		{"user", "用户错误"},
		{"role", "角色错误"},
		{"permission", "权限错误"},
		{"business", "业务错误"},
		{"system", "系统错误"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetErrorTypeName(tt.errorType)
			if result != tt.expected {
				t.Errorf("GetErrorTypeName(%s) = %v, want %v", tt.errorType, result, tt.expected)
			}
		})
	}
}

func TestGetErrorTypeOptions(t *testing.T) {
	options := GetErrorTypeOptions()

	if len(options) != 6 {
		t.Errorf("GetErrorTypeOptions() returned %d options, want 6", len(options))
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

func TestService(t *testing.T) {
	// 测试服务创建（不使用真实数据库）
	service := NewService(nil, DefaultConfig())
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

func TestGetDefaultValue(t *testing.T) {
	service := NewService(nil, DefaultConfig())

	tests := []struct {
		code     int
		expected string
	}{
		{200, "操作成功"},
		{201, "操作成功"},
		{400, "请求参数错误"},
		{401, "请求参数错误"},
		{500, "服务器内部错误"},
		{501, "服务器内部错误"},
		{999, "未知错误 999"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("code_%d", tt.code), func(t *testing.T) {
			result := service.getDefaultErrorMessage(tt.code)
			if result != tt.expected {
				t.Errorf("getDefaultErrorMessage(%d) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGetErrorMessageWithParams(t *testing.T) {
	service := NewService(nil, DefaultConfig())

	// 测试没有参数的情况
	message := service.GetErrorMessageWithParams(1001)
	expected := "未知错误，请联系管理员"
	if message != expected {
		t.Errorf("GetErrorMessageWithParams() without params = %v, want %v", message, expected)
	}

	// 测试有参数的情况（使用默认消息）
	message = service.GetErrorMessageWithParams(1002, "用户名")
	expected = "请求参数错误：用户名"
	if message != expected {
		t.Errorf("GetErrorMessageWithParams() with params = %v, want %v", message, expected)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	// 测试常量值是否正确
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{"CodeSuccess", CodeSuccess, 0},
		{"CodeUnknownError", CodeUnknownError, 1001},
		{"CodeParamError", CodeParamError, 1002},
		{"CodeDataNotFound", CodeDataNotFound, 1003},
		{"CodeDataExists", CodeDataExists, 1004},
		{"CodeOperationFailed", CodeOperationFailed, 1005},
		{"CodePermissionDenied", CodePermissionDenied, 1006},
		{"CodeTokenExpired", CodeTokenExpired, 1007},
		{"CodeTokenInvalid", CodeTokenInvalid, 1008},
		{"CodeRateLimitExceeded", CodeRateLimitExceeded, 1009},
		{"CodeUserNotFound", CodeUserNotFound, 2001},
		{"CodeUserExists", CodeUserExists, 2002},
		{"CodeUserDisabled", CodeUserDisabled, 2003},
		{"CodePasswordError", CodePasswordError, 2004},
		{"CodeAccountLocked", CodeAccountLocked, 2005},
		{"CodeLoginRequired", CodeLoginRequired, 2006},
		{"CodeRoleNotFound", CodeRoleNotFound, 3001},
		{"CodeRoleExists", CodeRoleExists, 3002},
		{"CodeRoleInUse", CodeRoleInUse, 3003},
		{"CodePermissionNotFound", CodePermissionNotFound, 4001},
		{"CodePermissionExists", CodePermissionExists, 4002},
		{"CodeBusinessError", CodeBusinessError, 5001},
		{"CodeDataInvalid", CodeDataInvalid, 5002},
		{"CodeConfigError", CodeConfigError, 5003},
		{"CodeServiceUnavailable", CodeServiceUnavailable, 5004},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
