package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Bot      BotConfig      `mapstructure:"bot"`
	Database DatabaseConfig `mapstructure:"database"`
	Web      WebConfig      `mapstructure:"web"`
	Download DownloadConfig `mapstructure:"download"`
	Search   SearchConfig   `mapstructure:"search"`
	Log      LogConfig      `mapstructure:"log"`
}

// BotConfig Telegram Bot 配置
type BotConfig struct {
	Token   string `mapstructure:"token"`
	AdminID int64  `mapstructure:"admin_id"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// WebConfig Web 管理端配置
type WebConfig struct {
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// DownloadConfig 下载配置
type DownloadConfig struct {
	WorkerCount int    `mapstructure:"worker_count"`
	MaxFileSize int    `mapstructure:"max_file_size"`
	TempDir     string `mapstructure:"temp_dir"`
}

// SearchConfig 搜索配置
type SearchConfig struct {
	APIURL  string `mapstructure:"api_url"`
	Timeout int    `mapstructure:"timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("bot.token", "")
	viper.SetDefault("bot.admin_id", 0)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "fish_music")
	viper.SetDefault("database.password", "fish_music_pass")
	viper.SetDefault("database.dbname", "fish_music")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("web.port", 9999)
	viper.SetDefault("web.username", "admin")
	viper.SetDefault("web.password", "admin123")
	viper.SetDefault("download.worker_count", 3)
	viper.SetDefault("download.max_file_size", 50)
	viper.SetDefault("download.temp_dir", "./tmp")
	viper.SetDefault("search.api_url", "")
	viper.SetDefault("search.timeout", 30)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", "")
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Bot.Token == "" {
		return fmt.Errorf("bot.token 不能为空")
	}
	if c.Bot.AdminID == 0 {
		return fmt.Errorf("bot.admin_id 不能为空")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database.host 不能为空")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database.dbname 不能为空")
	}
	return nil
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// EnsureTempDir 确保临时目录存在
func (d *DownloadConfig) EnsureTempDir() error {
	if d.TempDir == "" {
		d.TempDir = "./tmp"
	}
	if err := os.MkdirAll(d.TempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	return nil
}
