// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Client represents the tracing client and the configurations that were used to initialize it.
type Client struct {
	Config
	TracerProvider *sdktrace.TracerProvider
}

// Config represents the configurations for otel tracing.
type Config struct {
	EnableTracing      bool
	ServiceName        string
	ExporterURL        string
	CertPath           string
	ResourceAttributes map[string]string
	Sampler
}

// Sampler represents the configurations for the otel sampler.
// Used to determine if a trace should be sampled.
type Sampler struct {
	TraceStateAttributes map[string]string
	SpanAttributes       map[string]string
	Ratio                float64
	PerSecond            float64
}

// FromCLIContext takes cli context and returns a tracing config to supply to traceable services.
func FromCLIContext(c *cli.Context) (*Client, error) {
	cfg := Config{
		EnableTracing:      c.Bool("tracing.enable"),
		ServiceName:        c.String("tracing.service.name"),
		ExporterURL:        c.String("tracing.exporter.endpoint"),
		CertPath:           c.String("tracing.exporter.cert_path"),
		ResourceAttributes: map[string]string{},
		Sampler: Sampler{
			TraceStateAttributes: map[string]string{
				"sampler": c.String("tracing.sampler.tracestate"),
			},
			SpanAttributes: map[string]string{
				"w3c.tracestate": fmt.Sprintf("sampler=%s", c.String("tracing.sampler.tracestate")),
				"sampler.parent": c.String("tracing.sampler.parent"),
				"sampler.type":   c.String("tracing.sampler.type"),
			},
			Ratio:     c.Float64("tracing.sampler.ratio"),
			PerSecond: c.Float64("tracing.sampler.persecond"),
		},
	}

	// add resource attributes
	for _, attr := range c.StringSlice("tracing.resource.attributes") {
		kv := strings.Split(attr, "=")
		if len(kv) != 2 {
			continue
		}

		cfg.ResourceAttributes[kv[0]] = kv[1]
	}

	// add resource attributes from environment
	for _, attr := range c.StringSlice("tracing.resource.env_attributes") {
		kv := strings.Split(attr, "=")
		if len(kv) != 2 {
			continue
		}

		v, found := os.LookupEnv(kv[1])
		if found {
			cfg.ResourceAttributes[kv[0]] = v
		}
	}

	tracer, err := initTracer(c.Context, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:         cfg,
		TracerProvider: tracer,
	}, nil
}
