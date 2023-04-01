package zap

import (
	"context"

	logtypes "github.com/flamefatex/grpc-gateway-example/pkg/lib/log/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
}

func NewZapLogger(z *zap.Logger) *ZapLogger {
	logger := z.WithOptions(zap.AddCallerSkip(2))
	return &ZapLogger{
		Logger:        logger,
		SugaredLogger: logger.Sugar(),
	}
}

func NewZap(opts ...zap.Option) (*zap.Logger, error) {
	return zap.NewProduction(opts...)
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.SugaredLogger.Debug(args...)
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.SugaredLogger.Debugf(format, args...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.SugaredLogger.Info(args...)
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.SugaredLogger.Infof(format, args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.SugaredLogger.Warn(args...)
}

func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	z.SugaredLogger.Warnf(format, args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.SugaredLogger.Error(args...)
}

func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.SugaredLogger.Errorf(format, args...)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.SugaredLogger.Fatal(args...)
}

func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	z.SugaredLogger.Fatalf(format, args...)
}

func (z *ZapLogger) Panic(args ...interface{}) {
	z.SugaredLogger.Panic(args...)
}

func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	z.SugaredLogger.Panicf(format, args...)
}

func (z *ZapLogger) WithField(key string, value interface{}) logtypes.Logger {
	zapField := zap.Any(key, value)
	nz := z.Logger.With(zapField)
	return NewZapLogger(nz)
}
func (z *ZapLogger) WithFields(fields logtypes.Fields) logtypes.Logger {
	zapFields := make([]zap.Field, 0)

	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	nz := z.Logger.With(zapFields...)
	return NewZapLogger(nz)
}

// Extract 提取ctx的tags并返回logger
func (z *ZapLogger) Extract(ctx context.Context) logtypes.Logger {
	nz := ctxzap.Extract(ctx)
	return NewZapLogger(nz)
}
