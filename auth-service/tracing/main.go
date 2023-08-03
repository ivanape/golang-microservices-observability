package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
)

func InitTracer(serviceName string) {
	JAEGER_RECEIVER_ENDPOINT := "http://otel:14268/api/traces"
	// Create Jaeger exporter
	exporter, err := jaeger.New(

		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(JAEGER_RECEIVER_ENDPOINT)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Resource can be modelled after the service
	resource1 := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("v0.0.1"),
	)

	// Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource1),
		// Configure the sampler to always sample.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	// Set the propagator to tracecontext (the default is baggagetracecontext)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
