package logx

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	logtypes "github.com/flamefatex/grpc-gateway-example/pkg/lib/log/types"
)

func Debug(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Errorf(format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Fatalf(format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	log.Extract(ctx).Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	log.Extract(ctx).Panicf(format, args...)
}

func WithField(key string, value interface{}) logtypes.Logger {
	return log.WithField(key, value)
}

func WithFields(fields logtypes.Fields) logtypes.Logger {
	return log.WithFields(fields)
}
