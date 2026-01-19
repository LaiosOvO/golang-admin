package ai

import (
	"context"
	"time"
)

// Message 对话消息
type Message struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"` // user, assistant, system, function
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	TokenUsed int                    `json:"tokenUsed"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// 多模态支持
	ImageURL string `json:"imageUrl,omitempty"`
	AudioURL string `json:"audioUrl,omitempty"`

	// 函数调用
	FunctionCall *FunctionCall `json:"functionCall,omitempty"`

	// 流式响应
	Streaming    bool   `json:"streaming,omitempty"`
	DeltaContent string `json:"deltaContent,omitempty"`
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    interface{}            `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// Conversation 对话
type Conversation struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"userId"`
	Title     string                 `json:"title"`
	Messages  []Message              `json:"messages"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// 统计信息
	TotalMessages int           `json:"totalMessages"`
	TotalTokens   int           `json:"totalTokens"`
	Duration      time.Duration `json:"duration"`

	// 状态
	Status string `json:"status"` // active, archived, deleted
}

// ChatRequest 聊天请求
type ChatRequest struct {
	ConversationID string    `json:"conversationId,omitempty"`
	Messages       []Message `json:"messages"`
	Model          string    `json:"model,omitempty"`
	MaxTokens      int       `json:"maxTokens,omitempty"`
	Temperature    float64   `json:"temperature,omitempty"`
	TopP           float64   `json:"topP,omitempty"`
	Stream         bool      `json:"stream,omitempty"`

	// 系统配置
	SystemPrompt string               `json:"systemPrompt,omitempty"`
	Functions    []FunctionDefinition `json:"functions,omitempty"`

	// 元数据
	UserID    string                 `json:"userId,omitempty"`
	SessionID string                 `json:"sessionId,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string  `json:"id"`
	Message Message `json:"message"`
	Usage   Usage   `json:"usage"`
	Model   string  `json:"model"`
	Finish  string  `json:"finish"` // stop, length, function_call, content_filter

	// 流式响应
	Stream bool   `json:"stream,omitempty"`
	Delta  string `json:"delta,omitempty"`
	Done   bool   `json:"done,omitempty"`

	// 元数据
	RequestID string                 `json:"requestId,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// 错误信息
	Error *APIError `json:"error,omitempty"`
}

// Usage Token使用情况
type Usage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`

	// 成本计算
	Cost float64 `json:"cost,omitempty"`
}

// APIError API错误
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`

	// HTTP状态码
	StatusCode int `json:"statusCode,omitempty"`

	// 重试信息
	Retryable  bool          `json:"retryable,omitempty"`
	RetryAfter time.Duration `json:"retryAfter,omitempty"`
}

// AIProvider AI提供商接口
type AIProvider interface {
	// 基础方法
	Initialize(config *Config) error
	Validate() error

	// 聊天方法
	Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error)

	// 工具方法
	GetModels() []string
	GetModelInfo(model string) (*ModelInfo, error)

	// 统计方法
	GetUsage(ctx context.Context, startTime, endTime time.Time) (*Usage, error)
	GetCost(ctx context.Context, startTime, endTime time.Time) (float64, error)

	// 健康检查
	HealthCheck(ctx context.Context) error
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	MaxTokens   int       `json:"maxTokens"`
	InputCost   float64   `json:"inputCost"`
	OutputCost  float64   `json:"outputCost"`
	Features    []string  `json:"features"` // chat, image, voice, function_calling
	ContextSize int       `json:"contextSize"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AIService AI服务接口
type AIService interface {
	// 服务管理
	Initialize(config *Config) error
	Start() error
	Stop() error
	IsReady() bool

	// 对话管理
	CreateConversation(ctx context.Context, userID string, metadata map[string]interface{}) (*Conversation, error)
	GetConversation(ctx context.Context, conversationID string) (*Conversation, error)
	UpdateConversation(ctx context.Context, conversationID string, metadata map[string]interface{}) error
	DeleteConversation(ctx context.Context, conversationID string) error
	ListConversations(ctx context.Context, userID string, limit, offset int) ([]*Conversation, error)

	// 消息管理
	AddMessage(ctx context.Context, conversationID string, message *Message) error
	GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error)
	DeleteMessage(ctx context.Context, conversationID, messageID string) error

	// 聊天接口
	Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error)

	// 历史记录
	SearchConversations(ctx context.Context, userID, query string, limit, offset int) ([]*Conversation, error)

	// 统计接口
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	GetSystemStats(ctx context.Context) (*SystemStats, error)
}

// UserStats 用户统计
type UserStats struct {
	UserID            string        `json:"userId"`
	TotalMessages     int           `json:"totalMessages"`
	TotalTokens       int           `json:"totalTokens"`
	TotalCost         float64       `json:"totalCost"`
	ConversationCount int           `json:"conversationCount"`
	AvgResponseTime   time.Duration `json:"avgResponseTime"`
	LastActiveAt      time.Time     `json:"lastActiveAt"`
}

// SystemStats 系统统计
type SystemStats struct {
	TotalUsers         int           `json:"totalUsers"`
	TotalConversations int           `json:"totalConversations"`
	TotalMessages      int           `json:"totalMessages"`
	TotalTokens        int           `json:"totalTokens"`
	TotalCost          float64       `json:"totalCost"`
	ActiveUsers        int           `json:"activeUsers"`
	ActiveRequests     int           `json:"activeRequests"`
	AvgResponseTime    time.Duration `json:"avgResponseTime"`
	Timestamp          time.Time     `json:"timestamp"`
}

// CacheService 缓存服务接口
type CacheService interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Stats() CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	HitRatio    float64 `json:"hitRatio"`
	Count       int     `json:"count"`
	MemoryUsage int64   `json:"memoryUsage"`
}

// MetricsService 指标服务接口
type MetricsService interface {
	Counter(name string, tags map[string]string) Counter
	Histogram(name string, tags map[string]string) Histogram
	Gauge(name string, tags map[string]string) Gauge
	Timer(name string, tags map[string]string) Timer
}

// Counter 计数器
type Counter interface {
	Inc()
	Add(float64)
}

// Histogram 直方图
type Histogram interface {
	Observe(float64)
}

// Gauge 仪表盘
type Gauge interface {
	Set(float64)
	Inc()
	Dec()
}

// Timer 计时器
type Timer interface {
	Time(func())
	Since(time.Time) time.Duration
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	Handle(ctx context.Context, err error, request *ChatRequest) *APIError
}

// Middleware 中间件接口
type Middleware interface {
	Process(ctx context.Context, request *ChatRequest, next AIService) (*ChatResponse, error)
}

// Plugin 插件接口
type Plugin interface {
	Name() string
	Version() string
	Initialize(config *Config) error
	Process(ctx context.Context, message *Message) (*Message, error)
	Cleanup() error
}
