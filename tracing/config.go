// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Client represents the tracing client and the configurations that were used to initialize it.
type Client struct {
	Config
	TracerProvider *sdktrace.TracerProvider
}

// Config represents the configurations for otel tracing.
type Config struct {
	EnableTracing        bool
	ServiceName          string
	ExporterURL          string
	CertPath             string
	TLSMinVersion        string
	ResourceAttributes   map[string]string
	TraceStateAttributes map[string]string
	SpanAttributes       map[string]string
	Sampler
}

// Sampler represents the configurations for the otel sampler.
// Used to determine if a trace should be sampled.
type Sampler struct {
	PerSecond float64
	Tasks
}

// Tasks represents a map of task names to per-task configurations.
// A 'task name' is the endpoint or instrumentation scope, depending on the task.
// For example, database trace tasks could be 'gorm.query' and HTTP requests could be 'api/v1/:worker' depending on the endpoint.
type Tasks map[string]Task

// Task represents the sampler configurations on a per-task basis.
// 'Active' will disable/enable the task. If tracing encounters a task name not present in the map, it is considered Active (true).
type Task struct {
	Active bool
}

// FromCLICommand takes cli context and returns a tracing config to supply to traceable services.
func FromCLICommand(ctx context.Context, c *cli.Command) (*Client, error) {
	cfg := Config{
		EnableTracing:        c.Bool("tracing.enable"),
		ServiceName:          c.String("tracing.service.name"),
		ExporterURL:          c.String("tracing.exporter.endpoint"),
		CertPath:             c.String("tracing.exporter.cert_path"),
		TLSMinVersion:        c.String("tracing.exporter.tls-min-version"),
		ResourceAttributes:   c.StringMap("tracing.resource.attributes"),
		TraceStateAttributes: c.StringMap("tracing.tracestate.attributes"),
		SpanAttributes:       c.StringMap("tracing.span.attributes"),
		Sampler: Sampler{
			PerSecond: c.Float("tracing.sampler.persecond"),
			Tasks:     Tasks{},
		},
	}

	// read per-task configurations from file
	tasksConfigPath := c.String("tracing.sampler.tasks")
	if len(tasksConfigPath) > 0 {
		f, err := os.ReadFile(tasksConfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read tracing tasks config file from path %s: %w", tasksConfigPath, err)
		}

		err = json.Unmarshal(f, &cfg.Sampler.Tasks)
		if err != nil {
			return nil, fmt.Errorf("unable to parse tracing tasks config file from path %s: %w", tasksConfigPath, err)
		}
	}

	m := c.StringMap("tracing.resource.env_attributes")
	for k, v := range m {
		if len(v) > 0 {
			envVar := os.Getenv(v)
			if len(envVar) > 0 {
				cfg.ResourceAttributes[k] = envVar
			}
		}
	}

	client := &Client{
		Config: cfg,
	}

	if cfg.EnableTracing {
		// initialize the tracer provider and assign it to the client
		tracer, err := initTracer(ctx, cfg)
		if err != nil {
			return nil, err
		}

		client.TracerProvider = tracer
	}

	return client, nil
}
