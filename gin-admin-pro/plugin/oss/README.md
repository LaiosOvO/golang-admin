# OSS 文件存储插件

## 概述

OSS (Object Storage Service) 插件提供了统一的文件存储接口，支持多种存储后端，包括本地存储、MinIO、阿里云OSS等。该插件采用工厂模式设计，方便扩展和切换存储方式。

## 支持的存储方式

### 1. 本地存储 (Local Storage)
- 将文件存储在服务器本地文件系统
- 支持按日期目录结构存储
- 适合单机部署和开发环境

### 2. MinIO 存储 (MinIO Storage)
- 兼容 S3 API 的对象存储
- 支持分布式存储
- 适合生产环境部署（暂未完全实现）

### 3. 阿里云 OSS (Aliyun OSS)
- 阿里云对象存储服务
- 高可用、高可靠
- 适合企业级应用（暂未实现）

## 接口设计

### OSSInterface 统一接口

```go
type OSSInterface interface {
    // UploadFile 上传文件
    UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error)
    
    // DeleteFile 删除文件
    DeleteFile(key string) error
    
    // GetFileURL 获取文件访问URL
    GetFileURL(key string) string
    
    // IsExists 检查文件是否存在
    IsExists(key string) (bool, error)
}
```

## 配置说明

### 上传配置 (UploadConfig)

```yaml
upload:
  maxSize: 10485760          # 最大文件大小（字节），默认10MB
  allowedTypes:              # 允许的文件类型
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".gif"
    - ".pdf"
    - ".doc"
    - ".docx"
    - ".xls"
    - ".xlsx"
  path: "./uploads"         # 本地存储路径
  urlPrefix: "/uploads"      # URL前缀
```

### MinIO 配置 (环境变量)

```bash
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=gin-admin
MINIO_USE_SSL=false
```

## 使用方法

### 1. 初始化存储实例

```go
// 使用工厂模式创建存储实例
factory := oss.NewStorageFactory()
storage, err := factory.CreateStorage("local")
if err != nil {
    log.Fatal(err)
}

// 或使用默认存储
storage, err := oss.GetDefaultStorage()
```

### 2. 上传文件

```go
file, _ := os.Open("test.jpg")
defer file.Close()

url, err := storage.UploadFile("test.jpg", file, fileInfo.Size(), "image/jpeg")
if err != nil {
    log.Fatal(err)
}

fmt.Println("文件URL:", url)
```

### 3. 删除文件

```go
err := storage.DeleteFile("test.jpg")
if err != nil {
    log.Fatal(err)
}
```

### 4. 检查文件是否存在

```go
exists, err := storage.IsExists("test.jpg")
if err != nil {
    log.Fatal(err)
}
fmt.Println("文件存在:", exists)
```

## 本地存储实现

### 目录结构

```
uploads/
├── 2025/
│   ├── 01/
│   │   ├── 20/
│   │   │   ├── photo_20250120150405.jpg
│   │   │   └── document_20250120150406.pdf
│   │   └── 21/
│   │       └── ...
│   └── 02/
│       └── ...
└── 2025/
    └── ...
```

### 特性

- **自动目录创建**: 根据日期自动创建存储目录
- **文件重命名**: 自动生成带时间戳的文件名，避免冲突
- **大小限制**: 支持配置文件大小限制
- **类型验证**: 支持配置允许的文件类型
- **错误处理**: 完善的错误处理和回滚机制

## API 接口

### 文件上传

**接口地址**: `POST /api/v1/infra/file/upload`

**请求参数**:
- `file` (file): 上传的文件

**响应数据**:
```json
{
    "code": 0,
    "data": {
        "url": "/uploads/2025/01/20/photo_20250120150405.jpg",
        "fileName": "photo.jpg",
        "size": 1024000,
        "type": ".jpg"
    },
    "msg": "操作成功"
}
```

### 批量文件上传

**接口地址**: `POST /api/v1/infra/file/upload-multiple`

**请求参数**:
- `files` (file[]): 上传的文件数组

**响应数据**:
```json
{
    "code": 0,
    "data": [
        {
            "url": "/uploads/2025/01/20/photo1_20250120150405.jpg",
            "fileName": "photo1.jpg",
            "size": 1024000,
            "type": ".jpg"
        },
        {
            "url": "/uploads/2025/01/20/photo2_20250120150406.jpg",
            "fileName": "photo2.jpg",
            "size": 2048000,
            "type": ".jpg"
        }
    ],
    "msg": "操作成功"
}
```

### 删除文件

**接口地址**: `DELETE /api/v1/infra/file/delete`

**请求参数**:
```json
{
    "url": "/uploads/2025/01/20/photo_20250120150405.jpg"
}
```

**响应数据**:
```json
{
    "code": 0,
    "data": null,
    "msg": "操作成功"
}
```

## 安全考虑

### 1. 文件类型验证
- 基于扩展名的白名单验证
- 可配置允许的文件类型
- 防止上传危险文件

### 2. 文件大小限制
- 可配置最大文件大小
- 防止恶意大文件上传
- 保护服务器资源

### 3. 访问控制
- 所有文件操作接口需要认证
- 支持权限中间件集成
- 防止未授权访问

### 4. 路径安全
- 防止路径遍历攻击
- 文件名特殊字符过滤
- 安全的文件名生成

## 性能优化

### 1. 并发处理
- 支持并发文件上传
- 线程安全的存储操作
- 合理的资源管理

### 2. 内存管理
- 流式文件处理
- 避免大文件内存占用
- 及时释放资源

### 3. 目录优化
- 按日期分目录存储
- 避免单目录文件过多
- 提高文件检索效率

## 扩展指南

### 添加新的存储后端

1. 实现 `OSSInterface` 接口
2. 在工厂方法中添加创建逻辑
3. 添加配置结构
4. 编写单元测试

示例：
```go
type CustomStorage struct {
    config FileConfig
}

func (cs *CustomStorage) UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error) {
    // 实现上传逻辑
}

// 实现其他接口方法...

func NewCustomStorage() (*CustomStorage, error) {
    // 初始化逻辑
}
```

### 添加新功能

1. 在接口中定义新方法
2. 在所有实现中添加相应逻辑
3. 更新配置结构
4. 更新文档

## 最佳实践

### 1. 文件命名
- 使用有意义的文件名
- 包含时间戳避免冲突
- 过滤特殊字符

### 2. 目录规划
- 按业务模块分目录
- 按时间分目录
- 考虑数据量增长

### 3. 错误处理
- 记录详细错误日志
- 提供友好的错误信息
- 支持错误重试

### 4. 监控告警
- 监控存储空间使用
- 监控上传下载速度
- 设置异常告警

## 故障排除

### 常见问题

1. **文件上传失败**
   - 检查存储目录权限
   - 检查文件大小限制
   - 检查文件类型限制

2. **文件无法访问**
   - 检查URL路径配置
   - 检查文件是否存在
   - 检查静态文件服务配置

3. **存储空间不足**
   - 清理过期文件
   - 扩展存储容量
   - 启用文件压缩

### 调试方法

1. 启用详细日志
2. 检查网络连接
3. 验证配置参数
4. 测试存储权限

## 总结

OSS 插件提供了完整的文件存储解决方案，具有以下特点：

1. **统一接口**: 支持多种存储后端的统一访问接口
2. **易于扩展**: 工厂模式设计，方便添加新的存储方式
3. **安全可靠**: 完善的安全验证和错误处理
4. **性能优秀**: 优化的存储结构和并发处理
5. **配置灵活**: 支持多种配置方式和参数调整

该插件为企业级应用提供了可靠的文件存储基础，可以根据实际需求选择合适的存储方案。