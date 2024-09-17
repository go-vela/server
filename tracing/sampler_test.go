// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestTracing_Sampler_ShouldSampleTask(t *testing.T) {
	// setup tests
	tests := []struct {
		sampler       RateLimitSampler
		tasks         Tasks
		samplerParams sdktrace.SamplingParameters
		want          bool
	}{
		// no tasks
		{
			sampler: RateLimitSampler{
				Config: Config{
					Sampler: Sampler{},
				},
			},
			samplerParams: sdktrace.SamplingParameters{
				Name: "/health",
			},
			want: true,
		},
		// task is active
		{
			sampler: RateLimitSampler{
				Config: Config{
					Sampler: Sampler{
						Tasks: Tasks{
							"/health": {
								Active: true,
							},
						}},
				},
			},
			samplerParams: sdktrace.SamplingParameters{
				Name: "/health",
			},
			want: true,
		},
		// task is inactive
		{
			sampler: RateLimitSampler{
				Config: Config{
					Sampler: Sampler{
						Tasks: Tasks{
							"/health": {
								Active: false,
							},
						}},
				},
			},
			samplerParams: sdktrace.SamplingParameters{
				Name: "/health",
			},
			want: false,
		},
		// task is non-endpoint
		{
			sampler: RateLimitSampler{
				Config: Config{
					Sampler: Sampler{
						Tasks: Tasks{
							"gorm.query": {
								Active: false,
							},
						}},
				},
			},
			samplerParams: sdktrace.SamplingParameters{
				Name: "Gorm.Query",
			},
			want: false,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.sampler.ShouldSampleTask(test.samplerParams)

		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("ShouldSampleTask mismatch (-want +got):\n%s", diff)
		}
	}
}
