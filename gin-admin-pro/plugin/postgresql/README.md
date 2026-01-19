# PostgreSQL 插件

## 功能说明

PostgreSQL 数据库插件，提供了基于 GORM 的 PostgreSQL 数据库连接和管理功能，支持 PostGIS 地理信息系统扩展和 pgvector 向量数据库扩展。

## 特性

- ✅ 连接池管理
- ✅ 自动重连机制
- ✅ 慢查询日志
- ✅ 数据库健康检查
- ✅ 自动迁移支持
- ✅ 多级别日志控制
- ✅ **PostGIS 地理信息系统扩展**
- ✅ **pgvector 向量数据库扩展**

## 使用方法

### 1. 配置文件

在 `config.yaml` 中添加 PostgreSQL 配置：

```yaml
database:
  postgresql:
    host: localhost
    port: 5432
    database: gin_admin
    username: postgres
    password: password
    sslMode: disable
    timezone: Asia/Shanghai
    maxIdleConns: 10
    maxOpenConns: 100
    maxLifetime: 1h
    logLevel: info
    slowThreshold: 200ms
    extensions:
      - name: postgis
        version: "3.3"
        enabled: true
      - name: vector
        version: "0.5.1"
        enabled: true
      - name: uuid-ossp
        version: "1.1"
        enabled: true
      - name: btree_gin
        version: "1.0"
        enabled: true
```

### 2. 代码使用

```go
package main

import (
    "gin-admin-pro/plugin/postgresql"
    "github.com/spf13/viper"
)

func main() {
    // 加载配置
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()
    
    // 解析配置
    var cfg postgresql.Config
    viper.UnmarshalKey("database.postgresql", &cfg)
    
    // 创建客户端
    client, err := postgresql.NewClient(&cfg)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 获取 GORM DB 实例
    db := client.GetDB()
    
    // 检查扩展
    extensions, err := client.GetExtensions()
    if err != nil {
        panic(err)
    }
    fmt.Println("Enabled extensions:", extensions)
    
    // 健康检查
    if err := client.Ping(); err != nil {
        panic("数据库连接失败")
    }
}
```

## PostGIS 地理信息系统

### 基础地理操作

```go
// 获取 PostGIS 函数
pgFunc := client.PostGISFunc()

// 创建点
pointSQL := pgFunc.ST_Point(39.9042, 116.4074)

// 计算两点间距离
distanceSQL := pgFunc.ST_Distance(39.9042, 116.4074, 31.2304, 121.4737)

// 创建缓冲区
bufferSQL := pgFunc.ST_Buffer("geom", 1000.0) // 1公里缓冲区

// 转换为 GeoJSON
geoJSONSQL := pgFunc.ST_AsGeoJSON("geom")
```

### 地理计算工具

```go
// 计算两点间距离（米）
distance := DistanceBetween(lat1, lng1, lat2, lng2)

// 创建边界框
minLat, minLng, maxLat, maxLng := CreateBoundingBox(lat, lng, radiusKm)

// 判断点是否在多边形内
point := Point{1.0, 1.0}
polygon := []Point{{0,0}, {2,0}, {2,2}, {0,2}}
isInside := PointInPolygon(point, polygon)
```

## pgvector 向量数据库

### 基础向量操作

```go
// 获取向量函数
vecFunc := client.VectorFunc()

// 创建向量
vec := postgresql.Vector{1.0, 2.0, 3.0, 4.0}

// 计算余弦相似度
similarity, err := CosineSimilarity(vec1, vec2)

// 计算欧几里得距离
distance, err := EuclideanDistance(vec1, vec2)

// 向量归一化
normalized := vec.Normalize()
```

### 向量索引

```go
// 创建向量索引
err := vecFunc.CreateVectorIndex("products", "embedding", 768, "ivfflat")

// 创建 HNSW 索引（更高性能）
err = vecFunc.CreateVectorIndex("products", "embedding", 768, "hnsw")

// 删除索引
err = vecFunc.DropVectorIndex("products", "embedding")
```

### 相似性搜索

```go
// 相似性搜索 SQL
searchSQL := vecFunc.SimilaritySearch("products", "embedding", queryVec, 10, "ASC")

// 在查询中使用
db.Raw(searchSQL).Scan(&results)
```

### 模型定义

```go
type Product struct {
    ID        uint          `gorm:"primarykey" json:"id"`
    Name      string        `json:"name"`
    Embedding postgresql.Vector `gorm:"type:vector(768)" json:"embedding"`
    Location  postgresql.Geometry `gorm:"type:geometry" json:"location"`
}

// 自动迁移
err := client.AutoMigrate(&Product{})
```

## 配置参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| host | string | localhost | PostgreSQL 服务器地址 |
| port | int | 5432 | PostgreSQL 服务器端口 |
| database | string | gin_admin | 数据库名称 |
| username | string | postgres | 用户名 |
| password | string | password | 密码 |
| sslMode | string | disable | SSL 模式 |
| timezone | string | Asia/Shanghai | 时区 |
| maxIdleConns | int | 10 | 最大空闲连接数 |
| maxOpenConns | int | 100 | 最大打开连接数 |
| maxLifetime | duration | 1h | 连接最大生命周期 |
| logLevel | string | info | 日志级别 |
| slowThreshold | duration | 200ms | 慢查询阈值 |

### 扩展配置

| 扩展名 | 用途 | 默认版本 |
|--------|------|----------|
| postgis | 地理信息系统 | 3.3 |
| vector | 向量数据库 | 0.5.1 |
| uuid-ossp | UUID 生成 | 1.1 |
| btree_gin | GIN 索引 | 1.0 |

## API 参考

### Client

#### NewClient(cfg *Config) (*Client, error)
创建 PostgreSQL 客户端

#### GetDB() *gorm.DB
获取 GORM DB 实例

#### Close() error
关闭数据库连接

#### Ping() error
检查数据库连接

#### AutoMigrate(dst ...interface{}) error
自动迁移数据库表

#### GetExtensions() ([]string, error)
获取已启用的扩展列表

#### IsExtensionEnabled(name string) (bool, error)
检查扩展是否启用

### PostGISFunc

通过 `client.PostGISFunc()` 获取实例

#### ST_Point(lat, lng float64) string
创建点几何对象

#### ST_Distance(lat1, lng1, lat2, lng2 float64) string
计算两点间距离

#### ST_Buffer(geom string, radius float64) string
创建缓冲区

#### ST_AsGeoJSON(geom string) string
转换为 GeoJSON 格式

### VectorFunc

通过 `client.VectorFunc()` 获取实例

#### CreateVectorIndex(table, column, dimension, indexType) error
创建向量索引

#### SimilaritySearch(table, column, queryVec, limit, orderBy) string
生成相似性搜索 SQL

## 数据类型

### Vector 类型

```go
type Vector []float32
```

- 支持 768 维 OpenAI 向量
- 支持 1536 维 GPT-4 向量
- 支持自定义维度

### Geometry 类型

```go
type Geometry struct {
    Type        string      `json:"type"`
    Coordinates interface{} `json:"coordinates"`
}
```

- 支持点、线、面等几何类型
- 兼容 GeoJSON 格式

## 依赖

- gorm.io/gorm v1.25.10
- gorm.io/driver/postgres v1.6.0
- github.com/jackc/pgx/v5 v5.6.0

## 数据库要求

- PostgreSQL 12+ 
- PostGIS 3.3+
- pgvector 0.5.1+

## 注意事项

1. 确保数据库已创建
2. 扩展需要超级用户权限安装
3. 向量索引需要数据量大时才有明显效果
4. 地理计算建议使用 SRID 4326 (WGS84)

## 测试

```bash
go test ./plugin/postgresql/...
```

## 更新日志

- v1.0.0: 初始版本，支持基本连接池功能
- v1.1.0: 添加 PostGIS 地理信息系统支持
- v1.2.0: 添加 pgvector 向量数据库支持