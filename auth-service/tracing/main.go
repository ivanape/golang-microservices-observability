package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

//func InitTracer(serviceName string) {
//	JAEGER_ENDPOINT := "http://otel:14268/api/traces"
//	// Create Jaeger exporter
//	exporter, err := jaeger.New(
//
//		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(JAEGER_ENDPOINT)),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Resource can be modelled after the service
//	resource1 := resource.NewWithAttributes(
//		semconv.SchemaURL,
//		semconv.ServiceNameKey.String(serviceName),
//		semconv.ServiceVersionKey.String("v0.0.1"),
//	)
//
//	// Tracer Provider
//	tp := sdktrace.NewTracerProvider(
//		sdktrace.WithBatcher(exporter),
//		sdktrace.WithResource(resource1),
//		// Configure the sampler to always sample.
//		sdktrace.WithSampler(sdktrace.AlwaysSample()),
//	)
//
//	otel.SetTracerProvider(tp)
//
//	// Set the propagator to tracecontext (the default is baggagetracecontext)
//	otel.SetTextMapPropagator(propagation.TraceContext{})
//}

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
