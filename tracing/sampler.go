// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

// RateLimitSampler is a sampler that uses time-based rate limiting.
type RateLimitSampler struct {
	Config
	limiter *rate.Limiter
}

// NewRateLimitSampler returns a new rate limit sampler.
func NewRateLimitSampler(tc Config) *RateLimitSampler {
	return &RateLimitSampler{
		Config:  tc,
		limiter: rate.NewLimiter(rate.Every(time.Duration(1.0/tc.PerSecond)*time.Second), 1),
	}
}

// ShouldSample determines if a trace should be sampled.
func (s *RateLimitSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	ts := psc.TraceState()
	for k, v := range s.Config.TraceStateAttributes {
		ts, _ = ts.Insert(k, v)
	}

	attributes := []attribute.KeyValue{
		attribute.String("sampler.algorithm", "rate-limiting"),
		attribute.Float64("sampler.param", s.Config.PerSecond),
	}

	for k, v := range s.Config.SpanAttributes {
		attributes = append(attributes, attribute.String(k, v))
	}

	result := sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Attributes: attributes,
		Tracestate: ts,
	}

	if s.limiter.Allow() {
		result.Decision = sdktrace.RecordAndSample
	}

	return result
}

// Description returns the description of the rate limit sampler.
func (s *RateLimitSampler) Description() string {
	return fmt.Sprintf("rate-limit-sampler{%v}", s.Config.PerSecond)
}
