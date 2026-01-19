# Kafka消息队列集成实现

## 实现思路

本阶段实现了完整的Kafka消息队列插件，参考了业界最佳实践，提供了企业级消息队列解决方案。

## 技术选型

### 核心库
- **IBM Sarama**: 选择Sarama作为Kafka Go客户端，因其功能完整、性能优秀、社区活跃
- **支持特性**: 完整的Kafka协议支持、生产者/消费者管理、事务、压缩等

### 架构设计
```
plugin/kafka/
├── config.go          # 配置管理
├── kafka.go           # 核心客户端
├── producer.go        # 生产者实现
├── consumer.go        # 消费者实现
├── types.go           # 数据类型定义
├── kafka_test.go      # 单元测试
└── README.md          # 使用文档
```

## 核心功能实现

### 1. 配置管理 (config.go)

#### 设计思路
- 支持完整的Kafka配置参数
- 提供合理的默认值
- 支持YAML配置文件映射
- 分离安全配置和功能配置

#### 关键特性
```go
type Config struct {
    // 基础配置
    Brokers []string
    GroupID string
    
    // 生产者配置
    Producer struct {
        RequiredAcks int16
        Compression string
        Partitioner string
        // ... 更多配置
    }
    
    // 消费者配置
    Consumer struct {
        FetchMin int32
        FetchDefault int32
        SessionTimeout time.Duration
        // ... 更多配置
    }
    
    // 安全配置
    Security struct {
        Enabled bool
        Mechanism string
        TLSEnabled bool
        // ... 更多配置
    }
}
```

### 2. 核心客户端 (kafka.go)

#### 设计思路
- 统一的客户端接口
- 自动连接管理和重连机制
- 优雅关闭和错误处理
- 支持主题管理

#### 关键接口
```go
type KafkaClient struct {
    config    *Config
    producer  sarama.SyncProducer
    consumer  sarama.ConsumerGroup
    admin     sarama.ClusterAdmin
    consumers map[string]*ConsumerGroupManager
}

// 核心方法
- SendMessage(ctx, topic, key, value) error
- SendBatchMessages(ctx, topic, messages) error
- CreateConsumer(groupID, topics, handler) error
- CreateTopic(name, partitions, replication) error
- ListTopics() (metadata, error)
```

### 3. 生产者实现 (producer.go)

#### 设计思路
- 同步和异步两种模式
- 自动重试和错误处理
- 消息压缩和批处理
- 多种分区策略

#### 关键特性
```go
type Producer struct {
    client    *KafkaClient
    syncProd  sarama.SyncProducer
    asyncProd sarama.AsyncProducer
}

// 发送模式
- SendMessageSync() - 同步发送，等待确认
- SendMessageAsync() - 异步发送，不等待确认
- SendBatchMessagesSync() - 批量同步发送
- SendBatchMessagesAsync() - 批量异步发送

// 配置选项
- 压缩算法: none/gzip/snappy/lz4/zstd
- 分区策略: hash/random/roundrobin
- 确认级别: 0/1/-1 (NoResponse/WaitForLocal/WaitForAll)
```

### 4. 消费者实现 (consumer.go)

#### 设计思路
- 消费者组自动管理
- 支持单个和批量消息处理
- 自动重试和错误恢复
- 灵活的偏移量管理

#### 关键特性
```go
type Consumer struct {
    client *KafkaClient
    group  sarama.ConsumerGroup
}

// 处理器类型
- SimpleConsumerHandler - 单条消息处理
- BatchConsumerHandler - 批量消息处理

// 配置选项
- 偏移量策略: earliest/latest/none
- 自动提交: true/false
- 会话超时: 可配置
- 心跳间隔: 可配置
```

## 高级特性

### 1. 消息压缩

支持的压缩算法：
- **None**: 无压缩
- **GZIP**: 通用压缩，压缩率高
- **Snappy**: 快速压缩，CPU占用低
- **LZ4**: 极速压缩，实时性好
- **ZSTD**: 新一代压缩，平衡性能和压缩率

```go
// 配置示例
config.Producer.Compression = "gzip"
```

### 2. 分区策略

支持的分区策略：
- **Hash**: 基于消息Key的哈希值，保证相同Key到同一分区
- **Random**: 随机分区，负载均衡
- **RoundRobin**: 轮询分区，均匀分布

```go
// 配置示例
config.Producer.Partitioner = "hash"
```

### 3. 安全认证

支持的安全机制：
- **SASL/PLAIN**: 用户名密码认证
- **SASL/SCRAM-SHA-256**: 增强认证
- **SASL/SCRAM-SHA-512**: 最高安全认证
- **TLS**: 传输层加密

```go
// 配置示例
config.Security.Enabled = true
config.Security.Mechanism = "SCRAM-SHA-256"
config.Security.TLSEnabled = true
```

## 使用示例

### 生产者示例
```go
// 创建生产者
config := &kafka.ProducerConfig{
    ClientID: "my-producer",
    Brokers:  []string{"localhost:9092"},
    Compression: "gzip",
}

producer, err := kafka.NewProducer(config)
if err != nil {
    panic(err)
}

// 初始化
err = producer.Initialize()
if err != nil {
    panic(err)
}
defer producer.Close()

// 发送消息
_, _, err = producer.SendMessageSync(
    context.Background(),
    "test-topic",
    []byte("message-key"),
    []byte("message-value"),
)
```

### 消费者示例
```go
// 创建消费者
config := &kafka.ConsumerConfig{
    GroupID: "my-consumer-group",
    Brokers: []string{"localhost:9092"},
}

consumer, err := kafka.NewConsumer(config)
if err != nil {
    panic(err)
}

// 初始化
err = consumer.Initialize()
if err != nil {
    panic(err)
}
defer consumer.Close()

// 创建消息处理器
handler := &kafka.SimpleConsumerHandler{
    MessageHandler: func(topic string, partition int32, offset int64, key, value []byte) error {
        println("Received message:", string(value))
        return nil
    },
}

// 订阅主题
err = consumer.Subscribe([]string{"test-topic"}, handler)
```

## 测试策略

### 单元测试覆盖
- ✅ 配置构建和验证
- ✅ 生产者和消费者创建
- ✅ 消息处理器接口
- ✅ 默认配置验证
- ✅ 错误处理逻辑

### 集成测试
- 提供了完整的集成测试框架
- 支持外部Kafka实例测试
- 可选的测试模式（避免依赖外部服务）

## 性能优化

### 1. 批量处理
```go
config.Producer.FlushBytes = 65536    // 64KB
config.Producer.FlushMessages = 100     // 100条消息
config.Consumer.FetchDefault = 1048576 // 1MB
```

### 2. 压缩优化
```go
config.Producer.Compression = "snappy" // 平衡压缩率和性能
```

### 3. 连接池
```go
config.Net.MaxOpenRequests = 5  // 限制并发请求数
```

## 监控和运维

### 1. 健康检查
```go
if client.IsReady() {
    // Kafka客户端正常
}
```

### 2. 错误处理
```go
config.ErrorHandler = func(err error) {
    log.Printf("Kafka error: %v", err)
    // 可以集成监控系统
}
```

### 3. 主题管理
```go
// 创建主题
err := client.CreateTopic("new-topic", 6, 3)

// 列出主题
topics, err := client.ListTopics()
```

## 配置最佳实践

### 生产环境配置
```yaml
kafka:
  brokers:
    - kafka1:9092
    - kafka2:9092
    - kafka3:9092
  
  producer:
    requiredAcks: 1        # 确保持久化
    compression: snappy    # 平衡性能和压缩
    retries: 3
    retryBackoff: 2s
    
  consumer:
    sessionTimeout: 10s
    heartbeatInterval: 3s
    maxProcessingTime: 30s
    
  security:
    enabled: true
    mechanism: SCRAM-SHA-256
    tlsEnabled: true
```

### 开发环境配置
```yaml
kafka:
  brokers: [localhost:9092]
  groupId: gin-admin-dev-group
  
  producer:
    requiredAcks: 0        # 开发环境可以更快
    compression: none       # 方便调试
    
  consumer:
    autoOffset: earliest    # 从头开始消费
```

## 错误处理和恢复

### 自动重试
- 生产者自动重试发送失败的消息
- 消费者自动重连断开的连接
- 可配置的重试次数和退避策略

### 优雅关闭
- 等待正在处理的消息完成
- 正确提交偏移量
- 关闭所有网络连接

## 扩展性

### 插件化设计
- 独立的配置管理
- 可插拔的组件
- 标准化的接口

### 多租户支持
- 支持多个消费者组
- 隔离的主题命名空间
- 独立的配置管理

## 总结

Kafka消息队列插件提供了：
1. **完整的功能**: 生产者、消费者、管理客户端
2. **高性能**: 压缩、批处理、连接池
3. **高可用**: 自动重试、错误恢复、优雅关闭
4. **安全性**: SASL认证、TLS加密
5. **易用性**: 简单的API、丰富的配置
6. **可扩展**: 插件化设计、标准化接口

该实现为企业级应用提供了可靠、高性能的消息队列解决方案。