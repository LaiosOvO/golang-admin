package ai

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, "openai", config.Provider)
	assert.Equal(t, "gpt-3.5-turbo", config.Model)
	assert.Equal(t, 2048, config.MaxTokens)
	assert.Equal(t, 0.7, config.Temperature)
	assert.Equal(t, 1.0, config.TopP)
	assert.True(t, config.EnableHistory)
	assert.Equal(t, 20, config.MaxHistoryLength)
}

func TestConfigValidate(t *testing.T) {
	// Test valid config
	config := DefaultConfig()
	config.APIKey = "test-key" // Add API key to make config valid
	err := config.Validate()
	assert.NoError(t, err)

	// Test missing provider
	config.Provider = ""
	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provider is required")

	// Test missing API key for non-local provider
	config.Provider = "openai"
	config.APIKey = ""
	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API key is required")

	// Test invalid temperature
	config.APIKey = "test-key"
	config.Temperature = -1
	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "temperature must be between 0 and 2")

	// Test invalid maxTokens
	config.Temperature = 0.7
	config.MaxTokens = 0
	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maxTokens must be greater than 0")
}

func TestConfigIsLocalProvider(t *testing.T) {
	config := DefaultConfig()

	// Test non-local provider
	config.Provider = "openai"
	assert.False(t, config.IsLocalProvider())

	// Test local provider
	config.Provider = "ollama"
	assert.True(t, config.IsLocalProvider())

	config.Provider = "local"
	assert.True(t, config.IsLocalProvider())
}

func TestConfigGetEffectiveMaxTokens(t *testing.T) {
	config := DefaultConfig()

	// Test with config MaxTokens
	config.MaxTokens = 4096
	assert.Equal(t, 4096, config.GetEffectiveMaxTokens())

	// Test without config MaxTokens
	config.MaxTokens = 0
	assert.Greater(t, config.GetEffectiveMaxTokens(), 0)

	// Test with provider-specific config
	config.Provider = "claude"
	config.Claude.MaxTokens = 8192
	assert.Equal(t, 8192, config.GetEffectiveMaxTokens())
}

func TestConfigGetEffectiveBaseURL(t *testing.T) {
	config := DefaultConfig()

	// Test with custom base URL
	config.BaseURL = "https://custom.api.com/v1"
	assert.Equal(t, "https://custom.api.com/v1", config.GetEffectiveBaseURL())

	// Test without custom base URL
	config.BaseURL = ""
	config.Provider = "openai"
	assert.Equal(t, "https://api.openai.com/v1", config.GetEffectiveBaseURL())

	config.Provider = "claude"
	assert.Equal(t, "https://api.anthropic.com", config.GetEffectiveBaseURL())

	config.Provider = "ollama"
	config.Ollama.Host = "localhost"
	config.Ollama.Port = 11434
	assert.Equal(t, "http://localhost:11434", config.GetEffectiveBaseURL())
}

func TestNewDefaultAIService(t *testing.T) {
	config := DefaultConfig()
	config.Provider = "openai"
	config.APIKey = "test-key"

	service, err := NewDefaultAIService(config)
	require.NoError(t, err)
	assert.NotNil(t, service)
	// Note: service is not ready until Initialize is called
	err = service.Initialize(config)
	assert.NoError(t, err)
	assert.True(t, service.IsReady())
}

func TestNewDefaultAIServiceInvalidConfig(t *testing.T) {
	// Test nil config
	service, err := NewDefaultAIService(nil)
	assert.Error(t, err)
	assert.Nil(t, service)

	// Test invalid config
	config := DefaultConfig()
	config.Provider = "invalid"
	service, err = NewDefaultAIService(config)
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestDefaultAIServiceCreateConversation(t *testing.T) {
	service := createTestService(t)

	ctx := context.Background()
	userID := "test-user"
	metadata := map[string]interface{}{
		"title": "Test Conversation",
	}

	conversation, err := service.CreateConversation(ctx, userID, metadata)
	require.NoError(t, err)
	assert.NotNil(t, conversation)
	assert.NotEmpty(t, conversation.ID)
	assert.Equal(t, userID, conversation.UserID)
	assert.Equal(t, "Test Conversation", conversation.Title)
	assert.Equal(t, "active", conversation.Status)
	assert.Empty(t, conversation.Messages)
	assert.False(t, conversation.CreatedAt.IsZero())
	assert.False(t, conversation.UpdatedAt.IsZero())
}

func TestDefaultAIServiceGetConversation(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Create conversation
	conversation, err := service.CreateConversation(ctx, "test-user", nil)
	require.NoError(t, err)

	// Get conversation
	retrieved, err := service.GetConversation(ctx, conversation.ID)
	require.NoError(t, err)
	assert.Equal(t, conversation.ID, retrieved.ID)
	assert.Equal(t, conversation.UserID, retrieved.UserID)
}

func TestDefaultAIServiceGetConversationNotFound(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Get non-existent conversation
	conversation, err := service.GetConversation(ctx, "non-existent")
	assert.Error(t, err)
	assert.Nil(t, conversation)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestDefaultAIServiceListConversations(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create multiple conversations
	for i := 0; i < 5; i++ {
		_, err := service.CreateConversation(ctx, userID, nil)
		require.NoError(t, err)
	}

	// List conversations
	conversations, err := service.ListConversations(ctx, userID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, conversations, 5)
}

func TestDefaultAIServiceAddMessage(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Create conversation
	conversation, err := service.CreateConversation(ctx, "test-user", nil)
	require.NoError(t, err)

	// Add message
	message := &Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	err = service.AddMessage(ctx, conversation.ID, message)
	require.NoError(t, err)

	// Verify message was added
	retrieved, err := service.GetConversation(ctx, conversation.ID)
	require.NoError(t, err)
	assert.Len(t, retrieved.Messages, 1)
	assert.Equal(t, message.Role, retrieved.Messages[0].Role)
	assert.Equal(t, message.Content, retrieved.Messages[0].Content)
}

func TestDefaultAIServiceAddMessageAutoID(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Create conversation
	conversation, err := service.CreateConversation(ctx, "test-user", nil)
	require.NoError(t, err)

	// Add message without ID
	message := &Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	err = service.AddMessage(ctx, conversation.ID, message)
	require.NoError(t, err)

	// Verify ID was auto-generated
	retrieved, err := service.GetConversation(ctx, conversation.ID)
	require.NoError(t, err)
	assert.Len(t, retrieved.Messages, 1)
	assert.NotEmpty(t, retrieved.Messages[0].ID)
}

func TestDefaultAIServiceGetMessages(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Create conversation
	conversation, err := service.CreateConversation(ctx, "test-user", nil)
	require.NoError(t, err)

	// Add multiple messages
	for i := 0; i < 5; i++ {
		message := &Message{
			Role:    "user",
			Content: fmt.Sprintf("Message %d", i),
		}
		err := service.AddMessage(ctx, conversation.ID, message)
		require.NoError(t, err)
	}

	// Get messages
	messages, err := service.GetMessages(ctx, conversation.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, messages, 5)
}

func TestDefaultAIServiceGetMessagesPagination(t *testing.T) {
	service := createTestService(t)
	ctx := context.Background()

	// Create conversation
	conversation, err := service.CreateConversation(ctx, "test-user", nil)
	require.NoError(t, err)

	// Add messages
	for i := 0; i < 5; i++ {
		message := &Message{
			Role:    "user",
			Content: fmt.Sprintf("Message %d", i),
		}
		err := service.AddMessage(ctx, conversation.ID, message)
		require.NoError(t, err)
	}

	// Get first page
	messages, err := service.GetMessages(ctx, conversation.ID, 2, 0)
	require.NoError(t, err)
	assert.Len(t, messages, 2)

	// Get second page
	messages, err = service.GetMessages(ctx, conversation.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, messages, 2)

	// Get third page
	messages, err = service.GetMessages(ctx, conversation.ID, 2, 4)
	require.NoError(t, err)
	assert.Len(t, messages, 1)
}

func TestOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider()
	assert.NotNil(t, provider)
}

func TestOpenAIProviderInitialize(t *testing.T) {
	provider := NewOpenAIProvider()

	config := DefaultConfig()
	config.APIKey = "test-key"

	err := provider.Initialize(config)
	assert.NoError(t, err)

	// Validate provider
	err = provider.Validate()
	assert.NoError(t, err)
}

func TestOpenAIProviderInitializeNoKey(t *testing.T) {
	provider := NewOpenAIProvider()

	config := DefaultConfig()
	config.APIKey = ""

	// Initialize should fail if API key is required
	err := provider.Initialize(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API key is required")
}

func TestOpenAIProviderGetModels(t *testing.T) {
	provider := NewOpenAIProvider()
	models := provider.GetModels()
	assert.NotEmpty(t, models)
	assert.Contains(t, models, "gpt-3.5-turbo")
	assert.Contains(t, models, "gpt-4")
}

func TestOpenAIProviderGetModelInfo(t *testing.T) {
	provider := NewOpenAIProvider()

	// Test existing model
	info, err := provider.GetModelInfo("gpt-3.5-turbo")
	require.NoError(t, err)
	assert.Equal(t, "gpt-3.5-turbo", info.ID)
	assert.Equal(t, "openai", info.Provider)
	assert.Greater(t, info.MaxTokens, 0)
	assert.Greater(t, info.InputCost, 0.0)
	assert.Greater(t, info.OutputCost, 0.0)
	assert.NotEmpty(t, info.Features)

	// Test non-existing model
	info, err = provider.GetModelInfo("non-existent-model")
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestStubProviders(t *testing.T) {
	providers := []AIProvider{
		NewClaudeProvider(),
		NewDeepSeekProvider(),
		NewQwenProvider(),
		NewOllamaProvider(),
	}

	for _, provider := range providers {
		assert.NotNil(t, provider)

		// Test initialize
		err := provider.Initialize(DefaultConfig())
		assert.NoError(t, err)

		// Test validate
		err = provider.Validate()
		assert.NoError(t, err)

		// Test get models
		models := provider.GetModels()
		assert.NotEmpty(t, models)

		// Test health check
		err = provider.HealthCheck(context.Background())
		assert.Error(t, err) // Should return not implemented error
		assert.Contains(t, err.Error(), "not implemented")
	}
}

func TestMessage(t *testing.T) {
	message := &Message{
		ID:        "test-id",
		Role:      "user",
		Content:   "Hello, world!",
		Timestamp: time.Now(),
		TokenUsed: 10,
		ImageURL:  "http://example.com/image.jpg",
		FunctionCall: &FunctionCall{
			Name: "test_function",
			Arguments: map[string]interface{}{
				"param1": "value1",
			},
		},
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	assert.Equal(t, "test-id", message.ID)
	assert.Equal(t, "user", message.Role)
	assert.Equal(t, "Hello, world!", message.Content)
	assert.Equal(t, 10, message.TokenUsed)
	assert.Equal(t, "http://example.com/image.jpg", message.ImageURL)
	assert.Equal(t, "test_function", message.FunctionCall.Name)
	assert.Equal(t, "value1", message.FunctionCall.Arguments["param1"])
	assert.Equal(t, "test", message.Metadata["source"])
}

func TestConversation(t *testing.T) {
	now := time.Now()
	conversation := &Conversation{
		ID:        "test-conversation",
		UserID:    "test-user",
		Title:     "Test Conversation",
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
		Metadata: map[string]interface{}{
			"source": "test",
		},
		TotalMessages: 5,
		TotalTokens:   100,
		Status:        "active",
	}

	assert.Equal(t, "test-conversation", conversation.ID)
	assert.Equal(t, "test-user", conversation.UserID)
	assert.Equal(t, "Test Conversation", conversation.Title)
	assert.Empty(t, conversation.Messages)
	assert.Equal(t, now, conversation.CreatedAt)
	assert.Equal(t, now, conversation.UpdatedAt)
	assert.Equal(t, "test", conversation.Metadata["source"])
	assert.Equal(t, 5, conversation.TotalMessages)
	assert.Equal(t, 100, conversation.TotalTokens)
	assert.Equal(t, "active", conversation.Status)
}

func TestChatRequest(t *testing.T) {
	request := &ChatRequest{
		ConversationID: "test-conversation",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Model:        "gpt-3.5-turbo",
		MaxTokens:    100,
		Temperature:  0.7,
		TopP:         1.0,
		Stream:       false,
		SystemPrompt: "You are a helpful assistant.",
		Functions: []FunctionDefinition{
			{
				Name:        "test_function",
				Description: "A test function",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"param1": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		UserID:    "test-user",
		SessionID: "test-session",
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	assert.Equal(t, "test-conversation", request.ConversationID)
	assert.Len(t, request.Messages, 1)
	assert.Equal(t, "user", request.Messages[0].Role)
	assert.Equal(t, "Hello", request.Messages[0].Content)
	assert.Equal(t, "gpt-3.5-turbo", request.Model)
	assert.Equal(t, 100, request.MaxTokens)
	assert.Equal(t, 0.7, request.Temperature)
	assert.Equal(t, 1.0, request.TopP)
	assert.False(t, request.Stream)
	assert.Equal(t, "You are a helpful assistant.", request.SystemPrompt)
	assert.Len(t, request.Functions, 1)
	assert.Equal(t, "test_function", request.Functions[0].Name)
	assert.Equal(t, "test-user", request.UserID)
	assert.Equal(t, "test-session", request.SessionID)
	assert.Equal(t, "test", request.Metadata["source"])
}

func TestChatResponse(t *testing.T) {
	now := time.Now()
	response := &ChatResponse{
		ID: "test-response",
		Message: Message{
			Role:      "assistant",
			Content:   "Hello! How can I help you?",
			Timestamp: now,
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
			Cost:             0.01,
		},
		Model:     "gpt-3.5-turbo",
		Finish:    "stop",
		Stream:    false,
		Done:      true,
		RequestID: "test-request",
		Timestamp: now,
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	assert.Equal(t, "test-response", response.ID)
	assert.Equal(t, "assistant", response.Message.Role)
	assert.Equal(t, "Hello! How can I help you?", response.Message.Content)
	assert.Equal(t, now, response.Message.Timestamp)
	assert.Equal(t, 10, response.Usage.PromptTokens)
	assert.Equal(t, 20, response.Usage.CompletionTokens)
	assert.Equal(t, 30, response.Usage.TotalTokens)
	assert.Equal(t, 0.01, response.Usage.Cost)
	assert.Equal(t, "gpt-3.5-turbo", response.Model)
	assert.Equal(t, "stop", response.Finish)
	assert.False(t, response.Stream)
	assert.True(t, response.Done)
	assert.Equal(t, "test-request", response.RequestID)
	assert.Equal(t, now, response.Timestamp)
	assert.Equal(t, "test", response.Metadata["source"])
}

func TestFunctionCall(t *testing.T) {
	fnCall := &FunctionCall{
		Name: "test_function",
		Arguments: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
		Result: "success",
		Error:  "",
	}

	assert.Equal(t, "test_function", fnCall.Name)
	assert.Equal(t, "value1", fnCall.Arguments["param1"])
	assert.Equal(t, 42, fnCall.Arguments["param2"])
	assert.Equal(t, "success", fnCall.Result)
	assert.Equal(t, "", fnCall.Error)
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		Code:       "rate_limit_exceeded",
		Message:    "Rate limit exceeded",
		Type:       "rate_limit",
		Param:      "requests",
		StatusCode: 429,
		Retryable:  true,
		RetryAfter: time.Minute * 5,
	}

	assert.Equal(t, "rate_limit_exceeded", err.Code)
	assert.Equal(t, "Rate limit exceeded", err.Message)
	assert.Equal(t, "rate_limit", err.Type)
	assert.Equal(t, "requests", err.Param)
	assert.Equal(t, 429, err.StatusCode)
	assert.True(t, err.Retryable)
	assert.Equal(t, time.Minute*5, err.RetryAfter)
}

func TestUsage(t *testing.T) {
	usage := &Usage{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		Cost:             0.01,
	}

	assert.Equal(t, 100, usage.PromptTokens)
	assert.Equal(t, 50, usage.CompletionTokens)
	assert.Equal(t, 150, usage.TotalTokens)
	assert.Equal(t, 0.01, usage.Cost)
}

func TestModelInfo(t *testing.T) {
	now := time.Now()
	info := &ModelInfo{
		ID:          "gpt-3.5-turbo",
		Name:        "GPT-3.5 Turbo",
		Provider:    "openai",
		MaxTokens:   4096,
		InputCost:   0.0015,
		OutputCost:  0.002,
		Features:    []string{"chat", "function_calling"},
		ContextSize: 4096,
		CreatedAt:   now,
	}

	assert.Equal(t, "gpt-3.5-turbo", info.ID)
	assert.Equal(t, "GPT-3.5 Turbo", info.Name)
	assert.Equal(t, "openai", info.Provider)
	assert.Equal(t, 4096, info.MaxTokens)
	assert.Equal(t, 0.0015, info.InputCost)
	assert.Equal(t, 0.002, info.OutputCost)
	assert.Contains(t, info.Features, "chat")
	assert.Contains(t, info.Features, "function_calling")
	assert.Equal(t, 4096, info.ContextSize)
	assert.Equal(t, now, info.CreatedAt)
}

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(struct {
		Enabled bool          `yaml:"enabled" mapstructure:"enabled"`
		TTL     time.Duration `yaml:"ttl" mapstructure:"ttl"`
		MaxSize int           `yaml:"maxSize" mapstructure:"maxSize"`
	}{
		Enabled: true,
		TTL:     time.Hour,
		MaxSize: 100,
	})

	assert.NotNil(t, cache)

	// Test Set and Get
	err := cache.Set("test-key", "test-value", time.Minute)
	assert.NoError(t, err)

	value, exists := cache.Get("test-key")
	assert.True(t, exists)
	assert.Equal(t, "test-value", value)

	// Test Get non-existent
	_, exists = cache.Get("non-existent")
	assert.False(t, exists)

	// Test Delete
	err = cache.Delete("test-key")
	assert.NoError(t, err)

	_, exists = cache.Get("test-key")
	assert.False(t, exists)

	// Test Stats
	stats := cache.Stats()
	assert.NotNil(t, stats)
}

func TestDefaultMetrics(t *testing.T) {
	metrics := NewDefaultMetrics(struct {
		Enabled bool     `yaml:"enabled" mapstructure:"enabled"`
		Prefix  string   `yaml:"prefix" mapstructure:"prefix"`
		Tags    []string `yaml:"tags" mapstructure:"tags"`
	}{
		Enabled: true,
		Prefix:  "test",
		Tags:    []string{"service"},
	})

	assert.NotNil(t, metrics)

	// Test metric creation
	counter := metrics.Counter("test_counter", map[string]string{"env": "test"})
	assert.NotNil(t, counter)

	histogram := metrics.Histogram("test_histogram", map[string]string{"env": "test"})
	assert.NotNil(t, histogram)

	gauge := metrics.Gauge("test_gauge", map[string]string{"env": "test"})
	assert.NotNil(t, gauge)

	timer := metrics.Timer("test_timer", map[string]string{"env": "test"})
	assert.NotNil(t, timer)
}

// Helper functions

func createTestService(t *testing.T) *DefaultAIService {
	config := DefaultConfig()
	config.Provider = "openai"
	config.APIKey = "test-key"

	service, err := NewDefaultAIService(config)
	require.NoError(t, err)

	return service
}
