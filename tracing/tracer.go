package tracing

import (
	"context"

	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"go.opentelemetry.io/otel"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer returns the tracer provider supplied to the tracing config
func initTracer(c *cli.Context) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()

	// TODO: inject actual context
	exporter, err := otlp.New(context.TODO(), client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(context.TODO(), resource.WithAttributes(
		semconv.ServiceName(c.String("tracing.service.name")),
	))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
