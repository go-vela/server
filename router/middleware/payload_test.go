// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Payload(t *testing.T) {
	// setup types
	want := `{"foo": "bar"}`
	jsonBody, _ := json.Marshal(want)
	var got interface{}

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/health", bytes.NewBuffer(jsonBody))

	// setup mock server
	engine.Use(Payload())
	engine.POST("/health", func(c *gin.Context) {
		got = c.Value("payload")

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Payload returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Payload is %v, want %v", got, want)
	}
}
