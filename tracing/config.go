// SPDX-License-Identifier: Apache-2.0

package tracing

import (
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
	EnableTracing bool
	ServiceName   string
	Sampler
}

// Sampler represents the configurations for the otel sampler.
// Used to determine if a trace should be sampled.
type Sampler struct {
	Tags      []string
	Ratio     float64
	PerSecond float64
}

// FromCLIContext takes cli context and returns a tracing config to supply to traceable services.
func FromCLIContext(c *cli.Context) (*Client, error) {
	tCfg := Config{
		EnableTracing: c.Bool("tracing.enable"),
		ServiceName:   c.String("tracing.service.name"),
		Sampler: Sampler{
			Tags:      c.StringSlice("tracing.sample.tags"),
			Ratio:     c.Float64("tracing.sample.ratio"),
			PerSecond: c.Float64("tracing.sample.persecond"),
		},
	}

	tracer, err := initTracer(c.Context, tCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:         tCfg,
		TracerProvider: tracer,
	}, nil
}
