package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

func InitTracer(serviceName string) io.Closer {
	// Jaeger agent endpoint
	agentEndpoint := "jaeger:6831"

	// Configuration for the Jaeger tracer
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // Always sample
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentEndpoint,
		},
	}

	// Initialize the tracer
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic("Could not initialize jaeger tracer: " + err.Error())
	}
	opentracing.SetGlobalTracer(tracer) // Set as the global tracer
	return closer
}
