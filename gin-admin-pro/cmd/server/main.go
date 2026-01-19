package main

import (
	"fmt"
	"log"
	"os"

	"gin-admin-pro/internal/pkg/config"
)

func main() {
	// 获取环境变量
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // 默认开发环境
	}

	// 加载配置
	var err error
	if env == "dev" || env == "test" || env == "prod" {
		err = config.LoadWithEnv(env)
	} else {
		err = config.Load("")
	}

	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 获取配置
	cfg := config.GetConfig()

	fmt.Printf("=== %s 启动 ===\n", cfg.Server.Name)
	fmt.Printf("环境: %s\n", env)
	fmt.Printf("端口: %d\n", cfg.Server.Port)
	fmt.Printf("模式: %s\n", cfg.Server.Mode)

	// TODO: 初始化数据库连接
	// TODO: 初始化路由
	// TODO: 启动服务器

	fmt.Printf("服务器启动在端口 %d\n", cfg.Server.Port)
}
