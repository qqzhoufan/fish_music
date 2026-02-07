package config

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// StdLogger 标准日志实现
type StdLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

// NewLogger 创建日志器
func NewLogger(cfg *LogConfig) (Logger, func(), error) {
	var output io.Writer = os.Stdout

	// 如果指定了日志文件，输出到文件
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, nil, err
		}
		output = file

		// 返回清理函数
		cleanup := func() {
			file.Close()
		}

		flags := log.Ldate | log.Ltime | log.Lshortfile
		return &StdLogger{
			debug: log.New(output, "[DEBUG] ", flags),
			info:  log.New(output, "[INFO] ", flags),
			warn:  log.New(output, "[WARN] ", flags),
			error: log.New(output, "[ERROR] ", flags),
		}, cleanup, nil
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	return &StdLogger{
		debug: log.New(output, "[DEBUG] ", flags),
		info:  log.New(output, "[INFO] ", flags),
		warn:  log.New(output, "[WARN] ", flags),
		error: log.New(output, "[ERROR] ", flags),
	}, func() {}, nil
}

func (l *StdLogger) Debug(msg string, args ...interface{}) {
	l.debug.Output(2, fmt.Sprintf(msg, args...))
}

func (l *StdLogger) Info(msg string, args ...interface{}) {
	l.info.Output(2, fmt.Sprintf(msg, args...))
}

func (l *StdLogger) Warn(msg string, args ...interface{}) {
	l.warn.Output(2, fmt.Sprintf(msg, args...))
}

func (l *StdLogger) Error(msg string, args ...interface{}) {
	l.error.Output(2, fmt.Sprintf(msg, args...))
}
