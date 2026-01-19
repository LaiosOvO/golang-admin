package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis客户端
type Client struct {
	client redis.Cmdable
	config *Config
}

// NewClient 创建Redis客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	var client redis.Cmdable
	var err error

	if cfg.ClusterEnabled && len(cfg.ClusterAddrs) > 0 {
		// 集群模式
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           cfg.ClusterAddrs,
			Password:        cfg.Password,
			Username:        cfg.Username,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err = rdb.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to Redis cluster: %w", err)
		}

		client = rdb

	} else if cfg.SentinelEnabled && len(cfg.SentinelAddrs) > 0 {
		// 哨兵模式
		rdb := redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:      cfg.SentinelMaster,
			SentinelAddrs:   cfg.SentinelAddrs,
			Password:        cfg.Password,
			Username:        cfg.Username,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
			DB:              cfg.DB,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err = rdb.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to Redis via sentinel: %w", err)
		}

		client = rdb

	} else if cfg.ShardEnabled && len(cfg.ShardAddrs) > 0 {
		// 分片模式
		addrs := make(map[string]string)
		for i, addr := range cfg.ShardAddrs {
			addrs[fmt.Sprintf("shard%d", i)] = addr
		}

		rdb := redis.NewRing(&redis.RingOptions{
			Addrs:           addrs,
			Password:        cfg.Password,
			Username:        cfg.Username,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err = rdb.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to Redis sharded cluster: %w", err)
		}

		client = rdb

	} else {
		// 单机模式
		rdb := redis.NewClient(&redis.Options{
			Addr:            cfg.Addr,
			Password:        cfg.Password,
			Username:        cfg.Username,
			DB:              cfg.DB,
			MaxRetries:      cfg.MaxRetries,
			MinRetryBackoff: cfg.MinRetryBackoff,
			MaxRetryBackoff: cfg.MaxRetryBackoff,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err = rdb.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to Redis: %w", err)
		}

		client = rdb
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// GetClient 获取Redis客户端
func (c *Client) GetClient() redis.Cmdable {
	return c.client
}

// Close 关闭连接
func (c *Client) Close() error {
	switch v := c.client.(type) {
	case *redis.Client:
		return v.Close()
	case *redis.ClusterClient:
		return v.Close()
	case *redis.Ring:
		return v.Close()
	default:
		return nil
	}
}

// Ping 检查连接
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Set 设置键值对
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// HSet 设置哈希字段
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// SAdd 添加到集合
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SRem 从集合删除成员
func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, key, members...).Err()
}

// ZAdd 添加到有序集合
func (c *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围内的成员
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围内的成员和分数
func (c *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return c.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem 从有序集合删除成员
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

// LPush 从左侧推入列表
func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush 从右侧推入列表
func (c *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.RPush(ctx, key, values...).Err()
}

// LPop 从左侧弹出列表元素
func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	return c.client.LPop(ctx, key).Result()
}

// RPop 从右侧弹出列表元素
func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	return c.client.RPop(ctx, key).Result()
}

// LRange 获取列表范围内的元素
func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// Incr 递增
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Decr 递减
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// IncrBy 按指定值递增
func (c *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// DecrBy 按指定值递减
func (c *Client) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, key, value).Result()
}

// SetJSON 设置JSON对象
func (c *Client) SetJSON(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %w", err)
	}
	return c.Set(ctx, key, jsonData, expiration)
}

// GetJSON 获取JSON对象
func (c *Client) GetJSON(ctx context.Context, key string, obj interface{}) error {
	jsonData, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonData), obj)
}

// GetInfo 获取Redis信息
func (c *Client) GetInfo(ctx context.Context) (string, error) {
	return c.client.Info(ctx).Result()
}

// FlushDB 清空当前数据库
func (c *Client) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// FlushAll 清空所有数据库
func (c *Client) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}
