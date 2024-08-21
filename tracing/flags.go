// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	// Tracing Flags

	&cli.BoolFlag{
		EnvVars: []string{"VELA_ENABLE_TRACING", "TRACING_ENABLE"},
		Name:    "tracing.enable",
		Usage:   "enable otel tracing",
		Value:   false,
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_TRACING_SERVICE_NAME", "TRACING_SERVICE_NAME"},
		Name:    "tracing.service.name",
		Usage:   "set otel tracing service name",
		Value:   "vela-server",
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_TRACING_SAMPLE_TAGS", "OTEL_TRACING_SAMPLE_TAGS"},
		Name:    "tracing.sample.tags",
		Usage:   "set otel trace sample state tags. see: https://opentelemetry.io/docs/concepts/sampling/",
	},
	&cli.Float64Flag{
		EnvVars: []string{"VELA_TRACING_SAMPLE_RATIO", "OTEL_TRACE_SAMPLE_RATIO"},
		Name:    "tracing.sample.ratio",
		Usage:   "set otel tracing head-sampler acceptance ratio. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   0.5,
	},
	&cli.Float64Flag{
		EnvVars: []string{"VELA_TRACING_SAMPLE_RATELIMIT_PER_SECOND", "OTEL_TRACING_SAMPLE_RATELIMIT_PER_SECOND"},
		Name:    "tracing.sample.persecond",
		Usage:   "set otel tracing head-sampler rate-limiting to N per second. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   1,
	},
}
