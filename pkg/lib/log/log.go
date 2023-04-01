package log

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log/types"
	lib_zap "github.com/flamefatex/grpc-gateway-example/pkg/lib/log/zap"
)

var globalLogger types.Logger

func init() {
	logger, _ := lib_zap.NewZap()
	SetLogger(lib_zap.NewZapLogger(logger))
}

// SetLogger 设置全局logger
func SetLogger(logger types.Logger) {
	globalLogger = logger
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

func WithField(key string, value interface{}) types.Logger {
	return globalLogger.WithField(key, value)
}

func WithFields(fields types.Fields) types.Logger {
	return globalLogger.WithFields(fields)
}

// Extract 提取ctx的tags并返回logger
func Extract(ctx context.Context) types.Logger {
	return globalLogger.Extract(ctx)
}
