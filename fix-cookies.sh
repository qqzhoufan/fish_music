#!/bin/bash

echo "🔧 修复 YouTube Cookies 配置"
echo "================================"
echo ""

# 检查是否存在 cookies 文件
if [ ! -f "youtube-cookies.txt" ]; then
    echo "❌ 未找到 youtube-cookies.txt"
    echo ""
    echo "正在创建空文件..."
    touch youtube-cookies.txt
    echo "✅ 已创建 youtube-cookies.txt"
    echo ""
    echo "⚠️  现在需要通过 Bot 配置 cookies："
    echo "1. 发送 /cookies 给你的 Bot"
    echo "2. 按提示获取并发送 cookie 值"
    echo "3. 重启服务: docker compose restart bot"
else
    echo "✅ youtube-cookies.txt 已存在"
fi

echo ""
echo "📋 当前 cookies 文件内容："
echo "--------------------------------"
cat youtube-cookies.txt
echo "--------------------------------"
echo ""

echo "🔍 检查配置文件..."
if grep -q "cookies_file" config.yaml; then
    echo "✅ config.yaml 中已配置 cookies_file"
    grep "cookies_file" config.yaml
else
    echo "⚠️  config.yaml 中未配置 cookies_file"
    echo ""
    echo "请添加以下配置到 config.yaml 的 download 部分："
    echo ""
    echo "download:"
    echo "  worker_count: 3"
    echo "  max_file_size: 50"
    echo "  temp_dir: \"./tmp\""
    echo "  cookies_file: \"/app/youtube-cookies.txt\""
fi

echo ""
echo "================================"
echo "✅ 修复完成！"
echo ""
echo "下一步："
echo "1. 通过 Bot 发送 /cookies 命令配置 cookie"
echo "2. 重启服务: docker compose restart bot"
echo "3. 测试 YouTube 下载"
