package ai

import (
	"context"
	"fmt"
	"time"
)

// StubProvider 存根提供商实现
type StubProvider struct {
	name string
}

// NewClaudeProvider 创建Claude提供商
func NewClaudeProvider() AIProvider {
	return &StubProvider{name: "claude"}
}

// NewDeepSeekProvider 创建DeepSeek提供商
func NewDeepSeekProvider() AIProvider {
	return &StubProvider{name: "deepseek"}
}

// NewQwenProvider 创建Qwen提供商
func NewQwenProvider() AIProvider {
	return &StubProvider{name: "qwen"}
}

// NewOllamaProvider 创建Ollama提供商
func NewOllamaProvider() AIProvider {
	return &StubProvider{name: "ollama"}
}

// Initialize 初始化
func (p *StubProvider) Initialize(config *Config) error {
	return nil
}

// Validate 验证
func (p *StubProvider) Validate() error {
	return nil
}

// Chat 聊天
func (p *StubProvider) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	return nil, fmt.Errorf("provider %s not implemented", p.name)
}

// ChatStream 流式聊天
func (p *StubProvider) ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error) {
	return nil, fmt.Errorf("provider %s not implemented", p.name)
}

// GetModels 获取模型
func (p *StubProvider) GetModels() []string {
	return []string{p.name + "-model"}
}

// GetModelInfo 获取模型信息
func (p *StubProvider) GetModelInfo(model string) (*ModelInfo, error) {
	return &ModelInfo{
		ID:          model,
		Name:        model,
		Provider:    p.name,
		MaxTokens:   4096,
		InputCost:   0.01,
		OutputCost:  0.02,
		Features:    []string{"chat"},
		ContextSize: 4096,
		CreatedAt:   time.Now(),
	}, nil
}

// GetUsage 获取使用统计
func (p *StubProvider) GetUsage(ctx context.Context, startTime, endTime time.Time) (*Usage, error) {
	return &Usage{
		PromptTokens:     1000,
		CompletionTokens: 500,
		TotalTokens:      1500,
		Cost:             0.01,
	}, nil
}

// GetCost 获取成本
func (p *StubProvider) GetCost(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return 0.01, nil
}

// HealthCheck 健康检查
func (p *StubProvider) HealthCheck(ctx context.Context) error {
	return fmt.Errorf("provider %s not implemented", p.name)
}
