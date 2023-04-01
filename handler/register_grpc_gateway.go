package handler

import (
	"context"

	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GrpcGwRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// ExecRegisterGrpcGatewayEndpoint 注册grpc gateway 端点
func ExecRegisterGrpcGatewayEndpoint(ctx context.Context) []GrpcGwRegister {
	regs := []GrpcGwRegister{
		// 注册http端点
		proto_v1_example.RegisterExampleServiceHandlerFromEndpoint,
	}
	return regs
}
