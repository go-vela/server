// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	// Tracing Flags

	&cli.BoolFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_ENABLE"},
		Name:    "tracing.enable",
		Usage:   "enable otel tracing",
		Value:   false,
	},

	// Exporter Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SERVICE_NAME"},
		Name:    "tracing.service.name",
		Usage:   "set otel tracing service name",
		Value:   "vela-server",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_ENDPOINT"},
		Name:    "tracing.exporter.endpoint",
		Usage:   "set the otel exporter endpoint",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_EXPORTER_SSL_CERT_PATH"},
		Name:    "tracing.exporter.cert_path",
		Usage:   "set the path to certs used for communicating with the otel exporter. if not set, will use insecure communication",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_TLS_MIN_VERSION"},
		Name:    "tracing.exporter.tls-min-version",
		Usage:   "optional TLS minimum version requirement to set when communicating with the otel exporter",
		Value:   "1.2",
	},

	// Resource Flags

	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ATTRIBUTES"},
		Name:    "tracing.resource.attributes",
		Usage:   "set otel resource attributes as a list of key=value pairs. each one will be attached to each span as a resource attribute",
		Value:   cli.NewStringSlice("process.runtime.name=go"),
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ENV_ATTRIBUTES"},
		Name:    "tracing.resource.env_attributes",
		Usage:   "set otel resource attributes as a list of key=env_variable_key pairs. each one will be attached to each span as an attribute where the value is retrieved from the environment using the pair value",
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SPAN_ATTRIBUTES"},
		Name:    "tracing.span.attributes",
		Usage:   "set otel span attributes as a list of key=value pairs. each one will be attached to each span as a sampler attribute",
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_TRACESTATE_ATTRIBUTES"},
		Name:    "tracing.tracestate.attributes",
		Usage:   "set otel tracestate attributes as a list of key=value pairs. each one will be inserted into the tracestate for each sampled span",
	},

	// Sampler Flags

	&cli.Float64Flag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_RATELIMIT_PER_SECOND"},
		Name:    "tracing.sampler.persecond",
		Usage:   "set otel tracing head-sampler rate-limiting to N per second. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   0.2,
	},
}
