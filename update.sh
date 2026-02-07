#!/bin/bash

echo "🔄 Fish Music 更新脚本"
echo "===================="
echo ""

# 检查是否在正确的目录
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 错误：请在项目目录中运行此脚本"
    echo "   该目录应包含 docker-compose.yml 文件"
    exit 1
fi

echo "📋 步骤 1：备份当前配置"
if [ -f "config.yaml" ]; then
    cp config.yaml config.yaml.backup.$(date +%Y%m%d_%H%M%S)
    echo "✅ 配置已备份"
fi

echo ""
echo "📥 步骤 2：拉取最新镜像"
docker compose pull
if [ $? -ne 0 ]; then
    echo "❌ 拉取镜像失败"
    exit 1
fi
echo "✅ 镜像更新完成"

echo ""
echo "📝 步骤 3：检查 cookies 文件"
if [ ! -f "youtube-cookies.txt" ]; then
    touch youtube-cookies.txt
    echo "✅ 已创建 youtube-cookies.txt"
    echo "   请通过 /cookies 命令配置 cookie"
else
    echo "✅ youtube-cookies.txt 已存在"
fi

echo ""
echo "🔧 步骤 4：检查 config.yaml 配置"
if ! grep -q "cookies_file" config.yaml; then
    echo "⚠️  config.yaml 中缺少 cookies_file 配置"
    echo ""
    echo "正在添加配置..."
    # 使用 sed 添加配置
    if grep -q "temp_dir:" config.yaml; then
        sed -i '/temp_dir:/a\  cookies_file: "/app/youtube-cookies.txt"  # YouTube cookies（通过 /cookies 命令配置）' config.yaml
        echo "✅ 已添加 cookies_file 配置"
    else
        echo "❌ 无法自动添加配置，请手动编辑 config.yaml"
        echo "   在 download 部分添加：cookies_file: \"/app/youtube-cookies.txt\""
    fi
else
    echo "✅ config.yaml 已配置 cookies_file"
fi

echo ""
echo "🛑 步骤 5：停止旧容器"
docker compose down
echo "✅ 旧容器已停止"

echo ""
echo "🚀 步骤 6：启动新容器"
docker compose up -d
if [ $? -ne 0 ]; then
    echo "❌ 启动失败"
    exit 1
fi
echo "✅ 新容器已启动"

echo ""
echo "⏳ 步骤 7：等待服务启动"
sleep 5

echo ""
echo "✅ 验证服务状态"
docker compose ps

echo ""
echo "📝 验证 cookies 挂载"
if docker compose exec bot ls -la /app/youtube-cookies.txt >/dev/null 2>&1; then
    echo "✅ cookies 文件已正确挂载"
else
    echo "⚠️  cookies 文件挂载失败"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 更新完成！"
echo ""
echo "📝 下一步："
echo "1. 配置 cookies（如果还没有）"
echo "   发送 /cookies 命令给 Bot"
echo ""
echo "2. 查看 Bot 日志"
echo "   docker compose logs -f bot"
echo ""
echo "3. 测试 YouTube 下载"
echo "   发送一个 YouTube 链接试试"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
