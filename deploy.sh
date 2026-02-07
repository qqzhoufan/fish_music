#!/bin/bash

# Fish Music 一键部署脚本
# 使用 Docker Hub 镜像快速部署

set -e

echo "🎵 Fish Music - 一键部署脚本"
echo "================================"

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 创建项目目录
mkdir -p fish-music
cd fish-music

echo "📥 下载部署文件..."

# 下载配置文件
if [ ! -f "config.yaml.example" ]; then
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/config.yaml.example \
        -O config.yaml.example
    echo "✅ config.yaml.example"
fi

# 复制配置文件
if [ ! -f "config.yaml" ]; then
    cp config.yaml.example config.yaml
    echo "✅ config.yaml 已创建"
fi

# 下载 docker-compose.yml
if [ ! -f "docker-compose.yml" ]; then
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/docker-compose.yml \
        -O docker-compose.yml
    echo "✅ docker-compose.yml"
fi

# 创建 sql 目录并下载初始化脚本
mkdir -p sql
if [ ! -f "sql/init.sql" ]; then
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/sql/init.sql \
        -O sql/init.sql
    echo "✅ sql/init.sql"
fi

# 创建临时目录
mkdir -p tmp

echo ""
echo "📝 请编辑配置文件，填入你的 Bot Token 和 Admin ID："
echo "   nano config.yaml"
echo ""
echo "配置完成后，运行以下命令启动服务："
echo "   docker compose up -d"
echo ""
echo "查看日志："
echo "   docker compose logs -f bot"
echo ""
echo "停止服务："
echo "   docker compose down"
echo ""
echo "✅ 部署文件准备完成！"
