package obs

import (
	"context"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	DefaultServiceTags = map[string]string{
		"service": "broker-service",
		"app":     "example",
		"env":     "development",
	}
)

func InitTracer(serviceName string) io.Closer {
	agentEndpoint := os.Getenv("JAEGER_AGENT_HOST")

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

func LogErrorWithSpan(logger *logrus.Logger, span opentracing.Span, context context.Context, msg ...interface{}) {
	logger.WithContext(context).
		WithFields(logrus.Fields{
			"spanID":  span.Context().(jaeger.SpanContext).SpanID().String(),
			"traceID": span.Context().(jaeger.SpanContext).TraceID().String(),
		}).
		Error(msg...)
}

func LogInfoWithSpan(logger *logrus.Logger, span opentracing.Span, context context.Context, msg ...interface{}) {
	logger.WithContext(context).
		WithFields(logrus.Fields{
			"spanID":  span.Context().(jaeger.SpanContext).SpanID().String(),
			"traceID": span.Context().(jaeger.SpanContext).TraceID().String(),
		}).
		Info(msg...)
}

func SetSpanTags(span opentracing.Span) {
	for key, val := range DefaultServiceTags {
		span.SetTag(key, val)
	}
}
