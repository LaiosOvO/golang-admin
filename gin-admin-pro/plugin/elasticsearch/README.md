# Elasticsearch 插件

Elasticsearch 插件提供全文搜索功能，支持文档索引、搜索、聚合等特性。当前实现为简化版本，可以根据实际需要扩展完整的Elasticsearch功能。

## 功能特性

- 文档索引和搜索
- 索引管理
- 健康检查
- 配置管理
- 集群状态监控

## 使用方法

### 1. 初始化服务

```go
import "gin-admin-pro/plugin/elasticsearch"

// 创建Elasticsearch插件
esPlugin := elasticsearch.NewPlugin(nil)

// 初始化插件
err := esPlugin.Init()
if err != nil {
    log.Fatal("Failed to init elasticsearch plugin:", err)
}

// 获取客户端
client := esPlugin.GetClient()
```

### 2. 基础操作

```go
// 检查连接
err := client.Ping()

// 获取健康状态
health, err := client.Health()

// 创建索引
mapping := map[string]interface{}{
    "mappings": map[string]interface{}{
        "properties": map[string]interface{}{
            "title": map[string]interface{}{
                "type": "text",
            },
            "content": map[string]interface{}{
                "type": "text",
            },
        },
    },
}
err := client.CreateIndex("articles", mapping)

// 索引文档
doc := map[string]interface{}{
    "title":   "Elasticsearch入门",
    "content": "Elasticsearch是一个基于Lucene的搜索服务器。",
}
err := client.IndexDocument("articles", "1", doc)

// 搜索文档
query := map[string]interface{}{
    "query": map[string]interface{}{
        "match": map[string]interface{}{
            "title": "Elasticsearch",
        },
    },
}
result, err := client.SearchDocuments("articles", query)

// 获取文档
var article map[string]interface{}
err := client.GetDocument("articles", "1", &article)

// 删除文档
err := client.DeleteDocument("articles", "1")
```

### 3. 创建默认索引

```go
// 创建系统默认索引
err := esPlugin.CreateDefaultIndexes()
```

## 配置说明

### 基础配置

```yaml
elasticsearch:
  enabled: true                    # 是否启用Elasticsearch
  addresses:                       # Elasticsearch地址列表
    - "http://localhost:9200"
  username: ""                     # 用户名
  password: ""                     # 密码
  timeout: 30                      # 连接超时时间（秒）
  maxRetries: 3                    # 最大重试次数
  defaultIndexPrefix: "gin_admin"  # 默认索引前缀
```

### 自定义配置

```go
customConfig := &elasticsearch.Config{
    Addresses: []string{
        "http://es1:9200",
        "http://es2:9200",
        "http://es3:9200",
    },
    Username:           "elastic",
    Password:           "password123",
    Timeout:            60,
    MaxRetries:         5,
    Enabled:            true,
    DefaultIndexPrefix: "myapp",
}

esPlugin := elasticsearch.NewPlugin(customConfig)
```

## 索引设计

### 用户索引 (gin_admin_users)

用于搜索用户信息：

```json
{
  "mappings": {
    "properties": {
      "username": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "nickname": {
        "type": "text"
      },
      "email": {
        "type": "keyword"
      },
      "mobile": {
        "type": "keyword"
      },
      "status": {
        "type": "integer"
      },
      "dept_id": {
        "type": "integer"
      },
      "created_at": {
        "type": "date"
      }
    }
  }
}
```

### 操作日志索引 (gin_admin_oper_logs)

用于搜索操作日志：

```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text"
      },
      "business_type": {
        "type": "integer"
      },
      "method": {
        "type": "keyword"
      },
      "oper_name": {
        "type": "keyword"
      },
      "dept_name": {
        "type": "keyword"
      },
      "oper_url": {
        "type": "keyword"
      },
      "oper_ip": {
        "type": "ip"
      },
      "oper_location": {
        "type": "keyword"
      },
      "status": {
        "type": "integer"
      },
      "error_msg": {
        "type": "text"
      },
      "oper_time": {
        "type": "date"
      },
      "cost_time": {
        "type": "long"
      }
    }
  }
}
```

## 搜索示例

### 1. 全文搜索

```go
// 搜索用户
query := map[string]interface{}{
    "query": map[string]interface{}{
        "multi_match": map[string]interface{}{
            "query":  "张三",
            "fields": []string{"username", "nickname"},
        },
    },
}

result, err := client.SearchDocuments("gin_admin_users", query)
```

### 2. 精确匹配

```go
// 按状态搜索用户
query := map[string]interface{}{
    "query": map[string]interface{}{
        "term": map[string]interface{}{
            "status": 1,
        },
    },
}

result, err := client.SearchDocuments("gin_admin_users", query)
```

### 3. 范围搜索

```go
// 按时间范围搜索操作日志
query := map[string]interface{}{
    "query": map[string]interface{}{
        "range": map[string]interface{}{
            "oper_time": map[string]interface{}{
                "gte": "2023-01-01",
                "lte": "2023-12-31",
            },
        },
    },
}

result, err := client.SearchDocuments("gin_admin_oper_logs", query)
```

### 4. 布尔查询

```go
// 复合条件搜索
query := map[string]interface{}{
    "query": map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "term": map[string]interface{}{
                        "status": 1,
                    },
                },
                {
                    "range": map[string]interface{}{
                        "created_at": map[string]interface{}{
                            "gte": "2023-01-01",
                        },
                    },
                },
            },
            "must_not": []map[string]interface{}{
                {
                    "term": map[string]interface{}{
                        "username": "admin",
                    },
                },
            },
        },
    },
}

result, err := client.SearchDocuments("gin_admin_users", query)
```

## 集成到业务系统

### 1. 用户搜索服务

```go
type UserSearchService struct {
    esClient *elasticsearch.Client
}

func (s *UserSearchService) SearchUsers(keyword string, filters map[string]interface{}) ([]map[string]interface{}, error) {
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{},
            },
        },
    }
    
    // 添加关键词搜索
    if keyword != "" {
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
            query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
            map[string]interface{}{
                "multi_match": map[string]interface{}{
                    "query":  keyword,
                    "fields": []string{"username", "nickname", "email"},
                },
            },
        )
    }
    
    // 添加过滤条件
    for field, value := range filters {
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
            query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
            map[string]interface{}{
                "term": map[string]interface{}{
                    field: value,
                },
            },
        )
    }
    
    result, err := s.esClient.SearchDocuments("gin_admin_users", query)
    if err != nil {
        return nil, err
    }
    
    var users []map[string]interface{}
    for _, hit := range result.Hits.Hits {
        users = append(users, hit.Source)
    }
    
    return users, nil
}
```

### 2. 日志搜索服务

```go
type LogSearchService struct {
    esClient *elasticsearch.Client
}

func (s *LogSearchService) SearchOperLogs(filters map[string]interface{}, from, size int) (*elasticsearch.SearchResult, error) {
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{},
            },
        },
        "from": from,
        "size": size,
        "sort": []map[string]interface{}{
            {
                "oper_time": map[string]interface{}{
                    "order": "desc",
                },
            },
        },
    }
    
    // 添加过滤条件
    if operName, ok := filters["oper_name"].(string); ok && operName != "" {
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
            query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
            map[string]interface{}{
                "term": map[string]interface{}{
                    "oper_name": operName,
                },
            },
        )
    }
    
    if status, ok := filters["status"].(int); ok {
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
            query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
            map[string]interface{}{
                "term": map[string]interface{}{
                    "status": status,
                },
            },
        )
    }
    
    return s.esClient.SearchDocuments("gin_admin_oper_logs", query)
}
```

## API接口规范

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/infra/elasticsearch/health | 获取ES健康状态 |
| GET | /api/v1/infra/elasticsearch/search/users | 搜索用户 |
| GET | /api/v1/infra/elasticsearch/search/logs | 搜索操作日志 |
| POST | /api/v1/infra/elasticsearch/index/user | 索引用户数据 |
| DELETE | /api/v1/infra/elasticsearch/index/user/{id} | 删除用户索引 |

## 性能优化

### 1. 索引优化

- 合理设置分片数量
- 使用合适的字段类型
- 定期清理旧数据
- 优化索引映射

### 2. 查询优化

- 使用filter而非query
- 避免深度分页
- 合理使用缓存
- 优化查询语句

### 3. 集群优化

- 合理配置节点
- 监控集群状态
- 定期备份
- 优化JVM参数

## 注意事项

1. **连接管理**：确保正确管理Elasticsearch连接
2. **索引设计**：合理设计索引结构和字段类型
3. **数据同步**：确保数据库数据与Elasticsearch数据同步
4. **性能监控**：定期监控搜索性能和集群状态
5. **备份恢复**：定期备份索引数据

## 监控指标

### 集群健康指标

- **集群状态**：green/yellow/red
- **节点数量**：集群中的节点总数
- **分片数量**：主分片和副本分片数量
- **索引数量**：集群中的索引总数

### 性能指标

- **查询QPS**：每秒查询数量
- **索引QPS**：每秒索引数量
- **响应时间**：平均查询响应时间
- **错误率**：查询和索引操作的错误率

## 最佳实践

1. **索引规划**：按业务模块和时间范围规划索引
2. **字段映射**：为搜索字段设置合适的映射
3. **查询优化**：使用高效的查询语句
4. **缓存策略**：合理使用查询缓存
5. **监控告警**：设置关键指标的监控告警

## 扩展功能

当前实现为简化版本，可以根据需要扩展以下功能：

1. **完整客户端**：集成完整的Elasticsearch Go客户端
2. **聚合查询**：支持复杂的聚合分析
3. **索引模板**：支持索引模板管理
4. **别名管理**：支持索引别名功能
5. **滚动索引**：支持时间序列的滚动索引
6. **同步机制**：实现数据库到ES的自动同步

## 测试

```go
func TestElasticsearchClient(t *testing.T) {
    client := NewClient(DefaultConfig())
    
    // 测试连接
    err := client.Ping()
    assert.NoError(t, err)
    
    // 测试健康检查
    health, err := client.Health()
    assert.NoError(t, err)
    assert.Equal(t, "green", health["status"])
    
    // 测试索引创建
    mapping := map[string]interface{}{
        "mappings": map[string]interface{}{
            "properties": map[string]interface{}{
                "title": map[string]interface{}{
                    "type": "text",
                },
            },
        },
    }
    err = client.CreateIndex("test_index", mapping)
    assert.NoError(t, err)
    
    // 测试文档索引
    doc := map[string]interface{}{
        "title": "Test Document",
    }
    err = client.IndexDocument("test_index", "1", doc)
    assert.NoError(t, err)
}
```