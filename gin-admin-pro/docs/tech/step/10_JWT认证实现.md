# JWT 认证功能实现

## 概述

本文档详细记录了项目中JWT认证功能的完整实现，包括Token生成、验证、刷新、Redis存储和单点登录支持。

## 功能特性

### 1. 核心JWT功能
- ✅ JWT Token 生成和解析
- ✅ 访问Token和刷新Token分离
- ✅ Token 过期时间控制
- ✅ 自定义Claims结构

### 2. Redis 存储
- ✅ Token Redis 黑名单机制
- ✅ Token 过期时间管理
- ✅ 用户Token映射关系
- ✅ 高性能Token验证

### 3. 单点登录(SSO)
- ✅ 强制单设备登录
- ✅ 多设备登录限制
- ✅ 设备管理功能
- ✅ 在线状态检查

### 4. 安全特性
- ✅ 密钥签名验证
- ✅ 敏感数据过滤
- ✅ Token 撤销机制
- ✅ 防止Token重用

## 架构设计

### 组件结构
```
internal/
├── pkg/
│   ├── jwt/
│   │   └── jwt.go              # JWT核心工具包
│   ├── token/
│   │   └── service.go          # Token管理服务
│   └── sso/
│       └── manager.go          # 单点登录管理器
├── middleware/
│   └── auth.go                # 认证中间件
└── service/
    └── container.go           # 服务容器
```

### 数据流程
```
用户登录 -> Token生成 -> Redis存储 -> 客户端存储
     ↓
请求访问 -> Token验证 -> Redis检查 -> 用户信息注入
     ↓
Token刷新 -> 验证刷新Token -> 生成新Token -> 更新Redis
     ↓
用户登出 -> Token撤销 -> Redis清理 -> 单点登录处理
```

## 实现细节

### 1. JWT 工具包 (`internal/pkg/jwt/jwt.go`)

#### Claims 结构
```go
type Claims struct {
    UserID   uint   `json:"userId"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

#### 核心方法
```go
// 生成访问Token
func GenerateToken(userID uint, username string) (string, error)

// 解析Token
func ParseToken(tokenString string) (*Claims, error)

// 刷新Token
func RefreshToken(refreshTokenString string) (string, error)

// 验证Token
func ValidateToken(tokenString string) bool
```

#### 配置参数
```yaml
jwt:
  secret: "your-secret-key-here"
  accessTokenExpire: 7   # days
  refreshTokenExpire: 30 # days
```

### 2. Token管理服务 (`internal/pkg/token/service.go`)

#### Token 存储策略
- **访问Token**: `jwt:access:{token}` -> `userID` (TTL: 7天)
- **刷新Token**: `jwt:refresh:{token}` -> `userID` (TTL: 30天)  
- **用户Token映射**: `jwt:user_tokens:{userID}` -> `Set{tokens}` (TTL: 30天)

#### 核心方法
```go
// 生成Token对
func (s *TokenService) GenerateTokens(userID uint, username string) (*TokenPair, error)

// 验证Token
func (s *TokenService) ValidateToken(tokenString string) (*TokenInfo, error)

// 刷新Token
func (s *TokenService) RefreshToken(refreshToken string) (*TokenPair, error)

// 撤销Token
func (s *TokenService) RevokeToken(tokenString string) error

// 撤销用户所有Token
func (s *TokenService) RevokeAllUserTokens(userID uint) error
```

#### Token 结构
```go
// Token对
type TokenPair struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    ExpiresIn    int64  `json:"expiresIn"`
    TokenType    string `json:"tokenType"`
}

// Token信息
type TokenInfo struct {
    UserID    uint      `json:"userId"`
    Username  string    `json:"username"`
    Subject   string    `json:"subject"`
    ExpiresAt time.Time `json:"expiresAt"`
    IssuedAt  time.Time `json:"issuedAt"`
}
```

### 3. 单点登录管理器 (`internal/pkg/sso/manager.go`)

#### 登录选项
```go
type LoginOptions struct {
    EnableSSO        bool   `json:"enableSSO"`        // 是否启用单点登录
    MaxDevices       int    `json:"maxDevices"`       // 最大设备数量
    DeviceName       string `json:"deviceName"`       // 设备名称
    DeviceInfo       string `json:"deviceInfo"`       // 设备信息
}
```

#### 核心功能
```go
// 启用单点登录
func (s *SSOManager) EnableSingleSignOn(userID uint) error

// 允许多设备登录
func (s *SSOManager) EnableMultiDeviceLogin(userID uint) error

// 检查用户是否在线
func (s *SSOManager) IsUserOnline(userID uint) (bool, error)

// 获取活跃设备数
func (s *SSOManager) GetActiveDevices(userID uint) (int, error)

// 撤销指定设备
func (s *SSOManager) RevokeDevice(userID uint, deviceToken string) error

// 撤销所有设备
func (s *SSOManager) RevokeAllDevices(userID uint) error
```

### 4. 认证中间件 (`internal/middleware/auth.go`)

#### 认证流程
```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取Authorization头
        authHeader := c.GetHeader("Authorization")
        
        // 2. 检查Bearer格式
        if !strings.HasPrefix(authHeader, "Bearer ") {
            // 返回401错误
        }
        
        // 3. 提取Token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 4. 验证Token
        tokenInfo, err := service.Services.TokenService.ValidateToken(tokenString)
        if err != nil {
            // 返回401错误
        }
        
        // 5. 注入用户信息
        c.Set("userId", tokenInfo.UserID)
        c.Set("username", tokenInfo.Username)
        
        c.Next()
    }
}
```

## 使用示例

### 1. 用户登录
```go
// 生成Token
tokenPair, err := service.Services.TokenService.GenerateTokens(userID, username)
if err != nil {
    return error
}

// 返回给客户端
return gin.H{
    "code": 0,
    "data": tokenPair,
    "message": "登录成功",
}
```

### 2. 认证保护路由
```go
// 在路由中使用认证中间件
userGroup := system.Group("/user")
userGroup.Use(middleware.Auth())
{
    userGroup.GET("/profile", getUserProfile)
    userGroup.POST("/update", updateUserProfile)
}
```

### 3. Token 刷新
```go
func RefreshToken(c *gin.Context) {
    refreshToken := c.PostForm("refreshToken")
    
    tokenPair, err := service.Services.TokenService.RefreshToken(refreshToken)
    if err != nil {
        c.JSON(401, gin.H{"code": 401, "message": "Token刷新失败"})
        return
    }
    
    c.JSON(200, gin.H{"code": 0, "data": tokenPair})
}
```

### 4. 单点登录
```go
func Login(c *gin.Context) {
    // 验证用户名密码...
    
    // 配置单点登录选项
    loginOpts := sso.LoginOptions{
        EnableSSO:  true,  // 启用单点登录
        MaxDevices: 3,     // 最多3个设备
        DeviceName: "Web Browser",
    }
    
    // 执行登录
    ssoManager := sso.NewSSOManager()
    err := ssoManager.LoginWithOptions(userID, username, loginOpts)
    if err != nil {
        return error
    }
    
    // 生成Token...
}
```

## 安全考虑

### 1. 密钥管理
- 使用强随机密钥（至少32字节）
- 定期轮换密钥
- 密钥存储在环境变量或密钥管理系统中

### 2. Token 安全
- Token 设置合理的过期时间
- 实现Token黑名单机制
- 支持Token强制撤销

### 3. 传输安全
- 强制使用HTTPS传输
- 设置安全的Cookie属性
- 防止CSRF攻击

### 4. 存储安全
- Redis设置密码认证
- 敏感数据加密存储
- 定期清理过期Token

## 性能优化

### 1. Redis 优化
- 使用连接池管理Redis连接
- 设置合适的过期时间
- 定期清理无用的键

### 2. 验证优化
- 缓存频繁访问的Token信息
- 使用批量查询减少Redis调用
- 异步记录操作日志

### 3. 内存优化
- 避免大对象在内存中长时间驻留
- 及时释放不用的资源
- 监控内存使用情况

## 错误处理

### 1. 错误码定义
```go
const (
    ErrTokenInvalid     = 1001 // Token无效
    ErrTokenExpired     = 1002 // Token过期
    ErrTokenRevoked     = 1003 // Token已撤销
    ErrTooManyDevices   = 1004 // 设备数量超限
    ErrUserNotFound     = 1005 // 用户不存在
    ErrInvalidCredentials = 1006 // 认证信息无效
)
```

### 2. 错误响应格式
```json
{
    "code": 401,
    "message": "Token已过期",
    "data": {
        "errorCode": 1002,
        "details": "Your access token has expired"
    }
}
```

## 监控和日志

### 1. 关键指标
- Token 生成数量
- Token 验证成功率
- 用户在线数量
- 设备登录情况

### 2. 日志记录
```go
// Token生成日志
logger.Info("Token generated", 
    zap.Uint("userId", userID),
    zap.String("username", username),
    zap.String("clientIP", c.ClientIP()),
)

// Token验证日志
logger.Info("Token validated", 
    zap.Uint("userId", tokenInfo.UserID),
    zap.String("path", c.Request.URL.Path),
)
```

## 测试

### 1. 单元测试
```go
func TestTokenGeneration(t *testing.T) {
    service := NewTokenService(redisClient)
    
    tokenPair, err := service.GenerateTokens(1, "testuser")
    assert.NoError(t, err)
    assert.NotEmpty(t, tokenPair.AccessToken)
    assert.NotEmpty(t, tokenPair.RefreshToken)
}
```

### 2. 集成测试
```go
func TestAuthenticationFlow(t *testing.T) {
    // 1. 用户登录获取Token
    // 2. 使用Token访问受保护资源
    // 3. Token刷新
    // 4. 用户登出
    // 5. 验证Token已撤销
}
```

## 配置示例

### 开发环境
```yaml
jwt:
  secret: "dev-secret-key-change-in-production"
  accessTokenExpire: 1  # 1天，便于测试
  refreshTokenExpire: 7 # 7天

redis:
  host: localhost
  port: 6379
  database: 0
```

### 生产环境
```yaml
jwt:
  secret: "${JWT_SECRET}" # 从环境变量读取
  accessTokenExpire: 7
  refreshTokenExpire: 30

redis:
  host: "${REDIS_HOST}"
  port: "${REDIS_PORT}"
  password: "${REDIS_PASSWORD}"
  database: 0
```

## 下一步优化

1. **JWT Claim扩展**: 添加更多用户信息到Claims中
2. **设备指纹**: 实现设备指纹识别
3. **Token续期**: 实现滑动续期机制
4. **审计日志**: 完善认证相关的审计日志
5. **性能监控**: 添加详细的性能监控指标

## 最佳实践

1. **密钥安全**: 绝不要在代码中硬编码密钥
2. **Token大小**: 保持Token大小合理，避免过大
3. **错误处理**: 提供清晰的错误信息但不泄露敏感数据
4. **日志记录**: 记录关键操作但避免记录敏感信息
5. **定期审查**: 定期审查和更新安全配置