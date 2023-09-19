// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package tracing

import (
	"github.com/urfave/cli/v2"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config represents the configurations for otel tracing.
type Config struct {
	EnableTracing  bool
	ServiceName    string
	TracerProvider *sdktrace.TracerProvider
}

// New takes cli context and returns a tracing config to supply to traceable services.
func New(c *cli.Context) (*Config, error) {
	enable := c.Bool("tracing.enable")
	serviceName := c.String("tracing.service.name")

	// could skip creating the tracer if tracing is disabled
	tracerProvider, err := initTracer(c)
	if err != nil {
		return nil, err
	}

	return &Config{
		EnableTracing:  enable,
		ServiceName:    serviceName,
		TracerProvider: tracerProvider,
	}, nil
}
