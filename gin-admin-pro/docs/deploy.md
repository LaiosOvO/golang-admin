# Docker 部署指南

## 概述

Gin-Admin 支持使用 Docker Compose 进行一键部署，所有服务（数据库、缓存、消息队列等）都预先配置好。

## 系统要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 8GB RAM
- 至少 20GB 可用磁盘空间

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd gin-admin-pro
```

### 2. 配置文件

复制并修改配置文件：

```bash
# 开发环境
cp config/config.yaml config/config.dev.yaml

# 生产环境  
cp config/config.yaml config/config.prod.yaml
```

根据需要修改数据库密码、API密钥等配置。

### 3. 一键部署

```bash
# 开发环境
./scripts/deploy.sh dev

# 生产环境
./scripts/deploy.sh prod
```

## 服务架构

### 核心数据库

- **MySQL 8.0**: 主数据库，存储用户、角色、权限等核心数据
- **PostgreSQL 16**: 扩展数据库，支持 PostGIS 地理信息和 pgvector 向量搜索
- **MongoDB 6**: 文档数据库，存储日志、会话等非结构化数据

### 缓存和搜索

- **Redis 7**: 缓存和会话存储
- **Elasticsearch 8.11**: 全文搜索引擎

### 向量和AI

- **Milvus 2.3**: 向量数据库，支持 AI 向量搜索
- **MinIO**: 对象存储，用于 Milvus 数据持久化

### 消息队列

- **Kafka 7.4**: 消息队列，支持异步处理
- **Zookeeper**: Kafka 依赖服务

### 应用层

- **Gin-Admin**: 主应用（Go 服务）
- **Nginx**: 反向代理和负载均衡

## 配置说明

### 环境变量

| 服务 | 变量 | 默认值 | 说明 |
|------|------|--------|------|
| MySQL | MYSQL_ROOT_PASSWORD | root123 | 数据库root密码 |
| MySQL | MYSQL_DATABASE | gin_admin | 数据库名称 |
| PostgreSQL | POSTGRES_PASSWORD | gin_admin123 | PostgreSQL密码 |
| Redis | requirepass | redis123 | Redis密码 |
| MinIO | MINIO_ACCESS_KEY | minioadmin | 访问密钥 |
| MinIO | MINIO_SECRET_KEY | minioadmin123 | 秘密密钥 |

### 端口映射

| 服务 | 端口 | 说明 |
|------|------|------|
| app | 8080 | 主应用API |
| nginx | 80/443 | HTTP/HTTPS代理 |
| mysql | 3306 | MySQL数据库 |
| postgresql | 5432 | PostgreSQL数据库 |
| mongodb | 27017 | MongoDB数据库 |
| redis | 6379 | Redis缓存 |
| elasticsearch | 9200 | Elasticsearch HTTP |
| elasticsearch | 9300 | Elasticsearch TCP |
| milvus | 19530 | Milvus gRPC |
| milvus | 9091 | Milvus Web UI |
| kafka | 9092 | Kafka消息队列 |
| zookeeper | 2181 | Zookeeper协调服务 |
| minio | 9000 | MinIO API |
| minio | 9001 | MinIO Web UI |

### 数据持久化

所有服务数据都通过 Docker volumes 持久化：

- `mysql_data`: MySQL 数据文件
- `postgresql_data`: PostgreSQL 数据文件
- `mongodb_data`: MongoDB 数据文件
- `redis_data`: Redis 数据文件
- `elasticsearch_data`: Elasticsearch 索引文件
- `milvus_data`: Milvus 向量数据
- `kafka_data`: Kafka 消息数据
- `minio_data`: MinIO 对象存储

## 管理命令

### 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 启动特定服务
docker-compose up -d mysql redis
```

### 查看服务状态

```bash
docker-compose ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f mysql
```

### 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（注意：会删除所有数据）
docker-compose down -v
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart app
```

## 数据库初始化

### MySQL

启动时会自动执行 `./docs/sql` 目录下的 SQL 文件：

- `init.sql`: 基础表结构和初始数据
- `data.sql`: 测试数据（开发环境）

### PostgreSQL

启动时会自动执行：

- `./docs/sql/postgis.sql`: PostGIS 扩展
- `./docs/sql/pgvector.sql`: pgvector 扩展

### MongoDB

启动时会执行：

- `./docs/sql/mongo-init.js`: 初始化脚本

## 监控和维护

### 健康检查

应用提供健康检查接口：

```bash
curl http://localhost:8080/health
```

### 性能监控

建议使用以下工具进行监控：

- **Prometheus**: 指标收集
- **Grafana**: 可视化面板
- **ELK Stack**: 日志分析

### 备份策略

#### 数据库备份

```bash
# MySQL备份
docker-compose exec mysql mysqldump -u root -proot123 gin_admin > backup.sql

# PostgreSQL备份
docker-compose exec postgresql pg_dump -U gin_admin gin_admin > pg_backup.sql

# MongoDB备份
docker-compose exec mongodb mongodump --db gin_admin --out /tmp/backup
```

#### 恢复数据

```bash
# MySQL恢复
docker-compose exec -i mysql mysql -u root -proot123 gin_admin < backup.sql

# PostgreSQL恢复
docker-compose exec -i postgresql psql -U gin_admin gin_admin < pg_backup.sql

# MongoDB恢复
docker-compose exec mongodb mongorestore --db gin_admin /tmp/backup
```

## 故障排除

### 常见问题

1. **应用无法启动**
   ```bash
   # 检查日志
   docker-compose logs app
   
   # 检查配置文件
   cat config/config.dev.yaml
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose ps
   
   # 检查数据库日志
   docker-compose logs mysql
   docker-compose logs postgresql
   docker-compose logs mongodb
   ```

3. **内存不足**
   ```bash
   # 增加Docker内存限制
   # 或关闭不需要的服务
   docker-compose stop elasticsearch milvus
   ```

4. **端口冲突**
   ```bash
   # 修改 docker-compose.yaml 中的端口映射
   # 或停止占用端口的服务
   netstat -tulpn | grep :8080
   ```

### 日志位置

- 应用日志: `./logs/app.log`
- Nginx日志: Docker volume `nginx_logs`
- 数据库日志: 通过 `docker-compose logs` 查看

## 安全配置

### 生产环境安全建议

1. **修改默认密码**: 更改所有数据库的默认密码
2. **SSL证书**: 使用有效的SSL证书
3. **网络隔离**: 使用Docker网络隔离服务
4. **防火墙**: 只开放必要端口
5. **定期更新**: 保持Docker镜像和依赖更新

### SSL配置

```bash
# 生成SSL证书
mkdir -p docker/nginx/ssl
openssl req -x509 -newkey rsa:4096 -keyout docker/nginx/ssl/key.pem \
        -out docker/nginx/ssl/cert.pem -days 365 -nodes \
        -subj "/C=CN/ST=State/L=City/O=Organization/CN=domain.com"
```

## 扩展和定制

### 添加新服务

1. 在 `docker-compose.yaml` 中添加服务定义
2. 配置网络和数据卷
3. 更新部署脚本
4. 测试并验证

### 环境变量配置

可以通过 `.env` 文件管理环境变量：

```bash
# .env
MYSQL_ROOT_PASSWORD=your_password
REDIS_PASSWORD=your_redis_password
```

然后在 `docker-compose.yaml` 中引用：

```yaml
services:
  redis:
    environment:
      - requirepass=${REDIS_PASSWORD}
```

## 支持

如遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查项目的 Issues 页面
3. 提交新的 Issue 并包含详细的错误信息和环境信息