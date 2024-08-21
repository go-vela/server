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
	SamplerType      = "sampler.type"
	SamplerParam     = "sampler.param"
	SamplerPath      = "sampler.path"
	SamplerParent    = "sampler.parent"
	SamplerAlgorithm = "sampler.algorithm"
)

var _ sdktrace.Sampler = (*RateLimitSampler)(nil)

type RateLimitSampler struct {
	maxPerSecond float64
	limiter      *rate.Limiter
}

// NewRateLimitSampler returns a new rate limit sampler.
func NewRateLimitSampler(tc Config) *RateLimitSampler {
	return &RateLimitSampler{
		maxPerSecond: tc.PerSecond,
		limiter:      rate.NewLimiter(rate.Every(time.Duration(1.0/tc.PerSecond)*time.Second), 1),
	}
}

// ShouldSample determines if a trace should be sampled.
func (s *RateLimitSampler) ShouldSample(_ sdktrace.SamplingParameters) sdktrace.SamplingResult {
	// apply sampler attributes
	attributes := []attribute.KeyValue{
		attribute.String(SamplerAlgorithm, "rate-limiting"),
		attribute.Float64(SamplerParam, s.maxPerSecond),
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
func (s *RateLimitSampler) Description() string {
	return fmt.Sprintf("rate-limit-sampler{%v}", s.maxPerSecond)
}
