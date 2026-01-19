package router

import (
	apisystem "gin-admin-pro/internal/api/v1/system"
	apidao "gin-admin-pro/internal/dao/system"
	"gin-admin-pro/internal/middleware"
	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/internal/service"
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
	r.Use(middleware.Recovery())        // 自定义异常处理中间件
	r.Use(middleware.RateLimit())       // 限流中间件
	r.Use(middleware.OperationLogger()) // 操作日志中间件

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
				// 初始化DAO
				userDAO := apidao.NewUserDAO(service.Services.MySQLClient.GetDB())
				roleDAO := apidao.NewRoleDAO(service.Services.MySQLClient.GetDB())

				// 初始化控制器
				userCtrl := apisystem.NewUserController(userDAO, service.Services.TokenService)
				roleCtrl := apisystem.NewRoleController(roleDAO)

				// 用户管理路由（需要认证）
				user := system.Group("/user")
				user.Use(middleware.Auth()) // 认证中间件
				{
					user.GET("/page", userCtrl.Page)                                // 实现用户分页查询
					user.GET("/get", userCtrl.Get)                                  // 实现获取用户详情
					user.POST("/create", middleware.AdminOnly(), userCtrl.Create)   // 实现创建用户（仅管理员）
					user.PUT("/update", userCtrl.Update)                            // 实现更新用户
					user.DELETE("/delete", middleware.AdminOnly(), userCtrl.Delete) // 实现删除用户（仅管理员）
				}

				// 角色管理路由（需要认证）
				role := system.Group("/role")
				role.Use(middleware.Auth()) // 认证中间件
				{
					role.GET("/page", roleCtrl.Page)                                                         // 实现角色分页查询
					role.GET("/get", roleCtrl.Get)                                                           // 实现获取角色详情
					role.POST("/create", middleware.AdminOnly(), roleCtrl.Create)                            // 实现创建角色（仅管理员）
					role.PUT("/update", middleware.AdminOnly(), roleCtrl.Update)                             // 实现更新角色（仅管理员）
					role.DELETE("/delete", middleware.AdminOnly(), roleCtrl.Delete)                          // 实现删除角色（仅管理员）
					role.GET("/list-all-simple", roleCtrl.ListAllSimple)                                     // 实现获取角色精简列表
					role.PUT("/assign-menu/:roleId", middleware.AdminOnly(), roleCtrl.AssignMenuPermissions) // 分配菜单权限
					role.GET("/menu-permissions/:roleId", roleCtrl.GetMenuPermissions)                       // 获取菜单权限
				}

				// 菜单管理路由（需要认证）
				menu := system.Group("/menu")
				menu.Use(middleware.Auth()) // 认证中间件
				{
					menu.GET("/list", nil)                         // TODO: 实现菜单列表
					menu.GET("/get", nil)                          // TODO: 实现获取菜单详情
					menu.POST("/create", middleware.AdminOnly())   // TODO: 实现创建菜单（仅管理员）
					menu.PUT("/update", middleware.AdminOnly())    // TODO: 实现更新菜单（仅管理员）
					menu.DELETE("/delete", middleware.AdminOnly()) // TODO: 实现删除菜单（仅管理员）
				}

				// 部门管理路由（需要认证）
				dept := system.Group("/dept")
				dept.Use(middleware.Auth()) // 认证中间件
				{
					dept.GET("/list", nil)                         // TODO: 实现部门列表
					dept.GET("/get", nil)                          // TODO: 实现获取部门详情
					dept.POST("/create", middleware.AdminOnly())   // TODO: 实现创建部门（仅管理员）
					dept.PUT("/update", middleware.AdminOnly())    // TODO: 实现更新部门（仅管理员）
					dept.DELETE("/delete", middleware.AdminOnly()) // TODO: 实现删除部门（仅管理员）
				}

				// 认证路由（不需要认证）
				auth := system.Group("/auth")
				{
					auth.POST("/login", nil)                // TODO: 实现用户登录
					auth.POST("/logout", middleware.Auth()) // TODO: 实现用户登出（需要认证）
				}
			}

			// 基础设施模块（需要认证）
			infra := v1.Group("/infra")
			infra.Use(middleware.Auth()) // 认证中间件
			{
				file := infra.Group("/file")
				{
					file.POST("/upload", nil) // TODO: 实现文件上传
				}
			}

			// AI 模块（需要认证）
			ai := v1.Group("/ai")
			ai.Use(middleware.Auth()) // 认证中间件
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
