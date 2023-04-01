package grpc_gateway

import (
	"context"
	"time"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/handler"
	lib_grpc_client "github.com/flamefatex/grpc-gateway-example/pkg/lib/grpc/client"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetGrpcGatewayMux(ctx context.Context) (ggMux *runtime.ServeMux, err error) {
	// 设置grpc请求context的超时时间
	runtime.DefaultContextTimeout = 60 * time.Second

	// Note: Make sure the gRPC server is running properly and accessible
	ggMux = runtime.NewServeMux(
		GetGlobalServeMuxOptions()...,
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	// client中间件
	opts = append(opts, lib_grpc_client.GetGlobalMiddlewareDialOptions()...)

	// Register gRPC server endpoint
	regs := handler.ExecRegisterGrpcGatewayEndpoint(ctx)

	for _, reg := range regs {
		err = reg(ctx, ggMux, definition.GrpcAddr, opts)
		if err != nil {
			return nil, err
		}
	}
	return
}
