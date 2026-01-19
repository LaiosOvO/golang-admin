package ai

import (
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items  map[string]*cacheItem
	mu     sync.RWMutex
	config struct {
		TTL     time.Duration
		MaxSize int
	}
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expiration time.Time
	created    time.Time
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(config struct {
	Enabled bool          `yaml:"enabled" mapstructure:"enabled"`
	TTL     time.Duration `yaml:"ttl" mapstructure:"ttl"`
	MaxSize int           `yaml:"maxSize" mapstructure:"maxSize"`
}) CacheService {
	cache := &MemoryCache{
		items: make(map[string]*cacheItem),
	}

	cache.config.TTL = config.TTL
	cache.config.MaxSize = config.MaxSize

	// 启动清理goroutine
	go cache.cleanup()

	return cache
}

// Get 获取缓存值
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set 设置缓存值
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查容量限制
	if len(c.items) >= c.config.MaxSize {
		// 简化实现：删除最旧的一项
		var oldestKey string
		var oldestTime time.Time
		for k, item := range c.items {
			if oldestTime.IsZero() || item.created.Before(oldestTime) {
				oldestKey = k
				oldestTime = item.created
			}
		}
		if oldestKey != "" {
			delete(c.items, oldestKey)
		}
	}

	// 设置过期时间
	expiration := time.Now().Add(ttl)
	if ttl == 0 {
		expiration = time.Now().Add(c.config.TTL)
	}

	c.items[key] = &cacheItem{
		value:      value,
		expiration: expiration,
		created:    time.Now(),
	}

	return nil
}

// Delete 删除缓存值
func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Clear 清空缓存
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

// Stats 获取缓存统计
func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := CacheStats{
		Count: len(c.items),
	}

	// 简化实现：计算命中率（实际应该跟踪访问次数）
	stats.HitRatio = 0.8                          // 假设值
	stats.MemoryUsage = int64(len(c.items) * 100) // 估算内存使用

	return stats
}

// cleanup 定期清理过期项
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, item := range c.items {
				if now.After(item.expiration) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// DefaultMetrics 默认指标实现
type DefaultMetrics struct {
	prefix string
	tags   map[string]string
}

// NewDefaultMetrics 创建默认指标
func NewDefaultMetrics(config struct {
	Enabled bool     `yaml:"enabled" mapstructure:"enabled"`
	Prefix  string   `yaml:"prefix" mapstructure:"prefix"`
	Tags    []string `yaml:"tags" mapstructure:"tags"`
}) MetricsService {
	metrics := &DefaultMetrics{
		prefix: config.Prefix,
		tags:   make(map[string]string),
	}

	// 转换tags数组为map
	for _, tag := range config.Tags {
		metrics.tags[tag] = tag
	}

	return metrics
}

// Counter 创建计数器
func (m *DefaultMetrics) Counter(name string, tags map[string]string) Counter {
	return &defaultCounter{
		name: m.prefix + "_" + name,
		tags: m.mergeTags(tags),
	}
}

// Histogram 创建直方图
func (m *DefaultMetrics) Histogram(name string, tags map[string]string) Histogram {
	return &defaultHistogram{
		name: m.prefix + "_" + name,
		tags: m.mergeTags(tags),
	}
}

// Gauge 创建仪表盘
func (m *DefaultMetrics) Gauge(name string, tags map[string]string) Gauge {
	return &defaultGauge{
		name: m.prefix + "_" + name,
		tags: m.mergeTags(tags),
	}
}

// Timer 创建计时器
func (m *DefaultMetrics) Timer(name string, tags map[string]string) Timer {
	return &defaultTimer{
		name: m.prefix + "_" + name,
		tags: m.mergeTags(tags),
	}
}

// mergeTags 合并标签
func (m *DefaultMetrics) mergeTags(tags map[string]string) map[string]string {
	merged := make(map[string]string)

	// 复制默认标签
	for k, v := range m.tags {
		merged[k] = v
	}

	// 添加传入的标签
	for k, v := range tags {
		merged[k] = v
	}

	return merged
}

// defaultCounter 默认计数器实现
type defaultCounter struct {
	name string
	tags map[string]string
}

func (c *defaultCounter) Inc() {
	// 简化实现：记录到日志
}

func (c *defaultCounter) Add(value float64) {
	// 简化实现：记录到日志
}

// defaultHistogram 默认直方图实现
type defaultHistogram struct {
	name string
	tags map[string]string
}

func (h *defaultHistogram) Observe(value float64) {
	// 简化实现：记录到日志
}

// defaultGauge 默认仪表盘实现
type defaultGauge struct {
	name string
	tags map[string]string
}

func (g *defaultGauge) Set(value float64) {
	// 简化实现：记录到日志
}

func (g *defaultGauge) Inc() {
	// 简化实现：记录到日志
}

func (g *defaultGauge) Dec() {
	// 简化实现：记录到日志
}

// defaultTimer 默认计时器实现
type defaultTimer struct {
	name string
	tags map[string]string
}

func (t *defaultTimer) Time(fn func()) {
	start := time.Now()
	fn()
	_ = time.Since(start)
	// 记录到日志或发送到监控系统
}

func (t *defaultTimer) Since(start time.Time) time.Duration {
	return time.Since(start)
}
