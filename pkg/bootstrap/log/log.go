package log

import (
	"context"
	"strings"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	lib_log "github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	lib_zap "github.com/flamefatex/grpc-gateway-example/pkg/lib/log/zap"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
)

func BootstrapLog(ctx context.Context) {
	// 构建
	level, _ := zap.ParseAtomicLevel(config.Config().GetString("log.level"))
	development := config.Config().GetBool("log.development")
	zapConfig := zap.Config{
		Level:       level,
		Development: development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	opts := []zap.Option{
		zap.Fields(zap.String("service", strings.ToLower(definition.ServiceName))),
	}
	origZap, _ := zapConfig.Build(opts...)

	// 设置global logger
	zapLogger := lib_zap.NewZapLogger(origZap)
	lib_log.SetLogger(zapLogger)

	// grpc日志
	if config.Config().GetBool("log.enableGrpcLog") {
		grpc_zap.ReplaceGrpcLoggerV2(origZap.WithOptions(zap.AddCallerSkip(2)))
	}

	// 设置OrigZap
	definition.SetOrigZap(origZap)
}
