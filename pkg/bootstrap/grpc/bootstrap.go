package grpc

import (
	"context"
	"net"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/handler"
	lib_grpc_inject_ctx "github.com/flamefatex/grpc-gateway-example/pkg/lib/grpc/middleware/inject_ctx"
	lib_grpc_recovery "github.com/flamefatex/grpc-gateway-example/pkg/lib/grpc/middleware/recovery"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/logx"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// BootstrapGrpcServer 启动grpc服务
func BootstrapGrpcServer(ctx context.Context) {
	// grpc
	lis, err := net.Listen("tcp", definition.GrpcAddr)
	if err != nil {
		logx.Fatalf(ctx, "grpc failed to listen: %s", err)
	}

	// server and middleware
	// 设置prometheus
	grpc_prometheus.EnableHandlingTimeHistogram()
	// grpc中间件在上的先执行
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(definition.OrigZap, grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel)),
			lib_grpc_inject_ctx.StreamServerInterceptor(lib_grpc_inject_ctx.WithInjectCtxFunc(lib_grpc_inject_ctx.DefaultInjectFunc)), // 注入request_id
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(lib_grpc_recovery.PrintStackHandlerFuncContext)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(definition.OrigZap, grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel)),
			lib_grpc_inject_ctx.UnaryServerInterceptor(lib_grpc_inject_ctx.WithInjectCtxFunc(lib_grpc_inject_ctx.DefaultInjectFunc)), // 注入request_id
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(lib_grpc_recovery.PrintStackHandlerFuncContext)),
		)),
	)

	// graceful shutdown
	go func() {
		<-ctx.Done()
		// sig is a ^C, handle it
		log.Warn("shutting down grpc server...")
		grpcServer.GracefulStop()
	}()

	// 注册grpc服务处理器
	handler.ExecRegisterGrpcServiceServer(ctx, grpcServer)
	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(grpcServer)

	logx.Debug(ctx, "start grpc server")
	if err := grpcServer.Serve(lis); err != nil {
		logx.Fatalf(ctx, "grpc failed to serve: %s", err)
	}
	return
}
