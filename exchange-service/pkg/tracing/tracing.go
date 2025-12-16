package tracing

import (
	"context"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

// InitTracer initializes OpenTelemetry tracer with Zipkin exporter
func InitTracer(serviceName, zipkinEndpoint string, logger *zap.Logger) (io.Closer, error) {
	if zipkinEndpoint == "" {
		logger.Warn("Zipkin endpoint not configured, tracing disabled")
		return &noopCloser{}, nil
	}

	// Create Zipkin exporter
	exporter, err := zipkin.New(zipkinEndpoint)
	if err != nil {
		return nil, err
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator to tracecontext (W3C Trace Context)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	logger.Info("Tracing initialized",
		zap.String("service", serviceName),
		zap.String("zipkin", zipkinEndpoint),
	)

	return &tracerCloser{tp: tp}, nil
}

type tracerCloser struct {
	tp *sdktrace.TracerProvider
}

func (tc *tracerCloser) Close() error {
	return tc.tp.Shutdown(context.Background())
}

type noopCloser struct{}

func (nc *noopCloser) Close() error {
	return nil
}

