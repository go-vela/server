// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestMiddleware_Logger(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	svc := new(library.Service)
	svc.SetID(1)
	svc.SetRepoID(1)
	svc.SetBuildID(1)
	svc.SetNumber(1)
	svc.SetName("foo")

	s := new(library.Step)
	s.SetID(1)
	s.SetRepoID(1)
	s.SetBuildID(1)
	s.SetNumber(1)
	s.SetName("foo")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")

	w := new(library.Worker)
	w.SetID(1)
	w.SetHostname("worker_0")
	w.SetAddress("localhost")

	payload, _ := json.Marshal(`{"foo": "bar"}`)
	wantLevel := logrus.InfoLevel
	// need to set latency
	// to consistent message such as:
	wantMessage := "200 0µs  POST /foobar"

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/foobar", bytes.NewBuffer(payload))

	// setup mock server
	engine.Use(func(c *gin.Context) { build.ToContext(c, b) })
	engine.Use(func(c *gin.Context) { repo.ToContext(c, r) })
	engine.Use(func(c *gin.Context) { service.ToContext(c, svc) })
	engine.Use(func(c *gin.Context) { step.ToContext(c, s) })
	engine.Use(func(c *gin.Context) { user.ToContext(c, u) })
	engine.Use(func(c *gin.Context) { worker.ToContext(c, w) })
	engine.Use(Payload())
	engine.Use(Logger(logger, time.RFC3339))
	engine.POST("/foobar", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}

func TestMiddleware_Logger_Error(t *testing.T) {
	// setup types
	wantLevel := logrus.ErrorLevel
	wantMessage := "Error #01: test error\n"

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foobar", nil)

	// setup mock server
	engine.Use(Logger(logger, time.RFC3339))
	engine.GET("/foobar", func(c *gin.Context) {
		//nolint:errcheck // ignore checking error
		c.Error(fmt.Errorf("test error"))
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}

func TestMiddleware_Logger_Sanitize(t *testing.T) {
	var logBody, logWant map[string]interface{}

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	logRepo, _ := json.Marshal(r)

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEmail("octocat@github.com")
	logBuild, _ := json.Marshal(b)

	sanitizeBuild := *b
	sanitizeBuild.SetEmail("[secure]")
	logSBuild, _ := json.Marshal(&sanitizeBuild)

	tests := []struct {
		dataType string
		body     []byte
		want     []byte
	}{
		{
			dataType: "stringMap",
			body:     logRepo,
			want:     logRepo,
		},
		{
			dataType: "stringMap",
			body:     logBuild,
			want:     logSBuild,
		},
		{
			dataType: "string",
			body:     []byte("successfully updated step"),
			want:     []byte("successfully updated step"),
		},
	}

	for _, test := range tests {
		if strings.EqualFold(test.dataType, "stringMap") {
			err := json.Unmarshal(test.body, &logBody)
			if err != nil {
				t.Errorf("unable to unmarshal log body data")
			}

			err = json.Unmarshal(test.want, &logWant)
			if err != nil {
				t.Errorf("unable to unmarshal log want data")
			}
		}

		got := sanitize(logBody)

		if !reflect.DeepEqual(got, logWant) {
			t.Errorf("Logger returned %v, want %v", got, logWant)
		}
	}
}
