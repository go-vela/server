// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-vela/server/tracing"
)

func TestMiddleware_TracingClient(t *testing.T) {
	// setup types
	var got *tracing.Client
	want := &tracing.Client{
		Config: tracing.Config{
			EnableTracing: true,
		},
	}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(TracingClient(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("tracing").(*tracing.Client)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TracingClient is %v, want %v", got, want)
	}
}

func TestMiddleware_TracingInstrumentation(t *testing.T) {
	// setup types
	tt := []struct {
		tc     *tracing.Client
		assert func(trace.SpanContext) error
	}{
		{
			tc: &tracing.Client{
				Config: tracing.Config{
					EnableTracing: false,
					ServiceName:   "vela-test",
				},
				TracerProvider: sdktrace.NewTracerProvider(),
			},
			assert: func(got trace.SpanContext) error {
				if !reflect.DeepEqual(got, trace.SpanContext{}) {
					return errors.New("span context is not empty")
				}
				return nil
			},
		},
		{
			tc: &tracing.Client{
				Config: tracing.Config{
					EnableTracing: true,
					ServiceName:   "vela-test",
				},
				TracerProvider: sdktrace.NewTracerProvider(),
			},
			assert: func(got trace.SpanContext) error {
				if reflect.DeepEqual(got, trace.SpanContext{}) {
					return errors.New("span context is empty")
				}
				return nil
			},
		},
	}

	// setup context
	gin.SetMode(gin.TestMode)

	for _, test := range tt {
		got := trace.SpanContext{}
		resp := httptest.NewRecorder()
		context, engine := gin.CreateTestContext(resp)
		context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

		// setup mock server
		engine.Use(TracingInstrumentation(test.tc))
		engine.GET("/health", func(c *gin.Context) {
			got = trace.SpanContextFromContext(c.Request.Context())

			c.Status(http.StatusOK)
		})

		// run test
		engine.ServeHTTP(context.Writer, context.Request)

		if resp.Code != http.StatusOK {
			t.Errorf("TracingInstrumentation returned %v, want %v", resp.Code, http.StatusOK)
		}

		err := test.assert(got)
		if err != nil {
			t.Errorf("TracingInstrumentation test assertion failed: %s", err)
		}
	}
}
