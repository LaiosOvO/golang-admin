# Gin 框架配置实现

## 概述

本文档记录了 Gin 框架的配置和初始化过程，包括路由设置、中间件配置和优雅关闭功能的实现。

## 实现内容

### 1. 框架初始化

#### 1.1 依赖安装
```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/cors
```

#### 1.2 配置更新
在 `config/config.yaml` 中添加了服务器超时配置：
```yaml
server:
  port: 8080
  mode: debug  # debug, release, test
  name: "Gin Admin Pro"
  readTimeout: 30   # seconds
  writeTimeout: 30  # seconds
```

更新了配置结构体 `internal/pkg/config/types.go`：
```go
type ServerConfig struct {
    Port         int    `yaml:"port" json:"port"`
    Mode         string `yaml:"mode" json:"mode"`
    Name         string `yaml:"name" json:"name"`
    ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
    WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
}
```

### 2. 路由配置

#### 2.1 路由初始化 `internal/router/router.go`

实现了完整的路由初始化功能：
- Gin 引擎创建
- 中间件配置
- 健康检查接口
- API 分组路由

```go
func InitRouter() *gin.Engine {
    cfg := config.GetConfig()

    // 设置运行模式
    gin.SetMode(cfg.Server.Mode)

    // 创建 Gin 引擎
    r := gin.New()

    // 添加中间件
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // CORS 中间件
    if cfg.CORS.Enabled {
        r.Use(corsMiddleware(cfg.CORS))
    }

    // 健康检查
    r.GET("/health", healthHandler)

    // API 分组
    api := r.Group("/api")
    {
        v1 := api.Group("/v1")
        {
            // 系统管理模块
            system := v1.Group("/system")
            {
                // 用户管理路由
                user := system.Group("/user")
                {
                    user.GET("/page", nil)        // 用户分页查询
                    user.GET("/get", nil)         // 获取用户详情
                    user.POST("/create", nil)     // 创建用户
                    user.PUT("/update", nil)      // 更新用户
                    user.DELETE("/delete", nil)   // 删除用户
                }

                // 角色管理路由
                role := system.Group("/role")
                {
                    role.GET("/page", nil)
                    role.GET("/get", nil)
                    role.POST("/create", nil)
                    role.PUT("/update", nil)
                    role.DELETE("/delete", nil)
                    role.GET("/list-all-simple", nil)
                }

                // 菜单管理路由
                menu := system.Group("/menu")
                {
                    menu.GET("/list", nil)
                    menu.GET("/get", nil)
                    menu.POST("/create", nil)
                    menu.PUT("/update", nil)
                    menu.DELETE("/delete", nil)
                }

                // 部门管理路由
                dept := system.Group("/dept")
                {
                    dept.GET("/list", nil)
                    dept.GET("/get", nil)
                    dept.POST("/create", nil)
                    dept.PUT("/update", nil)
                    dept.DELETE("/delete", nil)
                }

                // 认证路由
                auth := system.Group("/auth")
                {
                    auth.POST("/login", nil)
                    auth.POST("/logout", nil)
                }
            }

            // 基础设施模块
            infra := v1.Group("/infra")
            {
                file := infra.Group("/file")
                {
                    file.POST("/upload", nil)
                }
            }

            // AI 模块
            ai := v1.Group("/ai")
            {
                ai.POST("/chat", nil)
            }
        }
    }

    return r
}
```

#### 2.2 CORS 中间件实现

基于配置文件的 CORS 中间件：
```go
func corsMiddleware(corsConfig config.CORSConfig) gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     corsConfig.AllowOrigins,
        AllowMethods:     corsConfig.AllowMethods,
        AllowHeaders:     corsConfig.AllowHeaders,
        ExposeHeaders:    corsConfig.ExposeHeaders,
        AllowCredentials: corsConfig.AllowCredentials,
        MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
    })
}
```

#### 2.3 健康检查接口

```go
r.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "服务正常",
        "data": gin.H{
            "status": "ok",
            "server": cfg.Server.Name,
        },
    })
})
```

### 3. 服务器管理

#### 3.1 服务器结构 `internal/server/server.go`

```go
type Server struct {
    httpServer *http.Server
}
```

#### 3.2 服务器启动

```go
func (s *Server) Start() error {
    cfg := config.GetConfig()
    
    // 初始化路由
    r := router.InitRouter()

    // 创建 HTTP 服务器
    s.httpServer = &http.Server{
        Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:        r,
        ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    fmt.Printf("服务器启动在端口 %d\n", cfg.Server.Port)
    fmt.Printf("健康检查: http://localhost:%d/health\n", cfg.Server.Port)
    fmt.Printf("API文档: http://localhost:%d/api/v1\n", cfg.Server.Port)

    if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        return fmt.Errorf("服务器启动失败: %v", err)
    }

    return nil
}
```

#### 3.3 优雅关闭

```go
func (s *Server) Stop() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    fmt.Println("正在关闭服务器...")

    if err := s.httpServer.Shutdown(ctx); err != nil {
        return fmt.Errorf("服务器关闭失败: %v", err)
    }

    fmt.Println("服务器已关闭")
    return nil
}

func (s *Server) WaitForShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    if err := s.Stop(); err != nil {
        fmt.Printf("关闭服务器时出错: %v\n", err)
    }
}
```

### 4. 主程序更新

更新了 `cmd/server/main.go`：
```go
func main() {
    // ... 配置加载代码 ...

    // 创建服务器实例
    s := server.NewServer()

    // 启动服务器（在 goroutine 中运行）
    go func() {
        if err := s.Start(); err != nil {
            log.Fatalf("服务器启动失败: %v", err)
        }
    }()

    // 等待关闭信号
    s.WaitForShutdown()

    fmt.Println("程序退出")
}
```

## 路由结构

### API 分组

```
/api/v1/
├── system/              # 系统管理模块
│   ├── user/           # 用户管理
│   ├── role/           # 角色管理
│   ├── menu/           # 菜单管理
│   ├── dept/           # 部门管理
│   └── auth/           # 认证模块
├── infra/               # 基础设施模块
│   └── file/           # 文件上传
└── ai/                  # AI 模块
    └── chat/           # AI 对话
```

### 健康检查

- **路径**: `/health`
- **方法**: GET
- **响应格式**:
```json
{
    "code": 0,
    "message": "服务正常",
    "data": {
        "status": "ok",
        "server": "Gin Admin Pro"
    }
}
```

## 中间件

### 已实现的中间件

1. **Logger**: HTTP 请求日志
2. **Recovery**: 错误恢复处理
3. **CORS**: 跨域资源共享

### 中间件顺序

```go
r.Use(gin.Logger())        // 请求日志
r.Use(gin.Recovery())      // 错误恢复
r.Use(corsMiddleware())    // CORS 处理
```

## 配置说明

### 运行模式

- **debug**: 开发模式，详细日志输出
- **release**: 生产模式，优化性能
- **test**: 测试模式

### CORS 配置

在 `config/config.yaml` 中配置：
```yaml
cors:
  enabled: true
  allowOrigins: ["*"]
  allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowHeaders: ["*"]
  exposeHeaders: []
  allowCredentials: false
  maxAge: 86400
```

## 测试验证

### 1. 编译测试

```bash
go build -o ./bin/server ./cmd/server
```

### 2. 运行测试

```bash
./bin/server
```

### 3. 健康检查测试

```bash
curl http://localhost:8080/health
```

### 4. 路由测试

所有路由都已定义，但目前返回 nil，将在后续阶段实现具体的控制器。

## 与 ruoyi-vue-pro 的对应

所有路由路径严格按照 ruoyi-vue-pro 的接口规范设计，确保前后端接口的兼容性：

- 用户管理：`/api/v1/system/user/*`
- 角色管理：`/api/v1/system/role/*`
- 菜单管理：`/api/v1/system/menu/*`
- 部门管理：`/api/v1/system/dept/*`
- 认证：`/api/v1/system/auth/*`

## 下一步计划

1. 实现具体的控制器和业务逻辑
2. 添加更多中间件（认证、限流等）
3. 完善错误处理
4. 添加 API 文档

## 技术要点

1. **模块化设计**: 路由、服务器、配置分离
2. **配置驱动**: 所有参数通过配置文件控制
3. **优雅关闭**: 支持信号处理和超时控制
4. **中间件扩展**: 易于添加新的中间件
5. **标准接口**: 严格按照 ruoyi-vue-pro 规范设计