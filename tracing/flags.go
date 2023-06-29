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
}
