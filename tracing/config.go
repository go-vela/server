package tracing

import (
	"github.com/urfave/cli/v2"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TracingConfig represents the configurations for otel tracing
type TracingConfig struct {
	EnableTracing  bool
	ServiceName    string
	TracerProvider *sdktrace.TracerProvider
}

// New takes cli context and returns a tracing config to supply to traceable services
func New(c *cli.Context) (*TracingConfig, error) {
	enable := c.Bool("tracing.enable")
	serviceName := c.String("tracing.service.name")

	// could skip creating the tracer if tracing is disabled
	tracerProvider, err := initTracer(c)
	if err != nil {
		return nil, err
	}

	return &TracingConfig{
		EnableTracing:  enable,
		ServiceName:    serviceName,
		TracerProvider: tracerProvider,
	}, nil
}
