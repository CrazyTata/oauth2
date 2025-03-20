package logging

import (
	"fmt"
	"os"
)

// ConsoleLogger 将日志输出到控制台
type ConsoleLogger struct {
	MinLevel LogLevel
}

// NewConsoleLogger 创建一个新的控制台日志记录器
func NewConsoleLogger(minLevel LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		MinLevel: minLevel,
	}
}

// Debug 记录调试级别的日志
func (l *ConsoleLogger) Debug(format string, args ...interface{}) {
	if l.MinLevel <= DEBUG {
		fmt.Fprintln(os.Stdout, formatLogMessage(DEBUG, format, args...))
	}
}

// Info 记录信息级别的日志
func (l *ConsoleLogger) Info(format string, args ...interface{}) {
	if l.MinLevel <= INFO {
		fmt.Fprintln(os.Stdout, formatLogMessage(INFO, format, args...))
	}
}

// Warn 记录警告级别的日志
func (l *ConsoleLogger) Warn(format string, args ...interface{}) {
	if l.MinLevel <= WARN {
		fmt.Fprintln(os.Stdout, formatLogMessage(WARN, format, args...))
	}
}

// Error 记录错误级别的日志
func (l *ConsoleLogger) Error(format string, args ...interface{}) {
	if l.MinLevel <= ERROR {
		fmt.Fprintln(os.Stderr, formatLogMessage(ERROR, format, args...))
	}
}

// Fatal 记录致命级别的日志并退出程序
func (l *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if l.MinLevel <= FATAL {
		fmt.Fprintln(os.Stderr, formatLogMessage(FATAL, format, args...))
		os.Exit(1)
	}
}
