# Fish Music 部署指南

> 详细的部署和安装说明，帮助你在自己的服务器上快速搭建 Fish Music 音乐机器人

[![Deploy](https://img.shields.io/badge/Deployment-Easy-success)]()
[![Time](https://img.shields.io/badge/Time-10min-blue)]()
[![Difficulty](https://img.shields.io/badge/Difficulty-Easy-green)]()

---

## 目录

- [系统要求](#系统要求)
- [部署前准备](#部署前准备)
- [方式一：Docker 部署（推荐）](#方式一docker-部署推荐)
- [方式二：手动部署](#方式二手动部署)
- [配置详解](#配置详解)
- [常见部署问题](#常见部署问题)
- [更新与维护](#更新与维护)
- [卸载说明](#卸载说明)

---

## 系统要求

### 最低配置

| 资源 | 最低要求 | 推荐配置 |
|------|---------|---------|
| 操作系统 | Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+) | Ubuntu 22.04 LTS |
| CPU | 1 核心 | 2 核心 |
| 内存 | 1 GB | 2 GB |
| 磁盘空间 | 5 GB | 10 GB |
| 网络 | 稳定连接 | 较快带宽 |

### 软件要求

- **Docker**: 20.10+
- **Docker Compose**: 2.0+

### 为什么配置要求这么低？

因为 Fish Music 使用 Telegram 云存储，所有音乐文件都存储在 Telegram 服务器上，你的服务器只保存元数据（每首歌约 1KB），所以几乎不占用存储空间！

---

## 部署前准备

### 1. 购买服务器（可选）

如果你没有服务器，可以考虑：
- **Vultr**：https://www.vultr.com (推荐，性价比高)
- **DigitalOcean**：https://www.digitalocean.com
- **腾讯云**：https://cloud.tencent.com
- **阿里云**：https://www.aliyun.com

**推荐配置**：1 核心 / 1GB 内存 / 10GB 磁盘（约 $5/月）

### 2. 获取 Telegram Bot Token

#### 步骤 1：创建 Bot

1. 在 Telegram 中搜索 [@BotFather](https://t.me/BotFather)
2. 发送 `/newbot` 命令
3. 按提示输入机器人名称（例如：`MyMusicBot`）
4. 输入机器人用户名（例如：`my_music_bot`，必须以 `_bot` 结尾）

#### 步骤 2：获取 Token

BotFather 会返回类似这样的 Token：

```
1234567890:ABCdefGhIJKlmNoPQRsTUVwxyZ-abc123
```

**⚠️ 重要**：请妥善保管这个 Token，不要泄露给他人！

### 3. 获取你的 Telegram User ID

1. 在 Telegram 中搜索 [@userinfobot](https://t.me/userinfobot)
2. 发送 `/start` 命令
3. 机器人会返回你的 ID（纯数字，例如：`123456789`）

**⚠️ 重要**：记下这个 ID，你将成为机器人的管理员！

---

## 方式一：Docker 部署（推荐）

### 为什么推荐 Docker？

- ✅ 一键部署，无需复杂配置
- ✅ 环境隔离，不会影响系统
- ✅ 自动安装所有依赖（ffmpeg, yt-dlp, PostgreSQL）
- ✅ 方便更新和维护

### 部署步骤

#### 步骤 1：安装 Docker 和 Docker Compose

**Ubuntu/Debian:**

```bash
# 更新软件包
sudo apt update

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo apt install docker-compose-plugin

# 验证安装
docker --version
docker compose version
```

**CentOS/RHEL:**

```bash
# 安装 Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 步骤 2：获取部署文件

**选项 A：使用 Docker Hub 镜像（推荐，最快）**

```bash
# 方法一：一键部署脚本（最简单）
curl -fsSL https://raw.githubusercontent.com/qqzhoufan/fish_music/main/deploy.sh -o deploy.sh
chmod +x deploy.sh
./deploy.sh
```

```bash
# 方法二：手动下载文件
mkdir fish-music && cd fish-music

# 下载配置文件
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/config.yaml.example
mv config.yaml.example config.yaml

# 下载 docker-compose.yml
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/docker-compose.yml

# 创建 sql 目录并下载初始化脚本（重要！）
mkdir -p sql
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/sql/init.sql -O sql/init.sql

# 创建临时目录
mkdir -p tmp
```

**选项 B：从 GitHub 克隆（完整源码）**

```bash
# 克隆项目
git clone https://github.com/qqzhoufan/fish_music.git
cd fish-music
```

#### 步骤 3：配置文件

```bash
# 复制配置模板
cp config.yaml.example config.yaml

# 编辑配置文件
nano config.yaml
# 或使用 vim
vim config.yaml
```

**配置内容示例：**

```yaml
# Telegram Bot 配置
bot:
  token: "1234567890:ABCdefGhIJKlmNoPQRsTUVwxyZ"  # 替换为你的 Bot Token
  admin_id: 123456789                             # 替换为你的 Telegram ID

# 数据库配置（Docker 部署保持默认即可）
database:
  host: "postgres"
  port: 5432
  user: "fish_music"
  password: "fish_music_pass"
  dbname: "fish_music"
  sslmode: "disable"

# Web 管理端配置
web:
  port: 9999                     # Web 服务端口
  username: "admin"              # 登录用户名（建议修改）
  password: "fishmusic2024"      # 登录密码（建议修改）

# 下载配置
download:
  worker_count: 3                # 并发下载数
  max_file_size: 50              # 最大文件大小（MB）
  temp_dir: "./tmp"              # 临时文件目录

# 日志配置
log:
  level: "info"                  # 日志级别
  file: ""                       # 留空输出到控制台
```

#### 步骤 4：启动服务

```bash
# 构建并启动所有服务
docker compose up -d

# 查看运行状态
docker compose ps

# 查看日志（首次启动会看到数据库初始化）
docker compose logs -f
```

#### 步骤 5：验证部署

**测试 Bot:**

1. 在 Telegram 中找到你的 Bot
2. 发送 `/start` 命令
3. 如果收到欢迎消息，说明 Bot 运行正常！

**测试 Web 面板:**

1. 浏览器访问：`http://你的服务器IP:9999`
2. 输入配置文件中设置的用户名和密码
3. 如果能看到管理面板，说明 Web 服务正常！

---

## 方式二：手动部署

如果你不想使用 Docker，可以手动部署：

### 步骤 1：安装依赖

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install -y \
    postgresql \
    postgresql-contrib \
    ffmpeg \
    yt-dlp \
    golang-1.21-go
```

**CentOS/RHEL:**

```bash
sudo yum install -y \
    postgresql-server \
    postgresql-contrib \
    ffmpeg \
    yt-dlp \
    golang
```

### 步骤 2：安装 Go

```bash
# 下载 Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# 解压安装
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证
go version
```

### 步骤 3：配置数据库

```bash
# 初始化数据库
sudo postgresql-setup initdb

# 启动 PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库和用户
sudo -u postgres psql << EOF
CREATE USER fish_music WITH PASSWORD 'fish_music_pass';
CREATE DATABASE fish_music OWNER fish_music;
GRANT ALL PRIVILEGES ON DATABASE fish_music TO fish_music;
\q
EOF

# 导入初始化脚本
psql -U fish_music -d fish_music -f sql/init.sql
```

### 步骤 4：编译项目

```bash
# 克隆项目
git clone https://github.com/yourusername/fish-music.git
cd fish-music

# 配置文件
cp config.yaml.example config.yaml
nano config.yaml  # 修改配置（同 Docker 部署）

# 修改数据库连接为 localhost
# database:
#   host: "localhost"

# 编译
go build -o bin/bot ./cmd/bot
go build -o bin/web ./cmd/web
```

### 步骤 5：启动服务

**使用 systemd 管理 Bot:**

```bash
# 创建 systemd 服务文件
sudo nano /etc/systemd/system/fish-music-bot.service
```

**服务文件内容：**

```ini
[Unit]
Description=Fish Music Bot
After=network.target postgresql.service

[Service]
Type=simple
User=your_username
WorkingDirectory=/path/to/fish-music
ExecStart=/path/to/fish-music/bin/bot
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**启动服务：**

```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启动 Bot
sudo systemctl start fish-music-bot

# 开机自启
sudo systemctl enable fish-music-bot

# 查看状态
sudo systemctl status fish-music-bot

# 查看日志
sudo journalctl -u fish-music-bot -f
```

**Web 服务类似，创建 `fish-music-web.service`**

---

## 配置详解

### Bot Token 配置

```yaml
bot:
  token: "1234567890:ABCdefGhIJKlmNoPQRsTUVwxyZ"
  admin_id: 123456789
```

| 参数 | 说明 | 获取方式 |
|------|------|---------|
| token | Bot 令牌 | 从 @BotFather 获取 |
| admin_id | 管理员 ID | 从 @userinfobot 获取 |

### 数据库配置

```yaml
database:
  host: "postgres"      # Docker: postgres, 手动: localhost
  port: 5432
  user: "fish_music"
  password: "fish_music_pass"
  dbname: "fish_music"
  sslmode: "disable"
```

### Web 管理端配置

```yaml
web:
  port: 9999
  username: "admin"        # 建议修改
  password: "fishmusic2024" # 建议修改为强密码
```

**⚠️ 安全提示**：
- 默认密码仅用于测试
- 生产环境务必修改为强密码
- 建议配置防火墙，限制 Web 端口访问

### 下载配置

```yaml
download:
  worker_count: 3          # 并发下载数（建议 2-5）
  max_file_size: 50        # 最大文件大小（MB）
  temp_dir: "./tmp"        # 临时文件目录
```

| 参数 | 说明 | 建议 |
|------|------|------|
| worker_count | 并发下载数 | 服务器性能好可设置为 5 |
| max_file_size | 文件大小限制 | Telegram 限制为 50MB |
| temp_dir | 临时目录 | 确保有足够磁盘空间 |

---

## 常见部署问题

### Q1: Docker 启动失败

**错误信息**：`Cannot connect to the Docker daemon`

**解决方案**：

```bash
# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 检查状态
sudo systemctl status docker
```

### Q2: 数据库连接失败

**错误信息**：`connection refused` 或 `database does not exist`

**解决方案（Docker）**：

```bash
# 检查数据库容器状态
docker compose ps

# 查看数据库日志
docker compose logs postgres

# 重启数据库
docker compose restart postgres
```

**解决方案（手动部署）**：

```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 检查数据库是否存在
sudo -u postgres psql -l

# 手动创建数据库
sudo -u postgres createdb fish_music
```

### Q3: Bot 无响应

**检查步骤**：

```bash
# 1. 查看 Bot 日志
docker compose logs bot | tail -50

# 2. 检查 Bot Token 是否正确
docker compose exec bot cat /app/config.yaml | grep token

# 3. 重启 Bot
docker compose restart bot

# 4. 检查网络连接
curl -I https://api.telegram.org
```

### Q4: YouTube 下载失败

**可能原因**：
- 服务器网络无法访问 YouTube
- yt-dlp 版本过旧

**解决方案**：

```bash
# 进入容器更新 yt-dlp
docker compose exec bot yt-dlp --update

# 或重启容器让系统自动更新
docker compose restart bot
```

### Q5: 端口被占用

**错误信息**：`port is already allocated`

**解决方案**：

```bash
# 查看端口占用
sudo netstat -tlnp | grep 9999

# 修改配置文件中的端口号
nano config.yaml
# 将 web.port 改为其他端口（如 9998）

# 重启服务
docker compose up -d
```

### Q6: 内存不足

**错误信息**：`Cannot allocate memory`

**解决方案**：

```bash
# 1. 创建 Swap 空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 2. 永久生效
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# 3. 减少并发下载数
nano config.yaml
# 将 download.worker_count 改为 1
```

---

## 更新与维护

### 更新到最新版本

**Docker 部署：**

```bash
# 1. 备份数据库
docker exec fish_music_db pg_dump -U fish_music fish_music > backup_$(date +%Y%m%d).sql

# 2. 停止服务
docker compose down

# 3. 拉取最新代码
git pull

# 4. 重新构建镜像
docker compose build --no-cache

# 5. 启动服务
docker compose up -d

# 6. 查看日志
docker compose logs -f
```

**手动部署：**

```bash
# 1. 备份数据库
pg_dump -U fish_music fish_music > backup_$(date +%Y%m%d).sql

# 2. 停止服务
sudo systemctl stop fish-music-bot
sudo systemctl stop fish-music-web

# 3. 拉取最新代码
git pull

# 4. 重新编译
go build -o bin/bot ./cmd/bot
go build -o bin/web ./cmd/web

# 5. 启动服务
sudo systemctl start fish-music-bot
sudo systemctl start fish-music-web
```

### 备份数据库

**定期备份脚本：**

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/path/to/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
docker exec fish_music_db pg_dump -U fish_music fish_music | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# 删除 30 天前的备份
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

**设置定时任务：**

```bash
# 编辑 crontab
crontab -e

# 每天凌晨 3 点备份
0 3 * * * /path/to/backup.sh
```

### 查看日志

**Docker 部署：**

```bash
# 实时查看所有日志
docker compose logs -f

# 查看 Bot 日志
docker compose logs -f bot

# 查看 Web 日志
docker compose logs -f web

# 查看最近 100 行日志
docker compose logs --tail=100 bot
```

**手动部署：**

```bash
# Bot 日志
sudo journalctl -u fish-music-bot -f

# Web 日志
sudo journalctl -u fish-music-web -f

# 查看最近 100 行
sudo journalctl -u fish-music-bot -n 100
```

---

## 卸载说明

### Docker 部署卸载

```bash
# 1. 停止并删除容器
docker compose down

# 2. 删除数据卷（会删除数据库数据！）
docker volume rm fish_music_postgres_data

# 3. 删除项目文件
cd ..
rm -rf fish-music
```

### 手动部署卸载

```bash
# 1. 停止服务
sudo systemctl stop fish-music-bot
sudo systemctl stop fish-music-web

# 2. 禁用开机自启
sudo systemctl disable fish-music-bot
sudo systemctl disable fish-music-web

# 3. 删除服务文件
sudo rm /etc/systemd/system/fish-music-bot.service
sudo rm /etc/systemd/system/fish-music-web.service

# 4. 重载 systemd
sudo systemctl daemon-reload

# 5. 删除数据库（可选）
sudo -u postgres psql -c "DROP DATABASE fish_music;"
sudo -u postgres psql -c "DROP USER fish_music;"

# 6. 删除项目文件
cd ..
rm -rf fish-music
```

---

## 性能优化建议

### 1. 启用 PostgreSQL 缓存

```bash
# 编辑 PostgreSQL 配置
docker compose exec postgres nano /var/lib/postgresql/data/postgresql.conf

# 添加以下配置
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1

# 重启数据库
docker compose restart postgres
```

### 2. 配置反向代理（可选）

使用 Nginx 作为反向代理：

```nginx
# /etc/nginx/sites-available/fish-music

server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:9999;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 配置 SSL 证书（推荐）

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

---

## 安全建议

### 1. 修改默认密码

```bash
# 编辑配置文件
nano config.yaml

# 修改以下字段
web:
  username: "your_username"    # 改为自定义用户名
  password: "strong_password"  # 改为强密码
```

### 2. 配置防火墙

```bash
# UFW 防火墙
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# 限制 Web 端口访问（仅允许特定 IP）
sudo ufw allow from YOUR_IP_ADDRESS to any port 9999
```

### 3. 定期更新

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 更新 Docker
docker compose pull
docker compose up -d --build
```

---

## 监控与告警

### 使用 Docker Health Check

```yaml
# docker-compose.yml

services:
  bot:
    # ... 其他配置
    healthcheck:
      test: ["CMD", "pgrep", "bot"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### 查看服务状态

```bash
# 检查所有容器状态
docker compose ps

# 查看资源使用
docker stats
```

---

## 📁 文件格式与命名说明

### 支持的音频格式

Fish Music 通过 Telegram Bot API 发送音频，支持以下格式：

| 格式 | 扩展名 | 推荐度 | 说明 |
|------|--------|--------|------|
| MP3 | `.mp3` | ⭐⭐⭐⭐⭐ | 最推荐，兼容性最好 |
| M4A | `.m4a` | ⭐⭐⭐⭐ | Apple 设备常用格式 |
| OGG | `.ogg` | ⭐⭐⭐ | 开源音频格式 |
| 其他 | - | ⭐⭐ | Telegram 支持的其他音频格式 |

### 文件命名规则

#### 推荐命名格式

```
歌手 - 歌曲名.mp3
```

**示例：**
- ✅ `周杰伦 - 稻香.mp3`
- ✅ `邓紫棋 - 光年之外.mp3`
- ✅ `Adele - Hello.mp3`

#### 其他支持的命名格式

系统会智能识别包含以下分隔符的文件名：

| 分隔符 | 示例 | 识别结果 |
|--------|------|----------|
| ` - ` | `周杰伦 - 稻香.mp3` | 歌手：周杰伦，歌名：稻香 |
| `.` | `周杰伦.稻香.mp3` | 歌手：周杰伦，歌名：稻香 |
| `_` | `周杰伦_稻香.mp3` | 歌手：周杰伦，歌名：稻香 |

#### 不推荐的命名

- ⚠️ `random_music.mp3` - 无法识别歌手和歌名
- ⚠️ `song1.mp3` - 无法识别歌手和歌名
- ⚠️ `音乐.mp3` - 无法识别歌手和歌名

**提示**：即使文件名不规范，文件也能正常保存，只是歌手和歌名会显示为文件名。你可以稍后在 Web 管理后台修改。

### 文件大小限制

- **单个文件**：最大 50MB（Telegram Bot API 限制）
- **推荐大小**：3MB - 10MB（平衡音质和大小）
- **常见比特率**：128kbps - 320kbps

### 如何获取音频文件

#### 方法一：在线转换工具

**YouTube 转 MP3：**
- ytmp3.cc
- y2mate.com
- 320ytmp3.com

**使用步骤：**
1. 复制 YouTube 视频链接
2. 粘贴到转换网站
3. 选择 MP3 格式和比特率
4. 下载转换后的文件
5. 发送给 Fish Music Bot

#### 方法二：本地音乐库

直接从你的电脑或手机选择音频文件发送给 Bot。

#### 方法三：其他音乐平台

从 QQ 音乐、网易云音乐、酷狗等平台下载后发送。

---

## 获取帮助

如果遇到问题：

1. **查看日志**：`docker compose logs -f`
2. **检查配置**：确认 `config.yaml` 配置正确
3. **查看 Issues**：[GitHub Issues](https://github.com/yourusername/fish-music/issues)
4. **提交问题**：详细描述问题并提供日志

---

## 下一步

部署完成后：

1. ✅ 访问 Web 管理面板：`http://你的服务器IP:9999`
2. ✅ 在 Telegram 中找到你的 Bot，发送 `/start`
3. ✅ 尝试添加第一首歌曲
4. ✅ 查看 [使用说明](./使用说明.md) 了解更多功能

---

**祝部署顺利！享受你的云端音乐体验！** 🎵

如有问题，欢迎反馈！
