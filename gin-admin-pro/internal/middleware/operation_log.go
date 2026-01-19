package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationLog 操作日志结构
type OperationLog struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Username  string    `json:"username"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Params    string    `json:"params"`
	Body      string    `json:"body"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Status    int       `json:"status"`
	Duration  int64     `json:"duration"`
	ErrorMsg  string    `json:"errorMsg"`
	CreatedAt time.Time `json:"createdAt"`
}

// OperationLogger 操作日志中间件
func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 获取用户信息
		userID, _ := c.Get("userId")
		username, _ := c.Get("username")
		usernameStr, _ := username.(string)
		userIDUint, _ := userID.(uint)
		if userIDUint == 0 {
			userIDUint = 0 // 未认证用户
		}
		if usernameStr == "" {
			usernameStr = "anonymous" // 匿名用户
		}

		// 创建操作日志
		log := OperationLog{
			UserID:    userIDUint,
			Username:  usernameStr,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Params:    c.Request.URL.RawQuery,
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			CreatedAt: startTime,
		}

		// 记录请求体（仅对 JSON 请求，且排除敏感信息）
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// 过滤敏感字段
				filteredBody := filterSensitiveData(bodyBytes)
				log.Body = string(filteredBody)
			}
		}

		// 使用自定义的 ResponseWriter 来捕获响应状态
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// 继续处理请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)
		log.Duration = duration.Milliseconds()

		// 记录响应状态
		log.Status = c.Writer.Status()

		// 记录错误信息
		if len(c.Errors) > 0 {
			log.ErrorMsg = c.Errors.Last().Error()
		}

		// TODO: 将日志保存到数据库
		// 这里先输出到控制台，后续实现数据库存储
		printOperationLog(log)
	}
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(data []byte) []byte {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return data
	}

	// 转换为 map 以便处理
	if m, ok := obj.(map[string]interface{}); ok {
		// 过滤的字段列表
		sensitiveFields := []string{"password", "password_confirmation", "secret", "token", "key"}

		for _, field := range sensitiveFields {
			if _, exists := m[field]; exists {
				m[field] = "***"
			}
		}

		// 重新编码为 JSON
		filtered, err := json.Marshal(m)
		if err != nil {
			return data
		}
		return filtered
	}

	return data
}

// printOperationLog 打印操作日志
func printOperationLog(log OperationLog) {
	// 转换为 JSON 格式输出
	logJSON, _ := json.Marshal(log)
	fmt.Printf("[OperationLog] %s\n", string(logJSON))
}

// responseBodyWriter 自定义 ResponseWriter 用于捕获响应
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
