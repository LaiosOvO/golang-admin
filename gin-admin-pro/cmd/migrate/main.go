package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gin-admin-pro/internal/migration"
	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/plugin/mysql"
)

var (
	action = flag.String("action", "migrate", "迁移动作: migrate, reset, drop")
	env    = flag.String("env", "dev", "环境: dev, test, prod")
)

func main() {
	flag.Parse()

	// 加载配置
	var err error
	if *env == "dev" || *env == "test" || *env == "prod" {
		err = config.LoadWithEnv(*env)
	} else {
		err = config.Load("")
	}

	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 获取配置
	cfg := config.GetConfig()

	// 转换配置类型
	mysqlConfig := &mysql.Config{
		Host:          cfg.Database.MySQL.Host,
		Port:          cfg.Database.MySQL.Port,
		Database:      cfg.Database.MySQL.Database,
		Username:      cfg.Database.MySQL.Username,
		Password:      cfg.Database.MySQL.Password,
		Charset:       cfg.Database.MySQL.Charset,
		ParseTime:     cfg.Database.MySQL.ParseTime,
		Loc:           cfg.Database.MySQL.Loc,
		MaxIdleConns:  cfg.Database.MySQL.MaxIdleConns,
		MaxOpenConns:  cfg.Database.MySQL.MaxOpenConns,
		MaxLifetime:   cfg.Database.MySQL.ConnMaxLifetime,
		LogLevel:      "info",
		SlowThreshold: 200 * time.Millisecond,
	}

	// 初始化MySQL连接
	mysqlClient, err := mysql.NewClient(mysqlConfig)
	if err != nil {
		log.Fatalf("连接MySQL失败: %v", err)
	}
	defer mysqlClient.Close()

	// 创建迁移器
	migrator := migration.NewMigrator(mysqlClient.GetDB())

	// 执行相应的动作
	switch *action {
	case "migrate":
		log.Println("开始数据库迁移...")
		if err := migrator.AutoMigrate(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		fmt.Println("数据库迁移完成!")

	case "reset":
		fmt.Println("警告：这将删除所有表并重新创建!")
		fmt.Print("确认继续吗? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("操作已取消")
			os.Exit(0)
		}

		log.Println("开始重置数据库...")
		if err := migrator.ResetDatabase(); err != nil {
			log.Fatalf("数据库重置失败: %v", err)
		}
		fmt.Println("数据库重置完成!")

	case "drop":
		fmt.Println("警告：这将删除所有表!")
		fmt.Print("确认继续吗? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("操作已取消")
			os.Exit(0)
		}

		log.Println("开始删除所有表...")
		if err := migrator.DropAllTables(); err != nil {
			log.Fatalf("删除表失败: %v", err)
		}
		fmt.Println("所有表删除完成!")

	default:
		fmt.Printf("不支持的动作: %s\n", *action)
		fmt.Println("支持的动作: migrate, reset, drop")
		os.Exit(1)
	}
}
