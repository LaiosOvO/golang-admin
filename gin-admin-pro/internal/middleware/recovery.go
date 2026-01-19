package middleware

import (
	"gin-admin-pro/internal/pkg/config"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 自定义恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()

				// 在开发模式下，返回详细的错误信息
				if config.GetConfig().Server.Mode == "debug" {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "服务器内部错误",
						"data": gin.H{
							"error": err,
							"stack": string(stack),
						},
					})
				} else {
					// 在生产模式下，只返回简单的错误信息
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "服务器内部错误",
						"data":    nil,
					})
				}

				// 终止请求处理
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last().Err

			// 根据错误类型返回不同的响应
			switch e := err.(type) {
			case *gin.Error:
				// Gin 错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": e.Error(),
					"data":    nil,
				})
			default:
				// 其他错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "请求处理失败",
					"data": gin.H{
						"error": e.Error(),
					},
				})
			}

			// 终止请求处理
			c.Abort()
		}
	}
}
