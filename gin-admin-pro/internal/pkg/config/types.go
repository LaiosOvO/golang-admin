package config

import (
	"time"
)

// Config 应用配置结构
type Config struct {
	Server    ServerConfig    `yaml:"server" json:"server"`
	Database  DatabaseConfig  `yaml:"database" json:"database"`
	Kafka     KafkaConfig     `yaml:"kafka" json:"kafka"`
	AI        AIConfig        `yaml:"ai" json:"ai"`
	JWT       JWTConfig       `yaml:"jwt" json:"jwt"`
	Log       LogConfig       `yaml:"log" json:"log"`
	CORS      CORSConfig      `yaml:"cors" json:"cors"`
	RateLimit RateLimitConfig `yaml:"rateLimit" json:"rateLimit"`
	Upload    UploadConfig    `yaml:"upload" json:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `yaml:"port" json:"port"`
	Mode         string `yaml:"mode" json:"mode"`
	Name         string `yaml:"name" json:"name"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL         MySQLConfig         `yaml:"mysql" json:"mysql"`
	PostgreSQL    PostgreSQLConfig    `yaml:"postgresql" json:"postgresql"`
	MongoDB       MongoDBConfig       `yaml:"mongodb" json:"mongodb"`
	Redis         RedisConfig         `yaml:"redis" json:"redis"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch" json:"elasticsearch"`
	Milvus        MilvusConfig        `yaml:"milvus" json:"milvus"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	Database        string        `yaml:"database" json:"database"`
	Username        string        `yaml:"username" json:"username"`
	Password        string        `yaml:"password" json:"password"`
	Charset         string        `yaml:"charset" json:"charset"`
	ParseTime       bool          `yaml:"parseTime" json:"parseTime"`
	Loc             string        `yaml:"loc" json:"loc"`
	MaxIdleConns    int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"`
}

// PostgreSQLConfig PostgreSQL配置
type PostgreSQLConfig struct {
	Host         string   `yaml:"host" json:"host"`
	Port         int      `yaml:"port" json:"port"`
	Database     string   `yaml:"database" json:"database"`
	Username     string   `yaml:"username" json:"username"`
	Password     string   `yaml:"password" json:"password"`
	SSLMode      string   `yaml:"sslmode" json:"sslmode"`
	MaxIdleConns int      `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int      `yaml:"maxOpenConns" json:"maxOpenConns"`
	Extensions   []string `yaml:"extensions" json:"extensions"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string        `yaml:"uri" json:"uri"`
	Database string        `yaml:"database" json:"database"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	Password     string `yaml:"password" json:"password"`
	Database     int    `yaml:"database" json:"database"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns" json:"minIdleConns"`
}

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	URLs     []string `yaml:"urls" json:"urls"`
	Username string   `yaml:"username" json:"username"`
	Password string   `yaml:"password" json:"password"`
}

// MilvusConfig Milvus配置
type MilvusConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
	GroupID string   `yaml:"groupId" json:"groupId"`
}

// AIConfig AI配置
type AIConfig struct {
	Enabled     bool    `yaml:"enabled" json:"enabled"`
	Provider    string  `yaml:"provider" json:"provider"`
	APIKey      string  `yaml:"apiKey" json:"apiKey"`
	BaseURL     string  `yaml:"baseUrl" json:"baseUrl"`
	Model       string  `yaml:"model" json:"model"`
	MaxTokens   int     `yaml:"maxTokens" json:"maxTokens"`
	Temperature float64 `yaml:"temperature" json:"temperature"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `yaml:"secret" json:"secret"`
	AccessTokenExpire  int    `yaml:"accessTokenExpire" json:"accessTokenExpire"`
	RefreshTokenExpire int    `yaml:"refreshTokenExpire" json:"refreshTokenExpire"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string   `yaml:"level" json:"level"`
	Format     string   `yaml:"format" json:"format"`
	Output     []string `yaml:"output" json:"output"`
	FilePath   string   `yaml:"filePath" json:"filePath"`
	MaxSize    int      `yaml:"maxSize" json:"maxSize"`
	MaxBackups int      `yaml:"maxBackups" json:"maxBackups"`
	MaxAge     int      `yaml:"maxAge" json:"maxAge"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled" json:"enabled"`
	AllowOrigins     []string `yaml:"allowOrigins" json:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods" json:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders" json:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders" json:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials" json:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge" json:"maxAge"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled  bool `yaml:"enabled" json:"enabled"`
	Requests int  `yaml:"requests" json:"requests"`
	Window   int  `yaml:"window" json:"window"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int      `yaml:"maxSize" json:"maxSize"`
	AllowedTypes []string `yaml:"allowedTypes" json:"allowedTypes"`
	Path         string   `yaml:"path" json:"path"`
	URLPrefix    string   `yaml:"urlPrefix" json:"urlPrefix"`
}
