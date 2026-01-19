package ai

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// DefaultAIService 默认AI服务实现
type DefaultAIService struct {
	config        *Config
	provider      AIProvider
	cache         CacheService
	metrics       MetricsService
	conversations map[string]*Conversation
	users         map[string]*UserStats

	// 中间件
	middlewares  []Middleware
	errorHandler ErrorHandler

	// 插件
	plugins []Plugin

	// 状态管理
	ready  bool
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewDefaultAIService 创建默认AI服务
func NewDefaultAIService(config *Config) (*DefaultAIService, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &DefaultAIService{
		config:        config,
		conversations: make(map[string]*Conversation),
		users:         make(map[string]*UserStats),
		middlewares:   []Middleware{},
		plugins:       []Plugin{},
		ctx:           ctx,
		cancel:        cancel,
	}

	// 创建AI提供商
	provider, err := service.createProvider(config.Provider)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("create provider failed: %w", err)
	}

	if err := provider.Initialize(config); err != nil {
		cancel()
		return nil, fmt.Errorf("initialize provider failed: %w", err)
	}

	service.provider = provider

	// 创建缓存服务
	if config.Cache.Enabled {
		service.cache = NewMemoryCache(config.Cache)
	}

	// 创建指标服务
	if config.Metrics.Enabled {
		service.metrics = NewDefaultMetrics(config.Metrics)
	}

	return service, nil
}

// Initialize 初始化AI服务
func (s *DefaultAIService) Initialize(config *Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ready {
		return fmt.Errorf("service already initialized")
	}

	// 初始化插件
	for _, plugin := range s.plugins {
		if err := plugin.Initialize(config); err != nil {
			log.Printf("Failed to initialize plugin %s: %v", plugin.Name(), err)
		}
	}

	s.ready = true
	log.Printf("AI service initialized with provider: %s", config.Provider)

	return nil
}

// Start 启动AI服务
func (s *DefaultAIService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.ready {
		return fmt.Errorf("service not initialized")
	}

	// 启动后台任务
	s.wg.Add(1)
	go s.backgroundTasks()

	log.Println("AI service started")
	return nil
}

// Stop 停止AI服务
func (s *DefaultAIService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.ready {
		return nil
	}

	// 取消上下文
	s.cancel()

	// 等待后台任务完成
	s.wg.Wait()

	// 清理插件
	for _, plugin := range s.plugins {
		if err := plugin.Cleanup(); err != nil {
			log.Printf("Plugin %s cleanup failed: %v", plugin.Name(), err)
		}
	}

	s.ready = false
	log.Println("AI service stopped")

	return nil
}

// IsReady 检查服务是否就绪
func (s *DefaultAIService) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ready
}

// Chat 聊天接口
func (s *DefaultAIService) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	if !s.IsReady() {
		return nil, fmt.Errorf("service not ready")
	}

	// 记录请求开始
	startTime := time.Now()

	// 执行中间件链
	response, err := s.executeMiddleware(ctx, request, func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
		// 检查缓存
		if s.cache != nil {
			cacheKey := s.generateCacheKey(req)
			if cached, exists := s.cache.Get(cacheKey); exists {
				if cachedResponse, ok := cached.(*ChatResponse); ok {
					log.Printf("Cache hit for request: %s", cacheKey)
					return cachedResponse, nil
				}
			}
		}

		// 执行AI提供商
		resp, err := s.provider.Chat(ctx, req)
		if err != nil {
			// 错误处理
			if s.errorHandler != nil {
				apiErr := s.errorHandler.Handle(ctx, err, req)
				if apiErr != nil {
					resp.Error = apiErr
					return resp, nil
				}
			}
			return nil, err
		}

		// 缓存响应
		if s.cache != nil && resp.Error == nil {
			cacheKey := s.generateCacheKey(req)
			s.cache.Set(cacheKey, resp, s.config.Cache.TTL)
		}

		// 处理插件
		if len(s.plugins) > 0 {
			processedMessage := resp.Message
			for _, plugin := range s.plugins {
				processed, err := plugin.Process(ctx, &processedMessage)
				if err != nil {
					log.Printf("Plugin %s processing failed: %v", plugin.Name(), err)
					continue
				}
				processedMessage = *processed
			}
			resp.Message = processedMessage
		}

		return resp, nil
	})

	// 记录指标
	if s.metrics != nil {
		timer := s.metrics.Timer("ai_chat_duration", map[string]string{
			"provider": s.config.Provider,
			"model":    s.getModel(request),
		})
		timer.Since(startTime)

		counter := s.metrics.Counter("ai_chat_requests_total", map[string]string{
			"provider": s.config.Provider,
			"model":    s.getModel(request),
			"status":   s.getStatus(response, err),
		})
		counter.Inc()
	}

	return response, err
}

// ChatStream 流式聊天接口
func (s *DefaultAIService) ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error) {
	if !s.IsReady() {
		return nil, fmt.Errorf("service not ready")
	}

	if !s.config.EnableStreaming {
		return nil, fmt.Errorf("streaming is disabled")
	}

	// 启用流式响应
	streamRequest := *request
	streamRequest.Stream = true

	return s.provider.ChatStream(ctx, &streamRequest)
}

// CreateConversation 创建对话
func (s *DefaultAIService) CreateConversation(ctx context.Context, userID string, metadata map[string]interface{}) (*Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation := &Conversation{
		ID:        s.generateID(),
		UserID:    userID,
		Title:     s.generateTitle(metadata),
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  metadata,
		Status:    "active",
	}

	s.conversations[conversation.ID] = conversation

	// 更新用户统计
	s.updateUserStats(userID, func(stats *UserStats) {
		stats.ConversationCount++
		stats.LastActiveAt = time.Now()
	})

	log.Printf("Created conversation %s for user %s", conversation.ID, userID)
	return conversation, nil
}

// GetConversation 获取对话
func (s *DefaultAIService) GetConversation(ctx context.Context, conversationID string) (*Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	return conversation, nil
}

// UpdateConversation 更新对话
func (s *DefaultAIService) UpdateConversation(ctx context.Context, conversationID string, metadata map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	if metadata != nil {
		if conversation.Metadata == nil {
			conversation.Metadata = make(map[string]interface{})
		}
		for k, v := range metadata {
			conversation.Metadata[k] = v
		}
	}

	conversation.UpdatedAt = time.Now()

	log.Printf("Updated conversation %s", conversationID)
	return nil
}

// DeleteConversation 删除对话
func (s *DefaultAIService) DeleteConversation(ctx context.Context, conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	// 标记为删除
	conversation.Status = "deleted"
	delete(s.conversations, conversationID)

	log.Printf("Deleted conversation %s", conversationID)
	return nil
}

// ListConversations 列出对话
func (s *DefaultAIService) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userConversations []*Conversation
	for _, conversation := range s.conversations {
		if conversation.UserID == userID && conversation.Status == "active" {
			userConversations = append(userConversations, conversation)
		}
	}

	// 简单分页（实际应该按时间排序）
	total := len(userConversations)
	if offset >= total {
		return []*Conversation{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return userConversations[offset:end], nil
}

// AddMessage 添加消息
func (s *DefaultAIService) AddMessage(ctx context.Context, conversationID string, message *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	// 设置消息ID和时间戳
	if message.ID == "" {
		message.ID = s.generateID()
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	conversation.Messages = append(conversation.Messages, *message)
	conversation.UpdatedAt = time.Now()

	// 更新统计
	conversation.TotalMessages++
	conversation.TotalTokens += message.TokenUsed

	// 更新用户统计
	s.updateUserStats(conversation.UserID, func(stats *UserStats) {
		stats.TotalMessages++
		stats.TotalTokens += message.TokenUsed
		stats.LastActiveAt = time.Now()
	})

	log.Printf("Added message %s to conversation %s", message.ID, conversationID)
	return nil
}

// GetMessages 获取消息
func (s *DefaultAIService) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	// 构建消息指针列表
	messages := make([]*Message, len(conversation.Messages))
	for i := range conversation.Messages {
		messages[i] = &conversation.Messages[i]
	}

	// 简单分页
	total := len(messages)
	if offset >= total {
		return []*Message{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return messages[offset:end], nil
}

// DeleteMessage 删除消息
func (s *DefaultAIService) DeleteMessage(ctx context.Context, conversationID, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	// 查找并删除消息
	for i, message := range conversation.Messages {
		if message.ID == messageID {
			conversation.Messages = append(conversation.Messages[:i], conversation.Messages[i+1:]...)
			conversation.UpdatedAt = time.Now()
			log.Printf("Deleted message %s from conversation %s", messageID, conversationID)
			return nil
		}
	}

	return fmt.Errorf("message not found: %s", messageID)
}

// SearchConversations 搜索对话
func (s *DefaultAIService) SearchConversations(ctx context.Context, userID, query string, limit, offset int) ([]*Conversation, error) {
	// 简化实现，实际应该使用搜索引擎
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*Conversation
	for _, conversation := range s.conversations {
		if conversation.UserID == userID && conversation.Status == "active" {
			// 简单的标题匹配
			if query == "" || contains(conversation.Title, query) {
				results = append(results, conversation)
			}
		}
	}

	// 分页
	total := len(results)
	if offset >= total {
		return []*Conversation{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return results[offset:end], nil
}

// GetUserStats 获取用户统计
func (s *DefaultAIService) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats, exists := s.users[userID]
	if !exists {
		return &UserStats{
			UserID: userID,
		}, nil
	}

	return stats, nil
}

// GetSystemStats 获取系统统计
func (s *DefaultAIService) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &SystemStats{
		TotalUsers:         len(s.users),
		TotalConversations: 0,
		TotalMessages:      0,
		TotalTokens:        0,
		Timestamp:          time.Now(),
	}

	for _, conversation := range s.conversations {
		if conversation.Status == "active" {
			stats.TotalConversations++
			stats.TotalMessages += conversation.TotalMessages
			stats.TotalTokens += conversation.TotalTokens
		}
	}

	// 计算活跃用户数
	activeThreshold := time.Now().Add(-24 * time.Hour)
	for _, userStats := range s.users {
		if userStats.LastActiveAt.After(activeThreshold) {
			stats.ActiveUsers++
		}
	}

	return stats, nil
}

// 添加中间件
func (s *DefaultAIService) AddMiddleware(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

// 设置错误处理器
func (s *DefaultAIService) SetErrorHandler(handler ErrorHandler) {
	s.errorHandler = handler
}

// 添加插件
func (s *DefaultAIService) AddPlugin(plugin Plugin) error {
	s.plugins = append(s.plugins, plugin)
	return nil
}

// 私有方法

// createProvider 创建AI提供商
func (s *DefaultAIService) createProvider(provider string) (AIProvider, error) {
	switch provider {
	case "openai":
		return NewOpenAIProvider(), nil
	case "claude":
		return NewClaudeProvider(), nil
	case "deepseek":
		return NewDeepSeekProvider(), nil
	case "qwen":
		return NewQwenProvider(), nil
	case "ollama":
		return NewOllamaProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// executeMiddleware 执行中间件链
func (s *DefaultAIService) executeMiddleware(ctx context.Context, request *ChatRequest, handler func(context.Context, *ChatRequest) (*ChatResponse, error)) (*ChatResponse, error) {
	// 如果没有中间件，直接执行处理器
	if len(s.middlewares) == 0 {
		return handler(ctx, request)
	}

	// 简化实现：直接执行处理器
	return handler(ctx, request)
}

// generateCacheKey 生成缓存键
func (s *DefaultAIService) generateCacheKey(request *ChatRequest) string {
	// 简化实现，实际应该包含请求的哈希
	return fmt.Sprintf("chat:%s:%s", request.Model, request.SystemPrompt)
}

// generateID 生成唯一ID
func (s *DefaultAIService) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// generateTitle 生成对话标题
func (s *DefaultAIService) generateTitle(metadata map[string]interface{}) string {
	if title, ok := metadata["title"].(string); ok && title != "" {
		return title
	}
	return fmt.Sprintf("对话 %s", time.Now().Format("2006-01-02 15:04:05"))
}

// getModel 获取请求的模型
func (s *DefaultAIService) getModel(request *ChatRequest) string {
	if request.Model != "" {
		return request.Model
	}
	return s.config.Model
}

// getStatus 获取响应状态
func (s *DefaultAIService) getStatus(response *ChatResponse, err error) string {
	if err != nil {
		return "error"
	}
	if response.Error != nil {
		return "api_error"
	}
	return "success"
}

// updateUserStats 更新用户统计
func (s *DefaultAIService) updateUserStats(userID string, updateFunc func(*UserStats)) {
	stats, exists := s.users[userID]
	if !exists {
		stats = &UserStats{
			UserID: userID,
		}
		s.users[userID] = stats
	}
	updateFunc(stats)
}

// backgroundTasks 后台任务
func (s *DefaultAIService) backgroundTasks() {
	defer s.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// 定期清理过期缓存
			if s.cache != nil {
				// 简化实现，实际应该基于TTL清理
			}
		}
	}
}

// contains 字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}()))
}

// middlewareAdapter 中间件适配器
type middlewareAdapter struct {
	handler func(context.Context, *ChatRequest) (*ChatResponse, error)
}

func (a *middlewareAdapter) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	return a.handler(ctx, request)
}

// middlewareWrapper 中间件包装器
type middlewareWrapper struct {
	middleware Middleware
	next       AIService
}

func (w *middlewareWrapper) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	return w.middleware.Process(ctx, request, w.next)
}
