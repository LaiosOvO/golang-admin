package mongodb

import (
	"fmt"
	"time"
)

// Config MongoDB配置结构
type Config struct {
	URI            string        `yaml:"uri" mapstructure:"uri"`
	Host           string        `yaml:"host" mapstructure:"host"`
	Port           int           `yaml:"port" mapstructure:"port"`
	Database       string        `yaml:"database" mapstructure:"database"`
	Username       string        `yaml:"username" mapstructure:"username"`
	Password       string        `yaml:"password" mapstructure:"password"`
	AuthSource     string        `yaml:"authSource" mapstructure:"authSource"`
	MaxPoolSize    uint64        `yaml:"maxPoolSize" mapstructure:"maxPoolSize"`
	MinPoolSize    uint64        `yaml:"minPoolSize" mapstructure:"minPoolSize"`
	MaxConnIdle    time.Duration `yaml:"maxConnIdle" mapstructure:"maxConnIdle"`
	ConnectTimeout time.Duration `yaml:"connectTimeout" mapstructure:"connectTimeout"`
	ServerTimeout  time.Duration `yaml:"serverTimeout" mapstructure:"serverTimeout"`
	Timeout        time.Duration `yaml:"timeout" mapstructure:"timeout"`
	CompressLevel  int           `yaml:"compressLevel" mapstructure:"compressLevel"`
	ReplicaSet     string        `yaml:"replicaSet" mapstructure:"replicaSet"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		URI:            "mongodb://localhost:27017",
		Host:           "localhost",
		Port:           27017,
		Database:       "gin_admin",
		Username:       "",
		Password:       "",
		AuthSource:     "admin",
		MaxPoolSize:    100,
		MinPoolSize:    10,
		MaxConnIdle:    time.Minute * 5,
		ConnectTimeout: time.Second * 10,
		ServerTimeout:  time.Second * 30,
		Timeout:        time.Second * 30,
		CompressLevel:  6,
		ReplicaSet:     "",
	}
}

// GetURI 获取连接URI
func (c *Config) GetURI() string {
	if c.URI != "" && c.URI != "mongodb://localhost:27017" {
		return c.URI
	}

	// 如果host为空，使用默认值
	host := c.Host
	if host == "" {
		host = "localhost"
	}

	// 如果port为0，使用默认值
	port := c.Port
	if port == 0 {
		port = 27017
	}

	uri := "mongodb://"

	if c.Username != "" && c.Password != "" {
		uri += c.Username + ":" + c.Password + "@"
	}

	uri += host + ":" + fmt.Sprintf("%d", port)

	if c.AuthSource != "" {
		uri += "/" + c.AuthSource
	}

	return uri
}
