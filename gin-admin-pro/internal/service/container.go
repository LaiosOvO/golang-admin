package service

import (
	"fmt"

	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/internal/pkg/token"
	"gin-admin-pro/plugin/redis"
)

// Services 全局服务实例
var Services *ServiceContainer

// ServiceContainer 服务容器
type ServiceContainer struct {
	TokenService *token.TokenService
	RedisClient  *redis.Client
}

// InitServices 初始化服务
func InitServices() error {
	cfg := config.GetConfig()

	// 初始化Redis客户端
	redisConfig := &redis.Config{
		Addr:         fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password:     cfg.Database.Redis.Password,
		DB:           cfg.Database.Redis.Database,
		PoolSize:     cfg.Database.Redis.PoolSize,
		MinIdleConns: cfg.Database.Redis.MinIdleConns,
	}

	redisClient, err := redis.NewClient(redisConfig)
	if err != nil {
		return fmt.Errorf("初始化Redis客户端失败: %w", err)
	}

	// 初始化Token服务
	tokenService := token.NewTokenService(redisClient)

	// 设置全局服务实例
	Services = &ServiceContainer{
		TokenService: tokenService,
		RedisClient:  redisClient,
	}

	return nil
}

// CleanupServices 清理服务
func CleanupServices() error {
	if Services != nil && Services.RedisClient != nil {
		return Services.RedisClient.Close()
	}
	return nil
}
