package grpc_gateway

import (
	"context"
	"net/http"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

const HttpMethodKey = "ffx-http-method"
const HttpRequestUriKey = "ffx-http-request-uri"

// GetGlobalServeMuxOptions 获取全局mux选项
func GetGlobalServeMuxOptions() []runtime.ServeMuxOption {
	opts := []runtime.ServeMuxOption{
		// http 数据往grpc传递
		// 参考https://grpc-ecosystem.github.io/grpc-gateway/docs/operations/annotated_context/#get-http-path-pattern
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := make(map[string]string)
			md[HttpMethodKey] = req.Method
			md[HttpRequestUriKey] = req.RequestURI
			if method, ok := runtime.RPCMethod(ctx); ok {
				md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
			}
			if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
				md["pattern"] = pattern // /v1/example/login
			}

			return metadata.New(md)
		}),
		// 自定义错误处理
		runtime.WithErrorHandler(CustomHTTPErrorHandler),
		// 序列化选项
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   false,
					UseEnumNumbers:  true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
	}
	return opts
}

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
