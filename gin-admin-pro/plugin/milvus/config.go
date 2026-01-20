package milvus

import (
	"fmt"
	"time"
)

// Config Milvus配置
type Config struct {
	// Milvus地址
	Address string `yaml:"address" json:"address"`
	// 端口
	Port int `yaml:"port" json:"port"`
	// 用户名
	Username string `yaml:"username" json:"username"`
	// 密码
	Password string `yaml:"password" json:"password"`
	// 数据库名称
	Database string `yaml:"database" json:"database"`
	// 连接超时时间（秒）
	Timeout int `yaml:"timeout" json:"timeout"`
	// 最大重试次数
	MaxRetries int `yaml:"maxRetries" json:"maxRetries"`
	// 是否启用
	Enabled bool `yaml:"enabled" json:"enabled"`
}

// GetAddress 获取完整地址
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

// GetTimeout 获取超时时间
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Address:    "localhost",
		Port:       19530,
		Username:   "",
		Password:   "",
		Database:   "gin_admin",
		Timeout:    30,
		MaxRetries: 3,
		Enabled:    true,
	}
}

// Client Milvus客户端（简化版）
type Client struct {
	config *Config
}

// NewClient 创建Milvus客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if !cfg.Enabled {
		return &Client{config: cfg}, nil
	}

	// 在实际使用中，这里会创建真实的Milvus客户端
	// 现在返回一个模拟客户端
	return &Client{
		config: cfg,
	}, nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// IsEnabled 是否启用
func (c *Client) IsEnabled() bool {
	return c.config.Enabled
}

// Ping 测试连接
func (c *Client) Ping() error {
	if !c.IsEnabled() {
		return fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会发送ping请求到Milvus
	// 现在返回模拟的成功响应
	return nil
}

// ListCollections 列出所有集合
func (c *Client) ListCollections() ([]string, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("milvus is disabled")
	}

	// 返回模拟的集合列表
	return []string{
		"user_embeddings",
		"content_embeddings",
		"product_embeddings",
	}, nil
}

// CreateCollection 创建集合
func (c *Client) CreateCollection(collectionName string, dimension int) error {
	if !c.IsEnabled() {
		return fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会在Milvus中创建集合
	// 现在返回模拟的成功响应
	fmt.Printf("Creating collection: %s with dimension: %d\n", collectionName, dimension)
	return nil
}

// DropCollection 删除集合
func (c *Client) DropCollection(collectionName string) error {
	if !c.IsEnabled() {
		return fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会删除Milvus中的集合
	// 现在返回模拟的成功响应
	fmt.Printf("Dropping collection: %s\n", collectionName)
	return nil
}

// HasCollection 检查集合是否存在
func (c *Client) HasCollection(collectionName string) (bool, error) {
	if !c.IsEnabled() {
		return false, fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会检查Milvus中集合是否存在
	// 现在返回模拟响应
	collections, _ := c.ListCollections()
	for _, name := range collections {
		if name == collectionName {
			return true, nil
		}
	}
	return false, nil
}

// InsertData 插入向量数据
func (c *Client) InsertData(collectionName string, IDs []int64, vectors [][]float32) error {
	if !c.IsEnabled() {
		return fmt.Errorf("milvus is disabled")
	}

	if len(IDs) != len(vectors) {
		return fmt.Errorf("IDs and vectors length mismatch")
	}

	// 在实际使用中，这里会向Milvus插入向量数据
	// 现在返回模拟的成功响应
	fmt.Printf("Inserting %d vectors into collection: %s\n", len(vectors), collectionName)
	return nil
}

// SearchVectors 搜索向量
func (c *Client) SearchVectors(collectionName string, queryVectors [][]float32, topK int) ([]SearchResult, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("milvus is disabled")
	}

	if len(queryVectors) == 0 {
		return nil, fmt.Errorf("query vectors is empty")
	}

	// 在实际使用中，这里会在Milvus中搜索向量
	// 现在返回模拟的搜索结果
	var results []SearchResult

	for i, queryVector := range queryVectors {
		result := SearchResult{
			QueryIndex: i,
			Results: []VectorResult{
				{
					ID:       int64(i + 1),
					Score:    0.95,
					Vector:   queryVector,
					Metadata: map[string]interface{}{"title": fmt.Sprintf("Result %d", i+1)},
				},
			},
		}
		results = append(results, result)
	}

	return results, nil
}

// DeleteData 删除数据
func (c *Client) DeleteData(collectionName string, IDs []int64) error {
	if !c.IsEnabled() {
		return fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会从Milvus删除数据
	// 现在返回模拟的成功响应
	fmt.Printf("Deleting %d records from collection: %s\n", len(IDs), collectionName)
	return nil
}

// GetCollectionStats 获取集合统计信息
func (c *Client) GetCollectionStats(collectionName string) (*CollectionStats, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("milvus is disabled")
	}

	// 在实际使用中，这里会获取Milvus集合的统计信息
	// 现在返回模拟的统计信息
	return &CollectionStats{
		CollectionName: collectionName,
		RowCount:       1000,
		VectorCount:    1000,
		Dimension:      128,
		Size:           1024 * 1024, // 1MB
	}, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	QueryIndex int            `json:"queryIndex"`
	Results    []VectorResult `json:"results"`
}

// VectorResult 向量搜索结果
type VectorResult struct {
	ID       int64                  `json:"id"`
	Score    float32                `json:"score"`
	Vector   []float32              `json:"vector,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CollectionStats 集合统计信息
type CollectionStats struct {
	CollectionName string `json:"collectionName"`
	RowCount       int64  `json:"rowCount"`
	VectorCount    int64  `json:"vectorCount"`
	Dimension      int    `json:"dimension"`
	Size           int64  `json:"size"`
}

// Plugin Milvus插件
type Plugin struct {
	config *Config
	client *Client
}

// NewPlugin 创建Milvus插件
func NewPlugin(cfg *Config) *Plugin {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	client, _ := NewClient(cfg)

	return &Plugin{
		config: cfg,
		client: client,
	}
}

// GetConfig 获取配置
func (p *Plugin) GetConfig() *Config {
	return p.config
}

// GetClient 获取客户端
func (p *Plugin) GetClient() *Client {
	return p.client
}

// IsEnabled 是否启用
func (p *Plugin) IsEnabled() bool {
	return p.config.Enabled
}

// Init 初始化插件
func (p *Plugin) Init() error {
	if !p.IsEnabled() {
		return nil
	}

	// 测试连接
	if err := p.client.Ping(); err != nil {
		return fmt.Errorf("failed to connect to milvus: %w", err)
	}

	// 创建默认集合
	if err := p.createDefaultCollections(); err != nil {
		return fmt.Errorf("failed to create default collections: %w", err)
	}

	return nil
}

// createDefaultCollections 创建默认集合
func (p *Plugin) createDefaultCollections() error {
	// 创建用户嵌入向量集合
	if err := p.client.CreateCollection("user_embeddings", 128); err != nil {
		// 如果集合已存在，忽略错误
		fmt.Printf("Collection user_embeddings may already exist: %v\n", err)
	}

	// 创建内容嵌入向量集合
	if err := p.client.CreateCollection("content_embeddings", 768); err != nil {
		// 如果集合已存在，忽略错误
		fmt.Printf("Collection content_embeddings may already exist: %v\n", err)
	}

	// 创建商品嵌入向量集合
	if err := p.client.CreateCollection("product_embeddings", 256); err != nil {
		// 如果集合已存在，忽略错误
		fmt.Printf("Collection product_embeddings may already exist: %v\n", err)
	}

	return nil
}

// GetCollectionInfo 获取集合信息
func (p *Plugin) GetCollectionInfo() (map[string]interface{}, error) {
	if !p.IsEnabled() {
		return nil, fmt.Errorf("milvus is disabled")
	}

	collections, err := p.client.ListCollections()
	if err != nil {
		return nil, err
	}

	var collectionInfos []map[string]interface{}
	for _, name := range collections {
		stats, err := p.client.GetCollectionStats(name)
		if err != nil {
			continue
		}

		collectionInfos = append(collectionInfos, map[string]interface{}{
			"name":        stats.CollectionName,
			"rowCount":    stats.RowCount,
			"vectorCount": stats.VectorCount,
			"dimension":   stats.Dimension,
			"size":        stats.Size,
		})
	}

	return map[string]interface{}{
		"totalCollections": len(collections),
		"collections":      collectionInfos,
	}, nil
}
