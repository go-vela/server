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
	&cli.Float64Flag{
		EnvVars: []string{"VELA_TRACING_SAMPLER_RATIO", "TRACING_SAMPLER_RATIO"},
		Name:    "tracing.sampler.ratio",
		Usage:   "set otel tracing sampler ratio. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   0.5,
	},
}
