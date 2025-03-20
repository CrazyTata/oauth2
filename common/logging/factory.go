package logging

import (
	"fmt"
	"os"
	"strings"
)

// LoggerFactory 根据环境生成日志记录器
type LoggerFactory struct{}

// LogConfig 存储日志配置
type LogConfig struct {
	Environment string // dev, test, prod
	Level       LogLevel
	Directory   string
	Filename    string
}

// 默认配置
var defaultLogConfig = LogConfig{
	Environment: "dev",
	Level:       DEBUG,
	Directory:   "logs",
	Filename:    "application",
}

// NewLogger 根据配置创建一个新的日志记录器
func (f *LoggerFactory) NewLogger(config LogConfig) (Logger, error) {
	// 如果没有指定环境，尝试从环境变量获取
	if config.Environment == "" {
		env := os.Getenv("APP_ENV")
		if env != "" {
			config.Environment = env
		} else {
			config.Environment = defaultLogConfig.Environment
		}
	}

	// 小写并去除空格
	env := strings.ToLower(strings.TrimSpace(config.Environment))

	// 根据环境确定日志实现
	switch env {
	case "dev", "development":
		return NewConsoleLogger(config.Level), nil
	case "test", "testing", "prod", "production":
		if config.Directory == "" {
			config.Directory = defaultLogConfig.Directory
		}
		if config.Filename == "" {
			config.Filename = defaultLogConfig.Filename
		}
		return NewFileLogger(config.Level, config.Directory, config.Filename)
	default:
		return nil, fmt.Errorf("未知的环境配置: %s", config.Environment)
	}
}

// 全局日志工厂和日志实例
var (
	factory      = &LoggerFactory{}
	globalLogger Logger
)

// 初始化默认的全局日志记录器
func init() {
	var err error
	globalLogger, err = factory.NewLogger(defaultLogConfig)
	if err != nil {
		// 如果创建失败，退回到控制台日志
		fmt.Printf("无法初始化日志系统: %v，使用默认控制台日志记录器\n", err)
		globalLogger = NewConsoleLogger(DEBUG)
	}
}

// 设置全局日志记录器
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// GetGlobalLogger 获取全局日志记录器
func GetGlobalLogger() Logger {
	return globalLogger
}

// 全局日志方法
func Debug(format string, args ...interface{}) {
	globalLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	globalLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	globalLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	globalLogger.Fatal(format, args...)
}
