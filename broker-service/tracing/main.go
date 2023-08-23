package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

func InitTracer(serviceName string) io.Closer {

	agentEndpoint := "jaeger:6831"

	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentEndpoint,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic("Could not initialize jaeger tracer: " + err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}
