# YouTube Cookies 配置指南

## 为什么需要 Cookies？

YouTube 会检测来自服务器的下载请求，可能会显示以下错误：

```
ERROR: [youtube] XXXXX: Sign in to confirm you're not a bot.
Use --cookies-from-browser or --cookies for the authentication.
```

为了解决这个问题，你可以提供 YouTube cookies 文件。

---

## 快速开始（推荐方法）

### 使用 Python 脚本导出 Cookies

**最简单的方式！不需要安装浏览器插件，只需要运行一个脚本。**

#### 1. 安装依赖

在你的**本地电脑**上（不是服务器）：

```bash
pip install browser_cookie3
```

如果提示 `pip` 不存在，先安装 Python：
- **Windows**: 从 https://python.org 下载安装
- **macOS**: `brew install python3`
- **Linux**: `sudo apt install python3-pip`

#### 2. 下载导出脚本

从项目目录获取 `export-cookies.py` 脚本，或者直接创建：

```bash
# 方法1: 如果已克隆项目
cp /opt/fish_music/export-cookies.py .

# 方法2: 下载脚本
wget https://raw.githubusercontent.com/qqzhoufan/fish_music/main/export-cookies.py
```

#### 3. 运行脚本导出 Cookies

**Chrome 用户**（推荐）:

```bash
python export-cookies.py chrome
```

**Firefox 用户**:

```bash
python export-cookies.py firefox
```

**Edge 用户**:

```bash
python export-cookies.py edge
```

**其他浏览器**:
- Safari (macOS): `python export-cookies.py safari`
- Brave: `python export-cookies.py brave`
- Opera: `python export-cookies.py opera`

#### 4. 上传到服务器

脚本会生成 `youtube-cookies.txt` 文件，将其上传到服务器的部署目录（与 `docker-compose.yml` 同级）：

```bash
# 使用 scp 上传（在你的本地电脑上执行）
scp youtube-cookies.txt root@your-server:/path/to/fish-music/
```

或者使用 FTP/SFTP 工具上传。

---

## 手动导出方法（不推荐）

如果无法运行 Python 脚本，可以使用以下方法：

### 方法 1: 使用浏览器开发者工具

1. **打开 YouTube 并登录**
   - 访问 https://www.youtube.com
   - 确保已登录你的账号

2. **打开开发者工具**
   - Chrome/Edge: 按 `F12` 或 `Ctrl+Shift+I` (Mac: `Cmd+Option+I`)
   - Firefox: 按 `F12` 或 `Ctrl+Shift+K`

3. **进入 Console 标签**

4. **复制并运行以下代码**（在 Console 中粘贴并回车）:

```javascript
// 复制这段代码并粘贴到浏览器 Console
document.cookie.split(';').forEach(c => {
    let [name, value] = c.trim().split('=');
    console.log(`${name}=${value}`);
});
```

5. **手动创建 cookies 文件**
   - 将输出的内容保存为 `youtube-cookies.txt`
   - 注意：这种方法可能不完整，建议使用上面的 Python 脚本

---

## 配置步骤

### 1. 确认文件位置

确保 `youtube-cookies.txt` 在部署目录中：

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

编辑 `config.yaml`，在 `download` 部分添加：

```yaml
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

查看 Bot 日志，确认没有验证错误：

```bash
docker compose logs -f bot
```

测试发送一个 YouTube 链接，应该能正常下载了！

---

## 常见问题

### Q: 脚本运行失败："无法读取 cookies"

**解决方案**：
1. 确保浏览器已关闭（Chrome 需要完全关闭，检查任务管理器）
2. 确保已经登录 YouTube
3. 尝试使用不同的浏览器（Firefox 通常更稳定）

### Q: 脚本运行失败："缺少依赖库"

**解决方案**：
```bash
pip install browser_cookie3
```

如果还是失败，尝试升级 pip：
```bash
python -m pip install --upgrade pip
pip install browser_cookie3
```

### Q: Cookies 文件导出成功，但仍然报错

**可能原因**：
1. Cookies 已过期（YouTube cookies 有效期通常为数周到数月）
2. 导出时浏览器中的 cookies 不完整
3. 文件格式错误

**解决方案**：
```bash
# 检查文件大小（正常应该是几KB到几十KB）
ls -lh youtube-cookies.txt

# 查看文件内容（应该包含多行数据）
head -5 youtube-cookies.txt
```

### Q: Windows 上运行脚本出错

**解决方案**：
```bash
# 确保使用 python 而不是 python3
python export-cookies.py chrome

# 或者明确指定 Python 版本
python3 export-cookies.py chrome
```

### Q: 不想使用 cookies，有其他方案吗？

**替代方案**：使用在线转换工具
1. 访问 https://y2mate.com 或 https://yt1s.com
2. 将 YouTube 链接转换为 MP3
3. 直接发送 MP3 文件给 Bot

---

## 安全建议

1. **隐私保护**：
   - Cookies 文件包含你的登录凭证
   - 不要分享给他人
   - 不要提交到 Git 仓库
   - 定期更新（建议每月更新一次）

2. **使用专用账号**：
   - 建议使用单独的 YouTube 账号
   - 不要使用个人主账号

3. **权限控制**：
   ```bash
   chmod 600 youtube-cookies.txt
   ```

---

## 技术说明

### Cookies 有效期

- YouTube cookies 通常有效期为 **数周到数月**
- 如果下载再次失败，重新导出 cookies 即可
- Bot 会在 cookies 过期时提示你重新配置

### 文件格式

脚本导出的是 **Netscape Cookie 格式**，这是 yt-dlp 官方支持的格式：

```
# Netscape HTTP Cookie File
.youtube.com	TRUE	/	TRUE	1234567890	SID	xxxxx
.youtube.com	TRUE	/	TRUE	1234567890	HSID	xxxxx
```

### 支持的浏览器

| 浏览器 | 支持状态 | 备注 |
|--------|---------|------|
| Chrome | ✅ | 包括 Edge、Brave、Opera |
| Firefox | ✅ | 完全支持 |
| Safari | ✅ | 仅 macOS |
| Edge | ✅ | 基于 Chromium |

---

## 更新日志

- **v1.3.1** - 添加 Python 导出脚本，无需浏览器插件
- **v1.3.0** - 新增 cookies_file 配置选项，支持 YouTube cookies 认证
