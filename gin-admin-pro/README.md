# Gin Admin Pro

基于 gin-vue-admin 架构，参考若依Vue Pro (ruoyi-vue-pro) 的实现规范，构建一个完整的企业级后台管理系统。

## 项目特点

- 🚀 **高性能**: 基于 Gin + GORM，性能优异
- 🛡️ **安全可靠**: JWT认证、RBAC权限控制、数据权限
- 🎯 **AI集成**: 支持多种AI模型，提供智能对话能力
- 📊 **多数据库**: MySQL、PostgreSQL、MongoDB、Redis、ES、Milvus
- 🔧 **插件化**: 数据库、中间件、AI服务等插件化设计
- 📝 **文档完善**: 详细的API文档和技术文档

## 技术栈

### 后端框架
- **语言**: Golang 1.21+
- **Web框架**: Gin
- **ORM**: GORM

### 数据库
- **关系型**: MySQL, PostgreSQL (PostGIS, pgvector)
- **非关系型**: MongoDB, Redis
- **搜索引擎**: Elasticsearch
- **向量数据库**: Milvus

### 中间件
- **消息队列**: Kafka
- **缓存**: Redis
- **认证**: JWT

### AI集成
- **语音识别**: 多种ASR服务
- **语音合成**: 多种TTS服务
- **大语言模型**: OpenAI, Ollama等
- **向量搜索**: Milvus

## 项目结构

```
gin-admin-pro/
├── cmd/                    # 程序入口
├── config/                 # 配置文件
├── internal/               # 内部代码
│   ├── api/                # API控制器
│   ├── service/            # 业务逻辑
│   ├── dao/                # 数据访问层
│   ├── model/              # 数据模型
│   ├── middleware/         # 中间件
│   ├── router/             # 路由
│   └── pkg/                # 内部工具包
├── plugin/                 # 插件
│   ├── redis/
│   ├── mysql/
│   ├── postgresql/
│   ├── mongodb/
│   ├── elasticsearch/
│   ├── milvus/
│   ├── kafka/
│   ├── oss/
│   └── ai/
├── docs/                   # 文档
└── scripts/                # 脚本
```

## 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- 其他数据库（可选）

### 安装依赖
```bash
go mod tidy
```

### 配置文件
```bash
cp config/config.yaml .config.yaml
# 编辑 .config.yaml 配置数据库连接等信息
```

### 启动服务
```bash
go run cmd/server/main.go
```

### 构建部署
```bash
make build
```

## 接口文档

启动服务后访问: http://localhost:8080/swagger/index.html

## 开发文档

- [技术调研报告](docs/tech/01_技术调研报告.md)
- [项目架构设计](docs/tech/02_项目架构设计.md)
- [接口规范文档](docs/tech/03_接口规范文档.md)
- [数据库设计文档](docs/tech/04_数据库设计文档.md)

## 开发进度

查看 [PROGRESS.md](PROGRESS.md) 了解项目开发进度。

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 支持

如有问题或建议，请提交 Issue。