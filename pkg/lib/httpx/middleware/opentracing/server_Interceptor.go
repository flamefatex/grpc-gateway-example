package opentracing

import (
	"context"
	"net/http"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

// NewHandler 链路跟踪http中间件
// 参考
// https://grpc-ecosystem.github.io/grpc-gateway/docs/operations/tracing/#opentracing-support
// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/tracing/opentracing/server_interceptors.go
// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/tracing/opentracing/id_extract.go
func NewHandler(h http.Handler, opts ...Option) http.Handler {
	o := evaluateOptions(opts)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 跳过链路跟踪机制
		if o.filterOutFunc != nil && !o.filterOutFunc(r.Context(), r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}

		// 操作名称原值为 "ServeHTTP", 设为url路径
		opName := r.URL.Path
		if o.opNameFunc != nil {
			opName = o.opNameFunc(r.URL.Path)
		}

		// 先从http头提取之前的SpanContext信息
		parentSpanContext, err := o.tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			// 开始链路跟踪新的span
			serverSpan := o.tracer.StartSpan(
				opName,
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			ext.HTTPMethod.Set(serverSpan, r.Method)
			ext.HTTPUrl.Set(serverSpan, r.URL.String())

			// 注入链路跟踪信息到ctx tags，方便日志输出
			injectOpentracingIdsToTags(o.traceHeaderName, serverSpan, grpc_ctxtags.Extract(r.Context()))

			// 传递到request里
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))

			defer finishServerSpan(r.Context(), serverSpan)
		}

		h.ServeHTTP(w, r)

	})
}

func finishServerSpan(ctx context.Context, serverSpan opentracing.Span) {
	// Log context information
	tags := grpc_ctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		// Don't tag errors, log them instead.
		if vErr, ok := v.(error); ok {
			serverSpan.LogKV(k, vErr.Error())
		} else {
			serverSpan.SetTag(k, v)
		}
	}

	serverSpan.Finish()
}
