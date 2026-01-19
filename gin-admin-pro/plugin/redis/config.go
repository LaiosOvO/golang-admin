package redis

import (
	"time"
)

// Config Redis配置结构
type Config struct {
	Addr            string        `yaml:"addr" mapstructure:"addr"`
	Password        string        `yaml:"password" mapstructure:"password"`
	DB              int           `yaml:"db" mapstructure:"db"`
	Username        string        `yaml:"username" mapstructure:"username"`
	MaxRetries      int           `yaml:"maxRetries" mapstructure:"maxRetries"`
	MinRetryBackoff time.Duration `yaml:"minRetryBackoff" mapstructure:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `yaml:"maxRetryBackoff" mapstructure:"maxRetryBackoff"`
	DialTimeout     time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	ReadTimeout     time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`
	PoolSize        int           `yaml:"poolSize" mapstructure:"poolSize"`
	MinIdleConns    int           `yaml:"minIdleConns" mapstructure:"minIdleConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" mapstructure:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" mapstructure:"connMaxLifetime"`

	// 集群配置
	ClusterEnabled bool     `yaml:"clusterEnabled" mapstructure:"clusterEnabled"`
	ClusterAddrs   []string `yaml:"clusterAddrs" mapstructure:"clusterAddrs"`

	// 哨兵配置
	SentinelEnabled bool     `yaml:"sentinelEnabled" mapstructure:"sentinelEnabled"`
	SentinelAddrs   []string `yaml:"sentinelAddrs" mapstructure:"sentinelAddrs"`
	SentinelMaster  string   `yaml:"sentinelMaster" mapstructure:"sentinelMaster"`

	// 分片配置
	ShardEnabled bool              `yaml:"shardEnabled" mapstructure:"shardEnabled"`
	ShardAddrs   []string          `yaml:"shardAddrs" mapstructure:"shardAddrs"`
	Sharding     ShardingAlgorithm `yaml:"sharding" mapstructure:"sharding"`
}

// ShardingAlgorithm 分片算法配置
type ShardingAlgorithm struct {
	Type    string                 `yaml:"type" mapstructure:"type"` // consistent, ketama, range, fnv1a
	KeyFunc string                 `yaml:"keyFunc" mapstructure:"keyFunc"`
	Options map[string]interface{} `yaml:"options" mapstructure:"options"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Addr:            "localhost:6379",
		Password:        "",
		DB:              0,
		Username:        "",
		MaxRetries:      3,
		MinRetryBackoff: time.Millisecond * 8,
		MaxRetryBackoff: time.Millisecond * 512,
		DialTimeout:     time.Second * 5,
		ReadTimeout:     time.Second * 3,
		WriteTimeout:    time.Second * 3,
		PoolSize:        10,
		MinIdleConns:    5,
		MaxIdleConns:    100,
		ConnMaxIdleTime: time.Minute * 5,
		ConnMaxLifetime: time.Hour,

		ClusterEnabled: false,
		ClusterAddrs:   []string{},

		SentinelEnabled: false,
		SentinelAddrs:   []string{},
		SentinelMaster:  "mymaster",

		ShardEnabled: false,
		ShardAddrs:   []string{},
		Sharding: ShardingAlgorithm{
			Type:    "consistent",
			KeyFunc: "key",
			Options: map[string]interface{}{},
		},
	}
}
