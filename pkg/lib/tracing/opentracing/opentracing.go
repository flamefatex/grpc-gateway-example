package opentracing

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

const InjectKeyRequestId = "request_id"

// GetTraceIdFromCtx 获取链路跟踪id
func GetTraceIdFromCtx(ctx context.Context) string {
	tags := grpc_ctxtags.Extract(ctx)
	tagMap := tags.Values()

	if traceId, ok := tagMap[grpc_opentracing.TagTraceId]; ok {
		return traceId.(string)
	}
	return ""
}

// InjectRequestId 注入链路跟踪请求id
func InjectRequestId(ctx context.Context) context.Context {
	requestId := GetTraceIdFromCtx(ctx)
	tags := grpc_ctxtags.Extract(ctx)
	tags = tags.Set(InjectKeyRequestId, requestId)

	newCtx := grpc_ctxtags.SetInContext(ctx, tags)
	return newCtx
}
