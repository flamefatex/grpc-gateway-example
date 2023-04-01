package inject_ctx

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/tracing/opentracing"
)

// DefaultInjectFunc 默认注入
func DefaultInjectFunc(ctx context.Context) context.Context {
	newCtx := ctx

	// 注入request_id
	newCtx = opentracing.InjectRequestId(newCtx)

	return newCtx
}
