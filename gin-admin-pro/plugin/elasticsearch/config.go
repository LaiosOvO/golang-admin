package elasticsearch

import (
	"fmt"
)

// Config Elasticsearch配置
type Config struct {
	// Elasticsearch地址
	Addresses []string `yaml:"addresses" json:"addresses"`
	// 用户名
	Username string `yaml:"username" json:"username"`
	// 密码
	Password string `yaml:"password" json:"password"`
	// 连接超时时间（秒）
	Timeout int `yaml:"timeout" json:"timeout"`
	// 最大重试次数
	MaxRetries int `yaml:"maxRetries" json:"maxRetries"`
	// 是否启用
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 默认索引前缀
	DefaultIndexPrefix string `yaml:"defaultIndexPrefix" json:"defaultIndexPrefix"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Addresses:          []string{"http://localhost:9200"},
		Username:           "",
		Password:           "",
		Timeout:            30,
		MaxRetries:         3,
		Enabled:            true,
		DefaultIndexPrefix: "gin_admin",
	}
}

// Client Elasticsearch客户端（简化版）
type Client struct {
	config *Config
}

// NewClient 创建Elasticsearch客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if !cfg.Enabled {
		return &Client{config: cfg}, nil
	}

	// 在实际使用中，这里会创建真实的Elasticsearch客户端
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
		return fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会发送ping请求到Elasticsearch
	// 现在返回模拟的成功响应
	return nil
}

// Health 检查健康状态
func (c *Client) Health() (map[string]interface{}, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("elasticsearch is disabled")
	}

	// 返回模拟的健康状态
	return map[string]interface{}{
		"status":                "green",
		"cluster_name":          "gin-admin-cluster",
		"number_of_nodes":       1,
		"active_primary_shards": 1,
		"active_shards":         1,
	}, nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(index string, mapping map[string]interface{}) error {
	if !c.IsEnabled() {
		return fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会创建Elasticsearch索引
	// 现在返回模拟的成功响应
	fmt.Printf("Creating index: %s with mapping: %v\n", index, mapping)
	return nil
}

// IndexDocument 索引文档
func (c *Client) IndexDocument(index string, docID string, doc interface{}) error {
	if !c.IsEnabled() {
		return fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会索引文档到Elasticsearch
	// 现在返回模拟的成功响应
	fmt.Printf("Indexing document %s to index %s: %v\n", docID, index, doc)
	return nil
}

// SearchDocuments 搜索文档
func (c *Client) SearchDocuments(index string, query map[string]interface{}) (*SearchResult, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会搜索Elasticsearch
	// 现在返回模拟的搜索结果
	return &SearchResult{
		Hits: Hits{
			Total: Total{
				Value:    0,
				Relation: "eq",
			},
			Hits: []Hit{},
		},
	}, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	Hits Hits `json:"hits"`
}

// Hits 命中结果
type Hits struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
}

// Total 总数
type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// Hit 单个命中结果
type Hit struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// GetDocument 获取文档
func (c *Client) GetDocument(index string, docID string, result interface{}) error {
	if !c.IsEnabled() {
		return fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会从Elasticsearch获取文档
	// 现在返回模拟响应
	fmt.Printf("Getting document %s from index %s\n", docID, index)
	return nil
}

// DeleteDocument 删除文档
func (c *Client) DeleteDocument(index string, docID string) error {
	if !c.IsEnabled() {
		return fmt.Errorf("elasticsearch is disabled")
	}

	// 在实际使用中，这里会从Elasticsearch删除文档
	// 现在返回模拟响应
	fmt.Printf("Deleting document %s from index %s\n", docID, index)
	return nil
}

// Plugin Elasticsearch插件
type Plugin struct {
	config *Config
	client *Client
}

// NewPlugin 创建Elasticsearch插件
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
		return fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}

	// 可以在这里创建默认索引
	return nil
}

// CreateDefaultIndexes 创建默认索引
func (p *Plugin) CreateDefaultIndexes() error {
	if !p.IsEnabled() {
		return nil
	}

	// 创建用户索引
	userMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"username": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"nickname": map[string]interface{}{
					"type": "text",
				},
				"email": map[string]interface{}{
					"type": "keyword",
				},
				"mobile": map[string]interface{}{
					"type": "keyword",
				},
				"status": map[string]interface{}{
					"type": "integer",
				},
				"dept_id": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	if err := p.client.CreateIndex(p.config.DefaultIndexPrefix+"_users", userMapping); err != nil {
		return fmt.Errorf("failed to create users index: %w", err)
	}

	// 创建操作日志索引
	operLogMapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"business_type": map[string]interface{}{
					"type": "integer",
				},
				"method": map[string]interface{}{
					"type": "keyword",
				},
				"oper_name": map[string]interface{}{
					"type": "keyword",
				},
				"dept_name": map[string]interface{}{
					"type": "keyword",
				},
				"oper_url": map[string]interface{}{
					"type": "keyword",
				},
				"oper_ip": map[string]interface{}{
					"type": "ip",
				},
				"oper_location": map[string]interface{}{
					"type": "keyword",
				},
				"status": map[string]interface{}{
					"type": "integer",
				},
				"error_msg": map[string]interface{}{
					"type": "text",
				},
				"oper_time": map[string]interface{}{
					"type": "date",
				},
				"cost_time": map[string]interface{}{
					"type": "long",
				},
			},
		},
	}

	if err := p.client.CreateIndex(p.config.DefaultIndexPrefix+"_oper_logs", operLogMapping); err != nil {
		return fmt.Errorf("failed to create oper_logs index: %w", err)
	}

	return nil
}

// GetHealthStatus 获取健康状态描述
func GetHealthStatus(status string) string {
	switch status {
	case "green":
		return "健康"
	case "yellow":
		return "警告"
	case "red":
		return "异常"
	default:
		return status
	}
}
