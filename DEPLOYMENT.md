# Fish Music 部署指南

## 快速开始

### 1. 使用 Docker Compose 部署（推荐）

这是最简单的部署方式，适合大多数场景。

#### 步骤 1: 克隆代码

```bash
git clone <repository-url> fish_music
cd fish_music
```

#### 步骤 2: 配置 Bot

1. 从 [@BotFather](https://t.me/BotFather) 获取 Bot Token
2. 获取你的 Telegram User ID（可以从 [@userinfobot](https://t.me/userinfobot) 获取）
3. 复制配置文件模板：

```bash
cp config.yaml.example config.yaml
```

4. 编辑 `config.yaml`，填入必要信息：

```yaml
bot:
  token: "YOUR_BOT_TOKEN_HERE"  # 从 BotFather 获取
  admin_id: 123456789           # 你的 Telegram ID

database:
  host: "postgres"              # Docker 容器名称
  port: 5432
  user: "fish_music"
  password: "fish_music_pass"
  dbname: "fish_music"

web:
  port: 9999
  username: "admin"
  password: "your_secure_password"
```

#### 步骤 3: 一键部署

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

或者使用 Make：

```bash
make docker-up
```

#### 步骤 4: 验证部署

1. 访问 Web 管理端: http://your-server-ip:9999
2. 在 Telegram 中搜索你的 Bot，发送 `/start` 测试

### 2. 手动部署

如果你需要更灵活的部署方式，可以手动部署各个组件。

#### 前置要求

- Go 1.21+
- PostgreSQL 12+
- FFmpeg
- yt-dlp

#### 安装依赖

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y golang postgresql ffmpeg python3-pip
pip3 install yt-dlp

# macOS
brew install go postgresql ffmpeg yt-dlp
```

#### 数据库初始化

```bash
# 创建数据库和用户
sudo -u postgres psql
CREATE DATABASE fish_music;
CREATE USER fish_music WITH PASSWORD 'fish_music_pass';
GRANT ALL PRIVILEGES ON DATABASE fish_music TO fish_music;
\q

# 导入初始表结构
psql -h localhost -U fish_music -d fish_music -f sql/init.sql
```

#### 构建和运行

```bash
# 安装 Go 依赖
go mod download

# 构建
make build

# 运行 Bot
./bin/bot

# 运行 Web 管理端（另开终端）
./bin/web
```

### 3. 使用 systemd 服务（生产环境）

创建服务文件 `/etc/systemd/system/fish-music-bot.service`：

```ini
[Unit]
Description=Fish Music Bot
After=network.target postgresql.service

[Service]
Type=simple
User=fishmusic
WorkingDirectory=/opt/fish_music
ExecStart=/opt/fish_music/bin/bot
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

创建服务文件 `/etc/systemd/system/fish-music-web.service`：

```ini
[Unit]
Description=Fish Music Web Admin
After=network.target postgresql.service

[Service]
Type=simple
User=fishmusic
WorkingDirectory=/opt/fish_music
ExecStart=/opt/fish_music/bin/web
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable fish-music-bot
sudo systemctl enable fish-music-web
sudo systemctl start fish-music-bot
sudo systemctl start fish-music-web
```

## 配置说明

### config.yaml 参数详解

```yaml
# Telegram Bot 配置
bot:
  token: "..."      # Bot Token（必填）
  admin_id: 123456  # 管理员 ID（必填）

# 数据库配置
database:
  host: "localhost"
  port: 5432
  user: "fish_music"
  password: "..."    # 数据库密码
  dbname: "fish_music"
  sslmode: "disable" # 生产环境建议使用 "require"

# Web 管理端配置
web:
  port: 9999              # Web 服务端口
  username: "admin"       # Basic Auth 用户名
  password: "..."         # Basic Auth 密码

# 下载配置
download:
  worker_count: 3         # 并发下载数（1-5）
  max_file_size: 50       # 最大文件大小（MB）
  temp_dir: "./tmp"       # 临时文件目录

# 搜索 API 配置
search:
  api_url: "..."          # 音乐搜索 API（可选）
  timeout: 30             # 搜索超时（秒）

# 日志配置
log:
  level: "info"           # debug / info / warn / error
  file: ""                # 留空输出到标准输出，或指定日志文件路径
```

## 常见问题

### 1. Bot 无法启动

- 检查 Bot Token 是否正确
- 检查网络连接是否正常
- 查看日志: `make docker-logs`

### 2. 数据库连接失败

- 检查 PostgreSQL 是否运行: `docker-compose ps`
- 检查数据库配置是否正确
- 检查防火墙设置

### 3. 下载失败

- 确保 yt-dlp 已安装: `yt-dlp --version`
- 检查网络连接
- 检查临时目录权限

### 4. Web 管理端无法访问

- 检查端口是否被占用
- 检查防火墙设置
- 检查 Basic Auth 用户名密码

## 维护命令

```bash
# 查看日志
make docker-logs

# 重启服务
docker-compose restart

# 停止服务
make docker-down

# 更新代码
git pull
docker-compose down
docker-compose build
docker-compose up -d

# 数据库备份
docker-compose exec postgres pg_dump -U fish_music fish_music > backup.sql

# 数据库恢复
docker-compose exec -T postgres psql -U fish_music fish_music < backup.sql

# 清理临时文件
rm -rf tmp/*
```

## 安全建议

1. **修改默认密码**: 修改数据库密码和 Web 管理端密码
2. **使用 SSL**: 生产环境建议使用 HTTPS 访问 Web 管理端
3. **限制访问**: 使用防火墙限制数据库和 Web 端口访问
4. **定期备份**: 定期备份数据库
5. **更新依赖**: 定期更新系统和依赖包

## 性能优化

1. **调整并发数**: 根据服务器性能调整 `worker_count`
2. **使用 CDN**: 将静态资源托管到 CDN
3. **数据库优化**: 定期清理历史记录，优化索引
4. **缓存策略**: 可以添加 Redis 缓存热门歌曲

## 监控建议

1. 使用 `docker stats` 监控容器资源使用
2. 配置日志收集（如 ELK Stack）
3. 设置告警（磁盘空间、服务状态等）
