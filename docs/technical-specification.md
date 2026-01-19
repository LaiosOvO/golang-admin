# Golang Admin 技术规范文档

> 版本：v1.0.0
> 更新日期：2026-01-10

---

## 目录

1. [项目总览](#1-项目总览)
2. [单体架构 Admin (go-admin-mono)](#2-单体架构-admin-go-admin-mono)
3. [微服务架构 Admin (go-admin-micro)](#3-微服务架构-admin-go-admin-micro)
4. [Starter 组件库 (go-admin-starter)](#4-starter-组件库-go-admin-starter)
5. [AI 语音交互项目 (go-xiaozhi)](#5-ai-语音交互项目-go-xiaozhi)
6. [数据库封装规范](#6-数据库封装规范)
7. [开发路线图](#7-开发路线图)

---

## 1. 项目总览

### 1.1 项目矩阵

| 项目名称 | 类型 | 描述 |
|----------|------|------|
| **go-admin-mono** | 单体应用 | 完整的单体 Admin 后台系统 |
| **go-admin-micro** | 微服务应用 | 分布式微服务 Admin 系统 |
| **go-admin-starter** | 组件库 | 类似 ruoyi-vue-pro 的 starter 组件 |
| **go-xiaozhi** | AI 应用 | 类似 xiaozhi-server 的语音交互系统 |

### 1.2 技术栈总览

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端层 (可选)                              │
│              Vue3 / React / 纯 API 服务                          │
├─────────────────────────────────────────────────────────────────┤
│                        网关层                                     │
│                    Gin / Hertz / gRPC                            │
├─────────────────────────────────────────────────────────────────┤
│                        业务层                                     │
│         单体: Gin + GORM    微服务: go-zero / kratos             │
├─────────────────────────────────────────────────────────────────┤
│                      Starter 组件层                               │
│    db-starter | cache-starter | mq-starter | search-starter     │
├─────────────────────────────────────────────────────────────────┤
│                        数据层                                     │
│  MySQL | PostgreSQL | MongoDB | Redis | ES | Milvus | Kafka     │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 核心设计原则

1. **插件化设计**：所有数据库操作通过 Starter 按需加载
2. **配置驱动**：YAML/Nacos 动态配置，零代码切换
3. **统一抽象**：定义标准接口，屏蔽底层实现差异
4. **代码生成**：提供 CLI 工具自动生成 CRUD 代码

---

## 2. 单体架构 Admin (go-admin-mono)

### 2.1 项目结构

```
go-admin-mono/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── config/                     # 配置管理
│   │   ├── config.go
│   │   └── config.yaml
│   ├── handler/                    # HTTP 处理器 (Controller)
│   │   ├── system/                 # 系统管理模块
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── menu.go
│   │   │   └── dept.go
│   │   ├── monitor/                # 监控模块
│   │   └── tool/                   # 工具模块
│   ├── service/                    # 业务逻辑层
│   │   ├── system/
│   │   └── ...
│   ├── repository/                 # 数据访问层
│   │   ├── mysql/
│   │   ├── postgres/
│   │   ├── mongo/
│   │   └── ...
│   ├── model/                      # 数据模型
│   │   ├── entity/                 # 数据库实体
│   │   ├── dto/                    # 数据传输对象
│   │   └── vo/                     # 视图对象
│   ├── middleware/                 # 中间件
│   │   ├── auth.go                 # JWT 认证
│   │   ├── permission.go           # RBAC 权限
│   │   ├── ratelimit.go            # 限流
│   │   └── logger.go               # 日志
│   └── pkg/                        # 内部工具包
│       ├── response/               # 统一响应
│       ├── validator/              # 参数校验
│       └── utils/                  # 工具函数
├── pkg/                            # 可导出的公共包
├── api/                            # API 定义 (OpenAPI/Swagger)
├── scripts/                        # 脚本
├── deployments/                    # 部署配置
│   ├── docker/
│   └── k8s/
├── go.mod
└── go.sum
```

### 2.2 核心模块设计

#### 2.2.1 系统管理模块

```go
// 用户管理
type UserHandler interface {
    Create(c *gin.Context)      // 创建用户
    Update(c *gin.Context)      // 更新用户
    Delete(c *gin.Context)      // 删除用户
    Get(c *gin.Context)         // 获取用户详情
    List(c *gin.Context)        // 用户列表 (分页)
    UpdateStatus(c *gin.Context) // 更新状态
    ResetPassword(c *gin.Context) // 重置密码
    Export(c *gin.Context)      // 导出 Excel
    Import(c *gin.Context)      // 导入 Excel
}

// 角色管理
type RoleHandler interface {
    Create(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    Get(c *gin.Context)
    List(c *gin.Context)
    UpdateStatus(c *gin.Context)
    AssignPermissions(c *gin.Context)  // 分配权限
    AssignUsers(c *gin.Context)        // 分配用户
}

// 菜单管理
type MenuHandler interface {
    Create(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    Get(c *gin.Context)
    Tree(c *gin.Context)        // 菜单树
    RoleMenuTree(c *gin.Context) // 角色菜单树
}

// 部门管理
type DeptHandler interface {
    Create(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    Get(c *gin.Context)
    Tree(c *gin.Context)        // 部门树
    UserTree(c *gin.Context)    // 部门用户树
}
```

#### 2.2.2 权限设计 (RBAC)

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  User   │────▶│  Role   │────▶│  Menu   │
└─────────┘     └─────────┘     └─────────┘
     │               │               │
     │               │               │
     ▼               ▼               ▼
┌─────────┐     ┌─────────┐     ┌─────────┐
│  Dept   │     │ DataScope│    │Permission│
└─────────┘     └─────────┘     └─────────┘
```

**数据权限范围 (DataScope)**：
- 全部数据权限
- 自定义数据权限
- 本部门数据权限
- 本部门及以下数据权限
- 仅本人数据权限

#### 2.2.3 配置示例

```yaml
# config.yaml
server:
  port: 8080
  mode: debug  # debug, release, test

# 数据库配置 - 按需启用
database:
  mysql:
    enabled: true
    host: localhost
    port: 3306
    username: root
    password: root
    database: go_admin
    max_idle_conns: 10
    max_open_conns: 100

  postgres:
    enabled: true
    host: localhost
    port: 5432
    username: postgres
    password: postgres
    database: go_admin
    extensions:
      - postgis       # GIS 插件
      - vector        # 向量插件

  mongodb:
    enabled: true
    uri: mongodb://localhost:27017
    database: go_admin

  redis:
    enabled: true
    host: localhost
    port: 6379
    password: ""
    db: 0

  elasticsearch:
    enabled: true
    addresses:
      - http://localhost:9200
    username: ""
    password: ""

  milvus:
    enabled: true
    host: localhost
    port: 19530

# 消息队列
mq:
  kafka:
    enabled: true
    brokers:
      - localhost:9092
    group_id: go-admin-group

# 定时任务
scheduler:
  enabled: true
  type: asynq  # gocron, asynq, xxl-job
  redis:
    addr: localhost:6379

# JWT 配置
jwt:
  secret: your-secret-key
  expire: 24h
  refresh_expire: 168h
```

### 2.3 经典案例封装

#### 案例1：MySQL CRUD 操作

```go
// repository/mysql/user_repo.go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// 创建用户
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// 分页查询
func (r *UserRepository) Page(ctx context.Context, query *dto.UserQuery) (*dto.PageResult[entity.User], error) {
    var users []entity.User
    var total int64

    db := r.db.WithContext(ctx).Model(&entity.User{})

    // 动态条件构建
    if query.Username != "" {
        db = db.Where("username LIKE ?", "%"+query.Username+"%")
    }
    if query.Status != nil {
        db = db.Where("status = ?", *query.Status)
    }
    if query.DeptID != 0 {
        db = db.Where("dept_id = ?", query.DeptID)
    }

    // 数据权限过滤
    db = r.applyDataScope(db, ctx)

    // 统计总数
    if err := db.Count(&total).Error; err != nil {
        return nil, err
    }

    // 分页查询
    offset := (query.PageNum - 1) * query.PageSize
    if err := db.Offset(offset).Limit(query.PageSize).Find(&users).Error; err != nil {
        return nil, err
    }

    return &dto.PageResult[entity.User]{
        List:     users,
        Total:    total,
        PageNum:  query.PageNum,
        PageSize: query.PageSize,
    }, nil
}
```

#### 案例2：PostgreSQL + PostGIS 地理查询

```go
// repository/postgres/location_repo.go
type LocationRepository struct {
    db *gorm.DB
}

// 查询附近的点 (指定半径内)
func (r *LocationRepository) FindNearby(ctx context.Context, lat, lng, radiusMeters float64) ([]entity.Location, error) {
    var locations []entity.Location

    // 使用 PostGIS 的 ST_DWithin 函数
    sql := `
        SELECT id, name, description,
               ST_X(geom::geometry) as longitude,
               ST_Y(geom::geometry) as latitude,
               ST_Distance(geom, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) as distance
        FROM locations
        WHERE ST_DWithin(
            geom,
            ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
            ?
        )
        ORDER BY distance
    `

    err := r.db.WithContext(ctx).Raw(sql, lng, lat, lng, lat, radiusMeters).Scan(&locations).Error
    return locations, err
}

// 区域内查询 (多边形)
func (r *LocationRepository) FindInPolygon(ctx context.Context, polygon string) ([]entity.Location, error) {
    var locations []entity.Location

    sql := `
        SELECT * FROM locations
        WHERE ST_Within(
            geom::geometry,
            ST_GeomFromGeoJSON(?)
        )
    `

    err := r.db.WithContext(ctx).Raw(sql, polygon).Scan(&locations).Error
    return locations, err
}
```

#### 案例3：PostgreSQL + pgvector 向量检索

```go
// repository/postgres/embedding_repo.go
type EmbeddingRepository struct {
    db *gorm.DB
}

// 向量相似度搜索
func (r *EmbeddingRepository) SimilaritySearch(ctx context.Context, embedding []float32, topK int) ([]entity.Document, error) {
    var docs []entity.Document

    // 将 embedding 转换为 pgvector 格式
    vectorStr := fmt.Sprintf("[%s]", floatsToString(embedding))

    sql := `
        SELECT id, content, metadata,
               1 - (embedding <=> ?::vector) as similarity
        FROM documents
        ORDER BY embedding <=> ?::vector
        LIMIT ?
    `

    err := r.db.WithContext(ctx).Raw(sql, vectorStr, vectorStr, topK).Scan(&docs).Error
    return docs, err
}

// 混合搜索 (向量 + 关键词)
func (r *EmbeddingRepository) HybridSearch(ctx context.Context, embedding []float32, keyword string, topK int) ([]entity.Document, error) {
    var docs []entity.Document

    vectorStr := fmt.Sprintf("[%s]", floatsToString(embedding))

    sql := `
        WITH vector_search AS (
            SELECT id, content, metadata,
                   1 - (embedding <=> ?::vector) as vector_score
            FROM documents
            ORDER BY embedding <=> ?::vector
            LIMIT ?
        ),
        keyword_search AS (
            SELECT id, content, metadata,
                   ts_rank(to_tsvector('chinese', content), plainto_tsquery('chinese', ?)) as text_score
            FROM documents
            WHERE to_tsvector('chinese', content) @@ plainto_tsquery('chinese', ?)
            LIMIT ?
        )
        SELECT COALESCE(v.id, k.id) as id,
               COALESCE(v.content, k.content) as content,
               COALESCE(v.metadata, k.metadata) as metadata,
               COALESCE(v.vector_score, 0) * 0.7 + COALESCE(k.text_score, 0) * 0.3 as score
        FROM vector_search v
        FULL OUTER JOIN keyword_search k ON v.id = k.id
        ORDER BY score DESC
        LIMIT ?
    `

    err := r.db.WithContext(ctx).Raw(sql, vectorStr, vectorStr, topK*2, keyword, keyword, topK*2, topK).Scan(&docs).Error
    return docs, err
}
```

#### 案例4：MongoDB 文档操作

```go
// repository/mongo/audit_log_repo.go
type AuditLogRepository struct {
    collection *mongo.Collection
}

// 创建审计日志
func (r *AuditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
    log.ID = primitive.NewObjectID()
    log.CreatedAt = time.Now()
    _, err := r.collection.InsertOne(ctx, log)
    return err
}

// 复杂聚合查询 - 统计操作类型分布
func (r *AuditLogRepository) AggregateByOperationType(ctx context.Context, startTime, endTime time.Time) ([]dto.OperationStats, error) {
    pipeline := mongo.Pipeline{
        // 时间范围过滤
        {{Key: "$match", Value: bson.D{
            {Key: "created_at", Value: bson.D{
                {Key: "$gte", Value: startTime},
                {Key: "$lte", Value: endTime},
            }},
        }}},
        // 按操作类型分组
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$operation_type"},
            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
            {Key: "users", Value: bson.D{{Key: "$addToSet", Value: "$user_id"}}},
        }}},
        // 计算独立用户数
        {{Key: "$project", Value: bson.D{
            {Key: "operation_type", Value: "$_id"},
            {Key: "count", Value: 1},
            {Key: "unique_users", Value: bson.D{{Key: "$size", Value: "$users"}}},
        }}},
        // 排序
        {{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []dto.OperationStats
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}
```

#### 案例5：Redis 缓存操作

```go
// repository/redis/cache_repo.go
type CacheRepository struct {
    client *redis.Client
}

// 缓存用户信息
func (r *CacheRepository) SetUser(ctx context.Context, userID int64, user *entity.User, ttl time.Duration) error {
    key := fmt.Sprintf("user:%d", userID)
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, data, ttl).Err()
}

// 分布式锁
func (r *CacheRepository) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    lockKey := fmt.Sprintf("lock:%s", key)
    return r.client.SetNX(ctx, lockKey, "1", ttl).Result()
}

func (r *CacheRepository) Unlock(ctx context.Context, key string) error {
    lockKey := fmt.Sprintf("lock:%s", key)
    return r.client.Del(ctx, lockKey).Err()
}

// 限流器 (滑动窗口)
func (r *CacheRepository) IsRateLimited(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    now := time.Now().UnixNano()
    windowStart := now - int64(window)

    pipe := r.client.Pipeline()

    // 移除窗口外的记录
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    // 添加当前请求
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
    // 获取窗口内请求数
    countCmd := pipe.ZCard(ctx, key)
    // 设置过期时间
    pipe.Expire(ctx, key, window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    return countCmd.Val() > int64(limit), nil
}

// 发布订阅
func (r *CacheRepository) Publish(ctx context.Context, channel string, message interface{}) error {
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }
    return r.client.Publish(ctx, channel, data).Err()
}

func (r *CacheRepository) Subscribe(ctx context.Context, channel string, handler func(message string)) error {
    pubsub := r.client.Subscribe(ctx, channel)
    defer pubsub.Close()

    ch := pubsub.Channel()
    for msg := range ch {
        handler(msg.Payload)
    }
    return nil
}
```

#### 案例6：Elasticsearch 搜索操作

```go
// repository/es/search_repo.go
type SearchRepository struct {
    client *elasticsearch.Client
    index  string
}

// 全文搜索
func (r *SearchRepository) Search(ctx context.Context, query *dto.SearchQuery) (*dto.SearchResult, error) {
    var buf bytes.Buffer

    searchBody := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "multi_match": map[string]interface{}{
                            "query":  query.Keyword,
                            "fields": []string{"title^2", "content", "tags"},
                            "type":   "best_fields",
                        },
                    },
                },
                "filter": r.buildFilters(query),
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "title":   map[string]interface{}{},
                "content": map[string]interface{}{"fragment_size": 150},
            },
            "pre_tags":  []string{"<em>"},
            "post_tags": []string{"</em>"},
        },
        "from": (query.PageNum - 1) * query.PageSize,
        "size": query.PageSize,
        "sort": r.buildSort(query),
    }

    if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
        return nil, err
    }

    res, err := r.client.Search(
        r.client.Search.WithContext(ctx),
        r.client.Search.WithIndex(r.index),
        r.client.Search.WithBody(&buf),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    var result dto.ESResponse
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, err
    }

    return r.parseSearchResult(&result), nil
}

// 聚合统计
func (r *SearchRepository) Aggregate(ctx context.Context, field string) (map[string]int64, error) {
    var buf bytes.Buffer

    aggBody := map[string]interface{}{
        "size": 0,
        "aggs": map[string]interface{}{
            "field_agg": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": field,
                    "size":  100,
                },
            },
        },
    }

    if err := json.NewEncoder(&buf).Encode(aggBody); err != nil {
        return nil, err
    }

    res, err := r.client.Search(
        r.client.Search.WithContext(ctx),
        r.client.Search.WithIndex(r.index),
        r.client.Search.WithBody(&buf),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    // 解析聚合结果
    var result map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, err
    }

    return r.parseAggregation(result), nil
}
```

#### 案例7：Milvus 向量数据库操作

```go
// repository/milvus/vector_repo.go
type VectorRepository struct {
    client client.Client
}

// 创建集合
func (r *VectorRepository) CreateCollection(ctx context.Context, collectionName string, dim int) error {
    schema := &entity.Schema{
        CollectionName: collectionName,
        Fields: []*entity.Field{
            {
                Name:       "id",
                DataType:   entity.FieldTypeInt64,
                PrimaryKey: true,
                AutoID:     true,
            },
            {
                Name:     "content",
                DataType: entity.FieldTypeVarChar,
                TypeParams: map[string]string{
                    "max_length": "65535",
                },
            },
            {
                Name:     "embedding",
                DataType: entity.FieldTypeFloatVector,
                TypeParams: map[string]string{
                    "dim": fmt.Sprintf("%d", dim),
                },
            },
            {
                Name:     "metadata",
                DataType: entity.FieldTypeJSON,
            },
        },
    }

    return r.client.CreateCollection(ctx, schema, 2)
}

// 插入向量
func (r *VectorRepository) Insert(ctx context.Context, collectionName string, docs []dto.VectorDoc) error {
    contents := make([]string, len(docs))
    embeddings := make([][]float32, len(docs))
    metadatas := make([][]byte, len(docs))

    for i, doc := range docs {
        contents[i] = doc.Content
        embeddings[i] = doc.Embedding
        metadata, _ := json.Marshal(doc.Metadata)
        metadatas[i] = metadata
    }

    contentColumn := entity.NewColumnVarChar("content", contents)
    embeddingColumn := entity.NewColumnFloatVector("embedding", len(docs[0].Embedding), embeddings)
    metadataColumn := entity.NewColumnJSONBytes("metadata", metadatas)

    _, err := r.client.Insert(ctx, collectionName, "", contentColumn, embeddingColumn, metadataColumn)
    return err
}

// 向量搜索
func (r *VectorRepository) Search(ctx context.Context, collectionName string, embedding []float32, topK int, filter string) ([]dto.VectorSearchResult, error) {
    // 加载集合到内存
    if err := r.client.LoadCollection(ctx, collectionName, false); err != nil {
        return nil, err
    }

    sp, _ := entity.NewIndexIvfFlatSearchParam(16)

    results, err := r.client.Search(
        ctx,
        collectionName,
        nil,
        filter,
        []string{"content", "metadata"},
        []entity.Vector{entity.FloatVector(embedding)},
        "embedding",
        entity.L2,
        topK,
        sp,
    )
    if err != nil {
        return nil, err
    }

    var searchResults []dto.VectorSearchResult
    for _, result := range results {
        for i := 0; i < result.ResultCount; i++ {
            searchResults = append(searchResults, dto.VectorSearchResult{
                ID:       result.IDs.(*entity.ColumnInt64).Data()[i],
                Score:    result.Scores[i],
                Content:  result.Fields.GetColumn("content").(*entity.ColumnVarChar).Data()[i],
                Metadata: result.Fields.GetColumn("metadata").(*entity.ColumnJSONBytes).Data()[i],
            })
        }
    }

    return searchResults, nil
}
```

#### 案例8：Kafka 消息队列操作

```go
// repository/kafka/mq_repo.go
type MQRepository struct {
    producer sarama.SyncProducer
    consumer sarama.ConsumerGroup
}

// 发送消息
func (r *MQRepository) SendMessage(ctx context.Context, topic string, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(data),
        Headers: []sarama.RecordHeader{
            {Key: []byte("trace_id"), Value: []byte(ctx.Value("trace_id").(string))},
        },
    }

    _, _, err = r.producer.SendMessage(msg)
    return err
}

// 批量发送
func (r *MQRepository) SendMessages(ctx context.Context, topic string, messages []dto.MQMessage) error {
    var msgs []*sarama.ProducerMessage

    for _, m := range messages {
        data, _ := json.Marshal(m.Value)
        msgs = append(msgs, &sarama.ProducerMessage{
            Topic: topic,
            Key:   sarama.StringEncoder(m.Key),
            Value: sarama.ByteEncoder(data),
        })
    }

    return r.producer.SendMessages(msgs)
}

// 消费者处理器
type ConsumerHandler struct {
    handler func(message *sarama.ConsumerMessage) error
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        if err := h.handler(msg); err != nil {
            // 记录错误，可以发送到死信队列
            log.Printf("消息处理失败: %v", err)
            continue
        }
        session.MarkMessage(msg, "")
    }
    return nil
}

// 订阅消费
func (r *MQRepository) Subscribe(ctx context.Context, topics []string, handler func(message *sarama.ConsumerMessage) error) error {
    h := &ConsumerHandler{handler: handler}

    for {
        select {
        case <-ctx.Done():
            return nil
        default:
            if err := r.consumer.Consume(ctx, topics, h); err != nil {
                log.Printf("消费错误: %v", err)
            }
        }
    }
}
```

---

## 3. 微服务架构 Admin (go-admin-micro)

### 3.1 服务拆分

```
go-admin-micro/
├── services/
│   ├── gateway/                    # API 网关
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── etc/
│   ├── system/                     # 系统服务
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── logic/
│   │   │   ├── svc/
│   │   │   └── server/
│   │   └── etc/
│   ├── auth/                       # 认证服务
│   ├── user/                       # 用户服务
│   ├── log/                        # 日志服务
│   ├── file/                       # 文件服务
│   ├── message/                    # 消息服务
│   ├── scheduler/                  # 调度服务
│   ├── search/                     # 搜索服务
│   └── ai/                         # AI 服务
├── pkg/                            # 公共包
│   ├── proto/                      # protobuf 定义
│   ├── interceptor/                # gRPC 拦截器
│   └── middleware/                 # HTTP 中间件
├── deploy/
│   ├── docker-compose/
│   └── k8s/
└── Makefile
```

### 3.2 服务架构图

```
                                    ┌─────────────┐
                                    │   Client    │
                                    └──────┬──────┘
                                           │
                                    ┌──────▼──────┐
                                    │   Gateway   │
                                    │  (Gin/gRPC) │
                                    └──────┬──────┘
                                           │
              ┌────────────────────────────┼────────────────────────────┐
              │                            │                            │
       ┌──────▼──────┐              ┌──────▼──────┐              ┌──────▼──────┐
       │    Auth     │              │   System    │              │    User     │
       │   Service   │              │   Service   │              │   Service   │
       └──────┬──────┘              └──────┬──────┘              └──────┬──────┘
              │                            │                            │
              └────────────────────────────┼────────────────────────────┘
                                           │
       ┌───────────────────────────────────┼───────────────────────────────────┐
       │                                   │                                   │
┌──────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐
│    MySQL    │   │  PostgreSQL │   │   MongoDB   │   │    Redis    │   │    Kafka    │
└─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘
```

### 3.3 基于 go-zero 的服务示例

#### 3.3.1 API 定义 (user.api)

```go
// services/user/api/user.api
syntax = "v1"

info (
    title: "用户服务"
    desc: "用户管理相关接口"
    version: "v1.0.0"
)

type (
    CreateUserReq {
        Username string `json:"username" validate:"required,min=3,max=50"`
        Password string `json:"password" validate:"required,min=6"`
        Nickname string `json:"nickname"`
        Email    string `json:"email" validate:"email"`
        Phone    string `json:"phone"`
        DeptId   int64  `json:"deptId"`
        RoleIds  []int64 `json:"roleIds"`
    }

    CreateUserResp {
        Id int64 `json:"id"`
    }

    UserInfo {
        Id       int64   `json:"id"`
        Username string  `json:"username"`
        Nickname string  `json:"nickname"`
        Email    string  `json:"email"`
        Phone    string  `json:"phone"`
        Status   int     `json:"status"`
        DeptId   int64   `json:"deptId"`
        Roles    []Role  `json:"roles"`
    }

    Role {
        Id   int64  `json:"id"`
        Name string `json:"name"`
        Code string `json:"code"`
    }

    PageReq {
        PageNum  int `form:"pageNum,default=1"`
        PageSize int `form:"pageSize,default=10"`
    }

    UserListReq {
        PageReq
        Username string `form:"username,optional"`
        Status   *int   `form:"status,optional"`
        DeptId   int64  `form:"deptId,optional"`
    }

    UserListResp {
        List     []UserInfo `json:"list"`
        Total    int64      `json:"total"`
        PageNum  int        `json:"pageNum"`
        PageSize int        `json:"pageSize"`
    }
)

@server (
    prefix: /api/v1/user
    group: user
    middleware: AuthMiddleware
)
service user-api {
    @doc "创建用户"
    @handler CreateUser
    post / (CreateUserReq) returns (CreateUserResp)

    @doc "获取用户详情"
    @handler GetUser
    get /:id returns (UserInfo)

    @doc "获取用户列表"
    @handler ListUser
    get /list (UserListReq) returns (UserListResp)

    @doc "更新用户"
    @handler UpdateUser
    put /:id (CreateUserReq)

    @doc "删除用户"
    @handler DeleteUser
    delete /:id
}
```

#### 3.3.2 RPC 定义 (user.proto)

```protobuf
// services/user/rpc/pb/user.proto
syntax = "proto3";

package user;

option go_package = "./user";

message CreateUserReq {
    string username = 1;
    string password = 2;
    string nickname = 3;
    string email = 4;
    string phone = 5;
    int64 dept_id = 6;
    repeated int64 role_ids = 7;
}

message CreateUserResp {
    int64 id = 1;
}

message GetUserReq {
    int64 id = 1;
}

message UserInfo {
    int64 id = 1;
    string username = 2;
    string nickname = 3;
    string email = 4;
    string phone = 5;
    int32 status = 6;
    int64 dept_id = 7;
    repeated Role roles = 8;
}

message Role {
    int64 id = 1;
    string name = 2;
    string code = 3;
}

message ListUserReq {
    int32 page_num = 1;
    int32 page_size = 2;
    string username = 3;
    optional int32 status = 4;
    int64 dept_id = 5;
}

message ListUserResp {
    repeated UserInfo list = 1;
    int64 total = 2;
    int32 page_num = 3;
    int32 page_size = 4;
}

message Empty {}

service UserService {
    rpc CreateUser(CreateUserReq) returns (CreateUserResp);
    rpc GetUser(GetUserReq) returns (UserInfo);
    rpc ListUser(ListUserReq) returns (ListUserResp);
    rpc UpdateUser(UserInfo) returns (Empty);
    rpc DeleteUser(GetUserReq) returns (Empty);
}
```

#### 3.3.3 服务配置

```yaml
# services/user/etc/user.yaml
Name: user-rpc
ListenOn: 0.0.0.0:8081

Etcd:
  Hosts:
    - localhost:2379
  Key: user.rpc

Mysql:
  DataSource: root:root@tcp(localhost:3306)/go_admin?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

Telemetry:
  Name: user-rpc
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger
```

### 3.4 服务间通信

```go
// services/gateway/internal/svc/servicecontext.go
type ServiceContext struct {
    Config      config.Config
    UserRpc     userclient.UserService
    AuthRpc     authclient.AuthService
    SystemRpc   systemclient.SystemService
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:      c,
        UserRpc:     userclient.NewUserService(zrpc.MustNewClient(c.UserRpc)),
        AuthRpc:     authclient.NewAuthService(zrpc.MustNewClient(c.AuthRpc)),
        SystemRpc:   systemclient.NewSystemService(zrpc.MustNewClient(c.SystemRpc)),
    }
}
```

---

## 4. Starter 组件库 (go-admin-starter)

### 4.1 组件架构

```
go-admin-starter/
├── db/
│   ├── mysql/
│   │   ├── starter.go              # MySQL Starter
│   │   ├── config.go               # 配置结构
│   │   └── options.go              # 可选配置
│   ├── postgres/
│   │   ├── starter.go              # PostgreSQL Starter
│   │   ├── gis.go                  # PostGIS 扩展
│   │   └── vector.go               # pgvector 扩展
│   ├── mongo/
│   │   └── starter.go
│   └── base/
│       └── interface.go            # 通用接口定义
├── cache/
│   ├── redis/
│   │   ├── starter.go
│   │   ├── distributed_lock.go
│   │   └── rate_limiter.go
│   └── local/
│       └── starter.go              # 本地缓存
├── mq/
│   ├── kafka/
│   │   └── starter.go
│   ├── rabbitmq/
│   │   └── starter.go
│   └── base/
│       └── interface.go
├── search/
│   ├── elasticsearch/
│   │   └── starter.go
│   └── milvus/
│       └── starter.go
├── scheduler/
│   ├── gocron/
│   │   └── starter.go
│   ├── asynq/
│   │   └── starter.go
│   └── base/
│       └── interface.go
├── log/
│   ├── zap/
│   │   └── starter.go
│   └── logrus/
│       └── starter.go
├── config/
│   ├── viper/
│   │   └── starter.go
│   └── nacos/
│       └── starter.go
├── boot/
│   └── bootstrap.go                # 自动装配
└── go.mod
```

### 4.2 Starter 接口设计

```go
// boot/starter.go
package boot

import "context"

// Starter 组件启动器接口
type Starter interface {
    // Name 返回组件名称
    Name() string

    // Init 初始化组件
    Init(ctx context.Context) error

    // Start 启动组件
    Start(ctx context.Context) error

    // Stop 停止组件
    Stop(ctx context.Context) error

    // Order 启动顺序，数字越小越先启动
    Order() int

    // DependsOn 依赖的组件
    DependsOn() []string
}

// StarterRegistry 组件注册中心
type StarterRegistry struct {
    starters map[string]Starter
    order    []string
}

// Register 注册组件
func (r *StarterRegistry) Register(starter Starter) {
    r.starters[starter.Name()] = starter
}

// StartAll 按序启动所有组件
func (r *StarterRegistry) StartAll(ctx context.Context) error {
    // 拓扑排序，处理依赖关系
    sorted := r.topologicalSort()

    for _, name := range sorted {
        starter := r.starters[name]
        if err := starter.Init(ctx); err != nil {
            return fmt.Errorf("init %s failed: %w", name, err)
        }
        if err := starter.Start(ctx); err != nil {
            return fmt.Errorf("start %s failed: %w", name, err)
        }
        log.Printf("[Starter] %s started", name)
    }
    return nil
}

// StopAll 逆序停止所有组件
func (r *StarterRegistry) StopAll(ctx context.Context) error {
    sorted := r.topologicalSort()

    // 逆序停止
    for i := len(sorted) - 1; i >= 0; i-- {
        starter := r.starters[sorted[i]]
        if err := starter.Stop(ctx); err != nil {
            log.Printf("[Starter] stop %s failed: %v", sorted[i], err)
        }
    }
    return nil
}
```

### 4.3 MySQL Starter 实现

```go
// db/mysql/starter.go
package mysql

import (
    "context"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type MySQLStarter struct {
    config *Config
    db     *gorm.DB
}

type Config struct {
    Enabled      bool   `yaml:"enabled"`
    Host         string `yaml:"host"`
    Port         int    `yaml:"port"`
    Username     string `yaml:"username"`
    Password     string `yaml:"password"`
    Database     string `yaml:"database"`
    MaxIdleConns int    `yaml:"max_idle_conns"`
    MaxOpenConns int    `yaml:"max_open_conns"`
    LogLevel     string `yaml:"log_level"`
}

func NewMySQLStarter(config *Config) *MySQLStarter {
    return &MySQLStarter{config: config}
}

func (s *MySQLStarter) Name() string {
    return "mysql"
}

func (s *MySQLStarter) Order() int {
    return 10
}

func (s *MySQLStarter) DependsOn() []string {
    return []string{"config"}
}

func (s *MySQLStarter) Init(ctx context.Context) error {
    if !s.config.Enabled {
        return nil
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        s.config.Username,
        s.config.Password,
        s.config.Host,
        s.config.Port,
        s.config.Database,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: s.getLogger(),
    })
    if err != nil {
        return err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    sqlDB.SetMaxIdleConns(s.config.MaxIdleConns)
    sqlDB.SetMaxOpenConns(s.config.MaxOpenConns)

    s.db = db
    return nil
}

func (s *MySQLStarter) Start(ctx context.Context) error {
    // 可以在这里执行自动迁移等操作
    return nil
}

func (s *MySQLStarter) Stop(ctx context.Context) error {
    if s.db != nil {
        sqlDB, _ := s.db.DB()
        return sqlDB.Close()
    }
    return nil
}

func (s *MySQLStarter) DB() *gorm.DB {
    return s.db
}

// 注册到全局
func init() {
    boot.RegisterStarterFactory("mysql", func(cfg interface{}) boot.Starter {
        return NewMySQLStarter(cfg.(*Config))
    })
}
```

### 4.4 PostgreSQL + 扩展 Starter

```go
// db/postgres/starter.go
package postgres

import (
    "context"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type PostgresStarter struct {
    config *Config
    db     *gorm.DB
}

type Config struct {
    Enabled    bool     `yaml:"enabled"`
    Host       string   `yaml:"host"`
    Port       int      `yaml:"port"`
    Username   string   `yaml:"username"`
    Password   string   `yaml:"password"`
    Database   string   `yaml:"database"`
    SSLMode    string   `yaml:"ssl_mode"`
    Extensions []string `yaml:"extensions"` // postgis, vector, etc.
}

func (s *PostgresStarter) Init(ctx context.Context) error {
    if !s.config.Enabled {
        return nil
    }

    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        s.config.Host,
        s.config.Port,
        s.config.Username,
        s.config.Password,
        s.config.Database,
        s.config.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    // 启用扩展
    for _, ext := range s.config.Extensions {
        if err := db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)).Error; err != nil {
            log.Printf("Warning: failed to create extension %s: %v", ext, err)
        }
    }

    s.db = db
    return nil
}

// GIS 操作封装
type GISOperations struct {
    db *gorm.DB
}

func (s *PostgresStarter) GIS() *GISOperations {
    return &GISOperations{db: s.db}
}

func (g *GISOperations) FindNearby(tableName string, lat, lng, radiusMeters float64, limit int) *gorm.DB {
    return g.db.Raw(`
        SELECT *, ST_Distance(geom, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) as distance
        FROM ?
        WHERE ST_DWithin(geom, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
        ORDER BY distance
        LIMIT ?
    `, lng, lat, gorm.Expr(tableName), lng, lat, radiusMeters, limit)
}

// Vector 操作封装
type VectorOperations struct {
    db *gorm.DB
}

func (s *PostgresStarter) Vector() *VectorOperations {
    return &VectorOperations{db: s.db}
}

func (v *VectorOperations) SimilaritySearch(tableName, vectorColumn string, embedding []float32, topK int) *gorm.DB {
    vectorStr := fmt.Sprintf("[%s]", floatsToString(embedding))
    return v.db.Raw(`
        SELECT *, 1 - (? <=> ?::vector) as similarity
        FROM ?
        ORDER BY ? <=> ?::vector
        LIMIT ?
    `, gorm.Expr(vectorColumn), vectorStr, gorm.Expr(tableName), gorm.Expr(vectorColumn), vectorStr, topK)
}
```

### 4.5 动态配置加载

```go
// config/nacos/starter.go
package nacos

import (
    "github.com/nacos-group/nacos-sdk-go/v2/clients"
    "github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfigStarter struct {
    config       *Config
    configClient config_client.IConfigClient
    listeners    map[string][]func(data string)
}

type Config struct {
    Enabled     bool   `yaml:"enabled"`
    ServerAddr  string `yaml:"server_addr"`
    Port        uint64 `yaml:"port"`
    NamespaceId string `yaml:"namespace_id"`
    DataId      string `yaml:"data_id"`
    Group       string `yaml:"group"`
}

func (s *NacosConfigStarter) Init(ctx context.Context) error {
    if !s.config.Enabled {
        return nil
    }

    sc := []constant.ServerConfig{
        {
            IpAddr: s.config.ServerAddr,
            Port:   s.config.Port,
        },
    }

    cc := constant.ClientConfig{
        NamespaceId: s.config.NamespaceId,
        TimeoutMs:   5000,
    }

    client, err := clients.NewConfigClient(
        vo.NacosClientParam{
            ClientConfig:  &cc,
            ServerConfigs: sc,
        },
    )
    if err != nil {
        return err
    }

    s.configClient = client
    return nil
}

// GetConfig 获取配置
func (s *NacosConfigStarter) GetConfig(dataId, group string) (string, error) {
    return s.configClient.GetConfig(vo.ConfigParam{
        DataId: dataId,
        Group:  group,
    })
}

// ListenConfig 监听配置变更
func (s *NacosConfigStarter) ListenConfig(dataId, group string, callback func(data string)) error {
    return s.configClient.ListenConfig(vo.ConfigParam{
        DataId: dataId,
        Group:  group,
        OnChange: func(namespace, group, dataId, data string) {
            callback(data)
        },
    })
}

// 动态重载组件
func (s *NacosConfigStarter) ReloadStarter(starterName string, newConfig interface{}) error {
    starter := boot.GetRegistry().Get(starterName)
    if starter == nil {
        return fmt.Errorf("starter %s not found", starterName)
    }

    // 停止旧组件
    ctx := context.Background()
    if err := starter.Stop(ctx); err != nil {
        return err
    }

    // 使用新配置重新初始化
    // 这里需要根据具体组件实现
    return starter.Init(ctx)
}
```

### 4.6 Bootstrap 自动装配

```go
// boot/bootstrap.go
package boot

import (
    "context"
    "os"
    "os/signal"
    "syscall"
)

type Application struct {
    name     string
    config   *Config
    registry *StarterRegistry
}

type Config struct {
    Name string `yaml:"name"`

    // 数据库配置
    MySQL    *mysql.Config    `yaml:"mysql"`
    Postgres *postgres.Config `yaml:"postgres"`
    MongoDB  *mongo.Config    `yaml:"mongodb"`

    // 缓存配置
    Redis *redis.Config `yaml:"redis"`

    // 消息队列配置
    Kafka *kafka.Config `yaml:"kafka"`

    // 搜索引擎配置
    Elasticsearch *es.Config     `yaml:"elasticsearch"`
    Milvus        *milvus.Config `yaml:"milvus"`

    // 定时任务配置
    Scheduler *scheduler.Config `yaml:"scheduler"`

    // 配置中心
    Nacos *nacos.Config `yaml:"nacos"`
}

func NewApplication(configPath string) (*Application, error) {
    // 加载配置
    config, err := loadConfig(configPath)
    if err != nil {
        return nil, err
    }

    app := &Application{
        name:     config.Name,
        config:   config,
        registry: NewStarterRegistry(),
    }

    // 自动注册已启用的组件
    app.autoRegister()

    return app, nil
}

func (a *Application) autoRegister() {
    // 根据配置自动注册组件
    if a.config.MySQL != nil && a.config.MySQL.Enabled {
        a.registry.Register(mysql.NewMySQLStarter(a.config.MySQL))
    }
    if a.config.Postgres != nil && a.config.Postgres.Enabled {
        a.registry.Register(postgres.NewPostgresStarter(a.config.Postgres))
    }
    if a.config.MongoDB != nil && a.config.MongoDB.Enabled {
        a.registry.Register(mongo.NewMongoStarter(a.config.MongoDB))
    }
    if a.config.Redis != nil && a.config.Redis.Enabled {
        a.registry.Register(redis.NewRedisStarter(a.config.Redis))
    }
    if a.config.Kafka != nil && a.config.Kafka.Enabled {
        a.registry.Register(kafka.NewKafkaStarter(a.config.Kafka))
    }
    if a.config.Elasticsearch != nil && a.config.Elasticsearch.Enabled {
        a.registry.Register(es.NewESStarter(a.config.Elasticsearch))
    }
    if a.config.Milvus != nil && a.config.Milvus.Enabled {
        a.registry.Register(milvus.NewMilvusStarter(a.config.Milvus))
    }
}

func (a *Application) Run() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 启动所有组件
    if err := a.registry.StartAll(ctx); err != nil {
        return err
    }

    log.Printf("[Application] %s started successfully", a.name)

    // 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("[Application] shutting down...")
    return a.registry.StopAll(ctx)
}

// 获取已注册的组件实例
func (a *Application) GetStarter(name string) Starter {
    return a.registry.Get(name)
}

func (a *Application) MySQL() *gorm.DB {
    if starter := a.registry.Get("mysql"); starter != nil {
        return starter.(*mysql.MySQLStarter).DB()
    }
    return nil
}

func (a *Application) Postgres() *gorm.DB {
    if starter := a.registry.Get("postgres"); starter != nil {
        return starter.(*postgres.PostgresStarter).DB()
    }
    return nil
}

func (a *Application) Redis() *redis.Client {
    if starter := a.registry.Get("redis"); starter != nil {
        return starter.(*redis.RedisStarter).Client()
    }
    return nil
}

// ... 其他组件的便捷方法
```

### 4.7 使用示例

```go
// main.go
package main

import (
    "github.com/yourname/go-admin-starter/boot"
    "github.com/yourname/go-admin-mono/internal/handler"
)

func main() {
    // 初始化应用
    app, err := boot.NewApplication("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 获取数据库连接
    db := app.MySQL()
    redisClient := app.Redis()

    // 初始化 Handler
    userHandler := handler.NewUserHandler(db, redisClient)

    // 设置路由
    router := gin.Default()
    router.GET("/users", userHandler.List)
    router.POST("/users", userHandler.Create)

    // 启动服务
    go router.Run(":8080")

    // 运行应用 (阻塞直到收到退出信号)
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 5. AI 语音交互项目 (go-xiaozhi)

### 5.1 项目结构

```
go-xiaozhi/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── websocket.go            # WebSocket 连接处理
│   │   ├── audio.go                # 音频处理
│   │   └── chat.go                 # 对话处理
│   ├── service/
│   │   ├── asr/                    # 语音识别
│   │   │   ├── interface.go
│   │   │   ├── whisper.go          # OpenAI Whisper
│   │   │   ├── funasr.go           # 阿里 FunASR
│   │   │   └── sherpa.go           # Sherpa-ONNX
│   │   ├── tts/                    # 语音合成
│   │   │   ├── interface.go
│   │   │   ├── edge_tts.go         # Edge TTS
│   │   │   ├── piper.go            # Piper TTS
│   │   │   └── cosyvoice.go        # CosyVoice
│   │   ├── llm/                    # 大语言模型
│   │   │   ├── interface.go
│   │   │   ├── openai.go           # OpenAI API
│   │   │   ├── claude.go           # Claude API
│   │   │   ├── ollama.go           # Ollama 本地
│   │   │   └── qwen.go             # 通义千问
│   │   ├── rag/                    # RAG 检索增强
│   │   │   ├── retriever.go
│   │   │   ├── embedder.go
│   │   │   └── reranker.go
│   │   ├── memory/                 # 对话记忆
│   │   │   ├── short_term.go       # 短期记忆
│   │   │   └── long_term.go        # 长期记忆
│   │   └── intent/                 # 意图识别
│   │       └── classifier.go
│   ├── protocol/
│   │   ├── xiaozhi.go              # 小智协议
│   │   └── message.go              # 消息定义
│   ├── agent/                      # Agent 能力
│   │   ├── tool.go                 # 工具定义
│   │   ├── executor.go             # 工具执行
│   │   └── planner.go              # 任务规划
│   └── middleware/
│       ├── auth.go
│       └── ratelimit.go
├── pkg/
│   ├── audio/                      # 音频处理工具
│   │   ├── opus.go                 # Opus 编解码
│   │   ├── pcm.go                  # PCM 处理
│   │   └── vad.go                  # 语音活动检测
│   ├── stream/                     # 流式处理
│   │   └── sse.go
│   └── utils/
├── web/                            # 前端页面
├── deployments/
└── go.mod
```

### 5.2 核心架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Client (Web/App/IoT)                        │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │ WebSocket / HTTP
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              Gateway Layer                               │
│                    (认证/限流/协议转换/会话管理)                           │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        ▼                         ▼                         ▼
┌───────────────┐         ┌───────────────┐         ┌───────────────┐
│      ASR      │         │      LLM      │         │      TTS      │
│   语音识别     │────────▶│   对话处理     │────────▶│   语音合成     │
│  Whisper/ASR  │         │ OpenAI/Claude │         │ Edge/Piper    │
└───────────────┘         └───────┬───────┘         └───────────────┘
                                  │
                    ┌─────────────┼─────────────┐
                    ▼             ▼             ▼
             ┌───────────┐ ┌───────────┐ ┌───────────┐
             │    RAG    │ │   Agent   │ │  Memory   │
             │ 知识检索   │ │  工具调用  │ │ 对话记忆   │
             └───────────┘ └───────────┘ └───────────┘
                    │             │             │
                    ▼             ▼             ▼
             ┌───────────┐ ┌───────────┐ ┌───────────┐
             │  Milvus   │ │  Tools    │ │   Redis   │
             │  向量库    │ │ API/函数  │ │   缓存    │
             └───────────┘ └───────────┘ └───────────┘
```

### 5.3 核心接口定义

```go
// service/asr/interface.go
package asr

import "context"

// ASRService 语音识别服务接口
type ASRService interface {
    // Recognize 识别音频，返回文本
    Recognize(ctx context.Context, audio []byte, format string) (*RecognizeResult, error)

    // StreamRecognize 流式识别
    StreamRecognize(ctx context.Context, audioStream <-chan []byte) (<-chan *RecognizeResult, error)
}

type RecognizeResult struct {
    Text       string    `json:"text"`
    Confidence float64   `json:"confidence"`
    Language   string    `json:"language"`
    Duration   float64   `json:"duration"`
    Words      []Word    `json:"words,omitempty"`
}

type Word struct {
    Text  string  `json:"text"`
    Start float64 `json:"start"`
    End   float64 `json:"end"`
}
```

```go
// service/tts/interface.go
package tts

import "context"

// TTSService 语音合成服务接口
type TTSService interface {
    // Synthesize 合成语音
    Synthesize(ctx context.Context, text string, opts *SynthesizeOptions) (*SynthesizeResult, error)

    // StreamSynthesize 流式合成
    StreamSynthesize(ctx context.Context, text string, opts *SynthesizeOptions) (<-chan []byte, error)
}

type SynthesizeOptions struct {
    Voice    string  `json:"voice"`    // 音色
    Speed    float64 `json:"speed"`    // 语速
    Pitch    float64 `json:"pitch"`    // 音调
    Volume   float64 `json:"volume"`   // 音量
    Format   string  `json:"format"`   // 输出格式
}

type SynthesizeResult struct {
    Audio    []byte  `json:"audio"`
    Format   string  `json:"format"`
    Duration float64 `json:"duration"`
}
```

```go
// service/llm/interface.go
package llm

import "context"

// LLMService 大语言模型服务接口
type LLMService interface {
    // Chat 对话
    Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*ChatResponse, error)

    // StreamChat 流式对话
    StreamChat(ctx context.Context, messages []Message, opts *ChatOptions) (<-chan *ChatChunk, error)

    // FunctionCall 函数调用
    FunctionCall(ctx context.Context, messages []Message, functions []Function) (*FunctionCallResponse, error)
}

type Message struct {
    Role       string      `json:"role"`    // system, user, assistant, function
    Content    string      `json:"content"`
    Name       string      `json:"name,omitempty"`
    ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
    ToolCallID string      `json:"tool_call_id,omitempty"`
}

type ChatOptions struct {
    Model       string    `json:"model"`
    Temperature float64   `json:"temperature"`
    MaxTokens   int       `json:"max_tokens"`
    TopP        float64   `json:"top_p"`
    Stream      bool      `json:"stream"`
}

type ChatResponse struct {
    ID      string    `json:"id"`
    Content string    `json:"content"`
    Usage   Usage     `json:"usage"`
}

type ChatChunk struct {
    Delta string `json:"delta"`
    Done  bool   `json:"done"`
}

type Function struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Parameters  interface{} `json:"parameters"`
}
```

### 5.4 WebSocket 处理器

```go
// internal/handler/websocket.go
package handler

import (
    "github.com/gorilla/websocket"
    "sync"
)

type WSHandler struct {
    asr       asr.ASRService
    tts       tts.TTSService
    llm       llm.LLMService
    rag       *rag.Retriever
    memory    *memory.Manager
    upgrader  websocket.Upgrader
}

type Session struct {
    ID        string
    Conn      *websocket.Conn
    UserID    string
    State     SessionState
    Messages  []llm.Message
    AudioBuf  *bytes.Buffer
    mu        sync.Mutex
}

type SessionState int

const (
    StateIdle SessionState = iota
    StateListening
    StateProcessing
    StateSpeaking
)

func (h *WSHandler) HandleConnection(c *gin.Context) {
    // 升级为 WebSocket
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // 创建会话
    session := &Session{
        ID:       uuid.New().String(),
        Conn:     conn,
        UserID:   c.GetString("user_id"),
        State:    StateIdle,
        Messages: h.memory.LoadHistory(c.GetString("user_id")),
        AudioBuf: new(bytes.Buffer),
    }

    // 处理消息
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            break
        }

        go h.handleMessage(session, message)
    }
}

func (h *WSHandler) handleMessage(session *Session, message []byte) {
    var msg protocol.Message
    if err := json.Unmarshal(message, &msg); err != nil {
        return
    }

    switch msg.Type {
    case protocol.MsgTypeAudio:
        h.handleAudio(session, msg.Data)
    case protocol.MsgTypeText:
        h.handleText(session, msg.Data)
    case protocol.MsgTypeControl:
        h.handleControl(session, msg.Data)
    }
}

func (h *WSHandler) handleAudio(session *Session, data []byte) {
    session.mu.Lock()
    session.State = StateListening
    session.mu.Unlock()

    // 语音识别
    result, err := h.asr.Recognize(context.Background(), data, "opus")
    if err != nil {
        h.sendError(session, err)
        return
    }

    // 处理识别结果
    h.processText(session, result.Text)
}

func (h *WSHandler) processText(session *Session, text string) {
    session.mu.Lock()
    session.State = StateProcessing
    session.mu.Unlock()

    // 添加用户消息
    session.Messages = append(session.Messages, llm.Message{
        Role:    "user",
        Content: text,
    })

    // RAG 检索 (如果需要)
    context := h.rag.Retrieve(context.Background(), text, 3)

    // 构建系统提示
    systemPrompt := h.buildSystemPrompt(context)
    messages := append([]llm.Message{{Role: "system", Content: systemPrompt}}, session.Messages...)

    // 流式对话
    stream, err := h.llm.StreamChat(context.Background(), messages, &llm.ChatOptions{
        Model:       "gpt-4",
        Temperature: 0.7,
        Stream:      true,
    })
    if err != nil {
        h.sendError(session, err)
        return
    }

    // 流式 TTS
    var fullResponse strings.Builder
    var sentenceBuffer strings.Builder

    for chunk := range stream {
        fullResponse.WriteString(chunk.Delta)
        sentenceBuffer.WriteString(chunk.Delta)

        // 检测句子边界，流式合成
        if h.isSentenceEnd(sentenceBuffer.String()) {
            sentence := sentenceBuffer.String()
            sentenceBuffer.Reset()

            // 异步合成并发送
            go h.synthesizeAndSend(session, sentence)
        }
    }

    // 处理剩余文本
    if sentenceBuffer.Len() > 0 {
        h.synthesizeAndSend(session, sentenceBuffer.String())
    }

    // 保存对话记录
    session.Messages = append(session.Messages, llm.Message{
        Role:    "assistant",
        Content: fullResponse.String(),
    })
    h.memory.Save(session.UserID, session.Messages)
}

func (h *WSHandler) synthesizeAndSend(session *Session, text string) {
    session.mu.Lock()
    session.State = StateSpeaking
    session.mu.Unlock()

    audioStream, err := h.tts.StreamSynthesize(context.Background(), text, &tts.SynthesizeOptions{
        Voice:  "zh-CN-XiaoxiaoNeural",
        Speed:  1.0,
        Format: "opus",
    })
    if err != nil {
        return
    }

    for audio := range audioStream {
        session.Conn.WriteMessage(websocket.BinaryMessage, audio)
    }
}
```

### 5.5 Agent 工具调用

```go
// internal/agent/tool.go
package agent

import "context"

// Tool 工具定义
type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Parameters  map[string]interface{} `json:"parameters"`
    Handler     ToolHandler            `json:"-"`
}

type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// ToolRegistry 工具注册中心
type ToolRegistry struct {
    tools map[string]*Tool
}

func NewToolRegistry() *ToolRegistry {
    r := &ToolRegistry{
        tools: make(map[string]*Tool),
    }

    // 注册内置工具
    r.registerBuiltinTools()

    return r
}

func (r *ToolRegistry) registerBuiltinTools() {
    // 天气查询
    r.Register(&Tool{
        Name:        "get_weather",
        Description: "获取指定城市的天气信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "city": map[string]interface{}{
                    "type":        "string",
                    "description": "城市名称",
                },
            },
            "required": []string{"city"},
        },
        Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            city := params["city"].(string)
            // 调用天气 API
            return getWeather(ctx, city)
        },
    })

    // 日程管理
    r.Register(&Tool{
        Name:        "create_reminder",
        Description: "创建提醒事项",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "content": map[string]interface{}{
                    "type":        "string",
                    "description": "提醒内容",
                },
                "time": map[string]interface{}{
                    "type":        "string",
                    "description": "提醒时间，格式：2006-01-02 15:04:05",
                },
            },
            "required": []string{"content", "time"},
        },
        Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            content := params["content"].(string)
            timeStr := params["time"].(string)
            return createReminder(ctx, content, timeStr)
        },
    })

    // 智能家居控制
    r.Register(&Tool{
        Name:        "control_device",
        Description: "控制智能家居设备",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "device": map[string]interface{}{
                    "type":        "string",
                    "description": "设备名称，如：客厅灯、空调、电视",
                },
                "action": map[string]interface{}{
                    "type":        "string",
                    "enum":        []string{"on", "off", "adjust"},
                    "description": "操作类型",
                },
                "value": map[string]interface{}{
                    "type":        "string",
                    "description": "调节值，如温度、亮度等",
                },
            },
            "required": []string{"device", "action"},
        },
        Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            device := params["device"].(string)
            action := params["action"].(string)
            value, _ := params["value"].(string)
            return controlDevice(ctx, device, action, value)
        },
    })

    // 知识库查询
    r.Register(&Tool{
        Name:        "search_knowledge",
        Description: "在知识库中搜索相关信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{
                    "type":        "string",
                    "description": "搜索关键词或问题",
                },
            },
            "required": []string{"query"},
        },
        Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            query := params["query"].(string)
            // 调用 RAG 检索
            return searchKnowledge(ctx, query)
        },
    })
}

// Executor 工具执行器
type Executor struct {
    registry *ToolRegistry
    llm      llm.LLMService
}

func (e *Executor) Execute(ctx context.Context, messages []llm.Message) (*ExecuteResult, error) {
    // 获取工具定义
    functions := e.registry.GetFunctions()

    // 调用 LLM 判断是否需要工具
    response, err := e.llm.FunctionCall(ctx, messages, functions)
    if err != nil {
        return nil, err
    }

    // 如果没有工具调用，直接返回
    if len(response.ToolCalls) == 0 {
        return &ExecuteResult{
            Type:    ResultTypeText,
            Content: response.Content,
        }, nil
    }

    // 执行工具调用
    var toolResults []llm.Message
    for _, call := range response.ToolCalls {
        tool := e.registry.Get(call.Function.Name)
        if tool == nil {
            continue
        }

        var params map[string]interface{}
        json.Unmarshal([]byte(call.Function.Arguments), &params)

        result, err := tool.Handler(ctx, params)
        if err != nil {
            result = map[string]interface{}{"error": err.Error()}
        }

        resultJSON, _ := json.Marshal(result)
        toolResults = append(toolResults, llm.Message{
            Role:       "tool",
            Content:    string(resultJSON),
            ToolCallID: call.ID,
        })
    }

    // 将工具结果反馈给 LLM
    messages = append(messages, llm.Message{
        Role:      "assistant",
        ToolCalls: response.ToolCalls,
    })
    messages = append(messages, toolResults...)

    // 递归处理（可能还需要更多工具调用）
    return e.Execute(ctx, messages)
}
```

### 5.6 配置示例

```yaml
# config.yaml
server:
  http_port: 8080
  ws_port: 8081

# ASR 配置
asr:
  provider: whisper  # whisper, funasr, sherpa
  whisper:
    model: whisper-1
    api_key: ${OPENAI_API_KEY}
  funasr:
    endpoint: ws://localhost:10095
    model: paraformer-zh

# TTS 配置
tts:
  provider: edge  # edge, piper, cosyvoice
  edge:
    voice: zh-CN-XiaoxiaoNeural
    rate: "+0%"
    pitch: "+0Hz"
  piper:
    model: /models/zh_CN-huayan-medium.onnx

# LLM 配置
llm:
  provider: openai  # openai, claude, ollama, qwen
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4-turbo-preview
    temperature: 0.7
  claude:
    api_key: ${ANTHROPIC_API_KEY}
    model: claude-3-opus-20240229
  ollama:
    endpoint: http://localhost:11434
    model: qwen:14b

# RAG 配置
rag:
  enabled: true
  embedding:
    provider: openai
    model: text-embedding-3-small
  vector_store:
    type: milvus
    host: localhost
    port: 19530
    collection: xiaozhi_knowledge
  reranker:
    enabled: true
    model: bge-reranker-base

# 记忆配置
memory:
  short_term:
    max_turns: 20
  long_term:
    enabled: true
    provider: redis
    ttl: 7d

# Agent 配置
agent:
  enabled: true
  tools:
    - weather
    - reminder
    - smart_home
    - knowledge_search
```

---

## 6. 数据库封装规范

### 6.1 统一接口规范

```go
// 通用 Repository 接口
type Repository[T any, ID comparable] interface {
    Create(ctx context.Context, entity *T) error
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id ID) error
    FindByID(ctx context.Context, id ID) (*T, error)
    FindAll(ctx context.Context) ([]T, error)
    FindByCondition(ctx context.Context, condition map[string]interface{}) ([]T, error)
    Page(ctx context.Context, pageNum, pageSize int, condition map[string]interface{}) (*PageResult[T], error)
    Count(ctx context.Context, condition map[string]interface{}) (int64, error)
}

// 分页结果
type PageResult[T any] struct {
    List     []T   `json:"list"`
    Total    int64 `json:"total"`
    PageNum  int   `json:"pageNum"`
    PageSize int   `json:"pageSize"`
    Pages    int   `json:"pages"`
}
```

### 6.2 各数据库经典操作清单

| 数据库 | 操作类型 | 封装方法 |
|--------|----------|----------|
| MySQL | CRUD | Create, Update, Delete, Find, Page |
| MySQL | 事务 | Transaction, BatchInsert |
| MySQL | 高级查询 | Join, Subquery, GroupBy |
| PostgreSQL | GIS | FindNearby, FindInPolygon, CalculateDistance |
| PostgreSQL | Vector | SimilaritySearch, HybridSearch |
| MongoDB | 文档操作 | Insert, Update, Aggregate, MapReduce |
| Redis | 缓存 | Get, Set, Delete, Expire, TTL |
| Redis | 分布式 | Lock, Unlock, RateLimit, PubSub |
| Elasticsearch | 搜索 | Search, Aggregate, Suggest, Highlight |
| Milvus | 向量 | Insert, Search, Delete, CreateIndex |
| Kafka | 消息 | Produce, Consume, BatchProduce |

---

## 7. 开发路线图

### Phase 1: 基础框架 (基础设施)
- [ ] 搭建项目骨架
- [ ] 实现 MySQL Starter
- [ ] 实现 Redis Starter
- [ ] 实现配置管理
- [ ] 实现日志系统

### Phase 2: 单体 Admin (go-admin-mono)
- [ ] 用户管理模块
- [ ] 角色权限模块
- [ ] 菜单管理模块
- [ ] 部门管理模块
- [ ] 日志审计模块
- [ ] 代码生成器

### Phase 3: Starter 扩展
- [ ] PostgreSQL + GIS Starter
- [ ] PostgreSQL + Vector Starter
- [ ] MongoDB Starter
- [ ] Elasticsearch Starter
- [ ] Milvus Starter
- [ ] Kafka Starter

### Phase 4: 微服务 Admin (go-admin-micro)
- [ ] 服务拆分设计
- [ ] API 网关
- [ ] 用户服务
- [ ] 认证服务
- [ ] 系统服务
- [ ] 服务治理

### Phase 5: AI 交互 (go-xiaozhi)
- [ ] WebSocket 服务
- [ ] ASR 集成
- [ ] TTS 集成
- [ ] LLM 集成
- [ ] RAG 系统
- [ ] Agent 工具

---

## 附录

### A. 推荐开源项目参考

| 项目 | 地址 | 参考价值 |
|------|------|----------|
| gin-vue-admin | https://github.com/flipped-aurora/gin-vue-admin | CRUD 封装、代码生成 |
| go-admin | https://github.com/go-admin-team/go-admin | 多租户、插件化 |
| go-zero | https://github.com/zeromicro/go-zero | 微服务框架 |
| kratos | https://github.com/go-kratos/kratos | 微服务最佳实践 |
| langchaingo | https://github.com/tmc/langchaingo | LLM 编排 |
| xiaozhi-esp32 | https://github.com/78/xiaozhi-esp32 | 语音交互协议 |

### B. 技术选型对比

| 场景 | 推荐方案 | 备选方案 |
|------|----------|----------|
| Web 框架 | Gin | Hertz, Echo |
| ORM | GORM | Ent, sqlx |
| 微服务框架 | go-zero | Kratos, go-micro |
| 配置中心 | Nacos | Apollo, Consul |
| 服务注册 | Etcd | Consul, Nacos |
| 链路追踪 | Jaeger | Zipkin, SkyWalking |
| ASR | Whisper API | FunASR, Sherpa |
| TTS | Edge TTS | Piper, CosyVoice |
| LLM | OpenAI GPT-4 | Claude, Qwen |

---

> 文档版本：v1.0.0
> 作者：AI Assistant
> 更新日期：2026-01-10
