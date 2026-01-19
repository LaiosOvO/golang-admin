package postgresql

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Client PostgreSQL客户端
type Client struct {
	db     *gorm.DB
	config *Config
}

// NewClient 创建PostgreSQL客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: getLogger(cfg.LogLevel, cfg.SlowThreshold),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// 获取底层的sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	client := &Client{
		db:     db,
		config: cfg,
	}

	// 启用扩展
	if err := client.enableExtensions(); err != nil {
		return nil, fmt.Errorf("failed to enable extensions: %w", err)
	}

	return client, nil
}

// GetDB 获取GORM DB实例
func (c *Client) GetDB() *gorm.DB {
	return c.db
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 检查数据库连接
func (c *Client) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// AutoMigrate 自动迁移数据库表
func (c *Client) AutoMigrate(dst ...interface{}) error {
	return c.db.AutoMigrate(dst...)
}

// HasTable 检查表是否存在
func (c *Client) HasTable(table string) bool {
	return c.db.Migrator().HasTable(table)
}

// DropTable 删除表
func (c *Client) DropTable(table string) error {
	return c.db.Migrator().DropTable(table)
}

// enableExtensions 启用数据库扩展
func (c *Client) enableExtensions() error {
	for _, ext := range c.config.Extensions {
		if !ext.Enabled {
			continue
		}

		if err := c.enableExtension(ext.Name, ext.Version); err != nil {
			return fmt.Errorf("failed to enable extension %s: %w", ext.Name, err)
		}
	}
	return nil
}

// enableExtension 启用单个扩展
func (c *Client) enableExtension(name, version string) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	// 检查扩展是否已存在
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = $1)`
	err = sqlDB.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check extension %s: %w", name, err)
	}

	if exists {
		return nil // 扩展已存在
	}

	// 创建扩展
	var createSQL string
	if version != "" {
		createSQL = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\" VERSION \"%s\"", name, version)
	} else {
		createSQL = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", name)
	}

	_, err = sqlDB.Exec(createSQL)
	if err != nil {
		// 尝试不指定版本
		if strings.Contains(err.Error(), "version") {
			createSQL = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", name)
			_, err = sqlDB.Exec(createSQL)
		}
		if err != nil {
			return fmt.Errorf("failed to create extension %s: %w", name, err)
		}
	}

	return nil
}

// GetExtensions 获取已启用的扩展列表
func (c *Client) GetExtensions() ([]string, error) {
	sqlDB, err := c.db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := sqlDB.Query("SELECT extname FROM pg_extension")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extensions []string
	for rows.Next() {
		var extName string
		if err := rows.Scan(&extName); err != nil {
			return nil, err
		}
		extensions = append(extensions, extName)
	}

	return extensions, nil
}

// IsExtensionEnabled 检查扩展是否启用
func (c *Client) IsExtensionEnabled(name string) (bool, error) {
	extensions, err := c.GetExtensions()
	if err != nil {
		return false, err
	}

	for _, ext := range extensions {
		if ext == name {
			return true, nil
		}
	}
	return false, nil
}

// getLogger 根据日志级别获取logger
func getLogger(level string, slowThreshold time.Duration) logger.Interface {
	switch level {
	case "silent":
		return logger.Default.LogMode(logger.Silent)
	case "error":
		return logger.Default.LogMode(logger.Error)
	case "warn":
		return logger.Default.LogMode(logger.Warn)
	case "info":
		return logger.Default.LogMode(logger.Info)
	default:
		return logger.Default.LogMode(logger.Info)
	}
}
