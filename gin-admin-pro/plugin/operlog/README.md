# 操作日志插件

操作日志插件提供系统操作的完整记录功能，支持记录用户操作、请求参数、响应结果、异常信息等，并提供数据脱敏、统计分析等功能。

## 功能特性

- 完整记录用户操作行为
- 支持多种请求方式记录
- 数据脱敏保护敏感信息
- 操作统计分析
- 日志导出和清理
- 支持IP地址定位
- 性能监控（响应时间）

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/operlog"

// 创建操作日志插件
operLogPlugin := operlog.NewPlugin(db, nil)

// 初始化插件
err := operLogPlugin.Init()
if err != nil {
    log.Fatal("Failed to init oper log plugin:", err)
}

// 获取操作日志服务
operLogService := operLogPlugin.GetService()
```

### 2. 记录操作日志

```go
// 创建日志上下文
logCtx := &operlog.LogContext{
    RequestID:   generateRequestID(),
    UserID:      userID,
    Username:    username,
    DeptID:      deptID,
    DeptName:    deptName,
    Method:      c.Request.Method,
    Path:        c.Request.URL.Path,
    IP:          c.ClientIP(),
    UserAgent:   c.GetHeader("User-Agent"),
    StartTime:   time.Now(),
    Status:      0, // 0正常 1异常
}

// 在中间件中记录日志
func OperationLogMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // 设置日志上下文
        logCtx := &operlog.LogContext{
            RequestID:   c.GetHeader("X-Request-ID"),
            Method:      c.Request.Method,
            Path:        c.Request.URL.Path,
            IP:          c.ClientIP(),
            UserAgent:   c.GetHeader("User-Agent"),
            StartTime:   start,
        }
        
        // 获取用户信息（需要根据实际情况实现）
        if userID := getUserID(c); userID > 0 {
            logCtx.UserID = userID
            logCtx.Username = getUsername(c)
            logCtx.DeptID = getDeptID(c)
            logCtx.DeptName = getDeptName(c)
        }
        
        // 继续处理请求
        c.Next()
        
        // 设置结束时间
        logCtx.EndTime = time.Now()
        
        // 检查是否应该记录此方法的日志
        if !operLogPlugin.ShouldRecordMethod(c.Request.Method) {
            return
        }
        
        // 获取响应数据
        if len(c.Errors) > 0 {
            logCtx.Status = 1 // 异常
            logCtx.Error = c.Errors.String()
        }
        
        // 获取响应体（需要根据实际情况实现）
        if responseWriter, ok := c.Writer.(*responseBodyWriter); ok {
            logCtx.Response = responseWriter.body
        }
        
        // 记录日志
        operLogService.CreateLogFromContext(logCtx)
    }
}
```

### 3. 查询操作日志

```go
// 分页查询操作日志
logs, total, err := operLogService.GetOperLogs(1, 20, "用户管理", "admin", "1", "0", "2023-01-01,2023-01-31")

// 根据ID获取详情
log, err := operLogService.GetOperLogByID(1)

// 获取统计信息
stats, err := operLogService.GetOperLogStatistics()
```

### 4. 日志管理

```go
// 删除指定日志
err := operLogService.DeleteOperLog(1)

// 批量删除日志
err := operLogService.DeleteOperLogs([]uint{1, 2, 3})

// 清空所有日志
err := operLogService.CleanOperLog()

// 导出日志
logs, err := operLogService.ExportOperLog([]uint{1, 2, 3})
```

## 配置说明

### 基础配置

```yaml
operlog:
  enabled: true                    # 是否启用操作日志
  enableMask: true                # 是否启用数据脱敏
  sensitiveFields:                # 敏感字段列表
    - "password"
    - "pwd"
    - "token"
    - "secret"
    - "key"
    - "access_key"
    - "secret_key"
  maxParamLength: 2000            # 请求参数最大长度
  maxResultLength: 2000           # 返回结果最大长度
  recordGet: false                # 是否记录GET请求
  recordPost: true                # 是否记录POST请求
  recordPut: true                 # 是否记录PUT请求
  recordDelete: true              # 是否记录DELETE请求
  retentionDays: 90               # 日志保留天数
  defaultPageSize: 20             # 默认分页大小
  maxPageSize: 100                # 最大分页大小
```

### 自定义配置

```go
customConfig := &operlog.Config{
    Enabled:         true,
    EnableMask:      true,
    SensitiveFields: []string{"password", "pwd", "token"},
    MaxParamLength:  1000,
    MaxResultLength: 1000,
    RecordGet:       true,  // 记录所有请求
    RecordPost:      true,
    RecordPut:       true,
    RecordDelete:    true,
    RetentionDays:   180,   // 保留半年
    DefaultPageSize: 10,
    MaxPageSize:     50,
}

operLogPlugin := operlog.NewPlugin(db, customConfig)
```

## 数据库表结构

### system_oper_log 表（操作日志表）

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint | 主键 | PRIMARY KEY |
| title | varchar(50) | 操作模块 | |
| business_type | tinyint | 业务类型 | DEFAULT 0 |
| method | varchar(100) | 请求方法 | |
| request_method | varchar(10) | 请求方式 | |
| operator_type | tinyint | 操作类别 | DEFAULT 0 |
| oper_name | varchar(50) | 操作人员 | |
| dept_name | varchar(50) | 部门名称 | |
| oper_url | varchar(255) | 请求URL | |
| oper_ip | varchar(128) | 操作地址 | |
| oper_location | varchar(255) | 操作地点 | |
| oper_param | varchar(2000) | 请求参数 | |
| json_result | varchar(2000) | 返回参数 | |
| status | tinyint | 操作状态 | DEFAULT 0 |
| error_msg | varchar(2000) | 错误消息 | |
| oper_time | datetime | 操作时间 | DEFAULT CURRENT_TIMESTAMP |
| cost_time | bigint | 消耗时间 | |

## 业务类型说明

| 值 | 名称 | 说明 |
|----|------|------|
| 0 | 其它 | 其它操作 |
| 1 | 新增 | 新增数据 |
| 2 | 修改 | 修改数据 |
| 3 | 删除 | 删除数据 |
| 4 | 授权 | 授权操作 |
| 5 | 导出 | 导出数据 |
| 6 | 导入 | 导入数据 |

## 操作类别说明

| 值 | 名称 | 说明 |
|----|------|------|
| 0 | 其它 | 其它来源 |
| 1 | 后台用户 | 系统后台用户 |
| 2 | 手机端用户 | 手机端用户 |

## API接口规范

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/oper-log/page | 分页查询操作日志 |
| GET | /api/v1/system/oper-log/get | 获取操作日志详情 |
| DELETE | /api/v1/system/oper-log/delete | 删除操作日志 |
| DELETE | /api/v1/system/oper-log/clean | 清空操作日志 |
| GET | /api/v1/system/oper-log/export | 导出操作日志 |
| GET | /api/v1/system/oper-log/statistics | 获取统计信息 |

## 中间件实现

### 响应体捕获中间件

```go
type responseBodyWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
    r.body.Write(b)
    return r.ResponseWriter.Write(b)
}

func ResponseCaptureMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        writer := &responseBodyWriter{
            ResponseWriter: c.Writer,
            body:           &bytes.Buffer{},
        }
        c.Writer = writer
        c.Next()
    }
}
```

### 用户信息获取

```go
func getUserID(c *gin.Context) uint {
    if userID, exists := c.Get("userID"); exists {
        if id, ok := userID.(uint); ok {
            return id
        }
    }
    return 0
}

func getUsername(c *gin.Context) string {
    if username, exists := c.Get("username"); exists {
        if name, ok := username.(string); ok {
            return name
        }
    }
    return ""
}

func getDeptID(c *gin.Context) uint {
    if deptID, exists := c.Get("deptID"); exists {
        if id, ok := deptID.(uint); ok {
            return id
        }
    }
    return 0
}

func getDeptName(c *gin.Context) string {
    if deptName, exists := c.Get("deptName"); exists {
        if name, ok := deptName.(string); ok {
            return name
        }
    }
    return ""
}
```

## 数据脱敏

### 默认敏感字段

- password（密码）
- pwd（密码缩写）
- token（令牌）
- secret（密钥）
- key（密钥）
- access_key（访问密钥）
- secret_key（秘密密钥）

### 自定义脱敏规则

```go
// 自定义敏感字段
customConfig := &operlog.Config{
    SensitiveFields: []string{
        "password",
        "pwd", 
        "token",
        "secret",
        "key",
        "bank_card",    // 银行卡号
        "id_card",      // 身份证号
        "phone",        // 手机号
        "email",        // 邮箱
    },
}
```

### 脱敏示例

```json
// 原始数据
{
    "username": "admin",
    "password": "123456",
    "email": "admin@example.com"
}

// 脱敏后数据
{
    "username": "admin",
    "password": "***",
    "email": "***"
}
```

## 统计分析

### 获取统计信息

```go
stats, err := operLogService.GetOperLogStatistics()
if err == nil {
    fmt.Printf("今日操作总数: %d\n", stats["todayTotal"])
    fmt.Printf("今日异常数: %d\n", stats["todayError"])
    fmt.Printf("本周操作总数: %d\n", stats["weekTotal"])
    fmt.Printf("本月操作总数: %d\n", stats["monthTotal"])
}
```

### 统计指标

- **今日操作总数**：当天产生的操作日志数量
- **今日异常数**：当天操作失败的数量
- **本周操作总数**：本周产生的操作日志数量
- **本月操作总数**：本月产生的操作日志数量

## 日志清理

### 定时清理

```go
// 定时清理过期日志
func startLogCleanScheduler() {
    ticker := time.NewTicker(24 * time.Hour) // 每天执行一次
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            cleanExpiredLogs()
        }
    }
}

func cleanExpiredLogs() {
    // 清理90天前的日志
    cutoffTime := time.Now().AddDate(0, 0, -90)
    err := db.Where("oper_time < ?", cutoffTime).Delete(&OperLog{}).Error
    if err != nil {
        log.Printf("Failed to clean expired logs: %v", err)
    }
}
```

### 手动清理

```go
// 清空所有日志
err := operLogService.CleanOperLog()

// 删除指定ID的日志
err := operLogService.DeleteOperLog(1)

// 批量删除日志
err := operLogService.DeleteOperLogs([]uint{1, 2, 3})
```

## 注意事项

1. **性能影响**：操作日志会记录所有请求，可能影响系统性能，建议合理配置记录范围
2. **存储空间**：日志会占用大量存储空间，定期清理过期日志
3. **敏感信息**：确保启用数据脱敏，避免记录敏感信息
4. **权限控制**：操作日志应该只有管理员可以查看
5. **异常处理**：日志记录失败不应该影响主业务流程

## 最佳实践

1. **合理配置**：根据实际需求配置记录的请求方法和数据长度
2. **定期清理**：设置合理的日志保留策略，避免占用过多存储
3. **监控告警**：设置异常日志监控，及时发现问题
4. **分析优化**：定期分析操作日志，优化系统性能

## 测试

```go
func TestOperLogService(t *testing.T) {
    service := NewService(testDB, DefaultConfig())
    
    // 测试创建日志
    log := &OperLog{
        Title:        "测试操作",
        BusinessType: 1,
        Method:       "POST",
        OperName:     "test_user",
        OperUrl:      "/api/test",
        OperIp:       "127.0.0.1",
        Status:       0,
    }
    
    err := service.CreateOperLog(log)
    assert.NoError(t, err)
    assert.NotZero(t, log.ID)
    
    // 测试查询
    logs, total, err := service.GetOperLogs(1, 10, "", "", "", "", "")
    assert.NoError(t, err)
    assert.Equal(t, int64(1), total)
    assert.Len(t, logs, 1)
}
```