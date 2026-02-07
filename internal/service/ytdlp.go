package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// YTDLPService yt-dlp 下载服务
type YTDLPService struct {
	bot        *tgbotapi.BotAPI
	songRepo   *database.SongRepository
	tempDir    string
	maxSize    int64
	cookiesFile string // YouTube cookies 文件路径（可选）
}

// NewYTDLPService 创建下载服务
func NewYTDLPService(
	bot *tgbotapi.BotAPI,
	songRepo *database.SongRepository,
	tempDir string,
	maxSize int,
	cookiesFile string,
) *YTDLPService {
	return &YTDLPService{
		bot:      bot,
		songRepo: songRepo,
		tempDir:  tempDir,
		maxSize:  int64(maxSize) * 1024 * 1024,
		cookiesFile: cookiesFile,
	}
}

// DownloadAndSave 下载并保存音乐
func (s *YTDLPService) DownloadAndSave(chatID int64, videoURL string, user *model.User) error {
	// 发送开始下载消息
	statusMsg := tgbotapi.NewMessage(chatID, "⏳ 开始下载...\n\n这可能需要几分钟，请稍候...")
	status, _ := s.bot.Send(statusMsg)

	// 检查是否已存在
	uniqueHash := s.generateHash(videoURL)
	existingSong, err := s.songRepo.FindByUniqueHash(uniqueHash)
	if err == nil && existingSong != nil {
		s.bot.Request(tgbotapi.NewDeleteMessage(chatID, status.MessageID))
		return s.sendDownloadedSong(chatID, existingSong, user)
	}

	// 下载音频
	tempFile, songInfo, err := s.downloadWithYTDLP(videoURL)
	if err != nil {
		s.bot.Request(tgbotapi.NewDeleteMessage(chatID, status.MessageID))
		// 发送错误消息给用户
		errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ 下载失败\n\n%s", err.Error()))
		errorMsg.ParseMode = "HTML"
		s.bot.Send(errorMsg)
		return err
	}
	defer os.Remove(tempFile)

	// 上传到 Telegram
	fileID, fileSize, err := s.uploadToTelegram(chatID, tempFile, songInfo)
	if err != nil {
		s.bot.Request(tgbotapi.NewDeleteMessage(chatID, status.MessageID))
		return fmt.Errorf("上传失败: %w", err)
	}

	// 保存到数据库
	song := &model.Song{
		UniqueHash:  uniqueHash,
		FileID:      fileID,
		SourceURL:   videoURL,
		Title:       songInfo.Title,
		Artist:      songInfo.Artist,
		Duration:    songInfo.Duration,
		FileSize:    fileSize,
		// CountryCode 不再根据歌手名自动判断，而是在 Web 后台编辑语言时自动设置
		Status:      "active",
	}

	if err := s.songRepo.Create(song); err != nil {
		return fmt.Errorf("保存失败: %w", err)
	}

	// 删除进度消息
	s.bot.Request(tgbotapi.NewDeleteMessage(chatID, status.MessageID))

	// 发送歌曲
	return s.sendDownloadedSong(chatID, song, user)
}

// SongInfo 歌曲信息
type SongInfo struct {
	Title    string
	Artist   string
	Album    string
	Duration int
}

// downloadWithYTDLP 使用 yt-dlp 下载
func (s *YTDLPService) downloadWithYTDLP(videoURL string) (string, *SongInfo, error) {
	// 生成唯一的文件名（不含扩展名）
	filename := fmt.Sprintf("%d_music", time.Now().UnixNano())
	tempBase := filepath.Join(s.tempDir, filename)
	tempFile := tempBase + ".mp3"

	// 第一步：获取标题
	titleArgs := []string{
		"--print", "title",
		"--no-playlist",
		"--no-warnings",
	}
	// 如果提供了 cookies 文件，添加到参数中
	if s.cookiesFile != "" {
		titleArgs = append([]string{"--cookies", s.cookiesFile}, titleArgs...)
	}
	titleArgs = append(titleArgs, videoURL)

	titleCmd := exec.Command("/usr/bin/yt-dlp", titleArgs...)
	// 设置环境变量
	titleCmd.Env = append(os.Environ(), "LANG=C.UTF-8", "LC_ALL=C.UTF-8")
	titleOutput, err := titleCmd.CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("获取标题失败: %w\n输出: %s", err, string(titleOutput))
	}
	title := strings.TrimSpace(string(titleOutput))

	// 第二步：下载音频
	downloadArgs := []string{
		"-x",                    // 仅提取音频
		"--audio-format", "mp3", // 转换为 MP3
		"--audio-quality", "0",  // 最佳质量
		"-o", filename,          // 使用相对路径，不带扩展名
		"--no-playlist",         // 不下载播放列表
		"--no-warnings",         // 不显示警告
	}
	// 如果提供了 cookies 文件，添加到参数中
	if s.cookiesFile != "" {
		downloadArgs = append([]string{"--cookies", s.cookiesFile}, downloadArgs...)
	}
	downloadArgs = append(downloadArgs, videoURL)

	downloadCmd := exec.Command("/usr/bin/yt-dlp", downloadArgs...)
	// 设置工作目录
	downloadCmd.Dir = s.tempDir
	// 设置环境变量
	downloadCmd.Env = append(os.Environ(), "LANG=C.UTF-8", "LC_ALL=C.UTF-8")
	// 设置工作目录
	downloadCmd.Dir = s.tempDir

	// 执行下载
	output, err := downloadCmd.CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("下载失败: %w\n输出: %s", err, string(output))
	}

	// 获取文件信息
	info, err := os.Stat(tempFile)
	if err != nil {
		// 如果带.mp3后缀的文件不存在，尝试不带后缀的
		if info2, err2 := os.Stat(tempBase); err2 == nil {
			tempFile = tempBase
			info = info2
		} else {
			// 列出目录中的所有文件，帮助调试
			files, _ := os.ReadDir(s.tempDir)
			var fileList []string
			for _, f := range files {
				fileList = append(fileList, f.Name())
			}
			return "", nil, fmt.Errorf("获取文件信息失败: %w\n下载的文件: %s 或 %s\n目录内容: %v\n输出: %s",
				err, tempFile, tempBase, fileList, string(output))
		}
	}

	// 检查文件大小
	if info.Size() > s.maxSize {
		os.Remove(tempFile)
		return "", nil, fmt.Errorf("文件过大: %d MB (最大 %d MB)", info.Size()/1024/1024, s.maxSize/1024/1024)
	}

	// 检查文件是否为空
	if info.Size() == 0 {
		os.Remove(tempFile)
		return "", nil, fmt.Errorf("下载的文件为空\n输出: %s", string(output))
	}

	// 解析标题作为歌曲信息
	songInfo := s.parseTitle(title)

	// 获取时长
	duration, _ := s.getDuration(tempFile)
	songInfo.Duration = duration

	return tempFile, songInfo, nil
}

// parseTitle 解析标题
func (s *YTDLPService) parseTitle(title string) *SongInfo {
	info := &SongInfo{
		Title:  strings.TrimSpace(title),
		Artist: "未知歌手",
		Album:  "",
	}

	// 尝试解析 "歌手 - 歌名" 格式
	if idx := strings.Index(title, " - "); idx != -1 {
		info.Artist = strings.TrimSpace(title[:idx])
		info.Title = strings.TrimSpace(title[idx+3:])
	}

	// 移除常见后缀
	info.Title = strings.TrimSuffix(info.Title, "Official Video")
	info.Title = strings.TrimSuffix(info.Title, "MV")
	info.Title = strings.TrimSuffix(info.Title, "Lyrics")
	info.Title = strings.TrimSpace(info.Title)

	return info
}

// getDuration 获取音频时长
func (s *YTDLPService) getDuration(filePath string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, nil // 忽略错误
	}

	var duration float64
	fmt.Sscanf(string(output), "%f", &duration)
	return int(duration), nil
}

// uploadToTelegram 上传到 Telegram
func (s *YTDLPService) uploadToTelegram(chatID int64, filePath string, songInfo *SongInfo) (string, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	upload := tgbotapi.NewAudio(chatID, tgbotapi.FileReader{
		Name:   fmt.Sprintf("%s - %s.mp3", songInfo.Artist, songInfo.Title),
		Reader: file,
	})
	upload.Title = songInfo.Title
	upload.Performer = songInfo.Artist
	upload.Caption = fmt.Sprintf("🎵 %s - %s\n\n⏰ %d秒",
		songInfo.Artist, songInfo.Title, songInfo.Duration)

	msg, err := s.bot.Send(upload)
	if err != nil {
		return "", 0, err
	}

	return msg.Audio.FileID, int64(msg.Audio.FileSize), nil
}

// sendDownloadedSong 发送已下载的歌曲
func (s *YTDLPService) sendDownloadedSong(chatID int64, song *model.Song, user *model.User) error {
	// 构建音频文件
	audio := tgbotapi.NewAudio(chatID, tgbotapi.FileID(song.FileID))
	audio.Title = song.Title
	audio.Performer = song.Artist

	// 构建说明文本
	var caption strings.Builder
	caption.WriteString(fmt.Sprintf("🎵 %s - %s", song.Artist, song.Title))
	if song.Album != "" {
		caption.WriteString(fmt.Sprintf("\n💿 %s", song.Album))
	}
	caption.WriteString(fmt.Sprintf("\n\n%s %s", song.GetCountryEmoji(), song.GetYearText()))
	audio.Caption = caption.String()

	// 创建操作按钮
	var keyboard [][]tgbotapi.InlineKeyboardButton
	favoriteBtn := tgbotapi.NewInlineKeyboardButtonData("❤️ 收藏", fmt.Sprintf("fav_%d", song.ID))
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{favoriteBtn})
	audio.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	// 发送音频
	_, err := s.bot.Send(audio)
	if err != nil {
		return err
	}

	// 记录历史
	historyRepo := database.NewHistoryRepository()
	historyRepo.Add(user.ID, song.ID)

	return nil
}

// generateHash 生成哈希
func (s *YTDLPService) generateHash(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// detectRegion 检测地区
func (s *YTDLPService) detectRegion(artist string) string {
	// 简单实现
	if s.containsCJK(artist) {
		return "CN"
	}
	return "US"
}

// containsCJK 检测是否包含中日韩文字
func (s *YTDLPService) containsCJK(str string) bool {
	for _, r := range str {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}
