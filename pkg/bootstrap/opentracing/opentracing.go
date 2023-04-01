package opentracing

import (
	"context"
	"io"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/opentracing/opentracing-go"
	jaeger_cfg "github.com/uber/jaeger-client-go/config"
	jaeger_prometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

func BootstrapOpentracing(ctx context.Context) io.Closer {
	c := &jaeger_cfg.Configuration{
		ServiceName: definition.ServiceName,
		Sampler: &jaeger_cfg.SamplerConfig{
			Type:  config.Config().GetString("opentracing.sampler.type"),
			Param: config.Config().GetFloat64("opentracing.sampler.param"),
		},
		Reporter: &jaeger_cfg.ReporterConfig{
			LogSpans:           config.Config().GetBool("opentracing.reporter.logSpans"),
			LocalAgentHostPort: config.Config().GetString("opentracing.reporter.localAgentHostPort"),
		},
	}
	opts := []jaeger_cfg.Option{
		jaeger_cfg.Gen128Bit(true),
		jaeger_cfg.Metrics(jaeger_prometheus.New()), // prometheus 指标
	}
	tracer, closer, err := c.NewTracer(opts...)

	if err != nil {
		log.Warnf("new tracer error: %s", err)
		return closer
	}
	opentracing.InitGlobalTracer(tracer)
	return closer
}
