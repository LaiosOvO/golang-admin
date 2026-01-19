.PHONY: help build run test clean lint fmt tidy docker-build docker-run

# 默认目标
help:
	@echo "Gin Admin Pro - 企业级后台管理系统"
	@echo ""
	@echo "可用命令:"
	@echo "  build        - 构建应用程序"
	@echo "  run          - 运行应用程序"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  lint         - 代码检查"
	@echo "  fmt          - 格式化代码"
	@echo "  tidy         - 整理依赖"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"
	@echo ""
	@echo "示例:"
	@echo "  make run     # 运行应用程序"

# 构建应用程序
build:
	@echo "构建应用程序..."
	go build -o bin/gin-admin cmd/server/main.go

# 运行应用程序
run:
	@echo "启动应用程序..."
	go run cmd/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f gin-admin

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 整理依赖
tidy:
	@echo "整理依赖..."
	go mod tidy
	go mod verify

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t gin-admin-pro:latest .

# 运行Docker容器
docker-run:
	@echo "运行Docker容器..."
	docker run -p 8080:8080 gin-admin-pro:latest

# 生成API文档
docs:
	@echo "生成API文档..."
	swag init -g cmd/server/main.go

# 安装开发工具
tools:
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 开发模式运行
dev:
	@echo "开发模式运行..."
	go run cmd/server/main.go -config .config.yaml