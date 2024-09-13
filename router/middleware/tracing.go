// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	tracingMiddleware "github.com/go-vela/server/router/middleware/tracing"
	"github.com/go-vela/server/tracing"
)

// TracingClient is a middleware function that attaches the tracing config
// to the context of every http.Request.
func TracingClient(tc *tracing.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracingMiddleware.ToContext(c, tc)

		c.Next()
	}
}

// TracingInstrumentation is a middleware function that attaches the tracing config
// to the context of every http.Request.
func TracingInstrumentation(tc *tracing.Client) gin.HandlerFunc {
	if tc.EnableTracing {
		return otelgin.Middleware(tc.ServiceName, otelgin.WithTracerProvider(tc.TracerProvider))
	}

	return func(_ *gin.Context) {}
}
