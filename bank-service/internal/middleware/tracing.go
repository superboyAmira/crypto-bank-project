package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracing returns a Fiber middleware that creates spans for HTTP requests
func Tracing(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	// Use TraceContext propagator explicitly
	propagator := propagation.TraceContext{}
	otel.SetTextMapPropagator(propagator)

	return func(c *fiber.Ctx) error {
		// Extract trace context from headers
		ctx := propagator.Extract(c.Context(), &fiberCarrier{c: c})

		// Start span
		spanName := c.Method() + " " + c.Path()
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.route", c.Route().Path),
				attribute.String("http.url", string(c.Request().URI().FullURI())),
				attribute.String("http.scheme", c.Protocol()),
				attribute.String("http.target", c.Path()),
				attribute.String("http.client_ip", c.IP()),
				attribute.String("http.user_agent", c.Get("User-Agent")),
			),
		)
		defer span.End()

		// Store context with span
		c.SetUserContext(ctx)

		// Execute request
		err := c.Next()

		// Record status code
		statusCode := c.Response().StatusCode()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
			if err != nil {
				span.RecordError(err)
			}
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}

// fiberCarrier adapts Fiber's context to propagation.TextMapCarrier
type fiberCarrier struct {
	c *fiber.Ctx
}

func (fc *fiberCarrier) Get(key string) string {
	return fc.c.Get(key)
}

func (fc *fiberCarrier) Set(key, val string) {
	fc.c.Set(key, val)
}

func (fc *fiberCarrier) Keys() []string {
	keys := make([]string, 0)
	fc.c.Request().Header.VisitAll(func(key, _ []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
