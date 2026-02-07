# YouTube Cookies 配置指南

## 为什么需要 Cookies？

YouTube 会检测来自服务器的下载请求，可能会显示以下错误：

```
ERROR: [youtube] XXXXX: Sign in to confirm you're not a bot.
Use --cookies-from-browser or --cookies for the authentication.
```

为了解决这个问题，你可以提供 YouTube cookies 文件。

## 获取 YouTube Cookies（三种方法）

### 方法 1：使用浏览器扩展（推荐）

#### Firefox 用户

1. **安装扩展**：
   - 安装 "Get cookies.txt LOCALLY" 扩展
   - 下载地址：https://addons.mozilla.org/en-US/firefox/addon/get-cookiestxt-locally/

2. **导出 Cookies**：
   - 访问 https://www.youtube.com
   - 确保已登录 YouTube 账号
   - 点击浏览器工具栏的扩展图标
   - 选择 "Export" → "Cleanup" → "Download"
   - 保存为 `youtube-cookies.txt`

#### Chrome/Edge 用户

1. **安装扩展**：
   - 安装 "Get cookies.txt LOCALLY" 扩展
   - Chrome：https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc
   - Edge：https://microsoftedge.microsoft.com/addons/detail/get-cookiestxt-locally/ejafafmemmkgbnfmdmjliifenbopcglg

2. **导出 Cookies**：
   - 访问 https://www.youtube.com
   - 确保已登录 YouTube 账号
   - 点击扩展图标
   - 点击 "Current Site" → "Export" → "Download"
   - 保存为 `youtube-cookies.txt`

### 方法 2：使用浏览器开发者工具

1. **打开 YouTube**：
   - 访问 https://www.youtube.com 并登录

2. **打开开发者工具**：
   - Windows/Linux: 按 `F12` 或 `Ctrl+Shift+I`
   - Mac: 按 `Cmd+Option+I`

3. **进入 Application 标签**：
   - 点击顶部的 "Application" 或 "应用程序" 标签
   - 在左侧找到 "Cookies" → "https://www.youtube.com"

4. **安装 yt-dlp（本地电脑）**：
   ```bash
   pip install yt-dlp
   ```

5. **使用浏览器导出命令**（在本地电脑执行）：
   ```bash
   # Chrome
   yt-dlp --cookies-from-browser chrome https://www.youtube.com/watch?v=VIDEO_ID

   # Firefox
   yt-dlp --cookies-from-browser firefox https://www.youtube.com/watch?v=VIDEO_ID

   # Edge
   yt-dlp --cookies-from-browser edge https://www.youtube.com/watch?v=VIDEO_ID
   ```

### 方法 3：手动创建 Cookies 文件（不推荐）

1. 使用上述方法 2 中的浏览器扩展
2. 或者使用在线工具导出 cookies

## 配置步骤

### 1. 放置 Cookies 文件

将 `youtube-cookies.txt` 文件放在你的部署目录中（与 `docker-compose.yml` 同级）：

```bash
your-deployment-directory/
├── config.yaml
├── docker-compose.yml
├── sql/
│   └── init.sql
├── tmp/
└── youtube-cookies.txt  ← 放在这里
```

### 2. 配置 config.yaml

在 `config.yaml` 中添加或修改 `cookies_file` 配置：

```yaml
# 下载配置
download:
  worker_count: 3
  max_file_size: 50
  temp_dir: "./tmp"
  cookies_file: "/app/youtube-cookies.txt"  # 启用 cookies
```

### 3. 重启服务

```bash
docker compose down
docker compose up -d
```

### 4. 验证配置

查看 Bot 日志，确认没有 "Sign in to confirm you're not a bot" 错误：

```bash
docker compose logs -f bot
```

## 不使用 Cookies 会怎样？

- 部分视频可能仍然可以下载
- 某些视频可能会遇到验证错误
- 对于受保护的内容或年龄限制内容，下载会失败

## 安全建议

1. **隐私保护**：
   - Cookies 文件包含你的登录凭证
   - 不要分享给他人
   - 不要提交到 Git 仓库
   - 定期更新（建议每月更新一次）

2. **权限控制**：
   ```bash
   chmod 600 youtube-cookies.txt
   ```

3. **使用专用账号**：
   - 建议使用单独的 YouTube 账号
   - 不要使用个人主账号

## 故障排除

### 问题 1：仍然显示验证错误

**解决方案**：
- Cookies 可能已过期，重新导出
- 确保导出时已登录 YouTube
- 检查文件路径是否正确

### 问题 2：文件格式错误

**解决方案**：
- 确保使用正确的扩展导出
- 文件格式应为 Netscape cookie 格式
- 第一行应该包含 `# Netscape HTTP Cookie File`

### 问题 3：下载其他视频正常，某个视频失败

**可能原因**：
- 该视频有年龄限制
- 该视频有地区限制
- 该视频需要付费订阅

## 替代方案

如果不想配置 Cookies，可以直接发送 MP3 文件给 Bot：

1. 在本地下载音频（使用其他工具）
2. 直接发送 MP3 文件给 Bot
3. Bot 会自动保存并提供搜索

## 更新日志

- **v1.3.0** - 新增 cookies_file 配置选项，支持 YouTube cookies 认证
