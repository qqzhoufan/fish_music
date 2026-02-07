package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/fish-music/internal/config"
	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/handler"
	"github.com/user/fish-music/internal/service"
	"github.com/user/fish-music/pkg/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	var logLevel logger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	case "warn":
		logLevel = logger.Error
	case "error":
		logLevel = logger.Silent
	default:
		logLevel = logger.Info
	}

	// 初始化数据库
	if err := database.Init(&cfg.Database, logLevel); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	// 自动迁移（已通过 SQL 初始化脚本完成，跳过）
	// if err := database.AutoMigrate(); err != nil {
	// 	log.Fatalf("数据库迁移失败: %v", err)
	// }

	// 确保临时目录存在
	if err := cfg.Download.EnsureTempDir(); err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}

	// 初始化 Bot
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("创建 Bot 失败: %v", err)
	}

	bot.Debug = cfg.Log.Level == "debug"
	log.Printf("Bot 已启动: %s", bot.Self.UserName)

	// 初始化处理器
	songRepo := database.NewSongRepository()
	userRepo := database.NewUserRepository()
	favoriteRepo := database.NewFavoriteRepository()
	historyRepo := database.NewHistoryRepository()

	// 初始化音乐 API 客户端
	musicAPI := api.NewNeteaseAPI(cfg.Search.APIURL)

	// 初始化 yt-dlp 下载服务
	ytdlpService := service.NewYTDLPService(
		bot,
		songRepo,
		cfg.Download.TempDir,
		cfg.Download.MaxFileSize,
		cfg.Download.CookiesFile,
	)

	botHandler := handler.NewBotHandler(
		bot,
		cfg.Bot.AdminID,
		songRepo,
		userRepo,
		favoriteRepo,
		historyRepo,
		musicAPI,
		ytdlpService,
		&cfg.Download,
	)

	// 设置更新配置
	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := bot.GetUpdatesChan(updateCfg)

	// 处理信号
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("收到停止信号，正在关闭...")
		bot.StopReceivingUpdates()
		cancel()
	}()

	// 主循环
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}

			// 处理私聊消息
			if update.Message != nil {
				if update.Message.Chat.Type == "private" {
					if err := botHandler.HandlePrivateMessage(update); err != nil {
						log.Printf("处理消息失败: %v", err)
					}
				}
			}

			// 处理回调查询
			if update.CallbackQuery != nil {
				if err := botHandler.HandleCallback(update.CallbackQuery); err != nil {
					log.Printf("处理回调失败: %v", err)
				}
			}
		}
	}
}
