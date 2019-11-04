// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/source"
	"github.com/go-vela/server/source/github"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Source(t *testing.T) {
	// setup types
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()
	var got source.Service
	want, _ := github.NewTest(s.URL)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Source(want))
	engine.GET("/health", func(c *gin.Context) {
		got = source.FromContext(c)

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
