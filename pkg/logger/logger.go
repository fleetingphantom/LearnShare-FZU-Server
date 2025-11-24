package logger

import (
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init 初始化日志系统，使用 zap 作为 hlog 的底层实现
func Init(logDir string, logLevel string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 解析日志级别
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置日志文件输出 - 普通日志
	infoLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    100, // MB
		MaxBackups: 30,  // 保留30个备份
		MaxAge:     7,   // 保留7天
		Compress:   true,
		LocalTime:  true,
	}

	// 配置日志文件输出 - 错误日志
	errorLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "error.log"),
		MaxSize:    100, // MB
		MaxBackups: 30,  // 保留30个备份
		MaxAge:     14,  // 保留14天
		Compress:   true,
		LocalTime:  true,
	}

	// 创建 Core
	// 控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 所有日志输出到 app.log
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	allFileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(infoLogFile),
		level,
	)

	// 错误日志输出到 error.log
	errorFileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(errorLogFile),
		zapcore.ErrorLevel, // 只记录 error 及以上级别
	)

	// 组合多个 Core
	core := zapcore.NewTee(consoleCore, allFileCore, errorFileCore)

	// 设置 hlog 使用 zap logger
	logger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core
	})))

	hlog.SetLogger(logger)
	hlog.SetLevel(convertToHlogLevel(level))

	return nil
}

// convertToHlogLevel 将 zap 日志级别转换为 hlog 日志级别
func convertToHlogLevel(level zapcore.Level) hlog.Level {
	switch level {
	case zapcore.DebugLevel:
		return hlog.LevelDebug
	case zapcore.InfoLevel:
		return hlog.LevelInfo
	case zapcore.WarnLevel:
		return hlog.LevelWarn
	case zapcore.ErrorLevel:
		return hlog.LevelError
	case zapcore.FatalLevel:
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}

// 以下是便捷函数，直接使用 hlog

// Info 信息日志
func Info(v ...interface{}) {
	hlog.Info(v...)
}

// Infof 格式化信息日志
func Infof(format string, v ...interface{}) {
	hlog.Infof(format, v...)
}

// Warn 警告日志
func Warn(v ...interface{}) {
	hlog.Warn(v...)
}

// Warnf 格式化警告日志
func Warnf(format string, v ...interface{}) {
	hlog.Warnf(format, v...)
}

// Error 错误日志
func Error(v ...interface{}) {
	hlog.Error(v...)
}

// Errorf 格式化错误日志
func Errorf(format string, v ...interface{}) {
	hlog.Errorf(format, v...)
}

// Fatal 致命错误日志
func Fatal(v ...interface{}) {
	hlog.Fatal(v...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, v ...interface{}) {
	hlog.Fatalf(format, v...)
}

// Debug 调试日志
func Debug(v ...interface{}) {
	hlog.Debug(v...)
}

// Debugf 格式化调试日志
func Debugf(format string, v ...interface{}) {
	hlog.Debugf(format, v...)
}

// CtxInfof 带上下文的格式化信息日志
func CtxInfof(format string, v ...interface{}) {
	hlog.CtxInfof(nil, format, v...)
}

// CtxErrorf 带上下文的格式化错误日志
func CtxErrorf(format string, v ...interface{}) {
	hlog.CtxErrorf(nil, format, v...)
}

// CtxWarnf 带上下文的格式化警告日志
func CtxWarnf(format string, v ...interface{}) {
	hlog.CtxWarnf(nil, format, v...)
}
