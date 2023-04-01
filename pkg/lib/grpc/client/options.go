package client

import (
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

// GetGlobalMiddlewareDialOptions 获取grpc client 全局中间件选项
func GetGlobalMiddlewareDialOptions() []grpc.DialOption {
	// grpc client 全局中间件选项
	// 注意不能当成全局变量，否则默认值，会最开始初始化会造成影响，
	// 如opentracing.InitGlobalTracer还没用，就执行了下面的方法
	var globalMiddlewareDialOptions = []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(),
		),
		grpc.WithChainStreamInterceptor(
			grpc_opentracing.StreamClientInterceptor(),
		),
	}

	return globalMiddlewareDialOptions
}
