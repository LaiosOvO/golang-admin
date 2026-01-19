package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client MongoDB客户端
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	config   *Config
}

// NewClient 创建MongoDB客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 配置连接选项
	clientOptions := options.Client()
	clientOptions.ApplyURI(cfg.GetURI())

	// 连接池设置
	clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	clientOptions.SetMaxConnIdleTime(cfg.MaxConnIdle)

	// 超时设置
	clientOptions.SetConnectTimeout(cfg.ConnectTimeout)
	clientOptions.SetServerSelectionTimeout(cfg.ServerTimeout)

	// 压缩设置 - 暂时禁用压缩，因为API可能不同
	// if cfg.CompressLevel > 0 {
	//     compressor := options.CompressionSnappy
	//     clientOptions.SetCompressors([]options.Compressor{compressor})
	// }

	// 副本集设置
	if cfg.ReplicaSet != "" {
		clientOptions.SetReplicaSet(cfg.ReplicaSet)
	}

	// 创建客户端
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.Database)

	return &Client{
		client:   client,
		database: database,
		config:   cfg,
	}, nil
}

// GetClient 获取MongoDB客户端实例
func (c *Client) GetClient() *mongo.Client {
	return c.client
}

// GetDatabase 获取数据库实例
func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}

// GetCollection 获取集合实例
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	return c.client.Disconnect(ctx)
}

// Ping 检查数据库连接
func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	return c.client.Ping(ctx, readpref.Primary())
}

// ListCollections 列出所有集合
func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	collections, err := c.database.ListCollectionNames(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return collections, nil
}

// CreateCollection 创建集合
func (c *Client) CreateCollection(ctx context.Context, name string) error {
	return c.database.CreateCollection(ctx, name)
}

// DropCollection 删除集合
func (c *Client) DropCollection(ctx context.Context, name string) error {
	collection := c.database.Collection(name)
	return collection.Drop(ctx)
}

// HasCollection 检查集合是否存在
func (c *Client) HasCollection(ctx context.Context, name string) (bool, error) {
	collections, err := c.database.ListCollectionNames(ctx, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return false, err
	}
	for _, col := range collections {
		if col == name {
			return true, nil
		}
	}
	return false, nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, collection string, index mongo.IndexModel) error {
	col := c.database.Collection(collection)
	_, err := col.Indexes().CreateOne(ctx, index)
	return err
}

// CreateIndexes 创建多个索引
func (c *Client) CreateIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error {
	col := c.database.Collection(collection)
	_, err := col.Indexes().CreateMany(ctx, indexes)
	return err
}

// DropIndex 删除索引
func (c *Client) DropIndex(ctx context.Context, collection string, indexName string) error {
	col := c.database.Collection(collection)
	_, err := col.Indexes().DropOne(ctx, indexName)
	return err
}

// ListIndexes 列出索引
func (c *Client) ListIndexes(ctx context.Context, collection string) ([]string, error) {
	col := c.database.Collection(collection)
	cursor, err := col.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	var indexes []string
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			continue
		}
		if name, ok := index["name"].(string); ok {
			indexes = append(indexes, name)
		}
	}

	return indexes, cursor.Close(ctx)
}

// GetServerStatus 获取服务器状态
func (c *Client) GetServerStatus(ctx context.Context) (interface{}, error) {
	result := c.database.RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}})
	var status bson.M
	if err := result.Decode(&status); err != nil {
		return nil, err
	}
	return status, nil
}

// GetDatabaseStats 获取数据库统计信息
func (c *Client) GetDatabaseStats(ctx context.Context) (interface{}, error) {
	result := c.database.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}})
	var stats bson.M
	if err := result.Decode(&stats); err != nil {
		return nil, err
	}
	return stats, nil
}

// GetCollectionStats 获取集合统计信息
func (c *Client) GetCollectionStats(ctx context.Context, collection string) (interface{}, error) {
	result := c.database.RunCommand(ctx, bson.D{
		{Key: "collStats", Value: collection},
	})
	var stats bson.M
	if err := result.Decode(&stats); err != nil {
		return nil, err
	}
	return stats, nil
}
