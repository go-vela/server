// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/pkg/errors"
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

// FromCLIContext takes cli context and returns a tracing config to supply to traceable services.
func FromCLIContext(c *cli.Context) (*Client, error) {
	cfg := Config{
		EnableTracing:        c.Bool("tracing.enable"),
		ServiceName:          c.String("tracing.service.name"),
		ExporterURL:          c.String("tracing.exporter.endpoint"),
		CertPath:             c.String("tracing.exporter.cert_path"),
		TLSMinVersion:        c.String("tracing.exporter.tls-min-version"),
		ResourceAttributes:   map[string]string{},
		TraceStateAttributes: map[string]string{},
		SpanAttributes:       map[string]string{},
		Sampler: Sampler{
			PerSecond: c.Float64("tracing.sampler.persecond"),
			Tasks:     Tasks{},
		},
	}

	// read per-task configurations from file
	tasksConfigPath := c.String("tracing.sampler.tasks")
	if len(tasksConfigPath) > 0 {
		f, err := os.ReadFile(tasksConfigPath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to read tracing tasks config file from path %s", tasksConfigPath))
		}

		err = json.Unmarshal(f, &cfg.Sampler.Tasks)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to parse tracing tasks config file from path %s", tasksConfigPath))
		}
	}

	// identity func used to map a string back to itself
	identityFn := func(s string) string { return s }

	// span attributes
	cfg.SpanAttributes = keyValueSliceToMap(c.StringSlice("tracing.span.attributes"), identityFn)

	// tracestate attributes
	cfg.TraceStateAttributes = keyValueSliceToMap(c.StringSlice("tracing.tracestate.attributes"), identityFn)

	// merge static resource attributes with those fetched from the environment using os.Getenv
	cfg.ResourceAttributes = keyValueSliceToMap(c.StringSlice("tracing.resource.attributes"), identityFn)
	m := keyValueSliceToMap(c.StringSlice("tracing.resource.env_attributes"), os.Getenv)
	maps.Copy(cfg.ResourceAttributes, m)

	client := &Client{
		Config: cfg,
	}

	if cfg.EnableTracing {
		// initialize the tracer provider and assign it to the client
		tracer, err := initTracer(c.Context, cfg)
		if err != nil {
			return nil, err
		}

		client.TracerProvider = tracer
	}

	return client, nil
}

// keyValueSliceToMap converts a slice of key=value strings to a map of key to value using the supplied map function.
func keyValueSliceToMap(kv []string, fn func(string) string) map[string]string {
	m := map[string]string{}

	for _, attr := range kv {
		parts := strings.SplitN(attr, "=", 2)

		if len(parts) != 2 || len(parts[1]) == 0 {
			continue
		}

		k := parts[0]
		v := fn(parts[1])

		if len(v) == 0 {
			continue
		}

		m[k] = v
	}

	return m
}
