package inject_ctx

import (
	"context"
	"net/http"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/tracing/opentracing"
)

// InjectCtxFunc 注入ctx方法
// 不要使用 req.Context() 和 req.WithContext()
type InjectCtxFunc func(ctx context.Context, req *http.Request) context.Context

// InjectRequestIdFunc RequestId注入
func InjectRequestIdFunc(ctx context.Context, req *http.Request) context.Context {
	newCtx := ctx

	// 注入request_id
	newCtx = opentracing.InjectRequestId(newCtx)

	return newCtx
}
