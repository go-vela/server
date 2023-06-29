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

type TracingConfig struct {
	EnableTracing  bool
	TracerProvider *sdktrace.TracerProvider
}

func DefaultTracingConfig() *TracingConfig {
	return &TracingConfig{
		EnableTracing: false,
	}
}

func SetupTracing(c *cli.Context) (*TracingConfig, error) {
	// initialize the shared tracing provider
	tp, err := setupTracerProvider()
	if err != nil {
		return nil, err
	}

	return &TracingConfig{
		EnableTracing:  c.Bool("enable-tracing"),
		TracerProvider: tp,
	}, nil
}

func setupTracerProvider() (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()

	exporter, err := otlp.New(context.TODO(), client)
	if err != nil {
		return nil, err
	}

	// TODO: make the serviceName configurable
	res, err := resource.New(context.TODO(), resource.WithAttributes(
		semconv.ServiceName("vela-server"),
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
