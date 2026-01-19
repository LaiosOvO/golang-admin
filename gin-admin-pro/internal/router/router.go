package router

import (
	"gin-admin-pro/internal/pkg/config"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	cfg := config.GetConfig()

	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS 中间件
	if cfg.CORS.Enabled {
		r.Use(corsMiddleware(cfg.CORS))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "服务正常",
			"data": gin.H{
				"status": "ok",
				"server": cfg.Server.Name,
			},
		})
	})

	// API 分组
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 系统管理模块
			system := v1.Group("/system")
			{
				// 用户管理路由
				user := system.Group("/user")
				{
					user.GET("/page", nil)      // TODO: 实现用户分页查询
					user.GET("/get", nil)       // TODO: 实现获取用户详情
					user.POST("/create", nil)   // TODO: 实现创建用户
					user.PUT("/update", nil)    // TODO: 实现更新用户
					user.DELETE("/delete", nil) // TODO: 实现删除用户
				}

				// 角色管理路由
				role := system.Group("/role")
				{
					role.GET("/page", nil)            // TODO: 实现角色分页查询
					role.GET("/get", nil)             // TODO: 实现获取角色详情
					role.POST("/create", nil)         // TODO: 实现创建角色
					role.PUT("/update", nil)          // TODO: 实现更新角色
					role.DELETE("/delete", nil)       // TODO: 实现删除角色
					role.GET("/list-all-simple", nil) // TODO: 实现获取角色精简列表
				}

				// 菜单管理路由
				menu := system.Group("/menu")
				{
					menu.GET("/list", nil)      // TODO: 实现菜单列表
					menu.GET("/get", nil)       // TODO: 实现获取菜单详情
					menu.POST("/create", nil)   // TODO: 实现创建菜单
					menu.PUT("/update", nil)    // TODO: 实现更新菜单
					menu.DELETE("/delete", nil) // TODO: 实现删除菜单
				}

				// 部门管理路由
				dept := system.Group("/dept")
				{
					dept.GET("/list", nil)      // TODO: 实现部门列表
					dept.GET("/get", nil)       // TODO: 实现获取部门详情
					dept.POST("/create", nil)   // TODO: 实现创建部门
					dept.PUT("/update", nil)    // TODO: 实现更新部门
					dept.DELETE("/delete", nil) // TODO: 实现删除部门
				}

				// 认证路由
				auth := system.Group("/auth")
				{
					auth.POST("/login", nil)  // TODO: 实现用户登录
					auth.POST("/logout", nil) // TODO: 实现用户登出
				}
			}

			// 基础设施模块
			infra := v1.Group("/infra")
			{
				// 文件上传
				file := infra.Group("/file")
				{
					file.POST("/upload", nil) // TODO: 实现文件上传
				}
			}

			// AI 模块
			ai := v1.Group("/ai")
			{
				ai.POST("/chat", nil) // TODO: 实现AI对话
			}
		}
	}

	return r
}

// corsMiddleware CORS 中间件
func corsMiddleware(corsConfig config.CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
	})
}
