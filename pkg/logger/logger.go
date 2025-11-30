package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapLogger *zap.Logger
)

// filterCore 过滤特定日志消息的 Core
type filterCore struct {
	zapcore.Core
	ignoreMessages []string
}

// newFilterCore 创建过滤 Core
func newFilterCore(core zapcore.Core, ignoreMessages []string) zapcore.Core {
	return &filterCore{
		Core:           core,
		ignoreMessages: ignoreMessages,
	}
}

// Check 检查是否应该记录日志
func (c *filterCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// 检查消息是否应该被过滤
	for _, ignoreMsg := range c.ignoreMessages {
		msgLen := len(entry.Message)
		ignoreLen := len(ignoreMsg)
		if entry.Message == ignoreMsg || (msgLen >= ignoreLen && entry.Message[:ignoreLen] == ignoreMsg) {
			// 过滤掉该日志，不记录
			return ce
		}
	}
	// 不过滤，正常记录
	return c.Core.Check(entry, ce)
}

// With 添加字段
func (c *filterCore) With(fields []zapcore.Field) zapcore.Core {
	return &filterCore{
		Core:           c.Core.With(fields),
		ignoreMessages: c.ignoreMessages,
	}
}

// Init 初始化日志系统，使用 zap 作为 hlog 的底层实现
// env 参数可选值: development, testing, production
func Init(logDir string, logLevel string, env ...string) error {
	environment := "production"
	if len(env) > 0 && env[0] != "" {
		environment = env[0]
	}
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

	// 生产环境优化：禁用调用者信息以提升性能
	if environment == "production" {
		encoderConfig.CallerKey = zapcore.OmitKey
	}

	// 开发环境使用彩色输出
	if environment == "development" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
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

	// 添加过滤器，过滤掉不需要的错误消息
	filteredCore := newFilterCore(core, []string{
		"HERTZ: Error=accept tcp", // 过滤服务器关闭时的错误
	})

	// 添加日志采样：每秒相同消息最多记录 5 条，之后每 100 条记录 1 条
	sampledCore := zapcore.NewSamplerWithOptions(
		filteredCore,
		time.Second,
		5,   // Initial: 每秒最初的 5 条相同消息都记录
		100, // Thereafter: 之后每 100 条记录 1 条
	)

	// 创建 zap logger（添加调用者信息跳过层级）
	zapLogger = zap.New(sampledCore, zap.AddCaller(), zap.AddCallerSkip(1))

	// 使用 hertzzap 库设置 hlog
	// 直接使用自定义的 core 和配置
	hlog.SetLogger(hertzzap.NewLogger(
		hertzzap.WithZapOptions(
			zap.WrapCore(func(c zapcore.Core) zapcore.Core {
				return sampledCore
			}),
			zap.AddCaller(),
			zap.AddCallerSkip(2), // hertzzap 需要跳过额外的层级
		),
	))
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

// WithFields 创建带结构化字段的日志记录器
func WithFields(fields ...zap.Field) *zap.Logger {
	if zapLogger == nil {
		// 如果未初始化，返回一个 nop logger
		return zap.NewNop()
	}
	return zapLogger.With(fields...)
}

// MaskSensitive 脱敏敏感信息
// 对于长度 <= 4 的字符串，完全脱敏
// 对于长度 > 4 的字符串，保留首尾各 2 个字符
func MaskSensitive(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

// MaskEmail 脱敏邮箱地址
// 示例: user@example.com -> us****@ex****
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return MaskSensitive(email)
	}
	return MaskSensitive(parts[0]) + "@" + MaskSensitive(parts[1])
}

// MaskPhone 脱敏手机号
// 示例: 13812345678 -> 138****5678
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}
	if len(phone) != 11 {
		return MaskSensitive(phone)
	}
	return phone[:3] + "****" + phone[7:]
}
