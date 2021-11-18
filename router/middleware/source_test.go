// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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

func TestMiddleware_Source(t *testing.T) {
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
	engine.Use(Source(want))
	engine.GET("/health", func(c *gin.Context) {
		got = scm.FromContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Source returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Source is %v, want %v", got, want)
	}
}
