# MySQL 插件

## 功能说明

MySQL 数据库插件，提供了基于 GORM 的 MySQL 数据库连接和管理功能。

## 特性

- ✅ 连接池管理
- ✅ 自动重连机制
- ✅ 慢查询日志
- ✅ 数据库健康检查
- ✅ 自动迁移支持
- ✅ 多级别日志控制

## 使用方法

### 1. 配置文件

在 `config.yaml` 中添加 MySQL 配置：

```yaml
database:
  mysql:
    host: localhost
    port: 3306
    database: gin_admin
    username: root
    password: password
    charset: utf8mb4
    parseTime: true
    loc: Local
    maxIdleConns: 10
    maxOpenConns: 100
    maxLifetime: 1h
    logLevel: info
    slowThreshold: 200ms
```

### 2. 代码使用

```go
package main

import (
    "gin-admin-pro/plugin/mysql"
    "github.com/spf13/viper"
)

func main() {
    // 加载配置
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()
    
    // 解析配置
    var cfg mysql.Config
    viper.UnmarshalKey("database.mysql", &cfg)
    
    // 创建客户端
    client, err := mysql.NewClient(&cfg)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 获取 GORM DB 实例
    db := client.GetDB()
    
    // 执行数据库操作
    // ...
    
    // 健康检查
    if err := client.Ping(); err != nil {
        panic("数据库连接失败")
    }
}
```

### 3. 模型定义

```go
type User struct {
    ID       uint   `gorm:"primarykey" json:"id"`
    Username string `gorm:"uniqueIndex" json:"username"`
    Email    string `gorm:"uniqueIndex" json:"email"`
    // ...
}

// 自动迁移
err := client.AutoMigrate(&User{})
```

## 配置参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| host | string | localhost | MySQL 服务器地址 |
| port | int | 3306 | MySQL 服务器端口 |
| database | string | gin_admin | 数据库名称 |
| username | string | root | 用户名 |
| password | string | password | 密码 |
| charset | string | utf8mb4 | 字符集 |
| parseTime | bool | true | 解析时间 |
| loc | string | Local | 时区 |
| maxIdleConns | int | 10 | 最大空闲连接数 |
| maxOpenConns | int | 100 | 最大打开连接数 |
| maxLifetime | duration | 1h | 连接最大生命周期 |
| logLevel | string | info | 日志级别(silent/error/warn/info) |
| slowThreshold | duration | 200ms | 慢查询阈值 |

## API 参考

### Client

#### NewClient(cfg *Config) (*Client, error)
创建 MySQL 客户端

#### GetDB() *gorm.DB
获取 GORM DB 实例

#### Close() error
关闭数据库连接

#### Ping() error
检查数据库连接

#### AutoMigrate(dst ...interface{}) error
自动迁移数据库表

#### HasTable(table string) bool
检查表是否存在

#### DropTable(table string) error
删除表

## 依赖

- gorm.io/gorm v1.25.5
- gorm.io/driver/mysql v1.5.2
- github.com/go-sql-driver/mysql v1.7.0

## 注意事项

1. 确保数据库已创建
2. 用户名密码正确
3. 网络连接正常
4. 连接池配置合理

## 测试

```bash
go test ./plugin/mysql/...
```

## 更新日志

- v1.0.0: 初始版本，支持基本连接池功能