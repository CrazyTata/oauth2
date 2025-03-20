package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileLogger 将日志输出到文件
type FileLogger struct {
	MinLevel  LogLevel
	directory string
	filename  string
	file      *os.File
	mu        sync.Mutex
	date      string
}

// NewFileLogger 创建一个新的文件日志记录器
func NewFileLogger(minLevel LogLevel, directory, filename string) (*FileLogger, error) {
	logger := &FileLogger{
		MinLevel:  minLevel,
		directory: directory,
		filename:  filename,
	}

	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("无法创建日志目录: %w", err)
	}

	if err := logger.rotateFileIfNeeded(); err != nil {
		return nil, err
	}

	return logger, nil
}

// rotateFileIfNeeded 如果需要，旋转日志文件（每天一个日志文件）
func (l *FileLogger) rotateFileIfNeeded() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if l.file != nil && l.date == currentDate {
		return nil
	}

	// 关闭旧文件
	if l.file != nil {
		l.file.Close()
	}

	l.date = currentDate
	fullFilename := filepath.Join(l.directory, fmt.Sprintf("%s_%s.log", l.filename, currentDate))
	file, err := os.OpenFile(fullFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("无法打开日志文件: %w", err)
	}

	l.file = file
	return nil
}

// writeLog 写入一条日志到文件
func (l *FileLogger) writeLog(level LogLevel, format string, args ...interface{}) {
	if l.MinLevel > level {
		return
	}

	if err := l.rotateFileIfNeeded(); err != nil {
		fmt.Fprintf(os.Stderr, "日志文件轮转失败: %v\n", err)
		return
	}

	logMessage := formatLogMessage(level, format, args...)

	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Fprintln(l.file, logMessage)
}

// Debug 记录调试级别的日志
func (l *FileLogger) Debug(format string, args ...interface{}) {
	l.writeLog(DEBUG, format, args...)
}

// Info 记录信息级别的日志
func (l *FileLogger) Info(format string, args ...interface{}) {
	l.writeLog(INFO, format, args...)
}

// Warn 记录警告级别的日志
func (l *FileLogger) Warn(format string, args ...interface{}) {
	l.writeLog(WARN, format, args...)
}

// Error 记录错误级别的日志
func (l *FileLogger) Error(format string, args ...interface{}) {
	l.writeLog(ERROR, format, args...)
}

// Fatal 记录致命级别的日志并退出程序
func (l *FileLogger) Fatal(format string, args ...interface{}) {
	l.writeLog(FATAL, format, args...)
	os.Exit(1)
}

// Close 关闭日志文件
func (l *FileLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
