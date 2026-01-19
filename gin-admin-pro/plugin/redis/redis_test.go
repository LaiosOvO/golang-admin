package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "localhost:6379", cfg.Addr)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
	assert.Equal(t, "", cfg.Username)
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Equal(t, time.Millisecond*8, cfg.MinRetryBackoff)
	assert.Equal(t, time.Millisecond*512, cfg.MaxRetryBackoff)
	assert.Equal(t, time.Second*5, cfg.DialTimeout)
	assert.Equal(t, time.Second*3, cfg.ReadTimeout)
	assert.Equal(t, time.Second*3, cfg.WriteTimeout)
	assert.Equal(t, 10, cfg.PoolSize)
	assert.Equal(t, 5, cfg.MinIdleConns)
	assert.Equal(t, 100, cfg.MaxIdleConns)
	assert.Equal(t, time.Minute*5, cfg.ConnMaxIdleTime)
	assert.Equal(t, time.Hour, cfg.ConnMaxLifetime)
	assert.False(t, cfg.ClusterEnabled)
	assert.False(t, cfg.SentinelEnabled)
	assert.False(t, cfg.ShardEnabled)
}

func TestShardingAlgorithm(t *testing.T) {
	sharding := ShardingAlgorithm{
		Type:    "consistent",
		KeyFunc: "key",
		Options: map[string]interface{}{
			"replicas": 150,
		},
	}

	assert.Equal(t, "consistent", sharding.Type)
	assert.Equal(t, "key", sharding.KeyFunc)
	assert.Equal(t, 150, sharding.Options["replicas"])
}

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test-key"
		value := "test-value"

		err := cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		retrieved, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrieved)
	})

	t.Run("Set and GetBytes", func(t *testing.T) {
		key := "test-bytes"
		value := []byte("test-data")

		err := cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		retrieved, err := cache.GetBytes(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrieved)
	})

	t.Run("SetJSON and GetJSON", func(t *testing.T) {
		key := "test-json"
		obj := map[string]interface{}{
			"name":  "test",
			"value": 123,
		}

		err := cache.SetJSON(ctx, key, obj, time.Minute)
		assert.NoError(t, err)

		var retrieved map[string]interface{}
		err = cache.GetJSON(ctx, key, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, obj["name"], retrieved["name"])
		assert.Equal(t, float64(123), retrieved["value"]) // JSON numbers are float64 by default
	})

	t.Run("Exists", func(t *testing.T) {
		key := "test-exists"
		value := "test-value"

		// 不存在的键
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)

		// 设置后
		err = cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		exists, err = cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test-delete"
		value := "test-value"

		err := cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		err = cache.Del(ctx, key)
		assert.NoError(t, err)

		exists, err = cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Expiration", func(t *testing.T) {
		key := "test-expiration"
		value := "test-value"

		err := cache.Set(ctx, key, value, time.Millisecond*100)
		assert.NoError(t, err)

		// 立即检查应该存在
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		// 等待过期
		time.Sleep(time.Millisecond * 150)

		exists, err = cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TTL", func(t *testing.T) {
		key := "test-ttl"
		value := "test-value"

		err := cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		ttl, err := cache.TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ttl > 0)
		assert.True(t, ttl <= time.Minute)

		// 不存在的键
		ttl, err = cache.TTL(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(-1), ttl)
	})
}

func TestGetOrSet(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache()

	t.Run("Cache miss", func(t *testing.T) {
		key := "test-or-set"
		called := false

		_, err := cache.Get(ctx, key)
		assert.Error(t, err) // 应该报错，因为键不存在

		result, err := cache.GetOrSet(ctx, key, func() (interface{}, error) {
			called = true
			return "callback-value", nil
		}, time.Minute)

		assert.NoError(t, err)
		assert.Equal(t, "callback-value", result)
		assert.True(t, called)
	})

	t.Run("Cache hit", func(t *testing.T) {
		key := "test-or-set-hit"

		// 先设置值
		err := cache.Set(ctx, key, "existing-value", time.Minute)
		assert.NoError(t, err)

		called := false

		result, err := cache.GetOrSet(ctx, key, func() (interface{}, error) {
			called = true
			return "callback-value", nil
		}, time.Minute)

		assert.NoError(t, err)
		assert.Equal(t, "existing-value", result)
		assert.False(t, called) // 回调不应该被调用
	})
}
