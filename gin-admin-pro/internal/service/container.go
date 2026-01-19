package service

import (
	"fmt"

	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/internal/pkg/token"
	"gin-admin-pro/plugin/mysql"
	"gin-admin-pro/plugin/oss"
	"gin-admin-pro/plugin/redis"
)

// Services 全局服务实例
var Services *ServiceContainer

// ServiceContainer 服务容器
type ServiceContainer struct {
	TokenService *token.TokenService
	RedisClient  *redis.Client
	MySQLClient  *mysql.Client
	OSSStorage   oss.OSSInterface
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

	// 初始化MySQL客户端
	mysqlConfig := &mysql.Config{
		Host:         cfg.Database.MySQL.Host,
		Port:         cfg.Database.MySQL.Port,
		Username:     cfg.Database.MySQL.Username,
		Password:     cfg.Database.MySQL.Password,
		Database:     cfg.Database.MySQL.Database,
		Charset:      cfg.Database.MySQL.Charset,
		ParseTime:    cfg.Database.MySQL.ParseTime,
		Loc:          cfg.Database.MySQL.Loc,
		MaxIdleConns: cfg.Database.MySQL.MaxIdleConns,
		MaxOpenConns: cfg.Database.MySQL.MaxOpenConns,
		MaxLifetime:  cfg.Database.MySQL.ConnMaxLifetime,
	}

	mysqlClient, err := mysql.NewClient(mysqlConfig)
	if err != nil {
		return fmt.Errorf("初始化MySQL客户端失败: %w", err)
	}

	// 初始化OSS存储
	ossStorage, err := oss.GetDefaultStorage()
	if err != nil {
		return fmt.Errorf("初始化OSS存储失败: %w", err)
	}

	// 设置全局服务实例
	Services = &ServiceContainer{
		TokenService: tokenService,
		RedisClient:  redisClient,
		MySQLClient:  mysqlClient,
		OSSStorage:   ossStorage,
	}

	return nil
}

// CleanupServices 清理服务
func CleanupServices() error {
	var err error

	if Services != nil {
		if Services.RedisClient != nil {
			if redisErr := Services.RedisClient.Close(); redisErr != nil {
				err = redisErr
			}
		}
		if Services.MySQLClient != nil {
			if mysqlErr := Services.MySQLClient.Close(); mysqlErr != nil {
				err = mysqlErr
			}
		}
	}

	return err
}
