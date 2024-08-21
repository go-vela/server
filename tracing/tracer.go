// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// initTracer returns the tracer provider supplied to the tracing config.
func initTracer(ctx context.Context, tCfg Config) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(tCfg.ServiceName),
	))
	if err != nil {
		return nil, err
	}

	ratioSampler := sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(tCfg.Sampler.Ratio),
	)
	rateLimitSampler := sdktrace.ParentBased(
		NewRateLimitSampler(tCfg),
	)

	// todo: apply the tags to the sampler

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(ratioSampler),
		sdktrace.WithSampler(rateLimitSampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
