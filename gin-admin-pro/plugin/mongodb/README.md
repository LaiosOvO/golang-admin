# MongoDB 插件

## 功能说明

MongoDB 数据库插件，提供了基于官方 Go 驱动的 MongoDB 数据库连接和管理功能，支持文档存储、索引管理、复制集等特性。

## 特性

- ✅ 连接池管理
- ✅ 自动重连机制
- ✅ 集合管理
- ✅ 索引管理
- ✅ 数据库健康检查
- ✅ 服务器状态监控
- ✅ 支持复制集
- ✅ 支持认证

## 使用方法

### 1. 配置文件

在 `config.yaml` 中添加 MongoDB 配置：

```yaml
database:
  mongodb:
    uri: mongodb://localhost:27017
    # 或者使用详细配置
    # host: localhost
    # port: 27017
    # database: gin_admin
    # username: admin
    # password: password
    # authSource: admin
    maxPoolSize: 100
    minPoolSize: 10
    maxConnIdle: 5m
    connectTimeout: 10s
    serverTimeout: 30s
    timeout: 30s
    compressLevel: 6
    replicaSet: rs0
```

### 2. 代码使用

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/mongodb"
    "github.com/spf13/viper"
)

func main() {
    // 加载配置
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()
    
    // 解析配置
    var cfg mongodb.Config
    viper.UnmarshalKey("database.mongodb", &cfg)
    
    // 创建客户端
    client, err := mongodb.NewClient(&cfg)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 获取数据库实例
    db := client.GetDatabase()
    
    // 获取集合实例
    collection := db.Collection("users")
    
    // 健康检查
    if err := client.Ping(); err != nil {
        panic("MongoDB 连接失败")
    }
    
    // 列出集合
    ctx := context.Background()
    collections, err := client.ListCollections(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println("Collections:", collections)
}
```

## 基础操作

### 创建文档

```go
type User struct {
    ID       string `bson:"_id,omitempty" json:"id"`
    Name     string `bson:"name" json:"name"`
    Email    string `bson:"email" json:"email"`
    Age      int    `bson:"age" json:"age"`
    Created  time.Time `bson:"created" json:"created"`
}

ctx := context.Background()
collection := client.GetCollection("users")

user := User{
    Name:    "张三",
    Email:   "zhangsan@example.com",
    Age:     30,
    Created: time.Now(),
}

result, err := collection.InsertOne(ctx, user)
if err != nil {
    panic(err)
}
fmt.Println("Inserted ID:", result.InsertedID)
```

### 查询文档

```go
// 单条查询
var user User
err := collection.FindOne(ctx, bson.M{"name": "张三"}).Decode(&user)
if err != nil {
    panic(err)
}
fmt.Println("Found user:", user)

// 多条查询
cursor, err := collection.Find(ctx, bson.M{"age": bson.M{"$gt": 25}})
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)

var users []User
if err = cursor.All(ctx, &users); err != nil {
    panic(err)
}
fmt.Println("Users older than 25:", users)
```

### 更新文档

```go
// 单条更新
filter := bson.M{"name": "张三"}
update := bson.M{"$set": bson.M{"age": 31}}

result, err := collection.UpdateOne(ctx, filter, update)
if err != nil {
    panic(err)
}
fmt.Println("Updated documents:", result.ModifiedCount)

// 多条更新
filter = bson.M{"age": bson.M{"$lt": 30}}
update = bson.M{"$inc": bson.M{"age": 1}}

result, err = collection.UpdateMany(ctx, filter, update)
if err != nil {
    panic(err)
}
fmt.Println("Updated documents:", result.ModifiedCount)
```

### 删除文档

```go
// 单条删除
filter := bson.M{"name": "张三"}
result, err := collection.DeleteOne(ctx, filter)
if err != nil {
    panic(err)
}
fmt.Println("Deleted documents:", result.DeletedCount)

// 多条删除
filter = bson.M{"age": bson.M{"$lt": 20}}
result, err = collection.DeleteMany(ctx, filter)
if err != nil {
    panic(err)
}
fmt.Println("Deleted documents:", result.DeletedCount)
```

## 索引管理

### 创建索引

```go
ctx := context.Background()

// 单字段索引
indexModel := mongo.IndexModel{
    Keys: bson.M{"email": 1}, // 1 表示升序，-1 表示降序
    Options: options.Index().
        SetUnique(true).
        SetName("email_unique"),
}
err := client.CreateIndex(ctx, "users", indexModel)
if err != nil {
    panic(err)
}

// 复合索引
compoundIndex := mongo.IndexModel{
    Keys: bson.M{"name": 1, "age": -1},
    Options: options.Index().SetName("name_age_idx"),
}

// 多个索引
indexes := []mongo.IndexModel{indexModel, compoundIndex}
err = client.CreateIndexes(ctx, "users", indexes)
if err != nil {
    panic(err)
}
```

### 列出和删除索引

```go
// 列出索引
indexes, err := client.ListIndexes(ctx, "users")
if err != nil {
    panic(err)
}
fmt.Println("Indexes:", indexes)

// 删除索引
err = client.DropIndex(ctx, "users", "email_unique")
if err != nil {
    panic(err)
}
```

## 集合管理

```go
ctx := context.Background()

// 创建集合
err := client.CreateCollection(ctx, "logs")
if err != nil {
    panic(err)
}

// 检查集合是否存在
exists, err := client.HasCollection(ctx, "users")
if err != nil {
    panic(err)
}
fmt.Println("Users collection exists:", exists)

// 删除集合
err = client.DropCollection(ctx, "logs")
if err != nil {
    panic(err)
}
```

## 监控和统计

```go
ctx := context.Background()

// 获取服务器状态
serverStatus, err := client.GetServerStatus(ctx)
if err != nil {
    panic(err)
}
fmt.Printf("Server Status: %+v\n", serverStatus)

// 获取数据库统计
dbStats, err := client.GetDatabaseStats(ctx)
if err != nil {
    panic(err)
}
fmt.Printf("Database Stats: %+v\n", dbStats)

// 获取集合统计
collStats, err := client.GetCollectionStats(ctx, "users")
if err != nil {
    panic(err)
}
fmt.Printf("Collection Stats: %+v\n", collStats)
```

## 高级查询示例

### 聚合管道

```go
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"age": bson.M{"$gte": 18}}}},
    {{"$group", bson.M{
        "_id": "$email",
        "count": bson.M{"$sum": 1},
        "avgAge": bson.M{"$avg": "$age"},
    }}},
    {{"$sort", bson.M{"count": -1}}},
    {{"$limit", 10}},
}

cursor, err := collection.Aggregate(ctx, pipeline)
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)

var results []bson.M
if err = cursor.All(ctx, &results); err != nil {
    panic(err)
}
```

### 文本搜索

```go
// 需要先创建文本索引
textIndex := mongo.IndexModel{
    Keys: bson.M{"name": "text", "email": "text"},
    Options: options.Index().SetName("text_search_idx"),
}
client.CreateIndex(ctx, "users", textIndex)

// 执行文本搜索
cursor, err := collection.Find(ctx, bson.M{
    "$text": bson.M{"$search": "张三"},
})
```

## 配置参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| uri | string | mongodb://localhost:27017 | MongoDB 连接URI |
| host | string | localhost | 服务器地址 |
| port | int | 27017 | 服务器端口 |
| database | string | gin_admin | 数据库名称 |
| username | string | - | 用户名 |
| password | string | - | 密码 |
| authSource | string | admin | 认证数据库 |
| maxPoolSize | uint64 | 100 | 最大连接池大小 |
| minPoolSize | uint64 | 10 | 最小连接池大小 |
| maxConnIdle | duration | 5m | 连接最大空闲时间 |
| connectTimeout | duration | 10s | 连接超时 |
| serverTimeout | duration | 30s | 服务器选择超时 |
| timeout | duration | 30s | 操作超时 |
| compressLevel | int | 6 | 压缩级别 |
| replicaSet | string | - | 复制集名称 |

## API 参考

### Client

#### NewClient(cfg *Config) (*Client, error)
创建 MongoDB 客户端

#### GetClient() *mongo.Client
获取官方 MongoDB 客户端实例

#### GetDatabase() *mongo.Database
获取数据库实例

#### GetCollection(name string) *mongo.Collection
获取集合实例

#### Close() error
关闭数据库连接

#### Ping() error
检查数据库连接

#### ListCollections(ctx context.Context) ([]string, error)
列出所有集合

#### CreateCollection(ctx context.Context, name string) error
创建集合

#### DropCollection(ctx context.Context, name string) error
删除集合

#### HasCollection(ctx context.Context, name string) (bool, error)
检查集合是否存在

#### CreateIndex(ctx context.Context, collection string, index mongo.IndexModel) error
创建单个索引

#### CreateIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error
创建多个索引

#### DropIndex(ctx context.Context, collection string, indexName string) error
删除索引

#### ListIndexes(ctx context.Context, collection string) ([]string, error)
列出索引

#### GetServerStatus(ctx context.Context) (interface{}, error)
获取服务器状态

#### GetDatabaseStats(ctx context.Context) (interface{}, error)
获取数据库统计信息

#### GetCollectionStats(ctx context.Context, collection string) (interface{}, error)
获取集合统计信息

## 数据类型映射

| Go 类型 | BSON 类型 |
|---------|-----------|
| string | String |
| int, int32, int64 | Int32, Int64 |
| float32, float64 | Double |
| bool | Boolean |
| time.Time | Date |
| []byte | Binary |
| nil | Null |
| map[string]interface{} | Object |
| []interface{} | Array |
| ObjectID | ObjectID |

## 依赖

- go.mongodb.org/mongo-driver v1.17.6
- github.com/golang/snappy v0.0.4
- github.com/klauspost/compress v1.17.0

## 数据库要求

- MongoDB 4.4+
- 建议使用 MongoDB 5.0+ 以获得最佳性能

## 注意事项

1. 确保数据库已创建
2. 集合会在首次插入时自动创建
3. 索引创建需要时间，大数据量时建议在低峰期操作
4. 复制集配置需要所有节点正常运行
5. 压缩功能会消耗额外CPU资源

## 测试

```bash
go test ./plugin/mongodb/...
```

## 更新日志

- v1.0.0: 初始版本，支持基本CRUD操作
- v1.1.0: 添加索引管理和监控功能
- v1.2.0: 添加复制集支持