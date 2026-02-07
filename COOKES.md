# YouTube Cookies 配置指南

## 为什么需要 Cookies？

YouTube 会检测来自服务器的下载请求，可能会显示以下错误：

```
ERROR: [youtube] XXXXX: Sign in to confirm you're not a bot.
```

配置 cookies 可以解决这个问题。

---

## 获取 Cookies（手动方法）

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

6. **创建 cookies 文件**
   - 在电脑上创建一个新文件 `youtube-cookies.txt`
   - 复制以下内容，替换 `<YOUR_COOKIE_VALUE>` 为刚才复制的值：

```
# Netscape HTTP Cookie File
.youtube.com	TRUE	/	TRUE	0	__Secure-3PSID	<YOUR_COOKIE_VALUE>
```

---

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

6. **创建 cookies 文件**（同上）

---

## 如果找不到上面的 Cookie

**尝试其他 cookie 名称**，按优先级：

1. `__Secure-3PSID`（最重要）
2. `SID`
3. `HSID`
4. `SSID`
5. `APISID`
6. `SAPISID`

可以添加多行，格式相同：

```
# Netscape HTTP Cookie File
.youtube.com	TRUE	/	TRUE	0	__Secure-3PSID	<COOKIE_1>
.youtube.com	TRUE	/	TRUE	0	SID	<COOKIE_2>
.youtube.com	TRUE	/	TRUE	0	HSID	<COOKIE_3>
```

---

## 配置步骤

### 1. 上传 cookies 文件

将 `youtube-cookies.txt` 上传到服务器的部署目录（与 `docker-compose.yml` 同级）：

```bash
# 使用 scp 上传
scp youtube-cookies.txt root@your-server:/path/to/fish-music/

# 或使用 FTP/SFTP 工具
```

### 2. 配置 config.yaml

编辑 `config.yaml`：

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

## 验证配置

发送一个 YouTube 链接给 Bot，如果能正常下载说明配置成功！

如果还有问题，查看日志：

```bash
docker compose logs -f bot
```

---

## 常见问题

### Q: 找不到 `__Secure-3PSID` 这个 cookie

**A:** 尝试使用 `SID`，或者尝试以下几个：
- `HSID`
- `SSID`
- `APISID`

通常 `SID` 就可以工作。

### Q: Cookies 过期了怎么办

**A:** YouTube cookies 通常有效期为数周到数月。过期后重新获取即可。

### Q: 还是不行

**A:** 替代方案：使用在线工具转换后直接发送 MP3 文件给 Bot
- https://y2mate.com
- https://yt1s.com

---

## 简化版（最快方法）

如果觉得上面太复杂，最简单的方法：

1. **安装 Chrome 扩展 "Get cookies.txt"**（虽然你不想用插件，但这是最简单的）
   - Chrome: https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc
   - 访问 YouTube
   - 点击扩展 → Export → Download
   - 得到 `youtube-cookies.txt`

2. **或者直接用在线转换工具**
   - 不配置 cookies
   - 用 https://y2mate.com 转 MP3
   - 直接发 MP3 给 Bot

---

## 文件格式说明

cookies 文件格式（Tab 分隔）：

```
域名    是否子域名    路径    是否HTTPS    过期时间    名称    值
```

示例：

```
.youtube.com	TRUE	/	TRUE	0	__Secure-3PSID	你的cookie值
```

注意：字段之间用 **Tab** 分隔，不是空格！
