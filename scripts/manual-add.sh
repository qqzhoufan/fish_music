#!/bin/bash

# Fish Music - 手动添加音乐工具
# 使用方法: ./scripts/manual-add.sh

echo "=== Fish Music 添加音乐工具 ==="
echo ""
echo "请准备好以下信息："
echo "  1. 歌曲标题"
echo "  2. 歌手名称"
echo "  3. Telegram File ID (从音频文件获取)"
echo ""
echo "获取 File ID 的方法："
echo "  方法1: 向 Bot 发送音频，然后查看数据库"
echo "  方法2: 使用 Telegram Bot API"
echo "  方法3: 从已有的 Telegram 音频链接获取"
echo ""
echo "----------------------------------------"
echo ""

# 交互式输入
read -p "歌曲标题: " TITLE
read -p "歌手名称: " ARTIST
read -p "专辑名称 (可选，直接回车跳过): " ALBUM
read -p "Telegram File ID: " FILE_ID
read -p "时长(秒) (可选，直接回车跳过): " DURATION
read -p "年份 (可选，直接回车跳过): " YEAR

if [ -z "$TITLE" ] || [ -z "$ARTIST" ] || [ -z "$FILE_ID" ]; then
    echo "❌ 错误: 标题、歌手和 File ID 不能为空"
    exit 1
fi

# 生成唯一哈希
UNIQUE_HASH="${TITLE}_${ARTIST}_$(date +%s)"

# 构建插入语句
if [ -z "$ALBUM" ]; then
    ALBUM="NULL"
else
    ALBUM="'$ALBUM'"
fi

if [ -z "$DURATION" ] || [ "$DURATION" = "0" ]; then
    DURATION="0"
fi

if [ -z "$YEAR" ] || [ "$YEAR" = "0" ]; then
    YEAR="NULL"
else
    YEAR="'$YEAR'"
fi

# 执行 SQL
SQL="INSERT INTO songs (unique_hash, file_id, source_url, title, artist, album, duration, file_size, country_code, year, status)
VALUES ('$UNIQUE_HASH', '$FILE_ID', '', '$TITLE', '$ARTIST', $ALBUM, $DURATION, 0, '', $YEAR, 'active');"

echo ""
echo "执行的 SQL:"
echo "$SQL"
echo ""

# 询问是否执行
read -p "确认添加？(y/n): " CONFIRM
if [ "$CONFIRM" = "y" ] || [ "$CONFIRM" = "Y" ]; then
    docker exec -i fish_music_db psql -U fish_music -d fish_music <<< "$SQL"
    echo ""
    echo "✅ 歌曲添加成功!"
    echo "   $ARTIST - $TITLE"
else
    echo "❌ 已取消"
fi
