// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/gin-gonic/gin"
)

func TestCookie_SecureCookie(t *testing.T) {
	type args struct {
		secure bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "secure cookie disabled",
			args: args{
				secure: false,
			},
			want: false,
		},
		{
			name: "secure cookie enabled",
			args: args{
				secure: true,
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

			engine.Use(SecureCookie(tt.args.secure))
			engine.GET("/health", func(c *gin.Context) {
				got = c.Value("securecookie").(bool)

				c.Status(http.StatusOK)
			})

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			assert.Equal(t, tt.want, got)
		})
	}
}
