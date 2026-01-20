#!/bin/bash

# Gin-Admin 部署脚本
# 使用方法: ./deploy.sh [dev|prod]

set -e

ENVIRONMENT=${1:-dev}
PROJECT_NAME="gin-admin"

echo "========================================"
echo "Gin-Admin 部署脚本"
echo "环境: $ENVIRONMENT"
echo "========================================"

# 检查Docker和Docker Compose
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    exit 1
fi

# 清理旧容器和网络
echo "清理旧容器..."
docker-compose down --remove-orphans || true
docker system prune -f

# 创建必要的目录
echo "创建目录结构..."
mkdir -p logs
mkdir -p uploads
mkdir -p docker/nginx/ssl

# 生成SSL证书（开发环境）
if [ "$ENVIRONMENT" = "dev" ]; then
    echo "生成开发用SSL证书..."
    if [ ! -f "docker/nginx/ssl/cert.pem" ]; then
        openssl req -x509 -newkey rsa:4096 -keyout docker/nginx/ssl/key.pem -out docker/nginx/ssl/cert.pem -days 365 -nodes \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=Development/OU=IT/CN=localhost"
    fi
fi

# 复制配置文件
echo "复制配置文件..."
cp config/config.yaml config/config.backup.yaml 2>/dev/null || true

if [ "$ENVIRONMENT" = "prod" ]; then
    if [ ! -f "config/config.prod.yaml" ]; then
        echo "错误: 生产环境配置文件 config/config.prod.yaml 不存在"
        exit 1
    fi
    CONFIG_FILE="config/config.prod.yaml"
else
    if [ ! -f "config/config.dev.yaml" ]; then
        echo "错误: 开发环境配置文件 config/config.dev.yaml 不存在"
        exit 1
    fi
    CONFIG_FILE="config/config.dev.yaml"
fi

# 构建应用镜像
echo "构建应用镜像..."
docker build -t $PROJECT_NAME:latest .

# 启动数据库服务（先启动依赖服务）
echo "启动数据库服务..."
docker-compose up -d mysql postgresql mongodb redis elasticsearch

# 等待数据库启动
echo "等待数据库启动..."
sleep 30

# 检查数据库连接
echo "检查数据库连接..."
docker-compose exec mysql mysqladmin ping -h localhost -u root -proot123 || {
    echo "MySQL 连接失败"
    docker-compose logs mysql
    exit 1
}

docker-compose exec postgresql pg_isready -U gin_admin || {
    echo "PostgreSQL 连接失败"
    docker-compose logs postgresql
    exit 1
}

# 启动剩余服务
echo "启动剩余服务..."
docker-compose up -d zookeeper kafka etcd minio milvus

# 等待Milvus启动
echo "等待Milvus启动..."
sleep 60

# 启动应用
echo "启动应用..."
docker-compose up -d app

# 等待应用启动
echo "等待应用启动..."
sleep 30

# 检查应用健康状态
echo "检查应用健康状态..."
for i in {1..10}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ 应用启动成功！"
        break
    fi
    
    if [ $i -eq 10 ]; then
        echo "❌ 应用启动失败"
        docker-compose logs app
        exit 1
    fi
    
    echo "等待应用启动... ($i/10)"
    sleep 10
done

# 启动Nginx（如果配置了）
if [ -f "docker/nginx/nginx.conf" ]; then
    echo "启动Nginx..."
    docker-compose up -d nginx
fi

echo "========================================"
echo "🎉 部署完成！"
echo "========================================"
echo "应用地址: http://localhost:8080"
echo "健康检查: http://localhost:8080/health"
if [ -f "docker/nginx/nginx.conf" ]; then
    echo "Nginx代理: http://localhost"
fi
echo ""
echo "数据库连接信息:"
echo "MySQL: localhost:3306 (root/root123)"
echo "PostgreSQL: localhost:5432 (gin_admin/gin_admin123)"
echo "MongoDB: localhost:27017 (admin/admin123)"
echo "Redis: localhost:6379 (密码: redis123)"
echo "Elasticsearch: localhost:9200"
echo "Milvus: localhost:19530"
echo "Kafka: localhost:9092"
echo ""
echo "管理命令:"
echo "查看日志: docker-compose logs -f [service]"
echo "停止服务: docker-compose down"
echo "重启服务: docker-compose restart [service]"
echo "========================================"