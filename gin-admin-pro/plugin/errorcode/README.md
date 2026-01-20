# 错误码管理插件

错误码管理插件提供系统错误码的统一管理功能，支持错误码的增删改查、分类管理、参数化消息、国际化支持等特性。

## 功能特性

- 错误码统一管理
- 错误类型分类
- 参数化错误消息
- 解决方案提供
- 错误码验证
- 批量导入导出
- 预定义错误码
- 缓存优化

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/errorcode"

// 创建错误码插件
errorCodePlugin := errorcode.NewPlugin(db, nil)

// 初始化插件
err := errorCodePlugin.Init()
if err != nil {
    log.Fatal("Failed to init error code plugin:", err)
}

// 获取错误码服务
errorCodeService := errorCodePlugin.GetService()
```

### 2. 错误码管理

```go
// 创建错误码
errorCode := &errorcode.ErrorCode{
    Type:     "business",
    Code:     5001,
    Name:     "订单不存在",
    Message:  "订单 %s 不存在",
    Solution: "请检查订单号是否正确",
    Status:   1,
}
err := errorCodeService.CreateErrorCode(errorCode)

// 查询错误码列表
errorCodes, total, err := errorCodeService.GetErrorCodes(1, 20, "business", "订单", "1")

// 根据错误码获取详情
errorCode, err := errorCodeService.GetErrorCodeByCode(5001)

// 更新错误码
errorCode.Message = "订单 %s 不存在或已删除"
err := errorCodeService.UpdateErrorCode(errorCode)

// 删除错误码
err := errorCodeService.DeleteErrorCode(1)
```

### 3. 错误消息获取

```go
// 获取错误消息
message := errorCodeService.GetErrorMessage(5001)

// 获取带参数的错误消息
message := errorCodeService.GetErrorMessageWithParams(5001, "ORD123456")

// 获取解决方案
solution := errorCodeService.GetErrorSolution(5001)
```

### 4. 错误码验证

```go
// 验证错误码是否存在
err := errorCodeService.ValidateErrorCode(5001)
if err != nil {
    // 错误码不存在
}
```

### 5. 搜索和分类

```go
// 根据类型获取错误码列表
businessErrors, err := errorCodeService.GetErrorCodesByType("business")

// 搜索错误码
searchResults, err := errorCodeService.SearchErrorCodes("订单")

// 获取所有错误类型
types, err := errorCodeService.GetErrorCodeTypes()
```

### 6. 导入导出

```go
// 导出错误码
errorCodes, err := errorCodeService.ExportErrorCodes([]int{1, 2, 3}, "business")

// 导入错误码
successCount, errorCount, err := errorCodeService.ImportErrorCodes(errorCodes)
```

## 配置说明

### 基础配置

```yaml
errorcode:
  enabled: true                   # 是否启用错误码管理
  enableCache: true               # 是否启用缓存
  cacheExpire: 3600              # 缓存过期时间（秒）
  cachePrefix: "errorcode:"      # 缓存前缀
  autoLoadPredefined: true       # 是否自动加载预定义错误码
  defaultPageSize: 20            # 默认分页大小
  maxPageSize: 100               # 最大分页大小
  enableValidation: true         # 是否启用错误码验证
```

### 自定义配置

```go
customConfig := &errorcode.Config{
    Enabled:            true,
    EnableCache:        true,
    CacheExpire:        7200,    // 2小时
    CachePrefix:        "myapp:errorcode:",
    AutoLoadPredefined: false,  // 不自动加载预定义错误码
    DefaultPageSize:    10,
    MaxPageSize:        50,
    EnableValidation:   true,
}

errorCodePlugin := errorcode.NewPlugin(db, customConfig)
```

## 数据库表结构

### system_error_code 表（错误码表）

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint | 主键 | PRIMARY KEY |
| type | varchar(100) | 错误类型 | NOT NULL |
| code | int | 错误码 | NOT NULL, UNIQUE |
| name | varchar(100) | 错误名称 | NOT NULL |
| message | varchar(500) | 错误消息 | NOT NULL |
| solution | varchar(1000) | 解决方案 | |
| status | tinyint | 状态 | DEFAULT 1 |
| remark | varchar(500) | 备注 | |
| create_by | bigint | 创建人 | |
| update_by | bigint | 更新人 | |
| create_at | datetime | 创建时间 | |
| update_at | datetime | 更新时间 | |

## 预定义错误码

### 通用错误码 (0, 1000-1999)

| 错误码 | 名称 | 消息 | 说明 |
|--------|------|------|------|
| 0 | 操作成功 | 操作成功 | 成功状态 |
| 1001 | 未知错误 | 未知错误，请联系管理员 | 系统未知错误 |
| 1002 | 参数错误 | 请求参数错误：%s | 参数验证失败 |
| 1003 | 数据不存在 | 数据不存在：%s | 查询的数据不存在 |
| 1004 | 数据已存在 | 数据已存在：%s | 创建的数据已存在 |
| 1005 | 操作失败 | 操作失败：%s | 业务操作失败 |
| 1006 | 权限不足 | 权限不足，无法执行此操作 | 权限验证失败 |
| 1007 | Token已过期 | 登录已过期，请重新登录 | 认证过期 |
| 1008 | Token无效 | 登录信息无效，请重新登录 | 认证无效 |
| 1009 | 请求频率超限 | 请求过于频繁，请稍后再试 | 频率限制 |

### 用户相关错误码 (2000-2999)

| 错误码 | 名称 | 消息 | 说明 |
|--------|------|------|------|
| 2001 | 用户不存在 | 用户不存在 | 用户查询失败 |
| 2002 | 用户已存在 | 用户已存在：%s | 用户创建冲突 |
| 2003 | 用户已禁用 | 用户已被禁用 | 用户状态异常 |
| 2004 | 密码错误 | 密码错误，请重新输入 | 登录密码错误 |
| 2005 | 账户已锁定 | 账户已锁定，请联系管理员 | 账户状态异常 |
| 2006 | 需要登录 | 请先登录后再进行操作 | 需要认证 |

### 角色相关错误码 (3000-3999)

| 错误码 | 名称 | 消息 | 说明 |
|--------|------|------|------|
| 3001 | 角色不存在 | 角色不存在 | 角色查询失败 |
| 3002 | 角色已存在 | 角色已存在：%s | 角色创建冲突 |
| 3003 | 角色正在使用 | 角色正在使用中，无法删除 | 角色删除限制 |

### 权限相关错误码 (4000-4999)

| 错误码 | 名称 | 消息 | 说明 |
|--------|------|------|------|
| 4001 | 权限不存在 | 权限不存在 | 权限查询失败 |
| 4002 | 权限已存在 | 权限已存在：%s | 权限创建冲突 |

### 业务相关错误码 (5000-9999)

| 错误码 | 名称 | 消息 | 说明 |
|--------|------|------|------|
| 5001 | 业务逻辑错误 | 业务逻辑错误：%s | 业务规则违反 |
| 5002 | 数据无效 | 数据无效：%s | 数据验证失败 |
| 5003 | 配置错误 | 配置错误：%s | 系统配置异常 |
| 5004 | 服务不可用 | 服务暂时不可用，请稍后再试 | 服务异常 |

## API接口规范

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/system/error-code/page | 分页查询错误码 |
| GET | /api/v1/system/error-code/get | 获取错误码详情 |
| POST | /api/v1/system/error-code/create | 创建错误码 |
| PUT | /api/v1/system/error-code/update | 更新错误码 |
| DELETE | /api/v1/system/error-code/delete | 删除错误码 |
| GET | /api/v1/system/error-code/types | 获取错误类型列表 |
| GET | /api/v1/system/error-code/search | 搜索错误码 |
| POST | /api/v1/system/error-code/import | 导入错误码 |
| GET | /api/v1/system/error-code/export | 导出错误码 |

## 在业务中使用

### 1. 返回错误响应

```go
// 在控制器中使用
func GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    // 查询用户
    user, err := userService.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 使用错误码
            message := errorCodeService.GetErrorMessageWithParams(errorcode.CodeDataNotFound, "用户")
            c.JSON(200, gin.H{
                "code": errorcode.CodeDataNotFound,
                "msg":  message,
                "data": nil,
            })
            return
        }
        
        // 未知错误
        message := errorCodeService.GetErrorMessage(errorcode.CodeUnknownError)
        c.JSON(200, gin.H{
            "code": errorcode.CodeUnknownError,
            "msg":  message,
            "data": nil,
        })
        return
    }
    
    c.JSON(200, gin.H{
        "code": errorcode.CodeSuccess,
        "msg":  errorCodeService.GetErrorMessage(errorcode.CodeSuccess),
        "data": user,
    })
}
```

### 2. 统一错误处理

```go
// 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 处理panic
                message := errorCodeService.GetErrorMessage(errorcode.CodeUnknownError)
                c.JSON(500, gin.H{
                    "code": errorcode.CodeUnknownError,
                    "msg":  message,
                    "data": nil,
                })
                c.Abort()
            }
        }()
        
        c.Next()
        
        // 处理业务错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            // 根据错误类型返回对应的错误码
            var code int
            var message string
            
            switch e := err.Err.(type) {
            case *BusinessError:
                code = e.Code
                message = errorCodeService.GetErrorMessageWithParams(e.Code, e.Params...)
            case *ValidationError:
                code = errorcode.CodeParamError
                message = errorCodeService.GetErrorMessageWithParams(errorcode.CodeParamError, e.Message)
            default:
                code = errorcode.CodeUnknownError
                message = errorCodeService.GetErrorMessage(errorcode.CodeUnknownError)
            }
            
            c.JSON(200, gin.H{
                "code": code,
                "msg":  message,
                "data": nil,
            })
        }
    }
}
```

### 3. 业务错误定义

```go
// 业务错误结构
type BusinessError struct {
    Code   int
    Params []interface{}
}

func (e *BusinessError) Error() string {
    return errorCodeService.GetErrorMessageWithParams(e.Code, e.Params...)
}

// 验证错误结构
type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

// 创建业务错误
func NewBusinessError(code int, params ...interface{}) *BusinessError {
    return &BusinessError{
        Code:   code,
        Params: params,
    }
}

// 创建验证错误
func NewValidationError(message string) *ValidationError {
    return &ValidationError{Message: message}
}
```

### 4. 在服务层使用

```go
// 用户服务示例
func (s *UserService) CreateUser(user *User) error {
    // 检查用户是否已存在
    var existingUser User
    err := s.db.Where("username = ?", user.Username).First(&existingUser).Error
    if err == nil {
        return NewBusinessError(errorcode.CodeUserExists, user.Username)
    }
    if err != gorm.ErrRecordNotFound {
        return err
    }
    
    // 创建用户
    if err := s.db.Create(user).Error; err != nil {
        return NewBusinessError(errorcode.CodeOperationFailed, "创建用户失败")
    }
    
    return nil
}
```

## 错误码规范

### 1. 错误码分配

- **0**: 成功状态
- **1000-1999**: 通用错误
- **2000-2999**: 用户相关错误
- **3000-3999**: 角色相关错误
- **4000-4999**: 权限相关错误
- **5000-5999**: 订单相关错误
- **6000-6999**: 支付相关错误
- **7000-7999**: 文件相关错误
- **8000-8999**: 消息相关错误
- **9000-9999**: 其他业务错误

### 2. 错误消息格式

- **简单消息**: 直接描述错误
- **参数化消息**: 使用占位符 `%s`，支持动态参数
- **国际化消息**: 支持多语言（扩展功能）

### 3. 错误类型分类

- **common**: 通用错误
- **user**: 用户相关错误
- **role**: 角色相关错误
- **permission**: 权限相关错误
- **business**: 业务逻辑错误
- **system**: 系统相关错误

## 缓存策略

### Redis缓存结构

```
errorcode:code:{code}           # 错误码信息
errorcode:type:{type}           # 错误类型列表
errorcode:all                    # 所有错误码列表
```

### 缓存更新策略

1. **创建/更新/删除时自动刷新**：修改错误码时自动清除相关缓存
2. **定时刷新**：定时从数据库重新加载错误码
3. **手动刷新**：提供接口手动刷新缓存

## 注意事项

1. **错误码唯一性**：确保错误码在整个系统中唯一
2. **消息描述性**：错误消息应该清晰、准确、有帮助
3. **解决方案提供**：为常见错误提供解决方案
4. **国际化支持**：考虑多语言支持的需求
5. **缓存一致性**：确保缓存与数据库的一致性

## 最佳实践

1. **合理分类**：按业务模块合理分类错误码
2. **消息参数化**：使用参数化消息提高复用性
3. **错误恢复**：提供明确的错误恢复指导
4. **监控告警**：对关键错误设置监控告警
5. **文档维护**：及时更新错误码文档

## 测试

```go
func TestErrorCodeService(t *testing.T) {
    service := NewService(testDB, DefaultConfig())
    
    // 测试创建错误码
    errorCode := &ErrorCode{
        Type:    "test",
        Code:    9999,
        Name:    "测试错误",
        Message: "这是一个测试错误：%s",
        Status:  1,
    }
    
    err := service.CreateErrorCode(errorCode)
    assert.NoError(t, err)
    
    // 测试获取错误消息
    message := service.GetErrorMessageWithParams(9999, "参数")
    assert.Equal(t, "这是一个测试错误：参数", message)
    
    // 测试验证错误码
    err = service.ValidateErrorCode(9999)
    assert.NoError(t, err)
    
    err = service.ValidateErrorCode(8888)
    assert.Error(t, err)
}
```