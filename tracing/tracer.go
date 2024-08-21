// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/go-vela/server/tracing/sampler"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

// initTracer returns the tracer provider supplied to the tracing config.
func initTracer(c *cli.Context) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()
	ctx := context.Background()

	// TODO: inject actual context
	exporter, err := otlp.New(ctx, client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(c.String("tracing.service.name")),
	))
	if err != nil {
		return nil, err
	}

	logrus.Infof("intializing tracing using samplers: ratio(%f), persecond(%f)", c.Float64("tracing.sampler.ratio"), c.Float64("tracing.sampler.persecond"))

	ratioSampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Float64("tracing.sampler.ratio")))
	rateLimitSampler := sdktrace.ParentBased(sampler.NewRateLimitSampler(c.Float64("tracing.sampler.persecond")))

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
