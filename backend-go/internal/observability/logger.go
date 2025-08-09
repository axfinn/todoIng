package observability

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *log.Logger

// InitLogger 初始化日志系统，支持日志滚动
func InitLogger() {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return
	}

	// 配置日志轮转
	logRotator := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"), // 日志文件路径
		MaxSize:    100,                              // 单个日志文件最大大小（MB）
		MaxBackups: 30,                               // 保留的旧日志文件数量
		MaxAge:     7,                                // 日志文件保留天数
		Compress:   true,                             // 是否压缩旧日志文件
		LocalTime:  true,                             // 使用本地时间
	}

	// 创建多重写入器，同时写入文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logRotator)

	// 创建自定义 logger
	logger = log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)

	// 替换全局 logger
	log.SetOutput(logRotator)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger.Println("Logger initialized with file rotation")
	logger.Printf("Log configuration: MaxSize=%dMB, MaxBackups=%d, MaxAge=%ddays",
		logRotator.MaxSize, logRotator.MaxBackups, logRotator.MaxAge)
}

// CtxLog 输出带时间戳和上下文的日志
func CtxLog(ctx context.Context, format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")

	// 如果 logger 没有初始化，使用默认 logger
	if logger == nil {
		log.Printf("[%s] %s", ts, fmt.Sprintf(format, args...))
		return
	}

	logger.Printf("[%s] %s", ts, fmt.Sprintf(format, args...))
}

// LogInfo 记录信息级别日志
func LogInfo(format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	if logger == nil {
		log.Printf("[%s] [INFO] %s", ts, fmt.Sprintf(format, args...))
		return
	}
	logger.Printf("[%s] [INFO] %s", ts, fmt.Sprintf(format, args...))
}

// LogError 记录错误级别日志
func LogError(format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	if logger == nil {
		log.Printf("[%s] [ERROR] %s", ts, fmt.Sprintf(format, args...))
		return
	}
	logger.Printf("[%s] [ERROR] %s", ts, fmt.Sprintf(format, args...))
}

// LogWarn 记录警告级别日志
func LogWarn(format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	if logger == nil {
		log.Printf("[%s] [WARN] %s", ts, fmt.Sprintf(format, args...))
		return
	}
	logger.Printf("[%s] [WARN] %s", ts, fmt.Sprintf(format, args...))
}

// LogDebug 记录调试级别日志
func LogDebug(format string, args ...any) {
	// 只在调试模式下输出
	if os.Getenv("DEBUG") == "true" {
		ts := time.Now().Format("2006-01-02 15:04:05.000")
		if logger == nil {
			log.Printf("[%s] [DEBUG] %s", ts, fmt.Sprintf(format, args...))
			return
		}
		logger.Printf("[%s] [DEBUG] %s", ts, fmt.Sprintf(format, args...))
	}
}
