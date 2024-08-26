// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// initTracer returns the tracer provider supplied to the tracing config.
func initTracer(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	withTLS := otlptracehttp.WithInsecure()

	if len(cfg.CertPath) > 0 {
		pem, err := os.ReadFile(cfg.CertPath)
		if err != nil {
			return nil, err
		}

		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(pem)

		withTLS = otlptracehttp.WithTLSClientConfig(
			&tls.Config{
				RootCAs:    certs,
				MinVersion: tls.VersionTLS12,
			})
	} else {
		logrus.Warn("no otel cert path set, exporting traces insecurely")
	}

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cfg.ExporterURL),
		withTLS,
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
	}

	for k, v := range cfg.ResourceAttributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(attrs...),
	)
	if err != nil {
		return nil, err
	}

	rateLimitSampler := NewRateLimitSampler(cfg)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(rateLimitSampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
