package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Log 是全局日志实例
	Log *logrus.Logger
)

// 初始化日志
func init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}

// SetLogLevel 设置日志级别
func SetLogLevel(level string) {
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}

// SetLogOutput 设置日志输出到文件和控制台
func SetLogOutput(logPath string) {
	if logPath == "" {
		return
	}

	// 确保日志目录存在
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		Log.Errorf("创建日志目录失败: %v", err)
		return
	}

	// 打开日志文件
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Errorf("打开日志文件失败: %v", err)
		return
	}

	// 同时输出到文件和控制台
	mw := io.MultiWriter(os.Stdout, file)
	Log.SetOutput(mw)
}

// WithCaller 添加调用者信息
func WithCaller() *logrus.Entry {
	_, file, line, _ := runtime.Caller(1)
	return Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	})
}

// WithModule 添加模块信息
func WithModule(module string) *logrus.Entry {
	return Log.WithField("module", module)
}

// Debug 调试日志
func Debug(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Debugf(format, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Info(args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Infof(format, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Warnf(format, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Error(args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Errorf(format, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Fatalf(format, args...)
}

// Panic 触发panic的日志
func Panic(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Panic(args...)
}

// Panicf 格式化触发panic的日志
func Panicf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}).Panicf(format, args...)
}
