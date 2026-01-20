# Milvus 插件

Milvus 插件提供向量数据库功能，支持向量存储、相似性搜索、向量分析等特性。当前实现为简化版本，可以根据实际需要扩展完整的Milvus功能。

## 功能特性

- 向量集合管理
- 向量数据插入和删除
- 相似性搜索
- 集合统计信息
- 连接健康检查
- 配置管理

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/milvus"

// 创建Milvus插件
milvusPlugin := milvus.NewPlugin(nil)

// 初始化插件
err := milvusPlugin.Init()
if err != nil {
    log.Fatal("Failed to init milvus plugin:", err)
}

// 获取客户端
client := milvusPlugin.GetClient()
```

### 2. 基础操作

```go
// 检查连接
err := client.Ping()

// 创建集合
err := client.CreateCollection("user_embeddings", 128)

// 检查集合是否存在
exists, err := client.HasCollection("user_embeddings")

// 列出所有集合
collections, err := client.ListCollections()

// 插入向量数据
ids := []int64{1, 2, 3}
vectors := [][]float32{
    {0.1, 0.2, 0.3, /* 125 more elements */},
    {0.4, 0.5, 0.6, /* 125 more elements */},
    {0.7, 0.8, 0.9, /* 125 more elements */},
}
err = client.InsertData("user_embeddings", ids, vectors)

// 搜索向量
queryVectors := [][]float32{
    {0.1, 0.2, 0.3, /* 125 more elements */},
}
results, err := client.SearchVectors("user_embeddings", queryVectors, 10)

// 获取集合统计信息
stats, err := client.GetCollectionStats("user_embeddings")

// 删除数据
deleteIDs := []int64{1, 2}
err = client.DeleteData("user_embeddings", deleteIDs)

// 删除集合
err = client.DropCollection("user_embeddings")
```

### 3. 获取插件信息

```go
// 获取集合信息
info, err := milvusPlugin.GetCollectionInfo()
```

## 配置说明

### 基础配置

```yaml
milvus:
  enabled: true            # 是否启用Milvus
  address: "localhost"     # Milvus地址
  port: 19530             # Milvus端口
  username: ""            # 用户名
  password: ""            # 密码
  database: "gin_admin"   # 数据库名称
  timeout: 30             # 连接超时时间（秒）
  maxRetries: 3           # 最大重试次数
```

### 自定义配置

```go
customConfig := &milvus.Config{
    Address:   "milvus.example.com",
    Port:      19530,
    Username:  "milvus",
    Password:  "password123",
    Database:  "myapp_db",
    Timeout:   60,
    MaxRetries: 5,
    Enabled:   true,
}

milvusPlugin := milvus.NewPlugin(customConfig)
```

## 集合设计

### 用户嵌入向量集合 (user_embeddings)

存储用户的特征向量，用于用户相似性分析：

- **维度**: 128
- **用途**: 用户推荐、用户聚类、相似用户查找

### 内容嵌入向量集合 (content_embeddings)

存储文章、评论等内容的向量，用于内容相似性分析：

- **维度**: 768
- **用途**: 内容推荐、相似内容查找、内容分类

### 商品嵌入向量集合 (product_embeddings)

存储商品的特征向量，用于商品相似性分析：

- **维度**: 256
- **用途**: 商品推荐、相似商品查找、商品分类

## 向量操作示例

### 1. 生成文本向量

```go
// 生成文本向量（需要集成文本嵌入模型）
func generateTextEmbedding(text string) ([]float32, error) {
    // 在实际使用中，这里会调用文本嵌入模型
    // 现在返回模拟向量
    vector := make([]float32, 768)
    for i := range vector {
        vector[i] = float32(i) / 1000.0
    }
    return vector, nil
}
```

### 2. 用户向量管理

```go
type UserVectorService struct {
    client *milvus.Client
}

func (s *UserVectorService) CreateUserVector(userID int64, features []float32) error {
    ids := []int64{userID}
    vectors := [][]float32{features}
    
    return s.client.InsertData("user_embeddings", ids, vectors)
}

func (s *UserVectorService) FindSimilarUsers(userID int64, topK int) ([]milvus.VectorResult, error) {
    // 首先获取用户向量（这里需要实现）
    userVector, err := s.getUserVector(userID)
    if err != nil {
        return nil, err
    }
    
    queryVectors := [][]float32{userVector}
    results, err := s.client.SearchVectors("user_embeddings", queryVectors, topK)
    if err != nil {
        return nil, err
    }
    
    if len(results) > 0 {
        return results[0].Results, nil
    }
    
    return nil, nil
}

func (s *UserVectorService) getUserVector(userID int64) ([]float32, error) {
    // 在实际使用中，这里需要从Milvus或缓存中获取用户向量
    // 现在返回模拟向量
    vector := make([]float32, 128)
    return vector, nil
}
```

### 3. 内容向量搜索

```go
type ContentVectorService struct {
    client *milvus.Client
}

func (s *ContentVectorService) SearchSimilarContent(content string, topK int) ([]milvus.VectorResult, error) {
    // 生成内容向量
    vector, err := generateTextEmbedding(content)
    if err != nil {
        return nil, err
    }
    
    queryVectors := [][]float32{vector}
    results, err := s.client.SearchVectors("content_embeddings", queryVectors, topK)
    if err != nil {
        return nil, err
    }
    
    if len(results) > 0 {
        return results[0].Results, nil
    }
    
    return nil, nil
}

func (s *ContentVectorService) IndexContent(contentID int64, content string) error {
    // 生成内容向量
    vector, err := generateTextEmbedding(content)
    if err != nil {
        return err
    }
    
    // 插入到向量数据库
    ids := []int64{contentID}
    vectors := [][]float32{vector}
    
    return s.client.InsertData("content_embeddings", ids, vectors)
}
```

## API接口规范

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/infra/milvus/health | 获取Milvus健康状态 |
| GET | /api/v1/infra/milvus/collections | 获取集合列表 |
| GET | /api/v1/infra/milvus/collection/{name}/stats | 获取集合统计信息 |
| POST | /api/v1/infra/milvus/collection/{name}/search | 搜索向量 |
| POST | /api/v1/infra/milvus/collection/{name}/insert | 插入向量 |
| DELETE | /api/v1/infra/milvus/collection/{name}/delete | 删除向量 |

## 向量搜索应用场景

### 1. 用户推荐系统

```go
// 基于用户相似性推荐内容
func RecommendContentByUser(userID int64, limit int) ([]int64, error) {
    // 找到相似用户
    similarUsers, err := userVectorService.FindSimilarUsers(userID, 10)
    if err != nil {
        return nil, err
    }
    
    // 获取相似用户喜欢的内容（这里需要实现）
    var recommendedContent []int64
    for _, user := range similarUsers {
        // 获取用户的内容偏好
        contentIDs := getUserLikedContent(user.ID)
        recommendedContent = append(recommendedContent, contentIDs...)
    }
    
    // 去重并限制数量
    return uniqueContent(recommendedContent, limit), nil
}
```

### 2. 内容相似性搜索

```go
// 查找相似文章
func FindSimilarArticles(articleID int64, limit int) ([]int64, error) {
    // 获取文章向量
    articleVector, err := getArticleVector(articleID)
    if err != nil {
        return nil, err
    }
    
    // 搜索相似文章
    queryVectors := [][]float32{articleVector}
    results, err := milvusClient.SearchVectors("content_embeddings", queryVectors, limit+1)
    if err != nil {
        return nil, err
    }
    
    if len(results) > 0 {
        var similarArticles []int64
        for _, result := range results[0].Results {
            // 排除自身
            if result.ID != articleID {
                similarArticles = append(similarArticles, result.ID)
            }
        }
        return similarArticles, nil
    }
    
    return nil, nil
}
```

### 3. 商品推荐

```go
// 基于商品相似性推荐
func RecommendSimilarProducts(productID int64, limit int) ([]int64, error) {
    // 获取商品向量
    productVector, err := getProductVector(productID)
    if err != nil {
        return nil, err
    }
    
    // 搜索相似商品
    queryVectors := [][]float32{productVector}
    results, err := milvusClient.SearchVectors("product_embeddings", queryVectors, limit+1)
    if err != nil {
        return nil, err
    }
    
    if len(results) > 0 {
        var similarProducts []int64
        for _, result := range results[0].Results {
            // 排除自身
            if result.ID != productID {
                similarProducts = append(similarProducts, result.ID)
            }
        }
        return similarProducts, nil
    }
    
    return nil, nil
}
```

## 性能优化

### 1. 向量索引

- 选择合适的索引类型（IVF_FLAT、IVF_PQ、HNSW等）
- 调整索引参数（nlist、nprobe等）
- 定期重建索引

### 2. 数据分区

- 按时间分区数据
- 按业务类型分区
- 使用Partition Key优化查询

### 3. 内存管理

- 合理设置缓存大小
- 监控内存使用情况
- 定期清理无用数据

## 监控指标

### 集合监控

- **数据量统计**：集合中的向量数量
- **存储空间**：集合占用的存储空间
- **查询性能**：搜索响应时间
- **索引状态**：索引构建和更新状态

### 系统监控

- **连接状态**：与Milvus的连接状态
- **内存使用**：Milvus服务内存使用情况
- **CPU使用**：Milvus服务CPU使用情况
- **磁盘IO**：磁盘读写性能

## 注意事项

1. **向量维度**：确保插入的向量维度与集合定义一致
2. **数据类型**：使用正确的数据类型（float32）
3. **ID唯一性**：确保向量ID的唯一性
4. **批量操作**：尽量使用批量操作提高性能
5. **索引选择**：根据数据特点选择合适的索引类型

## 最佳实践

1. **集合规划**：按业务模块合理规划集合
2. **维度选择**：根据模型输出确定向量维度
3. **索引优化**：根据查询特点优化索引参数
4. **数据清洗**：确保向量的质量和一致性
5. **性能监控**：定期监控查询性能和资源使用

## 扩展功能

当前实现为简化版本，可以根据需要扩展以下功能：

1. **完整客户端**：集成完整的Milvus Go客户端
2. **索引管理**：支持索引的创建、删除、更新
3. **分区管理**：支持数据分区功能
4. **向量预处理**：支持向量标准化、降维等预处理
5. **混合查询**：支持向量与标量的混合查询
6. **实时更新**：支持实时向量更新和增量索引

## 测试

```go
func TestMilvusClient(t *testing.T) {
    client := NewClient(DefaultConfig())
    
    // 测试连接
    err := client.Ping()
    assert.NoError(t, err)
    
    // 测试集合创建
    err = client.CreateCollection("test_collection", 128)
    assert.NoError(t, err)
    
    // 测试数据插入
    ids := []int64{1, 2, 3}
    vectors := [][]float32{
        {0.1, 0.2, 0.3},
        {0.4, 0.5, 0.6},
        {0.7, 0.8, 0.9},
    }
    err = client.InsertData("test_collection", ids, vectors)
    assert.NoError(t, err)
    
    // 测试向量搜索
    queryVectors := [][]float32{{0.1, 0.2, 0.3}}
    results, err := client.SearchVectors("test_collection", queryVectors, 10)
    assert.NoError(t, err)
    assert.NotEmpty(t, results)
    
    // 清理测试数据
    err = client.DropCollection("test_collection")
    assert.NoError(t, err)
}
```