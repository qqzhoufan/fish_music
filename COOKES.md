# YouTube Cookies 配置指南

## 为什么需要 Cookies？

YouTube 会检测来自服务器的下载请求，可能会显示以下错误：

```
ERROR: [youtube] XXXXX: Sign in to confirm you're not a bot.
```

配置 cookies 可以解决这个问题。

---

## 🚀 最简单方法：通过 Bot 配置（推荐）

**无需 SSH、无需手动编辑文件！**

### 步骤：

1. **获取 Cookie**
   - 打开 https://www.youtube.com 并登录
   - 按 `F12` 打开开发者工具
   - 点击 "Application" → "Cookies" → "https://www.youtube.com"
   - 找到 `__Secure-3PSID` 或 `SID`
   - 双击 Value 列，复制整个值

2. **发送给 Bot**
   - 在 Telegram 中给你的 Bot 发送：
   ```
   /cookies <复制的cookie值>
   ```

3. **重启服务**
   ```bash
   docker compose restart bot
   ```

4. **测试**
   - 发送一个 YouTube 链接试试！

**仅管理员可以使用此命令。**

---

## 获取 Cookies（详细步骤）

### Chrome / Edge

1. **打开 YouTube 并登录**
   - 访问 https://www.youtube.com
   - 确保已登录

2. **打开开发者工具**
   - 按 `F12` 或 `Ctrl+Shift+I` (Mac: `Cmd+Option+I`)

3. **进入 Application 标签**
   - 点击顶部的 "Application" 或 "应用程序"

4. **找到 Cookies**
   - 左侧菜单展开 "Cookies"
   - 点击 "https://www.youtube.com"

5. **复制 Cookie 值**
   - 找到名为 `__Secure-3PSID` 或 `SID` 的 cookie
   - 双击对应的 "Value" 列，复制整个值

### Firefox

1. **打开 YouTube 并登录**

2. **打开开发者工具**
   - 按 `F12` 或 `Ctrl+Shift+K`

3. **进入 Storage 标签**
   - 点击顶部的 "Storage"

4. **找到 Cookies**
   - 左侧展开 "Cookies"
   - 点击 "https://www.youtube.com"

5. **复制 Cookie 值**
   - 找到 `__Secure-3PSID` 或 `SID`
   - 右键点击 → 复制值

---

## 手动配置方法

如果 Bot 命令不可用，可以手动配置：

### 1. 创建 cookies 文件

在服务器部署目录创建 `youtube-cookies.txt`：

```
# Netscape HTTP Cookie File
.youtube.com	TRUE	/	TRUE	0	__Secure-3PSID	<你的cookie值>
```

### 2. 配置 config.yaml

```yaml
download:
  worker_count: 3
  max_file_size: 50
  temp_dir: "./tmp"
  cookies_file: "/app/youtube-cookies.txt"  # 添加这行
```

### 3. 重启服务

```bash
docker compose down
docker compose up -d
```

---

## 如果找不到 `__Secure-3PSID`

**尝试其他 cookie 名称**，按优先级：

1. `__Secure-3PSID`（最重要）
2. `SID`
3. `HSID`
4. `SSID`
5. `APISID`

通常 `SID` 就可以工作。

---

## 验证配置

发送一个 YouTube 链接给 Bot，如果能正常下载说明配置成功！

---

## 常见问题

### Q: Cookies 过期了怎么办

**A:** YouTube cookies 通常有效期为数周到数月。过期后重新发送 `/cookies` 命令即可。

### Q: Bot 命令提示权限不足

**A:** `/cookies` 命令仅管理员可用。确保使用配置文件中设置的 Admin ID 账号。

### Q: 还是下载失败

**A:** 可能原因：
1. Cookie 值不完整（确保复制了整个 Value）
2. Cookie 已过期
3. 视频有地区限制

**替代方案**：使用在线工具转换后直接发送 MP3 文件给 Bot
- https://y2mate.com
- https://yt1s.com

---

## 文件格式说明

cookies 文件格式（Tab 分隔）：

```
域名    子域名    路径    HTTPS    过期时间    名称    值
```

示例：

```
.youtube.com	TRUE	/	TRUE	0	__Secure-3PSID	你的cookie值
```

注意：字段之间用 **Tab** 分隔，不是空格！
