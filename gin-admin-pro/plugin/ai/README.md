# AI Plugin

AI插件提供了完整的智能对话功能，支持多种大模型提供商、流式响应、对话历史管理和插件扩展。

## 功能特性

- ✅ 多种AI提供商支持（OpenAI、Claude、DeepSeek、Qwen、Ollama）
- ✅ 流式响应和非流式响应
- ✅ 对话历史管理和持久化
- ✅ 多模态输入支持（文本、图像、音频）
- ✅ 函数调用能力
- ✅ 智能缓存机制
- ✅ 中间件和插件系统
- ✅ 成本控制和监控
- ✅ 并发控制和重试机制
- ✅ 完整的错误处理

## 快速开始

### 基础配置

```yaml
# config/config.yaml
ai:
  enabled: true
  provider: "openai"
  apiKey: "sk-your-api-key-here"
  baseUrl: "https://api.openai.com/v1"
  model: "gpt-3.5-turbo"
  maxTokens: 2048
  temperature: 0.7
  topP: 1.0
  systemPrompt: "You are a helpful AI assistant."
  timeout: 30s
  maxHistoryLength: 20
  enableHistory: true
  enableStreaming: true
  maxRetries: 3
  retryDelay: 2s
  
  # 成本控制
  dailyTokenLimit: 100000
  enableCostLimit: false
  
  # 功能开关
  enableFunctionCalling: false
  enableImageInput: false
  enableVoiceInput: false
  
  # 监控配置
  metrics:
    enabled: true
    prefix: "ai"
    tags: ["gin-admin"]
    
  # 缓存配置
  cache:
    enabled: true
    ttl: 1h
    maxSize: 1000
```

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "gin-admin-pro/plugin/ai"
)

func main() {
    // 创建AI服务
    config := ai.DefaultConfig()
    config.Provider = "openai"
    config.APIKey = "your-api-key"
    
    service, err := ai.NewDefaultAIService(config)
    if err != nil {
        panic(err)
    }
    
    // 初始化和启动
    err = service.Initialize(config)
    if err != nil {
        panic(err)
    }
    
    err = service.Start()
    if err != nil {
        panic(err)
    }
    defer service.Stop()
    
    // 创建对话
    ctx := context.Background()
    conversation, err := service.CreateConversation(ctx, "user-123", map[string]interface{}{
        "title": "测试对话",
    })
    if err != nil {
        panic(err)
    }
    
    // 发送消息
    request := &ai.ChatRequest{
        ConversationID: conversation.ID,
        Messages: []ai.Message{
            {Role: "user", Content: "你好，请介绍一下你自己"},
        },
        MaxTokens: 1000,
        Temperature: 0.7,
    }
    
    response, err := service.Chat(ctx, request)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("AI回复: %s\n", response.Message.Content)
}
```

## 提供商配置

### OpenAI

```yaml
ai:
  provider: "openai"
  apiKey: "sk-..."
  baseUrl: "https://api.openai.com/v1"
  model: "gpt-3.5-turbo"
  
  openai:
    organization: "org-..."
    projectId: "proj_..."
    maxTokens: 4096
```

### Claude

```yaml
ai:
  provider: "claude"
  apiKey: "sk-ant-..."
  baseUrl: "https://api.anthropic.com"
  model: "claude-3-sonnet-20240229"
  
  claude:
    apiVersion: "2023-06-01"
    maxTokens: 4096
```

### DeepSeek

```yaml
ai:
  provider: "deepseek"
  apiKey: "sk-..."
  baseUrl: "https://api.deepseek.com"
  model: "deepseek-chat"
  
  deepseek:
    model: "deepseek-chat"
    maxTokens: 4096
```

### Qwen

```yaml
ai:
  provider: "qwen"
  apiKey: "sk-..."
  baseUrl: "https://dashscope.aliyuncs.com/api/v1"
  model: "qwen-turbo"
  
  qwen:
    dashScopeApiKey: "sk-..."
    model: "qwen-turbo"
```

### Ollama

```yaml
ai:
  provider: "ollama"
  # apiKey not required for local models
  
  ollama:
    host: "localhost"
    port: 11434
    model: "llama2"
    keepAlive: true
    numPredict: 128
    numCtx: 2048
    repeatPenalty: 1.1
```

## 高级功能

### 流式响应

```go
// 启用流式响应
request := &ai.ChatRequest{
    ConversationID: conversation.ID,
    Messages: []ai.Message{
        {Role: "user", Content: "请写一个长篇故事"},
    },
    Stream: true,
}

// 接收流式响应
stream, err := service.ChatStream(ctx, request)
if err != nil {
    panic(err)
}

for response := range stream {
    if response.Done {
        break
    }
    if response.Delta != "" {
        fmt.Printf("增量内容: %s\n", response.Delta)
    }
}
```

### 函数调用

```go
// 定义函数
functions := []ai.FunctionDefinition{
    {
        Name: "get_weather",
        Description: "获取指定城市的天气信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "city": map[string]interface{}{
                    "type": "string",
                    "description": "城市名称",
                },
                "units": map[string]interface{}{
                    "type": "string",
                    "enum": []string{"celsius", "fahrenheit"},
                    "description": "温度单位",
                },
            },
            "required": []string{"city"},
        },
    },
}

request := &ai.ChatRequest{
    ConversationID: conversation.ID,
    Messages: []ai.Message{
        {Role: "user", Content: "北京今天天气怎么样？"},
    },
    Functions: functions,
}

response, err := service.Chat(ctx, request)
if err != nil {
    panic(err)
}

// 处理函数调用结果
if response.Message.FunctionCall != nil {
    fmt.Printf("AI调用了函数: %s\n", response.Message.FunctionCall.Name)
    fmt.Printf("参数: %+v\n", response.Message.FunctionCall.Arguments)
}
```

### 多模态输入

```go
// 图像输入
request := &ai.ChatRequest{
    Messages: []ai.Message{
        {
            Role: "user",
            Content: "这张图片里有什么？",
            ImageURL: "https://example.com/image.jpg",
        },
    },
}

// 音频输入
request = &ai.ChatRequest{
    Messages: []ai.Message{
        {
            Role: "user",
            Content: "请转录这段音频",
            AudioURL: "https://example.com/audio.wav",
        },
    },
}
```

### 中间件系统

```go
// 自定义中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Process(ctx context.Context, request *ai.ChatRequest, next ai.AIService) (*ai.ChatResponse, error) {
    log.Printf("开始处理请求: %+v", request)
    
    response, err := next.Chat(ctx, request)
    
    log.Printf("请求处理完成: %v, 错误: %v", response != nil, err)
    return response, err
}

// 添加中间件
service.AddMiddleware(&LoggingMiddleware{})
```

### 插件系统

```go
// 自定义插件
type ContentFilterPlugin struct{}

func (p *ContentFilterPlugin) Name() string {
    return "content_filter"
}

func (p *ContentFilterPlugin) Version() string {
    return "1.0.0"
}

func (p *ContentFilterPlugin) Initialize(config *ai.Config) error {
    return nil
}

func (p *ContentFilterPlugin) Process(ctx context.Context, message *ai.Message) (*ai.Message, error) {
    // 内容过滤逻辑
    if containsSensitiveContent(message.Content) {
        return &ai.Message{
            Role:    "assistant",
            Content: "抱歉，我无法回答这个问题。",
        }, nil
    }
    return message, nil
}

func (p *ContentFilterPlugin) Cleanup() error {
    return nil
}

// 添加插件
service.AddPlugin(&ContentFilterPlugin{})
```

## 对话管理

### 创建和查询对话

```go
// 创建新对话
conversation, err := service.CreateConversation(ctx, "user-123", map[string]interface{}{
    "title": "客服对话",
    "category": "customer_service",
})

// 列出用户对话
conversations, err := service.ListConversations(ctx, "user-123", 20, 0)

// 搜索对话
results, err := service.SearchConversations(ctx, "user-123", "客服", 10, 0)
```

### 消息管理

```go
// 添加消息
message := &ai.Message{
    Role:    "user",
    Content: "你好，请问有什么可以帮助您的？",
}
err := service.AddMessage(ctx, conversation.ID, message)

// 获取消息历史
messages, err := service.GetMessages(ctx, conversation.ID, 50, 0)

// 删除消息
err := service.DeleteMessage(ctx, conversation.ID, "message-id")
```

### 用户统计

```go
// 获取用户统计
stats, err := service.GetUserStats(ctx, "user-123")
if err == nil {
    fmt.Printf("总消息数: %d\n", stats.TotalMessages)
    fmt.Printf("总Token数: %d\n", stats.TotalTokens)
    fmt.Printf("总成本: $%.2f\n", stats.TotalCost)
    fmt.Printf("对话数量: %d\n", stats.ConversationCount)
    fmt.Printf("最后活跃时间: %v\n", stats.LastActiveAt)
}
```

## 监控和指标

### 指标收集

```go
// 启用指标收集
config := ai.DefaultConfig()
config.Metrics.Enabled = true
config.Metrics.Prefix = "my_app"
config.Metrics.Tags = []string{"service", "ai"}

service, _ := ai.NewDefaultAIService(config)
```

### 成本控制

```yaml
ai:
  enableCostLimit: true
  dailyTokenLimit: 10000
  maxConcurrentRequests: 5
```

### 缓存配置

```yaml
ai:
  cache:
    enabled: true
    ttl: 2h
    maxSize: 5000
```

## 错误处理

### 重试机制

```go
// 配置重试
config := ai.DefaultConfig()
config.MaxRetries = 5
config.RetryDelay = time.Second * 3

// 检查错误是否可重试
response, err := service.Chat(ctx, request)
if err != nil {
    if apiErr, ok := err.(*ai.APIError); ok {
        if apiErr.Retryable {
            fmt.Printf("可重试错误: %s\n", apiErr.Message)
        }
    }
}
```

### 错误分类

- **认证错误**: API密钥无效、权限不足
- **限流错误**: 请求频率超限、配额不足
- **模型错误**: 模型不存在、参数无效
- **网络错误**: 连接超时、服务不可用
- **内容错误**: 内容被拒绝、安全检查失败

## 性能优化

### 并发控制

```yaml
ai:
  maxConcurrentRequests: 10
  requestTimeout: 60s
```

### 缓存策略

- **相同问题缓存**: 相同用户问题缓存30分钟
- **模型结果缓存**: 复杂推理结果缓存2小时
- **分层缓存**: 热点问题永久缓存，常规问题短期缓存

### 流量控制

- **请求队列**: 使用队列控制并发请求数量
- **优先级队列**: VIP用户请求优先处理
- **熔断机制**: 错误率过高时暂停请求

## 安全配置

### API密钥管理

```yaml
ai:
  # 使用环境变量
  apiKey: "${AI_API_KEY}"
  
  # 或使用配置管理服务
  # apiKey: "vault://secret/ai-api-key"
```

### 内容过滤

```go
// 自定义内容过滤插件
type SecurityPlugin struct{}

func (p *SecurityPlugin) Process(ctx context.Context, message *ai.Message) (*ai.Message, error) {
    // 检查敏感词
    if containsSensitiveWords(message.Content) {
        return &ai.Message{
            Role:    "assistant",
            Content: "抱歉，我不能回答这个问题。",
        }, nil
    }
    
    // 检查恶意链接
    if containsMaliciousLinks(message.Content) {
        return nil, fmt.Errorf("message blocked: contains malicious links")
    }
    
    return message, nil
}
```

### 数据脱敏

```go
// 数据脱敏中间件
type DataMaskingMiddleware struct{}

func (m *DataMaskingMiddleware) Process(ctx context.Context, request *ai.ChatRequest, next ai.AIService) (*ai.ChatResponse, error) {
    // 脱敏用户输入中的敏感信息
    for i := range request.Messages {
        request.Messages[i].Content = maskSensitiveData(request.Messages[i].Content)
    }
    
    response, err := next.Chat(ctx, request)
    
    // 脱敏AI回复中的敏感信息
    if response != nil && response.Error == nil {
        response.Message.Content = maskSensitiveData(response.Message.Content)
    }
    
    return response, err
}
```

## 部署建议

### 开发环境

```yaml
ai:
  enabled: true
  provider: "openai"
  maxTokens: 100
  temperature: 0.1
  cache:
    enabled: false
  metrics:
    enabled: false
```

### 生产环境

```yaml
ai:
  enabled: true
  provider: "openai"
  maxTokens: 4096
  temperature: 0.7
  maxRetries: 3
  enableCostLimit: true
  dailyTokenLimit: 1000000
  cache:
    enabled: true
    ttl: 2h
    maxSize: 10000
  metrics:
    enabled: true
    prefix: "prod_ai"
    tags: ["production", "ai"]
```

### 高可用部署

```yaml
# 多提供商配置，自动故障转移
ai:
  primaryProvider: "openai"
  fallbackProvider: "claude"
  healthCheckInterval: 30s
```

## 故障排除

### 常见问题

1. **API密钥错误**
   ```
   错误: API key is required
   解决: 检查配置中的API密钥是否正确
   ```

2. **网络连接错误**
   ```
   错误: connection timeout
   解决: 检查网络连接和代理设置
   ```

3. **模型不存在**
   ```
   错误: model not found
   解决: 检查模型名称是否正确
   ```

4. **Token超限**
   ```
   错误: rate limit exceeded
   解决: 等待重置或增加配额
   ```

### 调试模式

```yaml
ai:
  debug: true
  logLevel: "debug"
  logRequests: true
  logResponses: true
```

### 健康检查

```go
// 检查AI服务健康状态
if !service.IsReady() {
    log.Println("AI service not ready")
}

// 检查提供商健康状态
err := provider.HealthCheck(ctx)
if err != nil {
    log.Printf("Provider health check failed: %v", err)
}
```

## 扩展开发

### 自定义提供商

```go
// 实现AIProvider接口
type CustomProvider struct{}

func (p *CustomProvider) Initialize(config *ai.Config) error {
    // 初始化自定义提供商
}

func (p *CustomProvider) Chat(ctx context.Context, request *ai.ChatRequest) (*ai.ChatResponse, error) {
    // 实现聊天逻辑
}

// 注册自定义提供商
func init() {
    ai.RegisterProvider("custom", func() ai.AIProvider {
        return &CustomProvider{}
    })
}
```

### 自定义缓存

```go
// 实现CacheService接口
type RedisCache struct{}

func (c *RedisCache) Get(key string) (interface{}, bool) {
    // Redis获取逻辑
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    // Redis设置逻辑
}
```

## 最佳实践

1. **配置管理**
   - 使用环境变量存储敏感信息
   - 分离开发和生产环境配置
   - 定期轮换API密钥

2. **性能优化**
   - 合理设置缓存TTL
   - 使用连接池
   - 监控关键指标

3. **安全防护**
   - 实施内容过滤
   - 启用请求限流
   - 记录访问日志

4. **错误处理**
   - 分类错误类型
   - 实现优雅降级
   - 提供用户友好的错误信息

## 总结

AI插件提供了完整的智能对话解决方案：

1. **多提供商支持**: 支持主流AI服务提供商
2. **灵活的配置**: 丰富的配置选项满足不同需求
3. **高性能**: 缓存、并发控制、流式响应
4. **高可用**: 重试机制、故障转移
5. **可扩展**: 插件系统、中间件支持
6. **生产就绪**: 监控、安全、成本控制

该插件为企业级应用提供了可靠、高性能的AI对话功能。