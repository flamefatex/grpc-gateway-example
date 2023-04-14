package handler

import (
	"context"

	v1 "github.com/flamefatex/grpc-gateway-example/handler/api/v1"
	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	proto_v1_hbe "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/http_body_example"
	"google.golang.org/grpc"
)

// ExecRegisterGrpcServiceServer 注册grpc服务处理器
func ExecRegisterGrpcServiceServer(ctx context.Context, grpcServer *grpc.Server) {
	// 注册
	proto_v1_example.RegisterExampleServiceServer(grpcServer, v1.NewExampleHandler())
	proto_v1_hbe.RegisterHttpBodyExampleServiceServer(grpcServer, v1.NewHttpBodyExampleHandler())
}
