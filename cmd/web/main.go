package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/fish-music/internal/config"
	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/handler"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置 Gin 模式
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
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

	// 初始化处理器
	songRepo := database.NewSongRepository()
	webHandler := handler.NewWebHandler(
		cfg.Web.Username,
		cfg.Web.Password,
		songRepo,
	)

	// 创建 Gin 路由
	router := gin.Default()

	// 加载 HTML 模板
	router.LoadHTMLGlob("web/templates/*")

	// 静态文件
	router.Static("/static", "web/static")

	// 注册路由
	webHandler.RegisterRoutes(router)

	// 创建服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Web.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		log.Printf("Web 服务器已启动: http://0.0.0.0:%d", cfg.Web.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Web 服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭 Web 服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Web 服务器关闭失败: %v", err)
	}

	log.Println("Web 服务器已关闭")
}
