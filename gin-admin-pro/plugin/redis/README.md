# Redis 插件

## 功能说明

Redis 缓存插件，提供了基于 go-redis/v9 的 Redis 数据库连接和管理功能，支持单机、集群、哨兵和分片等多种部署模式。

## 特性

- ✅ 连接池管理
- ✅ 自动重连机制
- ✅ 多种部署模式支持
- ✅ 缓存抽象层
- ✅ JSON 对象缓存
- ✅ 内存缓存实现
- ✅ 缓存回源（Cache Aside）
- ✅ 过期时间管理

## 使用方法

### 1. 配置文件

在 `config.yaml` 中添加 Redis 配置：

```yaml
database:
  redis:
    # 单机模式
    addr: localhost:6379
    password: ""
    db: 0
    username: ""
    maxRetries: 3
    minRetryBackoff: 8ms
    maxRetryBackoff: 512ms
    dialTimeout: 5s
    readTimeout: 3s
    writeTimeout: 3s
    poolSize: 10
    minIdleConns: 5
    maxIdleConns: 100
    connMaxIdleTime: 5m
    connMaxLifetime: 1h
    
    # 集群模式
    clusterEnabled: false
    clusterAddrs:
      - "redis-node1:6379"
      - "redis-node2:6379"
      - "redis-node3:6379"
    
    # 哨兵模式
    sentinelEnabled: false
    sentinelAddrs:
      - "redis-sentinel1:26379"
      - "redis-sentinel2:26379"
      - "redis-sentinel3:26379"
    sentinelMaster: mymaster
    
    # 分片模式
    shardEnabled: false
    shardAddrs:
      - "redis-shard1:6379"
      - "redis-shard2:6379"
      - "redis-shard3:6379"
    sharding:
      type: "consistent"  # consistent, ketama, range, fnv1a
      keyFunc: "key"
      options:
        replicas: 150
```

### 2. 代码使用

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/redis"
    "github.com/spf13/viper"
)

func main() {
    // 加载配置
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()
    
    // 解析配置
    var cfg redis.Config
    viper.UnmarshalKey("database.redis", &cfg)
    
    // 创建客户端
    client, err := redis.NewClient(&cfg)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 创建缓存实例
    cache := redis.NewRedisCache(client)
    
    ctx := context.Background()
    
    // 健康检查
    if err := client.Ping(ctx); err != nil {
        panic("Redis 连接失败")
    }
    
    // 设置缓存
    err = cache.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        panic(err)
    }
    
    // 获取缓存
    value, err := cache.Get(ctx, "key")
    if err != nil {
        panic(err)
    }
    fmt.Println("Value:", value)
}
```

## 基础操作

### 字符串操作

```go
ctx := context.Background()

// 设置键值对
err := client.Set(ctx, "user:1", "张三", time.Hour)

// 获取值
name, err := client.Get(ctx, "user:1")

// 删除键
err = client.Del(ctx, "user:1")

// 检查键是否存在
exists, err := client.Exists(ctx, "user:1")

// 设置过期时间
err = client.Expire(ctx, "user:1", time.Minute*30)

// 获取剩余过期时间
ttl, err := client.TTL(ctx, "user:1")
```

### 哈希操作

```go
// 设置哈希字段
err := client.HSet(ctx, "profile:1", map[string]interface{}{
    "name": "张三",
    "age": 30,
    "email": "zhangsan@example.com",
})

// 获取哈希字段
name, err := client.HGet(ctx, "profile:1", "name")

// 获取所有哈希字段
profile, err := client.HGetAll(ctx, "profile:1")

// 删除哈希字段
err = client.HDel(ctx, "profile:1", "email")
```

### 集合操作

```go
// 添加到集合
err := client.SAdd(ctx, "tags:1", "golang", "redis", "cache")

// 获取集合所有成员
tags, err := client.SMembers(ctx, "tags:1")

// 从集合删除成员
err = client.SRem(ctx, "tags:1", "cache")
```

### 有序集合操作

```go
// 添加到有序集合
err := client.ZAdd(ctx, "leaderboard", redis.Z{
    Score:  100,
    Member: "user1",
})

// 获取排行榜
leaders, err := client.ZRangeWithScores(ctx, "leaderboard", 0, -1)
```

### 列表操作

```go
// 从左侧推入
err := client.LPush(ctx, "queue:tasks", "task1", "task2")

// 从右侧弹出
task, err := client.RPop(ctx, "queue:tasks")

// 获取列表范围内的元素
tasks, err := client.LRange(ctx, "queue:tasks", 0, -1)
```

### 计数器操作

```go
// 递增
count, err := client.Incr(ctx, "counter:views")

// 按指定值递增
count, err = client.IncrBy(ctx, "counter:views", 10)

// 递减
count, err = client.Decr(ctx, "counter:likes")
```

## 缓存抽象层

### 基础缓存操作

```go
cache := redis.NewRedisCache(client)
ctx := context.Background()

// 设置字符串缓存
err := cache.Set(ctx, "cache:key", "value", time.Hour)

// 获取字符串缓存
value, err := cache.Get(ctx, "cache:key")

// 检查键是否存在
exists, err := cache.Exists(ctx, "cache:key")

// 设置JSON对象缓存
user := User{
    ID:    1,
    Name:  "张三",
    Email: "zhangsan@example.com",
}
err = cache.SetJSON(ctx, "cache:user:1", user, time.Hour)

// 获取JSON对象缓存
var retrieved User
err = cache.GetJSON(ctx, "cache:user:1", &retrieved)
```

### Cache Aside 模式

```go
// GetOrSet - 缓存回源模式
user, err := cache.GetOrSet(ctx, "user:1", func() (interface{}, error) {
    // 缓存不存在时，从数据库获取数据
    return getUserFromDatabase(1)
}, time.Hour)

// GetOrSetJSON - JSON对象回源模式
var user User
err = cache.GetOrSetJSON(ctx, "user:1", &user, func() (interface{}, error) {
    return getUserFromDatabase(1)
}, time.Hour)
```

### 内存缓存

```go
// 内存缓存（适用于测试或小型应用）
memCache := redis.NewMemoryCache()

// 操作与Redis缓存相同
err := memCache.Set(ctx, "key", "value", time.Minute)
value, err := memCache.Get(ctx, "key")
```

## 高级用法

### 分布式锁

```go
// 使用 SETNX 实现简单分布式锁
func acquireLock(client *redis.Client, key string, expiration time.Duration) (bool, error) {
    ctx := context.Background()
    result, err := client.SetNX(ctx, key, "locked", expiration).Result()
    return result, err
}

func releaseLock(client *redis.Client, key string) error {
    ctx := context.Background()
    return client.Del(ctx, key).Err()
}
```

### 限流器

```go
// 使用 Redis 实现滑动窗口限流
func isAllowed(client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
    ctx := context.Background()
    now := time.Now().Unix()
    
    // 清理过期的访问记录
    client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", now-window.Seconds()))
    
    // 添加当前访问
    client.ZAdd(ctx, key, redis.Z{
        Score:  float64(now),
        Member: now,
    })
    
    // 获取窗口内的访问次数
    count, err := client.ZCard(ctx, key).Result()
    if err != nil {
        return false, err
    }
    
    return count <= int64(limit), nil
}
```

### 消息队列

```go
// 使用 Redis List 实现简单消息队列
func publishMessage(client *redis.Client, queue string, message string) error {
    ctx := context.Background()
    return client.LPush(ctx, queue, message).Err()
}

func consumeMessage(client *redis.Client, queue string) (string, error) {
    ctx := context.Background()
    return client.RPop(ctx, queue).Result()
}
```

## 配置参数

### 单机模式

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| addr | string | localhost:6379 | 服务器地址 |
| password | string | - | 密码 |
| db | int | 0 | 数据库编号 |
| username | string | - | 用户名（Redis 6.0+） |
| maxRetries | int | 3 | 最大重试次数 |
| minRetryBackoff | duration | 8ms | 最小重试间隔 |
| maxRetryBackoff | duration | 512ms | 最大重试间隔 |
| dialTimeout | duration | 5s | 连接超时 |
| readTimeout | duration | 3s | 读取超时 |
| writeTimeout | duration | 3s | 写入超时 |
| poolSize | int | 10 | 连接池大小 |
| minIdleConns | int | 5 | 最小空闲连接 |
| maxIdleConns | int | 100 | 最大空闲连接 |
| connMaxIdleTime | duration | 5m | 连接最大空闲时间 |
| connMaxLifetime | duration | 1h | 连接最大生命周期 |

### 集群模式

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| clusterEnabled | bool | false | 是否启用集群模式 |
| clusterAddrs | []string | - | 集群节点地址列表 |

### 哨兵模式

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| sentinelEnabled | bool | false | 是否启用哨兵模式 |
| sentinelAddrs | []string | - | 哨兵地址列表 |
| sentinelMaster | string | mymaster | 主节点名称 |

### 分片模式

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| shardEnabled | bool | false | 是否启用分片模式 |
| shardAddrs | []string | - | 分片节点地址列表 |
| sharding.type | string | consistent | 分片算法 |
| sharding.keyFunc | string | key | 键函数 |
| sharding.options | map | - | 算法参数 |

## API 参考

### Client

#### NewClient(cfg *Config) (*Client, error)
创建 Redis 客户端

#### GetClient() redis.Cmdable
获取原生 Redis 客户端

#### Close() error
关闭连接

#### Ping(ctx context.Context) error
检查连接状态

#### 基础操作方法
Set, Get, Del, Exists, Expire, TTL

#### 哈希操作
HSet, HGet, HGetAll, HDel

#### 集合操作
SAdd, SMembers, SRem

#### 有序集合操作
ZAdd, ZRange, ZRangeWithScores, ZRem

#### 列表操作
LPush, RPush, LPop, RPop, LRange

#### 计数器操作
Incr, Decr, IncrBy, DecrBy

### Cache Interface

#### Set(ctx, key, value, expiration) error
设置缓存

#### Get(ctx, key) (string, error)
获取缓存

#### GetOrSet(ctx, key, callback, expiration) (interface{}, error)
获取或设置缓存（回源模式）

#### SetJSON(ctx, key, obj, expiration) error
设置JSON对象缓存

#### GetJSON(ctx, key, obj) error
获取JSON对象缓存

#### 其他方法
Del, Exists, Expire, TTL, GetBytes

## 分片算法

### consistent（一致性哈希）
- 默认算法，适合动态扩缩容
- 虚拟节点数可通过 options.replicas 配置

### ketama（Ketama哈希）
- Memcached 一致性哈希算法
- 兼容性好，性能优秀

### range（范围分片）
- 按键的范围分片
- 适合有序数据

### fnv1a（FNV-1a哈希）
- 高性能哈希算法
- 适合静态分片

## 依赖

- github.com/redis/go-redis/v9 v9.17.2
- github.com/cespare/xxhash/v2 v2.3.0
- github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f

## 数据库要求

- Redis 5.0+
- 推荐使用 Redis 6.0+ 以获得更好的 ACL 支持

## 注意事项

1. 集群模式下，多键操作需要在同一个分片
2. 哨兵模式下，需要正确配置主节点名称
3. 分片模式下的分片算法需要根据业务选择
4. 过期时间设置要合理，避免缓存雪崩
5. 大数据量对象建议使用压缩

## 测试

```bash
go test ./plugin/redis/...
```

## 更新日志

- v1.0.0: 初始版本，支持单机模式
- v1.1.0: 添加集群和哨兵模式支持
- v1.2.0: 添加缓存抽象层和内存缓存
- v1.3.0: 添加分片模式和多种分片算法