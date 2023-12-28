// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/gin-gonic/gin"
)

func TestWebhook_WebhookValidation(t *testing.T) {
	type args struct {
		validate bool
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "validation disabled",
			args: args{
				validate: false,
			},
			want: false,
		},
		{
			name: "validation enabled",
			args: args{
				validate: true,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup context
			gin.SetMode(gin.TestMode)

			var got bool

			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

			engine.Use(WebhookValidation(tt.args.validate))
			engine.GET("/health", func(c *gin.Context) {
				got = c.Value("webhookvalidation").(bool)

				c.Status(http.StatusOK)
			})

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			assert.Equal(t, tt.want, got)
		})
	}
}
