// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
)

func TestMiddleware_CompilerNative(t *testing.T) {
	// setup types
	defaultCloneImage := "target/vela-git-slim"
	wantCloneImage := "target/vela-git-slim:latest"

	c := new(cli.Command)

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: defaultCloneImage,
		},
	}

	want, _ := native.FromCLICommand(context.Background(), c)
	want.SetCloneImage(wantCloneImage)

	var got compiler.Engine

	got, _ = native.FromCLICommand(context.Background(), c)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			s := settings.Platform{
				Compiler: &settings.Compiler{},
			}
			s.SetCloneImage(wantCloneImage)

			sMiddleware.ToContext(c, &s)

			c.Next()
		}
	}(),
	)

	engine.Use(Compiler(got))

	context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)

	engine.GET("/health", func(c *gin.Context) {
		got = compiler.FromContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Compiler returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got.GetSettings(), want.GetSettings()) {
		t.Errorf("Compiler is %v, want %v", got, want)
	}
}
