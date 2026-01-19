package kafka

import (
	"time"
)

// Config Kafka配置结构
type Config struct {
	// 基础配置
	Brokers []string `yaml:"brokers" mapstructure:"brokers"`

	// 消费者配置
	GroupID    string `yaml:"groupId" mapstructure:"groupId"`
	AutoCommit bool   `yaml:"autoCommit" mapstructure:"autoCommit"`
	AutoOffset string `yaml:"autoOffset" mapstructure:"autoOffset"` // earliest, latest, none

	// 网络配置
	DialTimeout  time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`

	// 重试配置
	MaxRetries   int           `yaml:"maxRetries" mapstructure:"maxRetries"`
	RetryBackoff time.Duration `yaml:"retryBackoff" mapstructure:"retryBackoff"`

	// 生产者配置
	Producer struct {
		RequiredAcks    int16         `yaml:"requiredAcks" mapstructure:"requiredAcks"` // 0=NoResponse, 1=WaitForLocal, -1=WaitForAll
		Compression     string        `yaml:"compression" mapstructure:"compression"`   // none, gzip, snappy, lz4, zstd
		FlushFrequency  time.Duration `yaml:"flushFrequency" mapstructure:"flushFrequency"`
		FlushBytes      int           `yaml:"flushBytes" mapstructure:"flushBytes"`
		FlushMessages   int           `yaml:"flushMessages" mapstructure:"flushMessages"`
		MaxMessageBytes int           `yaml:"maxMessageBytes" mapstructure:"maxMessageBytes"`
		Partitioner     string        `yaml:"partitioner" mapstructure:"partitioner"` // hash, random, roundrobin, manual
		ReturnSuccesses bool          `yaml:"returnSuccesses" mapstructure:"returnSuccesses"`
		ReturnErrors    bool          `yaml:"returnErrors" mapstructure:"returnErrors"`
	} `yaml:"producer" mapstructure:"producer"`

	// 消费者配置
	Consumer struct {
		FetchMin             int32         `yaml:"fetchMin" mapstructure:"fetchMin"`
		FetchDefault         int32         `yaml:"fetchDefault" mapstructure:"fetchDefault"`
		FetchMax             int32         `yaml:"fetchMax" mapstructure:"fetchMax"`
		MaxProcessingTime    time.Duration `yaml:"maxProcessingTime" mapstructure:"maxProcessingTime"`
		OffsetCommitInterval time.Duration `yaml:"offsetCommitInterval" mapstructure:"offsetCommitInterval"`
		SessionTimeout       time.Duration `yaml:"sessionTimeout" mapstructure:"sessionTimeout"`
		HeartbeatInterval    time.Duration `yaml:"heartbeatInterval" mapstructure:"heartbeatInterval"`
		MaxWaitTime          time.Duration `yaml:"maxWaitTime" mapstructure:"maxWaitTime"`
	} `yaml:"consumer" mapstructure:"consumer"`

	// 安全配置
	Security struct {
		Enabled      bool   `yaml:"enabled" mapstructure:"enabled"`
		Mechanism    string `yaml:"mechanism" mapstructure:"mechanism"` // PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
		Username     string `yaml:"username" mapstructure:"username"`
		Password     string `yaml:"password" mapstructure:"password"`
		TLSEnabled   bool   `yaml:"tlsEnabled" mapstructure:"tlsEnabled"`
		TLSVersion   string `yaml:"tlsVersion" mapstructure:"tlsVersion"`
		TLSCertFile  string `yaml:"tlsCertFile" mapstructure:"tlsCertFile"`
		TLSKeyFile   string `yaml:"tlsKeyFile" mapstructure:"tlsKeyFile"`
		TLSCAFile    string `yaml:"tlsCaFile" mapstructure:"tlsCaFile"`
		InsecureSkip bool   `yaml:"insecureSkip" mapstructure:"insecureSkip"`
	} `yaml:"security" mapstructure:"security"`

	// 监控配置
	Metrics struct {
		Enabled bool     `yaml:"enabled" mapstructure:"enabled"`
		Prefix  string   `yaml:"prefix" mapstructure:"prefix"`
		Tags    []string `yaml:"tags" mapstructure:"tags"`
	} `yaml:"metrics" mapstructure:"metrics"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	config := &Config{
		Brokers:      []string{"localhost:9092"},
		GroupID:      "gin-admin-group",
		AutoCommit:   true,
		AutoOffset:   "latest",
		DialTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		MaxRetries:   3,
		RetryBackoff: time.Second * 2,
	}

	// 生产者默认配置
	config.Producer.RequiredAcks = 1 // WaitForLocal
	config.Producer.Compression = "none"
	config.Producer.FlushFrequency = time.Millisecond * 100
	config.Producer.FlushBytes = 16384 // 16KB
	config.Producer.FlushMessages = 0
	config.Producer.MaxMessageBytes = 1000000 // 1MB
	config.Producer.Partitioner = "hash"
	config.Producer.ReturnSuccesses = true
	config.Producer.ReturnErrors = true

	// 消费者默认配置
	config.Consumer.FetchMin = 1
	config.Consumer.FetchDefault = 1024 * 1024  // 1MB
	config.Consumer.FetchMax = 10 * 1024 * 1024 // 10MB
	config.Consumer.MaxProcessingTime = time.Second * 30
	config.Consumer.OffsetCommitInterval = time.Second
	config.Consumer.SessionTimeout = time.Second * 10
	config.Consumer.HeartbeatInterval = time.Second * 3
	config.Consumer.MaxWaitTime = time.Millisecond * 500

	// 安全默认配置
	config.Security.Enabled = false
	config.Security.Mechanism = "PLAIN"
	config.Security.TLSEnabled = false
	config.Security.TLSVersion = "1.2"
	config.Security.InsecureSkip = true

	// 监控默认配置
	config.Metrics.Enabled = false
	config.Metrics.Prefix = "kafka"
	config.Metrics.Tags = []string{"gin-admin"}

	return config
}
