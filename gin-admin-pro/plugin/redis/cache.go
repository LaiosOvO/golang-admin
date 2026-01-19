package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache 缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	GetOrSet(ctx context.Context, key string, callback func() (interface{}, error), expiration time.Duration) (interface{}, error)
	SetJSON(ctx context.Context, key string, obj interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, obj interface{}) error
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(client *Client) *RedisCache {
	return &RedisCache{client: client}
}

// GetOrSet 获取或设置缓存（使用回调函数）
func (c *RedisCache) GetOrSet(ctx context.Context, key string, callback func() (interface{}, error), expiration time.Duration) (interface{}, error) {
	// 尝试从缓存获取
	val, err := c.client.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// 缓存不存在，执行回调获取数据
	data, err := callback()
	if err != nil {
		return nil, fmt.Errorf("callback failed: %w", err)
	}

	// 设置缓存
	if err := c.client.Set(ctx, key, data, expiration); err != nil {
		// 记录错误但不返回，因为数据已经获取成功
		fmt.Printf("failed to set cache for key %s: %v\n", key, err)
	}

	return data, nil
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration)
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key)
}

// GetBytes 获取字节数组
func (c *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

// Del 删除缓存
func (c *RedisCache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...)
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration)
}

// TTL 获取剩余过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key)
}

// SetJSON 设置JSON对象缓存
func (c *RedisCache) SetJSON(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	return c.client.SetJSON(ctx, key, obj, expiration)
}

// GetJSON 获取JSON对象缓存
func (c *RedisCache) GetJSON(ctx context.Context, key string, obj interface{}) error {
	return c.client.GetJSON(ctx, key, obj)
}

// GetOrSetJSON 获取或设置JSON对象缓存
func (c *RedisCache) GetOrSetJSON(ctx context.Context, key string, obj interface{}, callback func() (interface{}, error), expiration time.Duration) error {
	// 尝试从缓存获取
	err := c.GetJSON(ctx, key, obj)
	if err == nil {
		return nil
	}

	// 缓存不存在，执行回调获取数据
	data, err := callback()
	if err != nil {
		return fmt.Errorf("callback failed: %w", err)
	}

	// 设置缓存
	if err := c.SetJSON(ctx, key, data, expiration); err != nil {
		// 记录错误但不返回，因为数据已经获取成功
		fmt.Printf("failed to set cache for key %s: %v\n", key, err)
	}

	// 解码到目标对象
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return json.Unmarshal(jsonData, obj)
}

// MemoryCache 内存缓存实现（用于测试或小型应用）
type MemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		data: make(map[string]cacheItem),
	}

	// 启动清理过期项的goroutine
	go cache.cleanup()

	return cache
}

// Set 设置内存缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expirationTime time.Time
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration)
	}

	c.data[key] = cacheItem{
		value:      value,
		expiration: expirationTime,
	}

	return nil
}

// Get 获取内存缓存
func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return "", fmt.Errorf("key expired")
	}

	switch v := item.value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetBytes 获取字节数组
func (c *MemoryCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

// Del 删除内存缓存
func (c *MemoryCache) Del(ctx context.Context, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.data, key)
	}

	return nil
}

// Exists 检查内存缓存是否存在
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Expire 设置过期时间
func (c *MemoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}

	item.expiration = time.Now().Add(expiration)
	c.data[key] = item

	return nil
}

// TTL 获取剩余过期时间
func (c *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return -1, nil
	}

	if item.expiration.IsZero() {
		return -1, nil
	}

	ttl := time.Until(item.expiration)
	if ttl <= 0 {
		delete(c.data, key)
		return -2, nil
	}

	return ttl, nil
}

// SetJSON 设置JSON对象缓存
func (c *MemoryCache) SetJSON(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %w", err)
	}
	return c.Set(ctx, key, jsonData, expiration)
}

// GetJSON 获取JSON对象缓存
func (c *MemoryCache) GetJSON(ctx context.Context, key string, obj interface{}) error {
	jsonData, err := c.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, obj)
}

// GetOrSet 获取或设置缓存（使用回调函数）
func (c *MemoryCache) GetOrSet(ctx context.Context, key string, callback func() (interface{}, error), expiration time.Duration) (interface{}, error) {
	// 尝试从缓存获取
	val, err := c.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// 缓存不存在，执行回调获取数据
	data, err := callback()
	if err != nil {
		return nil, fmt.Errorf("callback failed: %w", err)
	}

	// 设置缓存
	if err := c.Set(ctx, key, data, expiration); err != nil {
		// 记录错误但不返回，因为数据已经获取成功
		fmt.Printf("failed to set cache for key %s: %v\n", key, err)
	}

	return data, nil
}

// cleanup 清理过期项
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if !item.expiration.IsZero() && now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
