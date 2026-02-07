# 多阶段构建 Dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/bot ./cmd/bot
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/web ./cmd/web

# 运行阶段
FROM alpine:latest

# 安装运行依赖
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    yt-dlp \
    postgresql-client \
    curl

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/bot /app/bin/bot
COPY --from=builder /app/bin/web /app/bin/web

# 复制 web 静态文件
COPY --from=builder /app/web /app/web

# 创建临时目录
RUN mkdir -p /app/tmp && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 9999

# 默认命令（可以被 docker-compose 覆盖）
CMD ["/app/bin/bot"]
