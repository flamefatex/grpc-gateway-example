package main

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/definition"
	bs_grpc "github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/grpc"
	bs_http "github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/http"
	bs_log "github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/log"
	bs_opentracing "github.com/flamefatex/grpc-gateway-example/pkg/bootstrap/opentracing"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
)

func main() {
	ctx := context.Background()
	// 注入service name
	log.SetLogger(log.WithField("service", definition.ServiceName))
	// 初始化config
	config.Init(definition.ServiceName)

	// 引导加载自定义log
	bs_log.BootstrapLog(ctx)
	// 引导加载opentracing
	closer := bs_opentracing.BootstrapOpentracing(ctx)
	defer func() {
		if closer != nil {
			_ = closer.Close()
		}
	}()

	// 启动http
	bs_http.BootstrapHttpServer(ctx)
	// 启动grpc
	bs_grpc.BootstrapGrpcServer(ctx)

}
