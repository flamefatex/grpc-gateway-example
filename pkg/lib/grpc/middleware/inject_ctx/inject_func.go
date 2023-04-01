package inject_ctx

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/tracing/opentracing"
)

// InjectCtxFunc 注入ctx方法
type InjectCtxFunc func(ctx context.Context) context.Context

// InjectRequestIdFunc RequestId注入
func InjectRequestIdFunc(ctx context.Context) context.Context {
	newCtx := ctx

	// 注入request_id
	newCtx = opentracing.InjectRequestId(newCtx)

	return newCtx
}
