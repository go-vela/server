// SPDX-License-Identifier: Apache-2.0

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
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/dashboard"
	"github.com/go-vela/server/router/middleware/hook"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/schedule"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/types/library"
)

func TestMiddleware_Logger(t *testing.T) {
	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.GetOwner().SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

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

	h := new(library.Hook)
	h.SetID(1)
	h.SetRepoID(1)
	h.SetBuildID(1)

	sc := new(api.Schedule)
	sc.SetID(1)
	sc.SetRepo(r)
	sc.SetName("foo")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")

	d := new(api.Dashboard)
	d.SetID("1")
	d.SetName("foo")

	w := new(api.Worker)
	w.SetID(1)
	w.SetHostname("worker_0")
	w.SetAddress("localhost")

	p := new(library.Pipeline)
	p.SetID(1)
	p.SetRepoID(1)

	payload, _ := json.Marshal(`{"foo": "bar"}`)
	wantLevel := logrus.InfoLevel

	logger, loggerHook := test.NewNullLogger()
	defer loggerHook.Reset()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/repos/foo", bytes.NewBuffer(payload))

	// setup mock server
	engine.Use(Logger(logger, time.RFC3339))
	engine.Use(func(c *gin.Context) { build.ToContext(c, b) })
	engine.Use(func(c *gin.Context) { dashboard.ToContext(c, d) })
	engine.Use(func(c *gin.Context) { hook.ToContext(c, h) })
	engine.Use(func(c *gin.Context) { pipeline.ToContext(c, p) })
	engine.Use(func(c *gin.Context) { repo.ToContext(c, r) })
	engine.Use(func(c *gin.Context) { schedule.ToContext(c, sc) })
	engine.Use(func(c *gin.Context) { service.ToContext(c, svc) })
	engine.Use(func(c *gin.Context) { step.ToContext(c, s) })
	engine.Use(func(c *gin.Context) { user.ToContext(c, u) })
	engine.Use(func(c *gin.Context) { worker.ToContext(c, w) })
	engine.Use(org.Establish())
	engine.Use(Payload())
	engine.POST("/repos/:org", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := loggerHook.LastEntry().Level
	gotMessage := loggerHook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if gotMessage == "" {
		t.Errorf("Logger Message is %v, want non-empty string", gotMessage)
	}

	if strings.Contains(gotMessage, "GET") {
		t.Errorf("Logger Message is %v, want message to contain GET", gotMessage)
	}

	if !strings.Contains(gotMessage, "POST") {
		t.Errorf("Logger Message is %v, message shouldn't contain POST", gotMessage)
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

	r := new(api.Repo)
	r.SetID(1)
	r.GetOwner().SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	logRepo, _ := json.Marshal(r)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
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

func TestMiddleware_Format(t *testing.T) {
	wantLabels := "labels.vela"

	// setup data, fields, and logger
	formatter := &ECSFormatter{
		DataKey: wantLabels,
	}

	fields := logrus.Fields{
		"ip":         "123.4.5.6",
		"id":         "deadbeef",
		"method":     http.MethodGet,
		"path":       "/foobar",
		"latency":    0,
		"status":     http.StatusOK,
		"user-agent": "foobar",
		"version":    "v1.0.0",
		"org":        "foo",
		"user":       "octocat",
	}

	logger := logrus.NewEntry(logrus.StandardLogger())
	entry := logger.WithFields(fields)

	got, err := formatter.Format(entry)

	// run test

	if err != nil {
		t.Errorf("Format returned err: %v", err)
	}

	if got == nil {
		t.Errorf("Format returned nothing, want a log")
	}

	if !strings.Contains(string(got), "GET") {
		t.Errorf("Format returned %v, want to contain GET", string(got))
	}

	if !strings.Contains(string(got), "/foobar") {
		t.Errorf("Format returned %v, want to contain /foobar", string(got))
	}
}
