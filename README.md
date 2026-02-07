# 🎵 Fish Music

> 基于 Telegram 的个人云端音乐机器人 - 利用 Telegram 无限云存储构建你的专属音乐库

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)]
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)]
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)]
[![Deploy](https://img.shields.io/badge/Deployment-Easy-success)]()
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)]

---

## ⚡ 5 分钟快速部署

> 只需要一台服务器（1GB 内存即可），Docker 一键部署！

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

### 部署方式

#### 方式一：使用 Docker Hub 镜像（推荐）⭐

**无需克隆源码，直接使用已构建的镜像！**

```bash
# 下载并运行一键部署脚本
curl -fsSL https://raw.githubusercontent.com/qqzhoufan/fish_music/main/deploy.sh -o deploy.sh
chmod +x deploy.sh
./deploy.sh
```

**或者手动下载文件：**

```bash
# 1. 创建项目目录
mkdir fish-music && cd fish-music

# 2. 下载配置文件
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/config.yaml.example
mv config.yaml.example config.yaml

# 3. 编辑配置（填入 Bot Token 和 Admin ID）
nano config.yaml

# 4. 下载 docker-compose.yml
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/docker-compose.yml

# 5. 下载 sql 目录和初始化脚本
mkdir -p sql
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/sql/init.sql -O sql/init.sql

# 6. 创建临时目录
mkdir -p tmp

# 7. 启动服务
docker compose up -d

# 8. 查看日志
docker compose logs -f bot
```

**镜像地址**：`zhouwl/fish-music:latest`

---

#### 方式二：从源码部署

```bash
# 1️⃣ 获取 Telegram Bot Token

在 Telegram 中搜索 [@BotFather](https://t.me/BotFather)，发送 `/newbot` 创建机器人，获取 Token。

# 2️⃣ 获取你的 Telegram ID

在 Telegram 中搜索 [@userinfobot](https://t.me/userinfobot)，发送 `/start` 获取你的 User ID。

# 3️⃣ 克隆并配置

```bash
# 克隆项目
git clone https://github.com/qqzhoufan/fish_music.git
cd fish-music

# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置（填入 Bot Token 和 Admin ID）
nano config.yaml
```

**config.yaml 配置示例：**

```yaml
bot:
  token: "1234567890:ABCdefGhIJKlmNoPQRsTUVwxyZ"  # 你的 Bot Token
  admin_id: 123456789                             # 你的 Telegram ID

database:
  host: "postgres"      # Docker 部署保持默认
  port: 5432
  user: "fish_music"
  password: "fish_music_pass"
  dbname: "fish_music"

web:
  port: 9999
  username: "admin"       # 建议修改
  password: "fishmusic2024"  # 建议修改为强密码
```

#### 4️⃣ 启动服务

```bash
# 启动所有服务
docker compose up -d

# 查看运行状态
docker compose ps

# 查看日志
docker compose logs -f bot
```

#### 5️⃣ 验证部署

- **测试 Bot**: 在 Telegram 中找到你的 Bot，发送 `/start`
- **测试 Web**: 访问 `http://你的服务器IP:9999`

✅ 部署完成！开始添加音乐吧！

### 📖 详细部署文档

遇到问题？查看 [详细部署指南](./DEPLOY.md)

---

## 🐳 Docker 镜像说明

本项目提供两种部署方式：

### 使用 Docker Hub 镜像（推荐）

- **镜像地址**：`zhouwl/fish-music:latest`
- **镜像大小**：约 200MB
- **优点**：
  - ✅ 无需本地构建，下载即用
  - ✅ 部署速度快（约 1-2 分钟）
  - ✅ 自动更新，跟随最新版本
- **使用场景**：快速部署、生产环境

### 从源码构建

- **优点**：
  - ✅ 可自定义修改代码
  - ✅ 适合开发和调试
- **使用场景**：开发环境、自定义需求

**开发者使用**：
```bash
docker compose -f docker-compose.build.yml up -d
```

---

## ✨ 特色功能

- 🚀 **YouTube 自动下载** - 发送链接即可自动提取音频并添加到音乐库
- 💾 **无限云存储** - 音乐文件存储在 Telegram 云端，VPS 仅保存元数据，零存储成本
- 🎯 **智能元数据** - 支持手动设置歌手国籍（22 个国家）、歌曲类型、语言分类
- ⭐ **收藏系统** - 收藏喜欢的歌曲，随时查看
- 📜 **播放历史** - 自动记录播放历史，支持回听
- 🎲 **随机播放** - 发现音乐库中的惊喜
- 🎨 **分类筛选** - 按类型和语言浏览歌曲
- 🖥️ **Web 管理后台** - 完整的歌曲管理和补档系统
- 🌍 **多语言支持** - 完美支持中文、日语、韩语等多字节字符

---

## 📱 快速使用

### 添加音乐

**方式一：YouTube 自动下载（推荐）**

发送 YouTube 链接给 Bot，自动下载音频：

```
https://www.youtube.com/watch?v=xxxxx
```

**方式二：发送音频文件**

直接在 Telegram 发送音频文件给 Bot，秒速保存。

**支持的格式：**
- ✅ MP3（推荐）
- ✅ M4A
- ✅ OGG
- ✅ 其他 Telegram 支持的音频格式

**文件命名规则：**
- 📁 文件名任意，系统会自动识别
- 🎵 推荐格式：`歌手 - 歌曲名.mp3`
- 🔍 如果文件名是 `周杰伦 - 稻香.mp3`，系统会自动识别歌手和歌名
- ⚠️ 文件大小限制：50MB

### 搜索播放

私聊 Bot，发送歌曲名或歌手名：

```
周杰伦 稻香
```

### 机器人命令

| 命令 | 功能 |
|------|------|
| `/start` | 欢迎信息 |
| `/help` | 帮助文档 |
| `/songs` 或 `/list` | 浏览音乐库（随机 10 首） |
| `/random` | 随机播放 |
| `/favorites` | 收藏列表 |
| `/history` | 播放历史 |
| `/stats` | 统计信息 |

### Web 管理后台

访问 `http://你的服务器IP:9999`，使用配置的用户名和密码登录。

**功能：**
- 查看所有歌曲
- 编辑歌曲信息（类型、语言、国籍）
- 删除歌曲
- 缺失歌曲补档

📖 详细使用说明：[查看使用手册](./使用说明.md)

---

## 🏗️ 技术架构

### 存储架构

**重要**：音乐文件全部存储在 Telegram 云端，不占用 VPS 存储空间！

```
┌─────────────┐
│   Telegram  │
│  Cloud CDN  │  ← 音乐文件实际存储位置
└──────┬──────┘
       │
       │ FileID 引用
       ▼
┌───────────────────────────────────┐
│      PostgreSQL Database          │
│  - 只保存元数据（约 1KB/首）      │
│  - FileID、标题、歌手等信息        │
└───────────────────────────────────┘
```

### 系统要求

- **操作系统**：Linux / macOS / Windows
- **内存**：建议 1GB+（音乐文件不占用本地空间）
- **磁盘空间**：建议 5GB+（仅用于临时下载和数据库）
- **网络**：稳定连接

### 核心技术栈

| 组件 | 技术 |
|------|------|
| 后端语言 | Go 1.21+ |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL 15+ |
| Bot SDK | telegram-bot-api |
| 视频下载 | yt-dlp |
| 音频处理 | ffmpeg |
| 容器化 | Docker & Docker Compose |

---

## ❓ 常见问题

### Q: 音乐会占用我的 VPS 空间吗？

**A:** 不会！音乐文件存储在 Telegram 云端，你的 VPS 只保存元数据（文件 ID、标题、歌手等），每首歌仅占用约 1KB 数据库空间。

### Q: 文件大小限制是多少？

**A:** 单个文件最大 50MB，这是 Telegram Bot API 的限制。

### Q: YouTube 下载失败怎么办？

**A:** 可能的原因：
1. 视频有地区限制或版权保护
2. 服务器网络无法访问 YouTube

**推荐方案**：使用在线 YouTube 转 MP3 工具转换后，直接发送 MP3 文件给 Bot。

### Q: 日语、韩语歌名会显示乱码吗？

**A:** 不会！系统已优化 UTF-8 多字节字符处理，完美支持中文、日语、韩语等所有语言。

### Q: 可以在群组中使用吗？

**A:** 可以！使用行内模式：`@BotName 关键词` 即可在任何群组中搜索播放。

---

## 📚 文档

- 📖 **[详细部署指南](./DEPLOY.md)** - 完整的部署和安装说明
- 🎵 **[使用手册](./使用说明.md)** - 详细的使用说明和技巧
- 🔧 **[配置详解](./DEPLOY.md#配置详解)** - 所有配置参数说明

---

## 🔄 更新日志

### v1.2.0 (2026-02-08)

#### 新增功能
- ✨ Web 后台支持手动选择歌手国籍（22 个国家/地区）
- ✨ 歌曲列表新增"国家"列，显示国家 Emoji
- ✨ 智能国家代码更新：手动选择优先，未选时根据语言自动设置

#### 优化改进
- 🐛 修复 Telegram 按钮中日语、韩语等多字节字符显示乱码的问题
- 🎨 优化 UTF-8 字符处理逻辑
- 📝 更新文档，明确说明文件存储在 Telegram 云端

### v1.1.0 (2026-02-08)

#### 新增功能
- ✨ 添加 `/songs` 命令，浏览音乐库中的歌曲
- ✨ 歌曲支持类型和语言分类
- ✨ Web 后台添加类型/语言筛选功能
- ✨ 编辑歌曲时可设置类型和语言
- ✨ YouTube 自动下载功能优化

#### 优化改进
- 🎨 优化所有命令的提示信息
- 🎨 更新帮助文档，更加详细友好
- 🐛 修复播放历史不记录的问题
- 🐛 修复 Web 后台编辑功能的体验

### v1.0.0 (2026-02-07)

#### 初始版本
- ✅ 基础 Bot 功能（搜索、播放、收藏、历史）
- ✅ YouTube 自动下载
- ✅ Web 管理后台
- ✅ 元数据识别（地区、年份）
- ✅ 随机播放
- ✅ 统计信息

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发环境搭建

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/fish-music.git
cd fish-music

# 2. 复制配置文件
cp config.yaml.example config.yaml
# 编辑配置文件

# 3. 启动数据库
docker compose up -d postgres

# 4. 运行 Bot
go run cmd/bot/main.go

# 5. 运行 Web 服务
go run cmd/web/main.go
```

---

## 📄 许可证

MIT License

Copyright (c) 2026 Fish Music

---

## 🙏 鸣谢

感谢以下开源项目：

- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API SDK for Go
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - 视频下载工具
- [GORM](https://github.com/go-gorm/gorm) - Go 语言 ORM 库
- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架

---

## 📞 联系方式

- **问题反馈**：[GitHub Issues](https://github.com/yourusername/fish-music/issues)
- **功能建议**：[GitHub Discussions](https://github.com/yourusername/fish-music/discussions)

---

<div align="center">

**🎵 Fish Music - 你的个人云端音乐库**

[⭐ Star](https://github.com/yourusername/fish-music/stargazers) | [🍴 Fork](https://github.com/yourusername/fish-music/network/members) | [📖 文档](./DEPLOY.md)

Made with ❤️ by Fish Music Team

</div>
