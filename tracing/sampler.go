// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/time/rate"
)

const (
	SAMPLER_TYPE      = "sampler.type"
	SAMPLER_PARAM     = "sampler.param"
	SAMPLER_PATH      = "sampler.path"
	SAMPLER_PARENT    = "sampler.parent"
	SAMPLER_ALGORITHM = "sampler.algorithm"
)

var _ sdktrace.Sampler = (*rateLimitSampler)(nil)

type rateLimitSampler struct {
	maxPerSecond float64
	limiter      *rate.Limiter
}

// NewRateLimitSampler returns a new rate limit sampler.
func NewRateLimitSampler(tc Config) *rateLimitSampler {
	return &rateLimitSampler{
		maxPerSecond: tc.PerSecond,
		limiter:      rate.NewLimiter(rate.Every(time.Duration(1.0/tc.PerSecond)*time.Second), 1),
	}
}

// ShouldSample determines if a trace should be sampled.
func (s *rateLimitSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	// apply sampler attributes
	attributes := []attribute.KeyValue{
		attribute.String(SAMPLER_ALGORITHM, "rate-limiting"),
		attribute.Float64(SAMPLER_PARAM, s.maxPerSecond),
	}

	// default to drop
	result := sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Attributes: attributes,
	}

	if s.limiter.Allow() {
		result.Decision = sdktrace.RecordAndSample
	}

	return result
}

// Description returns the description of the rate limit sampler.
func (s *rateLimitSampler) Description() string {
	return fmt.Sprintf("rate-limit-sampler{%v}", s.maxPerSecond)
}
