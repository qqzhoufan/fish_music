package handler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/fish-music/internal/config"
	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/model"
	"github.com/user/fish-music/internal/service"
	"github.com/user/fish-music/pkg/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotHandler Bot 处理器
type BotHandler struct {
	bot            *tgbotapi.BotAPI
	adminID        int64
	songRepo       *database.SongRepository
	userRepo       *database.UserRepository
	favoriteRepo   *database.FavoriteRepository
	historyRepo    *database.HistoryRepository
	musicAPI       *api.NeteaseAPI
	ytdlpService   *service.YTDLPService
	downloadConfig *config.DownloadConfig
}

// NewBotHandler 创建 Bot 处理器
func NewBotHandler(
	bot *tgbotapi.BotAPI,
	adminID int64,
	songRepo *database.SongRepository,
	userRepo *database.UserRepository,
	favoriteRepo *database.FavoriteRepository,
	historyRepo *database.HistoryRepository,
	musicAPI *api.NeteaseAPI,
	ytdlpService *service.YTDLPService,
	downloadConfig *config.DownloadConfig,
) *BotHandler {
	return &BotHandler{
		bot:            bot,
		adminID:        adminID,
		songRepo:       songRepo,
		userRepo:       userRepo,
		favoriteRepo:   favoriteRepo,
		historyRepo:    historyRepo,
		musicAPI:       musicAPI,
		ytdlpService:   ytdlpService,
		downloadConfig: downloadConfig,
	}
}

// HandlePrivateMessage 处理私聊消息
func (h *BotHandler) HandlePrivateMessage(update tgbotapi.Update) error {
	message := update.Message
	if message == nil {
		return nil
	}

	// 获取或创建用户
	user, err := h.userRepo.FindOrCreate(
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 处理命令
	if message.IsCommand() {
		return h.handleCommand(message, user)
	}

	// 处理搜索关键词
	return h.handleSearch(message, user)
}

// handleCommand 处理命令
func (h *BotHandler) handleCommand(message *tgbotapi.Message, user *model.User) error {
	switch message.Command() {
	case "start":
		return h.cmdStart(message, user)
	case "help":
		return h.cmdHelp(message, user)
	case "history":
		return h.cmdHistory(message, user)
	case "favorites", "favs":
		return h.cmdFavorites(message, user)
	case "random":
		return h.cmdRandom(message, user)
	case "songs", "list":
		return h.cmdSongs(message, user)
	case "stats":
		return h.cmdStats(message, user)
	case "add":
		return h.cmdAdd(message, user)
	case "cookies":
		return h.cmdCookies(message, user)
	default:
		return h.cmdUnknown(message, user)
	}
}

// cmdAdd 添加音乐命令
func (h *BotHandler) cmdAdd(message *tgbotapi.Message, user *model.User) error {
	text := `📥 <b>如何添加音乐到 Fish Music</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>⭐ 方法一：YouTube 自动下载（最推荐）</b>

只需发送 YouTube 链接，自动下载并保存！

<b>支持的链接格式：</b>
• https://www.youtube.com/watch?v=xxxxx
• https://youtu.be/xxxxx

<b>使用步骤：</b>
1. 📺 在 YouTube 找到音乐视频
2. 📋 复制视频链接
3. 💬 直接发送给机器人
4. ⏳ 等待 1-3 分钟自动下载
5. ✅ 下载完成，自动保存！

<b>提示：</b>
• 可以下载任何 YouTube 音乐视频
• 自动提取音频为 MP3 格式
• 自动识别歌手和歌曲信息
• 单个文件最大 50MB

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>⭐⭐ 方法二：直接发送 MP3 文件</b>

<b>100% 成功率，最可靠的方式！</b>

<b>使用步骤：</b>
1. 📱 在 Telegram 点击发送文件
2. 🎵 选择 MP3 音频文件
3. 💬 发送给机器人
4. ✅ 立即保存成功！

<b>获取 MP3 的方法：</b>
• 使用在线 YouTube 转 MP3 工具
• 从电脑已有的音乐库选择
• 从其他音乐平台下载后发送

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>⭐⭐⭐ 方法三：手动添加 File ID</b>

如果你有 Telegram 文件的 File ID，可以手动添加。

<b>命令格式：</b>
<code>/add [歌曲名] [歌手名] [File_ID]</code>

<b>示例：</b>
<code>/add 稻香 周杰伦 AwADBwADgAD...</code>

<b>如何获取 File ID：</b>
1. 向机器人 @GetPublicIdBot 发送音频文件
2. 机器人会返回 File ID
3. 复制 File ID 使用上面的命令添加

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>💡 推荐使用方案</b>

<b>最佳方案：YouTube 自动下载</b>
• ✅ 全自动，最方便
• ✅ 自动识别歌曲信息
• ⚠️ 需要等待 1-3 分钟
• ⚠️ 部分 YouTube 视频可能下载失败

<b>最稳方案：发送 MP3 文件</b>
• ✅ 100% 成功率
• ✅ 秒速保存
• ✅ 不受平台限制
• ⚠️ 需要先获取 MP3 文件

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>❓ 常见问题</b>

Q: YouTube 下载失败怎么办？
A: 建议使用在线工具转换为 MP3 后发送给我

Q: 可以下载其他平台的视频吗？
A: 目前主要支持 YouTube，其他平台可能不稳定

Q: 下载需要多久？
A: 通常 1-3 分钟，取决于视频大小和网络速度

Q: 有文件大小限制吗？
A: 单个文件最大 50MB

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>🎉 现在就开始添加音乐吧！</b>

直接发送 YouTube 链接试试吧！ 🎵`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// cmdStart 开始命令
func (h *BotHandler) cmdStart(message *tgbotapi.Message, user *model.User) error {
	text := `🎵 <b>欢迎来到 Fish Music</b>

你的个人云端音乐库，基于 Telegram 无限存储空间！

<b>🚀 快速开始</b>
• 发送歌曲名或歌手名搜索音乐
• 发送 YouTube 链接自动下载音乐 ⭐
• 直接发送 MP3 文件保存

<b>📱 主要功能</b>
• <b>/songs</b> - 浏览音乐库 ⭐ 新功能
• <b>/random</b> - 随机播放一首歌
• <b>/favorites</b> - 我的收藏列表
• <b>/history</b> - 播放历史记录
• <b>/stats</b> - 音乐库统计
• <b>/add</b> - 添加音乐教程
• <b>/cookies</b> - 配置 YouTube 下载 ⭐ 新功能

<b>🌟 特色功能</b>
✅ YouTube 自动下载 - 发链接即可
✅ 元数据自动识别 - 歌手/地区/年份
✅ 收藏和历史 - 永久记录
✅ 无限存储 - 基于 Telegram 云端
✅ 歌曲分类 - 类型/语言筛选

<b>❓ YouTube 下载失败？</b>
发送 /cookies 查看配置教程

💡 <b>小技巧</b>
在任何群组中输入 @BotName 关键词 也能搜索！

需要帮助？使用 /help 查看完整指南`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// cmdHelp 帮助命令
func (h *BotHandler) cmdHelp(message *tgbotapi.Message, user *model.User) error {
	text := `📖 <b>Fish Music 完全使用指南</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>🎵 搜索与播放音乐</b>

<b>方式一：搜索歌曲</b>
直接发送歌曲名或歌手名
例如：<code>周杰伦 稻香</code>

<b>方式二：群组内搜索</b>
在任何群组输入：<code>@BotName 歌曲名</code>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>📥 添加音乐的三种方式</b>

<b>⭐ 方式一：YouTube 自动下载（推荐）</b>
1. 在 YouTube 找到音乐视频
2. 复制链接发送给我
3. 自动下载并保存到库中！

支持的链接格式：
• https://www.youtube.com/watch?v=xxx
• https://youtu.be/xxx

<b>⭐⭐ 方式二：直接发送 MP3 文件</b>
1. 在 Telegram 选择发送文件
2. 选择 MP3 音频文件
3. 发送给我即可保存

<b>⭐⭐⭐ 方式三：手动添加 File ID</b>
使用 <code>/add [歌名] [歌手] [File_ID]</code>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>📱 所有命令列表</b>

<b>/start</b> - 查看欢迎信息
<b>/help</b> - 显示本帮助信息
<b>/songs</b> 或 <b>/list</b> - 浏览音乐库 ⭐ 新功能
<b>/random</b> - 随机播放一首歌
<b>/favorites</b> 或 <b>/favs</b> - 收藏列表
<b>/history</b> - 播放历史（最近20首）
<b>/stats</b> - 音乐库统计数据
<b>/add</b> - 添加音乐详细教程
<b>/cookies</b> - 配置 YouTube 下载 ⭐ 新功能

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>✨ 功能亮点</b>

❤️ <b>收藏功能</b>
播放歌曲时点击 ❤️ 按钮即可收藏
随时查看收藏列表，不会丢失

📜 <b>历史记录</b>
自动记录所有播放过的歌曲
支持查看最近 20 首播放记录

🎲 <b>随机播放</b>
不知道听什么？试试随机播放
发现音乐库中的惊喜

🌍 <b>智能元数据</b>
• 自动识别歌手地区（国家 Emoji）
• 显示发行年份
• 完整的歌曲信息展示

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>💡 使用技巧</b>

1. <b>批量添加</b>：可以连续发送多个链接
2. <b>收藏整理</b>：喜欢的歌及时收藏
3. <b>搜索技巧</b>：歌名+歌手搜索更准确
4. <b>群组分享</b>：在任何群组都能搜索播放

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>❓ 常见问题</b>

Q: YouTube 下载失败怎么办？
A: 如果显示 "Sign in to confirm you're not a bot" 错误：
   1. 发送 <code>/cookies</code> 查看配置教程
   2. 按提示配置 cookies 即可解决
   3. 配置后需管理员重启服务

Q: 下载 YouTube 需要多久？
A: 通常 1-3 分钟，取决于视频大小

Q: 文件大小限制？
A: 单个文件最大 50MB

Q: 音乐会占用手机空间吗？
A: 不会！存储在 Telegram 云端

Q: 可以在电脑上用吗？
A: 可以！Telegram 桌面版同样支持

━━━━━━━━━━━━━━━━━━━━━━━━━

如有问题或建议，请联系管理员 🎵`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// cmdHistory 历史记录命令
func (h *BotHandler) cmdHistory(message *tgbotapi.Message, user *model.User) error {
	songs, err := h.historyRepo.GetByUser(user.ID, 20)
	if err != nil {
		return err
	}

	if len(songs) == 0 {
		text := `📜 <b>暂无播放历史</b>

你还没有播放过任何歌曲哦～

<b>💡 快速开始：</b>
• 搜索歌曲：直接发送歌名
• 随机播放：使用 <code>/random</code>
• 添加音乐：发送 YouTube 链接

开始播放后，这里会自动记录你的播放历史！`
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "HTML"
		_, err := h.bot.Send(msg)
		return err
	}

	var text strings.Builder
	text.WriteString("📜 <b>最近播放</b>\n\n")
	text.WriteString(fmt.Sprintf("显示最近 %d 首播放记录\n\n", len(songs)))

	for i, song := range songs {
		emoji := song.GetCountryEmoji()
		year := song.GetYearText()
		text.WriteString(fmt.Sprintf("%d. %s <b>%s</b> - %s (%s)\n", i+1, emoji, song.Title, song.Artist, year))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text.String())
	msg.ParseMode = "HTML"
	_, err = h.bot.Send(msg)
	return err
}

// cmdFavorites 收藏列表命令
func (h *BotHandler) cmdFavorites(message *tgbotapi.Message, user *model.User) error {
	songs, err := h.favoriteRepo.GetByUser(user.ID, 50)
	if err != nil {
		return err
	}

	if len(songs) == 0 {
		text := `⭐ <b>暂无收藏歌曲</b>

你还没有收藏任何歌曲哦～

<b>💡 如何收藏歌曲：</b>
播放任何歌曲时，点击播放卡片上的 <b>❤️ 收藏</b> 按钮即可！

收藏后的歌曲会永久保存在这里，随时可以查看和播放。

<b>🎵 现在就去搜索喜欢的歌曲吧！</b>`
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "HTML"
		_, err := h.bot.Send(msg)
		return err
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("⭐ <b>我的收藏</b> (共 %d 首)\n\n", len(songs)))

	for i, song := range songs {
		emoji := song.GetCountryEmoji()
		year := song.GetYearText()
		text.WriteString(fmt.Sprintf("%d. %s <b>%s</b> - %s (%s)\n", i+1, emoji, song.Title, song.Artist, year))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text.String())
	msg.ParseMode = "HTML"
	_, err = h.bot.Send(msg)
	return err
}

// cmdRandom 随机播放命令
func (h *BotHandler) cmdRandom(message *tgbotapi.Message, user *model.User) error {
	song, err := h.songRepo.GetRandom()
	if err != nil {
		text := `🎲 <b>音乐库暂无歌曲</b>

音乐库还是空的，添加一些歌曲吧！

<b>📥 添加音乐的方法：</b>

<b>⭐ YouTube 自动下载（推荐）</b>
发送 YouTube 链接，自动下载音乐
例如：https://www.youtube.com/watch?v=xxxxx

<b>⭐⭐ 发送 MP3 文件</b>
直接发送 MP3 文件，秒速保存！

<b>💡 使用教程：</b>
发送 <code>/add</code> 查看详细添加教程

🎵 开始添加你的第一首歌吧！`
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "HTML"
		_, err := h.bot.Send(msg)
		return err
	}

	return h.sendSong(message.Chat.ID, song, user)
}

// cmdStats 统计信息命令
func (h *BotHandler) cmdStats(message *tgbotapi.Message, user *model.User) error {
	stats, err := h.songRepo.GetStats()
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`📊 <b>音乐库统计信息</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

🎵 <b>总歌曲数</b>
   %v 首

🎤 <b>歌手数量</b>
   %v 位

❌ <b>缺失歌曲</b>
   %v 首

📅 <b>今日新增</b>
   %v 首

━━━━━━━━━━━━━━━━━━━━━━━━━

💡 <b>提示：</b>
缺失的歌曲需要重新补档
请使用管理后台处理`,
		stats["total_songs"],
		stats["total_artists"],
		stats["missing_songs"],
		stats["today_added"],
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err = h.bot.Send(msg)
	return err
}

// cmdUnknown 未知命令
func (h *BotHandler) cmdUnknown(message *tgbotapi.Message, user *model.User) error {
	text := `❓ <b>未知命令</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

我不认识这个命令哦～

<b>📱 可用命令列表：</b>

/start - 查看欢迎信息
/help - 完整使用指南
/random - 随机播放
/favorites - 收藏列表
/history - 播放历史
/stats - 统计信息
/add - 添加音乐教程

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>💡 或直接发送：</b>
• 歌曲名或歌手名搜索
• YouTube 链接自动下载
• MP3 文件直接保存

使用 <code>/help</code> 查看完整帮助 🎵`
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// handleSearch 处理搜索
func (h *BotHandler) handleSearch(message *tgbotapi.Message, user *model.User) error {
	keyword := strings.TrimSpace(message.Text)
	if keyword == "" {
		return nil
	}

	// 检测是否是 URL
	if strings.HasPrefix(keyword, "http://") || strings.HasPrefix(keyword, "https://") {
		return h.handleURL(message, keyword)
	}

	// 从数据库搜索
	songs, err := h.songRepo.Search(keyword, 10)
	if err != nil {
		return err
	}

	// 如果数据库有结果，直接返回
	if len(songs) > 0 {
		return h.sendSearchResults(message.Chat.ID, songs, keyword)
	}

	// 数据库无结果，提示用户如何添加
	text := fmt.Sprintf(`🔍 <b>未找到相关歌曲</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

关键词：<b>%s</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>💡 快速添加音乐的方法：</b>

<b>方法一：YouTube 自动下载 ⭐</b>
直接发送 YouTube 链接，自动下载音乐！

例如：
• https://www.youtube.com/watch?v=xxxxx
• https://youtu.be/xxxxx

<b>方法二：发送 MP3 文件 ⭐⭐⭐</b>
最可靠的方式，100% 成功！
直接在 Telegram 选择文件发送即可

<b>方法三：查看添加教程</b>
使用 <code>/add</code> 命令查看详细教程

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>🎵 现在就试试吧！</b>

发送一个 YouTube 链接，或者 MP3 文件～`, keyword)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err = h.bot.Send(msg)
	return err
}

// handleURL 处理音乐链接
func (h *BotHandler) handleURL(message *tgbotapi.Message, musicURL string) error {
	// 获取用户
	user, err := h.userRepo.FindByTelegramID(message.From.ID)
	if err != nil {
		return err
	}

	// 检测是否是支持的视频平台
	if h.isSupportedVideoPlatform(musicURL) {
		// 使用 yt-dlp 下载
		return h.ytdlpService.DownloadAndSave(message.Chat.ID, musicURL, user)
	}

	// 其他平台，提示用户
	return h.handleUnsupportedPlatform(message, musicURL)
}

// isSupportedVideoPlatform 检查是否支持的视频平台
func (h *BotHandler) isSupportedVideoPlatform(url string) bool {
	supportedPlatforms := []string{
		"youtube.com",
		"youtu.be",
		"bilibili.com",
		"b23.tv",
		"music.163.com",
		"y.qq.com",
		"kugou.com",
		"kuwo.cn",
	}

	urlLower := strings.ToLower(url)
	for _, platform := range supportedPlatforms {
		if strings.Contains(urlLower, platform) {
			return true
		}
	}

	return false
}

// handleUnsupportedPlatform 处理不支持的平台
func (h *BotHandler) handleUnsupportedPlatform(message *tgbotapi.Message, musicURL string) error {
	text := fmt.Sprintf(`📋 <b>收到链接</b>

%s

<b>💡 支持的平台：</b>

🎬 <b>视频平台</b>
• YouTube: youtube.com
• Bilibili: bilibili.com
• 其他 yt-dlp 支持的平台

🎵 <b>音乐平台</b>
• 网易云音乐
• QQ音乐、酷狗等

<b>✅ 推荐方法：</b>

1. <b>YouTube/B站</b>
   直接发送视频链接，我会自动提取音频！

2. <b>网易云等</b>
   • 发送链接获取歌曲信息
   • 然后手动下载 MP3 发给我

3. <b>直接发送 MP3</b>
   最简单可靠的方式！

---
💡 提示：支持的平台会自动下载并添加到音乐库`, musicURL)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// sendSearchResults 发送搜索结果（带分页）
func (h *BotHandler) sendSearchResults(chatID int64, songs []*model.Song, keyword string) error {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("🔍 <b>搜索结果</b>：%s\n\n", keyword))

	for i, song := range songs {
		emoji := song.GetCountryEmoji()
		year := song.GetYearText()
		text.WriteString(fmt.Sprintf("%d. %s <b>%s</b> - %s (%s)\n", i+1, emoji, song.Title, song.Artist, year))
	}

	// 创建 Inline Keyboard
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := []tgbotapi.InlineKeyboardButton{}

	for i, song := range songs {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s - %s", i+1, truncateString(song.Artist, 15), truncateString(song.Title, 20)),
			fmt.Sprintf("play_%d", song.ID),
		)
		row = append(row, btn)

		// 每行最多 2 个按钮
		if len(row) == 2 || i == len(songs)-1 {
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	_, err := h.bot.Send(msg)
	return err
}

// sendSong 发送歌曲
func (h *BotHandler) sendSong(chatID int64, song *model.Song, user *model.User) error {
	// 构建音频文件 - 使用 FileID 类型包装字符串
	audio := tgbotapi.NewAudio(chatID, tgbotapi.FileID(song.FileID))
	audio.Title = song.Title
	audio.Performer = song.Artist
	if song.Album != "" {
		audio.Caption = fmt.Sprintf("🎵 %s - %s\n%s %s", song.Artist, song.Title, song.GetCountryEmoji(), song.GetYearText())
	}

	// 检查是否已收藏
	isFavorited, _ := h.favoriteRepo.IsFavorited(user.ID, song.ID)

	// 创建操作按钮
	var keyboard [][]tgbotapi.InlineKeyboardButton

	favoriteBtn := tgbotapi.NewInlineKeyboardButtonData("❤️ 收藏", fmt.Sprintf("fav_%d", song.ID))
	if isFavorited {
		favoriteBtn = tgbotapi.NewInlineKeyboardButtonData("💔 取消收藏", fmt.Sprintf("unfav_%d", song.ID))
	}

	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{favoriteBtn})
	audio.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	// 发送音频
	_, err := h.bot.Send(audio)
	if err != nil {
		// 如果 FileID 失效，标记为需要补档
		if strings.Contains(err.Error(), "file") || strings.Contains(err.Error(), "invalid") {
			h.songRepo.MarkMissing(song.ID)
		}
		return err
	}

	// 记录历史
	h.historyRepo.Add(user.ID, song.ID)
	h.userRepo.UpdateLastSeen(user.ID)

	return nil
}

// HandleCallback 处理回调查询
func (h *BotHandler) HandleCallback(query *tgbotapi.CallbackQuery) error {
	// 获取用户
	user, err := h.userRepo.FindByTelegramID(query.From.ID)
	if err != nil {
		return h.answerCallback(query, "❌ 用户不存在", true)
	}

	data := query.Data

	// 解析回调数据
	if strings.HasPrefix(data, "play_") {
		return h.callbackPlay(query, user)
	}

	if strings.HasPrefix(data, "fav_") {
		return h.callbackFavorite(query, user, true)
	}

	if strings.HasPrefix(data, "unfav_") {
		return h.callbackFavorite(query, user, false)
	}

	return h.answerCallback(query, "❌ 未知操作", true)
}

// callbackPlay 播放回调
func (h *BotHandler) callbackPlay(query *tgbotapi.CallbackQuery, user *model.User) error {
	songIDStr := strings.TrimPrefix(query.Data, "play_")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		return h.answerCallback(query, "❌ 无效的歌曲ID", true)
	}

	song, err := h.getSongByID(uint(songID))
	if err != nil {
		return h.answerCallback(query, "❌ 歌曲不存在", true)
	}

	// 发送歌曲到聊天
	if err := h.sendSong(query.Message.Chat.ID, song, user); err != nil {
		return h.answerCallback(query, "❌ 发送失败", true)
	}

	return h.answerCallback(query, "✅ 播放成功", false)
}

// callbackFavorite 收藏回调
func (h *BotHandler) callbackFavorite(query *tgbotapi.CallbackQuery, user *model.User, add bool) error {
	songIDStr := strings.TrimPrefix(strings.TrimPrefix(query.Data, "fav_"), "unfav_")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		return h.answerCallback(query, "❌ 无效的歌曲ID", true)
	}

	if add {
		if err := h.favoriteRepo.Add(user.ID, uint(songID)); err != nil {
			return h.answerCallback(query, "❌ 收藏失败", true)
		}
		return h.answerCallback(query, "❤️ 已收藏", false)
	} else {
		if err := h.favoriteRepo.Remove(user.ID, uint(songID)); err != nil {
			return h.answerCallback(query, "❌ 取消收藏失败", true)
		}
		return h.answerCallback(query, "💔 已取消收藏", false)
	}
}

// answerCallback 回答回调查询
func (h *BotHandler) answerCallback(query *tgbotapi.CallbackQuery, text string, alert bool) error {
	callback := tgbotapi.NewCallback(query.ID, text)
	callback.ShowAlert = alert
	_, err := h.bot.Request(callback)
	return err
}

// getSongByID 根据 ID 获取歌曲
func (h *BotHandler) getSongByID(id uint) (*model.Song, error) {
	var song model.Song
	err := database.DB.Where("id = ?", id).First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	// 使用 rune 来正确处理 UTF-8 多字节字符
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// cmdSongs 显示歌曲列表命令
func (h *BotHandler) cmdSongs(message *tgbotapi.Message, user *model.User) error {
	// 随机获取最多10首歌曲
	songs, err := h.songRepo.GetRandomSongs(10)
	if err != nil {
		text := "❌ 获取歌曲列表失败"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		_, err := h.bot.Send(msg)
		return err
	}

	if len(songs) == 0 {
		text := `🎵 <b>音乐库是空的</b>

还没有任何歌曲哦～

<b>💡 快速添加音乐：</b>
• 发送 YouTube 链接自动下载
• 直接发送 MP3 文件

使用 <code>/add</code> 查看详细教程 🎵`
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "HTML"
		_, err := h.bot.Send(msg)
		return err
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("🎵 <b>音乐库歌曲列表</b>\n\n"))
	text.WriteString(fmt.Sprintf("随机展示 %d 首歌曲\n\n", len(songs)))

	for i, song := range songs {
		emoji := song.GetCountryEmoji()
		year := song.GetYearText()
		genre := song.GetGenreText()
		language := song.GetLanguageText()

		text.WriteString(fmt.Sprintf("<b>%d.</b> %s <b>%s</b> - %s\n", i+1, emoji, song.Title, song.Artist))
		text.WriteString(fmt.Sprintf("   %s · %s · %s\n", language, genre, year))
		text.WriteString("\n")
	}

	text.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	text.WriteString("💡 点击下方按钮播放歌曲")

	msg := tgbotapi.NewMessage(message.Chat.ID, text.String())
	msg.ParseMode = "HTML"

	// 创建 Inline Keyboard
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := []tgbotapi.InlineKeyboardButton{}

	for i, song := range songs {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", i+1, truncateString(song.Title, 20)),
			fmt.Sprintf("play_%d", song.ID),
		)
		row = append(row, btn)

		// 每行最多 2 个按钮
		if len(row) == 2 || i == len(songs)-1 {
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	_, err = h.bot.Send(msg)
	return err
}

// cmdCookies 配置 YouTube Cookies 命令（仅管理员）
func (h *BotHandler) cmdCookies(message *tgbotapi.Message, user *model.User) error {
	// 只有管理员可以使用此命令
	if message.From.ID != h.adminID {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ 此命令仅管理员可用")
		msg.ParseMode = "HTML"
		h.bot.Send(msg)
		return nil
	}

	// 获取命令参数
	args := message.CommandArguments()

	// 如果没有参数，发送使用说明
	if args == "" {
		text := `🍪 <b>YouTube Cookies 配置</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

<b>使用方法：</b>

1. 获取 Cookie：
   • 打开 https://www.youtube.com 并登录
   • 按 F12 → Application → Cookies
   • 找到 <code>__Secure-3PSID</code> 或 <code>SID</code>
   • 复制 Value 值

2. 发送给 Bot：
   <code>/cookies 你的cookie值</code>

<b>示例：</b>
<code>/cookies CgQihi...</code>

━━━━━━━━━━━━━━━━━━━━━━━━━

💡 配置后需要重启服务才能生效
   <code>docker compose restart bot</code>`

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "HTML"
		_, err := h.bot.Send(msg)
		return err
	}

	// 保存 cookies 到文件
	cookieValue := strings.TrimSpace(args)

	// 验证 cookie 不为空
	if cookieValue == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Cookie 值不能为空")
		h.bot.Send(msg)
		return nil
	}

	// 验证 cookie 格式（基本检查）
	if len(cookieValue) < 50 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Cookie 值格式不正确（太短）\n\n请确保复制了完整的 Value 值")
		h.bot.Send(msg)
		return nil
	}

	// 写入 cookies 文件
	cookiesContent := fmt.Sprintf("# Netscape HTTP Cookie File\n# Auto-generated by Fish Music Bot\n\n.youtube.com\tTRUE\t/\tTRUE\t0\t__Secure-3PSID\t%s\n", cookieValue)

	err := os.WriteFile("/app/youtube-cookies.txt", []byte(cookiesContent), 0644)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("❌ 保存失败：%v", err))
		h.bot.Send(msg)
		return err
	}

	// 发送成功消息
	text := fmt.Sprintf(`✅ <b>Cookie 配置成功！</b>

━━━━━━━━━━━━━━━━━━━━━━━━━

Cookie 已保存到服务器。

<b>下一步：</b>
重启 Bot 服务使配置生效：

<code>docker compose restart bot</code>

━━━━━━━━━━━━━━━━━━━━━━━━━

💡 测试：发送一个 YouTube 链接试试`)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	_, err = h.bot.Send(msg)
	return err
}
