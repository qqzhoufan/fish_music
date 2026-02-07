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

1. **获取 Cookie**（见下方详细教程）
2. **发送给 Bot**：`/cookies <复制的cookie值>`
3. **重启服务**：`docker compose restart bot`

---

## 📖 获取 Cookies 详细步骤

### 方法 1：Chrome / Edge（推荐）

#### 步骤 1：打开 YouTube 并登录
- 访问 https://www.youtube.com
- 确保已登录你的 Google 账号

#### 步骤 2：打开开发者工具
- **Windows/Linux**: 按 `F12` 或 `Ctrl + Shift + I`
- **Mac**: 按 `Cmd + Option + I`

#### 步骤 3：进入 Application 标签
- 点击顶部的 **"Application"** 或 **"应用程序"** 标签
- 如果看不到这个标签，点击 `>>` 更多按钮

#### 步骤 4：展开 Cookies
- 在左侧边栏找到 **"Storage"** 或 **"存储"**
- 展开 **"Cookies"**
- 点击 **"https://www.youtube.com"**

#### 步骤 5：查找 Cookie（按优先级尝试）

**优先尝试以下名称：**

1. **`__Secure-3PSID`** (最重要，通常以 `Cg` 开头)
2. **`SID`** (如果找不到上面这个)
3. **`HSID`**
4. **`SSID`**
5. **`APISID`**
6. **`SAPISID`**

#### 步骤 6：复制 Cookie 值
- 找到对应的 Cookie 名称
- 双击 **"Value"** 列的值（不是 Name）
- 右键 → 复制（或 Ctrl+C）

**⚠️ 重要：**
- 复制的是 **Value** 列，不是 Name 列
- 确保复制了完整的值（可能很长）
- 值通常以 `Cg`、`aSRj` 等字符开头

---

### 方法 2：Firefox

#### 步骤 1-2：同上（打开 YouTube 并登录，F12）

#### 步骤 3：进入 Storage 标签
- 点击顶部的 **"Storage"** 标签

#### 步骤 4-6：同 Chrome 方法

---

### 方法 3：如果找不到 youtube.com 的 cookies

尝试查看 **`https://google.com`** 的 cookies：

1. 在 Cookies 列表中找到 **"https://google.com"**
2. 查找以下 cookies：
   - `SAPISID`
   - `APISID`
   - `SID`
   - `HSID`

这些 Google cookies 也可以工作！

---

### 方法 4：使用浏览器扩展（最简单）

如果上面的方法太复杂，可以使用浏览器扩展：

#### Chrome / Edge 扩展
1. 安装 **"Get cookies.txt LOCALLY"** 扩展
   - Chrome: https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc
2. 访问 https://www.youtube.com 并登录
3. 点击浏览器工具栏的扩展图标
4. 选择 **"Current Site"** → **"Export"** → **"Download"**
5. 打开下载的文件，找到 `__Secure-3PSID` 或 `SID` 的值

#### Firefox 扩展
1. 安装 **"Get cookies.txt LOCALLY"** 扩展
   - Firefox: https://addons.mozilla.org/en-US/firefox/addon/get-cookiestxt-locally/
2. 后续步骤同上

---

## 🍪 使用 Bot 配置

获取到 cookie 值后：

1. **发送给 Bot**（仅管理员可用）：
   ```
   /cookies <你的cookie值>
   ```

2. **示例**：
   ```
   /cookies CgQihi...
   ```

3. **重启服务**：
   ```bash
   docker compose restart bot
   ```

---

## 🔍 常见问题

### Q: 我在 Application 标签下找不到 Cookies

**A:** 确认步骤：
1. 确保已经完全打开开发者工具（不是缩小版）
2. 在左侧最底部找到 "Storage" 或 "存储"
3. 展开 "Storage" 才能看到 "Cookies"
4. 如果还是没有，尝试刷新页面

### Q: 找到了 Cookies 列表，但没有 `__Secure-3PSID`

**A:** 尝试以下 cookie（按优先级）：
1. `SID`
2. `HSID`
3. `SSID`
4. `SAPISID`
5. `APISID`

通常 `SID` 就可以工作！

### Q: Cookie 值太短，看起来不对

**A:** 正确的 cookie 值通常：
- `__Secure-3PSID`: 很长（100+ 字符），以 `Cg` 开头
- `SID`: 中等长度（50-100 字符）
- `HSID`: 较短（20-50 字符）

如果值只有几个字符，可能复制错了。确保复制的是 **Value** 列！

### Q: 我没有登录 YouTube

**A:** 必须先登录！
1. 点击 YouTube 右上角的登录按钮
2. 使用你的 Google 账号登录
3. 登录后再查看 cookies

### Q: 配置后还是下载失败

**A:** 可能原因：
1. Cookie 值不完整 → 重新复制
2. Cookie 已过期 → 重新获取
3. 配置了错误的 cookie → 尝试其他 cookie 名称
4. 没有重启服务 → 运行 `docker compose restart bot`

---

## 🎯 图片教程（文字版）

### Chrome / Edge:

```
┌─────────────────────────────────────┐
│  YouTube 已登录                      │
│                                     │
│  [按 F12]                           │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│ 开发者工具                          │
│                                     │
│ Elements | Console | Sources | ...  │
│ Application >>                      │  ← 点击这个
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│ Application                         │
│                                     │
│ ▼ Storage                           │
│   ▼ Cookies                         │  ← 展开这个
│     ▼ https://www.youtube.com       │  ← 点击这个
│                                     │
│ Name    | Value  | Domain  | ...    │
│---------|--------|--------|--------│
│ ...     | ...    | ...    | ...    │
│ __Secure-3PSID | CgQihi...| .youtube│  ← 找这个！
│ SID     | aSRj...| .youtube|        │  ← 或这个
└─────────────────────────────────────┘
```

---

## 📱 手动配置方法（如果 Bot 不可用）

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
  cookies_file: "/app/youtube-cookies.txt"
```

### 3. 重启服务

```bash
docker compose down
docker compose up -d
```

---

## 🔒 安全提示

- Cookie 包含你的登录凭证，不要分享给他人
- 建议使用专用的 Google 账号
- Cookie 通常有效期为数周到数月
- 过期后重新获取即可

---

## ⚡ 快速参考

| Cookie 名称 | 位置 | 长度 | 优先级 |
|------------|------|------|--------|
| `__Secure-3PSID` | youtube.com | 很长 (100+) | ⭐⭐⭐ 最高 |
| `SID` | youtube.com | 中等 (50-100) | ⭐⭐ 高 |
| `HSID` | youtube.com | 较短 (20-50) | ⭐ 中 |
| `SAPISID` | google.com | 中等 | ⭐ 中 |
| `APISID` | google.com | 中等 | ⭐ 低 |

**提示**：如果找不到 `__Secure-3PSID`，使用 `SID` 通常也可以！
