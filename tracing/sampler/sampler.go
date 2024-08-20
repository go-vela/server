package sampler

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/time/rate"
)

const (
	SAMPLER_TYPE      = "sampler.type"
	SAMPLER_PARAM     = "sampler.param"
	SAMPLER_PATH      = "sampler.path"
	SAMPLER_PARENT    = "sampler.parent"
	SAMPLER_ALGORITHM = "sampler.algorithm"
)

var _ trace.Sampler = (*rateLimitSampler)(nil)

type rateLimitSampler struct {
	maxPerSecond float64
	limiter      *rate.Limiter
}

func newRateLimitSampler(perSec float64) *rateLimitSampler {
	return &rateLimitSampler{
		maxPerSecond: perSec,
		limiter:      rate.NewLimiter(rate.Every(time.Duration(1.0/perSec)*time.Second), 1),
	}
}

func (s *rateLimitSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if s.limiter.Allow() {
		return trace.SamplingResult{
			Decision: trace.RecordAndSample,
			Attributes: []attribute.KeyValue{
				attribute.String(SAMPLER_ALGORITHM, "rate-limiting"),
				attribute.Float64(SAMPLER_PARAM, s.maxPerSecond),
			},
		}
	} else {
		return trace.SamplingResult{
			Decision: trace.Drop,
			Attributes: []attribute.KeyValue{
				attribute.String(SAMPLER_ALGORITHM, "rate-limiting"),
				attribute.Float64(SAMPLER_PARAM, s.maxPerSecond),
			},
		}
	}

}

func (s *rateLimitSampler) Description() string {
	return fmt.Sprintf("RateLimitSampler{%v}", s.maxPerSecond)
}
