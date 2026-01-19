package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider OpenAI提供商实现
type OpenAIProvider struct {
	client *http.Client
	config *Config
}

// NewOpenAIProvider 创建OpenAI提供商
func NewOpenAIProvider() AIProvider {
	return &OpenAIProvider{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Initialize 初始化提供商
func (p *OpenAIProvider) Initialize(config *Config) error {
	p.config = config

	// 设置客户端超时
	p.client.Timeout = config.Timeout

	// 验证配置
	if config.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}

	return nil
}

// Validate 验证提供商
func (p *OpenAIProvider) Validate() error {
	if p.config == nil {
		return fmt.Errorf("config not initialized")
	}
	if p.config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

// Chat 聊天接口
func (p *OpenAIProvider) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	// 构建请求
	reqBody := p.buildChatRequest(request)

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 发送请求
	url := fmt.Sprintf("%s/chat/completions", p.config.GetEffectiveBaseURL())
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))
	if p.config.OpenAI.Organization != "" {
		req.Header.Set("OpenAI-Organization", p.config.OpenAI.Organization)
	}

	// 执行请求
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	return p.parseResponse(resp, respBody)
}

// ChatStream 流式聊天接口
func (p *OpenAIProvider) ChatStream(ctx context.Context, request *ChatRequest) (<-chan *ChatResponse, error) {
	// 启用流式响应
	streamRequest := *request
	streamRequest.Stream = true

	// 构建请求
	reqBody := p.buildChatRequest(&streamRequest)

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/chat/completions", p.config.GetEffectiveBaseURL())
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))
	if p.config.OpenAI.Organization != "" {
		req.Header.Set("OpenAI-Organization", p.config.OpenAI.Organization)
	}

	// 执行请求
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 创建响应通道
	respChan := make(chan *ChatResponse, 100)

	// 启动goroutine处理流式响应
	go func() {
		defer close(respChan)
		defer resp.Body.Close()

		p.handleStreamResponse(resp.Body, respChan)
	}()

	return respChan, nil
}

// GetModels 获取可用模型
func (p *OpenAIProvider) GetModels() []string {
	return []string{
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
		"gpt-4",
		"gpt-4-32k",
		"gpt-4-turbo-preview",
		"gpt-4-vision-preview",
	}
}

// GetModelInfo 获取模型信息
func (p *OpenAIProvider) GetModelInfo(model string) (*ModelInfo, error) {
	models := map[string]*ModelInfo{
		"gpt-3.5-turbo": {
			ID:          "gpt-3.5-turbo",
			Name:        "GPT-3.5 Turbo",
			Provider:    "openai",
			MaxTokens:   4096,
			InputCost:   0.0015,
			OutputCost:  0.002,
			Features:    []string{"chat"},
			ContextSize: 4096,
			CreatedAt:   time.Now(),
		},
		"gpt-4": {
			ID:          "gpt-4",
			Name:        "GPT-4",
			Provider:    "openai",
			MaxTokens:   8192,
			InputCost:   0.03,
			OutputCost:  0.06,
			Features:    []string{"chat"},
			ContextSize: 8192,
			CreatedAt:   time.Now(),
		},
		"gpt-4-vision-preview": {
			ID:          "gpt-4-vision-preview",
			Name:        "GPT-4 Vision Preview",
			Provider:    "openai",
			MaxTokens:   4096,
			InputCost:   0.03,
			OutputCost:  0.06,
			Features:    []string{"chat", "image"},
			ContextSize: 4096,
			CreatedAt:   time.Now(),
		},
	}

	info, exists := models[model]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", model)
	}

	return info, nil
}

// GetUsage 获取使用统计
func (p *OpenAIProvider) GetUsage(ctx context.Context, startTime, endTime time.Time) (*Usage, error) {
	// 简化实现，实际应该调用OpenAI的usage API
	return &Usage{
		PromptTokens:     1000,
		CompletionTokens: 500,
		TotalTokens:      1500,
		Cost:             0.01,
	}, nil
}

// GetCost 获取成本统计
func (p *OpenAIProvider) GetCost(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	usage, err := p.GetUsage(ctx, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return usage.Cost, nil
}

// HealthCheck 健康检查
func (p *OpenAIProvider) HealthCheck(ctx context.Context) error {
	// 简单的模型检查
	models := p.GetModels()
	if len(models) == 0 {
		return fmt.Errorf("no available models")
	}

	// 检查API连接
	testReq := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "test"},
		},
		MaxTokens: 1,
	}

	_, err := p.Chat(ctx, testReq)
	if err != nil {
		return fmt.Errorf("API health check failed: %w", err)
	}

	return nil
}

// 私有方法

// buildChatRequest 构建聊天请求
func (p *OpenAIProvider) buildChatRequest(request *ChatRequest) map[string]interface{} {
	req := map[string]interface{}{
		"model":    p.getModel(request.Model),
		"messages": p.convertMessages(request.Messages),
		"stream":   request.Stream,
	}

	// 添加可选参数
	if request.MaxTokens > 0 {
		req["max_tokens"] = request.MaxTokens
	} else if p.config.MaxTokens > 0 {
		req["max_tokens"] = p.config.MaxTokens
	}

	if request.Temperature > 0 {
		req["temperature"] = request.Temperature
	} else {
		req["temperature"] = p.config.Temperature
	}

	if request.TopP > 0 {
		req["top_p"] = request.TopP
	} else {
		req["top_p"] = p.config.TopP
	}

	if request.SystemPrompt != "" {
		req["system"] = request.SystemPrompt
	} else if p.config.SystemPrompt != "" {
		req["system"] = p.config.SystemPrompt
	}

	// 函数调用
	if len(request.Functions) > 0 {
		req["functions"] = request.Functions
		req["function_call"] = "auto"
	}

	return req
}

// convertMessages 转换消息格式
func (p *OpenAIProvider) convertMessages(messages []Message) []map[string]interface{} {
	converted := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		convertedMsg := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}

		// 添加函数调用
		if msg.FunctionCall != nil {
			convertedMsg["function_call"] = map[string]interface{}{
				"name":      msg.FunctionCall.Name,
				"arguments": msg.FunctionCall.Arguments,
			}
		}

		converted[i] = convertedMsg
	}
	return converted
}

// getModel 获取模型
func (p *OpenAIProvider) getModel(model string) string {
	if model != "" {
		return model
	}
	return p.config.Model
}

// parseResponse 解析响应
func (p *OpenAIProvider) parseResponse(resp *http.Response, body []byte) (*ChatResponse, error) {
	// 检查HTTP状态
	if resp.StatusCode != http.StatusOK {
		return &ChatResponse{
			Error: &APIError{
				Code:       fmt.Sprintf("HTTP_%d", resp.StatusCode),
				Message:    string(body),
				StatusCode: resp.StatusCode,
				Retryable:  resp.StatusCode >= 500,
			},
		}, nil
	}

	// 解析JSON响应
	var rawResp map[string]interface{}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	// 检查API错误
	if errObj, exists := rawResp["error"]; exists {
		return &ChatResponse{
			Error: p.parseAPIError(errObj),
		}, nil
	}

	// 解析成功响应
	return p.parseSuccessResponse(rawResp)
}

// parseAPIError 解析API错误
func (p *OpenAIProvider) parseAPIError(errObj interface{}) *APIError {
	errMap, ok := errObj.(map[string]interface{})
	if !ok {
		return &APIError{
			Code:    "unknown",
			Message: fmt.Sprintf("%v", errObj),
		}
	}

	code, _ := errMap["code"].(string)
	message, _ := errMap["message"].(string)

	return &APIError{
		Code:      code,
		Message:   message,
		Retryable: p.isRetryableError(code),
	}
}

// parseSuccessResponse 解析成功响应
func (p *OpenAIProvider) parseSuccessResponse(resp map[string]interface{}) (*ChatResponse, error) {
	choices, exists := resp["choices"].([]interface{})
	if !exists || len(choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})

	// 构建消息内容
	content, _ := message["content"].(string)
	role, _ := message["role"].(string)

	// 解析使用情况
	usage := &Usage{}
	if usageObj, exists := resp["usage"]; exists {
		usageMap := usageObj.(map[string]interface{})
		if promptTokens, ok := usageMap["prompt_tokens"].(float64); ok {
			usage.PromptTokens = int(promptTokens)
		}
		if completionTokens, ok := usageMap["completion_tokens"].(float64); ok {
			usage.CompletionTokens = int(completionTokens)
		}
		if totalTokens, ok := usageMap["total_tokens"].(float64); ok {
			usage.TotalTokens = int(totalTokens)
		}
	}

	// 解析函数调用
	var functionCall *FunctionCall
	if fnCall, exists := message["function_call"]; exists {
		if fnMap, ok := fnCall.(map[string]interface{}); ok {
			functionCall = &FunctionCall{
				Name:      fnMap["name"].(string),
				Arguments: fnMap["arguments"].(map[string]interface{}),
			}
		}
	}

	return &ChatResponse{
		ID: fmt.Sprintf("chat_%d", time.Now().Unix()),
		Message: Message{
			Role:         role,
			Content:      content,
			FunctionCall: functionCall,
			Timestamp:    time.Now(),
		},
		Usage:  *usage,
		Model:  p.config.Model,
		Finish: "stop",
		Stream: false,
		Done:   true,
	}, nil
}

// handleStreamResponse 处理流式响应
func (p *OpenAIProvider) handleStreamResponse(body io.ReadCloser, respChan chan<- *ChatResponse) {
	scanner := bufio.NewScanner(body)

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				// 发送错误
				respChan <- &ChatResponse{
					Error: &APIError{
						Code:    "stream_error",
						Message: err.Error(),
					},
					Stream: true,
					Done:   true,
				}
				return
			}
			break
		}

		line := scanner.Text()

		// 跳过空行
		if line == "" {
			continue
		}

		// 解析SSE格式
		if len(line) > 6 && line[:6] == "data: " {
			data := line[6:]

			if data == "[DONE]" {
				// 流结束
				respChan <- &ChatResponse{
					Stream: true,
					Done:   true,
				}
				return
			}

			// 解析JSON数据
			var streamResp map[string]interface{}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue // 跳过错误数据
			}

			// 转换为ChatResponse
			if resp := p.convertStreamResponse(streamResp); resp != nil {
				respChan <- resp
			}
		}
	}
}

// convertStreamResponse 转换流式响应
func (p *OpenAIProvider) convertStreamResponse(resp map[string]interface{}) *ChatResponse {
	choices, exists := resp["choices"].([]interface{})
	if !exists || len(choices) == 0 {
		return nil
	}

	choice := choices[0].(map[string]interface{})
	delta := choice["delta"].(map[string]interface{})

	content, _ := delta["content"].(string)
	role, _ := delta["role"].(string)

	return &ChatResponse{
		ID: fmt.Sprintf("stream_%d", time.Now().UnixNano()),
		Message: Message{
			Role:      role,
			Content:   content,
			Timestamp: time.Now(),
			Streaming: true,
		},
		Stream: true,
		Delta:  content,
		Done:   false,
	}
}

// isRetryableError 检查是否为可重试错误
func (p *OpenAIProvider) isRetryableError(code string) bool {
	retryableCodes := []string{
		"rate_limit_exceeded",
		"insufficient_quota",
		"engine_overloaded",
	}

	for _, retryableCode := range retryableCodes {
		if code == retryableCode {
			return true
		}
	}

	return false
}
