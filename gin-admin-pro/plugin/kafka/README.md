# Kafka Plugin

Kafka插件提供了完整的企业级消息队列解决方案，支持高可用、高性能的消息生产和消费。

## 功能特性

- ✅ 同步/异步消息发送
- ✅ 消费者组管理
- ✅ 自动重试机制
- ✅ 消息压缩支持（GZIP、Snappy、LZ4、ZSTD）
- ✅ 多种分区策略
- ✅ 安全认证（SASL、TLS）
- ✅ 批量消息处理
- ✅ 消息监控和指标
- ✅ 优雅关闭和错误处理

## 快速开始

### 基础配置

```yaml
# config/config.yaml
kafka:
  brokers:
    - localhost:9092
  groupId: gin-admin-group
  autoCommit: true
  autoOffset: latest
  
  producer:
    requiredAcks: 1
    compression: none
    partitioner: hash
    returnSuccesses: true
    returnErrors: true
    
  consumer:
    fetchMin: 1
    fetchDefault: 1048576  # 1MB
    fetchMax: 10485760     # 10MB
    sessionTimeout: 10s
    heartbeatInterval: 3s
    
  security:
    enabled: false
    mechanism: PLAIN
    username: ""
    password: ""
    tlsEnabled: false
```

### 生产者示例

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/kafka"
)

func main() {
    // 创建生产者
    config := &kafka.ProducerConfig{
        ClientID: "my-producer",
        Brokers:  []string{"localhost:9092"},
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
    if err != nil {
        panic(err)
    }
}
```

### 消费者示例

```go
package main

import (
    "gin-admin-pro/plugin/kafka"
)

func main() {
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
    if err != nil {
        panic(err)
    }
    
    // 保持程序运行
    select {}
}
```

### 批量消费者示例

```go
package main

import (
    "time"
    "github.com/IBM/sarama"
    "gin-admin-pro/plugin/kafka"
)

func main() {
    config := &kafka.ConsumerConfig{
        GroupID: "batch-consumer-group",
        Brokers: []string{"localhost:9092"},
    }
    
    consumer, err := kafka.NewConsumer(config)
    if err != nil {
        panic(err)
    }
    
    err = consumer.Initialize()
    if err != nil {
        panic(err)
    }
    defer consumer.Close()
    
    // 批量消息处理器
    handler := &kafka.BatchConsumerHandler{
        BatchSize:    100,
        BatchTimeout: time.Second * 5,
        MessageHandler: func(messages []*sarama.ConsumerMessage) error {
            println("Processing batch of", len(messages), "messages")
            // 批量处理逻辑
            return nil
        },
    }
    
    err = consumer.Subscribe([]string{"test-topic"}, handler)
    if err != nil {
        panic(err)
    }
    
    select {}
}
```

## 完整客户端示例

```go
package main

import (
    "context"
    "gin-admin-pro/plugin/kafka"
)

func main() {
    // 创建Kafka客户端
    config := kafka.DefaultConfig()
    config.Brokers = []string{"localhost:9092"}
    
    client, err := kafka.NewKafkaClient(config)
    if err != nil {
        panic(err)
    }
    
    // 初始化
    err = client.Initialize()
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 发送消息
    messages := []*kafka.Message{
        {Key: []byte("key1"), Value: []byte("value1")},
        {Key: []byte("key2"), Value: []byte("value2")},
    }
    
    err = client.SendBatchMessages(context.Background(), "test-topic", messages)
    if err != nil {
        panic(err)
    }
    
    // 创建主题
    err = client.CreateTopic("new-topic", 3, 2)
    if err != nil {
        panic(err)
    }
    
    // 列出主题
    topics, err := client.ListTopics()
    if err != nil {
        panic(err)
    }
    
    for topic := range topics {
        println("Topic:", topic)
    }
}
```

## 配置参数详解

### 基础配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `brokers` | []string | ["localhost:9092"] | Kafka broker地址列表 |
| `groupId` | string | "gin-admin-group" | 消费者组ID |
| `autoCommit` | bool | true | 是否自动提交偏移量 |
| `autoOffset` | string | "latest" | 初始偏移量策略：earliest/latest/none |

### 生产者配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `requiredAcks` | int16 | 1 | 确认级别：0=无确认,1=本地确认,-1=全部确认 |
| `compression` | string | "none" | 压缩算法：none/gzip/snappy/lz4/zstd |
| `partitioner` | string | "hash" | 分区策略：hash/random/roundrobin |
| `maxMessageBytes` | int | 1000000 | 最大消息大小（字节） |
| `returnSuccesses` | bool | true | 是否返回成功消息 |
| `returnErrors` | bool | true | 是否返回错误消息 |

### 消费者配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `sessionTimeout` | duration | 10s | 会话超时时间 |
| `heartbeatInterval` | duration | 3s | 心跳间隔 |
| `fetchMin` | int32 | 1 | 最小获取字节数 |
| `fetchDefault` | int32 | 1048576 | 默认获取字节数 |
| `fetchMax` | int32 | 10485760 | 最大获取字节数 |
| `maxProcessingTime` | duration | 30s | 最大处理时间 |

### 安全配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | false | 是否启用安全认证 |
| `mechanism` | string | "PLAIN" | SASL机制：PLAIN/SCRAM-SHA-256/SCRAM-SHA-512 |
| `username` | string | "" | 用户名 |
| `password` | string | "" | 密码 |
| `tlsEnabled` | bool | false | 是否启用TLS |

## 最佳实践

### 1. 生产者优化

```go
config := &kafka.ProducerConfig{
    // 启用压缩以减少网络传输
    Compression: "gzip",
    
    // 设置适当的批量大小
    FlushBytes: 65536, // 64KB
    FlushMessages: 100,
    
    // 启用重试机制
    Retries: 3,
    RetryBackoff: time.Second * 2,
    
    // 确保消息持久化
    RequiredAcks: 1, // WaitForLocal
}
```

### 2. 消费者优化

```go
config := &kafka.ConsumerConfig{
    // 设置适当的获取大小
    FetchMin: 1,
    FetchDefault: 1048576, // 1MB
    FetchMax: 10485760,    // 10MB
    
    // 调整会话和心跳时间
    SessionTimeout: time.Second * 10,
    HeartbeatInterval: time.Second * 3,
    
    // 设置处理超时
    MaxProcessingTime: time.Second * 30,
}
```

### 3. 错误处理

```go
config := &kafka.ProducerConfig{
    ErrorHandler: func(err error) {
        log.Printf("Producer error: %v", err)
        // 可以在这里添加监控告警
    },
}

consumerConfig := &kafka.ConsumerConfig{
    ErrorHandler: func(err error) {
        log.Printf("Consumer error: %v", err)
        // 可以在这里添加重试逻辑
    },
}
```

## 监控指标

启用监控指标：

```yaml
kafka:
  metrics:
    enabled: true
    prefix: "gin-admin-kafka"
    tags: ["service", "gin-admin"]
```

## 常见问题

### Q: 如何处理消息重复消费？

A: 消费者应该实现幂等性逻辑，可以通过消息ID或业务唯一标识来去重。

### Q: 如何确保消息不丢失？

A: 设置`requiredAcks: 1`或`-1`，并确保生产者正确处理错误重试。

### Q: 如何优化性能？

A: 启用压缩、调整批量大小、合理设置分区数。

## 运维指南

### 主题管理

```go
// 创建主题
err := client.CreateTopic("topic-name", 6, 3)

// 删除主题
err := client.DeleteTopic("topic-name")

// 列出所有主题
topics, err := client.ListTopics()
```

### 监控检查

```go
// 检查客户端状态
if client.IsReady() {
    println("Kafka client is ready")
}

// 获取配置信息
config := client.GetConfig()
```

## 依赖

- [IBM Sarama](https://github.com/IBM/sarama) - Go Kafka客户端库

## 许可证

本插件遵循项目主许可证。