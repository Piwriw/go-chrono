package chrono

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig 定义日志配置
type LoggerConfig struct {
	FilePath        string
	LogLevel        string
	MaxSize         int
	MaxBackups      int
	MaxAge          int
	IsCompress      bool
	ReplaceAttrFunc func(groups []string, a slog.Attr) slog.Attr
}

// 默认日志配置
var defaultLoggerConfig = LoggerConfig{
	FilePath:   "app.log",
	LogLevel:   "info",
	MaxSize:    10,   // 默认 10MB
	MaxBackups: 5,    // 默认最多保留 5 个备份
	MaxAge:     7,    // 默认 7 天
	IsCompress: true, // 默认启用压缩
	ReplaceAttrFunc: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.UTC().Format("2006-01-02 15:04:05"))
		}
		return a
	},
}

// Clone 复制默认配置
func (c *LoggerConfig) Clone() *LoggerConfig {
	cfg := *c
	return &cfg
}

// mergeConfig 合并用户配置和默认配置
func mergeConfig(config *LoggerConfig) *LoggerConfig {
	cfg := defaultLoggerConfig.Clone()
	if config == nil {
		return cfg
	}
	if config.FilePath != "" {
		cfg.FilePath = config.FilePath
	}
	if config.LogLevel != "" {
		cfg.LogLevel = config.LogLevel
	}
	if config.MaxSize > 0 {
		cfg.MaxSize = config.MaxSize
	}
	if config.MaxBackups > 0 {
		cfg.MaxBackups = config.MaxBackups
	}
	if config.MaxAge > 0 {
		cfg.MaxAge = config.MaxAge
	}
	if config.ReplaceAttrFunc != nil {
		cfg.ReplaceAttrFunc = config.ReplaceAttrFunc
	}
	cfg.IsCompress = config.IsCompress // bool 默认值为 false，不需要额外判断
	return cfg
}

// parseLogLevel 解析日志级别
func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// InitLogger 初始化日志
func InitLogger(config *LoggerConfig) *lumberjack.Logger {
	cfg := mergeConfig(config)

	// 设置日志文件轮转
	logFile := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.IsCompress,
	}

	// 设置控制台和文件输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 创建日志 Handler
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		AddSource:   true,
		Level:       parseLogLevel(cfg.LogLevel),
		ReplaceAttr: cfg.ReplaceAttrFunc,
	})

	// 创建 logger 实例
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logFile
}
