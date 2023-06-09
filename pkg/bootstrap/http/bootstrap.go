package http

import (
	"context"
	"net/http"
	"time"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/http/gin"
	"github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/http/grpc_gateway"
	"github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/http/pprof"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	httpx_ctxtags "github.com/flamefatex/grpc-gateway-example/pkg/lib/httpx/middleware/ctxtags"
	httpx_inject_ctx "github.com/flamefatex/grpc-gateway-example/pkg/lib/httpx/middleware/inject_ctx"
	httpx_zap "github.com/flamefatex/grpc-gateway-example/pkg/lib/httpx/middleware/logging/zap"
	httpx_ot "github.com/flamefatex/grpc-gateway-example/pkg/lib/httpx/middleware/opentracing"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/logx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func BootstrapHttpServer(ctx context.Context) {
	httpMux := http.NewServeMux()
	// gin
	ginRouter := gin.GetGinRouter(ctx)

	// grpc_gateway
	ggMux, err := grpc_gateway.GetGrpcGatewayMux(ctx)
	if err != nil {
		logx.Fatalf(ctx, "GetGrpcGatewayMux failed, err:%s", err)
	}

	// 注册路由
	httpMux.Handle("/g/", ginRouter)
	httpMux.Handle("/", ggMux)
	httpMux.Handle("/metrics", promhttp.Handler()) // Register Prometheus metrics handler.
	if config.Config().GetBool("pprof.enabled") {
		httpMux.Handle("/debug/pprof/", pprof.GetPprofMux(ctx)) // pprof
	}

	// http中间件 从下往上执行
	var handler http.Handler = httpMux

	handler = httpx_inject_ctx.NewHandler(handler,
		httpx_inject_ctx.WithInjectCtxFunc(httpx_inject_ctx.DefaultInjectFunc),
	) // 注入request_id
	handler = httpx_zap.NewHandler(handler, definition.OrigZap)
	handler = httpx_ot.NewHandler(handler)      // 链路跟踪
	handler = httpx_ctxtags.NewHandler(handler) // ctx tags

	httpServer := &http.Server{
		Addr:         definition.HttpAddr,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		// add handler with middleware
		Handler: handler,
	}

	// graceful shutdown
	go func() {
		<-ctx.Done()
		// sig is a ^C, handle it
		log.Warn("shutting down http server...")
		_ = httpServer.Shutdown(ctx)
	}()

	go func() {
		logx.Debug(ctx, "start http server")
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		if err = httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logx.Fatalf(ctx, "httpServer ListenAndServe failed, err:%s", err)
			}

		}
	}()
	return
}
