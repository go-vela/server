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
		Value:   "go-otel-server",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_EXPORTER_SSL_CERT_PATH"},
		Name:    "tracing.exporter.cert_path",
		Usage:   "set the path to certs used for communicating with the otel exporter. if not set, will use insecure communication",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_ENDPOINT"},
		Name:    "tracing.exporter.endpoint",
		Usage:   "set the otel exporter endpoint",
		Value:   "127.0.0.1:4318",
	},

	// Resource Flags

	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ATTRIBUTES"},
		Name:    "tracing.resource.attributes",
		Usage:   "set otel resource attributes as a list of key=value pairs. each one will be attached to each span as an attribute",
		Value:   cli.NewStringSlice("process.runtime.name=go"),
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ENV_ATTRIBUTES"},
		Name:    "tracing.resource.env_attributes",
		Usage:   "set otel resource attributes as a list of key=env_variable_key pairs. each one will be attached to each span as an attribute where the value is retrieved from the environment using the pair value",
		Value:   cli.NewStringSlice("deployment.environment=CLOUD_ENVIRONMENT"),
	},

	// Sampler Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_TRACESTATE"},
		Name:    "tracing.sampler.tracestate",
		Usage:   "set otel sampler trace state attached to each span as an attribute.",
		Value:   "sampler.tracestate",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_PARENT"},
		Name:    "tracing.sampler.parent",
		Usage:   "set otel sampler parent attribute attached to each span.",
		Value:   "sampler.parent",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_TYPE"},
		Name:    "tracing.sampler.type",
		Usage:   "set otel sampler type attribute attached to each span.",
		Value:   "sampler.type",
	},
	&cli.Float64Flag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_RATIO"},
		Name:    "tracing.sampler.ratio",
		Usage:   "set otel tracing head-sampler acceptance ratio. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   0.5,
	},
	&cli.Float64Flag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_RATELIMIT_PER_SECOND"},
		Name:    "tracing.sampler.persecond",
		Usage:   "set otel tracing head-sampler rate-limiting to N per second. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   1,
	},
}
