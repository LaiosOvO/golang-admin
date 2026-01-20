package operlog

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if !config.EnableMask {
		t.Error("Data masking should be enabled by default")
	}

	if len(config.SensitiveFields) == 0 {
		t.Error("Should have default sensitive fields")
	}

	if config.MaxParamLength != 2000 {
		t.Errorf("Default max param length should be 2000, got %d", config.MaxParamLength)
	}

	if config.MaxResultLength != 2000 {
		t.Errorf("Default max result length should be 2000, got %d", config.MaxResultLength)
	}

	if config.RecordGet {
		t.Error("GET requests should not be recorded by default")
	}

	if !config.RecordPost {
		t.Error("POST requests should be recorded by default")
	}

	if !config.RecordPut {
		t.Error("PUT requests should be recorded by default")
	}

	if !config.RecordDelete {
		t.Error("DELETE requests should be recorded by default")
	}

	if config.RetentionDays != 90 {
		t.Errorf("Default retention days should be 90, got %d", config.RetentionDays)
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
		Enabled:      false,
		EnableMask:   false,
		RecordGet:    true,
		RecordPost:   false,
		RecordPut:    false,
		RecordDelete: false,
	}

	plugin2 := NewPlugin(nil, customConfig)

	if plugin2.IsEnabled() {
		t.Error("Plugin2 should be disabled")
	}

	config2 := plugin2.GetConfig()
	if config2.Enabled {
		t.Error("Config should be disabled")
	}

	// 测试是否应该记录方法
	if !plugin2.ShouldRecordMethod("GET") {
		t.Error("Should record GET when enabled")
	}

	if plugin2.ShouldRecordMethod("POST") {
		t.Error("Should not record POST when disabled")
	}
}

func TestTableName(t *testing.T) {
	operLog := OperLog{}
	if operLog.TableName() != "system_oper_log" {
		t.Errorf("OperLog table name should be 'system_oper_log', got '%s'", operLog.TableName())
	}
}

func TestBusinessTypeFunctions(t *testing.T) {
	tests := []struct {
		businessType int
		expected     string
	}{
		{0, "其它"},
		{1, "新增"},
		{2, "修改"},
		{3, "删除"},
		{4, "授权"},
		{5, "导出"},
		{6, "导入"},
		{999, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetBusinessTypeName(tt.businessType)
			if result != tt.expected {
				t.Errorf("GetBusinessTypeName(%d) = %v, want %v", tt.businessType, result, tt.expected)
			}
		})
	}
}

func TestBusinessTypeOptions(t *testing.T) {
	options := GetBusinessTypeOptions()

	if len(options) != 7 {
		t.Errorf("GetBusinessTypeOptions() returned %d options, want 7", len(options))
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

func TestOperatorTypeFunctions(t *testing.T) {
	tests := []struct {
		operatorType int
		expected     string
	}{
		{0, "其它"},
		{1, "后台用户"},
		{2, "手机端用户"},
		{999, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetOperatorTypeName(tt.operatorType)
			if result != tt.expected {
				t.Errorf("GetOperatorTypeName(%d) = %v, want %v", tt.operatorType, result, tt.expected)
			}
		})
	}
}

func TestOperatorTypeOptions(t *testing.T) {
	options := GetOperatorTypeOptions()

	if len(options) != 3 {
		t.Errorf("GetOperatorTypeOptions() returned %d options, want 3", len(options))
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

func TestStatusFunctions(t *testing.T) {
	tests := []struct {
		status   int
		expected string
	}{
		{0, "正常"},
		{1, "异常"},
		{999, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetStatusName(tt.status)
			if result != tt.expected {
				t.Errorf("GetStatusName(%d) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestStatusOptions(t *testing.T) {
	options := GetStatusOptions()

	if len(options) != 2 {
		t.Errorf("GetStatusOptions() returned %d options, want 2", len(options))
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

func TestGetModuleTitle(t *testing.T) {
	service := NewService(nil, DefaultConfig())

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/system/user/page", "用户管理"},
		{"/api/v1/system/role/list", "角色管理"},
		{"/api/v1/system/menu/tree", "菜单管理"},
		{"/api/v1/system/dept/tree", "部门管理"},
		{"/api/v1/system/dict/type/list", "字典管理"},
		{"/api/v1/infra/file/upload", "文件管理"},
		{"/api/v1/unknown/path", "系统操作"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := service.getModuleTitle(tt.path)
			if result != tt.expected {
				t.Errorf("getModuleTitle(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetBusinessType(t *testing.T) {
	service := NewService(nil, DefaultConfig())

	tests := []struct {
		method   string
		path     string
		expected int
	}{
		{"POST", "/api/v1/system/user/create", 1},   // 新增
		{"PUT", "/api/v1/system/user/update", 2},    // 修改
		{"DELETE", "/api/v1/system/user/delete", 3}, // 删除
		{"GET", "/api/v1/system/user/page", 0},      // 其它
		{"POST", "/api/v1/auth/login", 0},           // 其它
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			result := service.getBusinessType(tt.path, tt.method)
			if result != tt.expected {
				t.Errorf("getBusinessType(%s, %s) = %v, want %v", tt.path, tt.method, result, tt.expected)
			}
		})
	}
}

func TestMaskSensitiveData(t *testing.T) {
	service := NewService(nil, &Config{
		SensitiveFields: []string{"password", "token", "email"},
	})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Password masking",
			input:    `{"username":"admin","password":"123456"}`,
			expected: `{"username":"admin","password":"***"}`,
		},
		{
			name:     "Token masking",
			input:    `{"token":"abc123xyz"}`,
			expected: `{"token":"***"}`,
		},
		{
			name:     "Email masking",
			input:    `{"email":"user@example.com"}`,
			expected: `{"email":"***"}`,
		},
		{
			name:     "No sensitive data",
			input:    `{"username":"admin","role":"user"}`,
			expected: `{"username":"admin","role":"user"}`,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.maskSensitiveData(tt.input)
			if result != tt.expected {
				t.Errorf("maskSensitiveData() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateLogFromContext(t *testing.T) {
	service := NewService(nil, DefaultConfig())

	now := time.Now()
	ctx := &LogContext{
		RequestID:   "test-123",
		UserID:      1,
		Username:    "testuser",
		DeptID:      1,
		DeptName:    "技术部",
		Method:      "POST",
		Path:        "/api/v1/system/user/create",
		IP:          "127.0.0.1",
		UserAgent:   "Mozilla/5.0",
		RequestBody: map[string]interface{}{"username": "test", "password": "123456"},
		StartTime:   now,
		EndTime:     now.Add(100 * time.Millisecond),
		Status:      0,
		Response:    map[string]interface{}{"code": 0, "msg": "success"},
	}

	// This test will not actually save to database since db is nil
	// but it will test the logic that doesn't require database
	err := service.CreateLogFromContext(ctx)
	// Should not return error even when db is nil due to disabled check
	if err != nil {
		t.Errorf("CreateLogFromContext() should not error with nil db, got: %v", err)
	}
}
