// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/scm/github"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Scm(t *testing.T) {
	// setup types
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	var got scm.Service

	want, _ := github.NewTest(s.URL)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Scm(want))
	engine.GET("/health", func(c *gin.Context) {
		got = scm.FromContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Scm returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Scm is %v, want %v", got, want)
	}
}
