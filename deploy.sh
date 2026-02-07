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

# 检查是否在正确的目录
if [ -f "deploy.sh" ] && [ -d "fish-music" ]; then
    echo "⚠️  检测到你在父目录，请不要在项目目录内运行此脚本！"
    echo "    请切换到其他目录，或删除现有的 fish-music 子目录"
    exit 1
fi

echo "📥 下载部署文件到当前目录 ($(pwd))..."

# 下载配置文件
if [ ! -f "config.yaml.example" ]; then
    echo "📥 下载 config.yaml.example..."
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/config.yaml.example \
        -O config.yaml.example
    echo "✅ config.yaml.example"
fi

# 复制配置文件
if [ ! -f "config.yaml" ]; then
    cp config.yaml.example config.yaml
    echo "✅ config.yaml 已创建"
else
    echo "ℹ️  config.yaml 已存在，跳过"
fi

# 下载 docker-compose.yml
if [ ! -f "docker-compose.yml" ]; then
    echo "📥 下载 docker-compose.yml..."
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/docker-compose.yml \
        -O docker-compose.yml
    echo "✅ docker-compose.yml"
else
    echo "ℹ️  docker-compose.yml 已存在，跳过"
fi

# 创建 sql 目录并下载初始化脚本
mkdir -p sql
if [ ! -f "sql/init.sql" ]; then
    echo "📥 下载 sql/init.sql..."
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/sql/init.sql \
        -O sql/init.sql
    echo "✅ sql/init.sql"
else
    echo "ℹ️  sql/init.sql 已存在，跳过"
fi

# 创建临时目录
mkdir -p tmp

# 下载 cookies 模板文件
if [ ! -f "youtube-cookies.txt.example" ]; then
    echo "📥 下载 youtube-cookies.txt.example..."
    wget -q --show-progress \
        https://raw.githubusercontent.com/qqzhoufan/fish_music/main/youtube-cookies.txt.example \
        -O youtube-cookies.txt.example
    echo "✅ youtube-cookies.txt.example (YouTube cookies 模板)"
fi

# 如果不存在 cookies 文件，从模板创建
if [ ! -f "youtube-cookies.txt" ]; then
    cp youtube-cookies.txt.example youtube-cookies.txt
    echo "ℹ️  youtube-cookies.txt 已从模板创建（需要填写 cookie 值）"
else
    echo "ℹ️  youtube-cookies.txt 已存在，跳过"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 部署文件准备完成！"
echo ""
echo "📝 下一步操作："
echo ""
echo "1️⃣  编辑配置文件："
echo "   nano config.yaml"
echo ""
echo "2️⃣  （可选）配置 YouTube cookies："
echo "   nano youtube-cookies.txt"
echo "   解决 YouTube 下载失败问题，见文件内说明"
echo ""
echo "3️⃣  启动服务："
echo "   docker compose up -d"
echo ""
echo "4️⃣  查看日志："
echo "   docker compose logs -f bot"
echo ""
echo "5️⃣  停止服务："
echo "   docker compose down"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 提示："
echo "   - Bot Token: @BotFather (发送 /newbot)"
echo "   - Admin ID: @userinfobot (发送 /start)"
echo "   - YouTube 下载失败? 编辑 youtube-cookies.txt"
echo ""
