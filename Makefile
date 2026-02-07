.PHONY: help build run-bot run-web test clean docker-build docker-up docker-down docker-logs

# 默认目标
help:
	@echo "Fish Music - 可用命令:"
	@echo ""
	@echo "  make build         - 构建 Bot 和 Web 二进制文件"
	@echo "  make run-bot       - 运行 Bot"
	@echo "  make run-web       - 运行 Web 管理端"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建产物"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-up     - 启动 Docker 服务"
	@echo "  make docker-down   - 停止 Docker 服务"
	@echo "  make docker-logs   - 查看 Docker 日志"
	@echo ""
	@echo "  make init-db       - 初始化数据库"
	@echo "  make migrate       - 运行数据库迁移"

# 构建二进制文件
build:
	@echo "构建 Bot..."
	@go build -o bin/bot ./cmd/bot
	@echo "构建 Web..."
	@go build -o bin/web ./cmd/web
	@echo "构建完成!"

# 运行 Bot
run-bot: build
	@echo "启动 Bot..."
	@./bin/bot

# 运行 Web
run-web: build
	@echo "启动 Web..."
	@./bin/web

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 清理
clean:
	@echo "清理构建产物..."
	@rm -rf bin/
	@rm -rf tmp/
	@echo "清理完成!"

# Docker 构建
docker-build:
	@echo "构建 Docker 镜像..."
	@docker-compose build

# Docker 启动
docker-up:
	@echo "启动 Docker 服务..."
	@docker-compose up -d

# Docker 停止
docker-down:
	@echo "停止 Docker 服务..."
	@docker-compose down

# Docker 日志
docker-logs:
	@docker-compose logs -f

# 初始化数据库
init-db:
	@echo "初始化数据库..."
	@psql -h localhost -U fish_music -d fish_music -f sql/init.sql

# 运行数据库迁移
migrate: build
	@echo "运行数据库迁移..."
	@./bin/bot -migrate

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod tidy
	@go mod download

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@golangci-lint run

# 生成 go.sum
go.sum:
	@go mod tidy
