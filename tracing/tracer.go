// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

		tlsCfg := &tls.Config{
			RootCAs:    certs,
			MinVersion: tls.VersionTLS12,
		}

		// if a TLS minimum version is supplied, set that in the config
		if len(cfg.TLSMinVersion) > 0 {
			var tlsVersion uint16

			switch cfg.TLSMinVersion {
			case "1.0":
				tlsVersion = tls.VersionTLS10
			case "1.1":
				tlsVersion = tls.VersionTLS11
			case "1.2":
				tlsVersion = tls.VersionTLS12
			case "1.3":
				tlsVersion = tls.VersionTLS13
			default:
				return nil, fmt.Errorf("invalid TLS minimum version supplied: %s", cfg.TLSMinVersion)
			}

			tlsCfg.MinVersion = tlsVersion
		}

		withTLS = otlptracehttp.WithTLSClientConfig(tlsCfg)
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
