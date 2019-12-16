// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Metadata(t *testing.T) {
	// setup types
	got := new(types.Metadata)
	want := &types.Metadata{}

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Metadata(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("metadata").(*types.Metadata)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Metadata returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Metadata is %v, want %v", got, want)
	}
}
