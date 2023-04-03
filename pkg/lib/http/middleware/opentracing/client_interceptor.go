package opentracing

import (
	"context"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/grpclog"
)

var (
	httpTag = opentracing.Tag{Key: string(ext.Component), Value: "HTTP"}
)

type opentracingTransport struct {
	nextRoundTripper http.RoundTripper //  链路跟踪header添加完，再调用nextRoundTripper的RoundTripp()
	o                *options
}

func OpentracingRoundTripper(nextRoundTripper http.RoundTripper, opts ...Option) http.RoundTripper {
	return &opentracingTransport{
		nextRoundTripper: nextRoundTripper,
		o:                evaluateOptions(opts),
	}
}

// ClientAddContextTags returns a context with specified opentracing tags, which
// are used by UnaryClientInterceptor/StreamClientInterceptor when creating a
// new span.
func ClientAddContextTags(ctx context.Context, tags opentracing.Tags) context.Context {
	return context.WithValue(ctx, clientSpanTagKey{}, tags)
}

type clientSpanTagKey struct{}

func (c *opentracingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// 获取父层span
	ctx := req.Context()
	var parentSpanCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentSpanCtx = parent.Context()
	}

	// 开启新的span
	opts := []opentracing.StartSpanOption{
		opentracing.ChildOf(parentSpanCtx),
		ext.SpanKindRPCClient,
		httpTag,
	}
	if tagx := ctx.Value(clientSpanTagKey{}); tagx != nil {
		if opt, ok := tagx.(opentracing.StartSpanOption); ok {
			opts = append(opts, opt)
		}
	}
	clientSpan := c.o.tracer.StartSpan(
		"Go-http-client", // 操作名
		opts...,
	)
	ext.HTTPMethod.Set(clientSpan, req.Method)
	ext.HTTPUrl.Set(clientSpan, req.URL.String())

	// 添加到http头
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	if iErr := clientSpan.Tracer().Inject(clientSpan.Context(), opentracing.HTTPHeaders, carrier); iErr != nil {
		grpclog.Infof("grpc_opentracing: failed serializing trace information: %v", iErr)
	}

	// 运行 nextRoundTripper
	resp, err = c.nextRoundTripper.RoundTrip(req)

	finishClientSpan(clientSpan, resp, err)

	return

}

func finishClientSpan(clientSpan opentracing.Span, resp *http.Response, err error) {
	if err != nil && err != io.EOF {
		ext.Error.Set(clientSpan, true)
		clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}
	if resp != nil {
		ext.HTTPStatusCode.Set(clientSpan, uint16(resp.StatusCode))
	}

	clientSpan.Finish()
}
