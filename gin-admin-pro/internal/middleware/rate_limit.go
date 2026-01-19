package middleware

import (
	"net/http"
	"time"

	"gin-admin-pro/internal/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimit 创建限流中间件
func RateLimit() gin.HandlerFunc {
	cfg := config.GetConfig()

	if !cfg.RateLimit.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 创建限流器配置
	rate := limiter.Rate{
		Period: time.Duration(cfg.RateLimit.Window) * time.Second,
		Limit:  int64(cfg.RateLimit.Requests),
	}

	// 创建内存存储
	store := memory.NewStore()

	// 创建限流器
	instance := limiter.New(store, rate)

	// 返回限流中间件
	middleware := stdlib.NewMiddleware(instance)

	return func(c *gin.Context) {
		middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)

		// 检查是否被限流
		if c.Writer.Status() == http.StatusTooManyRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
			return
		}
	}
}
