package ai

import (
	"fmt"
	"time"
)

// Config AI配置结构
type Config struct {
	// 基础配置
	Enabled  bool   `yaml:"enabled" mapstructure:"enabled"`
	Provider string `yaml:"provider" mapstructure:"provider"` // openai, claude, deepseek, qwen, ollama
	APIKey   string `yaml:"apiKey" mapstructure:"apiKey"`
	BaseURL  string `yaml:"baseUrl" mapstructure:"baseUrl"`
	Model    string `yaml:"model" mapstructure:"model"`

	// 对话配置
	MaxTokens   int     `yaml:"maxTokens" mapstructure:"maxTokens"`
	Temperature float64 `yaml:"temperature" mapstructure:"temperature"`
	TopP        float64 `yaml:"topP" mapstructure:"topP"`

	// 系统配置
	SystemPrompt string        `yaml:"systemPrompt" mapstructure:"systemPrompt"`
	Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout"`

	// 历史记录配置
	MaxHistoryLength int  `yaml:"maxHistoryLength" mapstructure:"maxHistoryLength"`
	EnableHistory    bool `yaml:"enableHistory" mapstructure:"enableHistory"`

	// 并发配置
	MaxConcurrentRequests int           `yaml:"maxConcurrentRequests" mapstructure:"maxConcurrentRequests"`
	RequestTimeout        time.Duration `yaml:"requestTimeout" mapstructure:"requestTimeout"`

	// 流式响应配置
	EnableStreaming bool `yaml:"enableStreaming" mapstructure:"enableStreaming"`

	// 重试配置
	MaxRetries int           `yaml:"maxRetries" mapstructure:"maxRetries"`
	RetryDelay time.Duration `yaml:"retryDelay" mapstructure:"retryDelay"`

	// 成本控制
	DailyTokenLimit int  `yaml:"dailyTokenLimit" mapstructure:"dailyTokenLimit"`
	EnableCostLimit bool `yaml:"enableCostLimit" mapstructure:"enableCostLimit"`

	// 功能配置
	EnableFunctionCalling bool `yaml:"enableFunctionCalling" mapstructure:"enableFunctionCalling"`
	EnableImageInput      bool `yaml:"enableImageInput" mapstructure:"enableImageInput"`
	EnableVoiceInput      bool `yaml:"enableVoiceInput" mapstructure:"enableVoiceInput"`

	// 提供商特定配置
	OpenAI   *OpenAIConfig   `yaml:"openai" mapstructure:"openai"`
	Claude   *ClaudeConfig   `yaml:"claude" mapstructure:"claude"`
	DeepSeek *DeepSeekConfig `yaml:"deepseek" mapstructure:"deepseek"`
	Qwen     *QwenConfig     `yaml:"qwen" mapstructure:"qwen"`
	Ollama   *OllamaConfig   `yaml:"ollama" mapstructure:"ollama"`

	// 监控配置
	Metrics struct {
		Enabled bool     `yaml:"enabled" mapstructure:"enabled"`
		Prefix  string   `yaml:"prefix" mapstructure:"prefix"`
		Tags    []string `yaml:"tags" mapstructure:"tags"`
	} `yaml:"metrics" mapstructure:"metrics"`

	// 缓存配置
	Cache struct {
		Enabled bool          `yaml:"enabled" mapstructure:"enabled"`
		TTL     time.Duration `yaml:"ttl" mapstructure:"ttl"`
		MaxSize int           `yaml:"maxSize" mapstructure:"maxSize"`
	} `yaml:"cache" mapstructure:"cache"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	Organization string `yaml:"organization" mapstructure:"organization"`
	ProjectID    string `yaml:"projectId" mapstructure:"projectId"`
	MaxTokens    int    `yaml:"maxTokens" mapstructure:"maxTokens"`
}

// ClaudeConfig Claude配置
type ClaudeConfig struct {
	APIVersion string `yaml:"apiVersion" mapstructure:"apiVersion"`
	MaxTokens  int    `yaml:"maxTokens" mapstructure:"maxTokens"`
}

// DeepSeekConfig DeepSeek配置
type DeepSeekConfig struct {
	Model     string `yaml:"model" mapstructure:"model"`
	MaxTokens int    `yaml:"maxTokens" mapstructure:"maxTokens"`
}

// QwenConfig Qwen配置
type QwenConfig struct {
	DashScopeAPIKey string `yaml:"dashScopeApiKey" mapstructure:"dashScopeApiKey"`
	Model           string `yaml:"model" mapstructure:"model"`
}

// OllamaConfig Ollama配置
type OllamaConfig struct {
	Host          string  `yaml:"host" mapstructure:"host"`
	Port          int     `yaml:"port" mapstructure:"port"`
	Model         string  `yaml:"model" mapstructure:"model"`
	KeepAlive     bool    `yaml:"keepAlive" mapstructure:"keepAlive"`
	NumPredict    int     `yaml:"numPredict" mapstructure:"numPredict"`
	NumCtx        int     `yaml:"numCtx" mapstructure:"numCtx"`
	RepeatPenalty float64 `yaml:"repeatPenalty" mapstructure:"repeatPenalty"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:               true,
		Provider:              "openai",
		APIKey:                "",
		BaseURL:               "https://api.openai.com/v1",
		Model:                 "gpt-3.5-turbo",
		MaxTokens:             2048,
		Temperature:           0.7,
		TopP:                  1.0,
		SystemPrompt:          "You are a helpful AI assistant.",
		Timeout:               time.Second * 30,
		MaxHistoryLength:      20,
		EnableHistory:         true,
		MaxConcurrentRequests: 10,
		RequestTimeout:        time.Second * 60,
		EnableStreaming:       true,
		MaxRetries:            3,
		RetryDelay:            time.Second * 2,
		DailyTokenLimit:       100000,
		EnableCostLimit:       false,
		EnableFunctionCalling: false,
		EnableImageInput:      false,
		EnableVoiceInput:      false,

		// 提供商默认配置
		OpenAI: &OpenAIConfig{
			Organization: "",
			ProjectID:    "",
		},
		Claude: &ClaudeConfig{
			APIVersion: "2023-06-01",
			MaxTokens:  4096,
		},
		DeepSeek: &DeepSeekConfig{
			Model:     "deepseek-chat",
			MaxTokens: 4096,
		},
		Qwen: &QwenConfig{
			Model: "qwen-turbo",
		},
		Ollama: &OllamaConfig{
			Host:          "localhost",
			Port:          11434,
			Model:         "llama2",
			KeepAlive:     true,
			NumPredict:    128,
			NumCtx:        2048,
			RepeatPenalty: 1.1,
		},

		// 监控默认配置
		Metrics: struct {
			Enabled bool     `yaml:"enabled" mapstructure:"enabled"`
			Prefix  string   `yaml:"prefix" mapstructure:"prefix"`
			Tags    []string `yaml:"tags" mapstructure:"tags"`
		}{
			Enabled: false,
			Prefix:  "ai",
			Tags:    []string{"gin-admin"},
		},

		// 缓存默认配置
		Cache: struct {
			Enabled bool          `yaml:"enabled" mapstructure:"enabled"`
			TTL     time.Duration `yaml:"ttl" mapstructure:"ttl"`
			MaxSize int           `yaml:"maxSize" mapstructure:"maxSize"`
		}{
			Enabled: false,
			TTL:     time.Hour,
			MaxSize: 1000,
		},
	}
}

// IsLocalProvider 检查是否为本地提供商
func (c *Config) IsLocalProvider() bool {
	return c.Provider == "ollama" || c.Provider == "local"
}

// IsStreamingEnabled 检查是否启用流式响应
func (c *Config) IsStreamingEnabled() bool {
	return c.EnableStreaming
}

// GetEffectiveMaxTokens 获取有效的最大token数
func (c *Config) GetEffectiveMaxTokens() int {
	if c.MaxTokens > 0 {
		return c.MaxTokens
	}

	switch c.Provider {
	case "openai":
		if c.OpenAI != nil && c.OpenAI.MaxTokens > 0 {
			return c.OpenAI.MaxTokens
		}
		return 4096
	case "claude":
		if c.Claude != nil && c.Claude.MaxTokens > 0 {
			return c.Claude.MaxTokens
		}
		return 4096
	case "deepseek":
		if c.DeepSeek != nil && c.DeepSeek.MaxTokens > 0 {
			return c.DeepSeek.MaxTokens
		}
		return 4096
	default:
		return 2048
	}
}

// GetEffectiveBaseURL 获取有效的API地址
func (c *Config) GetEffectiveBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}

	switch c.Provider {
	case "openai":
		return "https://api.openai.com/v1"
	case "claude":
		return "https://api.anthropic.com"
	case "deepseek":
		return "https://api.deepseek.com"
	case "qwen":
		return "https://dashscope.aliyuncs.com/api/v1"
	case "ollama":
		if c.Ollama != nil {
			return fmt.Sprintf("http://%s:%d", c.Ollama.Host, c.Ollama.Port)
		}
		return "http://localhost:11434"
	default:
		return ""
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.Provider == "" {
		return fmt.Errorf("AI provider is required")
	}

	if c.Provider != "ollama" && c.APIKey == "" {
		return fmt.Errorf("API key is required for provider %s", c.Provider)
	}

	if c.MaxTokens <= 0 {
		return fmt.Errorf("maxTokens must be greater than 0")
	}

	if c.Temperature < 0 || c.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if c.TopP < 0 || c.TopP > 1 {
		return fmt.Errorf("topP must be between 0 and 1")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	if c.MaxConcurrentRequests <= 0 {
		return fmt.Errorf("maxConcurrentRequests must be greater than 0")
	}

	return nil
}
