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
		Usage:   "enable otel tracing. see: https://opentelemetry.io/docs/concepts/signals/traces/",
		Value:   false,
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SERVICE_NAME"},
		Name:    "tracing.service.name",
		Usage:   "set otel tracing service name. see: https://opentelemetry.io/docs/languages/sdk-configuration/general/",
		Value:   "vela-server",
	},

	// Exporter Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_EXPORTER_OTLP_ENDPOINT"},
		Name:    "tracing.exporter.endpoint",
		Usage:   "set the otel exporter endpoint (ex. scheme://host:port). see: https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_EXPORTER_SSL_CERT_PATH"},
		Name:    "tracing.exporter.cert_path",
		Usage:   "set the path to certs used for communicating with the otel exporter. if not set, will use insecure communication. see: https://opentelemetry.io/docs/specs/otel/protocol/exporter/",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_TLS_MIN_VERSION"},
		Name:    "tracing.exporter.tls-min-version",
		Usage:   "optional TLS minimum version requirement to set when communicating with the otel exporter",
		Value:   "1.2",
	},

	// Attribute Flags

	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ATTRIBUTES"},
		Name:    "tracing.resource.attributes",
		Usage:   "set otel resource (span) attributes as a list of key=value pairs. each one will be attached to each span as a 'process' attribute. see: https://opentelemetry.io/docs/languages/go/instrumentation/#span-attributes",
		Value:   cli.NewStringSlice("process.runtime.name=go"),
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_RESOURCE_ENV_ATTRIBUTES"},
		Name:    "tracing.resource.env_attributes",
		Usage:   "set otel resource (span) attributes as a list of key=env_variable_key pairs. each one will be attached to each span as a 'process' attribute where the value is retrieved from the environment using the pair value. see: https://opentelemetry.io/docs/languages/go/instrumentation/#span-attributes",
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_SPAN_ATTRIBUTES"},
		Name:    "tracing.span.attributes",
		Usage:   "set otel span attributes as a list of key=value pairs. each one will be attached to each span as a 'tag' attribute. see: https://opentelemetry.io/docs/languages/go/instrumentation/#span-attributes",
	},
	&cli.StringSliceFlag{
		EnvVars: []string{"VELA_OTEL_TRACING_TRACESTATE_ATTRIBUTES"},
		Name:    "tracing.tracestate.attributes",
		Usage:   "set otel tracestate (span) attributes as a list of key=value pairs. each one will be inserted into the tracestate for each sampled span. see: https://www.w3.org/TR/trace-context",
	},

	// Sampler Flags

	&cli.Float64Flag{
		EnvVars: []string{"VELA_OTEL_TRACING_SAMPLER_RATELIMIT_PER_SECOND"},
		Name:    "tracing.sampler.persecond",
		Usage:   "set otel tracing head-sampler rate-limiting to N per second. see: https://opentelemetry.io/docs/concepts/sampling/",
		Value:   100,
	},
}
